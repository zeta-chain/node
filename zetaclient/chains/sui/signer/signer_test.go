package signer

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/testutil/sample"
	cc "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/sui/client"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
)

func TestSigner(t *testing.T) {
	t.Run("ProcessCCTX", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		const zetaHeight = 1000

		// Given cctx
		nonce := uint64(123)
		amount := math.NewUint(100_000)
		receiver := "0xdecb47015beebed053c19ef48fe4d722fa3870f567133d235ebe3a70da7b0000"

		cctx := sample.CrossChainTxV2(t, "0xABC123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams = []*cc.OutboundParams{{
			Receiver:        receiver,
			ReceiverChainId: ts.Chain.ChainId,
			CoinType:        coin.CoinType_Gas,
			Amount:          amount,
			TssNonce:        nonce,
			GasPrice:        "1000",
			CallOptions: &cc.CallOptions{
				GasLimit: 42,
			},
		}}

		// Given mocked gateway nonce
		ts.MockGatewayNonce(nonce)

		// Given mocked WithdrawCapID
		const withdrawCapID = "0xWithdrawCapID"
		ts.MockWithdrawCapID(withdrawCapID)

		// Given mocked MessageContextID
		const messageContextID = "0xMessageContextID"
		ts.MockMessageContextID(messageContextID)

		// Given expected MoveCall
		txBytes := base64.StdEncoding.EncodeToString([]byte("raw_tx_bytes"))

		ts.MockMoveCall(func(req models.MoveCallRequest) {
			require.Equal(t, ts.TSS.PubKey().AddressSui(), req.Signer)
			require.Equal(t, ts.Gateway.PackageID(), req.PackageObjectId)
			require.Equal(t, "withdraw", req.Function)

			expectedArgs := []any{
				ts.Gateway.ObjectID(),
				amount.String(),
				fmt.Sprintf("%d", nonce),
				receiver,
				"42000",
				withdrawCapID,
			}
			require.Equal(t, expectedArgs, req.Arguments)
		}, txBytes)

		// Given expected SuiExecuteTransactionBlock
		const digest = "0xTransactionBlockDigest"
		ts.MockExec(func(req models.SuiExecuteTransactionBlockRequest) {
			require.Equal(t, txBytes, req.TxBytes)
			require.NotEmpty(t, req.Signature)
		}, digest)

		// Given included tx from Sui RPC
		ts.SuiMock.
			On("SuiGetTransactionBlock", mock.Anything, mock.Anything).
			Return(models.SuiTransactionBlockResponse{
				Digest: digest,
				Effects: models.SuiEffects{
					Status: models.ExecutionStatus{
						Status: client.TxStatusSuccess,
					},
				},
				Checkpoint: "1000000",
			}, nil)

		// ACT
		err := ts.Signer.ProcessCCTX(ts.Ctx, cctx, zetaHeight)

		// ASSERT
		require.NoError(t, err)

		// Wait for vote posting
		wait := func() bool {
			if len(ts.TrackerBag) == 0 {
				return false
			}

			vote := ts.TrackerBag[0]
			return vote.hash == digest && vote.nonce == nonce
		}

		require.Eventually(t, wait, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("ProcessCCTX dry-mode", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)
		ts.ClientMode = mode.DryMode

		const zetaHeight = 1000

		// Given cctx
		nonce := uint64(123)
		amount := math.NewUint(100_000)
		receiver := "0xdecb47015beebed053c19ef48fe4d722fa3870f567133d235ebe3a70da7b0000"

		cctx := sample.CrossChainTxV2(t, "0xABC123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams = []*cc.OutboundParams{{
			Receiver:        receiver,
			ReceiverChainId: ts.Chain.ChainId,
			CoinType:        coin.CoinType_Gas,
			Amount:          amount,
			TssNonce:        nonce,
			GasPrice:        "1000",
			CallOptions: &cc.CallOptions{
				GasLimit: 42,
			},
		}}

		// Given mocked gateway nonce
		ts.MockGatewayNonce(nonce)

		// Given mocked WithdrawCapID
		const withdrawCapID = "0xWithdrawCapID"
		ts.MockWithdrawCapID(withdrawCapID)

		// ACT
		err := ts.Signer.ProcessCCTX(ts.Ctx, cctx, zetaHeight)

		// ASSERT
		require.NoError(t, err)

		// Wait for possible vote posting
		wait := func() bool {
			return len(ts.TrackerBag) == 0
		}

		require.Eventually(t, wait, 5*time.Second, 100*time.Millisecond)

		ts.SuiMock.AssertNotCalled(t, "MoveCall", mock.Anything, mock.Anything)
		ts.SuiMock.AssertNotCalled(t, "SuiExecuteTransactionBlock", mock.Anything, mock.Anything)
		ts.SuiMock.AssertNotCalled(t, "SuiGetTransactionBlock", mock.Anything, mock.Anything)
	})

	t.Run("ProcessCCTX restricted address", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		const zetaHeight = 1000

		// Given cctx
		nonce := uint64(123)
		amount := math.NewUint(100_000)
		receiver := "0xdecb47015beebed053c19ef48fe4d722fa3870f567133d235ebe3a70da7b0000"

		cctx := sample.CrossChainTxV2(t, "0xABC123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams = []*cc.OutboundParams{{
			Receiver:        receiver,
			ReceiverChainId: ts.Chain.ChainId,
			CoinType:        coin.CoinType_Gas,
			Amount:          amount,
			TssNonce:        nonce,
			GasPrice:        "1000",
			CallOptions: &cc.CallOptions{
				GasLimit: 42,
			},
		}}

		// Given compliance config
		cfg := config.Config{
			ComplianceConfig: config.ComplianceConfig{
				RestrictedAddresses: []string{receiver},
			},
		}
		config.SetRestrictedAddressesFromConfig(cfg)

		// Given mocked gateway nonce
		ts.MockGatewayNonce(nonce)

		// Given mocked WithdrawCapID
		const withdrawCapID = "0xWithdrawCapID"
		ts.MockWithdrawCapID(withdrawCapID)

		// Given expected MoveCall
		txBytes := base64.StdEncoding.EncodeToString([]byte("raw_tx_bytes"))

		ts.MockMoveCall(func(req models.MoveCallRequest) {
			require.Equal(t, ts.TSS.PubKey().AddressSui(), req.Signer)
			require.Equal(t, ts.Gateway.PackageID(), req.PackageObjectId)
			require.Equal(t, "increase_nonce", req.Function)

			expectedArgs := []any{
				ts.Gateway.ObjectID(),
				fmt.Sprintf("%d", nonce),
				"42000", // gasRefund = GasPrice (1000) * GasLimit (42)
				withdrawCapID,
			}
			require.Equal(t, expectedArgs, req.Arguments)
		}, txBytes)

		// Given expected SuiExecuteTransactionBlock
		const digest = "0xTransactionBlockDigest"
		ts.MockExec(func(req models.SuiExecuteTransactionBlockRequest) {
			require.Equal(t, txBytes, req.TxBytes)
			require.NotEmpty(t, req.Signature)
		}, digest)

		// Given included tx from Sui RPC
		ts.SuiMock.
			On("SuiGetTransactionBlock", mock.Anything, mock.Anything).
			Return(models.SuiTransactionBlockResponse{
				Digest: digest,
				Effects: models.SuiEffects{
					Status: models.ExecutionStatus{
						Status: client.TxStatusSuccess,
					},
				},
				Checkpoint: "1000000",
			}, nil)

		// ACT
		err := ts.Signer.ProcessCCTX(ts.Ctx, cctx, zetaHeight)

		// ASSERT
		require.NoError(t, err)

		// Wait for vote posting
		wait := func() bool {
			if len(ts.TrackerBag) == 0 {
				return false
			}

			vote := ts.TrackerBag[0]
			return vote.hash == digest && vote.nonce == nonce
		}

		require.Eventually(t, wait, 5*time.Second, 100*time.Millisecond)
	})

	t.Run("ProcessCCTX invalid receiver address", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		const zetaHeight = 1000

		// Given cctx
		nonce := uint64(123)
		amount := math.NewUint(100_000)

		// Given invalid receiver address, it's a EVM address
		receiver := "0x547a07f0564e0c8d48c4ae53305eabdef87e9610"

		cctx := sample.CrossChainTxV2(t, "0xABC123")
		cctx.InboundParams.CoinType = coin.CoinType_Gas
		cctx.OutboundParams = []*cc.OutboundParams{{
			Receiver:        receiver,
			ReceiverChainId: ts.Chain.ChainId,
			CoinType:        coin.CoinType_Gas,
			Amount:          amount,
			TssNonce:        nonce,
			GasPrice:        "1000",
			CallOptions: &cc.CallOptions{
				GasLimit: 42,
			},
		}}

		// Given mocked gateway nonce
		ts.MockGatewayNonce(nonce)

		// Given mocked WithdrawCapID
		const withdrawCapID = "0xWithdrawCapID"
		ts.MockWithdrawCapID(withdrawCapID)

		// Given expected MoveCall
		txBytes := base64.StdEncoding.EncodeToString([]byte("raw_tx_bytes"))

		ts.MockMoveCall(func(req models.MoveCallRequest) {
			require.Equal(t, ts.TSS.PubKey().AddressSui(), req.Signer)
			require.Equal(t, ts.Gateway.PackageID(), req.PackageObjectId)
			require.Equal(t, "increase_nonce", req.Function)

			expectedArgs := []any{
				ts.Gateway.ObjectID(),
				fmt.Sprintf("%d", nonce),
				"42000", // gasRefund = GasPrice (1000) * GasLimit (42)
				withdrawCapID,
			}
			require.Equal(t, expectedArgs, req.Arguments)
		}, txBytes)

		// Given expected SuiExecuteTransactionBlock
		const digest = "0xTransactionBlockDigest"
		ts.MockExec(func(req models.SuiExecuteTransactionBlockRequest) {
			require.Equal(t, txBytes, req.TxBytes)
			require.NotEmpty(t, req.Signature)
		}, digest)

		// Given included tx from Sui RPC
		ts.SuiMock.
			On("SuiGetTransactionBlock", mock.Anything, mock.Anything).
			Return(models.SuiTransactionBlockResponse{
				Digest: digest,
				Effects: models.SuiEffects{
					Status: models.ExecutionStatus{
						Status: client.TxStatusSuccess,
					},
				},
				Checkpoint: "1000000",
			}, nil)

		// ACT
		err := ts.Signer.ProcessCCTX(ts.Ctx, cctx, zetaHeight)

		// ASSERT
		require.NoError(t, err)

		// Wait for vote posting
		wait := func() bool {
			if len(ts.TrackerBag) == 0 {
				return false
			}

			vote := ts.TrackerBag[0]
			return vote.hash == digest && vote.nonce == nonce
		}

		require.Eventually(t, wait, 5*time.Second, 100*time.Millisecond)
	})
}

type testSuite struct {
	t   *testing.T
	Ctx context.Context

	Chain chains.Chain

	TSS      *mocks.TSS
	Zetacore *mocks.ZetacoreClient
	SuiMock  *mocks.SuiClient
	Gateway  *sui.Gateway

	*Signer

	TrackerBag []testTracker
}

type testSuiteConfig struct {
	withdrawCapID     string
	previousPackageID string
	originalPackageID string
}

func newTestSuite(t *testing.T, opts ...func(*testSuiteConfig)) *testSuite {
	var cfg testSuiteConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	var (
		ctx = context.Background()

		chain       = chains.SuiMainnet
		chainParams = mocks.MockChainParams(chain.ChainId, 10)

		tss      = mocks.NewTSS(t)
		zetacore = mocks.NewZetacoreClient(t).WithKeys(&keys.Keys{})

		testLogger = testlog.New(t)
		logger     = base.Logger{Std: testLogger.Logger, Compliance: testLogger.Logger}
	)

	// append withdraw cap ID, previous package ID and original package ID if provided
	if cfg.withdrawCapID != "" && cfg.previousPackageID != "" && cfg.originalPackageID != "" {
		chainParams.GatewayAddress = fmt.Sprintf(
			"%s,%s,%s,%s",
			chainParams.GatewayAddress,
			cfg.withdrawCapID,
			cfg.previousPackageID,
			cfg.originalPackageID,
		)
	}

	suiMock := mocks.NewSuiClient(t)

	gw, err := sui.NewGatewayFromPairID(chainParams.GatewayAddress)
	require.NoError(t, err)

	baseSigner := base.NewSigner(chain, tss, logger, mode.StandardMode)
	zetaRepo := zrepo.New(zetacore, chain, mode.StandardMode)
	signer := New(baseSigner, zetaRepo, suiMock, gw)

	ts := &testSuite{
		t:        t,
		Ctx:      ctx,
		Chain:    chain,
		TSS:      tss,
		Zetacore: zetacore,
		SuiMock:  suiMock,
		Gateway:  gw,
		Signer:   signer,
	}

	// Setup mocks
	ts.Zetacore.On("Chain").Return(chain).Maybe()

	ts.setupTrackersBag()

	return ts
}

func (ts *testSuite) MockGatewayNonce(nonce uint64) {
	ts.SuiMock.On("GetObjectParsedData", mock.Anything, mock.Anything).Return(models.SuiParsedData{
		SuiMoveObject: models.SuiMoveObject{
			Fields: map[string]any{"nonce": fmt.Sprintf("%d", nonce)},
		},
	}, nil)
}

func (ts *testSuite) MockWithdrawCapID(id string) {
	tss, structType := ts.TSS.PubKey().AddressSui(), ts.Gateway.WithdrawCapType()
	ts.SuiMock.On("GetOwnedObjectID", mock.Anything, tss, structType).Return(id, nil)
}

func (ts *testSuite) MockMessageContextID(id string) {
	ts.SuiMock.On("SuiXGetDynamicFieldObject", mock.Anything, mock.Anything).Maybe().Return(models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]any{"value": id},
				},
			},
		},
	}, nil)
}

func (ts *testSuite) MockMoveCall(assert func(req models.MoveCallRequest), txBytesBase64 string) {
	call := func(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error) {
		assert(req)
		return models.TxnMetaData{TxBytes: txBytesBase64}, nil
	}

	ts.SuiMock.On("MoveCall", mock.Anything, mock.Anything).Return(call)
}

func (ts *testSuite) MockExec(assert func(req models.SuiExecuteTransactionBlockRequest), digest string) {
	call := func(
		ctx context.Context,
		req models.SuiExecuteTransactionBlockRequest,
	) (models.SuiTransactionBlockResponse, error) {
		assert(req)
		return models.SuiTransactionBlockResponse{
			Effects: models.SuiEffects{
				Status: models.ExecutionStatus{
					Status: client.TxStatusSuccess,
				},
			},
			Digest: digest,
		}, nil
	}

	ts.SuiMock.On("SuiExecuteTransactionBlock", mock.Anything, mock.Anything).Return(call)
}

type testTracker struct {
	nonce uint64
	hash  string
}

func (ts *testSuite) setupTrackersBag() {
	catcher := func(args mock.Arguments) {
		require.Equal(ts.t, ts.Chain.ChainId, args.Get(1).(int64))
		nonce := args.Get(2).(uint64)
		txHash := args.Get(3).(string)

		ts.t.Logf("Adding outbound tracker: nonce=%d, hash=%s", nonce, txHash)

		ts.TrackerBag = append(ts.TrackerBag, testTracker{nonce, txHash})
	}

	ts.Zetacore.On(
		"PostOutboundTracker",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Maybe().Run(catcher).Return("", nil)
}
