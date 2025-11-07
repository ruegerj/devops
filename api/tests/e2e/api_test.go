//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	apiContainer tc.Container
	apiBaseURL   string
	validJwt     string
	invalidJwt   string
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	jwtSigningKey := "v3ryS3cure!"

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to obtain working dir:", err)
		os.Exit(1)
	}

	// build & start container from dockerfile
	req := &tc.ContainerRequest{
		FromDockerfile: tc.FromDockerfile{
			Context: "../../",
		},
		ExposedPorts: []string{"3000/tcp"},
		Env: map[string]string{
			"JWT_KEY": jwtSigningKey,
		},
		WaitingFor: wait.ForHTTP("/health").
			WithStartupTimeout(10 * time.Second),
	}

	if isInstrumentationEnabled() {
		coverageFlag := "true"
		fmt.Println("Building a testcontainer image using an instrumented binary...")

		workingDir = path.Join(workingDir, "..", "..")
		coverageDir := path.Join(workingDir, "bin", "coverage")
		if err := os.MkdirAll(coverageDir, 0o755); err != nil {
			log.Fatalf("Failed to create coverage dir: %v", err)
		}

		req.FromDockerfile = tc.FromDockerfile{
			Context:   "../..",
			BuildArgs: map[string]*string{"WITH_COVERAGE": &coverageFlag},
		}
		req.Env["GOCOVERDIR"] = "/app/coverage"
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.Binds = []string{
				fmt.Sprintf("%s:/app/coverage", coverageDir),
			}
		}
	}

	apiContainer, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: *req,
		Started:          true,
	})
	if err != nil {
		fmt.Println("Failed to start testcontainer:", err)
		os.Exit(1)
	}

	port, err := apiContainer.MappedPort(ctx, "3000")
	if err != nil {
		fmt.Println("Failed to obtain mapped port:", err)
		os.Exit(1)
	}

	apiBaseURL = fmt.Sprintf("http://localhost:%s", port.Port())
	fmt.Printf("Testcontainer listening on: %s...", apiBaseURL)

	validJwt = issueSignedJwtOrFailWith(jwtSigningKey)
	invalidJwt = issueSignedJwtOrFailWith("fak3K3y") // produce invalid jwt by signing with different key

	// run test suite
	code := m.Run()

	// teardown container gracefully -> allow dumpting of coverage data
	gracePeriod := 1 * time.Second
	if err = apiContainer.Terminate(ctx, tc.StopTimeout(gracePeriod)); err != nil {
		fmt.Println("Failed to stop container:", err)
	}

	os.Exit(code)
}

func TestHealthEndpoint(t *testing.T) {
	// WHEN
	resp, err := http.Get(fmt.Sprintf("%s/health", apiBaseURL))

	// THEN
	require.NoError(t, err, "request must succeed")
	require.Equal(t, http.StatusOK, resp.StatusCode, "must have status 200")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	rawBody, err := readBodyFrom(t, resp)
	assert.NoError(t, err)
	assert.Contains(t, rawBody, "\"status\":\"healthy\"")
}

func TestSecretEndoint_ValidJwt_ResponseReceived(t *testing.T) {
	// GIVEN
	expectedResponse := "{\"message\":\"Life, Universe and everything\",\"number\":42}"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/secret", apiBaseURL), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", validJwt))

	// WHEN
	resp, err := http.DefaultClient.Do(req)

	// THEN
	require.NoError(t, err, "request must succeed")
	require.Equal(t, http.StatusOK, resp.StatusCode, "must have status 200")
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	rawBody, err := readBodyFrom(t, resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, rawBody)
}

func TestSecretEndpoint_InvalidJwt_RequestDenied(t *testing.T) {
	// GIVEN
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/secret", apiBaseURL), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", invalidJwt))

	// WHEN
	resp, err := http.DefaultClient.Do(req)

	// THEN
	require.NoError(t, err, "request must succeed")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode, "must have status 401")
}

func issueSignedJwtOrFailWith(key string) string {
	token, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Subject("john.doe").
		Build()
	if err != nil {
		fmt.Println("Failed to create JWT:", err)
		os.Exit(1)
	}

	signedToken, err := jwt.Sign(token, jwt.WithKey(jwa.HS256(), []byte(key)))
	if err != nil {
		fmt.Println("Failed to sign JWT:", err)
	}

	return string(signedToken)
}

func readBodyFrom(t *testing.T, resp *http.Response) (string, error) {
	t.Helper()
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %w", err)
	}

	return strings.Trim(string(body), "\n"), nil
}

func isInstrumentationEnabled() bool {
	value := os.Getenv("INSTRUMENT_BINARY")
	return value == "true"
}
