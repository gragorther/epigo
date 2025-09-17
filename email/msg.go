package email

import (
	"fmt"
	"text/template"

	"github.com/wneessen/go-mail"
)

func (e *EmailService) newMsg(subject string, to string, opts ...mail.MsgOption) (*mail.Msg, error) {
	msg := mail.NewMsg(opts...)
	if err := msg.From(e.from); err != nil {
		return nil, err
	}
	if err := msg.To(to); err != nil {
		return nil, err
	}

	msg.Subject(subject)
	msg.SetDate()
	msg.SetMessageID()

	if err := msg.EnvelopeFrom(e.from); err != nil {
		return nil, err
	}
	return msg, nil
}

func (e *EmailService) newTextMsg(msg *mail.Msg, template *template.Template, templateData any, opts ...mail.PartOption) (*mail.Msg, error) {
	if err := msg.SetBodyTextTemplate(template, templateData, opts...); err != nil {
		return nil, fmt.Errorf("failed to set text template: %w", err)
	}
	return msg, nil
}
