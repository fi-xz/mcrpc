// Package protocol contains request parameter types for Minecraft Server Management Protocol methods.
package protocol

import (
	"github.com/fi-xz/mcrpc/internal/types"
)

// SetAllowlistParams contains parameters for setting the allowlist.
type SetAllowlistParams struct {
	Allowlist []types.Player `json:"players"`
}

// AddAllowlistParams contains parameters for adding players to the allowlist.
type AddAllowlistParams struct {
	AllowAdd []types.Player `json:"add"`
}

// RemoveAllowlistParams contains parameters for removing players from the allowlist.
type RemoveAllowlistParams struct {
	AllowRemove []types.Player `json:"remove"`
}

// SetBanlistParams contains parameters for setting the ban list.
type SetBanlistParams struct {
	Banlist []types.UserBan `json:"bans"`
}

// AddBanlistParams contains parameters for adding players to the ban list.
type AddBanlistParams struct {
	BanAdd []types.UserBan `json:"add"`
}

// RemoveBanlistParams contains parameters for removing players from the ban list.
type RemoveBanlistParams struct {
	BanRemove []types.Player `json:"remove"`
}

// SetIPBanlistParams contains parameters for setting the IP ban list.
type SetIPBanlistParams struct {
	IPBanlist []types.IPBan `json:"banlist"`
}

// AddIPBanlistParams contains parameters for adding IPs to the ban list.
type AddIPBanlistParams struct {
	IPBanAdd []types.IncomingIPBan `json:"add"`
}

// RemoveIPBanlistParams contains parameters for removing IPs from the ban list.
type RemoveIPBanlistParams struct {
	IPBanRemove []string `json:"ip"`
}

// KickPlayerParams contains parameters for kicking players from the server.
type KickPlayerParams struct {
	KickPlayers []types.KickPlayer `json:"kick"`
}

// ServerSaveParams contains parameters for saving the server state.
type ServerSaveParams struct {
	Flush bool `json:"flush"`
}

// SetOperatorParams contains parameters for setting the operator list.
type SetOperatorParams struct {
	Operators []types.Operator `json:"operators"`
}

// AddOperatorParams contains parameters for adding operators.
type AddOperatorParams struct {
	OperatorAdd []types.Operator `json:"add"`
}

// RemoveOperatorParams contains parameters for removing operators.
type RemoveOperatorParams struct {
	OperatorRemove []types.Player `json:"remove"`
}

// SystemMessageParams contains parameters for sending system messages to players.
type SystemMessageParams struct {
	Message types.SystemMessage `json:"message"`
}

// SetAutosaveParams contains parameters for enabling/disabling autosave.
type SetAutosaveParams struct {
	Enable bool `json:"enable"`
}

// SetDifficultyParams contains parameters for setting the difficulty level.
type SetDifficultyParams struct {
	Difficulty string `json:"difficulty"`
}

// SetEnforceAllowlistParams contains parameters for enabling/disabling allowlist enforcement.
type SetEnforceAllowlistParams struct {
	Enforce bool `json:"enforce"`
}

// SetUseAllowlistParams contains parameters for enabling/disabling the allowlist.
type SetUseAllowlistParams struct {
	Use bool `json:"use"`
}

// SetMaxPlayersParams contains parameters for setting the maximum number of players.
type SetMaxPlayersParams struct {
	Max int `json:"max"`
}

// SetPauseWhenEmptySecondsParams contains parameters for setting the pause when empty timeout.
type SetPauseWhenEmptySecondsParams struct {
	Seconds int `json:"seconds"`
}

// SetPlayerIdleTimeoutParams contains parameters for setting the player idle timeout.
type SetPlayerIdleTimeoutParams struct {
	Seconds int `json:"seconds"`
}

// SetAllowFlightParams contains parameters for allowing flight in Survival mode.
type SetAllowFlightParams struct {
	Allowed bool `json:"allowed"`
}

// SetMOTDParams contains parameters for setting the server message of the day.
type SetMOTDParams struct {
	Message string `json:"message"`
}

// SetSpawnProtectionRadiusParams contains parameters for setting the spawn protection radius.
type SetSpawnProtectionRadiusParams struct {
	Radius int `json:"radius"`
}

// SetForceGameModeParams contains parameters for forcing the default game mode.
type SetForceGameModeParams struct {
	Force bool `json:"force"`
}

// SetGameModeParams contains parameters for setting the default game mode.
type SetGameModeParams struct {
	Mode string `json:"mode"`
}

// SetViewDistanceParams contains parameters for setting the view distance.
type SetViewDistanceParams struct {
	Distance int `json:"distance"`
}

// SetSimulationDistanceParams contains parameters for setting the simulation distance.
type SetSimulationDistanceParams struct {
	Distance int `json:"distance"`
}

// SetAcceptTransfersParams contains parameters for enabling/disabling player transfers.
type SetAcceptTransfersParams struct {
	Accept bool `json:"accept"`
}

// SetStatusHeartbeatIntervalParams contains parameters for setting the status heartbeat interval.
type SetStatusHeartbeatIntervalParams struct {
	Seconds int `json:"seconds"`
}

// SetOperatorUserPermissionLevelParams contains parameters for setting the operator permission level.
type SetOperatorUserPermissionLevelParams struct {
	Level int `json:"level"`
}

// SetHideOnlinePlayersParams contains parameters for hiding/showing online players in status queries.
type SetHideOnlinePlayersParams struct {
	Hide bool `json:"hide"`
}

// SetStatusRepliesParams contains parameters for enabling/disabling status replies.
type SetStatusRepliesParams struct {
	Enable bool `json:"enable"`
}

// SetEntityBroadcastRangeParams contains parameters for setting the entity broadcast range.
type SetEntityBroadcastRangeParams struct {
	PercentagePoints int `json:"percentage_points"`
}

// UpdateGameRulesParams contains parameters for updating a game rule value.
type UpdateGameRulesParams struct {
	GameRules types.UntypedGameRule `json:"gamerule"`
}
