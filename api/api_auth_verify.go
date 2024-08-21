package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/vo"
	"github.com/sirjager/goth/worker"
)

type EmailVerificationResponse struct {
	Message string `json:"message,omitempty"`
} //	@name	EmailVerificationResponse

// Verify Email
//
//	@Summary		Verify
//	@Description	Email Verification
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/verify [post]
//	@Param			email	query		string	true	"Email to verify"
//	@Param			code	query		string	false	"Email verification code if already have any"
//	@Success		200		{string}	string	"Success message"
func (s *Server) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	queryParamVerifyEmail := r.URL.Query().Get("email")
	queryParamVerifyCode := r.URL.Query().Get("code")
	hasVerificationCode := len(queryParamVerifyCode) != 0

	email, err := vo.NewEmail(queryParamVerifyEmail)
	if err != nil {
		s.logr.Error().Err(err).Msg("invalid email")
		s.Failure(w, err, http.StatusBadRequest)
		return
	}

	res := s.repo.UserGetByEmail(r.Context(), email)
	user := res.User

	if res.Error == nil && user.Verified {
		s.SuccessOK(w, "email already verified")
		return
	}

	if res.Error != nil {
		if res.StatusCode != http.StatusNotFound {
			s.Failure(w, res.Error)
			return
		}
		if !hasVerificationCode {
			s.logr.Error().Msg("user not found, no verification code")
			message := EmailVerificationResponse{Message: "check email for further instructions"}
			s.Success(w, message)
			return
		}

		s.logr.Error().Msg("user not found, has verification code")
		s.Failure(w, errInvalidEmailVerificationCode, http.StatusBadRequest)
		return
	}

	verificationCodeKey := fmt.Sprintf("verify:%s", user.Email.Value())

	isAlreadyPending := true
	var pending payload.VerifyEmail
	if err = s.cache.Get(r.Context(), verificationCodeKey, &pending); err != nil {
		if !errors.Is(err, cache.ErrNoRecord) {
			s.Failure(w, err)
			return
		}
		isAlreadyPending = false
	}

	if isAlreadyPending && !hasVerificationCode {
		timeDifference := time.Since(pending.CreatedAt)
		if timeDifference < s.config.AuthEmailVerifyCooldown {
			tryAfter := s.config.AuthEmailVerifyCooldown - timeDifference
			err = fmt.Errorf("recently requested, please try again after %s", tryAfter)
			s.Failure(w, err, http.StatusBadRequest)
			return
		}

		s.logr.Error().Msg("pending verification, no email verification code")
		// since no email verification code is provided
		s.Failure(w, errInvalidEmailVerificationCode, http.StatusBadRequest)
		return
	}

	if isAlreadyPending && hasVerificationCode {
		// we dont have to verify payload or anything else
		// in cache after expration, cache automatically deletes it
		// so it exists means code is still valid, just have to match it
		if !email.IsEqual(pending.Email) {
			s.logr.Error().Msg("pending verification, mismatch email")
			s.Failure(w, errInvalidEmailVerificationCode, http.StatusBadRequest)
			return
		}
		if queryParamVerifyCode != pending.Code {
			s.logr.Error().Msg("pending verification, mismatch verification code")
			s.Failure(w, errInvalidEmailVerificationCode, http.StatusBadRequest)
			return
		}

		// now update it to users repository
		res = s.repo.UserUpdateVerified(r.Context(), user.ID, true)
		if res.Error != nil {
			s.Failure(w, res.Error, res.StatusCode)
			return
		}

		response := EmailVerificationResponse{Message: "email successfully verified"}
		s.Success(w, response)
		return
	}

	// here there is no pending email verification code
	if hasVerificationCode {
		// since there is no pending verification code, to match and verify to
		// so we directly return error, as code might have been expired.
		s.Failure(w, errInvalidEmailVerificationCode, http.StatusBadRequest)
		return
	}

	// here we only send email verification code

	payload := payload.NewVerifyEmailPayload(user)
	err = s.tasks.SendEmailVerification(r.Context(), payload,
		asynq.MaxRetry(2), asynq.Group(worker.PriorityLow),
		asynq.ProcessIn(time.Millisecond*time.Duration(utils.RandomInt(100, 600))),
	)
	if err != nil {
		s.Failure(w, err)
		return
	}

	s.SuccessOK(w, "check your inbox for further instructions")
}
