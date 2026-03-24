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
	"net/http"
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
}

// Create establishes a WebSocket connection to the Minecraft server without TLS.
// The secret is used for authentication with the server.
func Create(context context.Context, host string, port int, secret string) (*MCRPCClient, error) {
	return createMCRCPClient(context, host, port, secret, nil, false)
}

// CreateWithTLS establishes a WebSocket connection to the Minecraft server with TLS.
// The cert parameter is used for client certificate authentication.
// If insecure is true, the server's certificate will not be verified.
func CreateWithTLS(context context.Context, host string, port int, secret string, cert *tls.Certificate, insecure bool) (*MCRPCClient, error) {
	return createMCRCPClient(context, host, port, secret, cert, insecure)
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

func createMCRCPClient(ctx context.Context, host string, port int, secret string, cert *tls.Certificate, insecure bool) (*MCRPCClient, error) {
	header := http.Header{}

	// Determine WebSocket protocol based on TLS usage
	protocol := "ws"
	if cert != nil {
		protocol = "wss"
	}
	dialString := fmt.Sprintf("%s://%s:%d", protocol, host, port)

	header.Add("Authorization", "Bearer "+secret)
	header.Add("Sec-WebSocket-Protocol", "minecraft-v1,"+secret)

	dialer := websocket.Dialer{}

	if cert != nil {
		dialer.TLSClientConfig = &tls.Config{
			Certificates:       []tls.Certificate{*cert},
			InsecureSkipVerify: insecure,
		}
	}

	conn, response, err := dialer.Dial(dialString, header)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusSwitchingProtocols {
		return nil, err
	}

	objectStream := rpcws.NewObjectStream(conn)

	client := &MCRPCClient{
		WebsocketConn: conn,
	}

	client.JSONRPCConn = jsonrpc2.NewConn(ctx, objectStream, client.handleIncoming())

	return client, nil
}
