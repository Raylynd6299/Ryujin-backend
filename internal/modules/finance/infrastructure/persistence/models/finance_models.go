package models

import "time"

// CategoryModel is the GORM representation of a category row.
type CategoryModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    *string   `gorm:"type:uuid;index"` // nil = system category
	Name      string    `gorm:"type:varchar(100);not null"`
	Type      string    `gorm:"type:varchar(20);not null"` // income | expense | both
	Icon      string    `gorm:"type:varchar(50)"`
	Color     string    `gorm:"type:varchar(20)"`
	IsDefault bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`
}

func (CategoryModel) TableName() string { return "categories" }

// IncomeSourceModel is the GORM representation of an income source row.
type IncomeSourceModel struct {
	ID          string  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string  `gorm:"type:uuid;not null;index"`
	CategoryID  *string `gorm:"type:uuid;index"`
	Name        string  `gorm:"type:varchar(150);not null"`
	Description string  `gorm:"type:text"`
	// Money stored as int64 cents
	AmountCents int64     `gorm:"not null"`
	Currency    string    `gorm:"type:varchar(3);not null"`
	IncomeType  string    `gorm:"type:varchar(30);not null"`
	Recurrence  string    `gorm:"type:varchar(20);not null"`
	StartDate   time.Time `gorm:"not null"`
	EndDate     *time.Time
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

func (IncomeSourceModel) TableName() string { return "income_sources" }

// ExpenseModel is the GORM representation of an expense row.
type ExpenseModel struct {
	ID          string  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string  `gorm:"type:uuid;not null;index"`
	CategoryID  *string `gorm:"type:uuid;index"`
	Name        string  `gorm:"type:varchar(150);not null"`
	Description string  `gorm:"type:text"`
	// Money stored as int64 cents
	AmountCents int64     `gorm:"not null"`
	Currency    string    `gorm:"type:varchar(3);not null"`
	Priority    string    `gorm:"type:varchar(20);not null"`
	Recurrence  string    `gorm:"type:varchar(20);not null"`
	ExpenseDate time.Time `gorm:"not null"`
	EndDate     *time.Time
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

func (ExpenseModel) TableName() string { return "expenses" }

// DebtModel is the GORM representation of a debt row.
type DebtModel struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string `gorm:"type:uuid;not null;index"`
	Name        string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text"`
	DebtType    string `gorm:"type:varchar(30);not null"`
	// All amounts stored as int64 cents
	TotalAmountCents     int64   `gorm:"not null"`
	RemainingAmountCents int64   `gorm:"not null"`
	MonthlyPaymentCents  int64   `gorm:"not null"`
	Currency             string  `gorm:"type:varchar(3);not null"`
	InterestRate         float64 `gorm:"not null;default:0"`
	StartDate            *time.Time
	DueDate              *time.Time
	IsActive             bool      `gorm:"not null;default:true"`
	CreatedAt            time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"not null;autoUpdateTime"`
}

func (DebtModel) TableName() string { return "debts" }

// AccountModel is the GORM representation of a financial account row.
type AccountModel struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string `gorm:"type:uuid;not null;index"`
	Name        string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text"`
	AccountType string `gorm:"type:varchar(20);not null"`
	// Balance stored as int64 cents
	BalanceCents int64     `gorm:"not null;default:0"`
	Currency     string    `gorm:"type:varchar(3);not null"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime"`
}

func (AccountModel) TableName() string { return "accounts" }
