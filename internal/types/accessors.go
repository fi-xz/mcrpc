package types

import (
	"encoding/json"
	"strconv"
	"time"
)

// BanTimeFormat is the layout the management protocol uses for ban expiry
// timestamps.
const BanTimeFormat = time.RFC3339

// ExpiresAt reports when the ban lapses. The second return value is false for
// a permanent ban, and also when the server sent a timestamp this package
// cannot parse.
func (b UserBan) ExpiresAt() (time.Time, bool) {
	return parseExpiry(b.Expires)
}

// IsPermanent reports whether the ban never lapses.
func (b UserBan) IsPermanent() bool {
	return b.Expires == ""
}

// ExpiresAt reports when the IP ban lapses. The second return value is false
// for a permanent ban, and also when the server sent a timestamp this package
// cannot parse.
func (b IPBan) ExpiresAt() (time.Time, bool) {
	return parseExpiry(b.Expires)
}

// IsPermanent reports whether the IP ban never lapses.
func (b IPBan) IsPermanent() bool {
	return b.Expires == ""
}

func parseExpiry(raw string) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	parsed, err := time.Parse(BanTimeFormat, raw)
	if err != nil {
		return time.Time{}, false
	}
	return parsed, true
}

// UsesStringValues reports whether the server represented this rule's value as
// a JSON string rather than a native boolean or number.
//
// Servers up to 1.21.10 send every game rule value as a string; 1.21.11 and
// later send native types. An update has to match, so use WithBool, WithInt, or
// WithString to build one from a rule the server sent.
func (g UntypedGameRule) UsesStringValues() bool {
	_, isString := g.Value.(string)
	return isString
}

// WithBool builds an update for this rule carrying value, matching whichever
// representation the server used for it.
func (g UntypedGameRule) WithBool(value bool) UntypedGameRule {
	if g.UsesStringValues() {
		return UntypedGameRule{Key: g.Key, Value: strconv.FormatBool(value)}
	}
	return UntypedGameRule{Key: g.Key, Value: value}
}

// WithInt builds an update for this rule carrying value, matching whichever
// representation the server used for it.
func (g UntypedGameRule) WithInt(value int) UntypedGameRule {
	if g.UsesStringValues() {
		return UntypedGameRule{Key: g.Key, Value: strconv.Itoa(value)}
	}
	return UntypedGameRule{Key: g.Key, Value: value}
}

// WithString builds an update for this rule carrying value verbatim.
func (g UntypedGameRule) WithString(value string) UntypedGameRule {
	return UntypedGameRule{Key: g.Key, Value: value}
}

// Bool returns the game rule value as a boolean. The second return value is
// false if the rule does not hold a boolean.
//
// Servers up to 1.21.10 send booleans as the strings "true" and "false";
// 1.21.11 and later send native JSON booleans. Both are accepted.
func (g UntypedGameRule) Bool() (bool, bool) {
	switch v := g.Value.(type) {
	case bool:
		return v, true
	case string:
		parsed, err := strconv.ParseBool(v)
		return parsed, err == nil
	default:
		return false, false
	}
}

// Int returns the game rule value as an integer. The second return value is
// false if the rule does not hold an integer.
//
// Servers disagree about the representation: 1.21.9 and 1.21.10 send game rule
// values as strings ("3"), while 1.21.11 and later send native JSON numbers,
// which decode into any as float64. Both are accepted.
func (g UntypedGameRule) Int() (int, bool) {
	switch v := g.Value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		if v != float64(int(v)) {
			return 0, false
		}
		return int(v), true
	case json.Number:
		parsed, err := v.Int64()
		return int(parsed), err == nil
	case string:
		parsed, err := strconv.Atoi(v)
		return parsed, err == nil
	default:
		return 0, false
	}
}

// StringValue returns the game rule value as a string. The second return value
// is false if the rule does not hold a string.
//
// It is deliberately not named String, so that UntypedGameRule does not
// accidentally satisfy fmt.Stringer with a non-conforming signature.
func (g UntypedGameRule) StringValue() (string, bool) {
	v, ok := g.Value.(string)
	return v, ok
}
