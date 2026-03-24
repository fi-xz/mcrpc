// Package mcrpc provides player management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
	"github.com/fi-xz/mcrpc/internal/types"
)

// KickPlayer is an alias for types.KickPlayer, representing a player to be kicked with a message.
type KickPlayer = types.KickPlayer

// GetPlayers retrieves the list of currently online players.
func (c *MCRPCClient) GetPlayers(context context.Context) ([]Player, error) {
	var players []Player
	err := c.JSONRPCConn.Call(context, protocol.MethodPlayersGet, nil, &players)
	return players, err
}

// KickPlayers kicks the specified players from the server with custom messages.
func (c *MCRPCClient) KickPlayers(context context.Context, kick []KickPlayer) ([]Player, error) {
	var kicked []Player
	params := protocol.KickPlayerParams{KickPlayers: kick}
	err := c.JSONRPCConn.Call(context, protocol.MethodPlayersKick, params, &kicked)
	return kicked, err
}
