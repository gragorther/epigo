package email

import (
	"github.com/wneessen/go-mail"
)

type EmailService struct {
	client *mail.Client
	from   string
}

/*from is the address the email service is sending emails from*/
func NewEmailService(client *mail.Client, from string) *EmailService {
	return &EmailService{client: client, from: from}
}

func NewClient(host string, port int, password string, username string) (*mail.Client, error) {
	return mail.NewClient(host,
		mail.WithPort(port),
		mail.WithPassword(password),
		mail.WithUsername(username), mail.WithTLSPortPolicy(mail.TLSOpportunistic), mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover))
}
