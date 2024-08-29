package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/sirjager/goth/core"
	mockRepo "github.com/sirjager/goth/repository/mock"
	mockTask "github.com/sirjager/goth/worker/mock"
)

func TestWelcome(t *testing.T) {
	testCases := []struct {
		check    func(t *testing.T, recorder *httptest.ResponseRecorder, expected WelcomResponse)
		expected WelcomResponse
		name     string
	}{
		{
			name:     "OK",
			expected: WelcomResponse{Message: welcomeMessaage(testConfig.ServiceName), Docs: "/api/docs"},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, expected WelcomResponse) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response WelcomResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)
				require.Equal(t, expected, response)

				require.Equal(t, expected.Docs, response.Docs)
				require.Equal(t, expected.Message, response.Message)
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

			app := core.NewCoreApp(testConfig, testLogr, testCache, repo, testTokens, testMail, testTasks)
			server := NewServer(app)
			recoder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/api", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.check(t, recoder, tc.expected)
		})
	}
}
