package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mockRepo "github.com/sirjager/goth/repository/mock"
	mockTask "github.com/sirjager/goth/worker/mock"
)

func TestSwaggerDocs(t *testing.T) {
	testCases := []struct {
		check func(t *testing.T, recorder *httptest.ResponseRecorder)
		name  string
	}{
		{
			name: "OK",
			check: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockRepo.NewMockRepo(ctrl)
			testTasks := mockTask.NewMockTaskDistributor(ctrl)

			server := NewServer(repo, testLogr, testConfig, testCache, testTokens, testTasks)
			recoder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/swagger", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.check(t, recoder)
		})
	}
}
