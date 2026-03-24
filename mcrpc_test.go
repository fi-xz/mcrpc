package mcrpc

import (
	"context"
	"crypto/tls"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Test configuration
const (
	testHost = "nx-win.thrush-ide.ts.net"
	testPort = 12345
)

// getTestSecret reads the secret from the secrets file
func getTestSecret() (string, error) {
	secretPath := filepath.Join(".secrets", "secret.txt")
	data, err := os.ReadFile(secretPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// getTestCert loads the TLS certificate for testing
func getTestCert() (*tls.Certificate, error) {
	certPath := filepath.Join(".certs", "cert.crt")
	keyPath := filepath.Join(".certs", "cert.pem")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// createTestClient creates a client for integration tests
func createTestClient(t *testing.T) (*MCRPCClient, context.Context) {
	t.Helper()

	secret, err := getTestSecret()
	if err != nil {
		t.Skipf("Skipping test: cannot read secret: %v", err)
	}

	cert, err := getTestCert()
	if err != nil {
		t.Skipf("Skipping test: cannot load certificate: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	client, err := CreateWithTLS(ctx, testHost, testPort, secret, cert, true)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to server: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
	})

	return client, ctx
}

// TestClientCreation tests creating a client with TLS
func TestClientCreation(t *testing.T) {
	secret, err := getTestSecret()
	if err != nil {
		t.Skipf("Skipping test: cannot read secret: %v", err)
	}

	cert, err := getTestCert()
	if err != nil {
		t.Skipf("Skipping test: cannot load certificate: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := CreateWithTLS(ctx, testHost, testPort, secret, cert, true)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to server: %v", err)
	}
	defer client.Close()

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
