package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	mockRepo "github.com/sirjager/goth/repository/mock"
)

func TestWelcome(t *testing.T) {
	testCases := []struct {
		check    func(t *testing.T, recorder *httptest.ResponseRecorder, expected welcomeResponse)
		expected welcomeResponse
		name     string
	}{
		{
			name:     "OK",
			expected: welcomeResponse{Message: welcomeMessaage(testConfig.ServiceName), Docs: docsPath},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, expected welcomeResponse) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response welcomeResponse
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

			repo := mockRepo.NewMockRepository(ctrl)

			server := NewServer(repo, testLogr, testConfig, testCache, testTokens, testTasks)
			recoder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.check(t, recoder, tc.expected)
		})
	}
}