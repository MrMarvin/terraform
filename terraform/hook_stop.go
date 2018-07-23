package terraform

import (
	"sync/atomic"

	"github.com/zclconf/go-cty/cty"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states"
)

// stopHook is a private Hook implementation that Terraform uses to
// signal when to stop or cancel actions.
type stopHook struct {
	stop uint32
}

var _ Hook = (*stopHook)(nil)

func (h *stopHook) PreApply(addr addrs.ResourceInstance, gen states.Generation, priorState, plannedNewState cty.Value) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PostApply(addr addrs.ResourceInstance, gen states.Generation, newState cty.Value, err error) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PreDiff(addr addrs.ResourceInstance, priorState, proposedNewState cty.Value) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PostDiff(addr addrs.ResourceInstance, priorState, plannedNewState cty.Value) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PreProvisionInstance(addr addrs.ResourceInstance, state cty.Value) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PostProvisionInstance(addr addrs.ResourceInstance, state cty.Value) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PreProvisionInstanceStep(addr addrs.ResourceInstance, typeName string) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PostProvisionInstanceStep(addr addrs.ResourceInstance, typeName string, err error) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) ProvisionOutput(addr addrs.ResourceInstance, typeName string, line string) {
}

func (h *stopHook) PreRefresh(addr addrs.ResourceInstance, priorState cty.Value) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PostRefresh(addr addrs.ResourceInstance, priorState cty.Value, newState cty.Value) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PreImportState(addr addrs.ResourceInstance, importID string) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PostImportState(addr addrs.ResourceInstance, imported []*states.ImportedObject) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) PostStateUpdate(new *states.State) (HookAction, error) {
	return h.hook()
}

func (h *stopHook) hook() (HookAction, error) {
	if h.Stopped() {
		return HookActionHalt, nil
	}

	return HookActionContinue, nil
}

// reset should be called within the lock context
func (h *stopHook) Reset() {
	atomic.StoreUint32(&h.stop, 0)
}

func (h *stopHook) Stop() {
	atomic.StoreUint32(&h.stop, 1)
}

func (h *stopHook) Stopped() bool {
	return atomic.LoadUint32(&h.stop) == 1
}
