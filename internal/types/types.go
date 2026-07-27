// Package types contains data types used by the Minecraft Server Management Protocol.
// See https://minecraft.wiki/w/Minecraft_Server_Management_Protocol#Schemas for reference.
package types

// UntypedGameRule represents a game rule value that can be of any type (string, bool, or int).
type UntypedGameRule struct {
	Value any    `json:"value"` // Can be string, bool, or int
	Key   string `json:"key"`
}

// IncomingIPBan represents an IP ban with associated player information.
type IncomingIPBan struct {
	IPBan
	Player Player `json:"player"`
}

// SystemMessage represents a message to be sent to players on the server.
type SystemMessage struct {
	ReceivingPlayers []Player `json:"receivingPlayers"` // Players who will receive this message
	Overlay          bool     `json:"overlay"`          // Whether to display as an overlay
	Message          Message  `json:"message"`          // The message content
}

// KickPlayer represents a player to be kicked with a custom message.
type KickPlayer struct {
	Player  Player  `json:"player"`  // The player to kick
	Message Message `json:"message"` // The kick message
}

// IPBan represents a banned IP address.
type IPBan struct {
	Reason  string `json:"reason"`            // Reason for the ban
	Expires string `json:"expires,omitempty"` // Expiration time in ISO 8601 format, omitted if permanent
	IP      string `json:"ip"`                // The banned IP address
	Source  string `json:"source"`            // Who issued the ban
}

// TypedGameRule represents a game rule with a known type.
type TypedGameRule struct {
	UntypedGameRule
	Type string `json:"type"` // Type of the value: "integer" or "boolean"
}

// UserBan represents a banned player.
type UserBan struct {
	Reason  string `json:"reason"`            // Reason for the ban
	Expires string `json:"expires,omitempty"` // Expiration time in ISO 8601 format, omitted if permanent
	Source  string `json:"source"`            // Who issued the ban
	Player  Player `json:"player"`            // The banned player
}

// Message represents a message that can be sent to players.
// It can be either translatable or a literal string.
type Message struct {
	Translatable       string   `json:"translatable,omitempty"`       // Translation key
	TranslatableParams []string `json:"translatableParams,omitempty"` // Parameters for translation
	Literal            string   `json:"literal,omitempty"`            // Literal message text
}

// Version represents the Minecraft server version information.
type Version struct {
	Protocol int    `json:"protocol"` // Protocol version number
	Name     string `json:"name"`     // Version name (e.g., "1.21.4")
}

// ServerState represents the current state of the Minecraft server.
type ServerState struct {
	Players []Player `json:"players"` // List of currently online players
	Started bool     `json:"started"` // Whether the server is running
	Version Version  `json:"version"` // Server version information
}

// Operator represents a player with operator (admin) privileges.
type Operator struct {
	PermissionLevel     int    `json:"permissionLevel"`     // Operator permission level (0-4)
	BypassesPlayerLimit bool   `json:"bypassesPlayerLimit"` // Whether operator can bypass player limit
	Player              Player `json:"player"`              // The operator player
}

// Player represents a Minecraft player.
type Player struct {
	Name string `json:"name"` // Player's username
	UUID string `json:"id"`   // Player's UUID (note: field is named "id" in the protocol)
}
