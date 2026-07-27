// Package mcrpc provides server settings management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetAutosaveEnabled retrieves whether automatic world saving is enabled.
func (c *MCRPCClient) GetAutosaveEnabled(ctx context.Context) (bool, error) {
	var enabled bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsAutoSaveGet, nil, &enabled)
	return enabled, err
}

// SetAutosaveEnabled enables or disables automatic world saving.
func (c *MCRPCClient) SetAutosaveEnabled(ctx context.Context, enable bool) (bool, error) {
	var enabled bool
	params := protocol.SetAutosaveParams{Enable: enable}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsAutoSaveSet, params, &enabled)
	return enabled, err
}

// GetDifficulty retrieves the current difficulty level of the server.
func (c *MCRPCClient) GetDifficulty(ctx context.Context) (string, error) {
	var difficulty string
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsDifficultyGet, nil, &difficulty)
	return difficulty, err
}

// SetDifficulty sets the difficulty level of the server (peaceful, easy, normal, hard).
func (c *MCRPCClient) SetDifficulty(ctx context.Context, difficulty string) (string, error) {
	var result string
	params := protocol.SetDifficultyParams{Difficulty: difficulty}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsDifficultySet, params, &result)
	return result, err
}

// GetEnforceAllowlist retrieves whether allowlist enforcement is enabled.
func (c *MCRPCClient) GetEnforceAllowlist(ctx context.Context) (bool, error) {
	var enforced bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsEnforceAllowlistGet, nil, &enforced)
	return enforced, err
}

// SetEnforceAllowlist enables or disables allowlist enforcement.
func (c *MCRPCClient) SetEnforceAllowlist(ctx context.Context, enforce bool) (bool, error) {
	var enforced bool
	params := protocol.SetEnforceAllowlistParams{Enforce: enforce}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsEnforceAllowlistSet, params, &enforced)
	return enforced, err
}

// GetUseAllowlist retrieves whether the allowlist is enabled.
func (c *MCRPCClient) GetUseAllowlist(ctx context.Context) (bool, error) {
	var used bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsUseAllowlistGet, nil, &used)
	return used, err
}

// SetUseAllowlist enables or disables the allowlist.
func (c *MCRPCClient) SetUseAllowlist(ctx context.Context, use bool) (bool, error) {
	var used bool
	params := protocol.SetUseAllowlistParams{Use: use}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsUseAllowlistSet, params, &used)
	return used, err
}

// GetMaxPlayers retrieves the maximum number of players allowed on the server.
func (c *MCRPCClient) GetMaxPlayers(ctx context.Context) (int, error) {
	var limit int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsMaxPlayersGet, nil, &limit)
	return limit, err
}

// SetMaxPlayers sets the maximum number of players allowed on the server.
func (c *MCRPCClient) SetMaxPlayers(ctx context.Context, limit int) (int, error) {
	var result int
	params := protocol.SetMaxPlayersParams{Max: limit}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsMaxPlayersSet, params, &result)
	return result, err
}

// GetPauseWhenEmptySeconds retrieves the pause when empty timeout in seconds.
func (c *MCRPCClient) GetPauseWhenEmptySeconds(ctx context.Context) (int, error) {
	var seconds int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsPauseWhenEmptyGet, nil, &seconds)
	return seconds, err
}

// SetPauseWhenEmptySeconds sets the pause when empty timeout in seconds.
func (c *MCRPCClient) SetPauseWhenEmptySeconds(ctx context.Context, seconds int) (int, error) {
	var result int
	params := protocol.SetPauseWhenEmptySecondsParams{Seconds: seconds}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsPauseWhenEmptySet, params, &result)
	return result, err
}

// GetPlayerIdleTimeout retrieves the player idle timeout in seconds.
func (c *MCRPCClient) GetPlayerIdleTimeout(ctx context.Context) (int, error) {
	var seconds int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsPlayerIdleTimeoutGet, nil, &seconds)
	return seconds, err
}

// SetPlayerIdleTimeout sets the player idle timeout in seconds.
func (c *MCRPCClient) SetPlayerIdleTimeout(ctx context.Context, seconds int) (int, error) {
	var result int
	params := protocol.SetPlayerIdleTimeoutParams{Seconds: seconds}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsPlayerIdleTimeoutSet, params, &result)
	return result, err
}

// GetAllowFlight retrieves whether flight is allowed in Survival mode.
func (c *MCRPCClient) GetAllowFlight(ctx context.Context) (bool, error) {
	var allowed bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsAllowFlightGet, nil, &allowed)
	return allowed, err
}

// SetAllowFlight sets whether flight is allowed in Survival mode.
func (c *MCRPCClient) SetAllowFlight(ctx context.Context, allowed bool) (bool, error) {
	var result bool
	params := protocol.SetAllowFlightParams{Allowed: allowed}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsAllowFlightSet, params, &result)
	return result, err
}

// GetMOTD retrieves the server's message of the day.
func (c *MCRPCClient) GetMOTD(ctx context.Context) (string, error) {
	var message string
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsMOTDGet, nil, &message)
	return message, err
}

// SetMOTD sets the server's message of the day.
func (c *MCRPCClient) SetMOTD(ctx context.Context, message string) (string, error) {
	var result string
	params := protocol.SetMOTDParams{Message: message}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsMOTDSet, params, &result)
	return result, err
}

// GetSpawnProtectionRadius retrieves the spawn protection radius in blocks.
func (c *MCRPCClient) GetSpawnProtectionRadius(ctx context.Context) (int, error) {
	var radius int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsSpawnProtectionGet, nil, &radius)
	return radius, err
}

// SetSpawnProtectionRadius sets the spawn protection radius in blocks.
func (c *MCRPCClient) SetSpawnProtectionRadius(ctx context.Context, radius int) (int, error) {
	var result int
	params := protocol.SetSpawnProtectionRadiusParams{Radius: radius}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsSpawnProtectionSet, params, &result)
	return result, err
}

// GetForceGameMode retrieves whether players are forced to use the default game mode.
func (c *MCRPCClient) GetForceGameMode(ctx context.Context) (bool, error) {
	var forced bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsForceGameModeGet, nil, &forced)
	return forced, err
}

// SetForceGameMode sets whether players are forced to use the default game mode.
func (c *MCRPCClient) SetForceGameMode(ctx context.Context, force bool) (bool, error) {
	var forced bool
	params := protocol.SetForceGameModeParams{Force: force}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsForceGameModeSet, params, &forced)
	return forced, err
}

// GetGameMode retrieves the server's default game mode.
func (c *MCRPCClient) GetGameMode(ctx context.Context) (string, error) {
	var mode string
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsGameModeGet, nil, &mode)
	return mode, err
}

// SetGameMode sets the server's default game mode (creative, survival, spectator, adventure).
func (c *MCRPCClient) SetGameMode(ctx context.Context, mode string) (string, error) {
	var result string
	params := protocol.SetGameModeParams{Mode: mode}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsGameModeSet, params, &result)
	return result, err
}

// GetViewDistance retrieves the view distance in chunks.
func (c *MCRPCClient) GetViewDistance(ctx context.Context) (int, error) {
	var distance int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsViewDistanceGet, nil, &distance)
	return distance, err
}

// SetViewDistance sets the view distance in chunks.
func (c *MCRPCClient) SetViewDistance(ctx context.Context, distance int) (int, error) {
	var result int
	params := protocol.SetViewDistanceParams{Distance: distance}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsViewDistanceSet, params, &result)
	return result, err
}

// GetSimulationDistance retrieves the simulation distance in chunks.
func (c *MCRPCClient) GetSimulationDistance(ctx context.Context) (int, error) {
	var distance int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsSimulationDistanceGet, nil, &distance)
	return distance, err
}

// SetSimulationDistance sets the simulation distance in chunks.
func (c *MCRPCClient) SetSimulationDistance(ctx context.Context, distance int) (int, error) {
	var result int
	params := protocol.SetSimulationDistanceParams{Distance: distance}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsSimulationDistanceSet, params, &result)
	return result, err
}

// GetAcceptTransfers retrieves whether the server accepts player transfers.
func (c *MCRPCClient) GetAcceptTransfers(ctx context.Context) (bool, error) {
	var accepted bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsAcceptTransfersGet, nil, &accepted)
	return accepted, err
}

// SetAcceptTransfers sets whether the server accepts player transfers.
func (c *MCRPCClient) SetAcceptTransfers(ctx context.Context, accept bool) (bool, error) {
	var accepted bool
	params := protocol.SetAcceptTransfersParams{Accept: accept}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsAcceptTransfersSet, params, &accepted)
	return accepted, err
}

// GetStatusHeartbeatInterval retrieves the status heartbeat interval in seconds.
func (c *MCRPCClient) GetStatusHeartbeatInterval(ctx context.Context) (int, error) {
	var seconds int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsHeartbeatIntervalGet, nil, &seconds)
	return seconds, err
}

// SetStatusHeartbeatInterval sets the status heartbeat interval in seconds.
func (c *MCRPCClient) SetStatusHeartbeatInterval(ctx context.Context, seconds int) (int, error) {
	var result int
	params := protocol.SetStatusHeartbeatIntervalParams{Seconds: seconds}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsHeartbeatIntervalSet, params, &result)
	return result, err
}

// GetOperatorPermissionLevel retrieves the operator permission level.
func (c *MCRPCClient) GetOperatorPermissionLevel(ctx context.Context) (int, error) {
	var level int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsOperatorPermissionLevelGet, nil, &level)
	return level, err
}

// SetOperatorPermissionLevel sets the operator permission level.
func (c *MCRPCClient) SetOperatorPermissionLevel(ctx context.Context, level int) (int, error) {
	var result int
	params := protocol.SetOperatorUserPermissionLevelParams{Level: level}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsOperatorPermissionLevelSet, params, &result)
	return result, err
}

// GetHideOnlinePlayers retrieves whether online player info is hidden from status queries.
func (c *MCRPCClient) GetHideOnlinePlayers(ctx context.Context) (bool, error) {
	var hidden bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsHideOnlinePlayersGet, nil, &hidden)
	return hidden, err
}

// SetHideOnlinePlayers sets whether online player info is hidden from status queries.
func (c *MCRPCClient) SetHideOnlinePlayers(ctx context.Context, hide bool) (bool, error) {
	var hidden bool
	params := protocol.SetHideOnlinePlayersParams{Hide: hide}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsHideOnlinePlayersSet, params, &hidden)
	return hidden, err
}

// GetStatusReplies retrieves whether the server responds to status requests.
func (c *MCRPCClient) GetStatusReplies(ctx context.Context) (bool, error) {
	var enabled bool
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsStatusRepliesGet, nil, &enabled)
	return enabled, err
}

// SetStatusReplies sets whether the server responds to status requests.
func (c *MCRPCClient) SetStatusReplies(ctx context.Context, enable bool) (bool, error) {
	var enabled bool
	params := protocol.SetStatusRepliesParams{Enable: enable}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsStatusRepliesSet, params, &enabled)
	return enabled, err
}

// GetEntityBroadcastRange retrieves the entity broadcast range percentage.
func (c *MCRPCClient) GetEntityBroadcastRange(ctx context.Context) (int, error) {
	var percentage int
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsEntityBroadcastRangeGet, nil, &percentage)
	return percentage, err
}

// SetEntityBroadcastRange sets the entity broadcast range percentage.
func (c *MCRPCClient) SetEntityBroadcastRange(ctx context.Context, percentagePoints int) (int, error) {
	var result int
	params := protocol.SetEntityBroadcastRangeParams{PercentagePoints: percentagePoints}
	err := c.JSONRPCConn.Call(ctx, protocol.MethodServerSettingsEntityBroadcastRangeSet, params, &result)
	return result, err
}
