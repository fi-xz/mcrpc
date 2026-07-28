package mcrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

// describeParams renders a notification payload for the transcript, keeping
// "the field was absent" distinct from "the field was there and empty". A
// notification that declares no parameters could legitimately arrive either
// way, and a probe exists to find out which.
func describeParams(params json.RawMessage) string {
	if params == nil {
		return "(absent)"
	}
	if len(params) == 0 {
		return "(empty)"
	}
	return string(params)
}

// worldUpgradeLog records a timestamped transcript of selected notifications
// during a boot, plus the first raw params observed for each method.
type worldUpgradeLog struct {
	mu        sync.Mutex
	started   time.Time
	events    []string
	raw       map[string]string
	progress  []float64
	terminal  chan struct{}
	terminate sync.Once
}

func newWorldUpgradeLog() *worldUpgradeLog {
	return &worldUpgradeLog{
		started:  time.Now(),
		raw:      map[string]string{},
		terminal: make(chan struct{}),
	}
}

func (l *worldUpgradeLog) note(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events,
		time.Since(l.started).Truncate(time.Millisecond).String()+"  "+fmt.Sprintf(format, args...))
}

// done marks the upgrade as concluded, releasing the probe early.
func (l *worldUpgradeLog) done() {
	l.terminate.Do(func() { close(l.terminal) })
}

func (l *worldUpgradeLog) handler() Handler {
	return Handler{
		OnNotification: func(method string, params json.RawMessage) {
			l.mu.Lock()
			defer l.mu.Unlock()

			if _, seen := l.raw[method]; !seen {
				l.raw[method] = describeParams(params)
			}
		},
		OnError: func(method string, err error) {
			l.note("DECODE FAILURE  %s: %v", method, err)
		},

		OnWorldUpgradeStarted: func() {
			l.note("world/upgrade_started")
		},
		OnWorldUpgradeProgress: func(progress float64) {
			l.mu.Lock()
			l.progress = append(l.progress, progress)
			l.mu.Unlock()
			l.note("world/upgrade_progress  %.4f", progress)
		},
		OnWorldUpgradeFinished: func() {
			l.note("world/upgrade_finished")
			l.done()
		},
		OnWorldUpgradeFailed: func(reason string) {
			l.note("world/upgrade_failed  %q", reason)
			l.done()
		},

		// Boot-time context, so the transcript shows where the upgrade sits
		// relative to the server coming up.
		OnServerStarted:  func() { l.note("server/started") },
		OnServerStopping: func() { l.note("server/stopping") },
		OnServerStatus:   func(state ServerState) { l.note("server/status  started=%v", state.Started) },
	}
}

// TestProbeWorldUpgrade records the world upgrade notifications a server emits
// while it boots. They arrive only from management API 3.1.0 (Minecraft 26.3)
// and only when the world actually needs converting, so nothing else can
// confirm them.
//
// Start this first and boot the server afterwards. From API 3.0.0 the
// management server accepts connections before the Minecraft server spins up,
// so the probe simply waits for the port rather than racing the boot.
//
//	PROBE_WINDOW=900 TEST_HOST=127.0.0.1 TEST_PORT=8082 USE_TLS=true \
//	  TEST_TLS_SERVER_NAME=… go test -v -timeout 20m -run TestProbeWorldUpgrade ./...
func TestProbeWorldUpgrade(t *testing.T) {
	seconds, err := strconv.Atoi(os.Getenv("PROBE_WINDOW"))
	if err != nil || seconds <= 0 {
		t.Skip("set PROBE_WINDOW to the number of seconds to watch a server boot")
	}
	window := time.Duration(seconds) * time.Second

	// The upgrade can finish before the management server binds, so the gap
	// between attempts is the probe's resolution. PROBE_RETRY_MS tightens it.
	retry := 50 * time.Millisecond
	if raw := os.Getenv("PROBE_RETRY_MS"); raw != "" {
		if ms, err := strconv.Atoi(raw); err == nil && ms > 0 {
			retry = time.Duration(ms) * time.Millisecond
		}
	}

	host, port, _, _ := getTestConfig()

	log := newWorldUpgradeLog()
	options := append(testTLSOptions(t), WithHandler(log.handler()))
	options = append(options, testTraceOption(t)...)
	// A boot-time handshake that stalls should be abandoned and retried rather
	// than sitting on the default ten seconds.
	options = append(options, WithHandshakeTimeout(2*time.Second))

	client := NewHostPort(host, port, testSecret(), options...)
	t.Cleanup(func() { _ = client.Close() })

	ctx, stop := context.WithTimeout(context.Background(), window)
	defer stop()

	t.Logf("waiting up to %s for %s:%d, retrying every %s — start the server now",
		window, host, port, retry)

	// Keep a session up for the whole window. The management server may accept
	// before the world work begins, and may drop the connection as the server
	// finishes coming up, so reconnect rather than giving up.
	var (
		connections int
		lastDialErr error
	)
	for ctx.Err() == nil {
		if err := client.Start(ctx); err != nil {
			if errors.Is(err, ErrAlreadyStarted) {
				// Still connected; wait for the session to end.
				select {
				case <-client.DisconnectNotify():
				case <-log.terminal:
				case <-ctx.Done():
				}
				continue
			}

			lastDialErr = err
			select {
			case <-time.After(retry):
			case <-ctx.Done():
			}
			continue
		}

		connections++
		log.note("connected (attempt %d)", connections)

		// World upgrade notifications arrive from management API 3.1.0 on, so
		// record what this server reports rather than guessing at the run.
		if version, err := client.APIVersion(ctx); err == nil {
			log.note("management API %s (world upgrades need 3.1.0)", version)
		}

		select {
		case <-client.DisconnectNotify():
			log.note("disconnected")
		case <-log.terminal:
		case <-ctx.Done():
		}

		select {
		case <-log.terminal:
			stop()
		default:
		}
	}

	log.mu.Lock()
	defer log.mu.Unlock()

	if connections == 0 {
		t.Skipf("never connected to %s:%d within %s (last error: %v)", host, port, window, lastDialErr)
	}

	t.Logf("transcript (%d connection(s)):", connections)
	for _, event := range log.events {
		t.Logf("  %s", event)
	}

	t.Log("raw notification payloads:")
	for method, params := range log.raw {
		t.Logf("  %s\n      %s", method, params)
	}

	if len(log.progress) > 0 {
		t.Logf("progress values: %v", log.progress)
		for _, value := range log.progress {
			if value < 0 || value > 1 {
				t.Errorf("progress %v is outside the documented 0..1 range", value)
			}
		}
	}
}
