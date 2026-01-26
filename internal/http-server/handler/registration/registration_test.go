package registration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/registration"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/stretchr/testify/require"
)

func TestRegistrationHadler(t *testing.T) {
	cases := []struct {
		name      string
		username  string
		password  string
		email     string
		respError string
		mockError error
	}{
		{
			name:      "Missing username",
			username:  "",
			password:  "123456",
			email:     "kloader145@gmail.com",
			respError: "error",
		},
		{
			name:      "Missing password",
			username:  "Weit",
			password:  "",
			email:     "kloader145@gmail.com",
			respError: "error",
		},
		{
			name:      "Invalid username",
			username:  "Lg",
			password:  "123456",
			email:     "kloader145@gmail.com",
			respError: "error",
		},
		{
			name:      "Invalid password",
			username:  "Weit",
			password:  "123",
			email:     "kloader145@gmail.com",
			respError: "error",
		},
		{
			name:      "Missing email",
			username:  "Weit",
			password:  "123456",
			email:     "",
			respError: "error",
		},
		{
			name:      "Invalid email",
			username:  "Weit",
			password:  "123456",
			email:     "kloader145",
			respError: "error",
		},
		{
			name:      "Success",
			username:  "Weit",
			password:  "123456",
			email:     "kloader145@gmail.com",
			respError: "success",
		},
	}
	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := registration.New(slogdiscard.NewDiscardLogger())

			input := fmt.Sprintf(`{"email":"%s"}`, tc.email)

			req, err := http.NewRequest(http.MethodPost, "/order", strings.NewReader(input))
			require.NoError(t, err)

			req.SetBasicAuth(tc.username, tc.password)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()

			var resp response.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}
