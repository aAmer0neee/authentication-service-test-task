package notify

import (
	"net/smtp"

	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
)

type Notifyer struct {
	smtpHost string
	smtpPort string
	sender   string
	auth     smtp.Auth
}

func New(cfg *config.Cfg) *Notifyer {

	return &Notifyer{
		smtpHost: cfg.Notifyer.SmtpHost,
		smtpPort: cfg.Notifyer.SmtpPort,
		sender:   cfg.Notifyer.Email,
		auth: smtp.PlainAuth(
			"",
			cfg.Notifyer.Email,
			cfg.Notifyer.Password,
			cfg.Notifyer.SmtpHost),
	}
}

func (n *Notifyer) SendMail(recipient, message string) error {
	return smtp.SendMail(
		n.smtpHost+":"+n.smtpPort,
		n.auth,
		n.sender,
		[]string{recipient},
		[]byte("Subject: Warning Notification\r\n\r\n" + message))
}
