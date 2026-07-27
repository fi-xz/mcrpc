package mcrpc

import (
	"testing"
)

// TestGetOperators tests getting operators
func TestGetOperators(t *testing.T) {
	client, ctx := createTestClient(t)

	operators, err := client.GetOperators(ctx)
	if err != nil {
		t.Errorf("GetOperators failed: %v", err)
	}

	if operators == nil {
		t.Error("Expected non-nil operators list, got nil")
	}
}

// TestSetOperators tests setting operators
func TestSetOperators(t *testing.T) {
	client, ctx := createTestClient(t)

	originalOps, err := client.GetOperators(ctx)
	if err != nil {
		t.Fatalf("Failed to get current operators: %v", err)
	}

	// Set empty operators list
	updatedOps, err := client.SetOperators(ctx, []Operator{})
	if err != nil {
		t.Errorf("SetOperators failed: %v", err)
	}

	if updatedOps == nil {
		t.Error("Expected non-nil updated operators list")
	}

	// Restore original
	_, err = client.SetOperators(ctx, originalOps)
	if err != nil {
		t.Errorf("Failed to restore original operators: %v", err)
	}
}

// TestAddOperators tests adding operators
func TestAddOperators(t *testing.T) {
	client, ctx := createTestClient(t)

	opsToAdd := []Operator{
		{
			Player:              Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
			PermissionLevel:     4,
			BypassesPlayerLimit: false,
		},
	}

	updatedOps, err := client.AddOperators(ctx, opsToAdd...)
	if err != nil {
		t.Errorf("AddOperators failed: %v", err)
	}

	if updatedOps == nil {
		t.Error("Expected non-nil updated operators list")
	}
}

// TestRemoveOperators tests removing operators
func TestRemoveOperators(t *testing.T) {
	client, ctx := createTestClient(t)

	opsToRemove := []Player{
		{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"},
	}

	updatedOps, err := client.RemoveOperators(ctx, opsToRemove...)
	if err != nil {
		t.Errorf("RemoveOperators failed: %v", err)
	}

	if updatedOps == nil {
		t.Error("Expected non-nil updated operators list")
	}
}

// TestClearOperators tests clearing operators
func TestClearOperators(t *testing.T) {
	client, ctx := createTestClient(t)

	clearedOps, err := client.ClearOperators(ctx)
	if err != nil {
		t.Errorf("ClearOperators failed: %v", err)
	}

	if clearedOps == nil {
		t.Error("Expected non-nil cleared operators list")
	}
}
