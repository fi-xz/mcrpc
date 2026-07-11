package mcrpc

import (
	"testing"
)

// TestServerSettingsAutosave tests autosave settings
func TestServerSettingsAutosave(t *testing.T) {
	client, ctx := createTestClient(t)

	// Get current value
	enabled, err := client.GetAutosaveEnabled(ctx)
	if err != nil {
		t.Errorf("GetAutosaveEnabled failed: %v", err)
	}

	// Toggle value
	newEnabled, err := client.SetAutosaveEnabled(ctx, !enabled)
	if err != nil {
		t.Errorf("SetAutosaveEnabled failed: %v", err)
	}

	if newEnabled == enabled {
		t.Error("Expected autosave value to change")
	}

	// Restore original
	_, err = client.SetAutosaveEnabled(ctx, enabled)
	if err != nil {
		t.Errorf("Failed to restore autosave: %v", err)
	}
}

// TestServerSettingsDifficulty tests difficulty settings
func TestServerSettingsDifficulty(t *testing.T) {
	client, ctx := createTestClient(t)

	difficulty, err := client.GetDifficulty(ctx)
	if err != nil {
		t.Errorf("GetDifficulty failed: %v", err)
	}

	difficulties := []string{"peaceful", "easy", "normal", "hard"}
	found := false
	for _, d := range difficulties {
		if d == difficulty {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Unexpected difficulty value: %s", difficulty)
	}
}

// TestServerSettingsSetDifficulty tests setting the difficulty
func TestServerSettingsSetDifficulty(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetDifficulty(ctx)
	if err != nil {
		t.Fatalf("GetDifficulty failed: %v", err)
	}

	// Pick an alternative difficulty
	alternatives := map[string]string{"peaceful": "easy", "easy": "normal", "normal": "hard", "hard": "peaceful"}
	target := alternatives[original]

	result, err := client.SetDifficulty(ctx, target)
	if err != nil {
		t.Errorf("SetDifficulty failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected difficulty %q, got %q", target, result)
	}

	_, err = client.SetDifficulty(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore difficulty: %v", err)
	}
}

// TestServerSettingsMaxPlayers tests max players settings
func TestServerSettingsMaxPlayers(t *testing.T) {
	client, ctx := createTestClient(t)

	max, err := client.GetMaxPlayers(ctx)
	if err != nil {
		t.Errorf("GetMaxPlayers failed: %v", err)
	}

	if max < 1 {
		t.Errorf("Expected positive max players, got %d", max)
	}
}

// TestServerSettingsSetMaxPlayers tests setting max players
func TestServerSettingsSetMaxPlayers(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetMaxPlayers(ctx)
	if err != nil {
		t.Fatalf("GetMaxPlayers failed: %v", err)
	}

	target := original + 1
	result, err := client.SetMaxPlayers(ctx, target)
	if err != nil {
		t.Errorf("SetMaxPlayers failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected max players %d, got %d", target, result)
	}

	_, err = client.SetMaxPlayers(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore max players: %v", err)
	}
}

// TestServerSettingsMOTD tests MOTD settings
func TestServerSettingsMOTD(t *testing.T) {
	client, ctx := createTestClient(t)

	originalMOTD, err := client.GetMOTD(ctx)
	if err != nil {
		t.Errorf("GetMOTD failed: %v", err)
	}

	testMOTD := "Test MOTD from mcrpc"
	newMOTD, err := client.SetMOTD(ctx, testMOTD)
	if err != nil {
		t.Errorf("SetMOTD failed: %v", err)
	}

	if newMOTD != testMOTD {
		t.Errorf("Expected MOTD to be %q, got %q", testMOTD, newMOTD)
	}

	// Restore
	_, err = client.SetMOTD(ctx, originalMOTD)
	if err != nil {
		t.Errorf("Failed to restore MOTD: %v", err)
	}
}

// TestServerSettingsGameMode tests game mode settings
func TestServerSettingsGameMode(t *testing.T) {
	client, ctx := createTestClient(t)

	mode, err := client.GetGameMode(ctx)
	if err != nil {
		t.Errorf("GetGameMode failed: %v", err)
	}

	modes := []string{"creative", "survival", "spectator", "adventure"}
	found := false
	for _, m := range modes {
		if m == mode {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Unexpected game mode: %s", mode)
	}
}

// TestServerSettingsSetGameMode tests setting the game mode
func TestServerSettingsSetGameMode(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetGameMode(ctx)
	if err != nil {
		t.Fatalf("GetGameMode failed: %v", err)
	}

	alternatives := map[string]string{"survival": "creative", "creative": "survival", "spectator": "survival", "adventure": "survival"}
	target := alternatives[original]
	if target == "" {
		target = "survival"
	}

	result, err := client.SetGameMode(ctx, target)
	if err != nil {
		t.Errorf("SetGameMode failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected game mode %q, got %q", target, result)
	}

	_, err = client.SetGameMode(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore game mode: %v", err)
	}
}

// TestServerSettingsViewDistance tests view distance settings
func TestServerSettingsViewDistance(t *testing.T) {
	client, ctx := createTestClient(t)

	distance, err := client.GetViewDistance(ctx)
	if err != nil {
		t.Errorf("GetViewDistance failed: %v", err)
	}

	if distance < 1 {
		t.Errorf("Expected positive view distance, got %d", distance)
	}
}

// TestServerSettingsSetViewDistance tests setting view distance
func TestServerSettingsSetViewDistance(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetViewDistance(ctx)
	if err != nil {
		t.Fatalf("GetViewDistance failed: %v", err)
	}

	target := original + 1
	result, err := client.SetViewDistance(ctx, target)
	if err != nil {
		t.Errorf("SetViewDistance failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected view distance %d, got %d", target, result)
	}

	_, err = client.SetViewDistance(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore view distance: %v", err)
	}
}

// TestServerSettingsSimulationDistance tests simulation distance settings
func TestServerSettingsSimulationDistance(t *testing.T) {
	client, ctx := createTestClient(t)

	distance, err := client.GetSimulationDistance(ctx)
	if err != nil {
		t.Errorf("GetSimulationDistance failed: %v", err)
	}

	if distance < 1 {
		t.Errorf("Expected positive simulation distance, got %d", distance)
	}
}

// TestServerSettingsSetSimulationDistance tests setting simulation distance
func TestServerSettingsSetSimulationDistance(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetSimulationDistance(ctx)
	if err != nil {
		t.Fatalf("GetSimulationDistance failed: %v", err)
	}

	target := original + 1
	result, err := client.SetSimulationDistance(ctx, target)
	if err != nil {
		t.Errorf("SetSimulationDistance failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected simulation distance %d, got %d", target, result)
	}

	_, err = client.SetSimulationDistance(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore simulation distance: %v", err)
	}
}

// TestServerSettingsAllowlist tests allowlist settings
func TestServerSettingsAllowlist(t *testing.T) {
	client, ctx := createTestClient(t)

	enforced, err := client.GetEnforceAllowlist(ctx)
	if err != nil {
		t.Errorf("GetEnforceAllowlist failed: %v", err)
	}

	used, err := client.GetUseAllowlist(ctx)
	if err != nil {
		t.Errorf("GetUseAllowlist failed: %v", err)
	}

	t.Logf("Enforce allowlist: %v, Use allowlist: %v", enforced, used)
}

// TestServerSettingsSetEnforceAllowlist tests toggling enforce allowlist
func TestServerSettingsSetEnforceAllowlist(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetEnforceAllowlist(ctx)
	if err != nil {
		t.Fatalf("GetEnforceAllowlist failed: %v", err)
	}

	result, err := client.SetEnforceAllowlist(ctx, !original)
	if err != nil {
		t.Errorf("SetEnforceAllowlist failed: %v", err)
	}
	if result == original {
		t.Error("Expected enforce allowlist value to change")
	}

	_, err = client.SetEnforceAllowlist(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore enforce allowlist: %v", err)
	}
}

// TestServerSettingsSetUseAllowlist tests toggling use allowlist
func TestServerSettingsSetUseAllowlist(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetUseAllowlist(ctx)
	if err != nil {
		t.Fatalf("GetUseAllowlist failed: %v", err)
	}

	result, err := client.SetUseAllowlist(ctx, !original)
	if err != nil {
		t.Errorf("SetUseAllowlist failed: %v", err)
	}
	if result == original {
		t.Error("Expected use allowlist value to change")
	}

	_, err = client.SetUseAllowlist(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore use allowlist: %v", err)
	}
}

// TestServerSettingsAllowFlight tests allow flight setting
func TestServerSettingsAllowFlight(t *testing.T) {
	client, ctx := createTestClient(t)

	allowed, err := client.GetAllowFlight(ctx)
	if err != nil {
		t.Errorf("GetAllowFlight failed: %v", err)
	}

	t.Logf("Allow flight: %v", allowed)
}

// TestServerSettingsSetAllowFlight tests toggling allow flight
func TestServerSettingsSetAllowFlight(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetAllowFlight(ctx)
	if err != nil {
		t.Fatalf("GetAllowFlight failed: %v", err)
	}

	result, err := client.SetAllowFlight(ctx, !original)
	if err != nil {
		t.Errorf("SetAllowFlight failed: %v", err)
	}
	if result == original {
		t.Error("Expected allow flight value to change")
	}

	_, err = client.SetAllowFlight(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore allow flight: %v", err)
	}
}

// TestServerSettingsPauseWhenEmpty tests pause when empty setting
func TestServerSettingsPauseWhenEmpty(t *testing.T) {
	client, ctx := createTestClient(t)

	seconds, err := client.GetPauseWhenEmptySeconds(ctx)
	if err != nil {
		t.Errorf("GetPauseWhenEmptySeconds failed: %v", err)
	}

	t.Logf("Pause when empty seconds: %d", seconds)
}

// TestServerSettingsSetPauseWhenEmptySeconds tests setting pause when empty timeout
func TestServerSettingsSetPauseWhenEmptySeconds(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetPauseWhenEmptySeconds(ctx)
	if err != nil {
		t.Fatalf("GetPauseWhenEmptySeconds failed: %v", err)
	}

	target := original + 10
	result, err := client.SetPauseWhenEmptySeconds(ctx, target)
	if err != nil {
		t.Errorf("SetPauseWhenEmptySeconds failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected pause seconds %d, got %d", target, result)
	}

	_, err = client.SetPauseWhenEmptySeconds(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore pause when empty seconds: %v", err)
	}
}

// TestServerSettingsPlayerIdleTimeout tests player idle timeout setting
func TestServerSettingsPlayerIdleTimeout(t *testing.T) {
	client, ctx := createTestClient(t)

	seconds, err := client.GetPlayerIdleTimeout(ctx)
	if err != nil {
		t.Errorf("GetPlayerIdleTimeout failed: %v", err)
	}

	t.Logf("Player idle timeout: %d", seconds)
}

// TestServerSettingsSetPlayerIdleTimeout tests setting player idle timeout
func TestServerSettingsSetPlayerIdleTimeout(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetPlayerIdleTimeout(ctx)
	if err != nil {
		t.Fatalf("GetPlayerIdleTimeout failed: %v", err)
	}

	target := original + 5
	result, err := client.SetPlayerIdleTimeout(ctx, target)
	if err != nil {
		t.Errorf("SetPlayerIdleTimeout failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected idle timeout %d, got %d", target, result)
	}

	_, err = client.SetPlayerIdleTimeout(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore player idle timeout: %v", err)
	}
}

// TestServerSettingsForceGameMode tests force game mode setting
func TestServerSettingsForceGameMode(t *testing.T) {
	client, ctx := createTestClient(t)

	forced, err := client.GetForceGameMode(ctx)
	if err != nil {
		t.Errorf("GetForceGameMode failed: %v", err)
	}

	t.Logf("Force game mode: %v", forced)
}

// TestServerSettingsSetForceGameMode tests toggling force game mode
func TestServerSettingsSetForceGameMode(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetForceGameMode(ctx)
	if err != nil {
		t.Fatalf("GetForceGameMode failed: %v", err)
	}

	result, err := client.SetForceGameMode(ctx, !original)
	if err != nil {
		t.Errorf("SetForceGameMode failed: %v", err)
	}
	if result == original {
		t.Error("Expected force game mode value to change")
	}

	_, err = client.SetForceGameMode(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore force game mode: %v", err)
	}
}

// TestServerSettingsSpawnProtection tests spawn protection setting
func TestServerSettingsSpawnProtection(t *testing.T) {
	client, ctx := createTestClient(t)

	radius, err := client.GetSpawnProtectionRadius(ctx)
	if err != nil {
		t.Errorf("GetSpawnProtectionRadius failed: %v", err)
	}

	t.Logf("Spawn protection radius: %d", radius)
}

// TestServerSettingsSetSpawnProtectionRadius tests setting spawn protection radius
func TestServerSettingsSetSpawnProtectionRadius(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetSpawnProtectionRadius(ctx)
	if err != nil {
		t.Fatalf("GetSpawnProtectionRadius failed: %v", err)
	}

	target := original + 1
	result, err := client.SetSpawnProtectionRadius(ctx, target)
	if err != nil {
		t.Errorf("SetSpawnProtectionRadius failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected spawn protection radius %d, got %d", target, result)
	}

	_, err = client.SetSpawnProtectionRadius(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore spawn protection radius: %v", err)
	}
}

// TestServerSettingsAcceptTransfers tests accept transfers setting
func TestServerSettingsAcceptTransfers(t *testing.T) {
	client, ctx := createTestClient(t)

	accepted, err := client.GetAcceptTransfers(ctx)
	if err != nil {
		t.Errorf("GetAcceptTransfers failed: %v", err)
	}

	t.Logf("Accept transfers: %v", accepted)
}

// TestServerSettingsSetAcceptTransfers tests toggling accept transfers
func TestServerSettingsSetAcceptTransfers(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetAcceptTransfers(ctx)
	if err != nil {
		t.Fatalf("GetAcceptTransfers failed: %v", err)
	}

	result, err := client.SetAcceptTransfers(ctx, !original)
	if err != nil {
		t.Errorf("SetAcceptTransfers failed: %v", err)
	}
	if result == original {
		t.Error("Expected accept transfers value to change")
	}

	_, err = client.SetAcceptTransfers(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore accept transfers: %v", err)
	}
}

// TestServerSettingsStatusHeartbeat tests status heartbeat setting
func TestServerSettingsStatusHeartbeat(t *testing.T) {
	client, ctx := createTestClient(t)

	seconds, err := client.GetStatusHeartbeatInterval(ctx)
	if err != nil {
		t.Errorf("GetStatusHeartbeatInterval failed: %v", err)
	}

	t.Logf("Status heartbeat interval: %d", seconds)
}

// TestServerSettingsSetStatusHeartbeatInterval tests setting heartbeat interval
func TestServerSettingsSetStatusHeartbeatInterval(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetStatusHeartbeatInterval(ctx)
	if err != nil {
		t.Fatalf("GetStatusHeartbeatInterval failed: %v", err)
	}

	target := original + 1
	result, err := client.SetStatusHeartbeatInterval(ctx, target)
	if err != nil {
		t.Errorf("SetStatusHeartbeatInterval failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected heartbeat interval %d, got %d", target, result)
	}

	_, err = client.SetStatusHeartbeatInterval(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore heartbeat interval: %v", err)
	}
}

// TestServerSettingsOperatorPermissionLevel tests operator permission level setting
func TestServerSettingsOperatorPermissionLevel(t *testing.T) {
	client, ctx := createTestClient(t)

	level, err := client.GetOperatorPermissionLevel(ctx)
	if err != nil {
		t.Errorf("GetOperatorPermissionLevel failed: %v", err)
	}

	t.Logf("Operator permission level: %d", level)
}

// TestServerSettingsSetOperatorPermissionLevel tests setting operator permission level
func TestServerSettingsSetOperatorPermissionLevel(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetOperatorPermissionLevel(ctx)
	if err != nil {
		t.Fatalf("GetOperatorPermissionLevel failed: %v", err)
	}

	// Cycle between 3 and 4 to avoid going out of valid range
	target := 3
	if original == 3 {
		target = 4
	}

	result, err := client.SetOperatorPermissionLevel(ctx, target)
	if err != nil {
		t.Errorf("SetOperatorPermissionLevel failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected permission level %d, got %d", target, result)
	}

	_, err = client.SetOperatorPermissionLevel(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore operator permission level: %v", err)
	}
}

// TestServerSettingsHideOnlinePlayers tests hide online players setting
func TestServerSettingsHideOnlinePlayers(t *testing.T) {
	client, ctx := createTestClient(t)

	hidden, err := client.GetHideOnlinePlayers(ctx)
	if err != nil {
		t.Errorf("GetHideOnlinePlayers failed: %v", err)
	}

	t.Logf("Hide online players: %v", hidden)
}

// TestServerSettingsSetHideOnlinePlayers tests toggling hide online players
func TestServerSettingsSetHideOnlinePlayers(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetHideOnlinePlayers(ctx)
	if err != nil {
		t.Fatalf("GetHideOnlinePlayers failed: %v", err)
	}

	result, err := client.SetHideOnlinePlayers(ctx, !original)
	if err != nil {
		t.Errorf("SetHideOnlinePlayers failed: %v", err)
	}
	if result == original {
		t.Error("Expected hide online players value to change")
	}

	_, err = client.SetHideOnlinePlayers(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore hide online players: %v", err)
	}
}

// TestServerSettingsStatusReplies tests status replies setting
func TestServerSettingsStatusReplies(t *testing.T) {
	client, ctx := createTestClient(t)

	enabled, err := client.GetStatusReplies(ctx)
	if err != nil {
		t.Errorf("GetStatusReplies failed: %v", err)
	}

	t.Logf("Status replies: %v", enabled)
}

// TestServerSettingsSetStatusReplies tests toggling status replies
func TestServerSettingsSetStatusReplies(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetStatusReplies(ctx)
	if err != nil {
		t.Fatalf("GetStatusReplies failed: %v", err)
	}

	result, err := client.SetStatusReplies(ctx, !original)
	if err != nil {
		t.Errorf("SetStatusReplies failed: %v", err)
	}
	if result == original {
		t.Error("Expected status replies value to change")
	}

	_, err = client.SetStatusReplies(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore status replies: %v", err)
	}
}

// TestServerSettingsEntityBroadcastRange tests entity broadcast range setting
func TestServerSettingsEntityBroadcastRange(t *testing.T) {
	client, ctx := createTestClient(t)

	percentage, err := client.GetEntityBroadcastRange(ctx)
	if err != nil {
		t.Errorf("GetEntityBroadcastRange failed: %v", err)
	}

	t.Logf("Entity broadcast range: %d%%", percentage)
}

// TestServerSettingsSetEntityBroadcastRange tests setting entity broadcast range
func TestServerSettingsSetEntityBroadcastRange(t *testing.T) {
	client, ctx := createTestClient(t)

	original, err := client.GetEntityBroadcastRange(ctx)
	if err != nil {
		t.Fatalf("GetEntityBroadcastRange failed: %v", err)
	}

	// Stay within valid 0-500 range
	target := original + 1
	if target > 500 {
		target = original - 1
	}

	result, err := client.SetEntityBroadcastRange(ctx, target)
	if err != nil {
		t.Errorf("SetEntityBroadcastRange failed: %v", err)
	}
	if result != target {
		t.Errorf("Expected entity broadcast range %d, got %d", target, result)
	}

	_, err = client.SetEntityBroadcastRange(ctx, original)
	if err != nil {
		t.Errorf("Failed to restore entity broadcast range: %v", err)
	}
}
