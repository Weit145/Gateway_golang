package confirm_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	grpc "github.com/Weit145/proto-repo/auth"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/confirm"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/confirm/mocks"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConfirmHandler(t *testing.T) {
	CookieResponse := &grpc.CookieResponse{
		AccessToken: "123123",
		Cookie: &grpc.Cookie{
			Key:      "refresh_token",
			Value:    "",
			Httponly: true,
			Secure:   true,
			Samesite: "",
			MaxAge:   24,
		},
	}

	cases := []struct {
		name       string
		token      string
		respError  string
		mockError  error
		mockResult *grpc.CookieResponse
		shouldCall bool
	}{
		{
			name:       "Success",
			token:      "1231241231",
			respError:  "success",
			mockError:  nil,
			mockResult: CookieResponse,
			shouldCall: true,
		},
		{
			name:       "Error confirm user",
			token:      "1231241231",
			respError:  "error",
			mockError:  errors.New("Failed to confirm"),
			mockResult: nil,
			shouldCall: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockConfirmUser := mocks.NewGRPCConfirmUser(t)

			if tc.shouldCall {
				mockConfirmUser.On("RegistrationUser", mock.Anything, tc.token).
					Return(tc.mockResult, tc.mockError).
					Once()
			}

			handler := confirm.New(slogdiscard.NewDiscardLogger(), mockConfirmUser)

			req := httptest.NewRequest(http.MethodGet, "/registration/confirm/"+tc.token, nil)

			r := chi.NewRouter()
			r.Get("/registration/confirm/{token}", handler)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			var resp response.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, tc.respError, resp.Status)

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
		})
	}
}
