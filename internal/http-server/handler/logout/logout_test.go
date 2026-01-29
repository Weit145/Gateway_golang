package logout_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/logout"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/logout/mocks"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogOutHandler(t *testing.T) {
	cases := []struct {
		name       string
		tokenAsset string
		respError  string
		mockError  error
		shouldCall bool
	}{
		{
			name:       "Missing token",
			tokenAsset: "",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Invalid token",
			tokenAsset: "123123",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Error mock",
			tokenAsset: "Bearer 123123",
			respError:  "error",
			mockError:  errors.New("Faeld logout"),
			shouldCall: true,
		},
		{
			name:       "Success",
			tokenAsset: "Bearer 123123",
			respError:  "success",
			shouldCall: true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mocksGRPCLogOutUser := mocks.NewGRPCLogOutUser(t)

			if tc.shouldCall {
				mocksGRPCLogOutUser.On("LogOutUser", mock.Anything, tc.tokenAsset).
					Return(tc.mockError).
					Once()
			}

			handler := logout.New(slogdiscard.NewDiscardLogger(), mocksGRPCLogOutUser)

			req := httptest.NewRequest(http.MethodGet, "/logout", nil)
			req.Header.Add("Authorization", tc.tokenAsset)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if tc.respError != "error" {
				assetToken := rr.Result().Header.Get("Authorization")
				require.False(t, strings.HasPrefix(assetToken, "Bearer "))
			}

			var resp response.Response

			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}
