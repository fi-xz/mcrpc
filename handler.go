package mcrpc

import "encoding/json"

// Handler is the set of callbacks invoked for notifications pushed by the
// server. Every field is optional; a nil callback means the notification is
// ignored.
//
// A Handler is passed to New via WithHandler and copied into the client, so
// mutating your own copy afterwards has no effect on a running client. This is
// deliberate: the connection's receive goroutine reads these callbacks
// concurrently, and copying them once at construction keeps that read
// race-free.
//
// The context passed to a callback is the one given to Start. It is cancelled
// when the session ends, so a callback that makes further calls should pass it
// along and expect cancellation.
type Handler struct {
	// OnNotification is called for every notification, before any specific
	// callback below. It receives the raw params, which are nil if the
	// notification carried none.
	OnNotification func(method string, params json.RawMessage)

	// OnError is called when a notification arrives but its params cannot be
	// decoded into the expected type. Without it such notifications are
	// dropped silently.
	OnError func(method string, err error)

	// Server notifications
	OnServerStarted  func()
	OnServerStopping func()
	OnServerSaving   func()
	OnServerSaved    func()
	OnServerStatus   func(status ServerState)
	OnServerActivity func()

	// Player notifications
	OnPlayerJoined func(player Player)
	OnPlayerLeft   func(player Player)

	// Operator notifications
	OnOperatorAdded   func(operator Operator)
	OnOperatorRemoved func(operator Operator)

	// Allowlist notifications
	OnAllowlistAdded   func(player Player)
	OnAllowlistRemoved func(player Player)

	// Ban notifications
	OnBanAdded   func(ban UserBan)
	OnBanRemoved func(player Player)

	// IP ban notifications
	OnIPBanAdded   func(ban IPBan)
	OnIPBanRemoved func(ip string)

	// Gamerule notifications
	OnGameruleUpdated func(gamerule TypedGameRule)

	// World notifications
	OnWorldUpgradeStarted  func()
	OnWorldUpgradeProgress func(progress float64)
	OnWorldUpgradeFinished func()
	OnWorldUpgradeFailed   func(reason string)
}
