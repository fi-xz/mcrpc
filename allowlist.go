// Package mcrpc provides allowlist management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetAllowlist retrieves the current allowlist of players.
func (c *Client) GetAllowlist(ctx context.Context) ([]Player, error) {
	var allowlistPlayers []Player
	err := c.call(ctx, protocol.MethodAllowlistGet, nil, &allowlistPlayers)
	return allowlistPlayers, err
}

// SetAllowlist sets the allowlist to the specified list of players, replacing the existing list.
func (c *Client) SetAllowlist(ctx context.Context, players []Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.SetAllowlistParams{Allowlist: nonNilSlice(players)}
	err := c.call(ctx, protocol.MethodAllowlistSet, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// AddAllowlist adds the specified players to the allowlist.
func (c *Client) AddAllowlist(ctx context.Context, add ...Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.AddAllowlistParams{AllowAdd: nonNilSlice(add)}
	err := c.call(ctx, protocol.MethodAllowlistAdd, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// RemoveAllowlist removes the specified players from the allowlist.
func (c *Client) RemoveAllowlist(ctx context.Context, remove ...Player) ([]Player, error) {
	var updatedAllowlist []Player
	params := protocol.RemoveAllowlistParams{AllowRemove: nonNilSlice(remove)}
	err := c.call(ctx, protocol.MethodAllowlistRemove, params, &updatedAllowlist)
	return updatedAllowlist, err
}

// ClearAllowlist removes all players from the allowlist.
func (c *Client) ClearAllowlist(ctx context.Context) ([]Player, error) {
	var updatedAllowlist []Player
	err := c.call(ctx, protocol.MethodAllowlistClear, nil, &updatedAllowlist)
	return updatedAllowlist, err
}
