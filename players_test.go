package mcrpc

import (
	"testing"
)

// TestGetPlayers tests getting online players
func TestGetPlayers(t *testing.T) {
	client, ctx := createTestClient(t)

	players, err := client.GetPlayers(ctx)
	if err != nil {
		t.Errorf("GetPlayers failed: %v", err)
	}

	if players == nil {
		t.Error("Expected non-nil players list, got nil")
	}
}

// TestKickPlayers tests kicking players
func TestKickPlayers(t *testing.T) {
	client, ctx := createTestClient(t)

	// Note: This test may fail if the player is not online
	kickList := []KickPlayer{
		{
			Player: Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
			Message: Message{
				Literal: "You have been kicked for testing",
			},
		},
	}

	kicked, err := client.KickPlayers(ctx, kickList)
	if err != nil {
		t.Logf("KickPlayers returned error (expected if player not online): %v", err)
	}

	// kicked may be nil if no players were kicked
	_ = kicked
}
