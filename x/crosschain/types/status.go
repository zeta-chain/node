package types

import (
	"fmt"
	"slices"
)

type StatusMessages struct {
	StatusMessage        string `json:"status_message"`
	ErrorMessageOutbound string `json:"error_message_outbound"`
	ErrorMessageRevert   string `json:"error_message_revert"`
	ErrorMessageAbort    string `json:"error_message_abort"`
}

func (m *Status) SetAbortRefunded() {
	m.IsAbortRefunded = true
	m.StatusMessage = "CCTX aborted and Refunded"
}

func (m *Status) UpdateStatusAndErrorMessages(newStatus CctxStatus, messages StatusMessages) {
	m.UpdateStatus(newStatus)
	m.UpdateErrorMessages(messages)
}

// UpdateStatus updates the cctx status
func (m *Status) UpdateStatus(newStatus CctxStatus) {
	if m.ValidateTransition(newStatus) {
		m.Status = newStatus
	} else {
		m.StatusMessage = fmt.Sprintf(
			"Failed to transition status from %s to %s",
			m.Status.String(),
			newStatus.String(),
		)
		m.Status = CctxStatus_Aborted
	}
}

// UpdateErrorMessages updates cctx.status.error_message and cctx.status.error_message_revert.
func (m *Status) UpdateErrorMessages(messages StatusMessages) {
	// Always update the status message, status should contain only the most recent update
	m.StatusMessage = messages.StatusMessage

	if messages.ErrorMessageOutbound != "" {
		m.ErrorMessage = messages.ErrorMessageOutbound
	}
	if messages.ErrorMessageRevert != "" {
		m.ErrorMessageRevert = messages.ErrorMessageRevert
	}
	if messages.ErrorMessageAbort != "" {
		m.ErrorMessageAbort = messages.ErrorMessageAbort
	}
}

func (m *Status) ValidateTransition(newStatus CctxStatus) bool {
	stateTransitionMap := stateTransitionMap()
	oldStatus := m.Status
	nextStatusList, isOldStatusValid := stateTransitionMap[oldStatus]
	if !isOldStatusValid {
		return false
	}
	return slices.Contains(nextStatusList, newStatus)
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

// IsTerminal returns true if the status is terminal.
// The terminal states are
// CctxStatus_Aborted
// CctxStatus_Reverted
// CctxStatus_OutboundMined
func (c CctxStatus) IsTerminal() bool {
	return c == CctxStatus_Aborted || c == CctxStatus_Reverted || c == CctxStatus_OutboundMined
}

// IsPending returns true if the status is pending.
// The pending states are
// CctxStatus_PendingInbound
// CctxStatus_PendingOutbound
// CctxStatus_PendingRevert
func (c CctxStatus) IsPending() bool {
	return !c.IsTerminal()
}
