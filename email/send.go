package email

import (
	"context"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/wneessen/go-mail"
)

//go:embed templates/verification.txt
var verificationTemplate string

type User struct {
	Name  string
	Email string
}

func (e *EmailService) SendVerificationEmail(ctx context.Context, user User, verificationLink string) error {
	message := mail.NewMsg()
	if err := message.From(e.from); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}
	if err := message.To(user.Email); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}
	message.Subject("Verify your email address")
	if err := message.From(e.from); err != nil {
		return fmt.Errorf("failed to set message from: %w", err)
	}
	if err := message.EnvelopeFrom(e.from); err != nil {
		return fmt.Errorf("failed to set envelope from: %w", err)
	}

	tpl, err := template.New("verification").Parse(verificationTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse verification email text template: %w", err)
	}
	type templateData struct {
		Name             string
		VerificationLink string
	}
	if err := message.SetBodyTextTemplate(tpl, templateData{Name: user.Name, VerificationLink: verificationLink}); err != nil {
		return fmt.Errorf("failed to set body text template: %w", err)
	}
	return e.client.DialAndSendWithContext(ctx, message)

}
