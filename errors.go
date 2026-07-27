package mcrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
)

// ErrNotConnected is returned when a call is attempted on a client that has
// not been started, or whose session has ended.
var ErrNotConnected = errors.New("mcrpc: client is not connected")

// ErrAlreadyStarted is returned by Start when the client already holds a
// connection. Close it before starting a new session.
var ErrAlreadyStarted = errors.New("mcrpc: client is already started")

// Error describes a failed remote call. Every error returned by a Client
// method wraps the underlying cause, so errors.Is and errors.As both work:
//
//	var rpcErr *mcrpc.Error
//	if errors.As(err, &rpcErr) && rpcErr.Code == -32602 {
//	    // the server rejected the parameters
//	}
type Error struct {
	// Method is the JSON-RPC method that failed.
	Method string
	// Code is the JSON-RPC error code, or 0 if the call failed before the
	// server produced a response (transport errors, closed connection).
	Code int64
	// Message is the error message reported by the server, if any.
	Message string
	// Data is the optional payload the server attached to the error.
	Data *json.RawMessage

	err error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("mcrpc: %s failed: %s (code %d)", e.Method, e.Message, e.Code)
	}
	return fmt.Sprintf("mcrpc: %s failed: %v", e.Method, e.err)
}

// Unwrap returns the underlying cause, exposing it to errors.Is and errors.As.
func (e *Error) Unwrap() error { return e.err }

// call performs a JSON-RPC call and normalises any failure into an *Error.
// It returns a nil error interface on success, never a typed nil.
func (c *Client) call(ctx context.Context, method string, params, result any) error {
	if c == nil {
		return &Error{Method: method, err: ErrNotConnected}
	}

	conn := c.conn()
	if conn == nil {
		return &Error{Method: method, err: ErrNotConnected}
	}

	err := conn.Call(ctx, method, params, result)
	if err == nil {
		return nil
	}

	wrapped := &Error{Method: method, err: err}

	var rpcErr *jsonrpc2.Error
	if errors.As(err, &rpcErr) {
		wrapped.Code = rpcErr.Code
		wrapped.Message = rpcErr.Message
		wrapped.Data = rpcErr.Data
	}

	return wrapped
}
