package middleware_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/ruegerj/devops/api/middleware"
	"github.com/stretchr/testify/assert"
)

const testJwtSigningKey = "veryS3cure!"

func TestAuthenticate_ValidTokenGiven_RequestPass(t *testing.T) {
	// GIVEN
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	testUser := "john.doe"
	validJwt, err := buildAndSignJwt(t, testUser, testJwtSigningKey)
	if err != nil {
		t.Fatal(err)
	}

	var capturedCtx context.Context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()
		w.WriteHeader(http.StatusOK)
	})
	testHandler := middleware.Authenticate(testJwtSigningKey, logger, http.HandlerFunc(handler))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", validJwt))
	w := httptest.NewRecorder()

	// WHEN
	testHandler.ServeHTTP(w, req)

	// THEN
	assert.Equal(t, http.StatusOK, w.Code, "should pass through request")
	assert.Equal(t, testUser, capturedCtx.Value(middleware.AuthSubjectKey))
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Authenticated user: john.doe")
}

func TestAuthenticate_InvalidTokenGiven_RequestDenied(t *testing.T) {
	// GIVEN
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	testUser := "john.doe"
	invalidJwt, err := buildAndSignJwt(t, testUser, "not-the-og-key") // use alternate key to provoke invalid signature
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	testHandler := middleware.Authenticate(testJwtSigningKey, logger, http.HandlerFunc(handler))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", invalidJwt))
	w := httptest.NewRecorder()

	// WHEN
	testHandler.ServeHTTP(w, req)

	// THEN
	assert.Equal(t, http.StatusUnauthorized, w.Code, "should deny request")
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Failed user authentication")

}

func buildAndSignJwt(t *testing.T, username, key string) (string, error) {
	t.Helper()
	token, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Subject(username).
		Build()
	if err != nil {
		return "", nil
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(key)))
	if err != nil {
		return "", err
	}

	return string(signedToken), nil
}
