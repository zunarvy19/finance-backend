package transaction

import (
	"context"

	"github.com/google/uuid"
	"github.com/zunarvy19/finance-backend/internal/common"
)

type LedgerService interface {
	CreateTransaction(ctx context.Context, tx *Transaction) error
	UpdateTransaction(ctx context.Context, oldTx *Transaction, newTx *Transaction) error
	DeleteTransaction(ctx context.Context, tx *Transaction) error
}

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req *CreateTransactionRequest) (*TransactionResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateTransactionRequest) (*TransactionResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*TransactionResponse, error)
	List(ctx context.Context, userID uuid.UUID) ([]TransactionResponse, error)
}

type service struct {
	repo          Repository
	ledgerService LedgerService
}

func NewService(repo Repository, ledgerService LedgerService) Service {
	return &service{
		repo:          repo,
		ledgerService: ledgerService,
	}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, req *CreateTransactionRequest) (*TransactionResponse, error) {
	tx := &Transaction{
		UserID:      userID,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		ClientID:    req.ClientID,
		Type:        req.Type,
		Amount:      req.Amount,
		Date:        req.Date,
		Description: req.Description,
		Version:     1,
	}

	if err := s.ledgerService.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return toResponse(tx), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateTransactionRequest) (*TransactionResponse, error) {
	oldTx, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	if oldTx.Version != req.Version {
		return nil, common.ErrConflict
	}

	newTx := *oldTx
	newTx.AccountID = req.AccountID
	newTx.CategoryID = req.CategoryID
	newTx.Type = req.Type
	newTx.Amount = req.Amount
	newTx.Date = req.Date
	newTx.Description = req.Description

	if err := s.ledgerService.UpdateTransaction(ctx, oldTx, &newTx); err != nil {
		return nil, err
	}

	return toResponse(&newTx), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	tx, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	return s.ledgerService.DeleteTransaction(ctx, tx)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*TransactionResponse, error) {
	tx, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toResponse(tx), nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]TransactionResponse, error) {
	transactions, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = *toResponse(&tx)
	}
	return responses, nil
}

func toResponse(tx *Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:          tx.ID,
		AccountID:   tx.AccountID,
		CategoryID:  tx.CategoryID,
		Type:        tx.Type,
		Amount:      tx.Amount,
		Date:        tx.Date,
		Description: tx.Description,
		Version:     tx.Version,
	}
}
