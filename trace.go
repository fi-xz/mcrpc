package mcrpc

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

// TraceDirection reports which way a traced message travelled.
type TraceDirection int

const (
	// TraceSent marks a message the client wrote to the server.
	TraceSent TraceDirection = iota
	// TraceReceived marks a message the client read from the server.
	TraceReceived
)

// String renders the direction as an arrow, for log output.
func (d TraceDirection) String() string {
	if d == TraceSent {
		return "->"
	}
	return "<-"
}

// TraceMessage is one JSON-RPC message observed on the connection, carrying
// the params and result exactly as they were serialised. It is intended for
// diagnosing protocol mismatches: what a struct tag actually produces, and
// what JSON type a field actually arrives as.
type TraceMessage struct {
	// Direction is TraceSent or TraceReceived.
	Direction TraceDirection

	// Method names the call or notification. It is empty for a response,
	// unless the originating request is still known, in which case it names
	// that request's method.
	Method string

	// ID is the JSON-RPC request id, empty for a notification.
	ID string

	// Notification reports whether this message expects no response.
	Notification bool

	// Params is the request payload as sent or received, nil when absent.
	Params json.RawMessage

	// Result is the response payload, nil for requests and for error
	// responses.
	Result json.RawMessage

	// ErrorCode and ErrorMessage describe an error response. ErrorCode is 0
	// when the message is not an error.
	ErrorCode    int64
	ErrorMessage string
}

// IsError reports whether the message is an error response.
func (m TraceMessage) IsError() bool {
	return m.ErrorMessage != "" || m.ErrorCode != 0
}

// String renders the message as one line, suitable for passing straight to a
// logger:
//
//	-> minecraft:allowlist/add #3 {"add":[{"name":"fi_xz","id":""}]}
//	<- #3 [{"name":"fi_xz","id":"a0d8c884-..."}]
//	<- notify minecraft:notification/players/joined {"name":"fi_xz",...}
func (m TraceMessage) String() string {
	var line strings.Builder

	line.WriteString(m.Direction.String())

	if m.Notification {
		line.WriteString(" notify")
	}
	if m.Method != "" {
		line.WriteString(" " + m.Method)
	}
	if m.ID != "" {
		line.WriteString(" #" + m.ID)
	}

	switch {
	case m.IsError():
		line.WriteString(" error " + strconv.FormatInt(m.ErrorCode, 10) + " " + m.ErrorMessage)
	case m.Params != nil:
		line.WriteString(" " + string(m.Params))
	case m.Result != nil:
		line.WriteString(" " + string(m.Result))
	}

	return line.String()
}

// newTraceMessage flattens the request/response pair jsonrpc2 reports into a
// single observation. Exactly one of req and resp is set for a request; both
// are set for a response whose originating request is still known.
func newTraceMessage(direction TraceDirection, req *jsonrpc2.Request, resp *jsonrpc2.Response) TraceMessage {
	message := TraceMessage{Direction: direction}

	if req != nil {
		message.Method = req.Method
		message.Notification = req.Notif
		if req.Params != nil {
			message.Params = *req.Params
		}
		if !req.Notif {
			message.ID = req.ID.String()
		}
	}

	if resp != nil {
		message.ID = resp.ID.String()
		message.Params = nil
		if resp.Result != nil {
			message.Result = *resp.Result
		}
		if resp.Error != nil {
			message.ErrorCode = resp.Error.Code
			message.ErrorMessage = resp.Error.Message
		}
	}

	return message
}

// traceOptions builds the jsonrpc2 hooks that feed the configured tracer, or
// nil when tracing is off.
func (c *Client) traceOptions() []jsonrpc2.ConnOpt {
	if c.trace == nil {
		return nil
	}

	return []jsonrpc2.ConnOpt{
		jsonrpc2.OnSend(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			c.trace(newTraceMessage(TraceSent, req, resp))
		}),
		jsonrpc2.OnRecv(func(req *jsonrpc2.Request, resp *jsonrpc2.Response) {
			c.trace(newTraceMessage(TraceReceived, req, resp))
		}),
	}
}
