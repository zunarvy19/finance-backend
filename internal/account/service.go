package account

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req *CreateAccountRequest) (*AccountResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateAccountRequest) (*AccountResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*AccountResponse, error)
	List(ctx context.Context, userID uuid.UUID) ([]AccountResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, req *CreateAccountRequest) (*AccountResponse, error) {
	account := &Account{
		UserID:   userID,
		Name:     req.Name,
		Currency: req.Currency,
		Balance:  0,
		Version:  1,
	}

	if err := s.repo.Create(ctx, account); err != nil {
		return nil, err
	}

	return toResponse(account), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateAccountRequest) (*AccountResponse, error) {
	account, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	account.Name = req.Name
	account.Currency = req.Currency

	if err := s.repo.Update(ctx, account); err != nil {
		return nil, err
	}

	return toResponse(account), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*AccountResponse, error) {
	account, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toResponse(account), nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]AccountResponse, error) {
	accounts, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = *toResponse(&account)
	}
	return responses, nil
}

func toResponse(account *Account) *AccountResponse {
	return &AccountResponse{
		ID:       account.ID,
		Name:     account.Name,
		Currency: account.Currency,
		Balance:  account.Balance,
		Version:  account.Version,
	}
}
