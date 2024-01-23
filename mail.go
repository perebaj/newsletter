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
	Username string
}

// MailClient is the client that sends emails
type MailClient struct {
	cfg EmailConfig
}

// NewMailClient creates a new MailClient
func NewMailClient(cfg EmailConfig) *MailClient {
	return &MailClient{
		cfg: cfg,
	}
}

// Email is the interface that wraps the methods needed to deal with emails
type Email interface {
	Send(dest []string, bodyMessage string) error
}

// Send sends an email to the given destination
func (m MailClient) Send(dest []string, bodyMessage string) error {
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, SMTPServer)

	msg := []byte("To: " + dest[0] + "\r\n" +
		"Subject: Newsletter\r\n" +
		"\r\n" +
		bodyMessage + "\r\n")

	err := smtp.SendMail(SMTPServer+":587", auth, m.cfg.Username, dest, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}

// EmailTrigger get all newsletters and send an email to the user with new articles were found
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
