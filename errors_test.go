package mcrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/sourcegraph/jsonrpc2"
)

func TestCallWithoutConnectionReportsNotConnected(t *testing.T) {
	client := &MCRPCClient{}

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
