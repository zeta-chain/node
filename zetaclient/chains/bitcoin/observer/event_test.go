package observer_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/zeta-chain/node/testutil"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// createTestBtcEvent creates a test BTC inbound event
func createTestBtcEvent(
	t *testing.T,
	net *chaincfg.Params,
	memo []byte,
	memoStd *memo.InboundMemo,
) observer.BTCInboundEvent {
	return observer.BTCInboundEvent{
		FromAddress: sample.BtcAddressP2WPKH(t, net),
		ToAddress:   sample.EthAddress().Hex(),
		MemoBytes:   memo,
		MemoStd:     memoStd,
		TxHash:      sample.Hash().Hex(),
		BlockNumber: 123456,
	}
}

func Test_CheckProcessability(t *testing.T) {
	// setup compliance config
	cfg := config.Config{
		ComplianceConfig: sample.ComplianceConfig(),
	}
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
				FromAddress: sample.RestrictedBtcAddressTest,
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
						Receiver: common.HexToAddress(sample.RestrictedEVMAddressTest),
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
							RevertAddress: sample.RestrictedBtcAddressTest,
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
			result := tt.event.Processability()
			require.Equal(t, tt.expected, result)
		})
	}
}

func Test_DecodeEventMemoBytes(t *testing.T) {
	// test cases
	tests := []struct {
		name             string
		chainID          int64
		event            *observer.BTCInboundEvent
		expectedMemoStd  *memo.InboundMemo
		expectedReceiver common.Address
		donation         bool
		errMsg           string
	}{
		{
			name:    "should decode standard memo bytes successfully",
			chainID: chains.BitcoinTestnet.ChainId,
			event: &observer.BTCInboundEvent{
				// a deposit and call
				MemoBytes: testutil.HexToBytes(
					t,
					"5a0110032d07a9cbd57dcca3e2cf966c88bc874445b6e3b60d68656c6c6f207361746f736869",
				),
			},
			expectedMemoStd: &memo.InboundMemo{
				Header: memo.Header{
					Version:     0,
					EncodingFmt: memo.EncodingFmtCompactShort,
					OpCode:      memo.OpCodeDepositAndCall,
					DataFlags:   3, // reciever + payload
				},
				FieldsV0: memo.FieldsV0{
					Receiver: common.HexToAddress("0x2D07A9CBd57DCca3E2cF966C88Bc874445b6E3B6"),
					Payload:  []byte("hello satoshi"),
				},
			},
		},
		{
			name:    "should fall back to legacy memo successfully",
			chainID: chains.BitcoinTestnet.ChainId,
			event: &observer.BTCInboundEvent{
				// raw address + payload
				MemoBytes: testutil.HexToBytes(t, "2d07a9cbd57dcca3e2cf966c88bc874445b6e3b668656c6c6f207361746f736869"),
			},
			expectedReceiver: common.HexToAddress("0x2D07A9CBd57DCca3E2cF966C88Bc874445b6E3B6"),
		},
		{
			name:    "should disable standard memo for Bitcoin mainnet",
			chainID: chains.BitcoinMainnet.ChainId,
			event: &observer.BTCInboundEvent{
				// a deposit and call
				MemoBytes: testutil.HexToBytes(
					t,
					"5a0110032d07a9cbd57dcca3e2cf966c88bc874445b6e3b60d68656c6c6f207361746f736869",
				),
			},
			expectedReceiver: common.HexToAddress("0x5A0110032d07A9cbd57dcCa3e2Cf966c88bC8744"),
		},
		{
			name:    "should do nothing for donation message",
			chainID: chains.BitcoinTestnet.ChainId,
			event: &observer.BTCInboundEvent{
				MemoBytes: []byte(constant.DonationMessage),
			},
			donation: true,
		},
		{
			name:    "should return error if standard memo contains improper data",
			chainID: chains.BitcoinTestnet.ChainId,
			event: &observer.BTCInboundEvent{
				// a deposit and call, receiver is empty ZEVM address
				MemoBytes: testutil.HexToBytes(
					t,
					"5a01100300000000000000000000000000000000000000000d68656c6c6f207361746f736869",
				),
			},
			errMsg: "standard memo contains improper data",
		},
		{
			name:    "should return error if standard memo validation failed",
			chainID: chains.BitcoinTestnet.ChainId,
			event: &observer.BTCInboundEvent{
				// a no asset call opCode passed, not supported at the moment
				MemoBytes: testutil.HexToBytes(
					t,
					"5a0120032d07a9cbd57dcca3e2cf966c88bc874445b6e3b60d68656c6c6f207361746f736869",
				),
			},
			errMsg: "invalid standard memo for bitcoin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.DecodeMemoBytes(tt.chainID)
			if tt.errMsg != "" {
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)

			// donation message will skip decoding, so ToAddress will be left empty
			if tt.donation {
				require.Empty(t, tt.event.ToAddress)
				return
			}

			// if it's a standard memo
			if tt.expectedMemoStd != nil {
				require.NotNil(t, tt.event.MemoStd)
				require.Equal(t, tt.expectedMemoStd.Receiver.Hex(), tt.event.ToAddress)
				require.Equal(t, tt.expectedMemoStd, tt.event.MemoStd)
			} else {
				// if it's a legacy memo, check receiver address only
				require.Equal(t, tt.expectedReceiver.Hex(), tt.event.ToAddress)
			}
		})
	}
}

func Test_ValidateStandardMemo(t *testing.T) {
	// test cases
	tests := []struct {
		name   string
		memo   memo.InboundMemo
		errMsg string
	}{
		{
			name: "validation should pass for a valid standard memo",
			memo: memo.InboundMemo{
				Header: memo.Header{
					OpCode: memo.OpCodeDepositAndCall,
				},
				FieldsV0: memo.FieldsV0{
					RevertOptions: crosschaintypes.RevertOptions{
						RevertAddress: sample.BtcAddressP2WPKH(t, &chaincfg.TestNet3Params),
					},
				},
			},
		},
		{
			name: "NoAssetCall is disabled for Bitcoin",
			memo: memo.InboundMemo{
				Header: memo.Header{
					OpCode: memo.OpCodeCall,
				},
			},
			errMsg: "NoAssetCall is disabled for Bitcoin",
		},
		{
			name: "should return error on invalid revert address",
			memo: memo.InboundMemo{
				FieldsV0: memo.FieldsV0{
					RevertOptions: crosschaintypes.RevertOptions{
						// not a BTC address
						RevertAddress: "0x2D07A9CBd57DCca3E2cF966C88Bc874445b6E3B6",
					},
				},
			},
			errMsg: "invalid revert address in memo",
		},
		{
			name: "should return error if revert address is not a supported address type",
			memo: memo.InboundMemo{
				FieldsV0: memo.FieldsV0{
					RevertOptions: crosschaintypes.RevertOptions{
						// address not supported
						RevertAddress: "035e4ae279bd416b5da724972c9061ec6298dac020d1e3ca3f06eae715135cdbec",
					},
				},
			},
			errMsg: "unsupported revert address in memo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := observer.ValidateStandardMemo(tt.memo, chains.BitcoinTestnet.ChainId)
			if tt.errMsg != "" {
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func Test_CheckEventProcessability(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create test observer
	ob := MockBTCObserver(t, chain, params, nil)

	// setup compliance config
	cfg := config.Config{
		ComplianceConfig: sample.ComplianceConfig(),
	}
	config.LoadComplianceConfig(cfg)

	// test cases
	tests := []struct {
		name   string
		event  observer.BTCInboundEvent
		result bool
	}{
		{
			name:   "should return true for processable event",
			event:  createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("a memo"), nil),
			result: true,
		},
		{
			name:   "should return false on donation message",
			event:  createTestBtcEvent(t, &chaincfg.MainNetParams, []byte(constant.DonationMessage), nil),
			result: false,
		},
		{
			name: "should return false on compliance violation",
			event: createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("a memo"), &memo.InboundMemo{
				FieldsV0: memo.FieldsV0{
					Receiver: common.HexToAddress(sample.RestrictedEVMAddressTest),
				},
			}),
			result: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ob.CheckEventProcessability(tt.event)
			require.Equal(t, tt.result, result)
		})
	}
}

func Test_NewInboundVoteFromLegacyMemo(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create test observer
	ob := MockBTCObserver(t, chain, params, nil)
	zetacoreClient := mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{}).WithZetaChain()
	ob.WithZetacoreClient(zetacoreClient)

	t.Run("should create new inbound vote msg V1", func(t *testing.T) {
		// create test event
		event := createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("dummy memo"), nil)

		// test amount
		amountSats := big.NewInt(1000)

		// expected vote
		expectedVote := crosschaintypes.MsgVoteInbound{
			Sender:             event.FromAddress,
			SenderChainId:      chain.ChainId,
			TxOrigin:           event.FromAddress,
			Receiver:           event.ToAddress,
			ReceiverChain:      ob.ZetacoreClient().Chain().ChainId,
			Amount:             cosmosmath.NewUint(amountSats.Uint64()),
			Message:            hex.EncodeToString(event.MemoBytes),
			InboundHash:        event.TxHash,
			InboundBlockHeight: event.BlockNumber,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 0,
			},
			CoinType:                coin.CoinType_Gas,
			ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V1,
			RevertOptions:           crosschaintypes.NewEmptyRevertOptions(), // ignored by V1
		}

		// create new inbound vote V1
		vote := ob.NewInboundVoteFromLegacyMemo(&event, amountSats)
		require.Equal(t, expectedVote, *vote)
	})
}

func Test_NewInboundVoteFromStdMemo(t *testing.T) {
	// can use any bitcoin chain for testing
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create test observer
	ob := MockBTCObserver(t, chain, params, nil)
	zetacoreClient := mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{}).WithZetaChain()
	ob.WithZetacoreClient(zetacoreClient)

	t.Run("should create new inbound vote msg with standard memo", func(t *testing.T) {
		// create revert options
		revertOptions := crosschaintypes.NewEmptyRevertOptions()
		revertOptions.RevertAddress = sample.BtcAddressP2WPKH(t, &chaincfg.MainNetParams)

		// create test event
		receiver := sample.EthAddress()
		event := createTestBtcEvent(t, &chaincfg.MainNetParams, []byte("dymmy"), &memo.InboundMemo{
			FieldsV0: memo.FieldsV0{
				Receiver:      receiver,
				Payload:       []byte("some payload"),
				RevertOptions: revertOptions,
			},
		})

		// test amount
		amountSats := big.NewInt(1000)

		// expected vote
		memoBytesExpected := append(event.MemoStd.Receiver.Bytes(), event.MemoStd.Payload...)
		expectedVote := crosschaintypes.MsgVoteInbound{
			Sender:             revertOptions.RevertAddress, // should be overridden by revert address
			SenderChainId:      chain.ChainId,
			TxOrigin:           event.FromAddress,
			Receiver:           event.ToAddress,
			ReceiverChain:      ob.ZetacoreClient().Chain().ChainId,
			Amount:             cosmosmath.NewUint(amountSats.Uint64()),
			Message:            hex.EncodeToString(memoBytesExpected), // a simulated legacy memo
			InboundHash:        event.TxHash,
			InboundBlockHeight: event.BlockNumber,
			CallOptions: &crosschaintypes.CallOptions{
				GasLimit: 0,
			},
			CoinType:                coin.CoinType_Gas,
			ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V1,
			RevertOptions:           crosschaintypes.NewEmptyRevertOptions(), // ignored by V1
		}

		// create new inbound vote V1 with standard memo
		vote := ob.NewInboundVoteFromStdMemo(&event, amountSats)
		require.Equal(t, expectedVote, *vote)
	})
}
