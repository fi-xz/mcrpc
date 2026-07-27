// Package mcrpc provides allowlist management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetAllowlist retrieves the current allowlist of players.
func (c *MCRPCClient) GetAllowlist(ctx context.Context) ([]Player, error) {
	var allowlistPlayers []Player
	err := c.JSONRPCConn.Call(ctx, protocol.MethodAllowlistGet, nil, &allowlistPlayers)
	return allowlistPlayers, err
}

// SetAllowlist sets the allowlist to the specified list of players, replacing the existing list.
func (c *MCRPCClient) SetAllowlist(ctx context.Context, players []Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.SetAllowlistParams{Allowlist: players}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodAllowlistSet, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// AddAllowlist adds the specified players to the allowlist.
func (c *MCRPCClient) AddAllowlist(ctx context.Context, add []Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.AddAllowlistParams{AllowAdd: add}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodAllowlistAdd, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// RemoveAllowlist removes the specified players from the allowlist.
func (c *MCRPCClient) RemoveAllowlist(ctx context.Context, remove []Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.RemoveAllowlistParams{AllowRemove: remove}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodAllowlistRemove, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// ClearAllowlist removes all players from the allowlist.
func (c *MCRPCClient) ClearAllowlist(ctx context.Context) ([]Player, error) {
	var updatedAllowlist []Player
	err := c.JSONRPCConn.Call(ctx, protocol.MethodAllowlistClear, nil, &updatedAllowlist)
	return updatedAllowlist, err
}
