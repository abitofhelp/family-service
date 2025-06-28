// Copyright (c) 2025 A Bit of Help, Inc.

package profiling

import (
	"context"
	"fmt"
	"net/http"
	httppprof "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// Profiler provides functionality for profiling the application
type Profiler struct {
	logger         *logging.ContextLogger
	cpuProfilePath string
	memProfilePath string
	server         *http.Server
}

// NewProfiler creates a new Profiler
func NewProfiler(logger *logging.ContextLogger, cpuProfilePath, memProfilePath string) *Profiler {
	if logger == nil {
		panic("logger cannot be nil")
	}

	if cpuProfilePath == "" {
		cpuProfilePath = "cpu.pprof"
	}

	if memProfilePath == "" {
		memProfilePath = "mem.pprof"
	}

	return &Profiler{
		logger:         logger,
		cpuProfilePath: cpuProfilePath,
		memProfilePath: memProfilePath,
	}
}

// StartCPUProfiling starts CPU profiling
func (p *Profiler) StartCPUProfiling(ctx context.Context) error {
	p.logger.Info(ctx, "Starting CPU profiling", zap.String("path", p.cpuProfilePath))

	f, err := os.Create(p.cpuProfilePath)
	if err != nil {
		p.logger.Error(ctx, "Failed to create CPU profile file", zap.Error(err))
		return fmt.Errorf("failed to create CPU profile file: %w", err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		p.logger.Error(ctx, "Failed to start CPU profile", zap.Error(err))
		return fmt.Errorf("failed to start CPU profile: %w", err)
	}

	p.logger.Info(ctx, "CPU profiling started", zap.String("path", p.cpuProfilePath))
	return nil
}

// StopCPUProfiling stops CPU profiling
func (p *Profiler) StopCPUProfiling(ctx context.Context) {
	p.logger.Info(ctx, "Stopping CPU profiling")
	pprof.StopCPUProfile()
	p.logger.Info(ctx, "CPU profiling stopped", zap.String("path", p.cpuProfilePath))
}

// CaptureMemoryProfile captures a memory profile
func (p *Profiler) CaptureMemoryProfile(ctx context.Context) error {
	p.logger.Info(ctx, "Capturing memory profile", zap.String("path", p.memProfilePath))

	// Run garbage collection to get more accurate memory profile
	runtime.GC()

	f, err := os.Create(p.memProfilePath)
	if err != nil {
		p.logger.Error(ctx, "Failed to create memory profile file", zap.Error(err))
		return fmt.Errorf("failed to create memory profile file: %w", err)
	}
	defer f.Close()

	if err := pprof.WriteHeapProfile(f); err != nil {
		p.logger.Error(ctx, "Failed to write memory profile", zap.Error(err))
		return fmt.Errorf("failed to write memory profile: %w", err)
	}

	p.logger.Info(ctx, "Memory profile captured", zap.String("path", p.memProfilePath))
	return nil
}

// StartPprofServer starts an HTTP server for pprof
func (p *Profiler) StartPprofServer(ctx context.Context, addr string) error {
	if addr == "" {
		addr = "localhost:6060"
	}

	p.logger.Info(ctx, "Starting pprof HTTP server", zap.String("addr", addr))

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", httppprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", httppprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", httppprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", httppprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", httppprof.Trace)

	p.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			p.logger.Error(ctx, "Pprof HTTP server error", zap.Error(err))
		}
	}()

	p.logger.Info(ctx, "Pprof HTTP server started", 
		zap.String("addr", addr),
		zap.String("url", fmt.Sprintf("http://%s/debug/pprof/", addr)))
	return nil
}

// StopPprofServer stops the pprof HTTP server
func (p *Profiler) StopPprofServer(ctx context.Context) error {
	if p.server == nil {
		return nil
	}

	p.logger.Info(ctx, "Stopping pprof HTTP server")

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := p.server.Shutdown(ctxWithTimeout); err != nil {
		p.logger.Error(ctx, "Failed to stop pprof HTTP server", zap.Error(err))
		return fmt.Errorf("failed to stop pprof HTTP server: %w", err)
	}

	p.logger.Info(ctx, "Pprof HTTP server stopped")
	return nil
}

// ProfileOperation profiles a specific operation
func (p *Profiler) ProfileOperation(ctx context.Context, name string, operation func() error) error {
	p.logger.Info(ctx, "Profiling operation", zap.String("name", name))

	// Start CPU profiling with a unique name for this operation
	cpuProfilePath := fmt.Sprintf("%s_%s", name, p.cpuProfilePath)
	f, err := os.Create(cpuProfilePath)
	if err != nil {
		p.logger.Error(ctx, "Failed to create CPU profile file for operation", 
			zap.Error(err), 
			zap.String("name", name))
		return fmt.Errorf("failed to create CPU profile file for operation %s: %w", name, err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		p.logger.Error(ctx, "Failed to start CPU profile for operation", 
			zap.Error(err), 
			zap.String("name", name))
		return fmt.Errorf("failed to start CPU profile for operation %s: %w", name, err)
	}

	// Execute the operation
	opErr := operation()

	// Stop CPU profiling
	pprof.StopCPUProfile()
	f.Close()

	// Capture memory profile with a unique name for this operation
	memProfilePath := fmt.Sprintf("%s_%s", name, p.memProfilePath)
	runtime.GC()
	mf, err := os.Create(memProfilePath)
	if err != nil {
		p.logger.Error(ctx, "Failed to create memory profile file for operation", 
			zap.Error(err), 
			zap.String("name", name))
	} else {
		defer mf.Close()
		if err := pprof.WriteHeapProfile(mf); err != nil {
			p.logger.Error(ctx, "Failed to write memory profile for operation", 
				zap.Error(err), 
				zap.String("name", name))
		}
	}

	p.logger.Info(ctx, "Operation profiling completed", 
		zap.String("name", name),
		zap.String("cpu_profile", cpuProfilePath),
		zap.String("mem_profile", memProfilePath))

	return opErr
}
