package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	investErrors "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/errors"
	sharedVO "github.com/Raylynd6299/ryujin/internal/shared/domain/value_objects"
)

// Holding represents an investment position held by a user.
// It tracks a quantity of an asset purchased at a given price,
// along with optional real-time pricing data.
type Holding struct {
	ID     string
	UserID string

	Symbol    Symbol
	Name      string
	AssetType AssetType

	// Quantity stored as micro-units (1 share = 1_000_000)
	Quantity Quantity

	// Pricing
	BuyPrice     sharedVO.Money  // price per unit at purchase
	CurrentPrice *sharedVO.Money // nil until refreshed
	Currency     sharedVO.Currency

	Notes    string
	PricedAt *time.Time // when CurrentPrice was last updated

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewHolding creates and validates a new Holding entity
func NewHolding(
	id string,
	userID string,
	symbol string,
	name string,
	assetType string,
	quantityMicro int64,
	buyPriceCents int64,
	currency string,
	notes string,
) (*Holding, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, investErrors.NewHoldingValidationError("user_id", "user id cannot be empty")
	}

	if strings.TrimSpace(name) == "" {
		return nil, investErrors.NewHoldingValidationError("name", "holding name cannot be empty")
	}

	sym, err := NewSymbol(symbol)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("symbol", err.Error())
	}

	at, err := NewAssetType(assetType)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("asset_type", err.Error())
	}

	qty, err := NewQuantity(quantityMicro)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("quantity", err.Error())
	}

	if buyPriceCents <= 0 {
		return nil, investErrors.NewHoldingValidationError("buy_price", "buy price must be greater than zero")
	}

	cur, err := sharedVO.NewCurrency(currency)
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("currency", err.Error())
	}

	buyPrice, err := sharedVO.NewMoney(buyPriceCents, cur.Code())
	if err != nil {
		return nil, investErrors.NewHoldingValidationError("buy_price", err.Error())
	}

	holdingID := id
	if strings.TrimSpace(holdingID) == "" {
		holdingID = uuid.New().String()
	}

	now := time.Now()
	return &Holding{
		ID:           holdingID,
		UserID:       userID,
		Symbol:       sym,
		Name:         strings.TrimSpace(name),
		AssetType:    at,
		Quantity:     qty,
		BuyPrice:     *buyPrice,
		CurrentPrice: nil,
		Currency:     *cur,
		Notes:        strings.TrimSpace(notes),
		PricedAt:     nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update modifies the mutable fields of a Holding
func (h *Holding) Update(
	name string,
	quantityMicro int64,
	buyPriceCents int64,
	currency string,
	notes string,
) error {
	if strings.TrimSpace(name) == "" {
		return investErrors.NewHoldingValidationError("name", "holding name cannot be empty")
	}

	qty, err := NewQuantity(quantityMicro)
	if err != nil {
		return investErrors.NewHoldingValidationError("quantity", err.Error())
	}

	if buyPriceCents <= 0 {
		return investErrors.NewHoldingValidationError("buy_price", "buy price must be greater than zero")
	}

	cur, err := sharedVO.NewCurrency(currency)
	if err != nil {
		return investErrors.NewHoldingValidationError("currency", err.Error())
	}

	buyPrice, err := sharedVO.NewMoney(buyPriceCents, cur.Code())
	if err != nil {
		return investErrors.NewHoldingValidationError("buy_price", err.Error())
	}

	h.Name = strings.TrimSpace(name)
	h.Quantity = qty
	h.BuyPrice = *buyPrice
	h.Currency = *cur
	h.Notes = strings.TrimSpace(notes)
	h.UpdatedAt = time.Now()

	return nil
}

// RefreshPrice updates the current market price for this holding
func (h *Holding) RefreshPrice(priceCents int64, pricedAt time.Time) {
	price, _ := sharedVO.NewMoney(priceCents, h.Currency.Code())
	h.CurrentPrice = price
	h.PricedAt = &pricedAt
	h.UpdatedAt = time.Now()
}

// BelongsTo checks if this holding belongs to the given user
func (h *Holding) BelongsTo(userID string) bool {
	return h.UserID == userID
}

// MarketValue returns the total market value (CurrentPrice × Quantity).
// Returns nil if CurrentPrice has not been set.
func (h *Holding) MarketValue() *sharedVO.Money {
	if h.CurrentPrice == nil {
		return nil
	}
	mv := h.CurrentPrice.Multiply(h.Quantity.ToFloat())
	return mv
}

// UnrealizedGainLoss returns CurrentValue - CostBasis.
// Returns nil if CurrentPrice has not been set.
func (h *Holding) UnrealizedGainLoss() *sharedVO.Money {
	mv := h.MarketValue()
	if mv == nil {
		return nil
	}
	costBasis := h.BuyPrice.Multiply(h.Quantity.ToFloat())
	gl, _ := mv.Subtract(costBasis)
	return gl
}

// UnrealizedGainLossPct returns the unrealized gain/loss as a percentage.
// Returns nil if CurrentPrice is nil or BuyPrice is zero.
func (h *Holding) UnrealizedGainLossPct() *float64 {
	if h.CurrentPrice == nil {
		return nil
	}
	buyAmount := h.BuyPrice.Amount()
	if buyAmount == 0 {
		return nil
	}
	currentAmount := h.CurrentPrice.Amount()
	pct := (float64(currentAmount) - float64(buyAmount)) / float64(buyAmount) * 100.0
	return &pct
}
