package observer

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/testutil/sample"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

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
		ts := newTestSuite(t)

		evmBob := sample.EthAddress()
		evmAlice := sample.EthAddress()

		const usdc = "0x5d4b302506645c37ff133b98c4b50a5ae14841659738d6d733d59d0d217a93bf::coin::COIN"

		// Given gateway object from RPC (for "ensuring" the initial scroll cursor)
		gatewayRequest := models.SuiGetObjectRequest{
			ObjectId: ts.gateway.PackageID(),
			Options: models.SuiObjectDataOptions{
				ShowPreviousTransaction: true,
			},
		}

		gatewayObject := models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				ObjectId:            ts.gateway.PackageID(),
				PreviousTransaction: "ABC123_first_tx",
			},
		}

		ts.suiMock.
			On("SuiGetObject", mock.Anything, gatewayRequest).
			Return(gatewayObject, nil)

		// Given list of gateway events...
		expectedQuery := client.EventQuery{
			PackageID: ts.gateway.PackageID(),
			Module:    ts.gateway.Module(),
			Cursor:    "ABC123_first_tx,0",
			Limit:     client.DefaultEventsLimit,
		}

		// ...two of which are valid (1 & 3)
		events := []models.SuiEventResponse{
			ts.SampleEvent("TX_1_ok", string(sui.Deposit), map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "200",
				"sender":    "SUI_BOB",
				"receiver":  evmBob.String(),
			}),
			ts.SampleEvent("TX_2_unrelated_event", "something", map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "200",
				"sender":    "SUI_BOB",
				"receiver":  evmBob.String(),
			}),
			ts.SampleEvent("TX_3_ok", string(sui.DepositAndCall), map[string]any{
				// USDC
				"coin_type": usdc,
				"amount":    "300",
				"sender":    "SUI_ALICE",
				"receiver":  evmAlice.String(),
				"payload":   []any{float64(1), float64(2), float64(3)},
			}),
			ts.SampleEvent("TX_4_invalid_data", string(sui.Deposit), map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "hello",
				"sender":    "SUI_BOB",
				"receiver":  evmBob.String(),
			}),
		}

		ts.suiMock.On("QueryModuleEvents", mock.Anything, expectedQuery).Return(events, "", nil)

		// Given 2 transaction blocks
		ts.OnGetTx("TX_1_ok", "10000", false, nil)
		ts.OnGetTx("TX_3_ok", "20000", false, nil)

		// Given inbound votes catches so we can assert them later
		ts.CatchInboundVotes()

		// ACT
		err := ts.ObserveInbound(ts.ctx)

		// ASSERT
		require.NoError(t, err)

		// Check that final cursor is on INVALID event, that's expected
		assert.Equal(t, "TX_4_invalid_data,0", ts.LastTxScanned())

		// Check for transactions
		assert.Equal(t, 2, len(ts.inboundVotesBag))

		vote1 := ts.inboundVotesBag[0]
		assert.Equal(t, "TX_1_ok", vote1.InboundHash)
		assert.Equal(t, uint64(10_000), vote1.InboundBlockHeight)
		assert.Equal(t, coin.CoinType_Gas, vote1.CoinType)
		assert.Equal(t, false, vote1.IsCrossChainCall)
		assert.Equal(t, math.NewUint(200), vote1.Amount)
		assert.Equal(t, "", vote1.Asset)
		assert.Equal(t, evmBob.String(), vote1.Receiver)

		vote3 := ts.inboundVotesBag[1]
		assert.Equal(t, "TX_3_ok", vote3.InboundHash)
		assert.Equal(t, uint64(20_000), vote3.InboundBlockHeight)
		assert.Equal(t, coin.CoinType_ERC20, vote3.CoinType)
		assert.Equal(t, true, vote3.IsCrossChainCall)
		assert.Equal(t, usdc, vote3.Asset)
		assert.Equal(t, math.NewUint(300), vote3.Amount)
		assert.Equal(t, evmAlice.String(), vote3.Receiver)
		assert.Equal(t, "010203", vote3.Message)

		// Check that other 2 txs are skipped
		assert.Contains(
			t,
			ts.log.String(),
			`unable to parse amount: cannot convert \"hello\" to big.Int: event parse error","message":"Unable to parse event. Skipping"`,
		)
		assert.Contains(
			t,
			ts.log.String(),
			`cannot convert \"hello\" to big.Int: event parse error","message":"Unable to parse event. Skipping"`,
		)
	})

	t.Run("ProcessInboundTrackers", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given inbound tracker
		chainID := ts.Chain().ChainId
		ts.zetaMock.
			On("GetInboundTrackersForChain", mock.Anything, chainID).
			Return([]cctypes.InboundTracker{
				{
					ChainId:  chainID,
					TxHash:   "TX_TRACKER_1",
					CoinType: coin.CoinType_Gas,
				},
			}, nil)

		// Given underlying tx with event
		evmAlice := sample.EthAddress()

		ts.OnGetTx("TX_TRACKER_1", "15000", true, []models.SuiEventResponse{
			ts.SampleEvent("TX_TRACKER_1", string(sui.Deposit), map[string]any{
				"coin_type": string(sui.SUI),
				"amount":    "1000",
				"sender":    "SUI_ALICE",
				"receiver":  evmAlice.String(),
			}),
		})

		// Given votes catcher
		ts.CatchInboundVotes()

		// ACT
		err := ts.ProcessInboundTrackers(ts.ctx)

		// ASSERT
		require.NoError(t, err)

		require.Equal(t, 1, len(ts.inboundVotesBag))

		vote := ts.inboundVotesBag[0]

		assert.Equal(t, "TX_TRACKER_1", vote.InboundHash)
		assert.Equal(t, uint64(15_000), vote.InboundBlockHeight)
		assert.Equal(t, coin.CoinType_Gas, vote.CoinType)
		assert.Equal(t, false, vote.IsCrossChainCall)
		assert.Equal(t, math.NewUint(1000), vote.Amount)
		assert.Equal(t, evmAlice.String(), vote.Receiver)
	})
}

type testSuite struct {
	t        *testing.T
	ctx      context.Context
	zetaMock *mocks.ZetacoreClient
	suiMock  *mocks.SuiClient
	db       *db.DB
	log      *testlog.Log
	gateway  *sui.Gateway

	inboundVotesBag []*cctypes.MsgVoteInbound

	*Observer
}

func newTestSuite(t *testing.T) *testSuite {
	ctx := context.Background()

	chain := chains.SuiMainnet
	chainParams := mocks.MockChainParams(chain.ChainId, 10)
	require.NotEmpty(t, chainParams.GatewayAddress)

	// todo zctx with chain & params (in future PRs)

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

func (ts *testSuite) SampleEvent(txHash, event string, kv map[string]any) models.SuiEventResponse {
	eventType := fmt.Sprintf("%s::%s::%s", ts.gateway.PackageID(), ts.gateway.Module(), event)

	return models.SuiEventResponse{
		Id: models.EventId{
			TxDigest: txHash,
			EventSeq: "0",
		},
		PackageId:         ts.gateway.PackageID(),
		TransactionModule: "gateway",
		Sender:            "SENDER_ABC",
		Type:              eventType,
		ParsedJson:        kv,
	}
}

func (ts *testSuite) OnGetTx(digest, checkpoint string, showEvents bool, events []models.SuiEventResponse) {
	req := models.SuiGetTransactionBlockRequest{
		Digest:  digest,
		Options: models.SuiTransactionBlockOptions{ShowEvents: showEvents},
	}

	res := models.SuiTransactionBlockResponse{
		Digest:     digest,
		Events:     events,
		Checkpoint: checkpoint,
	}

	ts.suiMock.On("SuiGetTransactionBlock", mock.Anything, req).Return(res, nil).Once()
}

func (ts *testSuite) CatchInboundVotes() {
	callback := func(_ context.Context, _, _ uint64, msg *cctypes.MsgVoteInbound) (string, string, error) {
		ts.inboundVotesBag = append(ts.inboundVotesBag, msg)
		return "", "", nil
	}

	ts.zetaMock.
		On("PostVoteInbound", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(callback).
		Maybe()
}
