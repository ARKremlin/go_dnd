package middleware

import (
	"net/http"
	"strings"

	"github.com/ARKremlin/go_dnd/internal/pkg/token"
)

const bearerPrefix = "Bearer "

func Auth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := r.Header.Get("Authorization")
			if authz == "" {
				unauthorized(w, "missing authorization header")
				return
			}
			if !strings.HasPrefix(authz, bearerPrefix) {
				unauthorized(w, "invalid authorization header")
				return
			}

			tokStr := strings.TrimSpace(authz[len(bearerPrefix):])
			if tokStr == "" {
				unauthorized(w, "empty token")
				return
			}

			claims, err := token.ValidateToken(tokStr, secret)
			if err != nil {
				unauthorized(w, "invalid token")
				return
			}

			ctx := WithUserID(r.Context(), claims.UserID)
			ctx = WithRole(ctx, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func unauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
	http.Error(w, msg, http.StatusUnauthorized)
}
