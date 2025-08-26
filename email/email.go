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
