package types_test

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/types"
)

func Test_DecodeMemo(t *testing.T) {
	testReceiver := sample.EthAddress()

	// test cases
	tests := []struct {
		name             string
		event            *types.InboundEvent
		expectedReceiver string
		errMsg           string
	}{
		{
			name: "should decode receiver address successfully",
			event: &types.InboundEvent{
				Memo: testReceiver.Bytes(),
			},
			expectedReceiver: testReceiver.Hex(),
		},
		{
			name: "should skip decoding donation message",
			event: &types.InboundEvent{
				Memo: []byte(constant.DonationMessage),
			},
			expectedReceiver: "",
		},
		{
			name: "should return error if got an empty receiver address",
			event: &types.InboundEvent{
				Memo: []byte(""),
			},
			errMsg:           "got empty receiver address from memo",
			expectedReceiver: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.DecodeMemo()
			if tt.errMsg != "" {
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedReceiver, tt.event.Receiver)
		})
	}
}

func Test_Catetory(t *testing.T) {
	// setup compliance config
	cfg := config.Config{
		ComplianceConfig: sample.ComplianceConfig(),
	}
	config.SetRestrictedAddressesFromConfig(cfg)

	// test cases
	tests := []struct {
		name     string
		event    *types.InboundEvent
		expected types.InboundCategory
	}{
		{
			name: "should return InboundCategoryProcessable for a processable inbound event",
			event: &types.InboundEvent{
				Sender:   sample.SolanaAddress(t),
				Receiver: sample.EthAddress().Hex(),
			},
			expected: types.InboundCategoryProcessable,
		},
		{
			name: "should return InboundCategoryRestricted for a restricted sender address",
			event: &types.InboundEvent{
				Sender:   sample.RestrictedSolAddressTest,
				Receiver: sample.EthAddress().Hex(),
			},
			expected: types.InboundCategoryRestricted,
		},
		{
			name: "should return InboundCategoryRestricted for a restricted receiver address",
			event: &types.InboundEvent{
				Sender:   sample.SolanaAddress(t),
				Receiver: sample.RestrictedSolAddressTest,
			},
			expected: types.InboundCategoryRestricted,
		},
		{
			name: "should return InboundCategoryRestricted for a restricted receiver address in memo",
			event: &types.InboundEvent{
				Sender:   sample.SolanaAddress(t),
				Receiver: sample.EthAddress().Hex(),
				Memo:     ethcommon.HexToAddress(sample.RestrictedEVMAddressTest).Bytes(),
			},
			expected: types.InboundCategoryRestricted,
		},
		{
			name: "should return InboundCategoryDonation for a donation inbound event",
			event: &types.InboundEvent{
				Sender:   sample.SolanaAddress(t),
				Receiver: sample.EthAddress().Hex(),
				Memo:     []byte(constant.DonationMessage),
			},
			expected: types.InboundCategoryDonation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.Category()
			require.Equal(t, tt.expected, result)
		})
	}
}
