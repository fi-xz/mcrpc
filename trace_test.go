package mcrpc

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

// traceRecorder collects trace messages from the connection's read and write
// goroutines.
type traceRecorder struct {
	mu       sync.Mutex
	messages []TraceMessage
}

func (r *traceRecorder) record(m TraceMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, m)
}

// find returns the first recorded message satisfying match, waiting briefly
// for the asynchronous read side to catch up.
func (r *traceRecorder) find(t *testing.T, what string, match func(TraceMessage) bool) TraceMessage {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for {
		r.mu.Lock()
		for _, message := range r.messages {
			if match(message) {
				r.mu.Unlock()
				return message
			}
		}
		r.mu.Unlock()

		if time.Now().After(deadline) {
			t.Fatalf("no traced message matched %s; recorded: %s", what, r.dump())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (r *traceRecorder) dump() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	lines := make([]string, 0, len(r.messages))
	for _, message := range r.messages {
		lines = append(lines, message.String())
	}
	return "\n" + strings.Join(lines, "\n")
}

func TestTraceCapturesBothDirections(t *testing.T) {
	server := newFakeServer(t, func(string) (any, error) {
		return []Player{PlayerByName("fi_xz")}, nil
	})

	recorder := &traceRecorder{}
	client := New(server.addr(), "secret", WithTrace(recorder.record))

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	if _, err := client.AddAllowlist(context.Background(), PlayerByName("fi_xz")); err != nil {
		t.Fatalf("AddAllowlist failed: %v", err)
	}

	sent := recorder.find(t, "the outgoing request", func(m TraceMessage) bool {
		return m.Direction == TraceSent && m.Method != ""
	})

	if sent.ID == "" {
		t.Error("a call should be traced with a request id")
	}
	if sent.Notification {
		t.Error("a call should not be traced as a notification")
	}
	// The point of tracing: the params are the exact bytes put on the wire.
	if got := string(sent.Params); !strings.Contains(got, `"fi_xz"`) {
		t.Errorf("traced params = %s, want the player name in it", got)
	}

	received := recorder.find(t, "the response", func(m TraceMessage) bool {
		return m.Direction == TraceReceived && m.Result != nil
	})

	if received.ID != sent.ID {
		t.Errorf("response id = %q, want the request's id %q", received.ID, sent.ID)
	}
	if received.Method != sent.Method {
		t.Errorf("response method = %q, want %q carried over from the request", received.Method, sent.Method)
	}
	if received.Params != nil {
		t.Errorf("a response should not carry request params, got %s", received.Params)
	}
}

func TestTraceCapturesNotifications(t *testing.T) {
	server := newFakeServer(t, okResponder)

	recorder := &traceRecorder{}
	client := New(server.addr(), "secret", WithTrace(recorder.record))

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	conn := server.nextConn(t)
	if err := conn.Notify(context.Background(), "minecraft:notification/players/joined", []any{PlayerByName("fi_xz")}); err != nil {
		t.Fatalf("server failed to notify: %v", err)
	}

	notification := recorder.find(t, "the incoming notification", func(m TraceMessage) bool {
		return m.Direction == TraceReceived && m.Notification
	})

	if notification.ID != "" {
		t.Errorf("a notification should have no id, got %q", notification.ID)
	}
	if notification.Method != "minecraft:notification/players/joined" {
		t.Errorf("method = %q", notification.Method)
	}
	if !strings.Contains(notification.String(), "notify") {
		t.Errorf("String() should mark a notification, got %q", notification.String())
	}
}

func TestTraceCapturesServerErrors(t *testing.T) {
	server := newFakeServer(t, func(string) (any, error) {
		return nil, errServerRejected
	})

	recorder := &traceRecorder{}
	client := New(server.addr(), "secret", WithTrace(recorder.record))

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	if _, err := client.GetPlayers(context.Background()); err == nil {
		t.Fatal("expected the call to fail")
	}

	failure := recorder.find(t, "the error response", func(m TraceMessage) bool {
		return m.Direction == TraceReceived && m.IsError()
	})

	if failure.ErrorCode != -32602 {
		t.Errorf("ErrorCode = %d, want -32602", failure.ErrorCode)
	}
	if !strings.Contains(failure.String(), "error -32602") {
		t.Errorf("String() should render the error, got %q", failure.String())
	}
}

func TestTraceOffByDefault(t *testing.T) {
	server := newFakeServer(t, okResponder)
	client := New(server.addr(), "secret")

	if client.traceOptions() != nil {
		t.Error("a client without WithTrace should install no jsonrpc2 hooks")
	}

	if err := client.Start(context.Background()); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })

	if _, err := client.GetPlayers(context.Background()); err != nil {
		t.Errorf("GetPlayers failed: %v", err)
	}
}

func TestTraceMessageRendering(t *testing.T) {
	tests := []struct {
		name    string
		message TraceMessage
		want    string
	}{
		{
			name:    "outgoing call",
			message: TraceMessage{Direction: TraceSent, Method: "minecraft:players", ID: "1"},
			want:    "-> minecraft:players #1",
		},
		{
			name: "outgoing call with params",
			message: TraceMessage{
				Direction: TraceSent,
				Method:    "minecraft:allowlist/add",
				ID:        "2",
				Params:    []byte(`{"add":[]}`),
			},
			want: `-> minecraft:allowlist/add #2 {"add":[]}`,
		},
		{
			name: "incoming result",
			message: TraceMessage{
				Direction: TraceReceived,
				Method:    "minecraft:players",
				ID:        "1",
				Result:    []byte(`[]`),
			},
			want: "<- minecraft:players #1 []",
		},
		{
			name: "incoming notification",
			message: TraceMessage{
				Direction:    TraceReceived,
				Method:       "minecraft:notification/server/started",
				Notification: true,
			},
			want: "<- notify minecraft:notification/server/started",
		},
		{
			name: "incoming error",
			message: TraceMessage{
				Direction:    TraceReceived,
				Method:       "minecraft:players/kick",
				ID:           "3",
				ErrorCode:    -32602,
				ErrorMessage: "Invalid params",
			},
			want: "<- minecraft:players/kick #3 error -32602 Invalid params",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.message.String(); got != test.want {
				t.Errorf("String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTraceDirectionString(t *testing.T) {
	if got := TraceSent.String(); got != "->" {
		t.Errorf("TraceSent = %q, want ->", got)
	}
	if got := TraceReceived.String(); got != "<-" {
		t.Errorf("TraceReceived = %q, want <-", got)
	}
}

func TestTraceMessageIsError(t *testing.T) {
	if (TraceMessage{}).IsError() {
		t.Error("a message with no error should not report one")
	}
	if !(TraceMessage{ErrorCode: -32601}).IsError() {
		t.Error("a code alone should report an error")
	}
	if !(TraceMessage{ErrorMessage: "boom"}).IsError() {
		t.Error("a message alone should report an error")
	}
}
