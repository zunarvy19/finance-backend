package account

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) WithTx(tx *gorm.DB) Repository {
	return &repository{
		db: tx,
	}
}

// Create
func (r *repository) Create(
	ctx context.Context,
	account *Account,
) error {

	return r.db.
		WithContext(ctx).
		Create(account).
		Error

}

// Update
func (r *repository) Update(
	ctx context.Context,
	account *Account,
) error {

	return r.db.
		WithContext(ctx).
		Save(account).
		Error

}

// Delete
func (r *repository) Delete(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) error {

	return r.db.
		WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&Account{}).
		Error

}

// FindByID
func (r *repository) FindByID(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (*Account, error) {

	var account Account

	err := r.db.
		WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&account).
		Error

	if err != nil {
		return nil, err
	}

	return &account, nil

}

// List
func (r *repository) List(
	ctx context.Context,
	userID uuid.UUID,
) ([]Account, error) {

	var accounts []Account

	err := r.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at asc").
		Find(&accounts).
		Error

	if err != nil {
		return nil, err
	}

	return accounts, nil

}

// Exists
func (r *repository) Exists(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (bool, error) {

	var exists bool

	err := r.db.WithContext(ctx).
		Model(&Account{}).
		Where("id = ? AND user_id = ?", id, userID).
		Select("count(*) > 0").
		Scan(&exists).
		Error

	return exists, err

}

// IncreaseBalance
func (r *repository) IncreaseBalance(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	amount int64,
) error {

	return r.db.WithContext(ctx).
		Model(&Account{}).
		Where("id = ? AND user_id = ?", id, userID).
		UpdateColumn("balance", gorm.Expr("balance + ?", amount)).
		Error

}

// DecreaseBalance
func (r *repository) DecreaseBalance(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	amount int64,
) error {

	return r.db.WithContext(ctx).
		Model(&Account{}).
		Where("id = ? AND user_id = ?", id, userID).
		UpdateColumn("balance", gorm.Expr("balance - ?", amount)).
		Error

}
