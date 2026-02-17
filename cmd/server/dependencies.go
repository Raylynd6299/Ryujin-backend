package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Raylynd6299/ryujin/internal/config"
)

// AppDependencies holds all application dependencies
type AppDependencies struct {
	DB     *gorm.DB
	Engine *gin.Engine
}

// NewAppDependencies creates and initializes all application dependencies
func NewAppDependencies(cfg *config.Config) (*AppDependencies, error) {
	// Initialize database connection
	db, err := initDatabase(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Gin engine
	engine := initGinEngine(cfg)

	return &AppDependencies{
		DB:     db,
		Engine: engine,
	}, nil
}

// initDatabase initializes the database connection using GORM
func initDatabase(dbConfig config.DBConfig) (*gorm.DB, error) {
	var dsn string

	// Use DATABASE_URL if provided, otherwise build from individual params
	if dbConfig.URL != "" {
		dsn = dbConfig.URL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.SSLMode,
		)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✓ Database connection established")
	return db, nil
}

// initGinEngine creates and configures the Gin engine
func initGinEngine(cfg *config.Config) *gin.Engine {
	// Set Gin mode (can be configured via GIN_MODE env var)
	gin.SetMode(gin.DebugMode)

	// Create new Gin engine
	engine := gin.New()

	log.Println("✓ Gin engine initialized")
	return engine
}

// Close closes all application dependencies gracefully
func (deps *AppDependencies) Close() error {
	if deps.DB != nil {
		sqlDB, err := deps.DB.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		log.Println("✓ Database connection closed")
	}
	return nil
}
