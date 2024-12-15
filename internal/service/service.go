package service

import (
	"context"
	"email-verification-service/internal/pkg/smtp"
	"email-verification-service/internal/repository"
	"email-verification-service/model"
)

type Client interface {
	Register(ctx context.Context, client model.ClientInput) error
	Verify(ctx context.Context, client model.ClientVerification) error
	Refresh(ctx context.Context, client model.ClientVerification) (model.VerificationCode, error)
}

type Service struct {
	Client
}

func New(r *repository.Repository, smtpServer *smtp.Smtp) *Service {
	return &Service{
		Client: NewClientService(r, smtpServer),
	}
}
