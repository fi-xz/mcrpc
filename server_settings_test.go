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

// TestServerSettingsAllowFlight tests allow flight setting
func TestServerSettingsAllowFlight(t *testing.T) {
	client, ctx := createTestClient(t)

	allowed, err := client.GetAllowFlight(ctx)
	if err != nil {
		t.Errorf("GetAllowFlight failed: %v", err)
	}

	t.Logf("Allow flight: %v", allowed)
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

// TestServerSettingsPlayerIdleTimeout tests player idle timeout setting
func TestServerSettingsPlayerIdleTimeout(t *testing.T) {
	client, ctx := createTestClient(t)

	seconds, err := client.GetPlayerIdleTimeout(ctx)
	if err != nil {
		t.Errorf("GetPlayerIdleTimeout failed: %v", err)
	}

	t.Logf("Player idle timeout: %d", seconds)
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

// TestServerSettingsSpawnProtection tests spawn protection setting
func TestServerSettingsSpawnProtection(t *testing.T) {
	client, ctx := createTestClient(t)

	radius, err := client.GetSpawnProtectionRadius(ctx)
	if err != nil {
		t.Errorf("GetSpawnProtectionRadius failed: %v", err)
	}

	t.Logf("Spawn protection radius: %d", radius)
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

// TestServerSettingsStatusHeartbeat tests status heartbeat setting
func TestServerSettingsStatusHeartbeat(t *testing.T) {
	client, ctx := createTestClient(t)

	seconds, err := client.GetStatusHeartbeatInterval(ctx)
	if err != nil {
		t.Errorf("GetStatusHeartbeatInterval failed: %v", err)
	}

	t.Logf("Status heartbeat interval: %d", seconds)
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

// TestServerSettingsHideOnlinePlayers tests hide online players setting
func TestServerSettingsHideOnlinePlayers(t *testing.T) {
	client, ctx := createTestClient(t)

	hidden, err := client.GetHideOnlinePlayers(ctx)
	if err != nil {
		t.Errorf("GetHideOnlinePlayers failed: %v", err)
	}

	t.Logf("Hide online players: %v", hidden)
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

// TestServerSettingsEntityBroadcastRange tests entity broadcast range setting
func TestServerSettingsEntityBroadcastRange(t *testing.T) {
	client, ctx := createTestClient(t)

	percentage, err := client.GetEntityBroadcastRange(ctx)
	if err != nil {
		t.Errorf("GetEntityBroadcastRange failed: %v", err)
	}

	t.Logf("Entity broadcast range: %d%%", percentage)
}
