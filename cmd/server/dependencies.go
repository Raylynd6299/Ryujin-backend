package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Raylynd6299/Ryujin-backend/internal/config"
	financeAppServices "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/application/services"
	financeControllers "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/infrastructure/http/controllers"
	financeRepos "github.com/Raylynd6299/Ryujin-backend/internal/modules/finance/infrastructure/persistence/repositories"
	goalAppServices "github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/application/services"
	goalControllers "github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/infrastructure/http/controllers"
	goalRepos "github.com/Raylynd6299/Ryujin-backend/internal/modules/goal/infrastructure/persistence/repositories"
	investAppServices "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/application/services"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/external"
	investControllers "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/http/controllers"
	investRepos "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/persistence/repositories"
	"github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/infrastructure/worker"
	userAppServices "github.com/Raylynd6299/Ryujin-backend/internal/modules/user/application/services"
	userControllers "github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/http/controllers"
	userRepos "github.com/Raylynd6299/Ryujin-backend/internal/modules/user/infrastructure/persistence/repositories"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/infrastructure/persistence"
	"github.com/Raylynd6299/Ryujin-backend/internal/shared/utils"
	"github.com/Raylynd6299/Ryujin-backend/migrations"
)

// AppDependencies holds all application dependencies
type AppDependencies struct {
	DB         *gorm.DB
	Engine     *gin.Engine
	JWTService *utils.JWTService

	// User Module
	AuthController    *userControllers.AuthController
	ProfileController *userControllers.ProfileController

	// Finance Module
	CategoryController     *financeControllers.CategoryController
	IncomeSourceController *financeControllers.IncomeSourceController
	ExpenseController      *financeControllers.ExpenseController
	DebtController         *financeControllers.DebtController
	AccountController      *financeControllers.AccountController
	IndicesController      *financeControllers.IndicesController

	// Investment Module
	HoldingController       *investControllers.HoldingController
	PortfolioController     *investControllers.PortfolioController
	StockAnalysisController *investControllers.StockAnalysisController
	PriceRefreshWorker      *worker.PriceRefreshWorker

	// Goal Module
	GoalController *goalControllers.GoalController
}

// NewAppDependencies creates and initializes all application dependencies
func NewAppDependencies(cfg *config.Config) (*AppDependencies, error) {
	// Initialize database connection
	db, err := initDatabase(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Run pending migrations before wiring any module
	if err := persistence.RunMigrations(db, migrations.FS, migrations.Dir); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize Gin engine
	engine := initGinEngine(cfg)

	// ── Shared services ──────────────────────────────────────────────────────
	jwtService := utils.NewJWTService(cfg.JWT.Secret)

	// ── User Module ──────────────────────────────────────────────────────────
	userRepo := userRepos.NewUserRepositoryGorm(db)
	authService := userAppServices.NewAuthService(userRepo, jwtService, cfg.JWT.AccessTokenDuration, cfg.JWT.RefreshTokenDuration)
	profileService := userAppServices.NewProfileService(userRepo)
	authCtrl := userControllers.NewAuthController(authService)
	profileCtrl := userControllers.NewProfileController(profileService)

	log.Println("✓ User module initialized")

	// ── Finance Module ────────────────────────────────────────────────────────
	categoryRepo := financeRepos.NewCategoryRepositoryGorm(db)
	incomeRepo := financeRepos.NewIncomeSourceRepositoryGorm(db)
	expenseRepo := financeRepos.NewExpenseRepositoryGorm(db)
	debtRepo := financeRepos.NewDebtRepositoryGorm(db)
	accountRepo := financeRepos.NewAccountRepositoryGorm(db)

	categoryService := financeAppServices.NewCategoryService(categoryRepo)
	incomeService := financeAppServices.NewIncomeSourceService(incomeRepo)
	expenseService := financeAppServices.NewExpenseService(expenseRepo)
	debtService := financeAppServices.NewDebtService(debtRepo)
	accountService := financeAppServices.NewAccountService(accountRepo)

	categoryCtrl := financeControllers.NewCategoryController(categoryService)
	incomeCtrl := financeControllers.NewIncomeSourceController(incomeService)
	expenseCtrl := financeControllers.NewExpenseController(expenseService)
	debtCtrl := financeControllers.NewDebtController(debtService)
	accountCtrl := financeControllers.NewAccountController(accountService)

	log.Println("✓ Finance module initialized")

	// ── Investment Module ─────────────────────────────────────────────────────
	yahooClient := external.NewYahooFinanceClient()
	alphaVantageClient := external.NewAlphaVantageClient(cfg.ExternalAPI.AlphaVantageAPIKey)
	compositePriceProvider := external.NewCompositePriceProvider(yahooClient, alphaVantageClient)

	holdingRepo := investRepos.NewHoldingRepositoryGorm(db)
	stockQuoteRepo := investRepos.NewStockQuoteRepositoryGorm(db)
	stockPriceHistoryRepo := investRepos.NewStockPriceHistoryRepositoryGorm(db)

	holdingService := investAppServices.NewHoldingService(holdingRepo, compositePriceProvider, stockQuoteRepo, yahooClient)
	portfolioService := investAppServices.NewPortfolioService(holdingRepo)
	stockAnalysisService := investAppServices.NewStockAnalysisService(stockQuoteRepo, stockPriceHistoryRepo, yahooClient)

	holdingCtrl := investControllers.NewHoldingController(holdingService)
	portfolioCtrl := investControllers.NewPortfolioController(portfolioService)
	stockAnalysisCtrl := investControllers.NewStockAnalysisController(stockAnalysisService)

	priceRefreshWorker := worker.NewPriceRefreshWorker(
		stockQuoteRepo,
		stockPriceHistoryRepo,
		yahooClient,
		15*time.Minute,
	)

	log.Println("✓ Investment module initialized")

	// ── Indices (cross-cutting: finance + investment) ─────────────────────────
	indicesService := financeAppServices.NewIndicesCalculatorService(incomeRepo, expenseRepo, debtRepo, accountRepo, holdingRepo)
	indicesCtrl := financeControllers.NewIndicesController(indicesService)

	log.Println("✓ Finance indices initialized")

	// ── Goal Module ───────────────────────────────────────────────────────────
	goalRepo := goalRepos.NewGoalRepositoryGorm(db)
	goalContributionRepo := goalRepos.NewGoalContributionRepositoryGorm(db)
	goalService := goalAppServices.NewGoalService(goalRepo, goalContributionRepo)
	goalCtrl := goalControllers.NewGoalController(goalService)

	log.Println("✓ Goal module initialized")

	return &AppDependencies{
		DB:         db,
		Engine:     engine,
		JWTService: jwtService,

		AuthController:    authCtrl,
		ProfileController: profileCtrl,

		CategoryController:     categoryCtrl,
		IncomeSourceController: incomeCtrl,
		ExpenseController:      expenseCtrl,
		DebtController:         debtCtrl,
		AccountController:      accountCtrl,
		IndicesController:      indicesCtrl,

		HoldingController:       holdingCtrl,
		PortfolioController:     portfolioCtrl,
		StockAnalysisController: stockAnalysisCtrl,
		PriceRefreshWorker:      priceRefreshWorker,

		GoalController: goalCtrl,
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
