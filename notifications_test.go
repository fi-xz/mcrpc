package mcrpc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/fi-xz/mcrpc/internal/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

// makeNotif builds a jsonrpc2 notification request with optional JSON params.
func makeNotif(method string, params any) *jsonrpc2.Request {
	req := &jsonrpc2.Request{Method: method, Notif: true}
	if params != nil {
		raw, _ := json.Marshal(params)
		rawMsg := json.RawMessage(raw)
		req.Params = &rawMsg
	}
	return req
}

// TestHandleIncomingDispatch tests that handleIncoming routes each notification
// method to the correct handler without requiring a real server connection.
func TestHandleIncomingDispatch(t *testing.T) {
	ctx := context.Background()

	t.Run("NilHandlerNoParam_NoPanic", func(t *testing.T) {
		client := &Client{}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerStarted, nil))
	})

	t.Run("NonNotificationIgnored", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnServerStarted = func() { called = true }
		req := &jsonrpc2.Request{Method: protocol.NotificationServerStarted, Notif: false}
		client.handleIncoming().Handle(ctx, nil, req)
		if called {
			t.Error("handler called for non-notification request")
		}
	})

	t.Run("UnknownMethodNoPanic", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnNotification = func(_ string, _ json.RawMessage) { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif("minecraft:unknown/method", nil))
		if !called {
			t.Error("OnNotification not called for unknown method")
		}
	})

	t.Run("GenericCalledBeforeSpecific", func(t *testing.T) {
		client := &Client{}
		var order []string
		client.handler.OnNotification = func(_ string, _ json.RawMessage) { order = append(order, "generic") }
		client.handler.OnServerStarted = func() { order = append(order, "specific") }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerStarted, nil))
		if len(order) != 2 || order[0] != "generic" || order[1] != "specific" {
			t.Errorf("unexpected call order: %v", order)
		}
	})

	t.Run("InvalidParamsDoesNotCallTypedHandler", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnPlayerJoined = func(_ Player) { called = true }
		rawMsg := json.RawMessage(`{invalid}`)
		req := &jsonrpc2.Request{Method: protocol.NotificationPlayerJoined, Params: &rawMsg, Notif: true}
		client.handleIncoming().Handle(ctx, nil, req)
		if called {
			t.Error("typed handler called despite invalid JSON params")
		}
	})

	t.Run("ServerStarted", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnServerStarted = func() { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerStarted, nil))
		if !called {
			t.Error("OnServerStarted not dispatched")
		}
	})

	t.Run("ServerStopping", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnServerStopping = func() { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerStopping, nil))
		if !called {
			t.Error("OnServerStopping not dispatched")
		}
	})

	t.Run("ServerSaving", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnServerSaving = func() { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerSaving, nil))
		if !called {
			t.Error("OnServerSaving not dispatched")
		}
	})

	t.Run("ServerSaved", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnServerSaved = func() { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerSaved, nil))
		if !called {
			t.Error("OnServerSaved not dispatched")
		}
	})

	t.Run("ServerStatus", func(t *testing.T) {
		client := &Client{}
		var got ServerState
		client.handler.OnServerStatus = func(s ServerState) { got = s }
		state := ServerState{Started: true, Version: Version{Name: "1.21.4", Protocol: 769}}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerStatus, []any{state}))
		if !got.Started || got.Version.Name != "1.21.4" {
			t.Errorf("OnServerStatus got unexpected value: %+v", got)
		}
	})

	t.Run("ServerActivity", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnServerActivity = func() { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationServerActivity, nil))
		if !called {
			t.Error("OnServerActivity not dispatched")
		}
	})

	t.Run("PlayerJoined", func(t *testing.T) {
		client := &Client{}
		var got Player
		client.handler.OnPlayerJoined = func(p Player) { got = p }
		player := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationPlayerJoined, []any{player}))
		if got.Name != player.Name || got.UUID != player.UUID {
			t.Errorf("OnPlayerJoined got %+v, want %+v", got, player)
		}
	})

	t.Run("PlayerLeft", func(t *testing.T) {
		client := &Client{}
		var got Player
		client.handler.OnPlayerLeft = func(p Player) { got = p }
		player := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationPlayerLeft, []any{player}))
		if got.Name != player.Name {
			t.Errorf("OnPlayerLeft got %+v, want %+v", got, player)
		}
	})

	t.Run("OperatorAdded", func(t *testing.T) {
		client := &Client{}
		var got Operator
		client.handler.OnOperatorAdded = func(op Operator) { got = op }
		op := Operator{Player: Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}, PermissionLevel: 4}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationOperatorAdded, []any{op}))
		if got.Player.Name != op.Player.Name || got.PermissionLevel != op.PermissionLevel {
			t.Errorf("OnOperatorAdded got %+v, want %+v", got, op)
		}
	})

	t.Run("OperatorRemoved", func(t *testing.T) {
		client := &Client{}
		var got Operator
		client.handler.OnOperatorRemoved = func(op Operator) { got = op }
		op := Operator{Player: Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}, PermissionLevel: 4}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationOperatorRemoved, []any{op}))
		if got.Player.Name != op.Player.Name {
			t.Errorf("OnOperatorRemoved got %+v, want %+v", got, op)
		}
	})

	t.Run("AllowlistAdded", func(t *testing.T) {
		client := &Client{}
		var got Player
		client.handler.OnAllowlistAdded = func(p Player) { got = p }
		player := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationAllowlistAdded, []any{player}))
		if got.Name != player.Name {
			t.Errorf("OnAllowlistAdded got %+v, want %+v", got, player)
		}
	})

	t.Run("AllowlistRemoved", func(t *testing.T) {
		client := &Client{}
		var got Player
		client.handler.OnAllowlistRemoved = func(p Player) { got = p }
		player := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationAllowlistRemoved, []any{player}))
		if got.Name != player.Name {
			t.Errorf("OnAllowlistRemoved got %+v, want %+v", got, player)
		}
	})

	t.Run("BanAdded", func(t *testing.T) {
		client := &Client{}
		var got UserBan
		client.handler.OnBanAdded = func(b UserBan) { got = b }
		ban := UserBan{Player: Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}, Reason: "test", Source: "Test"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationBanAdded, []any{ban}))
		if got.Player.Name != ban.Player.Name || got.Reason != ban.Reason {
			t.Errorf("OnBanAdded got %+v, want %+v", got, ban)
		}
	})

	t.Run("BanRemoved", func(t *testing.T) {
		client := &Client{}
		var got Player
		client.handler.OnBanRemoved = func(p Player) { got = p }
		player := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationBanRemoved, []any{player}))
		if got.Name != player.Name {
			t.Errorf("OnBanRemoved got %+v, want %+v", got, player)
		}
	})

	t.Run("IPBanAdded", func(t *testing.T) {
		client := &Client{}
		var got IPBan
		client.handler.OnIPBanAdded = func(b IPBan) { got = b }
		ban := IPBan{IP: "192.168.1.100", Reason: "test", Source: "Test"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationIPBanAdded, []any{ban}))
		if got.IP != ban.IP || got.Reason != ban.Reason {
			t.Errorf("OnIPBanAdded got %+v, want %+v", got, ban)
		}
	})

	t.Run("IPBanRemoved", func(t *testing.T) {
		client := &Client{}
		var got string
		client.handler.OnIPBanRemoved = func(ip string) { got = ip }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationIPBanRemoved, []any{"192.168.1.100"}))
		if got != "192.168.1.100" {
			t.Errorf("OnIPBanRemoved got %q, want %q", got, "192.168.1.100")
		}
	})

	t.Run("GameruleUpdated", func(t *testing.T) {
		client := &Client{}
		var got TypedGameRule
		client.handler.OnGameruleUpdated = func(g TypedGameRule) { got = g }
		rule := TypedGameRule{UntypedGameRule: UntypedGameRule{Key: "minecraft:keep_inventory", Value: true}, Type: "boolean"}
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationGameruleUpdated, []any{rule}))
		if got.Key != rule.Key || got.Type != rule.Type {
			t.Errorf("OnGameruleUpdated got %+v, want %+v", got, rule)
		}
	})

	t.Run("WorldUpgradeStarted", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnWorldUpgradeStarted = func() { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationWorldUpgradeStarted, nil))
		if !called {
			t.Error("OnWorldUpgradeStarted not dispatched")
		}
	})

	t.Run("WorldUpgradeProgress", func(t *testing.T) {
		client := &Client{}
		var got float64
		client.handler.OnWorldUpgradeProgress = func(p float64) { got = p }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationWorldUpgradeProgress, []any{0.42}))
		if got != 0.42 {
			t.Errorf("OnWorldUpgradeProgress got %v, want 0.42", got)
		}
	})

	t.Run("WorldUpgradeFinished", func(t *testing.T) {
		client := &Client{}
		called := false
		client.handler.OnWorldUpgradeFinished = func() { called = true }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationWorldUpgradeFinished, nil))
		if !called {
			t.Error("OnWorldUpgradeFinished not dispatched")
		}
	})

	t.Run("WorldUpgradeFailed", func(t *testing.T) {
		client := &Client{}
		var got string
		client.handler.OnWorldUpgradeFailed = func(r string) { got = r }
		client.handleIncoming().Handle(ctx, nil, makeNotif(protocol.NotificationWorldUpgradeFailed, []any{"disk full"}))
		if got != "disk full" {
			t.Errorf("OnWorldUpgradeFailed got %q, want %q", got, "disk full")
		}
	})
}

// TestNotificationHandlersExist tests that notification handler fields exist and can be set
func TestNotificationHandlersExist(t *testing.T) {
	client, _ := createTestClient(t)

	// Test that all handlers can be set
	client.handler.OnNotification = func(method string, params json.RawMessage) {}
	client.handler.OnServerStarted = func() {}
	client.handler.OnServerStopping = func() {}
	client.handler.OnServerSaving = func() {}
	client.handler.OnServerSaved = func() {}
	client.handler.OnServerStatus = func(status ServerState) {}
	client.handler.OnServerActivity = func() {}
	client.handler.OnPlayerJoined = func(player Player) {}
	client.handler.OnPlayerLeft = func(player Player) {}
	client.handler.OnOperatorAdded = func(operator Operator) {}
	client.handler.OnOperatorRemoved = func(operator Operator) {}
	client.handler.OnAllowlistAdded = func(player Player) {}
	client.handler.OnAllowlistRemoved = func(player Player) {}
	client.handler.OnBanAdded = func(ban UserBan) {}
	client.handler.OnBanRemoved = func(player Player) {}
	client.handler.OnIPBanAdded = func(ban IPBan) {}
	client.handler.OnIPBanRemoved = func(ip string) {}
	client.handler.OnGameruleUpdated = func(gamerule TypedGameRule) {}
}

// TestNotificationHandlerInvocation tests that notification handlers can be invoked
func TestNotificationHandlerInvocation(t *testing.T) {
	client, _ := createTestClient(t)

	// Test ServerStarted
	t.Run("ServerStarted", func(t *testing.T) {
		called := make(chan bool, 1)
		client.handler.OnServerStarted = func() {
			called <- true
		}

		// Directly call the handler
		client.handler.OnServerStarted()

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
		client.handler.OnPlayerJoined = func(player Player) {
			received <- player
		}

		testPlayer := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}
		client.handler.OnPlayerJoined(testPlayer)

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

		client.handler.OnNotification = func(method string, params json.RawMessage) {
			received <- struct {
				method string
				params json.RawMessage
			}{method, params}
		}

		testMethod := protocol.NotificationServerStarted
		testParams := json.RawMessage("{}")
		client.handler.OnNotification(testMethod, testParams)

		select {
		case result := <-received:
			if result.method != testMethod {
				t.Errorf("Expected method %q, got %q", testMethod, result.method)
			}
		case <-time.After(time.Second):
			t.Error("OnNotification handler was not invoked")
		}
	})

	// World upgrade handlers
	t.Run("WorldUpgradeStarted", func(t *testing.T) {
		called := make(chan bool, 1)
		client.handler.OnWorldUpgradeStarted = func() { called <- true }
		client.handler.OnWorldUpgradeStarted()
		select {
		case <-called:
		case <-time.After(time.Second):
			t.Error("OnWorldUpgradeStarted handler was not invoked")
		}
	})

	t.Run("WorldUpgradeProgress", func(t *testing.T) {
		received := make(chan float64, 1)
		client.handler.OnWorldUpgradeProgress = func(p float64) { received <- p }
		client.handler.OnWorldUpgradeProgress(0.75)
		select {
		case p := <-received:
			if p != 0.75 {
				t.Errorf("Expected progress 0.75, got %v", p)
			}
		case <-time.After(time.Second):
			t.Error("OnWorldUpgradeProgress handler was not invoked")
		}
	})

	t.Run("WorldUpgradeFinished", func(t *testing.T) {
		called := make(chan bool, 1)
		client.handler.OnWorldUpgradeFinished = func() { called <- true }
		client.handler.OnWorldUpgradeFinished()
		select {
		case <-called:
		case <-time.After(time.Second):
			t.Error("OnWorldUpgradeFinished handler was not invoked")
		}
	})

	t.Run("WorldUpgradeFailed", func(t *testing.T) {
		received := make(chan string, 1)
		client.handler.OnWorldUpgradeFailed = func(r string) { received <- r }
		client.handler.OnWorldUpgradeFailed("disk full")
		select {
		case r := <-received:
			if r != "disk full" {
				t.Errorf("Expected reason %q, got %q", "disk full", r)
			}
		case <-time.After(time.Second):
			t.Error("OnWorldUpgradeFailed handler was not invoked")
		}
	})
}
