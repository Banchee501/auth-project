package middleware

import (
	jwtpkg "auth-project/pkg/jwt"
	"context"
	"net/http"
	"os"
)

type contextKey string

const userKey contextKey = "user_id"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := jwtpkg.Parse(cookie.Value, os.Getenv("JWT_SECRET"))
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userKey).(int)
	return id, ok
}
