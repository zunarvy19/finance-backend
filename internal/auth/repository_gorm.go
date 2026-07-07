package auth

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

func (r *repositoryGorm) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return common.ErrConflict
		}
		return err
	}
	return nil
}

func (r *repositoryGorm) FindUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repositoryGorm) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repositoryGorm) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *repositoryGorm) FindRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string) (*RefreshToken, error) {
	var token RefreshToken
	if err := r.db.WithContext(ctx).First(&token, "user_id = ? AND token_hash = ?", userID, tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &token, nil
}

func (r *repositoryGorm) DeleteRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}

func (r *repositoryGorm) DeleteRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", tokenID).Delete(&RefreshToken{}).Error
}
