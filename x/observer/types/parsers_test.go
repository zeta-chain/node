package types_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestConvertReceiveStatusToVoteType(t *testing.T) {
	tests := []struct {
		name     string
		status   chains.ReceiveStatus
		expected types.VoteType
	}{
		{"TestSuccessStatus", chains.ReceiveStatus_success, types.VoteType_SuccessObservation},
		{"TestFailedStatus", chains.ReceiveStatus_failed, types.VoteType_FailureObservation},
		{"TestDefaultStatus", chains.ReceiveStatus_created, types.VoteType_NotYetVoted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.ConvertReceiveStatusToVoteType(tt.status)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestParseStringToObservationType(t *testing.T) {
	tests := []struct {
		name            string
		observationType string
		expected        types.ObservationType
	}{
		{"TestValidObservationType1", "EmptyObserverType", types.ObservationType(0)},
		{"TestValidObservationType1", "InboundTx", types.ObservationType(1)},
		{"TestValidObservationType1", "OutboundTx", types.ObservationType(2)},
		{"TestValidObservationType1", "TSSKeyGen", types.ObservationType(3)},
		{"TestValidObservationType1", "TSSKeySign", types.ObservationType(4)},
		{"TestInvalidObservationType", "InvalidType", types.ObservationType(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := types.ParseStringToObservationType(tt.observationType)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetOperatorAddressFromAccAddress(t *testing.T) {
	tests := []struct {
		name    string
		accAddr string
		wantErr bool
	}{
		{"TestValidAccAddress", sample.AccAddress(), false},
		{"TestInvalidAccAddress", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := types.GetOperatorAddressFromAccAddress(tt.accAddr)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAccAddressFromOperatorAddress(t *testing.T) {
	// #nosec G404 test purpose - weak randomness is not an issue here
	r := rand.New(rand.NewSource(1))
	tests := []struct {
		name       string
		valAddress string
		wantErr    bool
	}{
		{"TestValidValAddress", sample.ValAddress(r).String(), false},
		{"TestInvalidValAddress", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := types.GetAccAddressFromOperatorAddress(tt.valAddress)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
