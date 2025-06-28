// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/family-service/infrastructure/adapters/sqlite"
	"github.com/abitofhelp/servicelib/logging"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"go.uber.org/zap"
)

const (
	defaultSQLiteURI = "file:data/dev/sqlite/family_service.db?cache=shared&mode=rwc"
	defaultTimeout   = 30 * time.Second
)

func main() {
	// Get SQLite URI from environment variable or use default
	sqliteURI := getEnv("SQLITE_URI", defaultSQLiteURI)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Connect to SQLite
	db, err := sql.Open("sqlite3", sqliteURI)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer db.Close()

	// Set connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Ping the database to verify connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping SQLite: %v", err)
	}

	fmt.Println("Connected to SQLite database")

	// Create logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer zapLogger.Sync()

	// Create context logger
	contextLogger := logging.NewContextLogger(zapLogger)

	// Create repository
	repo := sqlite.NewSQLiteFamilyRepository(db, contextLogger)

	// Create and save sample families
	if err := createSampleFamilies(ctx, repo); err != nil {
		log.Fatalf("Failed to create sample families: %v", err)
	}

	fmt.Println("SQLite initialization completed successfully!")
}

func createSampleFamilies(ctx context.Context, repo *sqlite.SQLiteFamilyRepository) error {
	// Family 1: Traditional family with two parents and two children
	family1, err := createFamily1()
	if err != nil {
		return fmt.Errorf("failed to create family 1: %w", err)
	}
	if err := repo.Save(ctx, family1); err != nil {
		return fmt.Errorf("failed to save family 1: %w", err)
	}
	fmt.Println("Created family 1: Traditional family with two parents and two children")

	// Family 2: Family with two parents and one child
	family2, err := createFamily2()
	if err != nil {
		return fmt.Errorf("failed to create family 2: %w", err)
	}
	if err := repo.Save(ctx, family2); err != nil {
		return fmt.Errorf("failed to save family 2: %w", err)
	}
	fmt.Println("Created family 2: Family with two parents and one child")

	// Family 3: Divorced family with one parent and no children
	family3, err := createFamily3()
	if err != nil {
		return fmt.Errorf("failed to create family 3: %w", err)
	}
	if err := repo.Save(ctx, family3); err != nil {
		return fmt.Errorf("failed to save family 3: %w", err)
	}
	fmt.Println("Created family 3: Divorced family with one parent and no children")

	// Family 4: Family with one living parent, one deceased parent, and two children
	family4, err := createFamily4()
	if err != nil {
		return fmt.Errorf("failed to create family 4: %w", err)
	}
	if err := repo.Save(ctx, family4); err != nil {
		return fmt.Errorf("failed to save family 4: %w", err)
	}
	fmt.Println("Created family 4: Family with one living parent, one deceased parent, and two children")

	return nil
}

func createFamily1() (*entity.Family, error) {
	// Create parents
	parent1, err := entity.NewParent(
		"00000000-0000-0000-0000-000000000101",
		"John",
		"Smith",
		parseDate("1980-05-15T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	parent2, err := entity.NewParent(
		"00000000-0000-0000-0000-000000000102",
		"Jane",
		"Smith",
		parseDate("1982-08-22T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Create children
	child1, err := entity.NewChild(
		"00000000-0000-0000-0000-000000000103",
		"Emily",
		"Smith",
		parseDate("2010-03-12T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	child2, err := entity.NewChild(
		"00000000-0000-0000-0000-000000000104",
		"Michael",
		"Smith",
		parseDate("2012-11-05T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Create family
	return entity.NewFamily(
		"00000000-0000-0000-0000-000000000100",
		entity.Married,
		[]*entity.Parent{parent1, parent2},
		[]*entity.Child{child1, child2},
	)
}

func createFamily2() (*entity.Family, error) {
	// Create parents
	parent1, err := entity.NewParent(
		"00000000-0000-0000-0000-000000000201",
		"Robert",
		"Johnson",
		parseDate("1975-12-10T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	parent2, err := entity.NewParent(
		"00000000-0000-0000-0000-000000000202",
		"Maria",
		"Johnson",
		parseDate("1978-04-28T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Create child
	child1, err := entity.NewChild(
		"00000000-0000-0000-0000-000000000203",
		"David",
		"Johnson",
		parseDate("2008-07-19T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Create family
	return entity.NewFamily(
		"00000000-0000-0000-0000-000000000200",
		entity.Married,
		[]*entity.Parent{parent1, parent2},
		[]*entity.Child{child1},
	)
}

func createFamily3() (*entity.Family, error) {
	// Create parent
	parent1, err := entity.NewParent(
		"00000000-0000-0000-0000-000000000301",
		"Thomas",
		"Williams",
		parseDate("1970-09-30T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Create family
	return entity.NewFamily(
		"00000000-0000-0000-0000-000000000300",
		entity.Divorced,
		[]*entity.Parent{parent1},
		[]*entity.Child{},
	)
}

func createFamily4() (*entity.Family, error) {
	// Create parents
	parent1, err := entity.NewParent(
		"00000000-0000-0000-0000-000000000401",
		"Sarah",
		"Brown",
		parseDate("1985-02-14T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	deathDate := parseDate("2020-04-15T00:00:00Z")
	parent2, err := entity.NewParent(
		"00000000-0000-0000-0000-000000000402",
		"James",
		"Brown",
		parseDate("1983-11-08T00:00:00Z"),
		&deathDate,
	)
	if err != nil {
		return nil, err
	}

	// Create children
	child1, err := entity.NewChild(
		"00000000-0000-0000-0000-000000000403",
		"Olivia",
		"Brown",
		parseDate("2015-06-23T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	child2, err := entity.NewChild(
		"00000000-0000-0000-0000-000000000404",
		"William",
		"Brown",
		parseDate("2017-09-11T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Create family
	return entity.NewFamily(
		"00000000-0000-0000-0000-000000000400",
		entity.Widowed,
		[]*entity.Parent{parent1, parent2},
		[]*entity.Child{child1, child2},
	)
}

func parseDate(dateStr string) time.Time {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		log.Fatalf("Failed to parse date %s: %v", dateStr, err)
	}
	return t
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
