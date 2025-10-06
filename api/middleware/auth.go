package middleware

import (
	"log/slog"
	"net/http"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
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
		logger.Info("Authenticated user: ", slog.String("username", username))

		next.ServeHTTP(w, r)
	})
}
