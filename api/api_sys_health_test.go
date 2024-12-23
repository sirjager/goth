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

	"github.com/sirjager/goth/modules"
	mockRepo "github.com/sirjager/goth/repository/mock"
	mockTask "github.com/sirjager/goth/worker/mock"
)

func TestHealth(t *testing.T) {
	testCases := []struct {
		check    func(t *testing.T, recorder *httptest.ResponseRecorder, expected healthResponse)
		expected healthResponse
		name     string
	}{
		{
			name: "OK",
			expected: healthResponse{
				Service: testConfig.ServiceName,
				Server:  testConfig.ServerName,
				Status:  healthpb.HealthCheckResponse_SERVING.String(),
				Started: testConfig.StartTime,
			},
			check: func(t *testing.T, recorder *httptest.ResponseRecorder, expected healthResponse) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response healthResponse
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

			adapters := modules.NewModules(testConfig, testLogr, testCache, repo, testTokens, testMail, testTasks)
			server := NewServer(adapters)

			recoder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/health", nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recoder, request)
			tc.check(t, recoder, tc.expected)
		})
	}
}
