package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func AllStatus() []CctxStatus {
	return []CctxStatus{
		CctxStatus_PendingInbound,
		CctxStatus_PendingOutbound,
		CctxStatus_OutboundMined,
		CctxStatus_PendingRevert,
		CctxStatus_Reverted,
		CctxStatus_Aborted,
	}
}

func (m *Status) ChangeStatus(ctx *sdk.Context, newStatus CctxStatus, msg, logIdentifier string) {
	oldStatus := m.Status
	m.StatusMessage = msg
	if !m.ValidateTransition(newStatus) {
		m.StatusMessage = fmt.Sprintf("Failed to transition : OldStatus %s , NewStatus %s , MSG : %s :", m.Status.String(), newStatus.String(), msg)
		m.Status = CctxStatus_Aborted
		return
	}
	m.Status = newStatus
	EmitStatusChangeEvent(ctx, oldStatus.String(), newStatus.String(), logIdentifier)
} //nolint:typecheck

func (m *Status) ValidateTransition(newStatus CctxStatus) bool {
	stateTransitionMap := stateTransitionMap()
	oldStatus := m.Status
	nextStatusList, isOldStatusValid := stateTransitionMap[oldStatus]
	if !isOldStatusValid {
		return false
	}
	for _, status := range nextStatusList {
		if status == newStatus {
			return true
		}
	}
	return false
}

func stateTransitionMap() map[CctxStatus][]CctxStatus {
	stateTransitionMap := make(map[CctxStatus][]CctxStatus)
	stateTransitionMap[CctxStatus_PendingInbound] = []CctxStatus{
		CctxStatus_PendingOutbound,
		CctxStatus_Aborted,
		CctxStatus_OutboundMined, // EVM Deposit
		CctxStatus_PendingRevert, // EVM Deposit contract call reverted; should refund
	}
	stateTransitionMap[CctxStatus_PendingOutbound] = []CctxStatus{
		CctxStatus_Aborted,
		CctxStatus_PendingRevert,
		CctxStatus_OutboundMined,
		CctxStatus_Reverted,
	}

	stateTransitionMap[CctxStatus_PendingRevert] = []CctxStatus{
		CctxStatus_Aborted,
		CctxStatus_OutboundMined,
		CctxStatus_Reverted,
	}
	return stateTransitionMap

}

func EmitStatusChangeEvent(ctx *sdk.Context, oldStatus, newStatus, logIdentifier string) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(StatusChanged,
			sdk.NewAttribute(OldStatus, oldStatus),
			sdk.NewAttribute(NewStatus, newStatus),
			sdk.NewAttribute(Identifiers, logIdentifier),
		),
	)
}
