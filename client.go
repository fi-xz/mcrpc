package mcrpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	rpcws "github.com/sourcegraph/jsonrpc2/websocket"
)

// DefaultHandshakeTimeout bounds the WebSocket handshake performed by Start
// unless WithHandshakeTimeout says otherwise.
const DefaultHandshakeTimeout = 10 * time.Second

// Client manages a Minecraft server over JSON-RPC 2.0 on a WebSocket.
//
// New only builds the value; nothing is dialled and no goroutine runs until
// Start. That split is what makes notification handlers safe to register:
// WithHandler is applied before any receive goroutine exists.
//
// A Client is safe for concurrent use, and can be restarted: after Close, a
// further Start dials again with the same configuration and handler.
type Client struct {
	addr             string
	secret           string
	handshakeTimeout time.Duration
	handler          Handler
	trace            func(TraceMessage)

	useTLS             bool
	tlsConfig          *tls.Config
	insecureSkipVerify bool
	clientCert         *tls.Certificate

	// mu guards the fields below, which change across Start and Close. It is
	// held for the duration of the handshake, so a Close concurrent with a
	// Start waits for that handshake to settle.
	mu     sync.RWMutex
	rpc    *jsonrpc2.Conn
	cancel context.CancelFunc
}

// New builds a client for the management server at addr, a "host:port" pair.
// It performs no I/O; call Start to connect.
//
// secret is the value of management-server-secret in server.properties.
// Without any option the connection is plaintext ws; see WithTLS.
func New(addr, secret string, opts ...Option) *Client {
	client := &Client{
		addr:             addr,
		secret:           secret,
		handshakeTimeout: DefaultHandshakeTimeout,
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// NewHostPort is New for callers that hold the host and port separately.
func NewHostPort(host string, port int, secret string, opts ...Option) *Client {
	return New(net.JoinHostPort(host, strconv.Itoa(port)), secret, opts...)
}

// Start dials the server and begins dispatching notifications. It returns
// ErrAlreadyStarted if the client is already connected.
//
// ctx bounds the whole session, not just the handshake: cancelling it closes
// the connection. To bound only the handshake use WithHandshakeTimeout, and
// pass a context that outlives the session here.
func (c *Client) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rpc != nil {
		return ErrAlreadyStarted
	}

	scheme := "ws"
	if c.useTLS {
		scheme = "wss"
	}
	url := fmt.Sprintf("%s://%s", scheme, c.addr)

	header := http.Header{}
	header.Add("Authorization", "Bearer "+c.secret)
	header.Add("Sec-WebSocket-Protocol", "minecraft-v1,"+c.secret)

	dialer := websocket.Dialer{
		TLSClientConfig: c.tlsSettings(),
	}
	if c.handshakeTimeout > 0 {
		dialer.HandshakeTimeout = c.handshakeTimeout
	}

	conn, response, err := dialer.DialContext(ctx, url, header)
	if err != nil {
		return fmt.Errorf("mcrpc: dial %s: %w", url, err)
	}
	if response.StatusCode != http.StatusSwitchingProtocols {
		conn.Close()
		return fmt.Errorf("mcrpc: handshake failed: unexpected status %s", response.Status)
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	rpc := jsonrpc2.NewConn(sessionCtx, rpcws.NewObjectStream(conn), c.handleIncoming(), c.traceOptions()...)

	c.rpc = rpc
	c.cancel = cancel

	// jsonrpc2's read loop does not watch the context, so cancellation has to
	// be turned into a Close here. The same goroutine also clears the client's
	// state when the server hangs up, which is what lets Start be called again.
	go func() {
		select {
		case <-sessionCtx.Done():
		case <-rpc.DisconnectNotify():
		}
		c.closeSession(rpc)
	}()

	return nil
}

// Close ends the session and releases the connection. It is safe to call
// multiple times and on a client that was never started, and it does not
// prevent a later Start from reconnecting.
func (c *Client) Close() error {
	return c.closeSession(nil)
}

// closeSession tears down the current session. When want is non-nil the
// teardown is skipped unless that connection is still the current one, so a
// watchdog left over from an earlier session cannot close its successor.
func (c *Client) closeSession(want *jsonrpc2.Conn) error {
	c.mu.Lock()
	if c.rpc == nil || (want != nil && c.rpc != want) {
		c.mu.Unlock()
		return nil
	}
	rpc, cancel := c.rpc, c.cancel
	c.rpc, c.cancel = nil, nil
	c.mu.Unlock()

	cancel()

	return rpc.Close()
}

// IsRunning reports whether the client currently holds a connection.
func (c *Client) IsRunning() bool {
	return c.conn() != nil
}

// DisconnectNotify returns a channel closed when the current session ends.
//
// The channel belongs to one session: after a Close and a further Start, call
// this again to observe the new one. For a client that is not running, the
// returned channel is already closed.
func (c *Client) DisconnectNotify() <-chan struct{} {
	if conn := c.conn(); conn != nil {
		return conn.DisconnectNotify()
	}

	closed := make(chan struct{})
	close(closed)
	return closed
}

// conn returns the live connection, or nil when the client is not running.
func (c *Client) conn() *jsonrpc2.Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rpc
}
