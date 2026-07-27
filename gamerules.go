// Package mcrpc provides game rule management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetGameRules retrieves all game rules and their current values.
func (c *MCRPCClient) GetGameRules(ctx context.Context) ([]TypedGameRule, error) {
	var gamerules []TypedGameRule
	err := c.JSONRPCConn.Call(ctx, protocol.MethodGameRulesGet, nil, &gamerules)
	return gamerules, err
}

// UpdateGameRule updates the value of a specific game rule.
func (c *MCRPCClient) UpdateGameRule(ctx context.Context, gamerule UntypedGameRule) (TypedGameRule, error) {
	var result TypedGameRule
	params := protocol.UpdateGameRulesParams{GameRules: gamerule}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodGameRulesUpdate, params, &result)
	return result, err
}
