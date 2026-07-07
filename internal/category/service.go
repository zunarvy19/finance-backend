package category

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, req *CreateCategoryRequest) (*CategoryResponse, error)
	Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateCategoryRequest) (*CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*CategoryResponse, error)
	List(ctx context.Context, userID uuid.UUID) ([]CategoryResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, req *CreateCategoryRequest) (*CategoryResponse, error) {
	category := &Category{
		UserID: userID,
		Name:   req.Name,
		Type:   req.Type,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return toResponse(category), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	category.Name = req.Name
	category.Type = req.Type

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return toResponse(category), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*CategoryResponse, error) {
	category, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return toResponse(category), nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]CategoryResponse, error) {
	categories, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = *toResponse(&cat)
	}
	return responses, nil
}

func toResponse(category *Category) *CategoryResponse {
	return &CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Type: category.Type,
	}
}
