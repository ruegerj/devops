package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ruegerj/devops/api/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogging(t *testing.T) {
	// GIVEN
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	testHandler := middleware.Logging(logger, handler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// WHEN
	testHandler.ServeHTTP(w, req)

	// THEN
	assert.Equal(t, http.StatusOK, w.Code, "should pass through status code")

	// output should contain the respective fields
	output := buf.String()
	require.NotEmpty(t, output, "log output should not be empty")
	assert.Contains(t, output, `"statusCode":200`)
	assert.Contains(t, output, `"method":"GET"`)
	assert.Contains(t, output, `"path":"/"`)
	assert.Contains(t, output, `"duration"`) // should have some duration field

}
