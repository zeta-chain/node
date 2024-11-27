package types

import (
	"fmt"
)

func (m *Status) AbortRefunded() {
	m.IsAbortRefunded = true
	m.StatusMessage = "CCTX aborted and Refunded"
}

// UpdateStatusAndErrorMessages transitions the Status and Error messages.
func (m *Status) UpdateStatusAndErrorMessages(newStatus CctxStatus, statusMsg, errorMsg string) {
	m.UpdateStatus(newStatus, statusMsg)

	if errorMsg != "" &&
		(newStatus == CctxStatus_Aborted || newStatus == CctxStatus_Reverted || newStatus == CctxStatus_PendingRevert) {
		m.UpdateErrorMessage(errorMsg)
	}
}

// UpdateStatus updates the cctx status and cctx.status.status_message.
func (m *Status) UpdateStatus(newStatus CctxStatus, statusMsg string) {
	if m.ValidateTransition(newStatus) {
		m.StatusMessage = fmt.Sprintf("Status changed from %s to %s", m.Status.String(), newStatus.String())
		m.Status = newStatus
	} else {
		m.StatusMessage = fmt.Sprintf(
			"Failed to transition status from %s to %s",
			m.Status.String(),
			newStatus.String(),
		)

		m.Status = CctxStatus_Aborted
	}

	if statusMsg != "" {
		m.StatusMessage += fmt.Sprintf(": %s", statusMsg)
	}
}

// UpdateErrorMessage updates cctx.status.error_message.
func (m *Status) UpdateErrorMessage(errorMsg string) {
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
