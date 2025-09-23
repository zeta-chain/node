package observer

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	zetaerrors "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/testutil/sample"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

var someArgStub = map[string]any{}

func TestObserver(t *testing.T) {
	t.Run("PostGasPrice", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given latest checkpoint from RPC
		checkpoint := models.CheckpointResponse{
			// should be used instead of block number
			Epoch:          "333",
			SequenceNumber: "123456",
		}

		ts.suiMock.On("GetLatestCheckpoint", mock.Anything).Return(checkpoint, nil)

		// Given ref price from RPC
		const refGasPrice = uint64(800)
		ts.suiMock.On("SuiXGetReferenceGasPrice", mock.Anything).Return(refGasPrice, nil)

		// Given expected vote for zetacore
		ts.zetaMock.
			On("PostVoteGasPrice", mock.Anything, chains.SuiMainnet, refGasPrice, uint64(0), uint64(333)).
			Return("", nil)

		// ACT
		err := ts.PostGasPrice(ts.ctx)

		// ASSERT
		require.NoError(t, err)
	})

	t.Run("ObserveInbound", func(t *testing.T) {
		// ARRANGE
		previousPackageID := sample.SuiAddress(t)
		originalPackageID := previousPackageID
		ts := newTestSuite(t, func(cfg *testSuiteConfig) {
			cfg.withdrawCapID = sample.SuiAddress(t)
			cfg.previousPackageID = previousPackageID
			cfg.originalPackageID = originalPackageID
		})
		packageID := ts.gateway.PackageID()

		evmBob := sample.EthAddress()
		evmAlice := sample.EthAddress()

		const usdc = "0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf::coin::COIN"

		// Given list of gateway events...
		expectedQuery := client.EventQuery{
			PackageID: originalPackageID,
			Module:    sui.GatewayModule,
			Cursor:    "",
			Limit:     client.DefaultEventsLimit,
		}

		// ...two of which are valid (1 & 3)
		events := []models.SuiEventResponse{
			ts.SampleEvent(packageID, "TX_1_ok", string(sui.DepositEvent), map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "200",
				"sender":    "SUI_BOB",
				"receiver":  evmBob.String(),
			}),
			ts.SampleEvent(packageID, "TX_2_unrelated_event", "something", map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "200",
				"sender":    "SUI_BOB",
				"receiver":  evmBob.String(),
			}),
			ts.SampleEvent(previousPackageID, "TX_3_ok", string(sui.DepositAndCallEvent), map[string]any{
				// USDC
				"coin_type": usdc,
				"amount":    "300",
				"sender":    "SUI_ALICE",
				"receiver":  evmAlice.String(),
				"payload":   preparePayload([]byte{1, 2, 3}),
			}),
			ts.SampleEvent(previousPackageID, "TX_4_invalid_data", string(sui.DepositEvent), map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "hello",
				"sender":    "SUI_BOB",
				"receiver":  evmBob.String(),
			}),
		}

		ts.suiMock.On("QueryModuleEvents", mock.Anything, expectedQuery).Return(events, "", nil)

		// Given 2 transaction blocks
		ts.OnGetTx("TX_1_ok", "10000", true, false, nil)
		ts.OnGetTx("TX_3_ok", "20000", true, false, nil)

		// Given inbound votes catches so we can assert them later
		ts.CatchInboundVotes()
		ts.zetaMock.MockGetCctxByHash(errors.New("not found"))

		// ACT
		err := ts.ObserveInbound(ts.ctx)

		// ASSERT
		require.NoError(t, err)

		// Check that final cursor is on INVALID event, that's expected
		require.Equal(t, "TX_4_invalid_data,0", ts.GetAuxString(originalPackageID))

		// Check for transactions
		require.Equal(t, 2, len(ts.inboundVotesBag))

		vote1 := ts.inboundVotesBag[0]
		require.Equal(t, "TX_1_ok", vote1.InboundHash)
		require.Equal(t, uint64(10_000), vote1.InboundBlockHeight)
		require.Equal(t, coin.CoinType_Gas, vote1.CoinType)
		require.Equal(t, false, vote1.IsCrossChainCall)
		require.Equal(t, math.NewUint(200), vote1.Amount)
		require.Equal(t, "", vote1.Asset)
		require.Equal(t, evmBob.String(), vote1.Receiver)

		vote3 := ts.inboundVotesBag[1]
		require.Equal(t, "TX_3_ok", vote3.InboundHash)
		require.Equal(t, uint64(20_000), vote3.InboundBlockHeight)
		require.Equal(t, coin.CoinType_ERC20, vote3.CoinType)
		require.Equal(t, true, vote3.IsCrossChainCall)
		require.Equal(t, usdc, vote3.Asset)
		require.Equal(t, math.NewUint(300), vote3.Amount)
		require.Equal(t, evmAlice.String(), vote3.Receiver)
		require.Equal(t, "010203", vote3.Message)

		// Check that other 2 txs are skipped
		require.Contains(
			t,
			ts.log.String(),
			`unable to parse amount: cannot convert \"hello\" to big.Int: event parse error","message":"unable to parse event; skipping"`,
		)
		require.Contains(
			t,
			ts.log.String(),
			`cannot convert \"hello\" to big.Int: event parse error","message":"unable to parse event; skipping"`,
		)
	})

	t.Run("ObserveInbound restricted address", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		evmBob := sample.EthAddress()

		// Given compliance config
		cfg := config.Config{
			ComplianceConfig: config.ComplianceConfig{
				RestrictedAddresses: []string{evmBob.String()},
			},
		}
		config.SetRestrictedAddressesFromConfig(cfg)

		// Given a deposit containing restricted address
		packageID := ts.gateway.PackageID()
		expectedQuery := client.EventQuery{
			PackageID: packageID,
			Module:    sui.GatewayModule,
			Cursor:    "",
			Limit:     client.DefaultEventsLimit,
		}

		events := []models.SuiEventResponse{
			ts.SampleEvent(packageID, "TX_restricted", string(sui.DepositEvent), map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "200",
				"sender":    "SUI_BOB",
				"receiver":  evmBob.String(),
			}),
		}

		ts.suiMock.On("QueryModuleEvents", mock.Anything, expectedQuery).Return(events, "", nil)

		// Given transaction block
		ts.OnGetTx("TX_restricted", "10000", true, false, nil)

		// Given inbound votes catches so we can assert them later
		ts.CatchInboundVotes()

		// ACT
		err := ts.ObserveInbound(ts.ctx)

		// ASSERT
		require.NoError(t, err)

		// Check that final cursor is expected on restricted tx
		require.Equal(t, "TX_restricted,0", ts.GetAuxString(packageID))

		// No inbound votes should be created
		require.Empty(t, ts.inboundVotesBag)
	})

	t.Run("ProcessInboundTrackers", func(t *testing.T) {
		// ARRANGE
		originalPackageID := sample.SuiAddress(t)
		ts := newTestSuite(t, func(cfg *testSuiteConfig) {
			cfg.withdrawCapID = sample.SuiAddress(t)
			cfg.previousPackageID = ""
			cfg.originalPackageID = originalPackageID
		})
		packageID := ts.gateway.PackageID()

		// Given inbound tracker
		chainID := ts.Chain().ChainId
		txHash := "TX_TRACKER_1"
		trackers := []cctypes.InboundTracker{
			{
				ChainId:  chainID,
				TxHash:   txHash,
				CoinType: coin.CoinType_Gas,
			},
		}

		ts.zetaMock.On("GetInboundTrackersForChain", mock.Anything, chainID).Return(trackers, nil)

		// Given underlying tx with event
		evmAlice := sample.EthAddress()

		ts.OnGetTx(txHash, "15000", true, true, []models.SuiEventResponse{
			ts.SampleEvent(packageID, txHash, string(sui.DepositEvent), map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "1000",
				"sender":    "SUI_ALICE",
				"receiver":  evmAlice.String(),
			}),
		})

		ts.zetaMock.MockGetCctxByHash(errors.New("not found"))

		// Given votes catcher
		ts.CatchInboundVotes()

		// ACT
		err := ts.ProcessInboundTrackers(ts.ctx)

		// ASSERT
		require.NoError(t, err)
		require.Equal(t, 1, len(ts.inboundVotesBag))

		vote := ts.inboundVotesBag[0]

		require.Equal(t, txHash, vote.InboundHash)
		require.Equal(t, uint64(15_000), vote.InboundBlockHeight)
		require.Equal(t, coin.CoinType_Gas, vote.CoinType)
		require.Equal(t, false, vote.IsCrossChainCall)
		require.Equal(t, math.NewUint(1000), vote.Amount)
		require.Equal(t, evmAlice.String(), vote.Receiver)
	})

	t.Run("ProcessOutboundTrackers", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given cctx
		const nonce = 333
		cctx := sample.CrossChainTxV2(t, "0x123")
		cctx.OutboundParams = []*cctypes.OutboundParams{{TssNonce: nonce}}

		ts.MockCCTXByNonce(cctx)

		// Given outbound tracker
		const digest = "0xSuiTxHash"
		tracker := cctypes.OutboundTracker{
			Index:    "0xAAA",
			ChainId:  ts.Chain().ChainId,
			Nonce:    nonce,
			HashList: []*cctypes.TxHash{{TxHash: digest}},
		}

		ts.MockOutboundTrackers([]cctypes.OutboundTracker{tracker})

		// Given Sui tx signature
		sigBase64, err := sui.SerializeSignatureECDSA([65]byte{1, 2, 3}, ts.TSS().PubKey().AsECDSA())
		require.NoError(t, err)

		// Given Sui tx
		eventNonce := fmt.Sprintf("%d", nonce)
		tx := models.SuiTransactionBlockResponse{
			Digest:     digest,
			Checkpoint: "123",
			Effects: models.SuiEffects{
				Status: models.ExecutionStatus{Status: client.TxStatusSuccess},
			},
			Transaction: models.SuiTransactionBlock{
				Data: models.SuiTransactionBlockData{
					Transaction: models.SuiTransactionBlockKind{
						Inputs: []models.SuiCallArg{
							someArgStub,
							someArgStub,
							map[string]any{
								"type":      "pure",
								"valueType": "u64",
								"value":     eventNonce,
							},
							someArgStub,
							someArgStub,
							someArgStub,
						},
					},
				},
				TxSignatures: []string{sigBase64},
			},
			Events: []models.SuiEventResponse{
				{
					Id:        models.EventId{TxDigest: digest, EventSeq: "1"},
					PackageId: ts.Gateway().PackageID(),
					Sender:    "0xSuiSender",
					Type:      ts.EventType(string(sui.WithdrawEvent)),
					ParsedJson: map[string]any{
						"coin_type": string(sui.SUI),
						"amount":    "200",
						"sender":    "0xSuiSender",
						"receiver":  "0xSuiReceiver",
						"nonce":     eventNonce,
					},
				},
			},
		}

		ts.MockGetTxOnce(tx)

		// ACT
		err = ts.ProcessOutboundTrackers(ts.ctx)

		// ASSERT
		require.NoError(t, err)
		require.True(t, ts.OutboundCreated(nonce))
		require.False(t, ts.OutboundCreated(nonce+1))
	})

	t.Run("VoteOutbound successful withdrawal", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given Sui Gateway
		gw := ts.Gateway()

		// Given cctx
		const nonce = 333
		cctx := sample.CrossChainTxV2(t, "0x123")
		cctx.OutboundParams = []*cctypes.OutboundParams{{TssNonce: nonce}}

		// Given Sui receiver
		const receiver = "0xAliceOnSui"

		// Given a valid Sui outbound tx with Withdrawal event
		const digest = "0xSuiTxDigest"
		tx := models.SuiTransactionBlockResponse{
			Digest:     digest,
			Checkpoint: "999",
			Effects: models.SuiEffects{
				Status: models.ExecutionStatus{Status: client.TxStatusSuccess},
				GasUsed: models.GasCostSummary{
					ComputationCost: "200",
					StorageCost:     "300",
					StorageRebate:   "50",
				},
			},
			Events: []models.SuiEventResponse{{
				Id:        models.EventId{TxDigest: digest, EventSeq: "1"},
				PackageId: gw.PackageID(),
				Sender:    ts.TSS().PubKey().AddressSui(),
				Type:      fmt.Sprintf("%s::%s::%s", gw.PackageID(), sui.GatewayModule, "WithdrawEvent"),
				ParsedJson: map[string]any{
					"coin_type": string(sui.SUI),
					"amount":    "200",
					"sender":    ts.TSS().PubKey().AddressSui(),
					"receiver":  receiver,
					"nonce":     fmt.Sprintf("%d", nonce),
				},
			}},
		}

		// What was fetched during ProcessOutboundTracker(...)
		ts.setTx(tx, nonce)

		// Given a gas price that was set during PostGasPrice(...)
		ts.setLatestGasPrice(1000)

		// Given outbound votes catcher
		ts.CatchOutboundVotes()

		// ACT
		err := ts.VoteOutbound(ts.ctx, cctx)

		// ASSERT
		require.NoError(t, err)
		require.Len(t, ts.outboundVotesBag, 1)

		vote := ts.outboundVotesBag[0]

		// common
		require.Equal(t, chains.ReceiveStatus_success, vote.Status) // success
		require.Equal(t, cctx.Index, vote.CctxHash)
		require.Equal(t, uint64(nonce), vote.OutboundTssNonce)
		require.Equal(t, ts.Chain().ChainId, vote.OutboundChain)

		// digest + checkpoint
		require.Equal(t, digest, vote.ObservedOutboundHash)
		require.Equal(t, uint64(999), vote.ObservedOutboundBlockHeight)

		// amount
		require.Equal(t, coin.CoinType_Gas, vote.CoinType)
		require.Equal(t, uint64(200), vote.ValueReceived.Uint64())

		// gas
		require.Equal(t, uint64(0), vote.ObservedOutboundEffectiveGasLimit)
		require.Equal(t, uint64(1000), vote.ObservedOutboundEffectiveGasPrice.Uint64())
		require.Equal(t, uint64(200+300-50), vote.ObservedOutboundGasUsed)
	})

	t.Run("VoteOutbound failed withdrawal", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given cctx
		const nonce = 333
		cctx := sample.CrossChainTxV2(t, "0x123")
		cctx.OutboundParams = []*cctypes.OutboundParams{
			{
				Amount:   math.NewUint(200),
				TssNonce: nonce,
			}}

		// Given a valid Sui outbound tx with Withdrawal event
		const digest = "0xSuiTxDigest"
		eventNonce := fmt.Sprintf("%d", nonce+1) // cancel tx event nonce == cctx nonce + 1
		tx := models.SuiTransactionBlockResponse{
			Digest:     digest,
			Checkpoint: "999",
			Effects: models.SuiEffects{
				Status: models.ExecutionStatus{Status: client.TxStatusSuccess},
				GasUsed: models.GasCostSummary{
					ComputationCost: "200",
					StorageCost:     "300",
					StorageRebate:   "50",
				},
			},
			Events: []models.SuiEventResponse{
				{
					Id:        models.EventId{TxDigest: digest, EventSeq: "1"},
					PackageId: ts.Gateway().PackageID(),
					Sender:    ts.TSS().PubKey().AddressSui(),
					Type:      ts.EventType(string(sui.CancelTxEvent)),
					ParsedJson: map[string]any{
						"sender": ts.TSS().PubKey().AddressSui(),
						"nonce":  eventNonce,
					},
				},
			},
		}

		// What was fetched during ProcessOutboundTracker(...)
		ts.setTx(tx, nonce)

		// Given a gas price that was set during PostGasPrice(...)
		ts.setLatestGasPrice(1000)

		// Given outbound votes catcher
		ts.CatchOutboundVotes()

		// ACT
		err := ts.VoteOutbound(ts.ctx, cctx)

		// ASSERT
		require.NoError(t, err)
		require.Len(t, ts.outboundVotesBag, 1)

		vote := ts.outboundVotesBag[0]

		// common
		require.Equal(t, chains.ReceiveStatus_failed, vote.Status) // failure
		require.Equal(t, cctx.Index, vote.CctxHash)
		require.Equal(t, uint64(nonce), vote.OutboundTssNonce)
		require.Equal(t, ts.Chain().ChainId, vote.OutboundChain)

		// digest + checkpoint
		require.Equal(t, digest, vote.ObservedOutboundHash)
		require.Equal(t, uint64(999), vote.ObservedOutboundBlockHeight)

		// amount
		require.Equal(t, coin.CoinType_Gas, vote.CoinType)
		require.Equal(t, uint64(200), vote.ValueReceived.Uint64())

		// gas
		require.Equal(t, uint64(0), vote.ObservedOutboundEffectiveGasLimit)
		require.Equal(t, uint64(1000), vote.ObservedOutboundEffectiveGasPrice.Uint64())
		require.Equal(t, uint64(200+300-50), vote.ObservedOutboundGasUsed)
	})
}

func Test_MigrateCursorForAuthenticatedCallUpgrade(t *testing.T) {
	withdrawCapID := sample.SuiAddress(t)
	previousPackageID := sample.SuiAddress(t)
	originalPackageID := previousPackageID

	tests := []struct {
		name              string
		cursor            string
		originalPackageID string
		wantKey           string
	}{
		{
			name:              "migration with empty original package ID",
			cursor:            "HBgprKGko6Kk1q7cjZpg1KHVUqFXE7PFhT2M9DaPAezi,0",
			originalPackageID: "",
		},
		{
			name:              "migration with non-empty original package ID",
			cursor:            "DVmb9QoJKRvSZfw6SkL8Zp1vRESxS7RuHy6XEYGT7WWn,0",
			originalPackageID: originalPackageID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			ts := newTestSuite(t, func(cfg *testSuiteConfig) {
				cfg.withdrawCapID = withdrawCapID
				cfg.previousPackageID = previousPackageID
				cfg.originalPackageID = tt.originalPackageID
			})
			packageID := ts.Gateway().PackageID()
			originalPackageID := ts.Gateway().Original().PackageID()

			// write a pre-v35 old cursor to the DB
			err := ts.WriteLastTxScannedToDB(tt.cursor)
			require.NoError(t, err)
			ts.WithLastTxScanned(tt.cursor)

			// ensure the old cursor is set
			oldCursor, err := ts.ReadLastTxScannedFromDB()
			require.NoError(t, err)
			require.Equal(t, tt.cursor, oldCursor)
			require.Equal(t, tt.cursor, ts.LastTxScanned())

			// ACT-1
			err = ts.MigrateCursorForAuthenticatedCallUpgrade()

			// ASSERT
			require.NoError(t, err)

			// ensure the new cursor is stored under original package ID
			newCursor, err := ts.ReadAuxStringFromDB(originalPackageID)
			require.NoError(t, err)
			require.Equal(t, tt.cursor, newCursor)
			require.Equal(t, tt.cursor, ts.GetAuxString(originalPackageID))

			// ensure nothing is stored under new package ID
			if packageID != originalPackageID {
				cursor, err := ts.ReadAuxStringFromDB(packageID)
				require.ErrorContains(t, err, "record not found")
				require.Empty(t, cursor)
				require.Empty(t, ts.GetAuxString(packageID))
			}

			// ensure the old cursor is set to empty
			oldCursor, err = ts.ReadLastTxScannedFromDB()
			require.NoError(t, err)
			require.Empty(t, oldCursor)
			require.Empty(t, ts.LastTxScanned())

			// ACT-2, migrate again
			err = ts.MigrateCursorForAuthenticatedCallUpgrade()

			// ASSERT
			require.NoError(t, err)

			// ensure the new cursor stay untouched
			cursor, err := ts.ReadAuxStringFromDB(originalPackageID)
			require.NoError(t, err)
			require.Equal(t, tt.cursor, cursor)
			require.Equal(t, tt.cursor, ts.GetAuxString(originalPackageID))
		})
	}
}

type testSuite struct {
	t        *testing.T
	ctx      context.Context
	zetaMock *mocks.ZetacoreClient
	suiMock  *mocks.SuiClient
	db       *db.DB
	log      *testlog.Log
	gateway  *sui.Gateway

	inboundVotesBag  []*cctypes.MsgVoteInbound
	outboundVotesBag []*cctypes.MsgVoteOutbound

	*Observer
}

type testSuiteConfig struct {
	withdrawCapID     string
	previousPackageID string
	originalPackageID string
}

func newTestSuite(t *testing.T, opts ...func(*testSuiteConfig)) *testSuite {
	ctx := context.Background()

	var cfg testSuiteConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	chain := chains.SuiMainnet
	chainParams := mocks.MockChainParams(chain.ChainId, 10)
	require.NotEmpty(t, chainParams.GatewayAddress)

	// append withdraw cap ID, previous package ID and original package ID if provided
	if cfg.withdrawCapID != "" && cfg.originalPackageID != "" {
		chainParams.GatewayAddress = fmt.Sprintf("%s,%s,%s,%s", chainParams.GatewayAddress, cfg.withdrawCapID, cfg.previousPackageID, cfg.originalPackageID)
	}

	zetacore := mocks.NewZetacoreClient(t).
		WithZetaChain().
		WithKeys(&keys.Keys{
			OperatorAddress: sample.Bech32AccAddress(),
		})

	tss := mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	log := testlog.New(t)
	logger := base.Logger{
		Std:        log.Logger,
		Compliance: log.Logger,
	}

	baseObserver, err := base.NewObserver(chain, chainParams, zetacore, tss, 1000, nil, database, logger)
	require.NoError(t, err)

	suiMock := mocks.NewSuiClient(t)

	gw, err := sui.NewGatewayFromPairID(chainParams.GatewayAddress)
	require.NoError(t, err)

	observer := New(baseObserver, suiMock, gw)

	return &testSuite{
		t:        t,
		ctx:      ctx,
		zetaMock: zetacore,
		suiMock:  suiMock,
		db:       database,
		log:      log,
		gateway:  gw,
		Observer: observer,
	}
}

func (ts *testSuite) SampleEvent(packageID, txHash, event string, kv map[string]any) models.SuiEventResponse {
	eventType := fmt.Sprintf("%s::%s::%s", packageID, sui.GatewayModule, event)

	return models.SuiEventResponse{
		Id: models.EventId{
			TxDigest: txHash,
			EventSeq: "0",
		},
		PackageId:         packageID,
		TransactionModule: "gateway",
		Sender:            "SENDER_ABC",
		Type:              eventType,
		ParsedJson:        kv,
	}
}

func (ts *testSuite) OnGetTx(digest, checkpoint string, showEffects, showEvents bool, events []models.SuiEventResponse) {
	req := models.SuiGetTransactionBlockRequest{
		Digest: digest,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: showEffects,
			ShowEvents:  showEvents,
		},
	}

	res := models.SuiTransactionBlockResponse{
		Digest: digest,
		Effects: models.SuiEffects{
			Status: models.ExecutionStatus{Status: client.TxStatusSuccess},
		},
		Events:     events,
		Checkpoint: checkpoint,
	}

	ts.suiMock.On("SuiGetTransactionBlock", mock.Anything, req).Return(res, nil).Once()
}

func (ts *testSuite) MockGetTxOnce(tx models.SuiTransactionBlockResponse) {
	ts.suiMock.On("SuiGetTransactionBlock", mock.Anything, mock.Anything).Return(tx, nil).Once()
}

func (ts *testSuite) CatchInboundVotes() {
	callback := func(_ context.Context, _, _ uint64, msg *cctypes.MsgVoteInbound, _ chan<- zetaerrors.ErrTxMonitor) (string, string, error) {
		ts.inboundVotesBag = append(ts.inboundVotesBag, msg)
		return "", "", nil
	}

	ts.zetaMock.
		On("PostVoteInbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(callback).
		Maybe()
}

func (ts *testSuite) CatchOutboundVotes() {
	callback := func(_ context.Context, _, _ uint64, msg *cctypes.MsgVoteOutbound) (string, string, error) {
		ts.outboundVotesBag = append(ts.outboundVotesBag, msg)
		return "", "", nil
	}

	ts.zetaMock.
		On("PostVoteOutbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(callback).
		Maybe()
}

func (ts *testSuite) MockCCTXByNonce(cctx *cctypes.CrossChainTx) *mock.Call {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	return ts.zetaMock.
		On("GetCctxByNonce", ts.ctx, ts.Chain().ChainId, nonce).
		Return(cctx, nil)
}

func (ts *testSuite) MockOutboundTrackers(trackers []cctypes.OutboundTracker) *mock.Call {
	return ts.zetaMock.
		On("GetOutboundTrackers", mock.Anything, ts.Chain().ChainId).
		Return(trackers, nil)
}

func (ts *testSuite) EventType(event string) string {
	return fmt.Sprintf("%s::%s::%s", ts.gateway.PackageID(), sui.GatewayModule, event)
}

func preparePayload(payload []byte) []any {
	payloadBytes := []byte(base64.StdEncoding.EncodeToString(payload))

	var out []any
	for _, p := range payloadBytes {
		out = append(out, float64(p))
	}

	return out
}
