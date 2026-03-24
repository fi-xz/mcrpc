package mcrpc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestConnectionDiagnostics checks if all required files and connections work
func TestConnectionDiagnostics(t *testing.T) {
	t.Run("CheckSecretFile", func(t *testing.T) {
		secretPath := filepath.Join(".secrets", "secret.txt")
		info, err := os.Stat(secretPath)
		if err != nil {
			t.Fatalf("Secret file not found at %s: %v", secretPath, err)
		}
		t.Logf("✓ Secret file exists: %s (size: %d bytes)", secretPath, info.Size())

		data, err := os.ReadFile(secretPath)
		if err != nil {
			t.Fatalf("Cannot read secret file: %v", err)
		}
		secret := string(data)
		if len(secret) == 0 {
			t.Fatal("Secret file is empty")
		}
		t.Logf("✓ Secret loaded successfully (length: %d)", len(secret))
	})

	t.Run("CheckCertFiles", func(t *testing.T) {
		certPath := filepath.Join(".certs", "cert.crt")
		keyPath := filepath.Join(".certs", "cert.pem")

		certInfo, err := os.Stat(certPath)
		if err != nil {
			t.Fatalf("Certificate file not found at %s: %v", certPath, err)
		}
		t.Logf("✓ Certificate file exists: %s (size: %d bytes)", certPath, certInfo.Size())

		keyInfo, err := os.Stat(keyPath)
		if err != nil {
			t.Fatalf("Key file not found at %s: %v", keyPath, err)
		}
		t.Logf("✓ Key file exists: %s (size: %d bytes)", keyPath, keyInfo.Size())
	})

	t.Run("CheckServerConnection", func(t *testing.T) {
		_, err := getTestSecret()
		if err != nil {
			t.Skipf("Cannot read secret: %v", err)
		}

		_, err = getTestCert()
		if err != nil {
			t.Skipf("Cannot load certificate: %v", err)
		}

		fmt.Printf("\nAttempting to connect to %s:%d...\n", testHost, testPort)

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
