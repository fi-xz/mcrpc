package mcrpc

import (
	"context"
	"os"
	"strconv"
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

// createTestClient creates a client for integration tests. The returned
// context carries a deadline for individual calls; the session itself is bound
// to a separate context so that a slow call cannot tear down the connection.
func createTestClient(t *testing.T) (*Client, context.Context) {
	t.Helper()

	host, port, secret, useTLS := getTestConfig()

	if useTLS {
		t.Skip("TLS connections require certificates - use environment variables to configure")
	}

	sessionCtx, stop := context.WithCancel(context.Background())
	t.Cleanup(stop)

	client := NewHostPort(host, port, secret, testTraceOption(t)...)
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
	host, port, secret, useTLS := getTestConfig()

	if useTLS {
		t.Skip("TLS connections require certificates")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewHostPort(host, port, secret, testTraceOption(t)...)
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
