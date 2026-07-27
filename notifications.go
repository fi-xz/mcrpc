// Package mcrpc provides notification handling functionality for Minecraft servers.
package mcrpc

import (
	"context"
	"encoding/json"

	"github.com/fi-xz/mcrpc/internal/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

type handlerFunc func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request)

func (f handlerFunc) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	f(ctx, conn, req)
}

// decodeNotification decodes a notification's payload into target.
func decodeNotification[T any](c *Client, method string, params json.RawMessage, target *T) bool {
	if err := json.Unmarshal(params, target); err != nil {
		c.reportDecodeError(method, err)
		return false
	}

	return true
}

// reportDecodeError surfaces a malformed notification through Handler.OnError
// rather than dropping it silently.
func (c *Client) reportDecodeError(method string, err error) {
	if c.handler.OnError != nil {
		c.handler.OnError(method, err)
	}
}

// handleIncoming creates a JSON-RPC handler that processes incoming notifications from the server.
// It dispatches notifications to the appropriate callbacks on the client's Handler.
func (c *Client) handleIncoming() jsonrpc2.Handler {
	return handlerFunc(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
		if !req.Notif {
			return
		}

		var params json.RawMessage
		if req.Params != nil {
			params = *req.Params
		}

		// Call generic notification handler if set
		if c.handler.OnNotification != nil {
			c.handler.OnNotification(req.Method, params)
		}

		// Handle specific notification types
		switch req.Method {
		// Server notifications
		case protocol.NotificationServerStarted:
			if c.handler.OnServerStarted != nil {
				c.handler.OnServerStarted()
			}

		case protocol.NotificationServerStopping:
			if c.handler.OnServerStopping != nil {
				c.handler.OnServerStopping()
			}

		case protocol.NotificationServerSaving:
			if c.handler.OnServerSaving != nil {
				c.handler.OnServerSaving()
			}

		case protocol.NotificationServerSaved:
			if c.handler.OnServerSaved != nil {
				c.handler.OnServerSaved()
			}

		case protocol.NotificationServerStatus:
			if c.handler.OnServerStatus != nil {
				var status ServerState
				if decodeNotification(c, req.Method, params, &status) {
					c.handler.OnServerStatus(status)
				}
			}

		case protocol.NotificationServerActivity:
			if c.handler.OnServerActivity != nil {
				c.handler.OnServerActivity()
			}

		// Player notifications
		case protocol.NotificationPlayerJoined:
			if c.handler.OnPlayerJoined != nil {
				var player Player
				if decodeNotification(c, req.Method, params, &player) {
					c.handler.OnPlayerJoined(player)
				}
			}

		case protocol.NotificationPlayerLeft:
			if c.handler.OnPlayerLeft != nil {
				var player Player
				if decodeNotification(c, req.Method, params, &player) {
					c.handler.OnPlayerLeft(player)
				}
			}

		// Operator notifications
		case protocol.NotificationOperatorAdded:
			if c.handler.OnOperatorAdded != nil {
				var operator Operator
				if decodeNotification(c, req.Method, params, &operator) {
					c.handler.OnOperatorAdded(operator)
				}
			}

		case protocol.NotificationOperatorRemoved:
			if c.handler.OnOperatorRemoved != nil {
				var operator Operator
				if decodeNotification(c, req.Method, params, &operator) {
					c.handler.OnOperatorRemoved(operator)
				}
			}

		// Allowlist notifications
		case protocol.NotificationAllowlistAdded:
			if c.handler.OnAllowlistAdded != nil {
				var player Player
				if decodeNotification(c, req.Method, params, &player) {
					c.handler.OnAllowlistAdded(player)
				}
			}

		case protocol.NotificationAllowlistRemoved:
			if c.handler.OnAllowlistRemoved != nil {
				var player Player
				if decodeNotification(c, req.Method, params, &player) {
					c.handler.OnAllowlistRemoved(player)
				}
			}

		// Ban notifications
		case protocol.NotificationBanAdded:
			if c.handler.OnBanAdded != nil {
				var ban UserBan
				if decodeNotification(c, req.Method, params, &ban) {
					c.handler.OnBanAdded(ban)
				}
			}

		case protocol.NotificationBanRemoved:
			if c.handler.OnBanRemoved != nil {
				var player Player
				if decodeNotification(c, req.Method, params, &player) {
					c.handler.OnBanRemoved(player)
				}
			}

		// IP Ban notifications
		case protocol.NotificationIPBanAdded:
			if c.handler.OnIPBanAdded != nil {
				var ban IPBan
				if decodeNotification(c, req.Method, params, &ban) {
					c.handler.OnIPBanAdded(ban)
				}
			}

		case protocol.NotificationIPBanRemoved:
			if c.handler.OnIPBanRemoved != nil {
				var ip string
				if decodeNotification(c, req.Method, params, &ip) {
					c.handler.OnIPBanRemoved(ip)
				}
			}

		// Gamerule notifications
		case protocol.NotificationGameruleUpdated:
			if c.handler.OnGameruleUpdated != nil {
				var gamerule TypedGameRule
				if decodeNotification(c, req.Method, params, &gamerule) {
					c.handler.OnGameruleUpdated(gamerule)
				}
			}

		// World Notifications
		case protocol.NotificationWorldUpgradeStarted:
			if c.handler.OnWorldUpgradeStarted != nil {
				c.handler.OnWorldUpgradeStarted()
			}

		case protocol.NotificationWorldUpgradeProgress:
			if c.handler.OnWorldUpgradeProgress != nil {
				var progress float64
				if decodeNotification(c, req.Method, params, &progress) {
					c.handler.OnWorldUpgradeProgress(progress)
				}
			}

		case protocol.NotificationWorldUpgradeFinished:
			if c.handler.OnWorldUpgradeFinished != nil {
				c.handler.OnWorldUpgradeFinished()
			}

		case protocol.NotificationWorldUpgradeFailed:
			if c.handler.OnWorldUpgradeFailed != nil {
				var reason string
				if decodeNotification(c, req.Method, params, &reason) {
					c.handler.OnWorldUpgradeFailed(reason)
				}
			}
		}
	})
}
