// Package mcrpc provides server management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// APIVersion returns the version of the management API the server implements,
// as advertised by rpc.discover.
//
// This is the version protocol behaviour is keyed on, not the Minecraft
// version: 1.0.0 ships with 1.21.9, 3.0.0 with 26.2, and 3.1.0 with 26.3, which
// is what added the world upgrade notifications. Game rule values, for example,
// are strings before 3.0.0 and native JSON types from 3.0.0 on.
func (c *Client) APIVersion(ctx context.Context) (string, error) {
	var schema struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}

	if err := c.call(ctx, protocol.MethodDiscover, nil, &schema); err != nil {
		return "", err
	}

	return schema.Info.Version, nil
}

// GetServerStatus retrieves the current status of the server including online players and version information.
func (c *Client) GetServerStatus(ctx context.Context) (ServerState, error) {
	var status ServerState
	err := c.call(ctx, protocol.MethodServerStatus, nil, &status)
	return status, err
}

// SaveServer saves the current server state. If flush is true, all pending changes are flushed to disk.
func (c *Client) SaveServer(ctx context.Context, flush bool) (bool, error) {
	var saving bool
	params := protocol.ServerSaveParams{Flush: flush}
	err := c.call(ctx, protocol.MethodServerSave, params, &saving)
	return saving, err
}

// StopServer stops the Minecraft server.
func (c *Client) StopServer(ctx context.Context) (bool, error) {
	var stopping bool
	err := c.call(ctx, protocol.MethodServerStop, nil, &stopping)
	return stopping, err
}

// SendSystemMessage sends a system message to the specified players on the server.
func (c *Client) SendSystemMessage(ctx context.Context, message SystemMessage) (bool, error) {
	var sent bool
	params := protocol.SystemMessageParams{Message: message}
	err := c.call(ctx, protocol.MethodServerSystemMessage, params, &sent)
	return sent, err
}
