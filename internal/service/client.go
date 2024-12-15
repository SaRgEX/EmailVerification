package service

import (
	"context"
	"email-verification-service/internal/pkg/hash"
	"email-verification-service/internal/pkg/smtp"
	"email-verification-service/internal/repository"
	"email-verification-service/model"
	"email-verification-service/utils"
	"fmt"
	"log/slog"
)

type ClientService struct {
	r          *repository.Repository
	smtpServer *smtp.Smtp
}

func NewClientService(r *repository.Repository, smtpServer *smtp.Smtp) *ClientService {
	return &ClientService{
		r:          r,
		smtpServer: smtpServer,
	}
}

func (s *ClientService) Register(ctx context.Context, client model.ClientInput) error {
	var err error
	client.VerificationCode = utils.GenerateCode()
	client.Password, err = hash.HashPassword(client.Password)
	err = s.r.Register(ctx, client)
	if err != nil {
		slog.With("error", err).Error("Error while registering client")
		return err
	}
	err = s.smtpServer.SendVerificationEmail(client.Email, client.VerificationCode)
	if err != nil {
		slog.With("error", err).Error("Error while sending verification code")
		return err
	}
	return nil
}

func (s *ClientService) Verify(ctx context.Context, client model.ClientVerification) error {
	return s.r.Verify(ctx, client)
}

func (s *ClientService) Refresh(ctx context.Context, client model.ClientVerification) (model.VerificationCode, error) {
	client.VerificationCode.Code = utils.GenerateCode()
	vc, err := s.r.IsVerified(ctx, client)
	if err != nil {
		slog.With("error", err).Error("Error while checking if client is verified")
		return model.VerificationCode{}, err
	}
	if vc.IsVerified {
		return model.VerificationCode{}, fmt.Errorf("client is verified")
	}
	output, err := s.r.Refresh(ctx, client)
	if err != nil {
		slog.With("error", err).Error("Error while refreshing verification code")
		return model.VerificationCode{}, err
	}
	err = s.smtpServer.SendVerificationEmail(client.Email, client.Code)
	if err != nil {
		slog.With("error", err).Error("Error while sending verification code")
		return model.VerificationCode{}, err
	}
	return output, nil
}
