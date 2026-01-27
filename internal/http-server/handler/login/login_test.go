package login_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/login"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	cases := []struct {
		name      string
		username  string
		password  string
		respError string
		mockError error
	}{
		{
			name:      "Missing username",
			username:  "",
			password:  "123123",
			respError: "error",
		},
		{
			name:      "Missing password",
			username:  "Weit",
			password:  "",
			respError: "error",
		},
		{
			name:      "Invalid username",
			username:  "12",
			password:  "123123",
			respError: "error",
		},
		{
			name:      "Invalid password",
			username:  "Weit",
			password:  "12",
			respError: "error",
		},
		{
			name:      "Success",
			username:  "Weit",
			password:  "123123",
			respError: "success",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := login.New(slogdiscard.NewDiscardLogger())

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
