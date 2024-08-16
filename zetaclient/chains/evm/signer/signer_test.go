package signer

import (
	"context"
	"math/big"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/db"
	"github.com/zeta-chain/zetacore/zetaclient/keys"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outboundprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

var (
	// Dummy addresses as they are just used as transaction data to be signed
	ConnectorAddress    = sample.EthAddress()
	ERC20CustodyAddress = sample.EthAddress()
)

// getNewEvmSigner creates a new EVM chain signer for testing
func getNewEvmSigner(tss interfaces.TSSSigner) (*Signer, error) {
	ctx := context.Background()

	// use default mock TSS if not provided
	if tss == nil {
		tss = mocks.NewTSSMainnet()
	}

	mpiAddress := ConnectorAddress
	erc20CustodyAddress := ERC20CustodyAddress
	logger := base.Logger{}

	return NewSigner(
		ctx,
		chains.BscMainnet,
		tss,
		nil,
		logger,
		mocks.EVMRPCEnabled,
		config.GetConnectorABI(),
		config.GetERC20CustodyABI(),
		mpiAddress,
		erc20CustodyAddress,
	)
}

// getNewEvmChainObserver creates a new EVM chain observer for testing
func getNewEvmChainObserver(t *testing.T, tss interfaces.TSSSigner) (*observer.Observer, error) {
	ctx := context.Background()

	// use default mock TSS if not provided
	if tss == nil {
		tss = mocks.NewTSSMainnet()
	}
	cfg := config.New(false)

	// prepare mock arguments to create observer
	evmcfg := config.EVMConfig{Chain: chains.BscMainnet, Endpoint: "http://localhost:8545"}
	evmClient := mocks.NewMockEvmClient().WithBlockNumber(1000)
	params := mocks.MockChainParams(evmcfg.Chain.ChainId, 10)
	cfg.EVMChainConfigs[chains.BscMainnet.ChainId] = evmcfg
	//appContext := context.New(cfg, zerolog.Nop())
	logger := base.Logger{}
	ts := &metrics.TelemetryServer{}

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	return observer.NewObserver(
		ctx,
		evmcfg,
		evmClient,
		params,
		mocks.NewZetacoreClient(t),
		tss,
		database,
		logger,
		ts,
	)
}

func getNewOutboundProcessor() *outboundprocessor.Processor {
	logger := zerolog.Logger{}
	return outboundprocessor.NewProcessor(logger)
}

func getCCTX(t *testing.T) *crosschaintypes.CrossChainTx {
	return testutils.LoadCctxByNonce(t, 56, 68270)
}

func getInvalidCCTX(t *testing.T) *crosschaintypes.CrossChainTx {
	cctx := getCCTX(t)
	// modify receiver chain id to make it invalid
	cctx.GetCurrentOutboundParam().ReceiverChainId = 13378337
	return cctx
}

// verifyTxSignature is a helper function to verify the signature of a transaction
func verifyTxSignature(t *testing.T, tx *ethtypes.Transaction, tssPubkey []byte, signer ethtypes.Signer) {
	_, r, s := tx.RawSignatureValues()
	signature := append(r.Bytes(), s.Bytes()...)
	hash := signer.Hash(tx)

	verified := crypto.VerifySignature(tssPubkey, hash.Bytes(), signature)
	require.True(t, verified)
}

// verifyTxBodyBasics is a helper function to verify 'to', 'nonce' and 'amount' of a transaction
func verifyTxBodyBasics(
	t *testing.T,
	tx *ethtypes.Transaction,
	to ethcommon.Address,
	nonce uint64,
	amount *big.Int,
) {
	require.Equal(t, to, *tx.To())
	require.Equal(t, nonce, tx.Nonce())
	require.True(t, amount.Cmp(tx.Value()) == 0)
}

func TestSigner_SetGetConnectorAddress(t *testing.T) {
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)
	// Get and compare
	require.Equal(t, ConnectorAddress, evmSigner.GetZetaConnectorAddress())

	// Update and get again
	newConnector := sample.EthAddress()
	evmSigner.SetZetaConnectorAddress(newConnector)
	require.Equal(t, newConnector, evmSigner.GetZetaConnectorAddress())
}

func TestSigner_SetGetERC20CustodyAddress(t *testing.T) {
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)
	// Get and compare
	require.Equal(t, ERC20CustodyAddress, evmSigner.GetERC20CustodyAddress())

	// Update and get again
	newCustody := sample.EthAddress()
	evmSigner.SetERC20CustodyAddress(newCustody)
	require.Equal(t, newCustody, evmSigner.GetERC20CustodyAddress())
}

func TestSigner_TryProcessOutbound(t *testing.T) {
	ctx := makeCtx(t)

	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)
	cctx := getCCTX(t)
	processor := getNewOutboundProcessor()
	mockObserver, err := getNewEvmChainObserver(t, nil)
	require.NoError(t, err)

	// Test with mock client that has keys
	client := mocks.NewZetacoreClient(t).
		WithKeys(&keys.Keys{}).
		WithZetaChain().
		WithPostVoteOutbound("", "")

	evmSigner.TryProcessOutbound(ctx, cctx, processor, "123", mockObserver, client, 123)

	// Check if cctx was signed and broadcasted
	list := evmSigner.GetReportedTxList()
	require.Len(t, *list, 1)
}

func TestSigner_SignOutbound(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct

	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignOutbound - should successfully sign LegacyTx", func(t *testing.T) {
		// Call SignOutbound
		tx, err := evmSigner.SignOutbound(ctx, txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// check that by default tx type is legacy tx
		assert.Equal(t, ethtypes.LegacyTxType, int(tx.Type()))
	})
	t.Run("SignOutbound - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignOutbound
		tx, err := evmSigner.SignOutbound(ctx, txData)
		require.ErrorContains(t, err, "sign onReceive error")
		require.Nil(t, tx)
	})

	t.Run("SignOutbound - should successfully sign DynamicFeeTx", func(t *testing.T) {
		// ARRANGE
		const (
			gwei        = 1_000_000_000
			priorityFee = 1 * gwei
			gasPrice    = 3 * gwei
		)

		// Given a CCTX with gas price and priority fee
		cctx := getCCTX(t)
		cctx.OutboundParams[0].GasPrice = big.NewInt(gasPrice).String()
		cctx.OutboundParams[0].GasPriorityFee = big.NewInt(priorityFee).String()

		// Given outbound data
		txData, skip, err := NewOutboundData(ctx, cctx, 123, makeLogger(t))
		require.False(t, skip)
		require.NoError(t, err)

		// Given a working TSS
		tss.Unpause()

		// ACT
		tx, err := evmSigner.SignOutbound(ctx, txData)
		require.NoError(t, err)

		// ASSERT
		verifyTxSignature(t, tx, mocks.NewTSSMainnet().Pubkey(), evmSigner.EvmSigner())

		// check that by default tx type is a dynamic fee tx
		assert.Equal(t, ethtypes.DynamicFeeTxType, int(tx.Type()))

		// check that the gasPrice & priorityFee are set correctly
		assert.Equal(t, int64(gasPrice), tx.GasFeeCap().Int64())
		assert.Equal(t, int64(priorityFee), tx.GasTipCap().Int64())
	})
}

func TestSigner_SignRevertTx(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignRevertTx - should successfully sign", func(t *testing.T) {
		// Call SignRevertTx
		tx, err := evmSigner.SignRevertTx(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls connector contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.zetaConnectorAddress, txData.nonce, big.NewInt(0))
	})
	t.Run("SignRevertTx - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignRevertTx
		tx, err := evmSigner.SignRevertTx(ctx, txData)
		require.ErrorContains(t, err, "sign onRevert error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignCancelTx(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCancelTx - should successfully sign", func(t *testing.T) {
		// Call SignRevertTx
		tx, err := evmSigner.SignCancelTx(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Cancel tx sends 0 gas token to TSS self address
		verifyTxBodyBasics(t, tx, evmSigner.TSS().EVMAddress(), txData.nonce, big.NewInt(0))
	})
	t.Run("SignCancelTx - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignCancelTx
		tx, err := evmSigner.SignCancelTx(ctx, txData)
		require.ErrorContains(t, err, "SignCancelTx error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignWithdrawTx(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignWithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignWithdrawTx
		tx, err := evmSigner.SignWithdrawTx(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})
	t.Run("SignWithdrawTx - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignWithdrawTx
		tx, err := evmSigner.SignWithdrawTx(ctx, txData)
		require.ErrorContains(t, err, "SignWithdrawTx error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignERC20WithdrawTx(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignERC20WithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(ctx, txData)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Withdraw tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, evmSigner.er20CustodyAddress, txData.nonce, big.NewInt(0))
	})

	t.Run("SignERC20WithdrawTx - should fail if keysign fails", func(t *testing.T) {
		// pause tss to make keysign fail
		tss.Pause()

		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(ctx, txData)
		require.ErrorContains(t, err, "sign withdraw error")
		require.Nil(t, tx)
	})
}

func TestSigner_BroadcastOutbound(t *testing.T) {
	ctx := makeCtx(t)

	// Setup evm signer
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	txData, skip, err := NewOutboundData(ctx, cctx, 123, zerolog.Logger{})
	require.NoError(t, err)
	require.False(t, skip)

	t.Run("BroadcastOutbound - should successfully broadcast", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(ctx, txData)
		require.NoError(t, err)

		evmSigner.BroadcastOutbound(
			ctx,
			tx,
			cctx,
			zerolog.Logger{},
			sdktypes.AccAddress{},
			mocks.NewZetacoreClient(t),
			txData,
		)

		//Check if cctx was signed and broadcasted
		list := evmSigner.GetReportedTxList()
		require.Len(t, *list, 1)
	})
}

func TestSigner_getEVMRPC(t *testing.T) {
	ctx := context.Background()

	t.Run("getEVMRPC error dialing", func(t *testing.T) {
		client, signer, err := getEVMRPC(ctx, "invalidEndpoint")
		require.Nil(t, client)
		require.Nil(t, signer)
		require.Error(t, err)
	})
}

func TestSigner_SignerErrorMsg(t *testing.T) {
	cctx := getCCTX(t)

	msg := ErrorMsg(cctx)
	require.Contains(t, msg, "nonce 68270 chain 56")
}

func makeCtx(t *testing.T) context.Context {
	app := zctx.New(config.New(false), nil, zerolog.Nop())

	bscParams := mocks.MockChainParams(chains.BscMainnet.ChainId, 10)

	err := app.Update(
		observertypes.Keygen{},
		[]chains.Chain{chains.BscMainnet, chains.ZetaChainMainnet},
		nil,
		map[int64]*observertypes.ChainParams{
			chains.BscMainnet.ChainId: &bscParams,
		},
		"tssPubKey",
		observertypes.CrosschainFlags{},
	)
	require.NoError(t, err, "unable to update app context")

	return zctx.WithAppContext(context.Background(), app)
}

func makeLogger(t *testing.T) zerolog.Logger {
	return zerolog.New(zerolog.NewTestWriter(t))
}
