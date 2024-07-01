package signer

import (
	"math/big"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
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
	// use default mock TSS if not provided
	if tss == nil {
		tss = mocks.NewTSSMainnet()
	}

	mpiAddress := ConnectorAddress
	erc20CustodyAddress := ERC20CustodyAddress
	logger := base.Logger{}
	cfg := config.NewConfig()

	return NewSigner(
		chains.BscMainnet,
		context.NewAppContext(cfg),
		tss,
		nil,
		logger,
		mocks.EVMRPCEnabled,
		config.GetConnectorABI(),
		config.GetERC20CustodyABI(),
		mpiAddress,
		erc20CustodyAddress)
}

// getNewEvmChainObserver creates a new EVM chain observer for testing
func getNewEvmChainObserver(t *testing.T, tss interfaces.TSSSigner) (*observer.Observer, error) {
	// use default mock TSS if not provided
	if tss == nil {
		tss = mocks.NewTSSMainnet()
	}
	cfg := config.NewConfig()

	// prepare mock arguments to create observer
	evmcfg := config.EVMConfig{Chain: chains.BscMainnet, Endpoint: "http://localhost:8545"}
	evmClient := mocks.NewMockEvmClient().WithBlockNumber(1000)
	params := mocks.MockChainParams(evmcfg.Chain.ChainId, 10)
	cfg.EVMChainConfigs[chains.BscMainnet.ChainId] = evmcfg
	appCTX := context.NewAppContext(cfg)
	logger := base.Logger{}
	ts := &metrics.TelemetryServer{}

	return observer.NewObserver(
		evmcfg,
		evmClient,
		params,
		appCTX,
		mocks.NewMockZetacoreClient(),
		tss,
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
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)
	cctx := getCCTX(t)
	processor := getNewOutboundProcessor()
	mockObserver, err := getNewEvmChainObserver(t, nil)
	require.NoError(t, err)

	// Test with mock client that has keys
	client := mocks.NewMockZetacoreClient().WithKeys(&keys.Keys{})
	evmSigner.TryProcessOutbound(cctx, processor, "123", mockObserver, client, 123)

	// Check if cctx was signed and broadcasted
	list := evmSigner.GetReportedTxList()
	require.Len(t, *list, 1)
}

func TestSigner_SignOutbound(t *testing.T) {
	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct

	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignOutbound - should successfully sign", func(t *testing.T) {
		// Call SignOutbound
		tx, err := evmSigner.SignOutbound(txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())
	})
	t.Run("SignOutbound - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignOutbound
		tx, err := evmSigner.SignOutbound(txData)
		require.ErrorContains(t, err, "sign onReceive error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignRevertTx(t *testing.T) {
	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignRevertTx - should successfully sign", func(t *testing.T) {
		// Call SignRevertTx
		tx, err := evmSigner.SignRevertTx(txData)
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
		tx, err := evmSigner.SignRevertTx(txData)
		require.ErrorContains(t, err, "sign onRevert error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignCancelTx(t *testing.T) {
	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCancelTx - should successfully sign", func(t *testing.T) {
		// Call SignRevertTx
		tx, err := evmSigner.SignCancelTx(txData)
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
		tx, err := evmSigner.SignCancelTx(txData)
		require.ErrorContains(t, err, "SignCancelTx error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignWithdrawTx(t *testing.T) {
	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignWithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignWithdrawTx
		tx, err := evmSigner.SignWithdrawTx(txData)
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
		tx, err := evmSigner.SignWithdrawTx(txData)
		require.ErrorContains(t, err, "SignWithdrawTx error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignCommandTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, nil)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCommandTx CmdWhitelistERC20", func(t *testing.T) {
		cmd := constant.CmdWhitelistERC20
		params := ConnectorAddress.Hex()
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(txData, cmd, params)
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		// Note: Revert tx calls erc20 custody contract with 0 gas token
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, big.NewInt(0))
	})

	t.Run("SignCommandTx CmdMigrateTssFunds", func(t *testing.T) {
		cmd := constant.CmdMigrateTssFunds
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(txData, cmd, "")
		require.NoError(t, err)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})
}

func TestSigner_SignERC20WithdrawTx(t *testing.T) {
	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignERC20WithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(txData)
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
		tx, err := evmSigner.SignERC20WithdrawTx(txData)
		require.ErrorContains(t, err, "sign withdraw error")
		require.Nil(t, tx)
	})
}

func TestSigner_BroadcastOutbound(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner(nil)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, nil)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("BroadcastOutbound - should successfully broadcast", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(txData)
		require.NoError(t, err)

		evmSigner.BroadcastOutbound(
			tx,
			cctx,
			zerolog.Logger{},
			sdktypes.AccAddress{},
			mocks.NewMockZetacoreClient(),
			txData,
		)

		//Check if cctx was signed and broadcasted
		list := evmSigner.GetReportedTxList()
		require.Len(t, *list, 1)
	})
}

func TestSigner_getEVMRPC(t *testing.T) {
	t.Run("getEVMRPC error dialing", func(t *testing.T) {
		client, signer, err := getEVMRPC("invalidEndpoint")
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

func TestSigner_SignWhitelistERC20Cmd(t *testing.T) {
	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignWhitelistERC20Cmd - should successfully sign", func(t *testing.T) {
		// Call SignWhitelistERC20Cmd
		tx, err := evmSigner.SignWhitelistERC20Cmd(txData, sample.EthAddress().Hex())
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, zeroValue)
	})
	t.Run("SignWhitelistERC20Cmd - should fail on invalid erc20 address", func(t *testing.T) {
		tx, err := evmSigner.SignWhitelistERC20Cmd(txData, "")
		require.Nil(t, tx)
		require.ErrorContains(t, err, "invalid erc20 address")
	})
	t.Run("SignWhitelistERC20Cmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignWhitelistERC20Cmd
		tx, err := evmSigner.SignWhitelistERC20Cmd(txData, sample.EthAddress().Hex())
		require.ErrorContains(t, err, "sign whitelist error")
		require.Nil(t, tx)
	})
}

func TestSigner_SignMigrateTssFundsCmd(t *testing.T) {
	// Setup evm signer
	tss := mocks.NewTSSMainnet()
	evmSigner, err := getNewEvmSigner(tss)
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver(t, tss)
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignMigrateTssFundsCmd - should successfully sign", func(t *testing.T) {
		// Call SignMigrateTssFundsCmd
		tx, err := evmSigner.SignMigrateTssFundsCmd(txData)
		require.NoError(t, err)
		require.NotNil(t, tx)

		// Verify tx signature
		tss := mocks.NewTSSMainnet()
		verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

		// Verify tx body basics
		verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
	})

	t.Run("SignMigrateTssFundsCmd - should fail if keysign fails", func(t *testing.T) {
		// Pause tss to make keysign fail
		tss.Pause()

		// Call SignMigrateTssFundsCmd
		tx, err := evmSigner.SignMigrateTssFundsCmd(txData)
		require.ErrorContains(t, err, "SignMigrateTssFundsCmd error")
		require.Nil(t, tx)
	})
}
