package models

import "time"

// HoldingModel is the GORM representation of an investment holding row.
type HoldingModel struct {
	ID     string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID string `gorm:"type:uuid;not null;index"`

	Symbol    string `gorm:"size:10;not null;index;constraint:fk_holdings_symbol,OnDelete:RESTRICT"`
	Name      string `gorm:"type:varchar(150);not null"`
	AssetType string `gorm:"type:varchar(20);not null"`

	// Quantity stored as micro-units (1 share = 1_000_000)
	QuantityMicro int64 `gorm:"not null"`

	// Buy price per unit in smallest currency unit (cents)
	BuyPriceCents int64  `gorm:"not null"`
	BuyCurrency   string `gorm:"type:varchar(3);not null"`

	// Current market price — nil until first price refresh
	CurrentPriceCents *int64     `gorm:"column:current_price_cents"`
	PricedAt          *time.Time `gorm:"column:priced_at"`

	Notes string `gorm:"type:text"`

	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`

	// Association — populated on explicit Preload
	StockQuote StockQuoteModel `gorm:"foreignKey:Symbol;references:Symbol"`
}

// TableName returns the GORM table name
func (HoldingModel) TableName() string { return "holdings" }

// StockQuoteModel is the GORM model for the global stock_quotes table.
// Symbol is the natural primary key (uppercase ticker, max 10 chars).
// All monetary fields are stored as int64 cents to avoid floating-point errors.
type StockQuoteModel struct {
	Symbol             string    `gorm:"primaryKey;size:10"`
	Name               string    `gorm:"not null"`
	Currency           string    `gorm:"size:3;not null"`
	PriceCents         int64     `gorm:"not null;default:0"`
	PreviousCloseCents int64     `gorm:"not null;default:0"`
	OpenCents          int64     `gorm:"not null;default:0"`
	DayHighCents       int64     `gorm:"not null;default:0"`
	DayLowCents        int64     `gorm:"not null;default:0"`
	Volume             int64     `gorm:"not null;default:0"`
	MarketCapCents     int64     `gorm:"not null;default:0"`
	Week52HighCents    int64     `gorm:"not null;default:0"`
	Week52LowCents     int64     `gorm:"not null;default:0"`
	TrailingPE         float64   `gorm:"not null;default:0"`
	ForwardPE          float64   `gorm:"not null;default:0"`
	DividendYield      float64   `gorm:"not null;default:0"`
	EPS                float64   `gorm:"not null;default:0"`
	FetchedAt          time.Time `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// TableName returns the GORM table name
func (StockQuoteModel) TableName() string { return "stock_quotes" }

// StockPriceHistoryModel is the GORM model for price history snapshots.
// Records are append-only — no update or delete operations.
type StockPriceHistoryModel struct {
	ID         string    `gorm:"primaryKey;size:36"`
	Symbol     string    `gorm:"size:10;not null;index"`
	PriceCents int64     `gorm:"not null"`
	Currency   string    `gorm:"size:3;not null"`
	RecordedAt time.Time `gorm:"not null;index"`
}

// TableName returns the GORM table name
func (StockPriceHistoryModel) TableName() string { return "stock_price_history" }
