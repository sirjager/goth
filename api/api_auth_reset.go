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
	"github.com/sirjager/goth/repository/users"
	"github.com/sirjager/goth/vo"
	"github.com/sirjager/goth/worker"
)

const checkYourInbox = "check your inbox for further instructions"

type ResetPasswordParams struct {
	Email       string `json:"email,omitempty"`
	Code        string `json:"code,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

// Reset Password
//
//	@Summary		Reset
//	@Description	Reset Password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/reset [post]
//	@Param			body	body	ResetPasswordParams	true	"ResetPasswordParams"
func (s *Server) Reset(w http.ResponseWriter, r *http.Request) {
	var param ResetPasswordParams
	if err := s.ParseAndValidate(r, &param, ValidationDisable); err != nil {
		s.Failure(w, err, http.StatusBadRequest)
		return
	}

	hasCode := len(param.Code) != 0
	emailAction := payload.PasswordReset
	cooldownTime := s.config.AuthPasswordResetCooldown
	codeExpiration := s.config.AuthPasswordResetExpire

	email, err := vo.NewEmail(param.Email)
	if err != nil {
		s.Failure(w, err, http.StatusBadRequest)
		return
	}

	res := s.repo.UserGetByEmail(r.Context(), email)
	if res.Error != nil {
		if res.StatusCode != http.StatusNotFound {
			s.Failure(w, res.Error, res.StatusCode)
			return
		}
		if hasCode {
			s.Failure(w, errInvalidCode, http.StatusForbidden)
			return
		}
		s.Success(w, MessageResponse{Message: checkYourInbox})
		return
	}

	if !res.User.Verified {
		s.Failure(w, errors.New("can not proceed without verified email"), http.StatusForbidden)
		return
	}

	// if request is already pending then return error with try again message
	isAlreadyPending := true
	var pending payload.EmailPayload
	actionKey := payload.EmailKey(res.User.Email.Value(), payload.PasswordReset)
	if err = s.cache.Get(r.Context(), actionKey, &pending); err != nil {
		if !errors.Is(err, cache.ErrNoRecord) {
			s.Failure(w, err)
			return
		}
		isAlreadyPending = false
	}

	if !isAlreadyPending {
		if hasCode {
			s.Failure(w, errInvalidCode, http.StatusForbidden)
			return
		}
		mailCode := utils.RandomNumberAsString(6)
		mailSub := "Password Reset Requested"
		payload := payload.EmailPayload{
			Email:     res.User.Email.Value(),
			Body:      resetPasswordEmailBody(mailCode, codeExpiration),
			Subject:   mailSub,
			Type:      emailAction,
			Code:      mailCode, // this will be used to validate code
			CacheKey:  payload.EmailKey(res.User.Email.Value(), emailAction),
			CacheExp:  codeExpiration,
			CreatedAt: time.Now(),
		}
		token, _, err := s.toknb.CreateToken(payload, codeExpiration)
		if err != nil {
			s.Failure(w, err)
			return
		}

		if err = s.tasks.SendEmail(r.Context(), worker.SendEmailParams{Token: token},
			asynq.MaxRetry(3), asynq.Group(worker.PriorityUrgent),
			asynq.ProcessIn(time.Millisecond*time.Duration(utils.RandomInt(1000, 5000))), // 1 to 5 seconds
		); err != nil {
			s.Failure(w, err)
			return
		}
		s.Success(w, MessageResponse{Message: checkYourInbox})
		return
	}

	// From here we have pending requests that needs to be completed
	if !hasCode {
		timeDifference := time.Since(pending.CreatedAt)
		if timeDifference < cooldownTime {
			tryAfter := cooldownTime - timeDifference
			err = fmt.Errorf("recently requested, please try again after %s", tryAfter)
			s.Failure(w, err, http.StatusBadRequest)
			return
		}
		// request is pending but no code provided so we reject with invalid code
		s.Failure(w, errInvalidCode, http.StatusForbidden)
		return
	}

	// Here we have pending request and code is also provided, we check and update
	if !email.IsEqual(pending.Email) {
		s.Failure(w, errInvalidCode, http.StatusForbidden)
		return
	}
	if param.Code != pending.Code {
		s.Failure(w, errInvalidCode, http.StatusForbidden)
		return
	}

	password, passErr := vo.NewPassword(param.NewPassword)
	if passErr != nil {
		s.Failure(w, passErr, http.StatusBadRequest)
		return
	}

	hashedPassword, hashErr := password.HashPassword()
	if hashErr != nil {
		s.Failure(w, hashErr)
		return
	}

	updateParams := users.UserUpdatePasswordTxParams{
		UserID:   res.User.ID,
		Password: hashedPassword,
		AfterUpdate: func() error {
			return s.cache.Delete(r.Context(), actionKey)
		},
	}

	res = s.repo.UserUpdatePasswordTx(r.Context(), updateParams)
	if res.Error != nil {
		s.Failure(w, res.Error, res.StatusCode)
		return
	}

	s.Success(w, MessageResponse{Message: "password reset successfully"})
}

func resetPasswordEmailBody(code string, validFor time.Duration) string {
	sb := &utils.StringBuilder{}
	sb.WriteLine("<p>Dear User,</p>")
	sb.WriteLine("<p>You have requested to reset your password.</p>")
	sb.WriteLine("<p>Please use the code below to reset your password:</p>")
	sb.WriteLine(fmt.Sprintf("<p><strong>Reset Code:</strong> %s</p>", code))
	sb.WriteLine(fmt.Sprintf("<p><em>This code is valid for the next %s.</em></p>", validFor))
	sb.WriteLine("<p>If you did not request this, please ignore this email.</p>")
	return sb.String()
}
