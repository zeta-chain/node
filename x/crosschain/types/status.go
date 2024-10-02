package types

import (
	"fmt"
)

func (m *Status) AbortRefunded() {
	m.IsAbortRefunded = true
	m.StatusMessage = "CCTX aborted and Refunded"
}

// UpdateCctxStatus transitions the Status.
// In case of an error, ErrorMessage is updated.
// In case of no error, StatusMessage is updated.
func (m *Status) UpdateCctxStatus(newStatus CctxStatus, isError bool, statusMsg, errorMsg string) {
	m.ChangeStatus(newStatus, statusMsg)

	if isError && errorMsg != "" {
		m.ErrorMessage = errorMsg
	} else if isError && errorMsg == "" {
		m.ErrorMessage = "unknown error"
	}
}

// ChangeStatus changes the status of the cross chain transaction.
func (m *Status) ChangeStatus(newStatus CctxStatus, statusMsg string) {
	if !m.ValidateTransition(newStatus) {
		m.StatusMessage = fmt.Sprintf(
			"Failed to transition status from %s to %s",
			m.Status.String(),
			newStatus.String(),
		)

		m.Status = CctxStatus_Aborted
		return
	}

	if statusMsg == "" {
		m.StatusMessage = fmt.Sprintf("Status changed from %s to %s", m.Status.String(), newStatus.String())
	} else {
		m.StatusMessage = statusMsg
	}

	m.Status = newStatus
}

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
