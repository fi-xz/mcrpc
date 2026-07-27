// Package mcrpc provides a JSON-RPC 2.0 over WebSocket client for managing
// Minecraft Java Edition servers. It implements the Minecraft Server Management
// Protocol (MSMP) and supports server management, player management, bans,
// operators, allowlists, server settings, and real-time notifications.
package mcrpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	rpcws "github.com/sourcegraph/jsonrpc2/websocket"
)

// MCRPCClient is a client for connecting to and managing Minecraft servers via
// JSON-RPC 2.0 over WebSocket. It provides methods for server management, player
// management, and real-time notifications.
type MCRPCClient struct {
	// JSONRPCConn is the underlying JSON-RPC 2.0 connection.
	JSONRPCConn *jsonrpc2.Conn
	// WebsocketConn is the underlying WebSocket connection.
	WebsocketConn *websocket.Conn
	closed        int32 // 0 = open, 1 = closed

	// OnNotification is called for all notifications when no specific handler is set.
	OnNotification func(method string, params json.RawMessage)

	// Server notification handlers
	OnServerStarted  func()
	OnServerStopping func()
	OnServerSaving   func()
	OnServerSaved    func()
	OnServerStatus   func(status ServerState)
	OnServerActivity func()

	// Player notification handlers
	OnPlayerJoined func(player Player)
	OnPlayerLeft   func(player Player)

	// Operator notification handlers
	OnOperatorAdded   func(operator Operator)
	OnOperatorRemoved func(operator Operator)

	// Allowlist notification handlers
	OnAllowlistAdded   func(player Player)
	OnAllowlistRemoved func(player Player)

	// Ban notification handlers
	OnBanAdded   func(ban UserBan)
	OnBanRemoved func(player Player)

	// IP Ban notification handlers
	OnIPBanAdded   func(ban IPBan)
	OnIPBanRemoved func(ip string)

	// Gamerule notification handlers
	OnGameruleUpdated func(gamerule TypedGameRule)

	// World notification handlers
	OnWorldUpgradeStarted  func()
	OnWorldUpgradeProgress func(progress float64)
	OnWorldUpgradeFinished func()
	OnWorldUpgradeFailed   func(reason string)
}

// Create establishes a plaintext WebSocket connection to the Minecraft server.
// The secret is used for authentication with the server.
func Create(ctx context.Context, host string, port int, secret string) (*MCRPCClient, error) {
	return createClient(ctx, host, port, secret, nil, false, false)
}

// CreateWithTLS establishes a TLS WebSocket (wss) connection to the Minecraft
// server. TLS is always used, whether or not a client certificate is supplied;
// cert may be nil for servers that do not require client authentication.
// If insecure is true, the server's certificate will not be verified.
func CreateWithTLS(ctx context.Context, host string, port int, secret string, cert *tls.Certificate, insecure bool) (*MCRPCClient, error) {
	return createClient(ctx, host, port, secret, cert, insecure, true)
}

// DisconnectNotify returns a channel that is closed when the connection is closed.
func (c *MCRPCClient) DisconnectNotify() <-chan struct{} {
	return c.JSONRPCConn.DisconnectNotify()
}

// Close closes the WebSocket connection and the underlying JSON-RPC connection.
// It is safe to call Close multiple times.
func (c *MCRPCClient) Close() error {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil // Connection is already closed.
	}
	return c.JSONRPCConn.Close()
}

// IsClosed returns true if the connection has been closed.
func (c *MCRPCClient) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func createClient(ctx context.Context, host string, port int, secret string, cert *tls.Certificate, insecure, useTLS bool) (*MCRPCClient, error) {
	header := http.Header{}

	// The scheme follows the caller's stated intent, not the presence of a
	// client certificate: a TLS server may well not require one.
	scheme := "ws"
	if useTLS {
		scheme = "wss"
	}
	dialString := fmt.Sprintf("%s://%s", scheme, net.JoinHostPort(host, strconv.Itoa(port)))

	header.Add("Authorization", "Bearer "+secret)
	header.Add("Sec-WebSocket-Protocol", "minecraft-v1,"+secret)

	dialer := websocket.Dialer{}

	if useTLS {
		dialer.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: insecure,
		}
		if cert != nil {
			dialer.TLSClientConfig.Certificates = []tls.Certificate{*cert}
		}
	}

	conn, response, err := dialer.Dial(dialString, header)

	if err != nil {
		return nil, fmt.Errorf("mcrpc: dial %s: %w", dialString, err)
	}

	if response.StatusCode != http.StatusSwitchingProtocols {
		conn.Close()
		return nil, fmt.Errorf("mcrpc: handshake failed: unexpected status %s", response.Status)
	}

	objectStream := rpcws.NewObjectStream(conn)

	client := &MCRPCClient{
		WebsocketConn: conn,
	}

	client.JSONRPCConn = jsonrpc2.NewConn(ctx, objectStream, client.handleIncoming())

	return client, nil
}
