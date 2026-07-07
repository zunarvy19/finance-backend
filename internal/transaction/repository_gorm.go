package transaction

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zunarvy19/finance-backend/internal/common"
	"gorm.io/gorm"
)

type repositoryGorm struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryGorm{db: db}
}

func (r *repositoryGorm) WithTx(tx *gorm.DB) Repository {
	return &repositoryGorm{db: tx}
}

func (r *repositoryGorm) Create(ctx context.Context, tx *Transaction) error {
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return common.ErrConflict
		}
		return err
	}
	return nil
}

func (r *repositoryGorm) Update(ctx context.Context, tx *Transaction) error {
	res := r.db.WithContext(ctx).Model(tx).Where("version = ?", tx.Version-1).Updates(tx)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return common.ErrConflict
	}
	return nil
}

func (r *repositoryGorm) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Transaction{}).Error
}

func (r *repositoryGorm) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Transaction, error) {
	var tx Transaction
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&tx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &tx, nil
}

func (r *repositoryGorm) List(ctx context.Context, userID uuid.UUID) ([]Transaction, error) {
	var transactions []Transaction
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date desc").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
