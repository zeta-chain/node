package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatus_ChangeStatus(t *testing.T) {
	tt := []struct {
		Name         string
		Status       Status
		NonErrStatus CctxStatus
		Msg          string
		IsErr        bool
		ErrStatus    CctxStatus
	}{
		{
			Name: "Transition on finalize Inbound",
			Status: Status{
				Status:              CctxStatus_PendingInbound,
				StatusMessage:       "Getting InTX Votes",
				LastUpdateTimestamp: 0,
			},
			Msg:          "Got super majority and finalized Inbound",
			NonErrStatus: CctxStatus_PendingOutbound,
			ErrStatus:    CctxStatus_Aborted,
			IsErr:        false,
		},
		{
			Name: "Transition on finalize Inbound Fail",
			Status: Status{
				Status:              CctxStatus_PendingInbound,
				StatusMessage:       "Getting InTX Votes",
				LastUpdateTimestamp: 0,
			},
			Msg:          "Got super majority and finalized Inbound",
			NonErrStatus: CctxStatus_OutboundMined,
			ErrStatus:    CctxStatus_Aborted,
			IsErr:        true,
		},
	}
	for _, test := range tt {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			test.Status.ChangeStatus(test.NonErrStatus, test.Msg)
			if test.IsErr {
				assert.Equal(t, test.ErrStatus, test.Status.Status)
			} else {
				assert.Equal(t, test.NonErrStatus, test.Status.Status)
			}
		})
	}
}
