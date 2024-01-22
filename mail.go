package newsletter

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
)

// SMTPServer is the SMTP server of gmail
const SMTPServer = "smtp.gmail.com"

// EmailConfig contains the necessary information to authenticate in the SMTP server
type EmailConfig struct {
	Password string
	UserName string
}

type MailClient struct {
	cfg EmailConfig
}

func NewMailClient(cfg EmailConfig) *MailClient {
	return &MailClient{
		cfg: cfg,
	}
}

type Email interface {
	Send(dest []string, bodyMessage string) error
}

// Send sends an email to the given destination
func (m MailClient) Send(dest []string, bodyMessage string) error {
	auth := smtp.PlainAuth("", m.cfg.UserName, m.cfg.Password, SMTPServer)

	msg := []byte("To: " + dest[0] + "\r\n" +
		"Subject: Newsletter\r\n" +
		"\r\n" +
		bodyMessage + "\r\n")

	err := smtp.SendMail(SMTPServer+":587", auth, m.cfg.UserName, dest, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}

func EmailTrigger(ctx context.Context, s Storage, e Email) error {
	nl, err := s.Newsletter()
	if err != nil {
		return fmt.Errorf("error getting newsletter: %v", err)
	}

	for _, n := range nl {
		pages, err := s.PageIn(ctx, n.URLs)
		if err != nil {
			slog.Error("error getting pages", "error", err)
		}
		var validURLS []string
		for _, p := range pages {
			if p.IsMostRecent {
				validURLS = append(validURLS, p.URL)
			}
		}
		if len(validURLS) > 0 {
			err = e.Send([]string{n.UserEmail}, fmt.Sprintf("Hi %s, \n\nWe have found %d new articles for you: \n\n%s", n.UserEmail, len(validURLS), validURLS))
			if err != nil {
				slog.Error("error sending email", "error", err)
			}
		}
	}

	return nil
}
