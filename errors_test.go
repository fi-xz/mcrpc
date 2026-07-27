package mcrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/sourcegraph/jsonrpc2"
)

func TestCallWithoutConnectionReportsNotConnected(t *testing.T) {
	client := &Client{}

	err := client.call(context.Background(), "minecraft:players", nil, nil)
	if err == nil {
		t.Fatal("expected an error from a client with no connection")
	}
	if !errors.Is(err, ErrNotConnected) {
		t.Errorf("errors.Is(err, ErrNotConnected) = false, err = %v", err)
	}

	var callErr *Error
	if !errors.As(err, &callErr) {
		t.Fatalf("errors.As(err, *mcrpc.Error) = false, err = %v", err)
	}
	if callErr.Method != "minecraft:players" {
		t.Errorf("Method = %q, want %q", callErr.Method, "minecraft:players")
	}
}

func TestErrorExposesServerCode(t *testing.T) {
	rpcErr := &jsonrpc2.Error{Code: -32602, Message: "invalid params"}
	err := error(&Error{Method: "minecraft:players/kick", Code: rpcErr.Code, Message: rpcErr.Message, err: rpcErr})

	var callErr *Error
	if !errors.As(err, &callErr) {
		t.Fatalf("errors.As(err, *mcrpc.Error) = false, err = %v", err)
	}
	if callErr.Code != -32602 {
		t.Errorf("Code = %d, want -32602", callErr.Code)
	}

	var unwrapped *jsonrpc2.Error
	if !errors.As(err, &unwrapped) {
		t.Error("expected the underlying jsonrpc2.Error to remain reachable")
	}

	want := "mcrpc: minecraft:players/kick failed: invalid params (code -32602)"
	if callErr.Error() != want {
		t.Errorf("Error() = %q, want %q", callErr.Error(), want)
	}
}

func TestErrorWithoutAServerMessage(t *testing.T) {
	// A transport failure produces no JSON-RPC message, so the wrapped cause is
	// what the string has to carry.
	err := &Error{Method: "minecraft:players", err: ErrNotConnected}

	got := err.Error()
	want := "mcrpc: minecraft:players failed: mcrpc: client is not connected"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if err.Code != 0 {
		t.Errorf("Code = %d, want 0 for a transport failure", err.Code)
	}
}

func TestCallOnANilClient(t *testing.T) {
	var client *Client

	err := client.call(context.Background(), "minecraft:players", nil, nil)
	if !errors.Is(err, ErrNotConnected) {
		t.Errorf("got %v, want ErrNotConnected", err)
	}
}
