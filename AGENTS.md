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
- **Exported functions/methods:** PascalCase (e.g., `GetServerStatus`, `WithHandler`)
- **Unexported functions:** camelCase (e.g., `createClient`, `handleIncoming`)
- **Types:** PascalCase, and never stuttering with the package name — `Client`,
  not `MCRPCClient` (e.g., `Client`, `ServerState`)
- **Context parameter:** always `ctx`, never `context` — naming it `context`
  shadows the package and blocks `context.WithTimeout` inside the function
- Don't shadow builtins (`max`, `min`, `len`, `cap`) with parameter names
- **Constants:** PascalCase for exported, camelCase for unexported
- **Interfaces:** PascalCase with "-er" suffix (e.g., `Notifier`, `Handler`)
- **Test functions:** `Test` + PascalCase (e.g., `TestGetServerStatus`)
- **JSON tags:** camelCase matching protocol spec (e.g., `json:"permissionLevel"`)

### Types
- Every type reachable from an exported signature must have an alias in `types.go`,
  the single re-export point for `internal/types`:
  ```go
  type Player = types.Player
  type Operator = types.Operator
  ```
- Never name an `internal/...` type in an exported signature. Consumers cannot
  import internal packages, so such a parameter is impossible to construct.
  `api_public_test.go` (external test package) enforces this at compile time.
- Use pointer receivers for methods that modify state
- Use value receivers for methods that only read
- Document exported types with comments

### Error Handling
- Return errors as the last return value
- Check errors immediately and return early
- Don't ignore errors with `_` unless explicitly intended
- Use descriptive error messages
- Wrap errors with context when appropriate
- **Never call `c.JSONRPCConn.Call` directly.** Route every remote call through
  the `(*Client).call` helper in `errors.go`, which normalises failures into
  `*Error` and handles the not-connected case.
- Helpers that can fail on malformed input return `(value, ok)` rather than an
  error, matching the type-assertion idiom (`Bool`, `Int`, `ExpiresAt`)

Example:
```go
func (c *Client) GetPlayers(ctx context.Context) ([]Player, error) {
    var players []Player
    err := c.call(ctx, protocol.MethodPlayersGet, nil, &players)
    return players, err
}
```

### Function Signatures
- First parameter should be `context.Context` for cancellable operations
- Return results followed by error
- Keep functions focused and single-purpose
- Maximum 4-5 parameters per function; past that, take an `Option`
- `Add*`, `Remove*`, and kick-style methods are variadic; `Set*` takes a slice,
  because replacing a whole list is naturally a list operation
- Wrap every list-valued request parameter in `nonNilSlice` so it serialises as
  `[]` rather than `null`

### Client Lifecycle
- `New` must stay free of I/O and goroutines. Everything that connects belongs
  in `Start`. This is what makes handler registration race-free, so don't move
  dialling back into the constructor.
- The context passed to `Start` bounds the session. `jsonrpc2`'s read loop does
  not watch it, so `client.go` runs a watchdog goroutine that turns
  cancellation and server hangup into a `Close`.
- Session state (`rpc`, `cancel`) is guarded by `Client.mu`. Read it through
  `conn()`, never directly.
- `Close` must remain idempotent and must not prevent a later `Start`.

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
func (c *Client) GetServerStatus(ctx context.Context) (ServerState, error) {
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
│       ├── types.go     # Player, Operator, ServerState, etc.
│       └── accessors.go # Methods on those types (ExpiresAt, Bool, Int, ...)
├── *.go                 # API implementations by category
│   ├── types.go        # Public aliases re-exporting internal/types
│   ├── errors.go       # *Error, ErrNotConnected, (*Client).call helper
│   ├── helpers.go      # Constructors: LiteralMessage, BoolRule, BanUntil, ...
│   ├── doc.go          # Package documentation
│   ├── client.go       # Client type, New, Start, Close, session lifecycle
│   ├── options.go      # Option type and With* constructors
│   ├── handler.go      # Handler struct (notification callbacks)
│   ├── trace.go        # WithTrace plumbing: TraceMessage, jsonrpc2 hooks
│   ├── allowlist.go    # Allowlist management
│   ├── ban.go          # Ban and IP ban management
│   ├── ops.go          # Operator management
│   ├── players.go      # Player management
│   ├── server.go       # Server management
│   ├── server_settings.go # Server settings
│   ├── gamerules.go    # Game rules
│   └── notifications.go # Notification handling
└── *_test.go           # Test files
    ├── client_test.go             # Lifecycle, against an httptest fake server
    ├── trace_test.go              # WithTrace, against the fake server
    ├── wire_test.go               # Live-server checks on the wire format
    ├── wire_notifications_test.go # Live-server checks on notification delivery
    └── wire_schema_test.go        # Conformance against the server's rpc.discover
```

## Testing Guidelines

### Test Structure
- One test file per implementation file (`foo.go` → `foo_test.go`)
- Prefer tests that need no Minecraft server. `client_test.go` has a
  `newFakeServer` helper — an `httptest` WebSocket endpoint speaking JSON-RPC —
  which covers connect, close, restart, notification dispatch, and error
  mapping. Reach for it before writing another skipped integration test.
- Use `createTestClient()` helper for the integration tests that genuinely need
  a live server; they `t.Skip` in CI and therefore prove nothing there
- Tests should clean up after themselves (restore original values)
- Use `t.Skip()` when prerequisites aren't met
- When a live run is needed, `TEST_TRACE=1 go test -v ...` prints the exact JSON
  exchanged. Reach for it before guessing at a wire-format mismatch — a success
  response does not mean the server read the payload the way you intended.
- A fake-server test only replays the assumptions you wrote into it. Nine
  notification handlers were dead for the life of this package because both the
  client and its tests assumed a single object where the server sends a list.
  Any claim about the wire needs `wire_test.go` / `wire_notifications_test.go`
  against a real server.

### Running against a live server

```bash
TEST_HOST=127.0.0.1 TEST_PORT=8080 USE_TLS=true   TEST_TLS_SERVER_NAME=<name-on-the-certificate>   TEST_TRACE=1 go test -v -race ./...

# Dump the server's own API description while you are there.
SCHEMA_OUT=schema.json go test -v -run TestWireNotificationsTakeOneParameter ./...
```

- `TEST_SECRET`, or `.secrets/secret.txt` when it is unset
- `TEST_TLS_SERVER_NAME` verifies against the certificate's name while dialling
  somewhere else, which is what a loopback-bound server with a real certificate
  needs. Prefer it to `TEST_TLS_INSECURE=true`, which turns verification off.
- `TEST_TLS_CA=path`, or `.certs/ca.crt` when it exists, adds a trusted root.
  Note that `.certs/cert.crt` is the *server's* identity, not a client
  certificate.

Integration tests mutate server state and restore it in `t.Cleanup`. Point them
at a throwaway server.

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
- **The server describes its own API.** `rpc.discover` returns an OpenRPC
  document listing every method, its parameter names, and its payload schemas.
  Check it before guessing. `TestWireRequestParameterNames` and
  `TestWireNotificationsTakeOneParameter` do this automatically against a live
  server; add new setters to that table.
  `allow_flight/set` shipped as `"allowed"` — the schema's name for the
  *result* — while the parameter is `"allow"`, so it failed every time.
- **Notification params are a positional argument list.** Every notification
  declares exactly one parameter, so the payload is element 0:
  `[{"name":…}]`, not `{"name":…}`. Decode with `decodeParam`. Reading the
  params directly as the payload left twelve handlers silently dead.
- **Behaviour differs across the supported range (1.21.9+).** Game rule keys and
  value representations both changed in 1.21.11. Verify wire-format claims
  against more than one server version before writing them down.
- **Key behaviour on the management API version, not the Minecraft version.**
  `(*Client).APIVersion` reads it from `rpc.discover`: 1.0.0 ships with 1.21.9,
  3.0.0 with 26.2, 3.1.0 with 26.3. A method missing from the schema may simply
  be newer than the server, as the `world/upgrade_*` notifications are.

### Type Definitions
- Define shared types in `internal/types/types.go`
- Reference: https://minecraft.wiki/w/Minecraft_Server_Management_Protocol#Schemas
- Use `any` for fields that can be multiple types (e.g., `UntypedGameRule.Value`),
  and pair it with typed accessors in `internal/types/accessors.go`
- Methods on protocol types belong in `internal/types/accessors.go`, not the root
  package: an alias is the *same* type, so Go forbids attaching methods to it
  from `package mcrpc`. They still surface to users through the alias.
- Use `omitempty` on fields the protocol treats as optional or mutually
  exclusive (`Message.Literal` vs `Message.Translatable`, `UserBan.Expires`).
  Do *not* use it on booleans where `false` is a meaningful value.

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
2. Add the callback field to the `Handler` struct in `handler.go`
3. Add case to switch in `notifications.go`
4. Parse JSON params appropriately
5. Call handler if set

### Adding a New Type
1. Add to `internal/types/types.go` if shared across packages
2. Add the alias to `types.go` if the type appears in any exported signature
   (directly, or as a field of another exported type)
3. Use JSON tags matching the protocol specification
4. Add it to `api_public_test.go` so the public reachability stays covered

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
- [ ] No `internal/...` type appears in an exported signature
- [ ] Remote calls go through `(*Client).call`, not `c.JSONRPCConn.Call`
- [ ] New behaviour is covered by a test that runs without a live server —
      integration tests `t.Skip` in CI and prove nothing there

## Resources

- **Protocol Documentation:** https://minecraft.wiki/w/Minecraft_Server_Management_Protocol
- **Go Style Guide:** https://go.dev/doc/effective_go
- **Project Repository:** https://github.com/fi-xz/mcrpc
