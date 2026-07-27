package mcrpc

import (
	"testing"
)

// TestBanlistGet tests getting the ban list
func TestBanlistGet(t *testing.T) {
	client, ctx := createTestClient(t)

	bans, err := client.GetBanlist(ctx)
	if err != nil {
		t.Errorf("GetBanlist failed: %v", err)
	}

	if bans == nil {
		t.Error("Expected non-nil banlist, got nil")
	}
}

// TestBanlistSet tests setting the ban list
func TestBanlistSet(t *testing.T) {
	client, ctx := createTestClient(t)

	// Get current banlist first
	originalBans, err := client.GetBanlist(ctx)
	if err != nil {
		t.Fatalf("Failed to get current banlist: %v", err)
	}

	// Set empty banlist
	updatedBans, err := client.SetBanlist(ctx, []UserBan{})
	if err != nil {
		t.Errorf("SetBanlist failed: %v", err)
	}

	if updatedBans == nil {
		t.Error("Expected non-nil updated banlist")
	}

	// Restore original
	_, err = client.SetBanlist(ctx, originalBans)
	if err != nil {
		t.Errorf("Failed to restore original banlist: %v", err)
	}
}

// TestBanlistAdd tests adding bans
func TestBanlistAdd(t *testing.T) {
	client, ctx := createTestClient(t)

	bansToAdd := []UserBan{
		{
			Player:  Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
			Reason:  "Test ban",
			Expires: "2026-03-24T01:45:59Z",
			Source:  "Test",
		},
	}

	updatedBans, err := client.AddBanlist(ctx, bansToAdd...)
	if err != nil {
		t.Errorf("AddBanlist failed: %v", err)
	}

	if updatedBans == nil {
		t.Error("Expected non-nil updated banlist")
	}
}

// TestBanlistRemove tests removing bans
func TestBanlistRemove(t *testing.T) {
	client, ctx := createTestClient(t)

	playersToRemove := []Player{
		{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
	}

	updatedBans, err := client.RemoveBanlist(ctx, playersToRemove...)
	if err != nil {
		t.Errorf("RemoveBanlist failed: %v", err)
	}

	if updatedBans == nil {
		t.Error("Expected non-nil updated banlist")
	}
}

// TestBanlistClear tests clearing the ban list
func TestBanlistClear(t *testing.T) {
	client, ctx := createTestClient(t)

	clearedBans, err := client.ClearBanlist(ctx)
	if err != nil {
		t.Errorf("ClearBanlist failed: %v", err)
	}

	if clearedBans == nil {
		t.Error("Expected non-nil cleared banlist")
	}
}

// TestIPBanlistGet tests getting the IP ban list
func TestIPBanlistGet(t *testing.T) {
	client, ctx := createTestClient(t)

	bans, err := client.GetIPBanlist(ctx)
	if err != nil {
		t.Errorf("GetIPBanlist failed: %v", err)
	}

	if bans == nil {
		t.Error("Expected non-nil IP banlist, got nil")
	}
}

// TestIPBanlistSet tests setting the IP ban list
func TestIPBanlistSet(t *testing.T) {
	client, ctx := createTestClient(t)

	originalBans, err := client.GetIPBanlist(ctx)
	if err != nil {
		t.Fatalf("Failed to get current IP banlist: %v", err)
	}

	updatedBans, err := client.SetIPBanlist(ctx, []IPBan{})
	if err != nil {
		t.Errorf("SetIPBanlist failed: %v", err)
	}

	if updatedBans == nil {
		t.Error("Expected non-nil updated IP banlist")
	}

	// Restore
	_, err = client.SetIPBanlist(ctx, originalBans)
	if err != nil {
		t.Errorf("Failed to restore original IP banlist: %v", err)
	}
}

// TestIPBanlistAdd tests adding IP bans
func TestIPBanlistAdd(t *testing.T) {
	client, ctx := createTestClient(t)

	bansToAdd := []IncomingIPBan{
		{
			IPBan: IPBan{
				IP:      "192.168.1.100",
				Reason:  "Test IP ban",
				Expires: "2026-03-24T01:45:59Z",
				Source:  "Test",
			},
			Player: Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
		},
	}

	updatedBans, err := client.AddIPBanlist(ctx, bansToAdd...)
	if err != nil {
		t.Errorf("AddIPBanlist failed: %v", err)
	}

	if updatedBans == nil {
		t.Error("Expected non-nil updated IP banlist")
	}
}

// TestIPBanlistRemove tests removing IP bans
func TestIPBanlistRemove(t *testing.T) {
	client, ctx := createTestClient(t)

	ipsToRemove := []string{"192.168.1.100"}

	updatedBans, err := client.RemoveIPBanlist(ctx, ipsToRemove...)
	if err != nil {
		t.Errorf("RemoveIPBanlist failed: %v", err)
	}

	if updatedBans == nil {
		t.Error("Expected non-nil updated IP banlist")
	}
}

// TestIPBanlistClear tests clearing the IP ban list
func TestIPBanlistClear(t *testing.T) {
	client, ctx := createTestClient(t)

	clearedBans, err := client.ClearIPBanlist(ctx)
	if err != nil {
		t.Errorf("ClearIPBanlist failed: %v", err)
	}

	if clearedBans == nil {
		t.Error("Expected non-nil cleared IP banlist")
	}
}
