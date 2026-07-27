// Package mcrpc provides ban management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetBanlist retrieves the current list of banned players.
func (c *MCRPCClient) GetBanlist(ctx context.Context) ([]UserBan, error) {
	var banlistPlayers []UserBan
	err := c.JSONRPCConn.Call(ctx, protocol.MethodBansGet, nil, &banlistPlayers)
	return banlistPlayers, err
}

// SetBanlist sets the ban list to the specified list of bans, replacing the existing list.
func (c *MCRPCClient) SetBanlist(ctx context.Context, bans []UserBan) ([]UserBan, error) {
	var updatedBanlist []UserBan
	params := protocol.SetBanlistParams{Banlist: bans}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodBansSet, params, &updatedBanlist)
	return updatedBanlist, err
}

// AddBanlist adds the specified bans to the ban list.
func (c *MCRPCClient) AddBanlist(ctx context.Context, add []UserBan) ([]UserBan, error) {
	var updatedBanlist []UserBan
	params := protocol.AddBanlistParams{BanAdd: add}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodBansAdd, params, &updatedBanlist)
	return updatedBanlist, err
}

// RemoveBanlist removes the specified players from the ban list.
func (c *MCRPCClient) RemoveBanlist(ctx context.Context, remove []Player) ([]UserBan, error) {
	var updatedBanlist []UserBan
	params := protocol.RemoveBanlistParams{BanRemove: remove}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodBansRemove, params, &updatedBanlist)
	return updatedBanlist, err
}

// ClearBanlist removes all bans from the ban list.
func (c *MCRPCClient) ClearBanlist(ctx context.Context) ([]UserBan, error) {
	var updatedBanlist []UserBan
	err := c.JSONRPCConn.Call(ctx, protocol.MethodBansClear, nil, &updatedBanlist)
	return updatedBanlist, err
}

// GetIPBanlist retrieves the current list of banned IP addresses.
func (c *MCRPCClient) GetIPBanlist(ctx context.Context) ([]IPBan, error) {
	var ipBanlist []IPBan
	err := c.JSONRPCConn.Call(ctx, protocol.MethodIPBansGet, nil, &ipBanlist)
	return ipBanlist, err
}

// SetIPBanlist sets the IP ban list to the specified list, replacing the existing list.
func (c *MCRPCClient) SetIPBanlist(ctx context.Context, banlist []IPBan) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	params := protocol.SetIPBanlistParams{IPBanlist: banlist}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodIPBansSet, params, &updatedIPBanlist)
	return updatedIPBanlist, err
}

// AddIPBanlist adds the specified IP bans to the ban list.
func (c *MCRPCClient) AddIPBanlist(ctx context.Context, add []IncomingIPBan) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	params := protocol.AddIPBanlistParams{IPBanAdd: add}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodIPBansAdd, params, &updatedIPBanlist)
	return updatedIPBanlist, err
}

// RemoveIPBanlist removes the specified IP addresses from the ban list.
func (c *MCRPCClient) RemoveIPBanlist(ctx context.Context, ip []string) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	params := protocol.RemoveIPBanlistParams{IPBanRemove: ip}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodIPBansRemove, params, &updatedIPBanlist)
	return updatedIPBanlist, err
}

// ClearIPBanlist removes all IP bans from the ban list.
func (c *MCRPCClient) ClearIPBanlist(ctx context.Context) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	err := c.JSONRPCConn.Call(ctx, protocol.MethodIPBansClear, nil, &updatedIPBanlist)
	return updatedIPBanlist, err
}
