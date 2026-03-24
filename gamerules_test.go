package mcrpc

import (
	"testing"
)

// TestGetGameRules tests getting gamerules
func TestGetGameRules(t *testing.T) {
	client, ctx := createTestClient(t)

	gamerules, err := client.GetGameRules(ctx)
	if err != nil {
		t.Errorf("GetGameRules failed: %v", err)
	}

	if gamerules == nil {
		t.Error("Expected non-nil gamerules list, got nil")
	}
}

// TestUpdateGameRule tests updating a gamerule
func TestUpdateGameRule(t *testing.T) {
	client, ctx := createTestClient(t)

	// First get current gamerules
	gamerules, err := client.GetGameRules(ctx)
	if err != nil {
		t.Fatalf("Failed to get current gamerules: %v", err)
	}

	if len(gamerules) == 0 {
		t.Skip("No gamerules available to test update")
	}

	// Try to update the first gamerule with its current value
	firstRule := gamerules[0]

	updatedRule, err := client.UpdateGameRule(ctx, firstRule.UntypedGameRule)
	if err != nil {
		t.Errorf("UpdateGameRule failed: %v", err)
	}

	if updatedRule.Key != firstRule.Key {
		t.Errorf("Expected updated rule key %q, got %q", firstRule.Key, updatedRule.Key)
	}
}
