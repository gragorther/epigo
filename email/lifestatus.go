package email

import (
	"context"
	_ "embed"
	"text/template"
)

type LifeStatusUser struct {
	Name  string
	Email string
}

//go:embed templates/lifestatus.txt
var userLifeStatusTpl string

func (e *EmailService) SendUserLifeStatusEmail(ctx context.Context, user LifeStatusUser, verificationURL string) error {
	msg, err := e.newMsg("verify your life status", user.Email)
	if err != nil {
		return err
	}
	tpl, err := template.New("lifeStatus").Parse(userLifeStatusTpl)
	if err != nil {
		return err
	}

	templateData := struct {
		VerificationURL string
		UserName        string
		Email           string
	}{
		VerificationURL: verificationURL,
		UserName:        user.Name,
		Email:           user.Email,
	}

	textMsg, err := e.newTextMsg(msg, tpl, templateData)
	if err != nil {
		return err
	}
	return e.client.DialAndSendWithContext(ctx, textMsg)
}
