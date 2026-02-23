package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Raylynd6299/ryujin/internal/modules/investment/domain/entities"
	investErrors "github.com/Raylynd6299/ryujin/internal/modules/investment/domain/errors"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/persistence/mappers"
	"github.com/Raylynd6299/ryujin/internal/modules/investment/infrastructure/persistence/models"
)

// allowedSortColumns whitelists valid column names to prevent SQL injection
var allowedSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"symbol":     true,
	"name":       true,
	"asset_type": true,
	"buy_price":  true,
	"priced_at":  true,
}

// allowedOrders whitelists valid sort directions
var allowedOrders = map[string]bool{
	"asc":  true,
	"desc": true,
}

// HoldingRepositoryGorm implements domain/repositories.HoldingRepository using GORM
type HoldingRepositoryGorm struct {
	db *gorm.DB
}

// NewHoldingRepositoryGorm creates a new GORM-backed holding repository
func NewHoldingRepositoryGorm(db *gorm.DB) *HoldingRepositoryGorm {
	return &HoldingRepositoryGorm{db: db}
}

// Create persists a new holding to the database
func (r *HoldingRepositoryGorm) Create(ctx context.Context, h *entities.Holding) error {
	model := mappers.HoldingToModel(h)
	return r.db.WithContext(ctx).Create(model).Error
}

// FindByID returns a holding by ID scoped to userID.
// Returns HoldingNotFoundError for both missing records and ownership mismatches
// to avoid leaking the existence of other users' data.
func (r *HoldingRepositoryGorm) FindByID(ctx context.Context, id, userID string) (*entities.Holding, error) {
	var model models.HoldingModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, investErrors.NewHoldingNotFoundError(id)
		}
		return nil, err
	}
	return mappers.ModelToHolding(&model)
}

// FindAllByUserID returns paginated holdings for a user along with the total count.
// sort must be a whitelisted column name; order must be "asc" or "desc".
// Invalid values fall back to defaults (created_at DESC).
func (r *HoldingRepositoryGorm) FindAllByUserID(ctx context.Context, userID string, page, limit int, sort, order string) ([]*entities.Holding, int64, error) {
	var modelList []models.HoldingModel
	var total int64

	// Sanitize pagination
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Sanitize sort/order to prevent SQL injection
	if !allowedSortColumns[sort] {
		sort = "created_at"
	}
	if !allowedOrders[order] {
		order = "desc"
	}

	orderClause := fmt.Sprintf("%s %s", sort, order)
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&models.HoldingModel{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order(orderClause).
		Offset(offset).
		Limit(limit).
		Find(&modelList).Error
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entities.Holding, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.ModelToHolding(&m)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, domain)
	}
	return result, total, nil
}

// FindActiveByUserID returns all holdings for a user without pagination
func (r *HoldingRepositoryGorm) FindActiveByUserID(ctx context.Context, userID string) ([]*entities.Holding, error) {
	var modelList []models.HoldingModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("symbol ASC").
		Find(&modelList).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Holding, 0, len(modelList))
	for _, m := range modelList {
		m := m
		domain, err := mappers.ModelToHolding(&m)
		if err != nil {
			return nil, err
		}
		result = append(result, domain)
	}
	return result, nil
}

// Update persists changes to an existing holding
func (r *HoldingRepositoryGorm) Update(ctx context.Context, h *entities.Holding) error {
	model := mappers.HoldingToModel(h)
	return r.db.WithContext(ctx).Save(model).Error
}

// Delete removes a holding scoped to userID.
// Returns HoldingNotFoundError if the holding does not exist or belongs to another user.
func (r *HoldingRepositoryGorm) Delete(ctx context.Context, id, userID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.HoldingModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return investErrors.NewHoldingNotFoundError(id)
	}
	return nil
}
