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

type MessageResponse struct {
	Message string `json:"message,omitempty"`
}

const checkYourInbox = "check your inbox for further instructions"

var errInvalidCode = errors.New("error invalid code")

type ResetPasswordParams struct {
	Code        string `json:"code,omitempty"`
	Email       string `json:"email,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

// Reset Password
//
//	@Summary		Reset
//	@Description	Reset Password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/refresh [get]
//	@Param			user	query		bool					false	"If true, returns User in body"
//	@Param			cookies	query		bool					false	"If true, returns AccessToken and SessionID in body"
//	@Success		200		{object}	RefreshTokenResponse	"RefreshTokenResponse"
func (s *Server) Reset(w http.ResponseWriter, r *http.Request) {
	var param ResetPasswordParams
	if err := s.ParseAndValidate(r, &param, ValidationDisable); err != nil {
		s.Failure(w, err, http.StatusBadRequest)
		return
	}
	// code can also be passed in query params  ?code=xxxxxx
	if param.Code == "" {
		param.Code = r.URL.Query().Get("code")
	}
	// email can also be passed in query params  ?email=xxxxxx
	if param.Email == "" {
		param.Code = r.URL.Query().Get("email")
	}

	hasCode := len(param.Code) != 0

	email, err := vo.NewEmail(param.Email)
	if err != nil {
		s.logr.Error().Err(err).Msg("invalid email")
		s.Failure(w, err, http.StatusBadRequest)
		return
	}

	res := s.repo.UserGetByEmail(r.Context(), email)
	if res.Error != nil {
		if res.StatusCode != http.StatusNotFound {
			s.Failure(w, res.Error, res.StatusCode)
			return
		}
		s.logr.Error().Err(res.Error).Msg("reset failed: no user")
		response := MessageResponse{Message: checkYourInbox}
		s.Success(w, response)
		return
	}

	if !res.User.Verified {
		response := MessageResponse{Message: "can not proceed without verified email"}
		s.Success(w, response, http.StatusForbidden)
		return
	}

	isAlreadyPending := true
	var pending payload.ResetPassword
	actionKey := fmt.Sprintf("reset:%s", res.User.Email.Value())
	if err = s.cache.Get(r.Context(), actionKey, &pending); err != nil {
		if !errors.Is(err, cache.ErrNoRecord) {
			s.Failure(w, err)
			return
		}
		isAlreadyPending = false
	}

	if isAlreadyPending && !hasCode {
		timeDifference := time.Since(pending.CreatedAt)
		if timeDifference < s.config.AuthPasswordResetCooldown {
			tryAfter := s.config.AuthPasswordResetCooldown - timeDifference
			err = fmt.Errorf("recently requested, please try again after %s", tryAfter)
			s.Failure(w, err, http.StatusForbidden)
			return
		}
		s.Failure(w, errInvalidCode, http.StatusForbidden)
		return
	}

	if isAlreadyPending && hasCode {
		// we dont have to verify payload or anything else
		// in cache after expiration, cache automatically deletes it
		// so it exists means code is still valid, just have to match it
		if !email.IsEqual(pending.Email) {
			s.logr.Error().Msg("pending verification, mismatch email")
			s.Failure(w, errInvalidCode, http.StatusForbidden)
			return
		}
		if param.Code != pending.Code {
			s.logr.Error().Msg("pending verification, mismatch verification code")
			s.Failure(w, errInvalidCode, http.StatusForbidden)
			return
		}

		password, err := vo.NewPassword(param.NewPassword)
		if err != nil {
			s.Failure(w, err, http.StatusBadRequest)
			return
		}

		hashedPassword, err := password.HashPassword()
		if err != nil {
			s.Failure(w, err)
			return
		}

		// now update it to users repository
		res = s.repo.UserUpdatePassword(r.Context(), res.User.ID, hashedPassword)
		if res.Error != nil {
			s.Failure(w, res.Error, res.StatusCode)
			return
		}

		response := MessageResponse{Message: "password reset successfully"}
		s.Success(w, response)
		return
	}

	// here there is no pending email verification code
	if hasCode {
		// since there is no pending verification code to match and verify to
		// so we directly return error, as code might have been expired.
		s.Failure(w, errInvalidCode, http.StatusForbidden)
		return
	}

	// here we only send email verification code

	emailPayload := payload.NewResetPasswordPayload(res.User)
	err = s.tasks.ResetPassword(r.Context(), emailPayload,
		asynq.MaxRetry(2), asynq.Group(worker.PriorityLow),
		asynq.ProcessIn(time.Millisecond*time.Duration(utils.RandomInt(100, 600))),
	)
	if err != nil {
		s.Failure(w, err)
		return
	}
	message := MessageResponse{Message: checkYourInbox}
	s.Success(w, message)
}
