package dto

import (
	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
)

// --- Request DTOs ---

// CreateHoldingRequest is used to create a new investment holding
type CreateHoldingRequest struct {
	Symbol        string `json:"symbol" binding:"required,min=1,max=10"`
	Name          string `json:"name" binding:"required"`
	AssetType     string `json:"assetType" binding:"required"`
	QuantityMicro int64  `json:"quantityMicro" binding:"required,min=1"`
	BuyPriceCents int64  `json:"buyPriceCents" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,len=3"`
	Notes         string `json:"notes"`
}

// UpdateHoldingRequest is used to update an existing holding
type UpdateHoldingRequest struct {
	Name          string `json:"name"`
	QuantityMicro int64  `json:"quantityMicro" binding:"min=1"`
	BuyPriceCents int64  `json:"buyPriceCents" binding:"min=1"`
	Currency      string `json:"currency" binding:"omitempty,len=3"`
	Notes         string `json:"notes"`
}

// --- Response DTOs ---

// HoldingResponse is returned for single holding operations
type HoldingResponse struct {
	ID                      string   `json:"id"`
	Symbol                  string   `json:"symbol"`
	Name                    string   `json:"name"`
	AssetType               string   `json:"assetType"`
	QuantityMicro           int64    `json:"quantityMicro"`
	QuantityFloat           float64  `json:"quantityFloat"`
	BuyPriceCents           int64    `json:"buyPriceCents"`
	Currency                string   `json:"currency"`
	CurrentPriceCents       *int64   `json:"currentPriceCents"`
	MarketValueCents        *int64   `json:"marketValueCents"`
	UnrealizedGainLossCents *int64   `json:"unrealizedGainLossCents"`
	UnrealizedGainLossPct   *float64 `json:"unrealizedGainLossPct"`
	PricedAt                *string  `json:"pricedAt"`
	Notes                   string   `json:"notes"`
	CreatedAt               string   `json:"createdAt"`
	UpdatedAt               string   `json:"updatedAt"`
}

// HoldingListResponse wraps a paginated list of holdings
type HoldingListResponse struct {
	Holdings []HoldingResponse `json:"holdings"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

// ToHoldingResponse maps a domain Holding entity to a HoldingResponse DTO
func ToHoldingResponse(h *entities.Holding) HoldingResponse {
	resp := HoldingResponse{
		ID:            h.ID,
		Symbol:        h.Symbol.Value(),
		Name:          h.Name,
		AssetType:     h.AssetType.String(),
		QuantityMicro: h.Quantity.MicroUnits(),
		QuantityFloat: h.Quantity.ToFloat(),
		BuyPriceCents: h.BuyPrice.Amount(),
		Currency:      h.Currency.Code(),
		Notes:         h.Notes,
		CreatedAt:     h.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     h.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}

	if h.CurrentPrice != nil {
		cents := h.CurrentPrice.Amount()
		resp.CurrentPriceCents = &cents
	}

	if mv := h.MarketValue(); mv != nil {
		mvCents := mv.Amount()
		resp.MarketValueCents = &mvCents
	}

	if gl := h.UnrealizedGainLoss(); gl != nil {
		glCents := gl.Amount()
		resp.UnrealizedGainLossCents = &glCents
	}

	resp.UnrealizedGainLossPct = h.UnrealizedGainLossPct()

	if h.PricedAt != nil {
		ts := h.PricedAt.UTC().Format("2006-01-02T15:04:05Z")
		resp.PricedAt = &ts
	}

	return resp
}
