package transaction

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	WithTx(tx *gorm.DB) Repository
	Create(ctx context.Context, tx *Transaction) error
	Update(ctx context.Context, tx *Transaction) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Transaction, error)
	List(ctx context.Context, userID uuid.UUID) ([]Transaction, error)
}
