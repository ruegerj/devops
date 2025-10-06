package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type AuthContextKey int // use dedicated key type to prevent context key collisions

const (
	AuthSubjectKey AuthContextKey = iota
)

func Authenticate(privateKey string, logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseRequest(r, jwt.WithKey(jwa.HS256(), []byte(privateKey)))
		if err != nil {
			logger.Warn("Failed user authentication", slog.Any("error", err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		username := "anom"
		if subject, exists := token.Subject(); exists {
			username = subject
		}
		ctx := context.WithValue(r.Context(), AuthSubjectKey, username)
		r = r.WithContext(ctx)

		logger.Info(fmt.Sprintf("Authenticated user: %s", username), slog.String("username", username))
		next.ServeHTTP(w, r)
	})
}
