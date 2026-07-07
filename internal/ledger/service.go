package ledger

import (
	"context"
	"time"

	"github.com/zunarvy19/finance-backend/internal/account"
	"github.com/zunarvy19/finance-backend/internal/transaction"
	"gorm.io/gorm"
)

type Service interface {
	CreateTransaction(ctx context.Context, tx *transaction.Transaction) error
	UpdateTransaction(ctx context.Context, oldTx *transaction.Transaction, newTx *transaction.Transaction) error
	DeleteTransaction(ctx context.Context, tx *transaction.Transaction) error
}

type service struct {
	db          *gorm.DB
	accountRepo account.Repository
	txRepo      transaction.Repository
}

func NewService(db *gorm.DB, accountRepo account.Repository, txRepo transaction.Repository) Service {
	return &service{
		db:          db,
		accountRepo: accountRepo,
		txRepo:      txRepo,
	}
}

func (s *service) CreateTransaction(ctx context.Context, tx *transaction.Transaction) error {
	return s.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		// Use Tx bound repositories
		accRepoTx := s.accountRepo.WithTx(dbTx)
		txRepoTx := s.txRepo.WithTx(dbTx)

		// 1. Insert Transaction
		if err := txRepoTx.Create(ctx, tx); err != nil {
			return err
		}

		// 2. Update Account Balance
		if tx.Type == "INCOME" {
			if err := accRepoTx.IncreaseBalance(ctx, tx.AccountID, tx.UserID, tx.Amount); err != nil {
				return err
			}
		} else if tx.Type == "EXPENSE" {
			if err := accRepoTx.DecreaseBalance(ctx, tx.AccountID, tx.UserID, tx.Amount); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *service) UpdateTransaction(ctx context.Context, oldTx *transaction.Transaction, newTx *transaction.Transaction) error {
	return s.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		accRepoTx := s.accountRepo.WithTx(dbTx)
		txRepoTx := s.txRepo.WithTx(dbTx)

		// Calculate the delta
		// To safely update, first we revert the old effect
		if oldTx.Type == "INCOME" {
			if err := accRepoTx.DecreaseBalance(ctx, oldTx.AccountID, oldTx.UserID, oldTx.Amount); err != nil {
				return err
			}
		} else if oldTx.Type == "EXPENSE" {
			if err := accRepoTx.IncreaseBalance(ctx, oldTx.AccountID, oldTx.UserID, oldTx.Amount); err != nil {
				return err
			}
		}

		// Increment version for optimistic locking
		newTx.Version = oldTx.Version + 1
		newTx.UpdatedAt = time.Now()

		// Update Transaction
		if err := txRepoTx.Update(ctx, newTx); err != nil {
			return err
		}

		// Apply the new effect
		if newTx.Type == "INCOME" {
			if err := accRepoTx.IncreaseBalance(ctx, newTx.AccountID, newTx.UserID, newTx.Amount); err != nil {
				return err
			}
		} else if newTx.Type == "EXPENSE" {
			if err := accRepoTx.DecreaseBalance(ctx, newTx.AccountID, newTx.UserID, newTx.Amount); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *service) DeleteTransaction(ctx context.Context, tx *transaction.Transaction) error {
	return s.db.WithContext(ctx).Transaction(func(dbTx *gorm.DB) error {
		accRepoTx := s.accountRepo.WithTx(dbTx)
		txRepoTx := s.txRepo.WithTx(dbTx)

		// Delete Transaction
		if err := txRepoTx.Delete(ctx, tx.ID, tx.UserID); err != nil {
			return err
		}

		// Revert Account Balance
		if tx.Type == "INCOME" {
			if err := accRepoTx.DecreaseBalance(ctx, tx.AccountID, tx.UserID, tx.Amount); err != nil {
				return err
			}
		} else if tx.Type == "EXPENSE" {
			if err := accRepoTx.IncreaseBalance(ctx, tx.AccountID, tx.UserID, tx.Amount); err != nil {
				return err
			}
		}

		return nil
	})
}
