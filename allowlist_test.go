package mcrpc

import (
	"testing"
)

// TestAllowlistGet tests getting the allowlist
func TestAllowlistGet(t *testing.T) {
	client, ctx := createTestClient(t)

	players, err := client.GetAllowlist(ctx)
	if err != nil {
		t.Errorf("GetAllowlist failed: %v", err)
	}

	// Allowlist should be a slice (may be empty)
	if players == nil {
		t.Error("Expected non-nil allowlist, got nil")
	}
}

// TestAllowlistSet tests setting the allowlist
func TestAllowlistSet(t *testing.T) {
	client, ctx := createTestClient(t)

	// First get current allowlist
	originalPlayers, err := client.GetAllowlist(ctx)
	if err != nil {
		t.Fatalf("Failed to get current allowlist: %v", err)
	}

	// Set a test allowlist
	testPlayers := []Player{
		{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
	}

	updatedPlayers, err := client.SetAllowlist(ctx, testPlayers)
	if err != nil {
		t.Errorf("SetAllowlist failed: %v", err)
	}

	if updatedPlayers == nil {
		t.Error("Expected non-nil updated allowlist")
	}

	// Restore original allowlist
	_, err = client.SetAllowlist(ctx, originalPlayers)
	if err != nil {
		t.Errorf("Failed to restore original allowlist: %v", err)
	}
}

// TestAllowlistAdd tests adding players to the allowlist
func TestAllowlistAdd(t *testing.T) {
	client, ctx := createTestClient(t)

	playersToAdd := []Player{
		{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
	}

	updatedPlayers, err := client.AddAllowlist(ctx, playersToAdd)
	if err != nil {
		t.Errorf("AddAllowlist failed: %v", err)
	}

	if updatedPlayers == nil {
		t.Error("Expected non-nil updated allowlist")
	}
}

// TestAllowlistRemove tests removing players from the allowlist
func TestAllowlistRemove(t *testing.T) {
	client, ctx := createTestClient(t)

	playersToRemove := []Player{
		{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
	}

	updatedPlayers, err := client.RemoveAllowlist(ctx, playersToRemove)
	if err != nil {
		t.Errorf("RemoveAllowlist failed: %v", err)
	}

	if updatedPlayers == nil {
		t.Error("Expected non-nil updated allowlist")
	}
}

// TestAllowlistClear tests clearing the allowlist
func TestAllowlistClear(t *testing.T) {
	client, ctx := createTestClient(t)

	// First add a player so we have something to clear
	testPlayer := []Player{
		{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
	}
	_, err := client.AddAllowlist(ctx, testPlayer)
	if err != nil {
		t.Logf("Warning: Failed to add test player: %v", err)
	}

	clearedPlayers, err := client.ClearAllowlist(ctx)
	if err != nil {
		t.Errorf("ClearAllowlist failed: %v", err)
	}

	if clearedPlayers == nil {
		t.Error("Expected non-nil cleared allowlist")
	}
}
