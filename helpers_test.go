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
	// Keys here are the 1.21.11+ registry form. Older servers use camelCase
	// (keepInventory); the library passes either through unchanged.
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

	// A live server sends real JSON types: true for a boolean rule, 3 for an
	// integer one. Decoded into any, an integer therefore arrives as float64.
	// The string cases below are defensive only.
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

func TestPlayerConstructors(t *testing.T) {
	if got := PlayerByName("fi_xz"); got.Name != "fi_xz" || got.UUID != "" {
		t.Errorf("PlayerByName = %+v", got)
	}

	const uuid = "a0d8c884-2a79-4c95-8617-a51d27a427ec"
	if got := PlayerByUUID(uuid); got.UUID != uuid || got.Name != "" {
		t.Errorf("PlayerByUUID = %+v", got)
	}
}

// TestGameRuleValueRepresentations covers the split between servers: values are
// strings before management API 3.0.0 and native JSON types from 3.0.0 on. Both
// must read, and an update must go back in the representation it came in.
func TestGameRuleValueRepresentations(t *testing.T) {
	stringly := TypedGameRule{
		UntypedGameRule: UntypedGameRule{Key: "keepInventory", Value: "true"},
		Type:            "boolean",
	}
	native := TypedGameRule{
		UntypedGameRule: UntypedGameRule{Key: "minecraft:keep_inventory", Value: true},
		Type:            "boolean",
	}

	if !stringly.UsesStringValues() {
		t.Error("a string value should report UsesStringValues")
	}
	if native.UsesStringValues() {
		t.Error("a native value should not report UsesStringValues")
	}

	if got := stringly.WithBool(false); got.Value != "false" {
		t.Errorf("WithBool on a string rule = %#v, want \"false\"", got.Value)
	}
	if got := native.WithBool(false); got.Value != false {
		t.Errorf("WithBool on a native rule = %#v, want false", got.Value)
	}

	counted := TypedGameRule{
		UntypedGameRule: UntypedGameRule{Key: "randomTickSpeed", Value: "3"},
		Type:            "integer",
	}
	if got := counted.WithInt(5); got.Value != "5" {
		t.Errorf("WithInt on a string rule = %#v, want \"5\"", got.Value)
	}
	if got := (UntypedGameRule{Key: "k", Value: 3}).WithInt(5); got.Value != 5 {
		t.Errorf("WithInt on a native rule = %#v, want 5", got.Value)
	}
	if got := counted.WithString("7"); got.Value != "7" {
		t.Errorf("WithString = %#v, want \"7\"", got.Value)
	}
}

func TestGameRuleAccessorsAcceptEveryRepresentation(t *testing.T) {
	bools := []struct {
		value any
		want  bool
		ok    bool
	}{
		{true, true, true},
		{"true", true, true},
		{"false", false, true},
		{"maybe", false, false},
		{3, false, false},
	}
	for _, test := range bools {
		got, ok := UntypedGameRule{Value: test.value}.Bool()
		if got != test.want || ok != test.ok {
			t.Errorf("Bool(%#v) = (%v, %v), want (%v, %v)", test.value, got, ok, test.want, test.ok)
		}
	}

	ints := []struct {
		value any
		want  int
		ok    bool
	}{
		{3, 3, true},
		{int64(3), 3, true},
		{float64(3), 3, true},
		{"3", 3, true},
		{json.Number("3"), 3, true},
		{"three", 0, false},
		{1.5, 0, false},
		{true, 0, false},
	}
	for _, test := range ints {
		got, ok := UntypedGameRule{Value: test.value}.Int()
		if got != test.want || ok != test.ok {
			t.Errorf("Int(%#v) = (%v, %v), want (%v, %v)", test.value, got, ok, test.want, test.ok)
		}
	}

	if _, ok := (UntypedGameRule{Value: 3}).StringValue(); ok {
		t.Error("StringValue should reject a non-string value")
	}
}

func TestExpiresAtRejectsAnUnparseableTimestamp(t *testing.T) {
	ban := UserBan{Expires: "not a timestamp"}

	if _, ok := ban.ExpiresAt(); ok {
		t.Error("an unparseable expiry should not be reported as a deadline")
	}
	if ban.IsPermanent() {
		t.Error("a ban with a malformed expiry is not permanent, just unreadable")
	}
}

func TestNonNilSlicePassesAPopulatedSliceThrough(t *testing.T) {
	players := []Player{PlayerByName("fi_xz")}

	got := nonNilSlice(players)
	if len(got) != 1 || got[0].Name != "fi_xz" {
		t.Errorf("nonNilSlice altered a populated slice: %+v", got)
	}
}
