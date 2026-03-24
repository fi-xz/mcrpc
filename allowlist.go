// Package mcrpc provides allowlist management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
	"github.com/fi-xz/mcrpc/internal/types"
)

// Player is an alias for types.Player, representing a Minecraft player.
type Player = types.Player

// GetAllowlist retrieves the current allowlist of players.
func (c *MCRPCClient) GetAllowlist(context context.Context) ([]Player, error) {
	var allowlistPlayers []Player
	err := c.JSONRPCConn.Call(context, protocol.MethodAllowlistGet, nil, &allowlistPlayers)
	return allowlistPlayers, err
}

// SetAllowlist sets the allowlist to the specified list of players, replacing the existing list.
func (c *MCRPCClient) SetAllowlist(context context.Context, players []Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.SetAllowlistParams{Allowlist: players}
	err := c.JSONRPCConn.Call(context, protocol.MethodAllowlistSet, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// AddAllowlist adds the specified players to the allowlist.
func (c *MCRPCClient) AddAllowlist(context context.Context, add []Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.AddAllowlistParams{AllowAdd: add}
	err := c.JSONRPCConn.Call(context, protocol.MethodAllowlistAdd, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// RemoveAllowlist removes the specified players from the allowlist.
func (c *MCRPCClient) RemoveAllowlist(context context.Context, remove []Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.RemoveAllowlistParams{AllowRemove: remove}
	err := c.JSONRPCConn.Call(context, protocol.MethodAllowlistRemove, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// ClearAllowlist removes all players from the allowlist.
func (c *MCRPCClient) ClearAllowlist(context context.Context) ([]Player, error) {
	var updatedAllowlist []Player
	err := c.JSONRPCConn.Call(context, protocol.MethodAllowlistClear, nil, &updatedAllowlist)
	return updatedAllowlist, err
}
