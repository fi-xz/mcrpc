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

// worldUpgradeLog records everything the server pushes during a boot, with the
// time since the probe started, so the ordering and cadence can be read back.
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
			if _, seen := l.raw[method]; !seen {
				l.raw[method] = string(params)
			}
			l.mu.Unlock()
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

	host, port, _, _ := getTestConfig()

	log := newWorldUpgradeLog()
	options := append(testTLSOptions(t), WithHandler(log.handler()))
	options = append(options, testTraceOption(t)...)

	client := NewHostPort(host, port, testSecret(), options...)
	t.Cleanup(func() { _ = client.Close() })

	ctx, stop := context.WithTimeout(context.Background(), window)
	defer stop()

	t.Logf("waiting up to %s for %s:%d — start the server now", window, host, port)

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
			case <-time.After(250 * time.Millisecond):
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
