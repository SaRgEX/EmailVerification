package smtp

import (
	"email-verification-service/internal/pkg/config"
	"fmt"
	"net/smtp"
)

type Smtp struct {
	Host     string
	Port     string
	Password string
	From     string
}

func New(cfg config.SMTPServer) *Smtp {
	return &Smtp{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Password: cfg.Password,
		From:     cfg.From,
	}
}

func (s *Smtp) SendVerificationEmail(to, code string) error {
	body := fmt.Sprintf("Subject: Verify\nYour verification code is: %s", code)
	return s.SendEmail(to, body)
}

func (s *Smtp) SendEmail(to, message string) error {
	auth := smtp.PlainAuth("", s.From, s.Password, s.Host)
	msg := []byte(message)
	return smtp.SendMail(s.Host+":"+s.Port, auth, s.From, []string{to}, msg)
}
