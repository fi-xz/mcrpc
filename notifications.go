// Package mcrpc provides notification handling functionality for Minecraft servers.
package mcrpc

import (
	"context"
	"encoding/json"

	"github.com/fi-xz/mcrpc/internal/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

// handleIncoming creates a JSON-RPC handler that processes incoming notifications from the server.
// It dispatches notifications to the appropriate handler callbacks set on the MCRPCClient.
func (c *MCRPCClient) handleIncoming() jsonrpc2.Handler {
	return jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result any, err error) {
		if req.Notif {
			var params json.RawMessage
			if req.Params != nil {
				params = *req.Params
			}

			// Call generic notification handler if set
			if c.OnNotification != nil {
				c.OnNotification(req.Method, params)
			}

			// Handle specific notification types
			switch req.Method {
			// Server notifications
			case protocol.NotificationServerStarted:
				if c.OnServerStarted != nil {
					c.OnServerStarted()
				}

			case protocol.NotificationServerStopping:
				if c.OnServerStopping != nil {
					c.OnServerStopping()
				}

			case protocol.NotificationServerSaving:
				if c.OnServerSaving != nil {
					c.OnServerSaving()
				}

			case protocol.NotificationServerSaved:
				if c.OnServerSaved != nil {
					c.OnServerSaved()
				}

			case protocol.NotificationServerStatus:
				if c.OnServerStatus != nil {
					var status ServerState
					if err := json.Unmarshal(params, &status); err == nil {
						c.OnServerStatus(status)
					}
				}

			case protocol.NotificationServerActivity:
				if c.OnServerActivity != nil {
					c.OnServerActivity()
				}

			// Player notifications
			case protocol.NotificationPlayerJoined:
				if c.OnPlayerJoined != nil {
					var player Player
					if err := json.Unmarshal(params, &player); err == nil {
						c.OnPlayerJoined(player)
					}
				}

			case protocol.NotificationPlayerLeft:
				if c.OnPlayerLeft != nil {
					var player Player
					if err := json.Unmarshal(params, &player); err == nil {
						c.OnPlayerLeft(player)
					}
				}

			// Operator notifications
			case protocol.NotificationOperatorAdded:
				if c.OnOperatorAdded != nil {
					var operator Operator
					if err := json.Unmarshal(params, &operator); err == nil {
						c.OnOperatorAdded(operator)
					}
				}

			case protocol.NotificationOperatorRemoved:
				if c.OnOperatorRemoved != nil {
					var operator Operator
					if err := json.Unmarshal(params, &operator); err == nil {
						c.OnOperatorRemoved(operator)
					}
				}

			// Allowlist notifications
			case protocol.NotificationAllowlistAdded:
				if c.OnAllowlistAdded != nil {
					var player Player
					if err := json.Unmarshal(params, &player); err == nil {
						c.OnAllowlistAdded(player)
					}
				}

			case protocol.NotificationAllowlistRemoved:
				if c.OnAllowlistRemoved != nil {
					var player Player
					if err := json.Unmarshal(params, &player); err == nil {
						c.OnAllowlistRemoved(player)
					}
				}

			// Ban notifications
			case protocol.NotificationBanAdded:
				if c.OnBanAdded != nil {
					var ban UserBan
					if err := json.Unmarshal(params, &ban); err == nil {
						c.OnBanAdded(ban)
					}
				}

			case protocol.NotificationBanRemoved:
				if c.OnBanRemoved != nil {
					var player Player
					if err := json.Unmarshal(params, &player); err == nil {
						c.OnBanRemoved(player)
					}
				}

			// IP Ban notifications
			case protocol.NotificationIPBanAdded:
				if c.OnIPBanAdded != nil {
					var ban IPBan
					if err := json.Unmarshal(params, &ban); err == nil {
						c.OnIPBanAdded(ban)
					}
				}

			case protocol.NotificationIPBanRemoved:
				if c.OnIPBanRemoved != nil {
					var ip string
					if err := json.Unmarshal(params, &ip); err == nil {
						c.OnIPBanRemoved(ip)
					}
				}

			// Gamerule notifications
			case protocol.NotificationGameruleUpdated:
				if c.OnGameruleUpdated != nil {
					var gamerule TypedGameRule
					if err := json.Unmarshal(params, &gamerule); err == nil {
						c.OnGameruleUpdated(gamerule)
					}
				}
			}
		}

		return nil, nil
	})
}
