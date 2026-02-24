package models

import "time"

// PurchaseGoalModel is the GORM representation of a purchase_goals row.
type PurchaseGoalModel struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string `gorm:"type:uuid;not null;index"`
	Name        string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text"`
	Icon        string `gorm:"type:varchar(10)"`
	// Target amount stored as int64 cents
	TargetAmountCents int64      `gorm:"not null"`
	Currency          string     `gorm:"type:varchar(3);not null"`
	Priority          string     `gorm:"type:varchar(10);not null;default:'medium'"`
	Deadline          *time.Time `gorm:"index"`
	IsCompleted       bool       `gorm:"not null;default:false"`
	CreatedAt         time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"not null;autoUpdateTime"`
}

func (PurchaseGoalModel) TableName() string { return "purchase_goals" }

// GoalContributionModel is the GORM representation of a goal_contributions row.
type GoalContributionModel struct {
	ID     string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	GoalID string `gorm:"type:uuid;not null;index"`
	UserID string `gorm:"type:uuid;not null;index"`
	// Amount stored as int64 cents
	AmountCents int64     `gorm:"not null"`
	Currency    string    `gorm:"type:varchar(3);not null"`
	Date        time.Time `gorm:"not null;index"`
	Notes       string    `gorm:"type:varchar(300)"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime"`
}

func (GoalContributionModel) TableName() string { return "goal_contributions" }
