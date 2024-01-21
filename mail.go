package newsletter

import (
	"fmt"
	"net/smtp"
)

// EmailConfig contains the necessary information to authenticate in the SMTP server
type EmailConfig struct {
	Password string
	UserName string
}

// SMTPServer is the SMTP server of gmail
const SMTPServer = "smtp.gmail.com"

// Send sends an email to the given destination
func Send(dest []string, bodyMessage string, cfg EmailConfig) error {
	auth := smtp.PlainAuth("", cfg.UserName, cfg.Password, SMTPServer)

	msg := []byte("To: " + dest[0] + "\r\n" +
		"Subject: Newsletter\r\n" +
		"\r\n" +
		bodyMessage + "\r\n")

	err := smtp.SendMail(SMTPServer+":587", auth, cfg.UserName, dest, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}
