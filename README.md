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

    // Load certificate (if using TLS)
    cert, err := tls.LoadX509KeyPair("cert.crt", "cert.pem")
    if err != nil {
        log.Fatal(err)
    }

    // Create client
    client, err := mcrpc.CreateWithTLS(ctx, "localhost", 8080, "your-secret", &cert, true)
    if err != nil {
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

## API Overview

### Server Management

- `GetServerStatus()` - Get current server state
- `SaveServer(ctx, flush)` - Save world data
- `StopServer(ctx)` - Stop the server
- `SendSystemMessage(ctx, message)` - Broadcast messages

### Player Management

- `GetPlayers(ctx)` - List online players
- `KickPlayers(ctx, players)` - Kick players from server

### Access Control

- **Allowlist**: `GetAllowlist`, `SetAllowlist`, `AddAllowlist`, `RemoveAllowlist`, `ClearAllowlist`
- **Bans**: `GetBanlist`, `SetBanlist`, `AddBanlist`, `RemoveBanlist`, `ClearBanlist`
- **IP Bans**: `GetIPBanlist`, `SetIPBanlist`, `AddIPBanlist`, `RemoveIPBanlist`, `ClearIPBanlist`
- **Operators**: `GetOperators`, `SetOperators`, `AddOperators`, `RemoveOperators`, `ClearOperators`

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

### Notifications

Register handlers for server events:

```go
client.OnPlayerJoined = func(player mcrpc.Player) {
    log.Printf("Player joined: %s", player.Name)
}

client.OnServerStatus = func(status mcrpc.ServerState) {
    log.Printf("Players online: %d", len(status.Players))
}
```

Available handlers:

- Server: `OnServerStarted`, `OnServerStopping`, `OnServerSaving`, `OnServerSaved`, `OnServerStatus`, `OnServerActivity`
- Players: `OnPlayerJoined`, `OnPlayerLeft`
- Operators: `OnOperatorAdded`, `OnOperatorRemoved`
- Allowlist: `OnAllowlistAdded`, `OnAllowlistRemoved`
- Bans: `OnBanAdded`, `OnBanRemoved`
- IP Bans: `OnIPBanAdded`, `OnIPBanRemoved`
- Game Rules: `OnGameruleUpdated`
- World: `OnWorldUpgradeStarted`, `OnWorldUpgradeProgress`, `OnWorldUpgradeFinished`, `OnWorldUpgradeFailed`

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

- Go 1.21 or later
- Minecraft Java Edition 1.21.9+ with management protocol enabled

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Resources

- [Minecraft Server Management Protocol Documentation](https://minecraft.wiki/w/Minecraft_Server_Management_Protocol)
- [Go Documentation](https://godoc.org/github.com/fi-xz/mcrpc)
