package mcrpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// getTestConfig returns test configuration from environment variables or defaults
func getTestConfig() (host string, port int, secret string, useTLS bool) {
	host = os.Getenv("TEST_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	portStr := os.Getenv("TEST_PORT")
	port = 12345
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	secret = os.Getenv("TEST_SECRET")
	if secret == "" {
		secret = "test-secret-for-ci-only"
	}

	useTLS = os.Getenv("USE_TLS") == "true"

	return host, port, secret, useTLS
}

// testSecret prefers TEST_SECRET, then the local .secrets/secret.txt, then the
// value CI uses. getTestConfig deliberately stays environment-only so that
// TestEnvironmentVariables remains deterministic; the file fallback lives here.
func testSecret() string {
	if secret := os.Getenv("TEST_SECRET"); secret != "" {
		return secret
	}

	if data, err := os.ReadFile(filepath.Join(".secrets", "secret.txt")); err == nil {
		if secret := strings.TrimSpace(string(data)); secret != "" {
			return secret
		}
	}

	return "test-secret-for-ci-only"
}

// testTLSOptions builds the TLS configuration for a live run:
//
//	USE_TLS=true                connect over wss
//	TEST_TLS_SERVER_NAME=name   verify against this name rather than the dialled
//	                            host, for a server bound to loopback whose
//	                            certificate names something else
//	TEST_TLS_CA=path            trust this PEM as a root, defaulting to
//	                            .certs/ca.crt when it exists
//	TEST_TLS_INSECURE=true      skip verification entirely (last resort)
//
// Note that .certs/cert.crt is the *server's* identity, not a client
// certificate, so it is not offered as one.
func testTLSOptions(t *testing.T) []Option {
	t.Helper()

	if os.Getenv("USE_TLS") != "true" {
		return nil
	}

	if os.Getenv("TEST_TLS_INSECURE") == "true" {
		return []Option{WithInsecureSkipVerify()}
	}

	config := &tls.Config{
		ServerName: os.Getenv("TEST_TLS_SERVER_NAME"),
	}

	caPath := os.Getenv("TEST_TLS_CA")
	if caPath == "" {
		if _, err := os.Stat(filepath.Join(".certs", "ca.crt")); err == nil {
			caPath = filepath.Join(".certs", "ca.crt")
		}
	}

	if caPath != "" {
		pem, err := os.ReadFile(caPath)
		if err != nil {
			t.Fatalf("cannot read TEST_TLS_CA %s: %v", caPath, err)
		}
		roots := x509.NewCertPool()
		if !roots.AppendCertsFromPEM(pem) {
			t.Fatalf("no certificates found in %s", caPath)
		}
		config.RootCAs = roots
	}

	return []Option{WithTLS(config)}
}

// testTraceOption enables wire tracing when TEST_TRACE is set, so that a run
// against a live server shows the exact JSON exchanged:
//
//	TEST_TRACE=1 go test -v -run TestBanlist ./...
//
// Traces are written with t.Log, so `-v` is required to see them.
func testTraceOption(t *testing.T) []Option {
	t.Helper()

	if os.Getenv("TEST_TRACE") == "" {
		return nil
	}

	return []Option{WithTrace(func(message TraceMessage) {
		t.Log(message)
	})}
}

// testClientOptions is the full option set an integration test connects with.
func testClientOptions(t *testing.T) []Option {
	t.Helper()

	options := testTLSOptions(t)
	return append(options, testTraceOption(t)...)
}

// createTestClient creates a client for integration tests. The returned
// context carries a deadline for individual calls; the session itself is bound
// to a separate context so that a slow call cannot tear down the connection.
func createTestClient(t *testing.T) (*Client, context.Context) {
	t.Helper()

	host, port, _, _ := getTestConfig()

	sessionCtx, stop := context.WithCancel(context.Background())
	t.Cleanup(stop)

	client := NewHostPort(host, port, testSecret(), testClientOptions(t)...)
	if err := client.Start(sessionCtx); err != nil {
		t.Skipf("Skipping test: cannot connect to server at %s:%d: %v", host, port, err)
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	callCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	return client, callCtx
}

// TestClientCreation tests connecting to a live server
func TestClientCreation(t *testing.T) {
	host, port, _, _ := getTestConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewHostPort(host, port, testSecret(), testClientOptions(t)...)
	if err := client.Start(ctx); err != nil {
		t.Skipf("Skipping test: cannot connect to server at %s:%d: %v", host, port, err)
	}
	defer func() { _ = client.Close() }()

	if !client.IsRunning() {
		t.Fatal("Expected client to be running")
	}
}

// TestClientClose tests closing the client
func TestClientClose(t *testing.T) {
	client, _ := createTestClient(t)

	err := client.Close()
	if err != nil {
		t.Errorf("Expected no error on close, got: %v", err)
	}

	if client.IsRunning() {
		t.Error("Expected client to be stopped")
	}

	// Closing again should not error
	err = client.Close()
	if err != nil {
		t.Errorf("Expected no error on second close, got: %v", err)
	}
}

// TestDisconnectNotify tests the disconnect notification channel
func TestDisconnectNotify(t *testing.T) {
	client, _ := createTestClient(t)

	disconnectCh := client.DisconnectNotify()
	if disconnectCh == nil {
		t.Fatal("Expected non-nil disconnect channel")
	}
}
