package terraform

import (
	"testing"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states"
)

func TestNilHook_impl(t *testing.T) {
	var _ Hook = new(NilHook)
}

// testHook is a Hook implementation that logs the calls it receives.
// It is intended for testing that core code is emitting the correct hooks
// for a given situation.
type testHook struct {
	Calls []*testHookCall
}

var _ Hook = (*testHook)(nil)

// testHookCall represents a single call in testHook.
// This hook just logs string names to make it easy to write "want" expressions
// in tests that can DeepEqual against the real calls.
type testHookCall struct {
	Action     string
	InstanceID string
}

func (h *testHook) PreApply(addr addrs.ResourceInstance, gen states.Generation, priorState, plannedNewState cty.Value) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PreApply", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PostApply(addr addrs.ResourceInstance, gen states.Generation, newState cty.Value, err error) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PostApply", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PreDiff(addr addrs.ResourceInstance, priorState, proposedNewState cty.Value) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PreDiff", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PostDiff(addr addrs.ResourceInstance, priorState, plannedNewState cty.Value) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PostDiff", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PreProvisionInstance(addr addrs.ResourceInstance, state cty.Value) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PreProvisionInstance", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PostProvisionInstance(addr addrs.ResourceInstance, state cty.Value) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PostProvisionInstance", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PreProvisionInstanceStep(addr addrs.ResourceInstance, typeName string) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PreProvisionInstanceStep", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PostProvisionInstanceStep(addr addrs.ResourceInstance, typeName string, err error) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PostProvisionInstanceStep", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) ProvisionOutput(addr addrs.ResourceInstance, typeName string, line string) {
	h.Calls = append(h.Calls, &testHookCall{"ProvisionOutput", i.ResourceAddress().String()})
}

func (h *testHook) PreRefresh(addr addrs.ResourceInstance, priorState cty.Value) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PreRefresh", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PostRefresh(addr addrs.ResourceInstance, priorState cty.Value, newState cty.Value) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PostRefresh", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PreImportState(addr addrs.ResourceInstance, importID string) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PreImportState", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PostImportState(addr addrs.ResourceInstance, imported []*states.ImportedObject) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PostImportState", i.ResourceAddress().String()})
	return HookActionContinue, nil
}

func (h *testHook) PostStateUpdate(new *states.State) (HookAction, error) {
	h.Calls = append(h.Calls, &testHookCall{"PostStateUpdate", ""})
	return HookActionContinue, nil
}
