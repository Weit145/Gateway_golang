package registration_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/registration"
	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/registration/mocks"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	grpc "github.com/Weit145/proto-repo/auth"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegistrationHandler(t *testing.T) {

	okResponse := &grpc.Okey{
		Success: true,
	}

	cases := []struct {
		name       string
		username   string
		login      string
		password   string
		email      string
		respError  string
		mockError  error
		mockResult *grpc.Okey
		shouldCall bool
	}{
		{
			name:       "Missing username",
			username:   "",
			login:      "weit145",
			password:   "123456",
			email:      "kloader145@gmail.com",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Missing password",
			username:   "Weit",
			login:      "weit145",
			password:   "",
			email:      "kloader145@gmail.com",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Invalid username (too short)",
			username:   "Lg",
			login:      "weit145",
			password:   "123456",
			email:      "kloader145@gmail.com",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Invalid password (too short)",
			username:   "Weit",
			login:      "weit145",
			password:   "123",
			email:      "kloader145@gmail.com",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Missing email",
			username:   "Weit",
			login:      "weit145",
			password:   "123456",
			email:      "",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Invalid email",
			username:   "Weit",
			login:      "weit145",
			password:   "123456",
			email:      "kloader145",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Missing login",
			username:   "Weit",
			login:      "",
			password:   "123456",
			email:      "kloader145@gmail.com",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Invalid login (too short)",
			username:   "Weit",
			login:      "we",
			password:   "123456",
			email:      "kloader145@gmail.com",
			respError:  "error",
			shouldCall: false,
		},
		{
			name:       "Error mock",
			username:   "Weit",
			login:      "weit145",
			password:   "123456",
			email:      "kloader145@gmail.com",
			respError:  "error",
			mockError:  errors.New("Failed to create"),
			mockResult: nil,
			shouldCall: true,
		},
		{
			name:       "Success",
			username:   "Weit",
			login:      "weit145",
			password:   "123456",
			email:      "kloader145@gmail.com",
			respError:  "success",
			mockError:  nil,
			mockResult: okResponse,
			shouldCall: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockCreateUser := mocks.NewGRPCRegistrationUser(t)

			if tc.shouldCall {
				mockCreateUser.On("CreateUser", mock.Anything, tc.login, tc.email, tc.password, tc.username).
					Return(tc.mockResult, tc.mockError).
					Once()
			}

			handler := registration.New(slogdiscard.NewDiscardLogger(), mockCreateUser)

			input := fmt.Sprintf(`{"email":"%s","login":"%s"}`, tc.email, tc.login)

			req, err := http.NewRequest(http.MethodPost, "/registration", strings.NewReader(input))
			require.NoError(t, err)

			req.SetBasicAuth(tc.username, tc.password)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()

			var resp response.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Status)

			mockCreateUser.AssertExpectations(t)
		})
	}
}
