package email

import (
	"auth-service/config"
	"context"
	"fmt"
	"github.com/wneessen/go-mail"
)

type Mailer struct {
	client *mail.Client
	from   string
}

func New(conf config.Config) (*Mailer, error) {
	client, err := mail.NewClient(
		conf.Email.SmtpHost,
		mail.WithPort(conf.Email.SmtpPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(conf.Email.FromMail),
		mail.WithPassword(conf.Email.MailPassword),
		mail.WithTimeout(conf.Email.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mail client: %w", err)
	}

	return &Mailer{
		client: client,
		from:   conf.Email.FromMail,
	}, nil
}

func (m *Mailer) Send(ctx context.Context, to, subject, body string) error {
	msg := mail.NewMsg()

	if err := msg.From(m.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := msg.To(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	if err := m.client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("failed to send mail: %w", err)
	}

	return nil
}
