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

// Bool returns the game rule value as a boolean. The second return value is
// false if the rule does not hold a boolean.
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
// Values decoded from JSON arrive as float64, and a value may also be
// transmitted as a string; both are accepted.
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
