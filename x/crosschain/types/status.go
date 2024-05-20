package types

import (
	"fmt"
)

func (m *Status) AbortRefunded(timeStamp int64) {
	m.IsAbortRefunded = true
	m.StatusMessage = "CCTX aborted and Refunded"
	m.LastUpdateTimestamp = timeStamp
}

// ChangeStatus changes the status of the cross chain transaction
// empty msg does not overwrite old status message
func (m *Status) ChangeStatus(newStatus CctxStatus, msg string) {
	if len(msg) > 0 {
		if m.StatusMessage != "" {
			m.StatusMessage = fmt.Sprintf("%s : %s", m.StatusMessage, msg)
		} else {
			m.StatusMessage = msg
		}
	}
	if !m.ValidateTransition(newStatus) {
		m.StatusMessage = fmt.Sprintf("Failed to transition : OldStatus %s , NewStatus %s , MSG : %s :", m.Status.String(), newStatus.String(), msg)
		m.Status = CctxStatus_Aborted
		return
	}
	m.Status = newStatus

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
