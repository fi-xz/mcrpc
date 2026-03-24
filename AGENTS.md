# Agent Guidelines for MCRPC

This document provides guidelines for AI coding agents working on the MCRPC (Minecraft Server Management Protocol Client) repository.

## Project Overview

**Language:** Go  
**Go Version:** 1.26.1  
**Purpose:** JSON-RPC 2.0 over WebSocket client for managing Minecraft Java Edition servers  
**Module:** `github.com/fi-xz/mcrpc`

## Build Commands

```bash
# Build the package
go build ./...

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -v -run TestFunctionName ./...

# Run tests with race detector
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Vet code for issues
go vet ./...

# Tidy dependencies
go mod tidy

# Download dependencies
go mod download
```

## Code Style Guidelines

### Imports
- Use standard library imports first
- Then third-party imports (github.com/...)
- Finally, internal package imports
- Group imports with blank lines between groups
- Use `goimports` style

Example:
```go
import (
    "context"
    "crypto/tls"
    "encoding/json"
    
    "github.com/gorilla/websocket"
    "github.com/sourcegraph/jsonrpc2"
    
    "github.com/fi-xz/mcrpc/internal/protocol"
    "github.com/fi-xz/mcrpc/internal/types"
)
```

### Naming Conventions
- **Exported functions/methods:** PascalCase (e.g., `GetServerStatus`, `CreateWithTLS`)
- **Unexported functions:** camelCase (e.g., `createMCRCPClient`, `handleIncoming`)
- **Types:** PascalCase (e.g., `MCRPCClient`, `ServerState`)
- **Constants:** PascalCase for exported, camelCase for unexported
- **Interfaces:** PascalCase with "-er" suffix (e.g., `Notifier`, `Handler`)
- **Test functions:** `Test` + PascalCase (e.g., `TestGetServerStatus`)
- **JSON tags:** camelCase matching protocol spec (e.g., `json:"permissionLevel"`)

### Types
- Use type aliases for exported types from internal packages:
  ```go
  type Player = types.Player
  type Operator = types.Operator
  ```
- Use pointer receivers for methods that modify state
- Use value receivers for methods that only read
- Document exported types with comments

### Error Handling
- Return errors as the last return value
- Check errors immediately and return early
- Don't ignore errors with `_` unless explicitly intended
- Use descriptive error messages
- Wrap errors with context when appropriate

Example:
```go
func (c *MCRPCClient) GetPlayers(ctx context.Context) ([]Player, error) {
    var players []Player
    err := c.JSONRPCConn.Call(ctx, protocol.MethodPlayersGet, nil, &players)
    if err != nil {
        return nil, fmt.Errorf("failed to get players: %w", err)
    }
    return players, nil
}
```

### Function Signatures
- First parameter should be `context.Context` for cancellable operations
- Return results followed by error
- Keep functions focused and single-purpose
- Maximum 4-5 parameters per function

### Struct Tags
- Use JSON tags for all struct fields that are serialized
- Follow the Minecraft protocol specification exactly
- Add comments for non-obvious fields

Example:
```go
type Player struct {
    Name string `json:"name"`
    UUID string `json:"id"` // Note: protocol uses "id" not "uuid"
}
```

### Comments
- All exported functions, types, and constants must have comments
- Comments should start with the name of the item being documented
- Use complete sentences
- Add documentation comments for method constants

Example:
```go
// GetServerStatus retrieves the current server state including online players
// and version information.
func (c *MCRPCClient) GetServerStatus(ctx context.Context) (ServerState, error) {
    // ...
}
```

## Project Structure

```
mcrpc/
├── .github/workflows/    # CI/CD configuration
├── internal/
│   ├── protocol/        # Protocol constants and request types
│   │   ├── methods.go   # Method and notification constants
│   │   └── requests.go  # Request parameter structs
│   └── types/           # Shared data types
│       └── types.go     # Player, Operator, ServerState, etc.
├── *.go                 # API implementations by category
│   ├── mcrpc.go        # Client creation and connection
│   ├── allowlist.go    # Allowlist management
│   ├── ban.go          # Ban and IP ban management
│   ├── ops.go          # Operator management
│   ├── players.go      # Player management
│   ├── server.go       # Server management
│   ├── server_settings.go # Server settings
│   ├── gamerules.go    # Game rules
│   └── notifications.go # Notification handling
└── *_test.go           # Test files
```

## Testing Guidelines

### Test Structure
- One test file per implementation file (`foo.go` → `foo_test.go`)
- Use `createTestClient()` helper for integration tests
- Tests should clean up after themselves (restore original values)
- Use `t.Skip()` when prerequisites aren't met

### Test Naming
- `Test<FunctionName>` for testing specific functions
- `Test<Type_Method>` for method tests
- Use subtests with `t.Run()` for related test cases

### Test Data
- Use test player: `Name: "fi_xz"`, `UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"`
- Don't use real credentials in tests
- Load secrets/certs from `.secrets/` and `.certs/` directories

### Example Test
```go
func TestGetPlayers(t *testing.T) {
    client, ctx := createTestClient(t)
    
    players, err := client.GetPlayers(ctx)
    if err != nil {
        t.Errorf("GetPlayers failed: %v", err)
    }
    
    if players == nil {
        t.Error("Expected non-nil players list")
    }
}
```

## Protocol Implementation

### Method Constants
- Define all method constants in `internal/protocol/methods.go`
- Use format: `Method<Category><Action>` (e.g., `MethodPlayersGet`)
- Use format for notifications: `Notification<Category><Event>` (e.g., `NotificationPlayerJoined`)
- Document each constant with its purpose

### Request Parameters
- Define parameter structs in `internal/protocol/requests.go`
- Use JSON tags matching the protocol specification
- Field names should match the JSON keys in the protocol

### Type Definitions
- Define shared types in `internal/types/types.go`
- Reference: https://minecraft.wiki/w/Minecraft_Server_Management_Protocol#Schemas
- Use `interface{}` for fields that can be multiple types (e.g., `UntypedGameRule.Value`)

## Dependencies

- `github.com/gorilla/websocket v1.4.1` - WebSocket client
- `github.com/sourcegraph/jsonrpc2 v0.2.1` - JSON-RPC 2.0 implementation

Keep dependencies minimal. Only add new dependencies if absolutely necessary.

## Security

- Never commit secrets, certificates, or credentials
- Use `.secrets/` and `.certs/` directories (already in `.gitignore`)
- Support TLS with client certificates
- Validate all inputs when possible

## Common Tasks

### Adding a New Method
1. Add method constant to `internal/protocol/methods.go`
2. Add request params struct if needed to `internal/protocol/requests.go`
3. Implement method in appropriate `*.go` file
4. Add test in corresponding `*_test.go` file

### Adding a New Notification Handler
1. Add constant to `internal/protocol/methods.go`
2. Add handler field to `MCRPCClient` struct in `mcrpc.go`
3. Add case to switch in `notifications.go`
4. Parse JSON params appropriately
5. Call handler if set

### Adding a New Type
1. Add to `internal/types/types.go` if shared across packages
2. Create type alias in relevant implementation file
3. Use JSON tags matching the protocol specification

## Code Review Checklist

Before submitting changes:
- [ ] Code compiles without errors: `go build ./...`
- [ ] All tests pass: `go test ./...`
- [ ] Code is formatted: `go fmt ./...`
- [ ] No vet issues: `go vet ./...`
- [ ] Dependencies are tidy: `go mod tidy`
- [ ] Comments are added for exported items
- [ ] Error handling is proper (no ignored errors)
- [ ] No sensitive data is committed

## Resources

- **Protocol Documentation:** https://minecraft.wiki/w/Minecraft_Server_Management_Protocol
- **Go Style Guide:** https://go.dev/doc/effective_go
- **Project Repository:** https://github.com/fi-xz/mcrpc
