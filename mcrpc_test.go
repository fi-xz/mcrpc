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

// createTestClient creates a client for integration tests
func createTestClient(t *testing.T) (*MCRPCClient, context.Context) {
	t.Helper()

	host, port, secret, useTLS := getTestConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	var client *MCRPCClient
	var err error

	if useTLS {
		t.Skip("TLS connections require certificates - use environment variables to configure")
	} else {
		client, err = Create(ctx, host, port, secret)
	}

	if err != nil {
		t.Skipf("Skipping test: cannot connect to server at %s:%d: %v", host, port, err)
	}

	t.Cleanup(func() {
		_ = client.Close()
	})

	return client, ctx
}

// TestClientCreation tests creating a client
func TestClientCreation(t *testing.T) {
	host, port, secret, useTLS := getTestConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var client *MCRPCClient
	var err error

	if useTLS {
		t.Skip("TLS connections require certificates")
	} else {
		client, err = Create(ctx, host, port, secret)
	}

	if err != nil {
		t.Skipf("Skipping test: cannot connect to server at %s:%d: %v", host, port, err)
	}
	defer func() { _ = client.Close() }()

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.IsClosed() {
		t.Fatal("Expected client to be open")
	}

	if client.JSONRPCConn == nil {
		t.Fatal("Expected JSONRPCConn to be initialized")
	}

	if client.WebsocketConn == nil {
		t.Fatal("Expected WebsocketConn to be initialized")
	}
}

// TestClientClose tests closing the client
func TestClientClose(t *testing.T) {
	client, _ := createTestClient(t)

	err := client.Close()
	if err != nil {
		t.Errorf("Expected no error on close, got: %v", err)
	}

	if !client.IsClosed() {
		t.Error("Expected client to be closed")
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
