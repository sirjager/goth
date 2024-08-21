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

const TaskSendEmailVerification = "task:sendEmailVerification"

func (d *dist) SendEmailVerification(
	ctx context.Context,
	payload *payload.VerifyEmail,
	opts ...asynq.Option,
) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling payload: %w", err)
	}
	task := asynq.NewTask(TaskSendEmailVerification, bytes, opts...)
	if _, err := d.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	d.logr.Info().Str("task", TaskSendEmailVerification).Msg("task enqueued")
	return nil
}

func (p *proc) SendEmailVerification(ctx context.Context, task *asynq.Task) (err error) {
	var _payload payload.VerifyEmail
	// if can't even unmarshal we will skip retring
	if err = json.Unmarshal(task.Payload(), &_payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	if time.Since(_payload.CreatedAt) > p.config.AuthEmailVerifyExpire {
		return fmt.Errorf("email verify payload expired: %w", asynq.SkipRetry)
	}

	userEmail, err := vo.NewEmail(_payload.Email)
	if err != nil {
		return fmt.Errorf("invalid user email: %w", asynq.SkipRetry)
	}

	uniqueCacheKey := fmt.Sprintf("verify:%s", userEmail.Value())
	codeExpireDuration := p.config.AuthEmailVerifyExpire

	email := mail.Mail{To: []string{userEmail.Value()}}
	email.Subject = "Complete Sign Up With Email Verification"
	email.Body = fmt.Sprintf(`
	Welcome to our community. <br><br>
	Complete your signup by verification <br>
	Your email verification code : <b>%s</b> <br>
	Verification code is only valid for : <b>%s</b> <br><br>`,
		_payload.Code,
		codeExpireDuration,
	)

	if err = p.mail.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// redis will automatically delete expired verification codes
	if err = p.cache.Set(ctx, uniqueCacheKey, _payload, codeExpireDuration); err != nil {
		return fmt.Errorf("failed to cache payload: %w", err)
	}

	p.logr.Info().Str("task", TaskSendEmailVerification).Msg("task processed successfully")
	return
}
