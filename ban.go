// Package mcrpc provides ban management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetBanlist retrieves the current list of banned players.
func (c *MCRPCClient) GetBanlist(context context.Context) ([]UserBan, error) {
	var banlistPlayers []UserBan
	err := c.JSONRPCConn.Call(context, protocol.MethodBansGet, nil, &banlistPlayers)
	return banlistPlayers, err
}

// SetBanlist sets the ban list to the specified list of bans, replacing the existing list.
func (c *MCRPCClient) SetBanlist(context context.Context, bans []UserBan) ([]UserBan, error) {
	var updatedBanlist []UserBan
	params := protocol.SetBanlistParams{Banlist: bans}
	err := c.JSONRPCConn.Call(context, protocol.MethodBansSet, params, &updatedBanlist)
	return updatedBanlist, err
}

// AddBanlist adds the specified bans to the ban list.
func (c *MCRPCClient) AddBanlist(context context.Context, add []UserBan) ([]UserBan, error) {
	var updatedBanlist []UserBan
	params := protocol.AddBanlistParams{BanAdd: add}
	err := c.JSONRPCConn.Call(context, protocol.MethodBansAdd, params, &updatedBanlist)
	return updatedBanlist, err
}

// RemoveBanlist removes the specified players from the ban list.
func (c *MCRPCClient) RemoveBanlist(context context.Context, remove []Player) ([]UserBan, error) {
	var updatedBanlist []UserBan
	params := protocol.RemoveBanlistParams{BanRemove: remove}
	err := c.JSONRPCConn.Call(context, protocol.MethodBansRemove, params, &updatedBanlist)
	return updatedBanlist, err
}

// ClearBanlist removes all bans from the ban list.
func (c *MCRPCClient) ClearBanlist(context context.Context) ([]UserBan, error) {
	var updatedBanlist []UserBan
	err := c.JSONRPCConn.Call(context, protocol.MethodBansClear, nil, &updatedBanlist)
	return updatedBanlist, err
}

// GetIPBanlist retrieves the current list of banned IP addresses.
func (c *MCRPCClient) GetIPBanlist(context context.Context) ([]IPBan, error) {
	var ipBanlist []IPBan
	err := c.JSONRPCConn.Call(context, protocol.MethodIPBansGet, nil, &ipBanlist)
	return ipBanlist, err
}

// SetIPBanlist sets the IP ban list to the specified list, replacing the existing list.
func (c *MCRPCClient) SetIPBanlist(context context.Context, banlist []IPBan) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	params := protocol.SetIPBanlistParams{IPBanlist: banlist}
	err := c.JSONRPCConn.Call(context, protocol.MethodIPBansSet, params, &updatedIPBanlist)
	return updatedIPBanlist, err
}

// AddIPBanlist adds the specified IP bans to the ban list.
func (c *MCRPCClient) AddIPBanlist(context context.Context, add []IncomingIPBan) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	params := protocol.AddIPBanlistParams{IPBanAdd: add}
	err := c.JSONRPCConn.Call(context, protocol.MethodIPBansAdd, params, &updatedIPBanlist)
	return updatedIPBanlist, err
}

// RemoveIPBanlist removes the specified IP addresses from the ban list.
func (c *MCRPCClient) RemoveIPBanlist(context context.Context, ip []string) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	params := protocol.RemoveIPBanlistParams{IPBanRemove: ip}
	err := c.JSONRPCConn.Call(context, protocol.MethodIPBansRemove, params, &updatedIPBanlist)
	return updatedIPBanlist, err
}

// ClearIPBanlist removes all IP bans from the ban list.
func (c *MCRPCClient) ClearIPBanlist(context context.Context) ([]IPBan, error) {
	var updatedIPBanlist []IPBan
	err := c.JSONRPCConn.Call(context, protocol.MethodIPBansClear, nil, &updatedIPBanlist)
	return updatedIPBanlist, err
}
