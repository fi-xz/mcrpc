package mcrpc

import (
	"context"
	"testing"
	"time"
)

// testBanDeadline is a stable, whole-second time in the future, so that a
// round-trip comparison is not affected by sub-second truncation.
func testBanDeadline() time.Time {
	return time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)
}

// The tests in this file assert facts about the wire format that only a live
// server can settle: whether an omitted field means what we assume, what JSON
// types untyped values actually use, and whether an empty list is accepted.
// They skip without a server, so treat a skipped run as "unverified", not
// "passing".

// restoreBanlist snapshots the ban list and puts it back when the test ends.
func restoreBanlist(t *testing.T, client *Client, ctx context.Context) {
	t.Helper()

	original, err := client.GetBanlist(ctx)
	if err != nil {
		t.Fatalf("GetBanlist failed: %v", err)
	}

	t.Cleanup(func() {
		if _, err := client.SetBanlist(ctx, original); err != nil {
			t.Errorf("could not restore the ban list: %v", err)
		}
	})
}

// TestWirePermanentBanOmitsExpires checks the assumption behind the omitempty
// on UserBan.Expires: that a ban sent without the field is stored as permanent
// rather than rejected or read as an empty timestamp.
func TestWirePermanentBanOmitsExpires(t *testing.T) {
	client, ctx := createTestClient(t)
	restoreBanlist(t, client, ctx)

	banned := UserBan{
		Player: Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
		Reason: "wire format check",
		Source: "mcrpc test",
	}

	updated, err := client.AddBanlist(ctx, banned)
	if err != nil {
		t.Fatalf("AddBanlist with no expiry failed: %v", err)
	}

	var stored *UserBan
	for i := range updated {
		if updated[i].Player.Name == banned.Player.Name {
			stored = &updated[i]
			break
		}
	}
	if stored == nil {
		t.Fatalf("the ban was not in the returned list: %+v", updated)
	}

	if !stored.IsPermanent() {
		t.Errorf("a ban sent without expires came back with Expires = %q, want it empty", stored.Expires)
	}
	if _, ok := stored.ExpiresAt(); ok {
		t.Error("a permanent ban should report no expiry timestamp")
	}
}

// TestWireTemporaryBanRoundTrips checks that BanUntil's format is the one the
// server stores and returns.
func TestWireTemporaryBanRoundTrips(t *testing.T) {
	client, ctx := createTestClient(t)
	restoreBanlist(t, client, ctx)

	deadline := testBanDeadline()

	banned := UserBan{
		Player:  Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
		Reason:  "wire format check",
		Source:  "mcrpc test",
		Expires: BanUntil(deadline),
	}

	updated, err := client.AddBanlist(ctx, banned)
	if err != nil {
		t.Fatalf("AddBanlist failed: %v", err)
	}

	for _, stored := range updated {
		if stored.Player.Name != banned.Player.Name {
			continue
		}

		parsed, ok := stored.ExpiresAt()
		if !ok {
			t.Fatalf("could not parse the expiry the server returned: %q", stored.Expires)
		}
		if !parsed.Equal(deadline) {
			t.Errorf("expiry round-tripped as %s, want %s", parsed, deadline)
		}
		return
	}

	t.Fatalf("the ban was not in the returned list: %+v", updated)
}

// TestWireEmptyListIsAccepted checks that a list-valued parameter serialised as
// [] is accepted, which is what nonNilSlice guarantees we send.
func TestWireEmptyListIsAccepted(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetAllowlist(ctx)
	if err != nil {
		t.Fatalf("GetAllowlist failed: %v", err)
	}
	t.Cleanup(func() {
		if _, err := client.SetAllowlist(ctx, original); err != nil {
			t.Errorf("could not restore the allowlist: %v", err)
		}
	})

	// A nil slice reaches the wire as [] rather than null.
	if _, err := client.SetAllowlist(ctx, nil); err != nil {
		t.Errorf("SetAllowlist with an empty list failed: %v", err)
	}

	// So does a variadic call with no arguments.
	if _, err := client.AddAllowlist(ctx); err != nil {
		t.Errorf("AddAllowlist with no players failed: %v", err)
	}
}

// TestWireGameRuleValueTypes checks that the accessors match the JSON types the
// server actually sends for each declared rule type.
func TestWireGameRuleValueTypes(t *testing.T) {
	client, ctx := createTestClient(t)

	rules, err := client.GetGameRules(ctx)
	if err != nil {
		t.Fatalf("GetGameRules failed: %v", err)
	}
	if len(rules) == 0 {
		t.Fatal("the server returned no game rules")
	}

	var booleans, integers int
	for _, rule := range rules {
		switch rule.Type {
		case "boolean":
			if _, ok := rule.Bool(); !ok {
				t.Errorf("%s is declared boolean but Bool() rejected %#v", rule.Key, rule.Value)
			}
			booleans++
		case "integer":
			if _, ok := rule.Int(); !ok {
				t.Errorf("%s is declared integer but Int() rejected %#v", rule.Key, rule.Value)
			}
			integers++
		default:
			t.Errorf("%s has unhandled type %q with value %#v", rule.Key, rule.Type, rule.Value)
		}
	}

	if booleans == 0 || integers == 0 {
		t.Errorf("expected both boolean and integer rules, got %d and %d", booleans, integers)
	}
}

// TestWireGameRuleUpdateRoundTrips checks that a rule written with BoolRule
// comes back with the same value and a declared type.
func TestWireGameRuleUpdateRoundTrips(t *testing.T) {
	client, ctx := createTestClient(t)

	rules, err := client.GetGameRules(ctx)
	if err != nil {
		t.Fatalf("GetGameRules failed: %v", err)
	}

	var target TypedGameRule
	for _, rule := range rules {
		if rule.Type == "boolean" {
			target = rule
			break
		}
	}
	if target.Key == "" {
		t.Skip("the server reports no boolean game rule to exercise")
	}

	current, ok := target.Bool()
	if !ok {
		t.Fatalf("%s is declared boolean but holds %#v", target.Key, target.Value)
	}
	t.Cleanup(func() {
		if _, err := client.UpdateGameRule(ctx, target.WithBool(current)); err != nil {
			t.Errorf("could not restore %s: %v", target.Key, err)
		}
	})

	// WithBool matches the representation the server used, which is a string on
	// 1.21.9 and 1.21.10 and a native boolean from 1.21.11 on.
	updated, err := client.UpdateGameRule(ctx, target.WithBool(!current))
	if err != nil {
		t.Fatalf("UpdateGameRule failed: %v", err)
	}

	got, ok := updated.Bool()
	if !ok {
		t.Fatalf("the updated rule holds %#v, which Bool() rejected", updated.Value)
	}
	if got == current {
		t.Errorf("%s still reads %v after being set to %v", target.Key, got, !current)
	}
	if updated.Type != "boolean" {
		t.Errorf("the updated rule came back with type %q, want boolean", updated.Type)
	}
}
