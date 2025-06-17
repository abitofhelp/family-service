// Package config provides configuration functionality for the application.
package config

import "time"

// Ensure Config implements the necessary interfaces
var (
	_ MongoDBConfigAdapter  = (*Config)(nil)
	_ PostgresConfigAdapter = (*Config)(nil)
	_ SQLiteConfigAdapter   = (*Config)(nil)
)

// MongoDBConfigAdapter adapts the Config to implement the ports.MongoDBConfig interface
type MongoDBConfigAdapter interface {
	// GetConnectionTimeout returns the timeout for establishing a database connection
	GetConnectionTimeout() time.Duration

	// GetPingTimeout returns the timeout for pinging the database to verify connection
	GetPingTimeout() time.Duration

	// GetDisconnectTimeout returns the timeout for disconnecting from the database
	GetDisconnectTimeout() time.Duration

	// GetIndexTimeout returns the timeout for creating database indexes
	GetIndexTimeout() time.Duration

	// GetURI returns the MongoDB connection URI
	GetURI() string
}

// PostgresConfigAdapter adapts the Config to implement the ports.PostgresConfig interface
type PostgresConfigAdapter interface {
	// GetConnectionTimeout returns the timeout for establishing a database connection
	GetConnectionTimeout() time.Duration

	// GetPingTimeout returns the timeout for pinging the database to verify connection
	GetPingTimeout() time.Duration

	// GetDisconnectTimeout returns the timeout for disconnecting from the database
	GetDisconnectTimeout() time.Duration

	// GetIndexTimeout returns the timeout for creating database indexes
	GetIndexTimeout() time.Duration

	// GetDSN returns the PostgreSQL connection data source name
	GetDSN() string
}

// SQLiteConfigAdapter adapts the Config to implement the ports.SQLiteConfig interface
type SQLiteConfigAdapter interface {
	// GetConnectionTimeout returns the timeout for establishing a database connection
	GetConnectionTimeout() time.Duration

	// GetPingTimeout returns the timeout for pinging the database to verify connection
	GetPingTimeout() time.Duration

	// GetDisconnectTimeout returns the timeout for disconnecting from the database
	GetDisconnectTimeout() time.Duration

	// GetMigrationTimeout returns the timeout for database migrations
	GetMigrationTimeout() time.Duration

	// GetSQLiteURI returns the SQLite connection URI
	GetSQLiteURI() string
}

// GetConnectionTimeout returns the MongoDB connection timeout
func (c *Config) GetConnectionTimeout() time.Duration {
	return c.Database.MongoDB.ConnectionTimeout
}

// GetPingTimeout returns the MongoDB ping timeout
func (c *Config) GetPingTimeout() time.Duration {
	return c.Database.MongoDB.PingTimeout
}

// GetDisconnectTimeout returns the MongoDB disconnect timeout
func (c *Config) GetDisconnectTimeout() time.Duration {
	return c.Database.MongoDB.DisconnectTimeout
}

// GetIndexTimeout returns the MongoDB index creation timeout
func (c *Config) GetIndexTimeout() time.Duration {
	return c.Database.MongoDB.IndexTimeout
}

// GetURI returns the MongoDB connection URI
func (c *Config) GetURI() string {
	return c.Database.MongoDB.URI
}

// GetDSN returns the PostgreSQL connection data source name
func (c *Config) GetDSN() string {
	return c.Database.Postgres.DSN
}

// GetSQLiteURI returns the SQLite connection URI
func (c *Config) GetSQLiteURI() string {
	return c.Database.SQLite.URI
}

// GetMigrationTimeout returns the migration timeout for the current database type
func (c *Config) GetMigrationTimeout() time.Duration {
	switch c.Database.Type {
	case "mongodb":
		return c.Database.MongoDB.MigrationTimeout
	case "postgres":
		return c.Database.Postgres.MigrationTimeout
	case "sqlite":
		return c.Database.SQLite.MigrationTimeout
	default:
		return 30 * time.Second // Default to 30 seconds
	}
}
