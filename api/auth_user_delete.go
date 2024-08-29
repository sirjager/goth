package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/gopkg/cache"
	"github.com/sirjager/gopkg/httpx"
	"github.com/sirjager/gopkg/utils"

	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/payload"
	"github.com/sirjager/goth/repository/users"
	"github.com/sirjager/goth/worker"
)

// Authenticated User Delete
//
//	@Summary		User Delete
//	@Description	Initiate Authenticated User Deletion
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/api/auth/user [delete]
//	@Param			code	query	string	false	"code if already have"
func (s *Server) authUserDelete(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)

	codeQueryParam := r.URL.Query().Get("code")
	emailAction := payload.UserDeletion
	cooldownTime := s.Config().AuthUserDeleteCooldown
	codeExpiration := s.Config().AuthUserDeleteExpire
	hasCode := len(codeQueryParam) != 0

	if !user.Verified {
		httpx.Error(w, errors.New("can not proceed without verified email"), http.StatusForbidden)
		return
	}

	// if request is already pending then return error with try again message
	isAlreadyPending := true
	var pending payload.EmailPayload
	actionKey := payload.EmailKey(user.Email.Value(), emailAction)
	if err := s.Cache().Get(r.Context(), actionKey, &pending); err != nil {
		if !errors.Is(err, cache.ErrNoRecord) {
			httpx.Error(w, err)
			return
		}
		isAlreadyPending = false
	}

	if !isAlreadyPending {
		if hasCode {
			httpx.Error(w, errInvalidCode, http.StatusForbidden)
			return
		}
		mailCode := utils.RandomNumberAsString(32)
		mailSub := "Account Deletion Requested"
		payload := payload.EmailPayload{
			Email:     user.Email.Value(),
			Body:      userDeleteEmailBody(mailCode, user.Email.Value(), codeExpiration),
			Subject:   mailSub,
			Type:      emailAction,
			Code:      mailCode, // this will be used to validate code
			CacheKey:  payload.EmailKey(user.Email.Value(), emailAction),
			CacheExp:  codeExpiration,
			CreatedAt: time.Now(),
		}
		token, _, err := s.Tokens().CreateToken(payload, codeExpiration)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		if err = s.Tasks().SendEmail(r.Context(), worker.SendEmailParams{Token: token},
			asynq.MaxRetry(2), asynq.Group(worker.PriorityLazy),
			asynq.ProcessIn(time.Millisecond*time.Duration(utils.RandomInt(3000, 6000))), // 1 to 5 seconds
		); err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.Success(w, MessageResponse{Message: checkYourInbox})
		return
	}

	// From here we have pending requests that needs to be completed
	if !hasCode {
		timeDifference := time.Since(pending.CreatedAt)
		if timeDifference < cooldownTime {
			tryAfter := cooldownTime - timeDifference
			err := fmt.Errorf("recently requested, please try again after %s", tryAfter)
			httpx.Error(w, err, http.StatusBadRequest)
			return
		}
		// request is pending but no code provided so we reject with invalid code
		httpx.Error(w, errInvalidCode, http.StatusForbidden)
		return
	}

	// Here we have pending request and code is also provided, we check and update
	if !user.Email.IsEqual(pending.Email) {
		httpx.Error(w, errInvalidCode, http.StatusForbidden)
		return
	}

	if codeQueryParam == pending.Code && len(codeQueryParam) != 0 && len(pending.Code) != 0 {
		deleteParams := users.UserDeleteTxParams{
			UserID: user.ID, AfterUpdate: func() error {
				return s.Cache().Delete(r.Context(), actionKey)
			},
		}
		res := s.Repo().UserDeleteTx(r.Context(), deleteParams)
		if res.Error != nil {
			httpx.Error(w, res.Error, res.StatusCode)
			return
		}

		httpx.Success(w, MessageResponse{Message: "user successfully deleted"})
		return
	}

	httpx.Error(w, errInvalidCode, http.StatusForbidden)
}

func userDeleteEmailBody(code, email string, validFor time.Duration) string {
	sb := &utils.StringBuilder{}
	sb.WriteLine("<p>Dear User,</p>")
	sb.WriteLine("<p>You have requested to delete your account.</p>")
	sb.WriteLine("<p>Please use the code below to delete account:</p>")
	sb.WriteLine(fmt.Sprintf("<p><strong>Deletion Code:</strong> %s</p>", code))
	sb.WriteLine(fmt.Sprintf("<p><strong>Email Associated:</strong> %s</p>", email))
	sb.WriteLine(fmt.Sprintf("<p><em>This code is valid for the next %s.</em></p>", validFor))
	sb.WriteLine("<p>If you did not request this, please ignore this email.</p>")
	return sb.String()
}
