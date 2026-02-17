package config

import (
	"os"
	"strconv"
	"time"
)

// DBConfig groups database-related environment values
type DBConfig struct {
	URL          string
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	SSLMode      string
	Params       string
	MaxOpenConns int
	MaxIdleConns int
}

// RedisConfig groups redis-related environment values
type RedisConfig struct {
	URL      string
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig groups JWT-related environment values
type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// ServerConfig groups server/runtime environment values
type ServerConfig struct {
	Port           string
	LogLevel       string
	RequestTimeout time.Duration
}

// CORSConfig groups CORS-related environment values
type CORSConfig struct {
	AllowedOrigins   string
	AllowedMethods   string
	AllowedHeaders   string
	AllowCredentials bool
}

// RateLimitConfig groups rate limit-related environment values
type RateLimitConfig struct {
	RPS     int
	Burst   int
	Enabled bool
}

// Config is the root configuration
type Config struct {
	DB        DBConfig
	JWT       JWTConfig
	Server    ServerConfig
	Redis     RedisConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
}

// App is the global config instance
var App Config

// Load reads environment variables into App and returns any load error
func Load() error {
	// Server
	App.Server.Port = getEnv("PORT", "8080")
	App.Server.LogLevel = getEnv("LOG_LEVEL", "info")
	App.Server.RequestTimeout = getEnvDuration("REQUEST_TIMEOUT", 30*time.Second)

	// Database
	App.DB.URL = getEnv("DATABASE_URL", "")
	App.DB.MaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25)
	App.DB.MaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 25)

	// JWT
	App.JWT.Secret = getEnv("JWT_SECRET", "")
	App.JWT.AccessTokenDuration = getEnvDuration("JWT_ACCESS_DURATION", 15*time.Minute)
	App.JWT.RefreshTokenDuration = getEnvDuration("JWT_REFRESH_DURATION", 24*time.Hour)

	// Database: support full parts or single DATABASE_URL
	if u := getEnv("DATABASE_URL", ""); u != "" {
		App.DB.URL = u
	} else {
		App.DB.Host = getEnv("DB_HOST", "localhost")
		App.DB.Port = getEnvInt("DB_PORT", 5432)
		App.DB.User = getEnv("DB_USER", "")
		App.DB.Password = getEnv("DB_PASSWORD", "")
		App.DB.Name = getEnv("DB_NAME", "")
		App.DB.SSLMode = getEnv("DB_SSLMODE", "disable")
		App.DB.Params = getEnv("DB_PARAMS", "")
	}
	App.DB.MaxOpenConns = getEnvInt("DB_MAX_OPEN_CONNS", 25)
	App.DB.MaxIdleConns = getEnvInt("DB_MAX_IDLE_CONNS", 25)

	// Redis / misc
	App.Redis.URL = getEnv("REDIS_URL", "")
	App.Redis.Host = getEnv("REDIS_HOST", "localhost")
	App.Redis.Port = getEnvInt("REDIS_PORT", 6379)
	App.Redis.Password = getEnv("REDIS_PASSWORD", "")
	App.Redis.DB = getEnvInt("REDIS_DB", 0)

	// CORS
	App.CORS.AllowedOrigins = getEnv("CORS_ALLOWED_ORIGINS", "*")
	App.CORS.AllowedMethods = getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	App.CORS.AllowedHeaders = getEnv("CORS_ALLOWED_HEADERS", "Authorization,Content-Type")
	App.CORS.AllowCredentials = getEnvBool("CORS_ALLOW_CREDENTIALS", false)

	// Rate limit
	App.RateLimit.RPS = getEnvInt("RATE_LIMIT_RPS", 10)
	App.RateLimit.Burst = getEnvInt("RATE_LIMIT_BURST", 20)
	App.RateLimit.Enabled = getEnvBool("RATE_LIMIT_ENABLED", true)

	return nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

// getEnvDuration attempts to parse a duration string (e.g. "30s") or an integer
// number of seconds. Falls back to def on parse errors.
func getEnvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
		if i, err := strconv.Atoi(v); err == nil {
			return time.Duration(i) * time.Second
		}
	}
	return def
}
