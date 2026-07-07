package auth

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	FindRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string) (*RefreshToken, error)
	DeleteRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteRefreshToken(ctx context.Context, tokenID uuid.UUID) error
}
