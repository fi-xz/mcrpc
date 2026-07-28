package mcrpc

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestProbePlayerNotifications waits for a real player to join and leave, and
// reports the raw params alongside whether the typed handlers fired. Nothing
// else can confirm these two notifications: they need an actual client, so this
// only runs when PROBE_WINDOW gives it a number of seconds to listen for.
//
//	PROBE_WINDOW=120 go test -v -timeout 300s -run TestProbePlayerNotifications ./...
func TestProbePlayerNotifications(t *testing.T) {
	seconds, err := strconv.Atoi(os.Getenv("PROBE_WINDOW"))
	if err != nil || seconds <= 0 {
		t.Skip("set PROBE_WINDOW to the number of seconds to wait for a player to join")
	}
	window := time.Duration(seconds) * time.Second

	host, port, _, _ := getTestConfig()

	var mu sync.Mutex
	rawParams := map[string]string{}
	typed := map[string]Player{}
	var decodeErrors []string

	record := func(method string, params json.RawMessage) {
		mu.Lock()
		defer mu.Unlock()
		if _, seen := rawParams[method]; !seen {
			rawParams[method] = describeParams(params)
		}
	}

	options := append(testTLSOptions(t), WithHandler(Handler{
		OnNotification: record,
		OnError: func(method string, err error) {
			mu.Lock()
			defer mu.Unlock()
			decodeErrors = append(decodeErrors, method+": "+err.Error())
		},
		OnPlayerJoined: func(player Player) {
			mu.Lock()
			defer mu.Unlock()
			typed["joined"] = player
		},
		OnPlayerLeft: func(player Player) {
			mu.Lock()
			defer mu.Unlock()
			typed["left"] = player
		},
		OnServerStatus: func(state ServerState) {
			mu.Lock()
			defer mu.Unlock()
			typed["status"] = Player{Name: state.Version.Name}
		},
	}))

	client := NewHostPort(host, port, testSecret(), options...)

	ctx := t.Context()
	if err := client.Start(ctx); err != nil {
		t.Skipf("Skipping probe: cannot connect to %s:%d: %v", host, port, err)
	}
	t.Cleanup(func() { _ = client.Close() })

	t.Logf("listening on %s:%d for %s — join the Minecraft server now", host, port, window)

	deadline := time.Now().Add(window)
	for time.Now().Before(deadline) {
		mu.Lock()
		_, joined := typed["joined"]
		_, left := typed["left"]
		mu.Unlock()
		if joined && left {
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(rawParams) == 0 {
		t.Log("RESULT: no notifications arrived at all")
		return
	}

	for method, params := range rawParams {
		t.Logf("RAW  %s\n      %s", method, params)
	}
	for event, player := range typed {
		t.Logf("TYPED  %s -> %+v", event, player)
	}
	for _, failure := range decodeErrors {
		t.Logf("DECODE FAILURE  %s", failure)
	}
}
