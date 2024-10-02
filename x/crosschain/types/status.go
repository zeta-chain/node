package types

import (
	"fmt"
)

func (m *Status) AbortRefunded() {
	m.IsAbortRefunded = true
	m.StatusMessage = "CCTX aborted and Refunded"
}

// UpdateCctxMessages transitions the Status and Error messages.
func (m *Status) UpdateCctxMessages(newStatus CctxStatus, isError bool, statusMsg, errorMsg string) {
	m.UpdateStatusMessage(newStatus, statusMsg)
	m.UpdateErrorMessage(isError, errorMsg)
}

// UpdateStatusMessage updates the cctx status and cctx.status.status_message.
func (m *Status) UpdateStatusMessage(newStatus CctxStatus, statusMsg string) {
	if !m.ValidateTransition(newStatus) {
		m.StatusMessage = fmt.Sprintf(
			"Failed to transition status from %s to %s",
			m.Status.String(),
			newStatus.String(),
		)

		m.Status = CctxStatus_Aborted
		return
	}

	m.StatusMessage = fmt.Sprintf("Status changed from %s to %s", m.Status.String(), newStatus.String())

	if statusMsg != "" {
		m.StatusMessage += fmt.Sprintf(": %s", statusMsg)
	}

	m.Status = newStatus
}

// UpdateErrorMessage updates cctx.status.error_message.
func (m *Status) UpdateErrorMessage(isError bool, errorMsg string) {
	if !isError {
		return
	}

	errMsg := errorMsg
	if errMsg == "" {
		errMsg = "unknown error"
	}

	m.ErrorMessage = errMsg
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
