package account

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	WithTx(tx *gorm.DB) Repository
	Create(
		ctx context.Context,
		account *Account,
	) error

	Update(
		ctx context.Context,
		account *Account,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) error

	FindByID(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) (*Account, error)

	List(
		ctx context.Context,
		userID uuid.UUID,
	) ([]Account, error)

	IncreaseBalance(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
		amount int64,
	) error

	DecreaseBalance(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
		amount int64,
	) error

	Exists(
		ctx context.Context,
		id uuid.UUID,
		userID uuid.UUID,
	) (bool, error)
}
