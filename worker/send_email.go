package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirjager/gopkg/mail"

	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/vo"
)

const TaskSendEmail = "task:sendemail"

type SendEmailParams struct {
	Token string `json:"token,omitempty"`
}

func (d *dist) SendEmail(ctx context.Context, p SendEmailParams, opts ...asynq.Option) error {
	bytes, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TaskSendEmail, bytes, opts...)
	if _, err := d.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}
	d.logr.Info().Str("task", TaskSendEmail).Msg("task enqueued")
	return nil
}

func (p *proc) SendEmail(ctx context.Context, task *asynq.Task) error {
	var _payload SendEmailParams
	if err := json.Unmarshal(task.Payload(), &_payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	var emailPayload payload.EmailPayload
	if _, err := p.toknb.VerifyToken(_payload.Token, &emailPayload); err != nil {
		return fmt.Errorf("failed to verify email token: %w", asynq.SkipRetry)
	}

	userEmail, err := vo.NewEmail(emailPayload.Email)
	if err != nil {
		return fmt.Errorf("invalid user email: %w", asynq.SkipRetry)
	}

	email := mail.Mail{To: []string{userEmail.Value()}}
	email.Subject = emailPayload.Subject
	email.Body = emailPayload.Body
	if err = p.mail.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if err = p.cache.Set(ctx, emailPayload.CacheKey, emailPayload, emailPayload.CacheExp); err != nil {
		return fmt.Errorf("failed to cache payload: %w", err)
	}

	p.logr.Info().
		Str("task", TaskSendEmail).
		Str("type", emailPayload.Type.String()).
		Msg("task processed successfully")

	return nil
}
