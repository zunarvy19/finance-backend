package transaction

import (
	"time"

	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	AccountID   uuid.UUID `json:"account_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	ClientID    uuid.UUID `json:"client_id" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof=INCOME EXPENSE"`
	Amount      int64     `json:"amount" validate:"required,gt=0"`
	Date        time.Time `json:"date" validate:"required"`
	Description string    `json:"description"`
}

type UpdateTransactionRequest struct {
	AccountID   uuid.UUID `json:"account_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof=INCOME EXPENSE"`
	Amount      int64     `json:"amount" validate:"required,gt=0"`
	Date        time.Time `json:"date" validate:"required"`
	Description string    `json:"description"`
	Version     int64     `json:"version" validate:"required,gt=0"`
}

type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	AccountID   uuid.UUID `json:"account_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Type        string    `json:"type"`
	Amount      int64     `json:"amount"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Version     int64     `json:"version"`
}
