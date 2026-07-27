package mcrpc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	rpcws "github.com/sourcegraph/jsonrpc2/websocket"
)

// fakeServer is a minimal JSON-RPC over WebSocket endpoint standing in for a
// Minecraft management server, so the client's lifecycle can be exercised
// without one.
type fakeServer struct {
	http *httptest.Server
	// conns receives every accepted connection, letting a test push
	// notifications or hang up.
	conns chan *jsonrpc2.Conn
}

func newFakeServer(t *testing.T, respond func(method string) (any, error)) *fakeServer {
	t.Helper()

	server := &fakeServer{conns: make(chan *jsonrpc2.Conn, 4)}

	upgrader := websocket.Upgrader{}
	server.http = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		handler := handlerFunc(func(ctx context.Context, rpc *jsonrpc2.Conn, req *jsonrpc2.Request) {
			if req.Notif {
				return
			}
			result, err := respond(req.Method)
			if err != nil {
				_ = rpc.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{Code: -32602, Message: err.Error()})
				return
			}
			_ = rpc.Reply(ctx, req.ID, result)
		})

		rpc := jsonrpc2.NewConn(context.Background(), rpcws.NewObjectStream(conn), handler)
		select {
		case server.conns <- rpc:
		default:
		}
	}))

	t.Cleanup(server.http.Close)

	return server
}

// addr returns the host:port the fake server listens on.
func (s *fakeServer) addr() string {
	return strings.TrimPrefix(s.http.URL, "http://")
}

// nextConn waits for the server side of the next accepted connection.
func (s *fakeServer) nextConn(t *testing.T) *jsonrpc2.Conn {
	t.Helper()
	select {
	case conn := <-s.conns:
		return conn
	case <-time.After(2 * time.Second):
		t.Fatal("server never accepted a connection")
		return nil
	}
}

func okResponder(string) (any, error) { return []Player{}, nil }

// errServerRejected stands in for a server-side rejection; the fake server maps
// it to JSON-RPC code -32602.
var errServerRejected = errors.New("invalid params")

func TestNewPerformsNoIO(t *testing.T) {
	// An address nothing listens on: if New dialled, this would fail or block.
	client := New("127.0.0.1:1", "secret")

	if client.IsRunning() {
		t.Error("a client that was never started reports itself running")
	}

	if _, err := client.GetPlayers(context.Background()); !errors.Is(err, ErrNotConnected) {
		t.Errorf("call before Start: got %v, want ErrNotConnected", err)
	}

	select {
	case <-client.DisconnectNotify():
	default:
		t.Error("DisconnectNotify should already be closed before Start")
	}
}

func TestStartAndClose(t *testing.T) {
	server := newFakeServer(t, okResponder)
	client := New(server.addr(), "secret")

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if !client.IsRunning() {
		t.Error("client should be running after Start")
	}

	if _, err := client.GetPlayers(context.Background()); err != nil {
		t.Errorf("GetPlayers failed: %v", err)
	}

	if err := client.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
	if client.IsRunning() {
		t.Error("client should not be running after Close")
	}
	if err := client.Close(); err != nil {
		t.Errorf("second Close should be a no-op, got %v", err)
	}

	if _, err := client.GetPlayers(context.Background()); !errors.Is(err, ErrNotConnected) {
		t.Errorf("call after Close: got %v, want ErrNotConnected", err)
	}
}

func TestStartTwiceIsRejected(t *testing.T) {
	server := newFakeServer(t, okResponder)
	client := New(server.addr(), "secret")

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	if err := client.Start(context.Background()); !errors.Is(err, ErrAlreadyStarted) {
		t.Errorf("second Start: got %v, want ErrAlreadyStarted", err)
	}
}

func TestRestartAfterClose(t *testing.T) {
	server := newFakeServer(t, okResponder)
	client := New(server.addr(), "secret")

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("first Start failed: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// The configuration and handler survive, so reconnecting is one call.
	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("restart failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	if _, err := client.GetPlayers(context.Background()); err != nil {
		t.Errorf("GetPlayers after restart failed: %v", err)
	}
}

func TestContextCancellationClosesSession(t *testing.T) {
	server := newFakeServer(t, okResponder)
	client := New(server.addr(), "secret")

	ctx, cancel := context.WithCancel(context.Background())
	if err := client.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	disconnected := client.DisconnectNotify()
	cancel()

	select {
	case <-disconnected:
	case <-time.After(2 * time.Second):
		t.Fatal("cancelling the context did not close the session")
	}

	// Give the watchdog a moment to clear the client's state.
	deadline := time.Now().Add(2 * time.Second)
	for client.IsRunning() && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if client.IsRunning() {
		t.Error("client still reports itself running after its context was cancelled")
	}
}

// TestHandlerRegisteredBeforeConnection is the point of splitting New from
// Start: a notification that arrives immediately cannot outrun handler
// registration, because no connection exists until Start.
func TestHandlerRegisteredBeforeConnection(t *testing.T) {
	server := newFakeServer(t, okResponder)

	joined := make(chan Player, 1)
	client := New(server.addr(), "secret", WithHandler(Handler{
		OnPlayerJoined: func(player Player) { joined <- player },
	}))

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	conn := server.nextConn(t)
	if err := conn.Notify(context.Background(), "minecraft:notification/players/joined", PlayerByName("fi_xz")); err != nil {
		t.Fatalf("server failed to notify: %v", err)
	}

	select {
	case player := <-joined:
		if player.Name != "fi_xz" {
			t.Errorf("got player %q, want fi_xz", player.Name)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OnPlayerJoined was never called")
	}
}

func TestOnErrorReportsUndecodableNotification(t *testing.T) {
	server := newFakeServer(t, okResponder)

	failures := make(chan error, 1)
	client := New(server.addr(), "secret", WithHandler(Handler{
		OnPlayerJoined: func(Player) {},
		OnError:        func(_ string, err error) { failures <- err },
	}))

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	conn := server.nextConn(t)
	if err := conn.Notify(context.Background(), "minecraft:notification/players/joined", "not-a-player"); err != nil {
		t.Fatalf("server failed to notify: %v", err)
	}

	select {
	case err := <-failures:
		if err == nil {
			t.Error("OnError called with a nil error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OnError was never called for an undecodable notification")
	}
}

func TestServerErrorSurfacesAsError(t *testing.T) {
	server := newFakeServer(t, func(string) (any, error) {
		return nil, errServerRejected
	})
	client := New(server.addr(), "secret")

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	_, err := client.GetPlayers(context.Background())

	var callErr *Error
	if !errors.As(err, &callErr) {
		t.Fatalf("errors.As(err, *mcrpc.Error) = false, err = %v", err)
	}
	if callErr.Code != -32602 {
		t.Errorf("Code = %d, want -32602", callErr.Code)
	}
	if callErr.Method == "" {
		t.Error("Method should name the failing call")
	}

	var rpcErr *jsonrpc2.Error
	if !errors.As(err, &rpcErr) {
		t.Error("the underlying jsonrpc2.Error should stay reachable")
	}
}

func TestListParametersSerialiseAsArrays(t *testing.T) {
	// A variadic call with no arguments must not send null where the protocol
	// expects a list.
	params := struct {
		Players []Player `json:"players"`
	}{Players: nonNilSlice[Player](nil)}

	encoded, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	if string(encoded) != `{"players":[]}` {
		t.Errorf("got %s, want {\"players\":[]}", encoded)
	}
}
