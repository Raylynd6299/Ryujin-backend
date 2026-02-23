package models

import "time"

// UserModel is the GORM representation of a user row.
// Kept strictly separate from the domain entity to protect domain purity.
type UserModel struct {
	ID                        string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email                     string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	HashedPassword            string     `gorm:"type:text;not null"`
	FirstName                 string     `gorm:"type:varchar(100);not null"`
	LastName                  string     `gorm:"type:varchar(100);not null"`
	DefaultSavingsCurrency    string     `gorm:"type:varchar(3);not null;default:'USD'"`
	DefaultInvestmentCurrency string     `gorm:"type:varchar(3);not null;default:'USD'"`
	Locale                    string     `gorm:"type:varchar(10);not null;default:'en'"`
	CreatedAt                 time.Time  `gorm:"not null;autoCreateTime"`
	UpdatedAt                 time.Time  `gorm:"not null;autoUpdateTime"`
	DeletedAt                 *time.Time `gorm:"index"`
}

// TableName sets the table name for the UserModel.
func (UserModel) TableName() string {
	return "users"
}
