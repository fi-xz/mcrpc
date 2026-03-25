package mcrpc

import (
	"fmt"
	"os"
	"testing"
)

// TestConnectionDiagnostics checks if all required configuration and connections work
func TestConnectionDiagnostics(t *testing.T) {
	t.Run("CheckConfiguration", func(t *testing.T) {
		host, port, secret, useTLS := getTestConfig()

		t.Logf("✓ Configuration loaded:")
		t.Logf("  - Host: %s", host)
		t.Logf("  - Port: %d", port)
		t.Logf("  - Secret: %s (length: %d)", maskSecret(secret), len(secret))
		t.Logf("  - Use TLS: %v", useTLS)

		if secret == "" {
			t.Fatal("Secret is empty - set TEST_SECRET environment variable")
		}
	})

	t.Run("CheckServerConnection", func(t *testing.T) {
		host, port, _, _ := getTestConfig()

		fmt.Printf("\nAttempting to connect to %s:%d...\n", host, port)

		client, ctx := createTestClient(t)

		t.Logf("✓ Successfully connected to server")
		t.Logf("  - WebSocket connected: %v", client.WebsocketConn != nil)
		t.Logf("  - JSON-RPC connected: %v", client.JSONRPCConn != nil)
		t.Logf("  - Client closed: %v", client.IsClosed())

		// Try to get server status
		status, err := client.GetServerStatus(ctx)
		if err != nil {
			t.Errorf("Failed to get server status: %v", err)
		} else {
			t.Logf("✓ Server status received:")
			t.Logf("  - Started: %v", status.Started)
			t.Logf("  - Version: %s (protocol: %d)", status.Version.Name, status.Version.Protocol)
			t.Logf("  - Players online: %d", len(status.Players))
		}
	})
}

// maskSecret returns a masked version of the secret for logging
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

// TestEnvironmentVariables tests that environment variables are properly read
func TestEnvironmentVariables(t *testing.T) {
	// Save original values
	origHost := os.Getenv("TEST_HOST")
	origPort := os.Getenv("TEST_PORT")
	origSecret := os.Getenv("TEST_SECRET")
	origTLS := os.Getenv("USE_TLS")

	// Restore after test
	defer func() {
		_ = os.Setenv("TEST_HOST", origHost)
		_ = os.Setenv("TEST_PORT", origPort)
		_ = os.Setenv("TEST_SECRET", origSecret)
		_ = os.Setenv("USE_TLS", origTLS)
	}()

	// Test with custom values
	_ = os.Setenv("TEST_HOST", "test.example.com")
	_ = os.Setenv("TEST_PORT", "99999")
	_ = os.Setenv("TEST_SECRET", "my-test-secret")
	_ = os.Setenv("USE_TLS", "true")

	host, port, secret, useTLS := getTestConfig()

	if host != "test.example.com" {
		t.Errorf("Expected host 'test.example.com', got '%s'", host)
	}

	if port != 99999 {
		t.Errorf("Expected port 99999, got %d", port)
	}

	if secret != "my-test-secret" {
		t.Errorf("Expected secret 'my-test-secret', got '%s'", secret)
	}

	if !useTLS {
		t.Error("Expected useTLS to be true")
	}

	// Test defaults when env vars are empty
	_ = os.Unsetenv("TEST_HOST")
	_ = os.Unsetenv("TEST_PORT")
	_ = os.Unsetenv("TEST_SECRET")
	_ = os.Unsetenv("USE_TLS")

	host, port, secret, useTLS = getTestConfig()

	if host != "127.0.0.1" {
		t.Errorf("Expected default host '127.0.0.1', got '%s'", host)
	}

	if port != 12345 {
		t.Errorf("Expected default port 12345, got %d", port)
	}

	if secret != "test-secret-for-ci-only" {
		t.Errorf("Expected default secret, got '%s'", secret)
	}

	if useTLS {
		t.Error("Expected useTLS to be false by default")
	}
}
