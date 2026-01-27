package refresh_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/refresh"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestRefreshHandler(t *testing.T) {
	cases := []struct {
		name         string
		tokenRefresh string
		respError    string
		mockError    error
	}{
		{
			name:         "Missing tokenRefresh",
			tokenRefresh: "",
			respError:    "error",
		},
		{
			name:         "Success",
			tokenRefresh: "123123",
			respError:    "success",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := refresh.New(slogdiscard.NewDiscardLogger())

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

			if tc.tokenRefresh != "" {
				assetToken := rr.Result().Header.Get("Authorization")
				require.True(t, strings.HasPrefix(assetToken, "Bearer "))
			}
			var resp response.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}
