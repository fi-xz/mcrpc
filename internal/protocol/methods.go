package protocol

const (
	// Request methods

	// Allowlist management

	// Get the current allowlist
	MethodAllowlistGet = "minecraft:allowlist"
	// Set the allowlist to a specific list of players
	MethodAllowlistSet = "minecraft:allowlist/set"
	// Add players to the allowlist
	MethodAllowlistAdd = "minecraft:allowlist/add"
	// Remove players from the allowlist
	MethodAllowlistRemove = "minecraft:allowlist/remove"
	// Clear the allowlist
	MethodAllowlistClear = "minecraft:allowlist/clear"

	// Ban management

	// Get the current ban list
	MethodBansGet = "minecraft:bans"
	// Set the ban list to a specific list of players
	MethodBansSet = "minecraft:bans/set"
	// Add players to the ban list
	MethodBansAdd = "minecraft:bans/add"
	// Remove players from the ban list
	MethodBansRemove = "minecraft:bans/remove"
	// Clear the ban list
	MethodBansClear = "minecraft:bans/clear"

	// IP Ban management

	// Get the current IP ban list
	MethodIPBansGet = "minecraft:ip_bans"
	// Set the IP ban list to a specific list of IPs
	MethodIPBansSet = "minecraft:ip_bans/set"
	// Add IPs to the IP ban list
	MethodIPBansAdd = "minecraft:ip_bans/add"
	// Remove IPs from the IP ban list
	MethodIPBansRemove = "minecraft:ip_bans/remove"
	// Clear the IP ban list
	MethodIPBansClear = "minecraft:ip_bans/clear"

	// Players Management

	// Get the list of currently online players
	MethodPlayersGet = "minecraft:players"
	// Kick a player from the server
	MethodPlayersKick = "minecraft:players/kick"

	// Operators Management

	// Get the list of operators
	MethodOperatorsGet = "minecraft:operators"
	// Set the list of operators
	MethodOperatorsSet = "minecraft:operators/set"
	// Add players as operators
	MethodOperatorsAdd = "minecraft:operators/add"
	// Remove players as operators
	MethodOperatorsRemove = "minecraft:operators/remove"
	// Clear all operators
	MethodOperatorsClear = "minecraft:operators/clear"

	// Server Management

	// Get the server status
	MethodServerStatus = "minecraft:server/status"
	// Save the server state
	MethodServerSave = "minecraft:server/save"
	// Stop the server
	MethodServerStop = "minecraft:server/stop"
	// Send a system message to players
	MethodServerSystemMessage = "minecraft:server/system_message"

	// Server Setting Management

	// Get whether automatic world saving is enabled
	MethodServerSettingsAutoSaveGet = "minecraft:serversettings/autosave"
	// Enable or disable automatic world saving
	MethodServerSettingsAutoSaveSet = "minecraft:serversettings/autosave/set"
	// Get the current difficulty level
	MethodServerSettingsDifficultyGet = "minecraft:serversettings/difficulty"
	// Set the difficulty level (peaceful, easy, normal, hard)
	MethodServerSettingsDifficultySet = "minecraft:serversettings/difficulty/set"
	// Get whether allowlist enforcement is enabled
	MethodServerSettingsEnforceAllowlistGet = "minecraft:serversettings/enforce_allowlist"
	// Enable or disable allowlist enforcement
	MethodServerSettingsEnforceAllowlistSet = "minecraft:serversettings/enforce_allowlist/set"
	// Get whether the allowlist is enabled
	MethodServerSettingsUseAllowlistGet = "minecraft:serversettings/use_allowlist"
	// Enable or disable the allowlist
	MethodServerSettingsUseAllowlistSet = "minecraft:serversettings/use_allowlist/set"
	// Get the maximum number of players allowed
	MethodServerSettingsMaxPlayersGet = "minecraft:serversettings/max_players"
	// Set the maximum number of players allowed
	MethodServerSettingsMaxPlayersSet = "minecraft:serversettings/max_players/set"
	// Get the pause when empty timeout in seconds
	MethodServerSettingsPauseWhenEmptyGet = "minecraft:serversettings/pause_when_empty_seconds"
	// Set the pause when empty timeout in seconds
	MethodServerSettingsPauseWhenEmptySet = "minecraft:serversettings/pause_when_empty_seconds/set"
	// Get the player idle timeout in seconds
	MethodServerSettingsPlayerIdleTimeoutGet = "minecraft:serversettings/player_idle_timeout"
	// Set the player idle timeout in seconds
	MethodServerSettingsPlayerIdleTimeoutSet = "minecraft:serversettings/player_idle_timeout/set"
	// Get whether flight is allowed in Survival mode
	MethodServerSettingsAllowFlightGet = "minecraft:serversettings/allow_flight"
	// Set whether flight is allowed in Survival mode
	MethodServerSettingsAllowFlightSet = "minecraft:serversettings/allow_flight/set"
	// Get the server's message of the day
	MethodServerSettingsMOTDGet = "minecraft:serversettings/motd"
	// Set the server's message of the day
	MethodServerSettingsMOTDSet = "minecraft:serversettings/motd/set"
	// Get the spawn protection radius in blocks
	MethodServerSettingsSpawnProtectionGet = "minecraft:serversettings/spawn_protection_radius"
	// Set the spawn protection radius in blocks
	MethodServerSettingsSpawnProtectionSet = "minecraft:serversettings/spawn_protection_radius/set"
	// Get whether players are forced to use the default game mode
	MethodServerSettingsForceGameModeGet = "minecraft:serversettings/force_game_mode"
	// Set whether players are forced to use the default game mode
	MethodServerSettingsForceGameModeSet = "minecraft:serversettings/force_game_mode/set"
	// Get the server's default game mode
	MethodServerSettingsGameModeGet = "minecraft:serversettings/game_mode"
	// Set the server's default game mode (creative, survival, spectator, adventure)
	MethodServerSettingsGameModeSet = "minecraft:serversettings/game_mode/set"
	// Get the view distance in chunks
	MethodServerSettingsViewDistanceGet = "minecraft:serversettings/view_distance"
	// Set the view distance in chunks
	MethodServerSettingsViewDistanceSet = "minecraft:serversettings/view_distance/set"
	// Get the simulation distance in chunks
	MethodServerSettingsSimulationDistanceGet = "minecraft:serversettings/simulation_distance"
	// Set the simulation distance in chunks
	MethodServerSettingsSimulationDistanceSet = "minecraft:serversettings/simulation_distance/set"
	// Get whether the server accepts player transfers
	MethodServerSettingsAcceptTransfersGet = "minecraft:serversettings/accept_transfers"
	// Set whether the server accepts player transfers
	MethodServerSettingsAcceptTransfersSet = "minecraft:serversettings/accept_transfers/set"
	// Get the status heartbeat interval in seconds
	MethodServerSettingsHeartbeatIntervalGet = "minecraft:serversettings/status_heartbeat_interval"
	// Set the status heartbeat interval in seconds
	MethodServerSettingsHeartbeatIntervalSet = "minecraft:serversettings/status_heartbeat_interval/set"
	// Get the operator permission level
	MethodServerSettingsOperatorPermissionLevelGet = "minecraft:serversettings/operator_user_permission_level"
	// Set the operator permission level
	MethodServerSettingsOperatorPermissionLevelSet = "minecraft:serversettings/operator_user_permission_level/set"
	// Get whether online player info is hidden from status queries
	MethodServerSettingsHideOnlinePlayersGet = "minecraft:serversettings/hide_online_players"
	// Set whether online player info is hidden from status queries
	MethodServerSettingsHideOnlinePlayersSet = "minecraft:serversettings/hide_online_players/set"
	// Get whether the server responds to status requests
	MethodServerSettingsStatusRepliesGet = "minecraft:serversettings/status_replies"
	// Set whether the server responds to status requests
	MethodServerSettingsStatusRepliesSet = "minecraft:serversettings/status_replies/set"
	// Get the entity broadcast range percentage
	MethodServerSettingsEntityBroadcastRangeGet = "minecraft:serversettings/entity_broadcast_range"
	// Set the entity broadcast range percentage
	MethodServerSettingsEntityBroadcastRangeSet = "minecraft:serversettings/entity_broadcast_range/set"

	// GameRules management

	// Get all gamerules and their values
	MethodGameRulesGet = "minecraft:gamerules"
	// Update a gamerule value
	MethodGameRulesUpdate = "minecraft:gamerules/update"

	// Notification methods

	// Server notifications

	// Notification sent when the server starts
	NotificationServerStarted = "minecraft:notification/server/started"
	// Notification sent when the server is shutting down
	NotificationServerStopping = "minecraft:notification/server/stopping"
	// Notification sent when the server begins saving
	NotificationServerSaving = "minecraft:notification/server/saving"
	// Notification sent when the server finishes saving
	NotificationServerSaved = "minecraft:notification/server/saved"
	// Periodic server status heartbeat notification
	NotificationServerStatus = "minecraft:notification/server/status"
	// Notification sent when network connection is initialized
	NotificationServerActivity = "minecraft:notification/server/activity"

	// Player notifications

	// Notification sent when a player joins the server
	NotificationPlayerJoined = "minecraft:notification/players/joined"
	// Notification sent when a player leaves the server
	NotificationPlayerLeft = "minecraft:notification/players/left"

	// Operator notifications

	// Notification sent when a player is promoted to operator
	NotificationOperatorAdded = "minecraft:notification/operators/added"
	// Notification sent when a player is demoted from operator
	NotificationOperatorRemoved = "minecraft:notification/operators/removed"

	// Allowlist notifications

	// Notification sent when a player is added to the allowlist
	NotificationAllowlistAdded = "minecraft:notification/allowlist/added"
	// Notification sent when a player is removed from the allowlist
	NotificationAllowlistRemoved = "minecraft:notification/allowlist/removed"

	// IP Ban notifications

	// Notification sent when an IP is added to the ban list
	NotificationIPBanAdded = "minecraft:notification/ip_bans/added"
	// Notification sent when an IP is removed from the ban list
	NotificationIPBanRemoved = "minecraft:notification/ip_bans/removed"

	// Ban notifications

	// Notification sent when a player is added to the ban list
	NotificationBanAdded = "minecraft:notification/bans/added"
	// Notification sent when a player is removed from the ban list
	NotificationBanRemoved = "minecraft:notification/bans/removed"

	// Gamerule notifications

	// Notification sent when a gamerule is updated
	NotificationGameruleUpdated = "minecraft:notification/gamerules/updated"

	// World Notifications

	// Notification sent when a world upgrade started
	NotificationWorldUpgradeStarted = "minecraft:notification/world/upgrade_started"
	// Notification sent when a world upgrade is in progress. Rate limited to 1 notification per second.
	NotificationWorldUpgradeProgress = "minecraft:notification/world/upgrade_progress"
	// Notification sent when a world upgrade is finished
	NotificationWorldUpgradeFinished = "minecraft:notification/world/upgrade_finished"
	// Notification sent when a world upgrade failed
	NotificationWorldUpgradeFailed = "minecraft:notification/world/upgrade_failed"
)
