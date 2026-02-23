package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"

	financeErrors "github.com/Raylynd6299/ryujin/internal/modules/finance/domain/errors"
)

// CategoryType defines whether a category is for income or expense
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
	CategoryTypeBoth    CategoryType = "both"
)

// Category represents a classification for financial transactions.
// Categories can be system-defaults or user-created.
type Category struct {
	ID        string
	UserID    *string // nil = system/global category, non-nil = user-created
	Name      string
	Type      CategoryType
	Icon      string // emoji or icon identifier
	Color     string // hex color code
	IsDefault bool   // system default categories cannot be deleted

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewCategory creates a new user-defined Category
func NewCategory(userID string, name string, categoryType CategoryType, icon string, color string) (*Category, error) {
	if strings.TrimSpace(name) == "" {
		return nil, financeErrors.NewCategoryInvalidError("category name cannot be empty")
	}

	if categoryType != CategoryTypeIncome && categoryType != CategoryTypeExpense && categoryType != CategoryTypeBoth {
		return nil, financeErrors.NewCategoryInvalidError("invalid category type")
	}

	now := time.Now()
	return &Category{
		ID:        uuid.New().String(),
		UserID:    &userID,
		Name:      strings.TrimSpace(name),
		Type:      categoryType,
		Icon:      icon,
		Color:     color,
		IsDefault: false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NewSystemCategory creates a system-wide default category (no user ownership)
func NewSystemCategory(name string, categoryType CategoryType, icon string, color string) *Category {
	now := time.Now()
	return &Category{
		ID:        uuid.New().String(),
		UserID:    nil,
		Name:      name,
		Type:      categoryType,
		Icon:      icon,
		Color:     color,
		IsDefault: true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates the category metadata
func (c *Category) Update(name string, icon string, color string) error {
	if strings.TrimSpace(name) == "" {
		return financeErrors.NewCategoryInvalidError("category name cannot be empty")
	}
	if c.IsDefault {
		return financeErrors.NewCategoryInvalidError("cannot modify a system default category")
	}

	c.Name = strings.TrimSpace(name)
	c.Icon = icon
	c.Color = color
	c.UpdatedAt = time.Now()
	return nil
}

// IsUserOwned returns true if this category belongs to a specific user
func (c *Category) IsUserOwned() bool {
	return c.UserID != nil
}

// BelongsTo returns true if the category belongs to the given user or is a system category
func (c *Category) BelongsTo(userID string) bool {
	if c.UserID == nil {
		return true // system categories are accessible to everyone
	}
	return *c.UserID == userID
}
