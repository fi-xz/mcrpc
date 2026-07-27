// This file lives in the external test package on purpose: it may only use the
// public API of mcrpc, exactly like a downstream consumer. It fails to compile
// if any exported method requires a type that a consumer cannot name, so it
// guards against internal types leaking back into the public surface.
package mcrpc_test

import (
	"testing"

	"github.com/fi-xz/mcrpc"
)

// TestPublicTypesAreConstructible verifies every value an exported method
// accepts as a parameter can be built from outside the module.
func TestPublicTypesAreConstructible(t *testing.T) {
	player := mcrpc.Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}

	message := mcrpc.Message{
		Translatable:       "multiplayer.disconnect.kicked",
		TranslatableParams: []string{},
		Literal:            "See you later",
	}

	values := []any{
		player,
		message,
		mcrpc.Operator{Player: player, PermissionLevel: 4},
		mcrpc.KickPlayer{Player: player, Message: message},
		mcrpc.UserBan{Player: player, Reason: "griefing", Source: "test"},
		mcrpc.IPBan{IP: "127.0.0.1", Reason: "griefing", Source: "test"},
		mcrpc.IncomingIPBan{IPBan: mcrpc.IPBan{IP: "127.0.0.1"}, Player: player},
		mcrpc.SystemMessage{ReceivingPlayers: []mcrpc.Player{player}, Message: message},
		mcrpc.UntypedGameRule{Key: "minecraft:keep_inventory", Value: true},
		mcrpc.TypedGameRule{
			UntypedGameRule: mcrpc.UntypedGameRule{Key: "minecraft:keep_inventory", Value: true},
			Type:            "boolean",
		},
		mcrpc.ServerState{
			Started: true,
			Players: []mcrpc.Player{player},
			Version: mcrpc.Version{Name: "1.21.9", Protocol: 774},
		},
	}

	for _, value := range values {
		if value == nil {
			t.Errorf("unexpected nil value in public type set")
		}
	}
}
