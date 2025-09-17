package email

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
	"text/template"

	"github.com/wneessen/go-mail"
)

type UserDeathEmailAndRecipients struct {
	Title           string
	Content         string
	RecipientEmails []string
}

//go:embed templates/death.txt
var deathTemplateString string

// name is the name of the person who died
func (e *EmailService) SendUserDeathEmails(ctx context.Context, name string, emails []UserDeathEmailAndRecipients) error {
	if emails == nil {
		return errors.New("nil emails")
	}
	type templateData struct {
		Email   string
		Name    string
		Message string
	}
	tpl, err := template.New("deathTemplate").Parse(deathTemplateString)
	if err != nil {
		return err
	}

	var messages []*mail.Msg
	for _, email := range emails {
		for _, recipient := range email.RecipientEmails {
			msg := mail.NewMsg()
			if err := msg.FromFormat(e.fromFormat, e.from); err != nil {
				return err
			}

			if err := msg.EnvelopeFrom(fmt.Sprintf("%s+%d", e.from, rand.Int31())); err != nil {
				return err
			}

			if err := msg.To(recipient); err != nil {
				return err
			}
			msg.SetDate()
			msg.SetMessageID()
			msg.Subject(fmt.Sprintf("Message from %s: %s", name, email.Title))
			if err := msg.SetBodyTextTemplate(tpl, templateData{
				Email:   recipient,
				Name:    name,
				Message: email.Content,
			}); err != nil {
				return err
			}
			messages = append(messages, msg)
		}
	}

	return e.client.DialAndSendWithContext(ctx, messages...)
}
