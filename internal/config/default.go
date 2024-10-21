package config

import "time"

const (
	DefaultDatabaseMaxIdleConns    = 50
	DefaultDatabaseMaxOpenConns    = 100
	DefaultDatabaseConnMaxLifetime = 1 * time.Hour
	DefaultDatabasePingInterval    = 1 * time.Second
	DefaultDatabaseRetryAttempts   = 3
	DefaultDatabaseTimeout         = 120

	DefaultRedisCacheTTL = 15 * time.Minute
)
