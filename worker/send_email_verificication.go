package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/sirjager/gopkg/mail"

	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/vo"
)

const TaskSendEmailVerification = "task:sendEmailVerification"

type SendEmailVerificationPayload struct {
	Token string `json:"token,omitempty"`
}

func (d *distributor) SendEmailVerification(
	ctx context.Context,
	payload SendEmailVerificationPayload,
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

func (p *processor) SendEmailVerification(ctx context.Context, task *asynq.Task) (err error) {
	var payload SendEmailVerificationPayload

	// if can't even unmarshal we will skip retring
	if err = json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	var payloadData mw.TokenCustomPayload
	// NOTE: if token is invalid, we wont even retry the task
	_, err = p.tokens.VerifyToken(payload.Token, &payloadData)
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", asynq.SkipRetry)
	}

	userEmail, err := vo.NewEmail(payloadData.UserEmail)
	if err != nil {
		return fmt.Errorf("invalid user email: %w", asynq.SkipRetry)
	}

	_, err = vo.NewIDFrom(payloadData.UserID)
	if err != nil {
		return fmt.Errorf("invalid user id: %w", asynq.SkipRetry)
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
		payloadData.EmailVerificationCode,
		codeExpireDuration,
	)

	if err = p.mail.SendMail(email); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	// redis will automatically delete expired verification codes
	if err = p.cache.Set(ctx, uniqueCacheKey, payloadData, codeExpireDuration); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	p.logr.Info().Str("task", TaskSendEmailVerification).Msg("task processed successfully")
	return
}
