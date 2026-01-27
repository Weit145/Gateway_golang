package confirm_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Weit145/GATEWAY_golang/internal/http-server/handler/confirm"
	"github.com/Weit145/GATEWAY_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/GATEWAY_golang/internal/lib/response"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

func TestConfirmHandler(t *testing.T) {
	cases := []struct {
		name      string
		token     string
		respError string
		mockError error
	}{
		{
			name:      "Success",
			token:     "1231241231",
			respError: "success",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := confirm.New(slogdiscard.NewDiscardLogger())

			req := httptest.NewRequest(http.MethodGet, "/registration/confirm/"+tc.token, nil)
			r := chi.NewRouter()
			r.Get("/registration/confirm/{token}", handler)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

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

			var resp response.Response
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.Equal(t, tc.respError, resp.Status)
		})
	}
}
