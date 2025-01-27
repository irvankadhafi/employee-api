package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"time"
)

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.SetConfigName("config")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Warningf("%v", err)
	}
}

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}

// HTTPPort :nodoc:
func HTTPPort() string {
	return viper.GetString("ports.http")
}

// GRPCPort :nodoc:
func GRPCPort() string {
	return viper.GetString("ports.grpc")
}

// DatabaseDSN :nodoc:
func DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		DatabaseUsername(),
		DatabasePassword(),
		DatabaseHost(),
		DatabaseName(),
		DatabaseSSLMode())
}

// DatabaseHost :nodoc:
func DatabaseHost() string {
	return viper.GetString("postgres.host")
}

func DatabasePort() string {
	return viper.GetString("postgres.port")
}

// DatabaseName :nodoc:
func DatabaseName() string {
	return viper.GetString("postgres.database")
}

// DatabaseUsername :nodoc:
func DatabaseUsername() string {
	return viper.GetString("postgres.username")
}

// DatabasePassword :nodoc:
func DatabasePassword() string {
	return viper.GetString("postgres.password")
}

// DatabaseSSLMode :nodoc:
func DatabaseSSLMode() string {
	if viper.IsSet("postgres.sslmode") {
		return viper.GetString("postgres.sslmode")
	}
	return "disable"
}

// DatabaseTimezone :nodoc:
func DatabaseTimezone() string {
	if viper.IsSet("postgres.timezone") {
		return viper.GetString("postgres.timezone")
	}
	return "Asia/Jakarta"
}

// DatabasePingInterval :nodoc:
func DatabasePingInterval() time.Duration {
	if viper.GetInt("postgres.ping_interval") <= 0 {
		return DefaultDatabasePingInterval
	}
	return time.Duration(viper.GetInt("postgres.ping_interval")) * time.Millisecond
}

// DatabaseRetryAttempts :nodoc:
func DatabaseRetryAttempts() float64 {
	if viper.GetInt("postgres.retry_attempts") > 0 {
		return float64(viper.GetInt("postgres.retry_attempts"))
	}
	return DefaultDatabaseRetryAttempts
}

// DatabaseMaxIdleConns :nodoc:
func DatabaseMaxIdleConns() int {
	if viper.GetInt("postgres.max_idle_conns") <= 0 {
		return DefaultDatabaseMaxIdleConns
	}
	return viper.GetInt("postgres.max_idle_conns")
}

// DatabaseMaxOpenConns :nodoc:
func DatabaseMaxOpenConns() int {
	if viper.GetInt("postgres.max_open_conns") <= 0 {
		return DefaultDatabaseMaxOpenConns
	}
	return viper.GetInt("postgres.max_open_conns")
}

// DatabaseConnMaxLifetime :nodoc:
func DatabaseConnMaxLifetime() time.Duration {
	if !viper.IsSet("postgres.conn_max_lifetime") {
		return DefaultDatabaseConnMaxLifetime
	}
	return time.Duration(viper.GetInt("postgres.conn_max_lifetime")) * time.Millisecond
}

// DatabaseRetryAttempts :nodoc:
func DatabaseTimeout() int {
	if viper.GetInt("postgres.timeout") > 0 {
		return viper.GetInt("postgres.timeout")
	}
	return DefaultDatabaseTimeout
}

// RedisCacheHost :nodoc:
func RedisCacheHost() string {
	return viper.GetString("redis.cache_host")
}

// RedisLockHost :nodoc:
func RedisLockHost() string {
	return viper.GetString("redis.lock_host")
}

// DisableCaching :nodoc:
func DisableCaching() bool {
	return viper.GetBool("disable_caching")
}

// RedisDialTimeout :nodoc:
func RedisDialTimeout() time.Duration {
	cfg := viper.GetString("redis.dial_timeout")
	return parseDuration(cfg, 5*time.Second)
}

// RedisWriteTimeout :nodoc:
func RedisWriteTimeout() time.Duration {
	cfg := viper.GetString("redis.write_timeout")
	return parseDuration(cfg, 2*time.Second)
}

// RedisReadTimeout :nodoc:
func RedisReadTimeout() time.Duration {
	cfg := viper.GetString("redis.read_timeout")
	return parseDuration(cfg, 2*time.Second)
}

// RedisMaxIdleConn :nodoc:
func RedisMaxIdleConn() int {
	if viper.GetInt("redis.max_idle_conn") > 0 {
		return viper.GetInt("redis.max_idle_conn")
	}
	return 20
}

// RedisMaxActiveConn :nodoc:
func RedisMaxActiveConn() int {
	if viper.GetInt("redis.max_active_conn") > 0 {
		return viper.GetInt("redis.max_active_conn")
	}
	return 50
}

// RedisCacheTTL :nodoc:
func CacheTTL() time.Duration {
	cfg := viper.GetString("cache_ttl")
	return parseDuration(cfg, DefaultRedisCacheTTL)
}

func parseDuration(in string, defaultDuration time.Duration) time.Duration {
	dur, err := time.ParseDuration(in)
	if err != nil {
		return defaultDuration
	}
	return dur
}
