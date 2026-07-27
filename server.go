// Package mcrpc provides server management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetServerStatus retrieves the current status of the server including online players and version information.
func (c *MCRPCClient) GetServerStatus(ctx context.Context) (ServerState, error) {
	var status ServerState
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerStatus, nil, &status)
	return status, err
}

// SaveServer saves the current server state. If flush is true, all pending changes are flushed to disk.
func (c *MCRPCClient) SaveServer(ctx context.Context, flush bool) (bool, error) {
	var saving bool
	params := protocol.ServerSaveParams{Flush: flush}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSave, params, &saving)
	return saving, err
}

// StopServer stops the Minecraft server.
func (c *MCRPCClient) StopServer(ctx context.Context) (bool, error) {
	var stopping bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerStop, nil, &stopping)
	return stopping, err
}

// SendSystemMessage sends a system message to the specified players on the server.
func (c *MCRPCClient) SendSystemMessage(ctx context.Context, message SystemMessage) (bool, error) {
	var sent bool
	params := protocol.SystemMessageParams{Message: message}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSystemMessage, params, &sent)
	return sent, err
}
