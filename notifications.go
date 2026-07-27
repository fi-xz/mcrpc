// Package mcrpc provides notification handling functionality for Minecraft servers.
package mcrpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fi-xz/mcrpc/internal/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

type handlerFunc func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request)

func (f handlerFunc) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	f(ctx, conn, req)
}

// errNoParams reports a notification that should have carried a payload but
// arrived with an empty argument list.
var errNoParams = errors.New("notification carried no parameters")

// decodeParam decodes a notification's payload into target.
//
// JSON-RPC params are a *positional argument list*, and every notification the
// server advertises through rpc.discover declares exactly one parameter. The
// payload is therefore element 0 of an array:
//
//	minecraft:notification/players/joined  ->  [{"name":"fi_xz","id":"…"}]
//
// not the bare object. Decoding the params directly into the payload type
// fails, which is how every payload-carrying handler in this package came to be
// silently dead.
func decodeParam[T any](c *Client, method string, params json.RawMessage, target *T) bool {
	var positional []json.RawMessage
	if err := json.Unmarshal(params, &positional); err != nil {
		c.reportDecodeError(method, err)
		return false
	}

	if len(positional) == 0 {
		c.reportDecodeError(method, errNoParams)
		return false
	}

	if err := json.Unmarshal(positional[0], target); err != nil {
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
				if decodeParam(c, req.Method, params, &status) {
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
				if decodeParam(c, req.Method, params, &player) {
					c.handler.OnPlayerJoined(player)
				}
			}

		case protocol.NotificationPlayerLeft:
			if c.handler.OnPlayerLeft != nil {
				var player Player
				if decodeParam(c, req.Method, params, &player) {
					c.handler.OnPlayerLeft(player)
				}
			}

		// Operator notifications
		case protocol.NotificationOperatorAdded:
			if c.handler.OnOperatorAdded != nil {
				var operator Operator
				if decodeParam(c, req.Method, params, &operator) {
					c.handler.OnOperatorAdded(operator)
				}
			}

		case protocol.NotificationOperatorRemoved:
			if c.handler.OnOperatorRemoved != nil {
				var operator Operator
				if decodeParam(c, req.Method, params, &operator) {
					c.handler.OnOperatorRemoved(operator)
				}
			}

		// Allowlist notifications
		case protocol.NotificationAllowlistAdded:
			if c.handler.OnAllowlistAdded != nil {
				var player Player
				if decodeParam(c, req.Method, params, &player) {
					c.handler.OnAllowlistAdded(player)
				}
			}

		case protocol.NotificationAllowlistRemoved:
			if c.handler.OnAllowlistRemoved != nil {
				var player Player
				if decodeParam(c, req.Method, params, &player) {
					c.handler.OnAllowlistRemoved(player)
				}
			}

		// Ban notifications
		case protocol.NotificationBanAdded:
			if c.handler.OnBanAdded != nil {
				var ban UserBan
				if decodeParam(c, req.Method, params, &ban) {
					c.handler.OnBanAdded(ban)
				}
			}

		case protocol.NotificationBanRemoved:
			if c.handler.OnBanRemoved != nil {
				var player Player
				if decodeParam(c, req.Method, params, &player) {
					c.handler.OnBanRemoved(player)
				}
			}

		// IP Ban notifications
		case protocol.NotificationIPBanAdded:
			if c.handler.OnIPBanAdded != nil {
				var ban IPBan
				if decodeParam(c, req.Method, params, &ban) {
					c.handler.OnIPBanAdded(ban)
				}
			}

		case protocol.NotificationIPBanRemoved:
			if c.handler.OnIPBanRemoved != nil {
				var ip string
				if decodeParam(c, req.Method, params, &ip) {
					c.handler.OnIPBanRemoved(ip)
				}
			}

		// Gamerule notifications
		case protocol.NotificationGameruleUpdated:
			if c.handler.OnGameruleUpdated != nil {
				var gamerule TypedGameRule
				if decodeParam(c, req.Method, params, &gamerule) {
					c.handler.OnGameruleUpdated(gamerule)
				}
			}

		// World notifications, added in API version 3.1.0 (Minecraft 26.3).
		// Servers below that do not advertise them through rpc.discover and
		// never send them.
		//
		// Observed on a 26.3 server converting a world: upgrade_started and
		// upgrade_finished carry no parameters, and upgrade_progress arrives as
		// [0.0] — the number wrapped in the positional argument list, same as
		// every other payload. It is rate limited to one notification per
		// second. upgrade_failed carries a reason string and has not been seen,
		// since provoking a failed conversion was not attempted.
		case protocol.NotificationWorldUpgradeStarted:
			if c.handler.OnWorldUpgradeStarted != nil {
				c.handler.OnWorldUpgradeStarted()
			}

		case protocol.NotificationWorldUpgradeProgress:
			if c.handler.OnWorldUpgradeProgress != nil {
				var progress float64
				if decodeParam(c, req.Method, params, &progress) {
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
				if decodeParam(c, req.Method, params, &reason) {
					c.handler.OnWorldUpgradeFailed(reason)
				}
			}
		}
	})
}
