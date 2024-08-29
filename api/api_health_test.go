package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/sirjager/goth/core"
	mockRepo "github.com/sirjager/goth/repository/mock"
	mockTask "github.com/sirjager/goth/worker/mock"
)

func TestHealth(t *testing.T) {
	testCases := []struct {
		check    func(t *testing.T, recorder *httptest.ResponseRecorder, expected HealthResponse)
		expected HealthResponse
		name     string
	}{
		{
			name: "OK",
			expected: HealthResponse{
				Service: testConfig.ServiceName,
				Server:  testConfig.ServerName,
				Status:  healthpb.HealthCheckResponse_SERVING.String(),
				Started: testConfig.StartTime,
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, expected HealthResponse) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response HealthResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				require.Equal(t, expected.Server, response.Server)
				require.Equal(t, expected.Status, response.Status)
				require.Equal(t, expected.Service, response.Service)
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

			request, err := http.NewRequest(http.MethodGet, "/api/health", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.check(t, recoder, tc.expected)
		})
	}
}
