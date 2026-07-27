// Package mcrpc provides player management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetPlayers retrieves the list of currently online players.
func (c *MCRPCClient) GetPlayers(ctx context.Context) ([]Player, error) {
	var players []Player
	err := c.JSONRPCConn.Call(ctx, protocol.MethodPlayersGet, nil, &players)
	return players, err
}

// KickPlayers kicks the specified players from the server with custom messages.
func (c *MCRPCClient) KickPlayers(ctx context.Context, kick []KickPlayer) ([]Player, error) {
	var kicked []Player
	params := protocol.KickPlayerParams{KickPlayers: kick}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodPlayersKick, params, &kicked)
	return kicked, err
}
