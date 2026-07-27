package mcrpc

import (
	"encoding/json"
	"sort"
	"sync"
	"testing"
	"time"
)

// notificationLog records which handlers fired, and the raw params of every
// notification, so a failure can show what actually arrived.
type notificationLog struct {
	mu     sync.Mutex
	fired  map[string]int
	params map[string]string
	errors map[string]error
}

func newNotificationLog() *notificationLog {
	return &notificationLog{
		fired:  map[string]int{},
		params: map[string]string{},
		errors: map[string]error{},
	}
}

func (l *notificationLog) mark(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fired[name]++
}

func (l *notificationLog) count(name string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.fired[name]
}

func (l *notificationLog) handler(t *testing.T) Handler {
	t.Helper()

	return Handler{
		OnNotification: func(method string, params json.RawMessage) {
			l.mu.Lock()
			defer l.mu.Unlock()
			if _, seen := l.params[method]; !seen {
				l.params[method] = string(params)
			}
		},
		OnError: func(method string, err error) {
			l.mu.Lock()
			defer l.mu.Unlock()
			l.errors[method] = err
		},

		OnAllowlistAdded:   func(Player) { l.mark("allowlist/added") },
		OnAllowlistRemoved: func(Player) { l.mark("allowlist/removed") },
		OnOperatorAdded:    func(Operator) { l.mark("operators/added") },
		OnOperatorRemoved:  func(Operator) { l.mark("operators/removed") },
		OnBanAdded:         func(UserBan) { l.mark("bans/added") },
		OnBanRemoved:       func(Player) { l.mark("bans/removed") },
		OnIPBanAdded:       func(IPBan) { l.mark("ip_bans/added") },
		OnIPBanRemoved:     func(string) { l.mark("ip_bans/removed") },
		OnGameruleUpdated:  func(TypedGameRule) { l.mark("gamerules/updated") },
	}
}

// report renders everything observed, for use in a failure message.
func (l *notificationLog) report() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	methods := make([]string, 0, len(l.params))
	for method := range l.params {
		methods = append(methods, method)
	}
	sort.Strings(methods)

	report := "\nobserved notifications:"
	for _, method := range methods {
		report += "\n  " + method + "\n      " + l.params[method]
		if err, failed := l.errors[method]; failed {
			report += "\n      DECODE FAILURE: " + err.Error()
		}
	}
	return report
}

// TestWireNotificationsReachHandlers checks that every notification the server
// pushes is decoded and delivered.
//
// This is the regression test for a real defect. JSON-RPC params are a
// positional argument list and every notification declares exactly one
// parameter, so adding one player to the allowlist sends [{"name":…}] — the
// payload wrapped in an argument list. The client decoded the params directly
// as the payload, so every handler here silently never fired.
func TestWireNotificationsReachHandlers(t *testing.T) {
	host, port, _, _ := getTestConfig()

	log := newNotificationLog()
	options := append(testTLSOptions(t), WithHandler(log.handler(t)))
	options = append(options, testTraceOption(t)...)

	client := NewHostPort(host, port, testSecret(), options...)

	ctx := t.Context()
	if err := client.Start(ctx); err != nil {
		t.Skipf("Skipping test: cannot connect to server at %s:%d: %v", host, port, err)
	}
	t.Cleanup(func() { _ = client.Close() })

	player := Player{Name: "fi_xz", UUID: "a0d8c884-2a79-4c95-8617-a51d27a427ec"}

	// Snapshot everything this test disturbs, and put it back afterwards.
	allowlist, err := client.GetAllowlist(ctx)
	if err != nil {
		t.Fatalf("GetAllowlist failed: %v", err)
	}
	operators, err := client.GetOperators(ctx)
	if err != nil {
		t.Fatalf("GetOperators failed: %v", err)
	}
	bans, err := client.GetBanlist(ctx)
	if err != nil {
		t.Fatalf("GetBanlist failed: %v", err)
	}
	ipBans, err := client.GetIPBanlist(ctx)
	if err != nil {
		t.Fatalf("GetIPBanlist failed: %v", err)
	}

	t.Cleanup(func() {
		_, _ = client.SetAllowlist(ctx, allowlist)
		_, _ = client.SetOperators(ctx, operators)
		_, _ = client.SetBanlist(ctx, bans)
		_, _ = client.SetIPBanlist(ctx, ipBans)
	})

	gamerules, err := client.GetGameRules(ctx)
	if err != nil {
		t.Fatalf("GetGameRules failed: %v", err)
	}
	var rule TypedGameRule
	for _, candidate := range gamerules {
		if candidate.Type == "boolean" {
			rule = candidate
			break
		}
	}
	if rule.Key == "" {
		t.Fatal("the server reports no boolean game rule to exercise")
	}
	ruleValue, _ := rule.Bool()
	t.Cleanup(func() {
		_, _ = client.UpdateGameRule(ctx, rule.WithBool(ruleValue))
	})

	steps := []struct {
		notification string
		mutate       func() error
	}{
		{"allowlist/added", func() error { _, err := client.AddAllowlist(ctx, player); return err }},
		{"allowlist/removed", func() error { _, err := client.RemoveAllowlist(ctx, player); return err }},
		{"operators/added", func() error {
			_, err := client.AddOperators(ctx, Operator{Player: player, PermissionLevel: 4})
			return err
		}},
		{"operators/removed", func() error { _, err := client.RemoveOperators(ctx, player); return err }},
		{"bans/added", func() error {
			_, err := client.AddBanlist(ctx, UserBan{Player: player, Reason: "notification check", Source: "mcrpc test"})
			return err
		}},
		{"bans/removed", func() error { _, err := client.RemoveBanlist(ctx, player); return err }},
		{"ip_bans/added", func() error {
			_, err := client.AddIPBanlist(ctx, IncomingIPBan{
				IPBan:  IPBan{IP: "192.0.2.1", Reason: "notification check", Source: "mcrpc test"},
				Player: player,
			})
			return err
		}},
		{"ip_bans/removed", func() error { _, err := client.RemoveIPBanlist(ctx, "192.0.2.1"); return err }},
		{"gamerules/updated", func() error {
			_, err := client.UpdateGameRule(ctx, rule.WithBool(!ruleValue))
			return err
		}},
	}

	for _, step := range steps {
		before := log.count(step.notification)

		if err := step.mutate(); err != nil {
			t.Errorf("%s: the triggering call failed: %v", step.notification, err)
			continue
		}

		// Notifications are pushed asynchronously, after the call returns.
		deadline := time.Now().Add(3 * time.Second)
		for log.count(step.notification) == before && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
		}

		if log.count(step.notification) == before {
			t.Errorf("%s was never delivered to its handler%s", step.notification, log.report())
		}
	}
}
