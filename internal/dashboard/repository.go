package dashboard

import (
	"context"

	"github.com/google/uuid"
	"github.com/zunarvy19/finance-backend/internal/transaction"
)

type Repository interface {
	GetTotalBalance(ctx context.Context, userID uuid.UUID) (int64, error)
	GetIncomeExpenseByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate string) (income int64, expense int64, err error)
	GetLatestTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]transaction.Transaction, error)
}
