package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	appServices "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/application/services"
	investErrors "github.com/Raylynd6299/Ryujin-backend/internal/modules/investment/domain/errors"
	sharedHTTP "github.com/Raylynd6299/Ryujin-backend/internal/shared/infrastructure/http"
)

const (
	defaultPriceHistoryLimit = 30
	maxPriceHistoryLimit     = 100
)

// StockAnalysisController handles stock quote and analysis endpoints
type StockAnalysisController struct {
	stockAnalysisService *appServices.StockAnalysisService
}

// NewStockAnalysisController creates a new StockAnalysisController
func NewStockAnalysisController(service *appServices.StockAnalysisService) *StockAnalysisController {
	return &StockAnalysisController{stockAnalysisService: service}
}

// GetStockQuote godoc
// GET /api/v1/stocks/:symbol/quote
func (ctrl *StockAnalysisController) GetStockQuote(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))
	if symbol == "" {
		sharedHTTP.BadRequestResponse(c, "symbol is required", nil)
		return
	}

	quote, err := ctrl.stockAnalysisService.GetStockQuote(c.Request.Context(), symbol)
	if err != nil {
		if errors.Is(err, investErrors.ErrStockQuoteNotFound) || strings.Contains(err.Error(), "not found") {
			sharedHTTP.NotFoundResponse(c, "stock quote not found for symbol: "+symbol)
			return
		}
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, quote, "")
}

// ListStockQuotes godoc
// GET /api/v1/stocks
func (ctrl *StockAnalysisController) ListStockQuotes(c *gin.Context) {
	quotes, err := ctrl.stockAnalysisService.ListStockQuotes(c.Request.Context())
	if err != nil {
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, quotes, "")
}

// GetStockPriceHistory godoc
// GET /api/v1/stocks/:symbol/history?limit=N
func (ctrl *StockAnalysisController) GetStockPriceHistory(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))
	if symbol == "" {
		sharedHTTP.BadRequestResponse(c, "symbol is required", nil)
		return
	}

	limit := defaultPriceHistoryLimit
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed <= 0 {
			sharedHTTP.BadRequestResponse(c, "limit must be a positive integer", nil)
			return
		}
		if parsed > maxPriceHistoryLimit {
			parsed = maxPriceHistoryLimit
		}
		limit = parsed
	}

	history, err := ctrl.stockAnalysisService.GetPriceHistory(c.Request.Context(), symbol, limit)
	if err != nil {
		if errors.Is(err, investErrors.ErrStockQuoteNotFound) {
			sharedHTTP.NotFoundResponse(c, "no price history found for symbol: "+symbol)
			return
		}
		handleInvestmentError(c, err)
		return
	}

	sharedHTTP.OkResponse(c, history, "")
}

// handleStockError maps stock-specific domain errors to HTTP responses.
// Falls back to handleInvestmentError for all other error types.
func handleStockError(c *gin.Context, err error, symbol string) {
	if errors.Is(err, investErrors.ErrStockQuoteNotFound) {
		sharedHTTP.NotFoundResponse(c, "stock quote not found for symbol: "+symbol)
		return
	}
	if errors.Is(err, investErrors.ErrInvalidSymbol) {
		sharedHTTP.ErrorResponse(c, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}
	handleInvestmentError(c, err)
}
