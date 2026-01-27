package logout_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/logout"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/stretchr/testify/require"
)

func TestLogOutHandler(t *testing.T) {
	cases := []struct {
		name       string
		tokenAsset string
		respError  string
		mockError  error
	}{
		{
			name:       "Missing token",
			tokenAsset: "",
			respError:  "error",
		},
		{
			name:       "Invalid token",
			tokenAsset: "123123",
			respError:  "error",
		},
		{
			name:       "Success",
			tokenAsset: "Bearer 123123",
			respError:  "success",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := logout.New(slogdiscard.NewDiscardLogger())

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
