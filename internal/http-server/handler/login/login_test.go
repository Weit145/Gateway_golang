package login_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	grpc "github.com/Weit145/proto-repo/auth"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/login"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/login/mocks"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	CookieResponse := &grpc.CookieResponse{
		AccessToken: "123123",
		Cookie: &grpc.Cookie{
			Key:      "refresh_token",
			Value:    "123123",
			Httponly: false,
			Secure:   false,
			Samesite: "lax",
			MaxAge:   24,
		},
	}
	cases := []struct {
		name       string
		username   string
		password   string
		respError  string
		mockError  error
		mockResult *grpc.CookieResponse
		shouldCall bool
	}{
		{
			name:       "Missing username",
			username:   "",
			password:   "123123",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Missing password",
			username:   "Weit",
			password:   "",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Invalid username",
			username:   "12",
			password:   "123123",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Invalid password",
			username:   "Weit",
			password:   "12",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Error mocks",
			username:   "Weit",
			password:   "123123",
			respError:  "error",
			mockResult: CookieResponse,
			mockError:  errors.New("Faeld login"),
			shouldCall: true,
		},
		{
			name:       "Success",
			username:   "Weit",
			password:   "123123",
			respError:  "success",
			mockResult: CookieResponse,
			shouldCall: true,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mocksAuthenticate := mocks.NewGRPCLoginUser(t)

			if tc.shouldCall {
				mocksAuthenticate.On("Authenticate", mock.Anything, tc.username, tc.password).
					Return(tc.mockResult, tc.mockError).
					Once()
			}

			handler := login.New(slogdiscard.NewDiscardLogger(), mocksAuthenticate)

			req := httptest.NewRequest(http.MethodPost, "/login", nil)

			req.SetBasicAuth(tc.username, tc.password)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			if tc.respError != "error" {
				cookies := rr.Result().Cookies()
				require.NotEmpty(t, cookies)

				var name string

				for _, cookie := range cookies {
					if cookie.Name == "refresh_token" {
						name = cookie.Name
						break
					}
				}
				require.NotEqual(t, name, "")

				assetToken := rr.Result().Header.Get("Authorization")
				require.True(t, strings.HasPrefix(assetToken, "Bearer "))
			}
			var resp response.Response

			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}
