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

type VerificationEmailPayload struct {
	Email string `json:"email"`
}

func NewVerificationEmailTask(email string) (*asynq.Task, error) {
	payload, err := sonic.Marshal(VerificationEmailPayload{Email: email})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task payload: %w", err)
	}
	return asynq.NewTask(TypeVerificationEmail, payload), nil
}

func HandleVerificationEmailTask(jwtSecret []byte, emailService *email.EmailService, registrationRoute string, baseURL string) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		var p VerificationEmailPayload
		if err := sonic.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("failed to unmarshal task payload: %w", err)
		}
		token, err := tokens.CreateEmailVerification(jwtSecret, p.Email, baseURL, baseURL)
		if err != nil {
			return err
		}
		return emailService.SendVerificationEmail(ctx, email.User{Email: p.Email}, fmt.Sprintf("%v?token=%v", registrationRoute, token))

	}

}
