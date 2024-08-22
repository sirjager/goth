package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirjager/gopkg/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sirjager/goth/entity"
	repoerrors "github.com/sirjager/goth/repository/errors"
	mockRepo "github.com/sirjager/goth/repository/mock"
	"github.com/sirjager/goth/repository/users"
	"github.com/sirjager/goth/vo"
	mockTask "github.com/sirjager/goth/worker/mock"
)

func TestSignup(t *testing.T) {
	user, password, hashedPassword := randomUser(t)
	testCases := []struct {
		check func(t *testing.T, recorder *httptest.ResponseRecorder)
		stubs func(repo *mockRepo.MockRepo)
		body  SignUpRequestParams
		name  string
	}{
		{
			name: "OK",
			body: SignUpRequestParams{
				Username: user.User.Username,
				Email:    user.User.Email,
				Password: password.Value(),
			},
			stubs: func(repo *mockRepo.MockRepo) {
				newUser := &entity.User{
					ID:       vo.MustParseID(user.User.ID),
					Email:    vo.MustParseEmail(user.User.Email),
					Username: vo.MustParseUsername(user.User.Username),
					Password: hashedPassword,
					Provider: "credentials",
				}

				masterRes := users.UserReadResult{
					User:       nil,
					StatusCode: http.StatusNotFound,
					Error:      repoerrors.ErrUserNotFound,
				}
				res := users.UserReadResult{
					User:       newUser,
					StatusCode: http.StatusCreated,
					Error:      nil,
				}

				repo.EXPECT().UserGetMaster(gomock.Any()).Times(1).Return(masterRes)
				repo.EXPECT().
					UserCreate(gomock.Any(), EqCreateUserParams(newUser, password)).
					Times(1).
					Return(res)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "InvalidParams",
			body: SignUpRequestParams{
				Email:    user.User.Email,
				Password: password.Value(),
			},
			stubs: func(repo *mockRepo.MockRepo) {
				repo.EXPECT().UserGetMaster(gomock.Any()).Times(0)
				repo.EXPECT().UserCreate(gomock.Any(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: SignUpRequestParams{
				Email:    "galat-email-daal-deta-hu-kya-he-pata-$chalega@gmail.com",
				Username: user.User.Username,
				Password: password.Value(),
			},
			stubs: func(repo *mockRepo.MockRepo) {
				repo.EXPECT().UserGetMaster(gomock.Any()).Times(0)
				repo.EXPECT().UserCreate(gomock.Any(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: SignUpRequestParams{
				Email:    user.User.Email,
				Username: "cool-username.me",
				Password: password.Value(),
			},
			stubs: func(repo *mockRepo.MockRepo) {
				repo.EXPECT().UserGetMaster(gomock.Any()).Times(0)
				repo.EXPECT().UserCreate(gomock.Any(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: SignUpRequestParams{
				Email:    user.User.Email,
				Username: user.User.Username,
				Password: "missing-1-symbol-and-uppercase",
			},
			stubs: func(repo *mockRepo.MockRepo) {
				repo.EXPECT().UserGetMaster(gomock.Any()).Times(0)
				repo.EXPECT().UserCreate(gomock.Any(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "AlreadyExists",
			body: SignUpRequestParams{
				Email:    user.User.Email,
				Username: user.User.Username,
				Password: password.Value(),
			},
			stubs: func(repo *mockRepo.MockRepo) {
				masterRes := users.UserReadResult{
					Error:      repoerrors.ErrUserNotFound,
					User:       nil,
					StatusCode: http.StatusNotFound,
				}
				repo.EXPECT().UserGetMaster(gomock.Any()).Times(1).Return(masterRes)
				repo.EXPECT().UserCreate(gomock.Any(), gomock.Any()).Times(1).Return(
					users.UserReadResult{
						User:       nil,
						StatusCode: http.StatusConflict,
						Error:      repoerrors.ErrUserAlreadyExists,
					})
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: SignUpRequestParams{
				Email:    user.User.Email,
				Username: user.User.Username,
				Password: password.Value(),
			},
			stubs: func(repo *mockRepo.MockRepo) {
				repo.EXPECT().UserGetMaster(gomock.Any()).Times(1).Return(users.UserReadResult{
					User:       nil,
					StatusCode: http.StatusInternalServerError,
					Error:      pgx.ErrTooManyRows,
				})
				repo.EXPECT().UserCreate(gomock.Any(), gomock.Any()).Times(0)
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockRepo.NewMockRepo(ctrl)
			tc.stubs(repo)

			testTasks := mockTask.NewMockTaskDistributor(ctrl)

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			server := NewServer(repo, testLogr, testConfig, testCache, testTokens, testTasks)
			recoder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.check(t, recoder)
		})
	}
}

func randomUser(t *testing.T) (UserResponse, *vo.Password, *vo.HashedPassword) {
	id, err := vo.NewID()
	require.NoError(t, err)
	password, err := vo.NewPassword(utils.RandomPassword() + "A")
	require.NoError(t, err)
	username, err := vo.NewUsername(utils.RandomUserName())
	require.NoError(t, err)
	email, err := vo.NewEmail(utils.RandomEmail())
	require.NoError(t, err)

	hashedPassword, err := password.HashPassword()
	require.NoError(t, err)

	return UserResponse{
		User: &entity.Profile{
			ID:         id.Value().String(),
			Email:      email.Value(),
			Username:   username.Value(),
			FullName:   utils.RandomUserName(),
			FirstName:  utils.RandomUserName(),
			LastName:   utils.RandomUserName(),
			PictureURL: utils.RandomString(64),
			Verified:   true,
			Blocked:    false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}, password, hashedPassword
}

type eqUserCreateParamsMatcher struct {
	newUser  *entity.User
	password *vo.Password
}

func (e eqUserCreateParamsMatcher) Matches(x interface{}) bool {
	userCreated, ok := x.(*entity.User)
	if !ok {
		return false
	}
	if err := userCreated.Password.VerifyPassword(e.password.Value()); err != nil {
		return false
	}

	if !e.newUser.Email.IsEqual(userCreated.Email.Value()) {
		return false
	}

	if !e.newUser.Username.IsEqual(userCreated.Username.Value()) {
		return false
	}

	return true
}

func (e eqUserCreateParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.newUser, e.password.Value())
}

func EqCreateUserParams(arg *entity.User, password *vo.Password) gomock.Matcher {
	return eqUserCreateParamsMatcher{arg, password}
}
