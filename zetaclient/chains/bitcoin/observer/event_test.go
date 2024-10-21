package observer_test

import (
	"testing"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

func Test_CheckProcessability(t *testing.T) {
	// setup compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	// add restricted address
	restrictedBtcAddress := sample.RestrictedBtcAddressTest
	restrictedEvmAddress := sample.RestrictedEVMAddressTest
	cfg.ComplianceConfig.RestrictedAddresses = []string{restrictedBtcAddress, restrictedEvmAddress}
	config.LoadComplianceConfig(cfg)

	// test cases
	tests := []struct {
		name     string
		event    *observer.BTCInboundEvent
		expected observer.InboundProcessability
	}{
		{
			name: "should return InboundProcessabilityGood for a processable inbound event",
			event: &observer.BTCInboundEvent{
				FromAddress: "tb1quhassyrlj43qar0mn0k5sufyp6mazmh2q85lr6ex8ehqfhxpzsksllwrsu",
				ToAddress:   testutils.TSSAddressBTCAthens3,
			},
			expected: observer.InboundProcessabilityGood,
		},
		{
			name: "should return InboundProcessabilityComplianceViolation for a restricted sender address",
			event: &observer.BTCInboundEvent{
				FromAddress: restrictedBtcAddress,
				ToAddress:   testutils.TSSAddressBTCAthens3,
			},
			expected: observer.InboundProcessabilityComplianceViolation,
		},
		{
			name: "should return InboundProcessabilityComplianceViolation for a restricted receiver address in standard memo",
			event: &observer.BTCInboundEvent{
				FromAddress: "tb1quhassyrlj43qar0mn0k5sufyp6mazmh2q85lr6ex8ehqfhxpzsksllwrsu",
				ToAddress:   testutils.TSSAddressBTCAthens3,
				MemoStd: &memo.InboundMemo{
					FieldsV0: memo.FieldsV0{
						Receiver: common.HexToAddress(restrictedEvmAddress),
					},
				},
			},
			expected: observer.InboundProcessabilityComplianceViolation,
		},
		{
			name: "should return InboundProcessabilityComplianceViolation for a restricted revert address in standard memo",
			event: &observer.BTCInboundEvent{
				FromAddress: "tb1quhassyrlj43qar0mn0k5sufyp6mazmh2q85lr6ex8ehqfhxpzsksllwrsu",
				ToAddress:   testutils.TSSAddressBTCAthens3,
				MemoStd: &memo.InboundMemo{
					FieldsV0: memo.FieldsV0{
						RevertOptions: crosschaintypes.RevertOptions{
							RevertAddress: restrictedBtcAddress,
						},
					},
				},
			},
			expected: observer.InboundProcessabilityComplianceViolation,
		},
		{
			name: "should return InboundProcessabilityDonation for a donation inbound event",
			event: &observer.BTCInboundEvent{
				FromAddress: "tb1quhassyrlj43qar0mn0k5sufyp6mazmh2q85lr6ex8ehqfhxpzsksllwrsu",
				ToAddress:   testutils.TSSAddressBTCAthens3,
				MemoBytes:   []byte(constant.DonationMessage),
			},
			expected: observer.InboundProcessabilityDonation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.CheckProcessability()
			require.Equal(t, tt.expected, result)
		})
	}
}
