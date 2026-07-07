package dashboard

import (
	"context"

	"github.com/google/uuid"
	"github.com/zunarvy19/finance-backend/internal/transaction"
	"gorm.io/gorm"
)

type repositoryGorm struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryGorm{db: db}
}

func (r *repositoryGorm) GetTotalBalance(ctx context.Context, userID uuid.UUID) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Table("accounts").
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(balance), 0)").Scan(&total).Error
	return total, err
}

func (r *repositoryGorm) GetIncomeExpenseByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) (int64, int64, error) {
	var result struct {
		Income  int64
		Expense int64
	}

	err := r.db.WithContext(ctx).Table("transactions").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Select(`
			COALESCE(SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END), 0) as expense
		`).Scan(&result).Error

	return result.Income, result.Expense, err
}

func (r *repositoryGorm) GetLatestTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]transaction.Transaction, error) {
	var transactions []transaction.Transaction
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("date desc").Limit(limit).Find(&transactions).Error
	return transactions, err
}
