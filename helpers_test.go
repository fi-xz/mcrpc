package mcrpc

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMessageConstructorsOmitUnusedFields(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		want    string
	}{
		{
			name:    "literal",
			message: LiteralMessage("See you later"),
			want:    `{"literal":"See you later"}`,
		},
		{
			name:    "translatable",
			message: TranslatableMessage("multiplayer.disconnect.kicked"),
			want:    `{"translatable":"multiplayer.disconnect.kicked"}`,
		},
		{
			name:    "translatable with params",
			message: TranslatableMessage("chat.type.text", "fi_xz", "hello"),
			want:    `{"translatable":"chat.type.text","translatableParams":["fi_xz","hello"]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := json.Marshal(test.message)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(encoded) != test.want {
				t.Errorf("got %s, want %s", encoded, test.want)
			}
		})
	}
}

func TestBanExpiryRoundTrip(t *testing.T) {
	deadline := time.Date(2026, 8, 1, 12, 0, 0, 0, time.UTC)

	ban := UserBan{Player: PlayerByName("fi_xz"), Expires: BanUntil(deadline)}

	got, ok := ban.ExpiresAt()
	if !ok {
		t.Fatalf("ExpiresAt reported no expiry for %q", ban.Expires)
	}
	if !got.Equal(deadline) {
		t.Errorf("got %s, want %s", got, deadline)
	}
	if ban.IsPermanent() {
		t.Error("expected a temporary ban")
	}
}

func TestPermanentBanOmitsExpires(t *testing.T) {
	ban := IPBan{IP: "127.0.0.1", Reason: "griefing", Source: "test"}

	if !ban.IsPermanent() {
		t.Error("expected a permanent ban")
	}
	if _, ok := ban.ExpiresAt(); ok {
		t.Error("expected no expiry timestamp")
	}

	encoded, err := json.Marshal(ban)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if _, present := fields["expires"]; present {
		t.Errorf("permanent ban should omit expires, got %s", encoded)
	}
}

func TestGameRuleAccessors(t *testing.T) {
	t.Run("bool", func(t *testing.T) {
		if v, ok := BoolRule("minecraft:keep_inventory", true).Bool(); !ok || !v {
			t.Errorf("got (%v, %v), want (true, true)", v, ok)
		}
		if _, ok := BoolRule("minecraft:keep_inventory", true).Int(); ok {
			t.Error("boolean rule should not read as an integer")
		}
	})

	t.Run("int", func(t *testing.T) {
		if v, ok := IntRule("minecraft:random_tick_speed", 3).Int(); !ok || v != 3 {
			t.Errorf("got (%v, %v), want (3, true)", v, ok)
		}
	})

	t.Run("string", func(t *testing.T) {
		if v, ok := StringRule("mode", "hard").StringValue(); !ok || v != "hard" {
			t.Errorf("got (%v, %v), want (hard, true)", v, ok)
		}
	})

	// Values arriving over the wire are decoded into any, so integers show up
	// as float64, and a value may also arrive as a string.
	t.Run("decoded from json", func(t *testing.T) {
		var rule TypedGameRule
		if err := json.Unmarshal([]byte(`{"key":"minecraft:random_tick_speed","value":3,"type":"integer"}`), &rule); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if v, ok := rule.Int(); !ok || v != 3 {
			t.Errorf("got (%v, %v), want (3, true)", v, ok)
		}

		var stringly TypedGameRule
		if err := json.Unmarshal([]byte(`{"key":"minecraft:keep_inventory","value":"true","type":"boolean"}`), &stringly); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if v, ok := stringly.Bool(); !ok || !v {
			t.Errorf("got (%v, %v), want (true, true)", v, ok)
		}
	})

	t.Run("non-integral float is rejected", func(t *testing.T) {
		rule := UntypedGameRule{Key: "minecraft:respawn_radius", Value: 1.5}
		if _, ok := rule.Int(); ok {
			t.Error("1.5 should not read as an integer")
		}
	})
}
