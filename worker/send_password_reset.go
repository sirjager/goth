package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/gopkg/mail"

	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/vo"
)

const TaskSendResetPassword = "task:send:resetpassword"

func (d *dist) ResetPassword(
	ctx context.Context,
	payload *payload.ResetPassword,
	opts ...asynq.Option,
) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TaskSendResetPassword, bytes, opts...)
	if _, err := d.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}
	d.logr.Info().Str("task", TaskSendResetPassword).Msg("task enqueued")
	return nil
}

func (p *proc) ResetPassword(ctx context.Context, task *asynq.Task) error {
	var _payload payload.ResetPassword
	if err := json.Unmarshal(task.Payload(), &_payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	if time.Since(_payload.CreatedAt) > p.config.AuthEmailVerifyExpire {
		return fmt.Errorf("email payload expired: %w", asynq.SkipRetry)
	}

	userEmail, err := vo.NewEmail(_payload.Email)
	if err != nil {
		return fmt.Errorf("invalid user email: %w", asynq.SkipRetry)
	}

	uniqueCacheKey := fmt.Sprintf("reset:%s", userEmail.Value())
	codeIsValidTill := _payload.CreatedAt.Add(p.config.AuthRefreshTokenExpire)
	email := mail.Mail{To: []string{userEmail.Value()}}
	email.Subject = "Password Reset Code"
	email.Body = fmt.Sprintf(`
	Password reset code: <b>%s</b><br>
	Code is valid till: <b>%s</b>`,
		_payload.Code, codeIsValidTill)

	if err = p.mail.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// redis will automatically delete expired verification codes
	if err = p.cache.Set(ctx, uniqueCacheKey, _payload, p.config.AuthPasswordResetExpire); err != nil {
		return fmt.Errorf("failed to cache payload: %w", err)
	}

	p.logr.Info().Str("task", TaskSendEmailVerification).Msg("task processed successfully")

	return nil
}
