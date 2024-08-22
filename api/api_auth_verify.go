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
//	@Router			/auth/verify [get]
//	@Param			email	query		string	true	"Email to verify"
//	@Param			code	query		string	false	"Email verification code if already have any"
//	@Success		200		{string}	string	"Success message"
func (s *Server) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	emailQueryParam := r.URL.Query().Get("email")
	codeQueryParam := r.URL.Query().Get("code")

	hasCode := len(codeQueryParam) != 0
	emailAction := payload.EmailVerification
	cooldownTime := s.config.AuthEmailVerifyCooldown
	codeExpiration := s.config.AuthEmailVerifyExpire

	email, err := vo.NewEmail(emailQueryParam)
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
		if hasCode {
			s.Failure(w, errInvalidCode, http.StatusForbidden)
			return
		}
		s.Success(w, MessageResponse{Message: checkYourInbox})
		return
	}

	if res.User.Verified {
		s.Success(w, MessageResponse{Message: "email already verified"})
		return
	}

	// if request is already pending then return error with try again message
	isAlreadyPending := true
	var pending payload.EmailPayload
	actionKey := payload.EmailKey(res.User.Email.Value(), emailAction)
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
		mailSub := "Email Verification Requested"
		payload := payload.EmailPayload{
			Email:     res.User.Email.Value(),
			Body:      emailVerificationBody(mailCode, codeExpiration),
			Subject:   mailSub,
			Type:      emailAction,
			Code:      mailCode, // this will be used to validate code
			CacheKey:  payload.EmailKey(res.User.Email.Value(), emailAction),
			CacheExp:  codeExpiration,
			CreatedAt: time.Now(),
		}
		token, _, tokenErr := s.toknb.CreateToken(payload, codeExpiration)
		if tokenErr != nil {
			s.Failure(w, tokenErr)
			return
		}

		if err = s.tasks.SendEmail(r.Context(), worker.SendEmailParams{Token: token},
			asynq.MaxRetry(3), asynq.Group(worker.PriorityLow),
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
	if codeQueryParam != pending.Code {
		s.Failure(w, errInvalidCode, http.StatusForbidden)
		return
	}

	// now update it to users repository
	res = s.repo.UserUpdateVerified(r.Context(), res.User.ID, true)
	if res.Error != nil {
		s.Failure(w, res.Error, res.StatusCode)
		return
	}

	response := EmailVerificationResponse{Message: "email successfully verified"}
	s.Success(w, response)
}

func emailVerificationBody(code string, validFor time.Duration) string {
	sb := &utils.StringBuilder{}
	sb.WriteLine("<p>Dear User,</p>")
	sb.WriteLine("<p>Please use the code to verify your email:</p>")
	sb.WriteLine(fmt.Sprintf("<p><strong>Code:</strong> %s</p>", code))
	sb.WriteLine(fmt.Sprintf("<p><em>This code is valid for the next %s.</em></p>", validFor))
	sb.WriteLine("<p>If you did not request this, please ignore this email.</p>")
	return sb.String()
}
