package account

import "github.com/google/uuid"

type CreateAccountRequest struct {
	Name     string `json:"name" validate:"required"`
	Currency string `json:"currency" validate:"required,len=3"`
}

type UpdateAccountRequest struct {
	Name     string `json:"name" validate:"required"`
	Currency string `json:"currency" validate:"required,len=3"`
}

type AccountResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Currency string    `json:"currency"`
	Balance  int64     `json:"balance"`
	Version  int64     `json:"version"`
}
