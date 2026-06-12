package mailer

import (
	"fmt"
	"notification-service/config"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendEmail(to, subject, body string) error
}

type mailer struct {
	config *config.Config
}

func NewMailer(cfg *config.Config) Mailer {
	return &mailer{config: cfg}
}

func (m *mailer) SendEmail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", m.config.SMTP.FromName, m.config.SMTP.From))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		m.config.SMTP.Host,
		m.config.SMTP.Port,
		m.config.SMTP.Username,
		m.config.SMTP.Password,
	)

	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
