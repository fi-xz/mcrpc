# MCRPC - Minecraft Server Management Protocol Client

[![CI](https://github.com/fi-xz/mcrpc/actions/workflows/ci.yml/badge.svg)](https://github.com/fi-xz/mcrpc/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fi-xz/mcrpc)](https://goreportcard.com/report/github.com/fi-xz/mcrpc)
[![GoDoc](https://godoc.org/github.com/fi-xz/mcrpc?status.svg)](https://godoc.org/github.com/fi-xz/mcrpc)

A Go client library for the [Minecraft Server Management Protocol](https://minecraft.wiki/w/Minecraft_Server_Management_Protocol) (MSMP), providing JSON-RPC 2.0 over WebSocket support for managing Minecraft Java Edition dedicated servers.

> [!Warning]
>
> Part of this project is handled by the AI.
>
> Most code structure is mostly handled by myself. Only repetitive tasks, tests, and long length codes are written by the AI.
>
> While I try my best not to make vibe-coded library, there may some errors occur.

## Features

- **Full Protocol Support**: Complete implementation of all MSMP methods
- **Secure Connections**: TLS support with client certificates
- **Event Notifications**: Real-time event handling for player joins, server status, etc.
- **Type-Safe API**: Strongly typed request/response structures
- **Concurrent-Safe**: Thread-safe operations for multi-goroutine use

## Installation

```bash
go get github.com/fi-xz/mcrpc
```

## Quick Start

```go
package main

import (
    "context"
    "crypto/tls"
    "fmt"
    "log"

    "github.com/fi-xz/mcrpc"
)

func main() {
    ctx := context.Background()

    // Load certificate (only if the server requires client authentication)
    cert, err := tls.LoadX509KeyPair("cert.crt", "cert.pem")
    if err != nil {
        log.Fatal(err)
    }

    // Build the client. No I/O happens here, so handlers registered now are
    // guaranteed to be in place before any notification can arrive.
    client := mcrpc.New("localhost:8080", "your-secret",
        mcrpc.WithClientCertificate(cert),
        mcrpc.WithInsecureSkipVerify(), // self-signed certificates only
        mcrpc.WithHandler(mcrpc.Handler{
            OnPlayerJoined: func(player mcrpc.Player) {
                log.Printf("Player joined: %s", player.Name)
            },
        }),
    )

    // Connect. ctx bounds the whole session: cancelling it closes the client.
    if err := client.Start(ctx); err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Get server status
    status, err := client.GetServerStatus(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Server: %s, Players: %d\n", status.Version.Name, len(status.Players))
}
```

For a plaintext connection, drop the TLS options:

```go
client := mcrpc.New("localhost:8080", "your-secret")
```

`NewHostPort(host, port, secret, opts...)` is available when the host and port
are held separately.

### Lifecycle

| Call | Effect |
|---|---|
| `New` / `NewHostPort` | Builds the value. No connection, no goroutine. |
| `Start(ctx)` | Dials and begins dispatching. `ErrAlreadyStarted` if running. |
| `Close()` | Ends the session. Idempotent, and safe before `Start`. |
| `Start(ctx)` again | Reconnects with the same configuration and handler. |
| `IsRunning()` | Whether a connection is currently held. |
| `DisconnectNotify()` | Channel closed when the *current* session ends. |

`ctx` passed to `Start` bounds the session and is the context handed to every
notification callback; cancelling it closes the connection. To bound only the
handshake, use `WithHandshakeTimeout` and pass a longer-lived context to
`Start`.

### Server versions

The management API is versioned separately from Minecraft, and it is the API
version that protocol behaviour follows:

```go
version, err := client.APIVersion(ctx) // "1.0.0", "3.0.0", …
```

| API | Minecraft | Notable for this library |
|---|---|---|
| 1.0.0 | 1.21.9 | Game rule values are strings; keys are camelCase |
| 3.0.0 | 26.2 | Game rule values are native JSON; keys are namespaced. `notification/server/activity` added |
| 3.1.0 | 26.3 | `notification/world/upgrade_*` added |

Every method and payload shape a server supports can be read from its own
OpenRPC document via `rpc.discover`; `APIVersion` reads the version field out
of it.

## API Overview

### Server Management

- `GetServerStatus()` - Get current server state
- `SaveServer(ctx, flush)` - Save world data
- `StopServer(ctx)` - Stop the server
- `SendSystemMessage(ctx, message)` - Broadcast messages

### Player Management

- `GetPlayers(ctx)` - List online players
- `KickPlayers(ctx, players...)` - Kick players from server

### Access Control

- **Allowlist**: `GetAllowlist`, `SetAllowlist`, `AddAllowlist`, `RemoveAllowlist`, `ClearAllowlist`
- **Bans**: `GetBanlist`, `SetBanlist`, `AddBanlist`, `RemoveBanlist`, `ClearBanlist`
- **IP Bans**: `GetIPBanlist`, `SetIPBanlist`, `AddIPBanlist`, `RemoveIPBanlist`, `ClearIPBanlist`
- **Operators**: `GetOperators`, `SetOperators`, `AddOperators`, `RemoveOperators`, `ClearOperators`

`Add*`, `Remove*`, and `KickPlayers` are variadic, so a single item needs no
slice literal. `Set*` keeps taking a slice, because replacing a whole list is
naturally a list operation:

```go
client.AddAllowlist(ctx, mcrpc.PlayerByName("fi_xz"))
client.RemoveIPBanlist(ctx, "1.2.3.4", "5.6.7.8")
client.SetAllowlist(ctx, roster)
```

### Server Settings

All settings support both Get and Set methods:

- `Autosave`, `Difficulty`, `MaxPlayers`, `MOTD`
- `GameMode`, `ViewDistance`, `SimulationDistance`
- `AllowFlight`, `Allowlist`, `ForceGameMode`
- `SpawnProtectionRadius`, `PlayerIdleTimeout`
- `AcceptTransfers`, `StatusReplies`, and more

### Game Rules

- `GetGameRules(ctx)` - Get all game rules
- `UpdateGameRule(ctx, rule)` - Update a game rule

Game rule values are untyped on the wire, so accessors are provided:

```go
rule, err := client.UpdateGameRule(ctx, mcrpc.BoolRule("minecraft:keep_inventory", true))
if err != nil {
    log.Fatal(err)
}

if enabled, ok := rule.Bool(); ok {
    log.Printf("minecraft:keep_inventory is now %v", enabled)
}
```

`BoolRule`, `IntRule`, and `StringRule` build updates; `Bool()`, `Int()`, and
`StringValue()` read them back, each returning an `ok` flag.

> [!Important]
>
> **Game rules changed shape in 1.21.11.** Both the key and the value
> representation differ, verified against live 1.21.9 and 26.2 servers:
>
> | | 1.21.9 – 1.21.10 | 1.21.11 and later |
> |---|---|---|
> | Key | `keepInventory` | `minecraft:keep_inventory` |
> | Value on the wire | `"true"`, `"3"` (strings) | `true`, `3` (native JSON) |
>
> Rules were also renamed to match the "Edit Game Rules" menu, and **some had
> their values inverted**, so a key is not always a straight rename.

`Bool()`, `Int()`, and `StringValue()` accept either representation, so reading
is version-independent. For writing, build the update from the rule the server
sent and the representation is matched for you:

```go
rules, err := client.GetGameRules(ctx)
// ... find the rule you want ...

updated, err := client.UpdateGameRule(ctx, rule.WithBool(true))
```

`WithBool`, `WithInt`, and `WithString` mirror whatever the server used;
`UsesStringValues()` reports which it is. The standalone `BoolRule`, `IntRule`,
and `StringRule` constructors send exactly what you give them, so use those only
when you know the server's representation.

### Constructing values

Helpers cover the awkward literals:

```go
mcrpc.LiteralMessage("Server restarting in 5 minutes")
mcrpc.TranslatableMessage("chat.type.text", "fi_xz", "hello")

mcrpc.PlayerByName("fi_xz")
mcrpc.PlayerByUUID("a0d8c884-2a79-4c95-8617-a51d27a427ec")

// Ban expiry. Leave Expires empty for a permanent ban.
ban := mcrpc.UserBan{
    Player:  mcrpc.PlayerByName("fi_xz"),
    Reason:  "griefing",
    Expires: mcrpc.BanUntil(time.Now().Add(24 * time.Hour)),
}

if deadline, ok := ban.ExpiresAt(); ok {
    log.Printf("ban lifts at %s", deadline)
}
```

### Error handling

Every method returns errors wrapping `*mcrpc.Error`, which carries the failing
method and the server's JSON-RPC error code:

```go
players, err := client.GetPlayers(ctx)
if err != nil {
    var rpcErr *mcrpc.Error
    if errors.As(err, &rpcErr) {
        log.Printf("%s failed with code %d: %s", rpcErr.Method, rpcErr.Code, rpcErr.Message)
    }
    if errors.Is(err, mcrpc.ErrNotConnected) {
        // the client was never connected, or has been closed
    }
}
```

The original `*jsonrpc2.Error` stays reachable through `errors.As` if you need it.

### Tracing the wire

`WithTrace` reports every JSON-RPC message with its params and result exactly as
serialised. Use it to confirm what a struct tag actually produces, what an
omitted field looks like on the wire, and what JSON type an untyped value
arrives as:

```go
client := mcrpc.New(addr, secret, mcrpc.WithTrace(func(m mcrpc.TraceMessage) {
    log.Print(m)
}))
```

```
-> minecraft:allowlist/add #1 {"add":[{"name":"fi_xz","id":""}]}
<- minecraft:allowlist/add #1 [{"name":"fi_xz","id":"a0d8c884-…"}]
<- notify minecraft:notification/allowlist/added [{"name":"fi_xz","id":"a0d8c884-…"}]
-> minecraft:players/kick #2 {"kick":[…]}
<- minecraft:players/kick #2 error -32602 Invalid params
```

`TraceMessage` exposes `Direction`, `Method`, `ID`, `Notification`, `Params`,
`Result`, `ErrorCode`, and `ErrorMessage`; `String()` produces the lines above.
The callback runs on the connection's read and write paths, so it must not
block or call back into the client. Trace output should be treated as
sensitive.

The integration tests wire this to `TEST_TRACE`:

```bash
TEST_TRACE=1 go test -v -run TestBanlist ./...
```

### Notifications

Callbacks are passed to `New` as a `Handler`, so they are in place before the
connection exists:

```go
client := mcrpc.New(addr, secret, mcrpc.WithHandler(mcrpc.Handler{
    OnPlayerJoined: func(player mcrpc.Player) {
        log.Printf("Player joined: %s", player.Name)
    },
    OnServerStatus: func(status mcrpc.ServerState) {
        log.Printf("Players online: %d", len(status.Players))
    },
    OnError: func(method string, err error) {
        log.Printf("could not decode %s: %v", method, err)
    },
}))
```

The `Handler` is copied into the client, so there is no shared mutable state for
the receive goroutine to race against, and the same value can be reused across
reconnects or across clients.

`OnNotification` receives the raw params of every notification, which is the
escape hatch when a payload needs handling this package does not provide.

Available handlers:

- Server: `OnServerStarted`, `OnServerStopping`, `OnServerSaving`, `OnServerSaved`,
  `OnServerStatus`, `OnServerActivity` (API 3.0.0+, rate limited to one per 30s)
- Players: `OnPlayerJoined`, `OnPlayerLeft`
- Operators: `OnOperatorAdded`, `OnOperatorRemoved`
- Allowlist: `OnAllowlistAdded`, `OnAllowlistRemoved`
- Bans: `OnBanAdded`, `OnBanRemoved`
- IP Bans: `OnIPBanAdded`, `OnIPBanRemoved`
- Game Rules: `OnGameruleUpdated`
- World (API 3.1.0+, Minecraft 26.3): `OnWorldUpgradeStarted`,
  `OnWorldUpgradeProgress` (0–1, rate limited to one per second),
  `OnWorldUpgradeFinished`, `OnWorldUpgradeFailed`. These fire while the server
  boots, before it finishes coming up, so register them and connect early
- Catch-all: `OnNotification` (every notification, with raw params), `OnError`
  (a notification whose params could not be decoded)

### Reconnecting

The client keeps its configuration and handler across sessions, so reconnecting
is `Start` again:

```go
for {
    if err := client.Start(ctx); err != nil {
        log.Printf("connect failed: %v", err)
    } else {
        <-client.DisconnectNotify()
    }

    select {
    case <-ctx.Done():
        return
    case <-time.After(5 * time.Second):
    }
}
```

## Configuration

### Enable management protocol on your server

Modify `server.properties`:

```properties
management-server-enabled=true
management-server-host=0.0.0.0
management-server-port=8080
management-server-secret=your-40-character-secret-here
management-server-tls-enabled=true
management-server-tls-keystore=keystore.p12
management-server-tls-keystore-password=your-password
```

### TLS Certificate

Generate a PKCS12 keystore:

```bash
keytool -genkeypair -alias mcrpc -keyalg RSA -keysize 2048 \
  -storetype PKCS12 -keystore keystore.p12 -validity 3650
```

## Testing

Run the test suite:

```bash
# Run all tests
go test -v ./...

# Run with race detector
go test -race ./...

# Run with coverage
go test -cover ./...
```

## Requirements

- Go 1.26.1 or later (see `go.mod`)
- Minecraft Java Edition 1.21.9+ with management protocol enabled

## Compatibility notes

The API is pre-v1. Connection setup was reworked to remove a data race on
notification handlers; the changes below are breaking.

| Before | Now |
|---|---|
| `Create(ctx, host, port, secret)` | `NewHostPort(host, port, secret)` + `Start(ctx)` |
| `CreateWithTLS(ctx, host, port, secret, &cert, insecure)` | `NewHostPort(host, port, secret, WithClientCertificate(cert), WithInsecureSkipVerify())` + `Start(ctx)` |
| `MCRPCClient` | `Client` |
| `client.OnPlayerJoined = fn` | `WithHandler(Handler{OnPlayerJoined: fn})` |
| `client.IsClosed()` | `!client.IsRunning()` |
| `client.JSONRPCConn`, `client.WebsocketConn` | unexported |
| `AddAllowlist(ctx, []Player{p})` | `AddAllowlist(ctx, p)` |

Behavioural changes:

- **The connection is no longer opened by the constructor.** `New` performs no
  I/O; nothing is dialled until `Start`. This is what makes handler
  registration race-free, and it lets the same client reconnect after `Close`.
- **The context passed to `Start` bounds the session.** Cancelling it closes the
  connection. Previously the context given to `Create` only reached notification
  callbacks and never closed anything, so a short-lived context silently handed
  every later callback an already-cancelled context.
- **TLS follows the options, not the certificate.** `WithTLS(nil)` connects over
  `wss` without a client certificate; previously TLS was only used when a
  certificate was supplied.
- **Permanent bans omit `expires`** rather than sending an empty string.
- **List parameters serialise as `[]`, never `null`.**

Remaining known issues and their rationale are tracked in
[docs/ux-proposal.md](docs/ux-proposal.md).

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Resources

- [Minecraft Server Management Protocol Documentation](https://minecraft.wiki/w/Minecraft_Server_Management_Protocol)
- [Go Documentation](https://godoc.org/github.com/fi-xz/mcrpc)
