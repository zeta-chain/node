package types_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestStatus_AbortRefunded(t *testing.T) {
	t.Run("should set status to aborted and message to CCTX aborted and Refunded", func(t *testing.T) {
		status := types.Status{
			Status:              0,
			StatusMessage:       "",
			LastUpdateTimestamp: 0,
			IsAbortRefunded:     false,
		}
		status.AbortRefunded()
		require.Equal(t, status.IsAbortRefunded, true)
		require.Equal(t, status.StatusMessage, "CCTX aborted and Refunded")
	})
}

func TestStatus_ValidateTransition(t *testing.T) {
	tests := []struct {
		name          string
		oldStatus     types.CctxStatus
		newStatus     types.CctxStatus
		expectedValid bool
	}{
		{
			"Valid - PendingInbound to PendingOutbound",
			types.CctxStatus_PendingInbound,
			types.CctxStatus_PendingOutbound,
			true,
		},
		{"Valid - PendingInbound to Aborted", types.CctxStatus_PendingInbound, types.CctxStatus_Aborted, true},
		{
			"Valid - PendingInbound to OutboundMined",
			types.CctxStatus_PendingInbound,
			types.CctxStatus_OutboundMined,
			true,
		},
		{
			"Valid - PendingInbound to PendingRevert",
			types.CctxStatus_PendingInbound,
			types.CctxStatus_PendingRevert,
			true,
		},

		{"Valid - PendingOutbound to Aborted", types.CctxStatus_PendingOutbound, types.CctxStatus_Aborted, true},
		{
			"Valid - PendingOutbound to PendingRevert",
			types.CctxStatus_PendingOutbound,
			types.CctxStatus_PendingRevert,
			true,
		},
		{
			"Valid - PendingOutbound to OutboundMined",
			types.CctxStatus_PendingOutbound,
			types.CctxStatus_OutboundMined,
			true,
		},
		{"Valid - PendingOutbound to Reverted", types.CctxStatus_PendingOutbound, types.CctxStatus_Reverted, true},

		{"Valid - PendingRevert to Aborted", types.CctxStatus_PendingRevert, types.CctxStatus_Aborted, true},
		{
			"Valid - PendingRevert to OutboundMined",
			types.CctxStatus_PendingRevert,
			types.CctxStatus_OutboundMined,
			true,
		},
		{"Valid - PendingRevert to Reverted", types.CctxStatus_PendingRevert, types.CctxStatus_Reverted, true},

		{"Invalid - PendingInbound to Reverted", types.CctxStatus_PendingInbound, types.CctxStatus_Reverted, false},
		{
			"Invalid - PendingInbound to PendingInbound",
			types.CctxStatus_PendingInbound,
			types.CctxStatus_PendingInbound,
			false,
		},

		{
			"Invalid - PendingOutbound to PendingInbound",
			types.CctxStatus_PendingOutbound,
			types.CctxStatus_PendingInbound,
			false,
		},
		{
			"Invalid - PendingOutbound to PendingOutbound",
			types.CctxStatus_PendingOutbound,
			types.CctxStatus_PendingOutbound,
			false,
		},

		{
			"Invalid - PendingRevert to PendingInbound",
			types.CctxStatus_PendingRevert,
			types.CctxStatus_PendingInbound,
			false,
		},
		{
			"Invalid - PendingRevert to PendingOutbound",
			types.CctxStatus_PendingRevert,
			types.CctxStatus_PendingOutbound,
			false,
		},
		{
			"Invalid - PendingRevert to PendingRevert",
			types.CctxStatus_PendingRevert,
			types.CctxStatus_PendingRevert,
			false,
		},

		{"Invalid old status - CctxStatus_Aborted", types.CctxStatus_Aborted, types.CctxStatus_PendingRevert, false},
		{"Invalid old status - CctxStatus_Reverted", types.CctxStatus_Reverted, types.CctxStatus_PendingRevert, false},
		{
			"Invalid old status - CctxStatus_OutboundMined",
			types.CctxStatus_OutboundMined,
			types.CctxStatus_PendingRevert,
			false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := types.Status{Status: tc.oldStatus}
			valid := m.ValidateTransition(tc.newStatus)
			if valid != tc.expectedValid {
				t.Errorf("expected %v, got %v", tc.expectedValid, valid)
			}
		})
	}
}

func TestStatus_ChangeStatus(t *testing.T) {
	t.Run("should change status and msg if transition is valid", func(t *testing.T) {
		s := types.Status{Status: types.CctxStatus_PendingInbound}

		s.ChangeStatus(types.CctxStatus_PendingOutbound, "msg")
		assert.Equal(t, s.Status, types.CctxStatus_PendingOutbound)
		assert.Equal(t, s.StatusMessage, "msg")
	})

	t.Run("should change status if transition is valid", func(t *testing.T) {
		s := types.Status{Status: types.CctxStatus_PendingInbound}

		s.ChangeStatus(types.CctxStatus_PendingOutbound, "")
		assert.Equal(t, s.Status, types.CctxStatus_PendingOutbound)
		assert.Equal(t, s.StatusMessage, "")
	})

	t.Run("should change status to aborted and msg if transition is invalid", func(t *testing.T) {
		s := types.Status{Status: types.CctxStatus_PendingOutbound}

		s.ChangeStatus(types.CctxStatus_PendingInbound, "msg")
		assert.Equal(t, s.Status, types.CctxStatus_Aborted)
		assert.Equal(
			t,
			fmt.Sprintf(
				"Failed to transition : OldStatus %s , NewStatus %s , MSG : %s :",
				types.CctxStatus_PendingOutbound.String(),
				types.CctxStatus_PendingInbound.String(),
				"msg",
			),
			s.StatusMessage,
		)
	})
}
