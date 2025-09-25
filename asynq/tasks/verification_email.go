package tasks

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/gragorther/epigo/email"
	"github.com/gragorther/epigo/tokens"
	"github.com/hibiken/asynq"
)

const TypeVerificationEmail = "email:verification"

type verificationEmailPayload struct {
	Email string `json:"email"`
}

func (q *queue) SendVerificationEmail(email string) error {
	payload, err := sonic.Marshal(verificationEmailPayload{Email: email})
	if err != nil {
		return err
	}
	task := asynq.NewTask(TypeVerificationEmail, payload)
	_, err = q.enqueueTask(task)
	return err
}

func HandleVerificationEmailTask(createEmailVerification tokens.CreateEmailVerificationFunc,
	unmarshal UnmarshalFunc,
	emailService verificationEmailSender, registrationRoute string,
) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p verificationEmailPayload
		if err := unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("failed to unmarshal task payload: %w", err)
		}
		token, err := createEmailVerification(p.Email)
		if err != nil {
			return err
		}
		return emailService.SendVerificationEmail(ctx, email.User{Email: p.Email}, fmt.Sprintf("%v?token=%v", registrationRoute, token))
	}
}
