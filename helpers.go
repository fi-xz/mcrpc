package mcrpc

import (
	"time"

	"github.com/fi-xz/mcrpc/internal/types"
)

// BanTimeFormat is the layout the management protocol uses for ban expiry
// timestamps.
const BanTimeFormat = types.BanTimeFormat

// LiteralMessage builds a message carrying plain text.
func LiteralMessage(text string) Message {
	return Message{Literal: text}
}

// TranslatableMessage builds a message the client renders from a translation
// key, substituting params in order.
func TranslatableMessage(key string, params ...string) Message {
	return Message{Translatable: key, TranslatableParams: params}
}

// PlayerByName identifies a player by username, for calls where the server
// resolves the UUID itself.
func PlayerByName(name string) Player {
	return Player{Name: name}
}

// PlayerByUUID identifies a player by UUID, for calls where the username is
// unknown or may have changed.
func PlayerByUUID(uuid string) Player {
	return Player{UUID: uuid}
}

// BanUntil formats t as a ban expiry timestamp. Pass the result as the Expires
// field of a UserBan or IPBan; leave that field empty for a permanent ban.
func BanUntil(t time.Time) string {
	return t.Format(BanTimeFormat)
}

// BoolRule builds a game rule update carrying a boolean value.
func BoolRule(key string, value bool) UntypedGameRule {
	return UntypedGameRule{Key: key, Value: value}
}

// IntRule builds a game rule update carrying an integer value.
func IntRule(key string, value int) UntypedGameRule {
	return UntypedGameRule{Key: key, Value: value}
}

// StringRule builds a game rule update carrying a string value.
func StringRule(key, value string) UntypedGameRule {
	return UntypedGameRule{Key: key, Value: value}
}

// nonNilSlice substitutes an empty slice for a nil one, so that a list-valued
// request parameter serialises as [] rather than null. A variadic call with no
// arguments, and an explicit "replace the list with nothing", both produce nil.
func nonNilSlice[T any](values []T) []T {
	if values == nil {
		return []T{}
	}
	return values
}
