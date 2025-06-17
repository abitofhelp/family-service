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
	_ "github.com/mattn/go-sqlite3" // SQLite driver
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

	// Create repository
	repo := sqlite.NewSQLiteFamilyRepository(db)

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
		"p1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
		"John",
		"Smith",
		parseDate("1980-05-15T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	parent2, err := entity.NewParent(
		"p2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"c1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
		"Emily",
		"Smith",
		parseDate("2010-03-12T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	child2, err := entity.NewChild(
		"c2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"f1a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
		entity.Married,
		[]*entity.Parent{parent1, parent2},
		[]*entity.Child{child1, child2},
	)
}

func createFamily2() (*entity.Family, error) {
	// Create parents
	parent1, err := entity.NewParent(
		"p3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
		"Robert",
		"Johnson",
		parseDate("1975-12-10T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	parent2, err := entity.NewParent(
		"p4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"c3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"f2a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
		entity.Married,
		[]*entity.Parent{parent1, parent2},
		[]*entity.Child{child1},
	)
}

func createFamily3() (*entity.Family, error) {
	// Create parent
	parent1, err := entity.NewParent(
		"p5a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"f3a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
		entity.Divorced,
		[]*entity.Parent{parent1},
		[]*entity.Child{},
	)
}

func createFamily4() (*entity.Family, error) {
	// Create parents
	parent1, err := entity.NewParent(
		"p6a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"p7a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"c4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
		"Olivia",
		"Brown",
		parseDate("2015-06-23T00:00:00Z"),
		nil,
	)
	if err != nil {
		return nil, err
	}

	child2, err := entity.NewChild(
		"c5a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
		"f4a2b3c4-d5e6-f7a8-b9c0-d1e2f3a4b5c6",
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
