package mcrpc

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/fi-xz/mcrpc/internal/protocol"
)

// openRPCSchema is the part of the rpc.discover document this package relies
// on. The server advertises its own API, which makes protocol drift something
// the test suite can detect rather than something a user reports.
type openRPCSchema struct {
	Methods []struct {
		Name   string `json:"name"`
		Params []struct {
			Name string `json:"name"`
		} `json:"params"`
	} `json:"methods"`
}

// discoverSchema asks the server for its API description.
func discoverSchema(t *testing.T) (openRPCSchema, *Client) {
	t.Helper()

	client, ctx := createTestClient(t)

	var raw json.RawMessage
	if err := client.call(ctx, protocol.MethodDiscover, nil, &raw); err != nil {
		t.Fatalf("rpc.discover failed: %v", err)
	}

	var schema openRPCSchema
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("could not decode the schema: %v", err)
	}
	if len(schema.Methods) == 0 {
		t.Fatal("the server advertised no methods")
	}

	if out := os.Getenv("SCHEMA_OUT"); out != "" {
		indented, err := json.MarshalIndent(raw, "", "  ")
		if err == nil && os.WriteFile(out, indented, 0o644) == nil {
			t.Logf("wrote %s", out)
		}
	}

	return schema, client
}

// TestWireNotificationsTakeOneParameter guards the assumption every
// notification handler is built on: params are a positional argument list
// carrying exactly one value, so the payload is element 0.
//
// A future notification with two parameters would silently deliver the wrong
// payload, so it should fail here first.
func TestWireNotificationsTakeOneParameter(t *testing.T) {
	schema, _ := discoverSchema(t)

	checked := 0
	for _, method := range schema.Methods {
		if !strings.Contains(method.Name, "notification/") {
			continue
		}
		checked++
		if len(method.Params) > 1 {
			names := make([]string, 0, len(method.Params))
			for _, param := range method.Params {
				names = append(names, param.Name)
			}
			t.Errorf("%s declares %d parameters (%s); decodeParam only reads the first",
				method.Name, len(method.Params), strings.Join(names, ", "))
		}
	}

	if checked == 0 {
		t.Fatal("the schema advertised no notifications")
	}
	t.Logf("checked %d notifications", checked)
}

// TestWireRequestParameterNames compares the JSON field names this package
// sends against the ones the server advertises.
//
// SetAllowFlightParams shipped as "allowed" — the schema's name for the
// *result* — while the parameter is "allow", so every SetAllowFlight call
// failed with Invalid params. This is the check that catches that class of bug.
func TestWireRequestParameterNames(t *testing.T) {
	schema, _ := discoverSchema(t)

	advertised := map[string][]string{}
	for _, method := range schema.Methods {
		names := make([]string, 0, len(method.Params))
		for _, param := range method.Params {
			names = append(names, param.Name)
		}
		advertised[method.Name] = names
	}

	cases := []struct {
		method string
		params any
	}{
		{protocol.MethodAllowlistAdd, protocol.AddAllowlistParams{}},
		{protocol.MethodAllowlistRemove, protocol.RemoveAllowlistParams{}},
		{protocol.MethodAllowlistSet, protocol.SetAllowlistParams{}},
		{protocol.MethodBansAdd, protocol.AddBanlistParams{}},
		{protocol.MethodBansRemove, protocol.RemoveBanlistParams{}},
		{protocol.MethodBansSet, protocol.SetBanlistParams{}},
		{protocol.MethodGameRulesUpdate, protocol.UpdateGameRulesParams{}},
		{protocol.MethodIPBansAdd, protocol.AddIPBanlistParams{}},
		{protocol.MethodIPBansRemove, protocol.RemoveIPBanlistParams{}},
		{protocol.MethodIPBansSet, protocol.SetIPBanlistParams{}},
		{protocol.MethodOperatorsAdd, protocol.AddOperatorParams{}},
		{protocol.MethodOperatorsRemove, protocol.RemoveOperatorParams{}},
		{protocol.MethodOperatorsSet, protocol.SetOperatorParams{}},
		{protocol.MethodPlayersKick, protocol.KickPlayerParams{}},
		{protocol.MethodServerSave, protocol.ServerSaveParams{}},
		{protocol.MethodServerSettingsAcceptTransfersSet, protocol.SetAcceptTransfersParams{}},
		{protocol.MethodServerSettingsAllowFlightSet, protocol.SetAllowFlightParams{}},
		{protocol.MethodServerSettingsAutoSaveSet, protocol.SetAutosaveParams{}},
		{protocol.MethodServerSettingsDifficultySet, protocol.SetDifficultyParams{}},
		{protocol.MethodServerSettingsEnforceAllowlistSet, protocol.SetEnforceAllowlistParams{}},
		{protocol.MethodServerSettingsEntityBroadcastRangeSet, protocol.SetEntityBroadcastRangeParams{}},
		{protocol.MethodServerSettingsForceGameModeSet, protocol.SetForceGameModeParams{}},
		{protocol.MethodServerSettingsGameModeSet, protocol.SetGameModeParams{}},
		{protocol.MethodServerSettingsHeartbeatIntervalSet, protocol.SetStatusHeartbeatIntervalParams{}},
		{protocol.MethodServerSettingsHideOnlinePlayersSet, protocol.SetHideOnlinePlayersParams{}},
		{protocol.MethodServerSettingsMOTDSet, protocol.SetMOTDParams{}},
		{protocol.MethodServerSettingsMaxPlayersSet, protocol.SetMaxPlayersParams{}},
		{protocol.MethodServerSettingsOperatorPermissionLevelSet, protocol.SetOperatorUserPermissionLevelParams{}},
		{protocol.MethodServerSettingsPauseWhenEmptySet, protocol.SetPauseWhenEmptySecondsParams{}},
		{protocol.MethodServerSettingsPlayerIdleTimeoutSet, protocol.SetPlayerIdleTimeoutParams{}},
		{protocol.MethodServerSettingsSimulationDistanceSet, protocol.SetSimulationDistanceParams{}},
		{protocol.MethodServerSettingsSpawnProtectionSet, protocol.SetSpawnProtectionRadiusParams{}},
		{protocol.MethodServerSettingsStatusRepliesSet, protocol.SetStatusRepliesParams{}},
		{protocol.MethodServerSettingsUseAllowlistSet, protocol.SetUseAllowlistParams{}},
		{protocol.MethodServerSettingsViewDistanceSet, protocol.SetViewDistanceParams{}},
		{protocol.MethodServerSystemMessage, protocol.SystemMessageParams{}},
	}

	for _, test := range cases {
		want, known := advertised[test.method]
		if !known {
			t.Errorf("%s is not advertised by this server", test.method)
			continue
		}

		got := jsonFieldNames(t, test.params)
		sort.Strings(got)

		sorted := append([]string(nil), want...)
		sort.Strings(sorted)

		if strings.Join(got, ",") != strings.Join(sorted, ",") {
			t.Errorf("%s sends %v, schema declares %v", test.method, got, want)
		}
	}
}

// jsonFieldNames returns the top-level JSON keys a params struct serialises to.
func jsonFieldNames(t *testing.T, params any) []string {
	t.Helper()

	encoded, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("could not marshal %T: %v", params, err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil {
		t.Fatalf("could not inspect %T: %v", params, err)
	}

	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	return names
}

// TestWireAPIVersion records the management API version this server implements.
// It is the version protocol behaviour follows, not the Minecraft version.
func TestWireAPIVersion(t *testing.T) {
	client, ctx := createTestClient(t)

	version, err := client.APIVersion(ctx)
	if err != nil {
		t.Fatalf("APIVersion failed: %v", err)
	}
	if version == "" {
		t.Error("the server advertised no API version")
	}

	status, err := client.GetServerStatus(ctx)
	if err != nil {
		t.Fatalf("GetServerStatus failed: %v", err)
	}

	t.Logf("Minecraft %s (protocol %d) implements management API %s",
		status.Version.Name, status.Version.Protocol, version)
}
