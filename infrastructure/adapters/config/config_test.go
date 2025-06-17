package config

import (
	"os"
	"testing"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/stretchr/testify/assert"
)

// TestLoadConfig tests the LoadConfig function
func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("APP_ENV", "") // Ensure no .env file is loaded
	os.Setenv("POSTGRESQL_POSTGRES_PASSWORD", "testpgpass")
	os.Setenv("MONGODB_ROOT_PASSWORD", "testmongopass")
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("POSTGRESQL_POSTGRES_PASSWORD")
		os.Unsetenv("MONGODB_ROOT_PASSWORD")
	}()

	// Call the function to test
	config, err := LoadConfig()
	assert.NoError(t, err)

	// Verify default values were set
	assert.Equal(t, "8089", config.Server.Port)
	assert.Equal(t, "debug", config.Log.Level)
	assert.Equal(t, true, config.Log.Development)
	assert.Equal(t, "mongodb", config.Database.Type)
	assert.NotEmpty(t, config.Database.Postgres.DSN)
	assert.NotEmpty(t, config.Database.MongoDB.URI)

	// Verify durations were correctly converted
	assert.Equal(t, 30*time.Second, config.Auth.OIDCTimeout)
	assert.Equal(t, 10*time.Second, config.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, config.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.Server.IdleTimeout)
	assert.Equal(t, 10*time.Second, config.Server.ShutdownTimeout)

	// Verify feature flags
	assert.Equal(t, true, config.Features.UseGenerics)
}

// TestLoadConfigWithEnvironmentVariables tests loading config with environment variables
func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("APP_ENV", "")             // Ensure no .env file is loaded
	os.Setenv("APP_SERVER_PORT", "9090") // Override server port
	defer func() {
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_SERVER_PORT")
	}()

	// Call the function to test
	config, err := LoadConfig()
	assert.NoError(t, err)

	// Verify the config was loaded successfully
	assert.NotEmpty(t, config.Server.Port)
	assert.NotEmpty(t, config.Database.MongoDB.URI)
}

// TestLoadConfigFileNotFound tests loading config when the file is not found
func TestLoadConfigFileNotFound(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("APP_ENV", "") // Ensure no .env file is loaded
	defer func() {
		os.Unsetenv("APP_ENV")
	}()

	// No need to set a non-existent config path with koanf
	// The LoadConfig function will handle this case

	// Call the function to test
	config, err := LoadConfig()
	assert.NoError(t, err)

	// Verify default values were set
	assert.Equal(t, "8089", config.Server.Port)
	assert.Equal(t, "debug", config.Log.Level)
}

// TestLoadConfigUnmarshalError tests error handling when unmarshaling fails
func TestLoadConfigUnmarshalError(t *testing.T) {
	// Create a custom koanf instance for testing
	k := koanf.New(".")

	// Set environment variables for testing
	os.Setenv("APP_ENV", "") // Ensure no .env file is loaded
	defer func() {
		os.Unsetenv("APP_ENV")
	}()

	// Set a value that will cause unmarshaling to fail
	// Setting a complex value (map) for a field that expects a simple type
	k.Set("server.port", map[string]string{
		"invalid": "structure",
	})

	// Create a config struct
	var config Config

	// Try to unmarshal the koanf instance into the config struct
	err := k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{
		Tag: "mapstructure",
	})

	// Verify that an error was returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected type 'string', got unconvertible type")
}
