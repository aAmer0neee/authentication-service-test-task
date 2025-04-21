package notify

import (
	"net/smtp"

	"github.com/aAmer0neee/authentication-service-test-task/internal/config"
)

//go:generate mockgen -source=notify.go -destination=mocks/notify_mock.go -package=notify_mock
type Notifyer interface {
	SendMail(recipient, message string) error
}

type email struct {
	smtpHost string
	smtpPort string
	sender   string
	auth     smtp.Auth
}

func New(cfg *config.Cfg) Notifyer {

	return &email{
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

func (n *email) SendMail(recipient, message string) error {
	return smtp.SendMail(
		n.smtpHost+":"+n.smtpPort,
		n.auth,
		n.sender,
		[]string{recipient},
		[]byte("Subject: Warning Notification\r\n\r\n"+message))
}
