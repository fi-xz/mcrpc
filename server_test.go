package mcrpc

import (
	"testing"
)

// TestGetServerStatus tests getting server status
func TestGetServerStatus(t *testing.T) {
	client, ctx := createTestClient(t)

	status, err := client.GetServerStatus(ctx)
	if err != nil {
		t.Errorf("GetServerStatus failed: %v", err)
	}

	// ServerState should have valid fields
	_ = status
}

// TestSaveServer tests saving the server
func TestSaveServer(t *testing.T) {
	client, ctx := createTestClient(t)

	saving, err := client.SaveServer(ctx, true)
	if err != nil {
		t.Errorf("SaveServer failed: %v", err)
	}

	// saving should be true if save was initiated
	_ = saving
}

// TestSendSystemMessage tests sending a system message
func TestSendSystemMessage(t *testing.T) {
	client, ctx := createTestClient(t)

	message := SystemMessage{
		ReceivingPlayers: []Player{
			{
				Name: "fi_xz",
				UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec",
			},
		},
		Overlay: false,
		Message: Message{
			Literal: "Test system message from mcrpc",
		},
	}

	sent, err := client.SendSystemMessage(ctx, message)
	if err != nil {
		t.Errorf("SendSystemMessage failed: %v", err)
	}

	if !sent {
		t.Error("Expected message to be sent successfully")
	}
}

// TestStopServer tests stopping the server
// NOTE: This test is commented out as it would actually stop the server
/*
func TestStopServer(t *testing.T) {
	client, ctx := createTestClient(t)

	stopping, err := client.StopServer(ctx)
	if err != nil {
		t.Errorf("StopServer failed: %v", err)
	}

	if !stopping {
		t.Error("Expected server to be stopping")
	}
}
*/
