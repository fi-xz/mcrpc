package mcrpc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// TestNotificationHandlersExist tests that notification handler fields exist and can be set
func TestNotificationHandlersExist(t *testing.T) {
	client, _ := createTestClient(t)

	// Test that all handlers can be set
	client.OnNotification = func(method string, params json.RawMessage) {}
	client.OnServerStarted = func() {}
	client.OnServerStopping = func() {}
	client.OnServerSaving = func() {}
	client.OnServerSaved = func() {}
	client.OnServerStatus = func(status ServerState) {}
	client.OnServerActivity = func() {}
	client.OnPlayerJoined = func(player Player) {}
	client.OnPlayerLeft = func(player Player) {}
	client.OnOperatorAdded = func(operator Operator) {}
	client.OnOperatorRemoved = func(operator Operator) {}
	client.OnAllowlistAdded = func(player Player) {}
	client.OnAllowlistRemoved = func(player Player) {}
	client.OnBanAdded = func(ban UserBan) {}
	client.OnBanRemoved = func(player Player) {}
	client.OnIPBanAdded = func(ban IPBan) {}
	client.OnIPBanRemoved = func(ip string) {}
	client.OnGameruleUpdated = func(gamerule TypedGameRule) {}
}

// TestNotificationHandlerInvocation tests that notification handlers can be invoked
func TestNotificationHandlerInvocation(t *testing.T) {
	client, _ := createTestClient(t)

	// Test ServerStarted
	t.Run("ServerStarted", func(t *testing.T) {
		called := make(chan bool, 1)
		client.OnServerStarted = func() {
			called <- true
		}

		// Directly call the handler
		client.OnServerStarted()

		select {
		case <-called:
			// Success
		case <-time.After(time.Second):
			t.Error("OnServerStarted handler was not invoked")
		}
	})

	// Test PlayerJoined
	t.Run("PlayerJoined", func(t *testing.T) {
		received := make(chan Player, 1)
		client.OnPlayerJoined = func(player Player) {
			received <- player
		}

		testPlayer := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}
		client.OnPlayerJoined(testPlayer)

		select {
		case player := <-received:
			if player.Name != testPlayer.Name {
				t.Errorf("Expected player name %q, got %q", testPlayer.Name, player.Name)
			}
		case <-time.After(time.Second):
			t.Error("OnPlayerJoined handler was not invoked")
		}
	})

	// Test generic OnNotification
	t.Run("GenericNotification", func(t *testing.T) {
		received := make(chan struct {
			method string
			params json.RawMessage
		}, 1)

		client.OnNotification = func(method string, params json.RawMessage) {
			received <- struct {
				method string
				params json.RawMessage
			}{method, params}
		}

		testMethod := protocol.NotificationServerStarted
		testParams := json.RawMessage("{}")
		client.OnNotification(testMethod, testParams)

		select {
		case result := <-received:
			if result.method != testMethod {
				t.Errorf("Expected method %q, got %q", testMethod, result.method)
			}
		case <-time.After(time.Second):
			t.Error("OnNotification handler was not invoked")
		}
	})
}
