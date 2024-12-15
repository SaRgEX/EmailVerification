package repository

import (
	"context"
	"email-verification-service/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Register(ctx context.Context, client model.ClientInput) error
	Verify(ctx context.Context, client model.ClientVerification) error
	IsVerified(ctx context.Context, client model.ClientVerification) (model.VerificationCode, error)
	Refresh(ctx context.Context, client model.ClientVerification) (model.VerificationCode, error)
}

type Repository struct {
	Client
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Client: NewClientPostgres(pool),
	}
}
