package app

import (
	"context"
	"email-verification-service/internal/handler"
	"email-verification-service/internal/pkg/config"
	"email-verification-service/internal/pkg/database"
	"email-verification-service/internal/pkg/logger"
	"email-verification-service/internal/pkg/server"
	"email-verification-service/internal/pkg/smtp"
	"email-verification-service/internal/repository"
	"email-verification-service/internal/service"
	"log/slog"
)

type App struct {
	server     *server.Server
	config     *config.Config
	database   *database.Database
	logger     *logger.Logger
	handler    *handler.Handler
	service    *service.Service
	repository *repository.Repository
	smtpServer *smtp.Smtp
}

func New() *App {
	a := &App{}
	a.config = config.InitConfig()
	a.logger = logger.New(a.config.Logger)
	a.smtpServer = smtp.New(a.config.SMTPServer)
	a.database = database.New(context.TODO(), a.config.Database)
	a.repository = repository.New(a.database.Pool)
	a.service = service.New(a.repository, a.smtpServer)
	a.handler = handler.New(a.service)
	a.server = server.New(a.config.HTTPServer, a.handler.InitRoutes())
	return a
}

func (a *App) Run() error {
	return a.server.Run()
}

func (a *App) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		slog.With("error", err).Error("Server shutdown failed")
		return err
	}
	slog.Debug("Server shutdown")
	a.database.Pool.Close()
	slog.Debug("Database connection closed")
	return nil
}
