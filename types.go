// Package mcrpc re-exports the protocol data types so that callers never need
// to reference the internal packages that define them.
package mcrpc

import "github.com/fi-xz/mcrpc/internal/types"

// Player is an alias for types.Player, representing a Minecraft player.
type Player = types.Player

// Operator is an alias for types.Operator, representing a server operator.
type Operator = types.Operator

// ServerState is an alias for types.ServerState, representing the current state of the server.
type ServerState = types.ServerState

// Version is an alias for types.Version, representing the server version information.
type Version = types.Version

// Message is an alias for types.Message, representing a translatable or literal
// message that can be sent to players.
type Message = types.Message

// SystemMessage is an alias for types.SystemMessage, representing a message
// broadcast to a set of players.
type SystemMessage = types.SystemMessage

// KickPlayer is an alias for types.KickPlayer, representing a player to be kicked with a message.
type KickPlayer = types.KickPlayer

// UserBan is an alias for types.UserBan, representing a banned player.
type UserBan = types.UserBan

// IPBan is an alias for types.IPBan, representing a banned IP address.
type IPBan = types.IPBan

// IncomingIPBan is an alias for types.IncomingIPBan, representing an incoming IP ban with player information.
type IncomingIPBan = types.IncomingIPBan

// UntypedGameRule is an alias for types.UntypedGameRule, representing a game
// rule value whose type is not declared.
type UntypedGameRule = types.UntypedGameRule

// TypedGameRule is an alias for types.TypedGameRule, representing a game rule with a typed value.
type TypedGameRule = types.TypedGameRule
