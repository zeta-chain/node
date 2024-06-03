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
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/common"
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

func getNewEvmSigner() (*Signer, error) {
	mpiAddress := ConnectorAddress
	erc20CustodyAddress := ERC20CustodyAddress
	logger := common.ClientLogger{}
	ts := &metrics.TelemetryServer{}
	cfg := config.NewConfig()
	return NewSigner(
		chains.BscMainnet,
		mocks.EVMRPCEnabled,
		mocks.NewTSSMainnet(),
		config.GetConnectorABI(),
		config.GetERC20CustodyABI(),
		mpiAddress,
		erc20CustodyAddress,
		context.NewZetacoreContext(cfg),
		logger,
		ts)
}

func getNewEvmChainObserver() (*observer.Observer, error) {
	logger := common.ClientLogger{}
	ts := &metrics.TelemetryServer{}
	cfg := config.NewConfig()
	tss := mocks.NewTSSMainnet()

	evmcfg := config.EVMConfig{Chain: chains.BscMainnet, Endpoint: "http://localhost:8545"}
	cfg.EVMChainConfigs[chains.BscMainnet.ChainId] = evmcfg
	coreCTX := context.NewZetacoreContext(cfg)
	appCTX := context.NewAppContext(coreCTX, cfg)

	return observer.NewObserver(appCTX, mocks.NewMockZetacoreClient(), tss, "", logger, evmcfg, ts)
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
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)
	// Get and compare
	require.Equal(t, ConnectorAddress, evmSigner.GetZetaConnectorAddress())

	// Update and get again
	newConnector := sample.EthAddress()
	evmSigner.SetZetaConnectorAddress(newConnector)
	require.Equal(t, newConnector, evmSigner.GetZetaConnectorAddress())
}

func TestSigner_SetGetERC20CustodyAddress(t *testing.T) {
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)
	// Get and compare
	require.Equal(t, ERC20CustodyAddress, evmSigner.GetERC20CustodyAddress())

	// Update and get again
	newCustody := sample.EthAddress()
	evmSigner.SetERC20CustodyAddress(newCustody)
	require.Equal(t, newCustody, evmSigner.GetERC20CustodyAddress())
}

func TestSigner_TryProcessOutbound(t *testing.T) {
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)
	cctx := getCCTX(t)
	processorManager := getNewOutboundProcessor()
	mockObserver, err := getNewEvmChainObserver()
	require.NoError(t, err)

	// Test with mock client that has keys
	client := mocks.NewMockZetacoreClient().WithKeys(&keys.Keys{})
	evmSigner.TryProcessOutbound(cctx, processorManager, "123", mockObserver, client, 123)

	// Check if cctx was signed and broadcasted
	list := evmSigner.GetReportedTxList()
	require.Len(t, *list, 1)
}

func TestSigner_SignOutbound(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct

	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
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
}

func TestSigner_SignRevertTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
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
}

func TestSigner_SignCancelTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
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
		verifyTxBodyBasics(t, tx, evmSigner.tssSigner.EVMAddress(), txData.nonce, big.NewInt(0))
	})
}

func TestSigner_SignWithdrawTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
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
}

func TestSigner_SignCommandTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
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
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
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
}

func TestSigner_BroadcastOutbound(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
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
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	tx, err := evmSigner.SignWhitelistERC20Cmd(txData, "")
	require.Nil(t, tx)
	require.ErrorContains(t, err, "invalid erc20 address")
}

func TestSigner_SignMigrateTssFundsCmd(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	cctx := getCCTX(t)
	mockObserver, err := getNewEvmChainObserver()
	require.NoError(t, err)
	txData, skip, err := NewOutboundData(cctx, mockObserver, evmSigner.EvmClient(), zerolog.Logger{}, 123)
	require.False(t, skip)
	require.NoError(t, err)

	tx, err := evmSigner.SignMigrateTssFundsCmd(txData)
	require.NoError(t, err)
	require.NotNil(t, tx)

	// Verify tx signature
	tss := mocks.NewTSSMainnet()
	verifyTxSignature(t, tx, tss.Pubkey(), evmSigner.EvmSigner())

	// Verify tx body basics
	verifyTxBodyBasics(t, tx, txData.to, txData.nonce, txData.amount)
}
