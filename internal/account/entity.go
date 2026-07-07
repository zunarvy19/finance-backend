package account

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	Name string `gorm:"size:100;not null"`

	Currency string `gorm:"type:char(3);default:'IDR'"`

	Balance int64 `gorm:"not null;default:0"`

	Version int64 `gorm:"not null;default:1"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (Account) TableName() string {
	return "accounts"
}
