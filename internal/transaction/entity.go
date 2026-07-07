package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	AccountID   uuid.UUID `gorm:"type:uuid;not null;index"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null"`
	ClientID    uuid.UUID `gorm:"type:uuid;not null"`
	Type        string    `gorm:"size:20;not null"` // 'INCOME' or 'EXPENSE'
	Amount      int64     `gorm:"not null"`
	Date        time.Time `gorm:"not null;index"`
	Description string    `gorm:"type:text"`
	Version     int64     `gorm:"not null;default:1"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Transaction) TableName() string {
	return "transactions"
}
