package dashboard

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zunarvy19/finance-backend/internal/transaction"
)

type Service interface {
	GetDashboard(ctx context.Context, userID uuid.UUID) (*DashboardResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetDashboard(ctx context.Context, userID uuid.UUID) (*DashboardResponse, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location()).Format(time.RFC3339)

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format(time.RFC3339)
	monthEnd := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location()).Format(time.RFC3339)

	totalBalance, err := s.repo.GetTotalBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	todayInc, todayExp, err := s.repo.GetIncomeExpenseByDateRange(ctx, userID, todayStart, todayEnd)
	if err != nil {
		return nil, err
	}

	monthInc, monthExp, err := s.repo.GetIncomeExpenseByDateRange(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	latestTx, err := s.repo.GetLatestTransactions(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	var latestTxRes []transaction.TransactionResponse
	for _, tx := range latestTx {
		latestTxRes = append(latestTxRes, transaction.TransactionResponse{
			ID:          tx.ID,
			AccountID:   tx.AccountID,
			CategoryID:  tx.CategoryID,
			Type:        tx.Type,
			Amount:      tx.Amount,
			Date:        tx.Date,
			Description: tx.Description,
			Version:     tx.Version,
		})
	}

	return &DashboardResponse{
		TotalBalance:       totalBalance,
		TodayIncome:        todayInc,
		TodayExpense:       todayExp,
		MonthlyIncome:      monthInc,
		MonthlyExpense:     monthExp,
		LatestTransactions: latestTxRes,
	}, nil
}
