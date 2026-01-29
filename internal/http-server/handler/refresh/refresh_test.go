package refresh_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	grpc "github.com/Weit145/proto-repo/auth"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/refresh"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/refresh/mocks"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRefreshHandler(t *testing.T) {
	AssetTokenResponse := &grpc.AccessTokenResponse{
		AccessToken: "123123",
	}
	cases := []struct {
		name         string
		tokenRefresh string
		respError    string
		mockError    error
		mockResult   *grpc.AccessTokenResponse
		shouldCall   bool
	}{
		{
			name:         "Missing tokenRefresh",
			tokenRefresh: "",
			respError:    "error",
			shouldCall:   false,
		},
		{
			name:         "Error mock",
			tokenRefresh: "123123",
			respError:    "error",
			mockError:    errors.New("Failed to refresh"),
			mockResult:   nil,
			shouldCall:   true,
		},
		{
			name:         "Success",
			tokenRefresh: "123123",
			respError:    "success",
			mockResult:   AssetTokenResponse,
			shouldCall:   true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRefreshToken := mocks.NewGRPCRefreshToken(t)

			if tc.shouldCall {
				mockRefreshToken.On("RefreshToken", mock.Anything, tc.tokenRefresh).
					Return(tc.mockResult, tc.mockError).
					Once()
			}

			handler := refresh.New(slogdiscard.NewDiscardLogger(), mockRefreshToken)

			req := httptest.NewRequest(http.MethodGet, "/refresh", nil)
			req.AddCookie(&http.Cookie{
				Name:     "refresh_token",
				Value:    tc.tokenRefresh,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})

			r := chi.NewRouter()
			r.Get("/refresh", handler)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if tc.tokenRefresh != "" && tc.mockError == nil {
				assetToken := rr.Result().Header.Get("Authorization")
				require.True(t, strings.HasPrefix(assetToken, "Bearer "))
			}
			var resp response.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}
