package category

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

func (r *repositoryGorm) Create(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *repositoryGorm) Update(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *repositoryGorm) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Category{}).Error
}

func (r *repositoryGorm) FindByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Category, error) {
	var category Category
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *repositoryGorm) List(ctx context.Context, userID uuid.UUID) ([]Category, error) {
	var categories []Category
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("name asc").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
