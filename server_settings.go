// Package mcrpc provides server settings management functionality for Minecraft servers.
package mcrpc

import (
	"context"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// GetAutosaveEnabled retrieves whether automatic world saving is enabled.
func (c *MCRPCClient) GetAutosaveEnabled(context context.Context) (bool, error) {
	var enabled bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsAutoSaveGet, nil, &enabled)
	return enabled, err
}

// SetAutosaveEnabled enables or disables automatic world saving.
func (c *MCRPCClient) SetAutosaveEnabled(context context.Context, enable bool) (bool, error) {
	var enabled bool
	params := protocol.SetAutosaveParams{Enable: enable}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsAutoSaveSet, params, &enabled)
	return enabled, err
}

// GetDifficulty retrieves the current difficulty level of the server.
func (c *MCRPCClient) GetDifficulty(context context.Context) (string, error) {
	var difficulty string
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsDifficultyGet, nil, &difficulty)
	return difficulty, err
}

// SetDifficulty sets the difficulty level of the server (peaceful, easy, normal, hard).
func (c *MCRPCClient) SetDifficulty(context context.Context, difficulty string) (string, error) {
	var result string
	params := protocol.SetDifficultyParams{Difficulty: difficulty}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsDifficultySet, params, &result)
	return result, err
}

// GetEnforceAllowlist retrieves whether allowlist enforcement is enabled.
func (c *MCRPCClient) GetEnforceAllowlist(context context.Context) (bool, error) {
	var enforced bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsEnforceAllowlistGet, nil, &enforced)
	return enforced, err
}

// SetEnforceAllowlist enables or disables allowlist enforcement.
func (c *MCRPCClient) SetEnforceAllowlist(context context.Context, enforce bool) (bool, error) {
	var enforced bool
	params := protocol.SetEnforceAllowlistParams{Enforce: enforce}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsEnforceAllowlistSet, params, &enforced)
	return enforced, err
}

// GetUseAllowlist retrieves whether the allowlist is enabled.
func (c *MCRPCClient) GetUseAllowlist(context context.Context) (bool, error) {
	var used bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsUseAllowlistGet, nil, &used)
	return used, err
}

// SetUseAllowlist enables or disables the allowlist.
func (c *MCRPCClient) SetUseAllowlist(context context.Context, use bool) (bool, error) {
	var used bool
	params := protocol.SetUseAllowlistParams{Use: use}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsUseAllowlistSet, params, &used)
	return used, err
}

// GetMaxPlayers retrieves the maximum number of players allowed on the server.
func (c *MCRPCClient) GetMaxPlayers(context context.Context) (int, error) {
	var max int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsMaxPlayersGet, nil, &max)
	return max, err
}

// SetMaxPlayers sets the maximum number of players allowed on the server.
func (c *MCRPCClient) SetMaxPlayers(context context.Context, max int) (int, error) {
	var result int
	params := protocol.SetMaxPlayersParams{Max: max}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsMaxPlayersSet, params, &result)
	return result, err
}

// GetPauseWhenEmptySeconds retrieves the pause when empty timeout in seconds.
func (c *MCRPCClient) GetPauseWhenEmptySeconds(context context.Context) (int, error) {
	var seconds int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsPauseWhenEmptyGet, nil, &seconds)
	return seconds, err
}

// SetPauseWhenEmptySeconds sets the pause when empty timeout in seconds.
func (c *MCRPCClient) SetPauseWhenEmptySeconds(context context.Context, seconds int) (int, error) {
	var result int
	params := protocol.SetPauseWhenEmptySecondsParams{Seconds: seconds}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsPauseWhenEmptySet, params, &result)
	return result, err
}

// GetPlayerIdleTimeout retrieves the player idle timeout in seconds.
func (c *MCRPCClient) GetPlayerIdleTimeout(context context.Context) (int, error) {
	var seconds int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsPlayerIdleTimeoutGet, nil, &seconds)
	return seconds, err
}

// SetPlayerIdleTimeout sets the player idle timeout in seconds.
func (c *MCRPCClient) SetPlayerIdleTimeout(context context.Context, seconds int) (int, error) {
	var result int
	params := protocol.SetPlayerIdleTimeoutParams{Seconds: seconds}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsPlayerIdleTimeoutSet, params, &result)
	return result, err
}

// GetAllowFlight retrieves whether flight is allowed in Survival mode.
func (c *MCRPCClient) GetAllowFlight(context context.Context) (bool, error) {
	var allowed bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsAllowFlightGet, nil, &allowed)
	return allowed, err
}

// SetAllowFlight sets whether flight is allowed in Survival mode.
func (c *MCRPCClient) SetAllowFlight(context context.Context, allowed bool) (bool, error) {
	var result bool
	params := protocol.SetAllowFlightParams{Allowed: allowed}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsAllowFlightSet, params, &result)
	return result, err
}

// GetMOTD retrieves the server's message of the day.
func (c *MCRPCClient) GetMOTD(context context.Context) (string, error) {
	var message string
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsMOTDGet, nil, &message)
	return message, err
}

// SetMOTD sets the server's message of the day.
func (c *MCRPCClient) SetMOTD(context context.Context, message string) (string, error) {
	var result string
	params := protocol.SetMOTDParams{Message: message}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsMOTDSet, params, &result)
	return result, err
}

// GetSpawnProtectionRadius retrieves the spawn protection radius in blocks.
func (c *MCRPCClient) GetSpawnProtectionRadius(context context.Context) (int, error) {
	var radius int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsSpawnProtectionGet, nil, &radius)
	return radius, err
}

// SetSpawnProtectionRadius sets the spawn protection radius in blocks.
func (c *MCRPCClient) SetSpawnProtectionRadius(context context.Context, radius int) (int, error) {
	var result int
	params := protocol.SetSpawnProtectionRadiusParams{Radius: radius}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsSpawnProtectionSet, params, &result)
	return result, err
}

// GetForceGameMode retrieves whether players are forced to use the default game mode.
func (c *MCRPCClient) GetForceGameMode(context context.Context) (bool, error) {
	var forced bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsForceGameModeGet, nil, &forced)
	return forced, err
}

// SetForceGameMode sets whether players are forced to use the default game mode.
func (c *MCRPCClient) SetForceGameMode(context context.Context, force bool) (bool, error) {
	var forced bool
	params := protocol.SetForceGameModeParams{Force: force}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsForceGameModeSet, params, &forced)
	return forced, err
}

// GetGameMode retrieves the server's default game mode.
func (c *MCRPCClient) GetGameMode(context context.Context) (string, error) {
	var mode string
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsGameModeGet, nil, &mode)
	return mode, err
}

// SetGameMode sets the server's default game mode (creative, survival, spectator, adventure).
func (c *MCRPCClient) SetGameMode(context context.Context, mode string) (string, error) {
	var result string
	params := protocol.SetGameModeParams{Mode: mode}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsGameModeSet, params, &result)
	return result, err
}

// GetViewDistance retrieves the view distance in chunks.
func (c *MCRPCClient) GetViewDistance(context context.Context) (int, error) {
	var distance int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsViewDistanceGet, nil, &distance)
	return distance, err
}

// SetViewDistance sets the view distance in chunks.
func (c *MCRPCClient) SetViewDistance(context context.Context, distance int) (int, error) {
	var result int
	params := protocol.SetViewDistanceParams{Distance: distance}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsViewDistanceSet, params, &result)
	return result, err
}

// GetSimulationDistance retrieves the simulation distance in chunks.
func (c *MCRPCClient) GetSimulationDistance(context context.Context) (int, error) {
	var distance int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsSimulationDistanceGet, nil, &distance)
	return distance, err
}

// SetSimulationDistance sets the simulation distance in chunks.
func (c *MCRPCClient) SetSimulationDistance(context context.Context, distance int) (int, error) {
	var result int
	params := protocol.SetSimulationDistanceParams{Distance: distance}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsSimulationDistanceSet, params, &result)
	return result, err
}

// GetAcceptTransfers retrieves whether the server accepts player transfers.
func (c *MCRPCClient) GetAcceptTransfers(context context.Context) (bool, error) {
	var accepted bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsAcceptTransfersGet, nil, &accepted)
	return accepted, err
}

// SetAcceptTransfers sets whether the server accepts player transfers.
func (c *MCRPCClient) SetAcceptTransfers(context context.Context, accept bool) (bool, error) {
	var accepted bool
	params := protocol.SetAcceptTransfersParams{Accept: accept}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsAcceptTransfersSet, params, &accepted)
	return accepted, err
}

// GetStatusHeartbeatInterval retrieves the status heartbeat interval in seconds.
func (c *MCRPCClient) GetStatusHeartbeatInterval(context context.Context) (int, error) {
	var seconds int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsHeartbeatIntervalGet, nil, &seconds)
	return seconds, err
}

// SetStatusHeartbeatInterval sets the status heartbeat interval in seconds.
func (c *MCRPCClient) SetStatusHeartbeatInterval(context context.Context, seconds int) (int, error) {
	var result int
	params := protocol.SetStatusHeartbeatIntervalParams{Seconds: seconds}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsHeartbeatIntervalSet, params, &result)
	return result, err
}

// GetOperatorPermissionLevel retrieves the operator permission level.
func (c *MCRPCClient) GetOperatorPermissionLevel(context context.Context) (int, error) {
	var level int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsOperatorPermissionLevelGet, nil, &level)
	return level, err
}

// SetOperatorPermissionLevel sets the operator permission level.
func (c *MCRPCClient) SetOperatorPermissionLevel(context context.Context, level int) (int, error) {
	var result int
	params := protocol.SetOperatorUserPermissionLevelParams{Level: level}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsOperatorPermissionLevelSet, params, &result)
	return result, err
}

// GetHideOnlinePlayers retrieves whether online player info is hidden from status queries.
func (c *MCRPCClient) GetHideOnlinePlayers(context context.Context) (bool, error) {
	var hidden bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsHideOnlinePlayersGet, nil, &hidden)
	return hidden, err
}

// SetHideOnlinePlayers sets whether online player info is hidden from status queries.
func (c *MCRPCClient) SetHideOnlinePlayers(context context.Context, hide bool) (bool, error) {
	var hidden bool
	params := protocol.SetHideOnlinePlayersParams{Hide: hide}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsHideOnlinePlayersSet, params, &hidden)
	return hidden, err
}

// GetStatusReplies retrieves whether the server responds to status requests.
func (c *MCRPCClient) GetStatusReplies(context context.Context) (bool, error) {
	var enabled bool
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsStatusRepliesGet, nil, &enabled)
	return enabled, err
}

// SetStatusReplies sets whether the server responds to status requests.
func (c *MCRPCClient) SetStatusReplies(context context.Context, enable bool) (bool, error) {
	var enabled bool
	params := protocol.SetStatusRepliesParams{Enable: enable}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsStatusRepliesSet, params, &enabled)
	return enabled, err
}

// GetEntityBroadcastRange retrieves the entity broadcast range percentage.
func (c *MCRPCClient) GetEntityBroadcastRange(context context.Context) (int, error) {
	var percentage int
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsEntityBroadcastRangeGet, nil, &percentage)
	return percentage, err
}

// SetEntityBroadcastRange sets the entity broadcast range percentage.
func (c *MCRPCClient) SetEntityBroadcastRange(context context.Context, percentagePoints int) (int, error) {
	var result int
	params := protocol.SetEntityBroadcastRangeParams{PercentagePoints: percentagePoints}
	err := c.JSONRPCConn.Call(context, protocol.MethodServerSettingsEntityBroadcastRangeSet, params, &result)
	return result, err
}
