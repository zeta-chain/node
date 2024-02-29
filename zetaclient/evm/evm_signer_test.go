package evm

import (
	"math/big"
	"path"
	"testing"

	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	corecommon "github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/outtxprocessor"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mock"
)

const (
	// Dummy addresses as they are just used as transaction data to be signed
	ConnectorAddress    = "0x00000000219ab540356cbb839cbe05303d7705fa"
	ERC20CustodyAddress = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
)

func getNewEvmSigner() (*Signer, error) {
	mpiAddress := ethcommon.HexToAddress(ConnectorAddress)
	erc20CustodyAddress := ethcommon.HexToAddress(ERC20CustodyAddress)
	logger := common.ClientLogger{}
	ts := &metrics.TelemetryServer{}
	return NewEVMSigner(
		corecommon.BscMainnetChain(),
		mock.EVMRPCEnabled,
		mock.NewTSSMainnet(),
		config.GetConnectorABI(),
		config.GetERC20CustodyABI(),
		mpiAddress,
		erc20CustodyAddress,
		logger,
		ts)
}

func getNewEvmChainClient() (*ChainClient, error) {
	logger := common.ClientLogger{}
	ts := &metrics.TelemetryServer{}
	cfg := config.NewConfig()
	tss := mock.NewTSSMainnet()

	evmcfg := config.EVMConfig{Chain: corecommon.BscMainnetChain(), Endpoint: "http://localhost:8545"}
	cfg.EVMChainConfigs[corecommon.BscMainnetChain().ChainId] = &evmcfg
	coreCTX := corecontext.NewZetaCoreContext(cfg)
	appCTX := appcontext.NewAppContext(coreCTX, cfg)

	return NewEVMChainClient(appCTX, mock.NewZetaCoreBridge(), tss, "", logger, evmcfg, ts)
}

func getNewOutTxProcessor() *outtxprocessor.Processor {
	logger := zerolog.Logger{}
	return outtxprocessor.NewOutTxProcessorManager(logger)
}

func getCCTX() (*types.CrossChainTx, error) {
	var cctx crosschaintypes.CrossChainTx
	err := testutils.LoadObjectFromJSONFile(&cctx, path.Join("../", testutils.TestDataPathCctx, "cctx_56_68270.json"))
	return &cctx, err
}

func TestSigner_TryProcessOutTx(t *testing.T) {
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)
	cctx, err := getCCTX()
	require.NoError(t, err)
	processorManager := getNewOutTxProcessor()
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)

	evmSigner.TryProcessOutTx(cctx, processorManager, "123", mockChainClient, mock.NewZetaCoreBridge(), 123)

	//Check if cctx was signed and broadcasted
	list := evmSigner.GetReportedTxList()
	found := false
	for range *list {
		found = true
	}
	require.True(t, found)
}

func TestSigner_SignOutboundTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := TransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, mock.NewZetaCoreBridge(), &txData)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignOutboundTx - should successfully sign", func(t *testing.T) {
		// Call SignOutboundTx
		tx, err := evmSigner.SignOutboundTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mock.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignRevertTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := TransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, mock.NewZetaCoreBridge(), &txData)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignRevertTx - should successfully sign", func(t *testing.T) {
		// Call SignRevertTx
		tx, err := evmSigner.SignRevertTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mock.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignWithdrawTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := TransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, mock.NewZetaCoreBridge(), &txData)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignWithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignWithdrawTx
		tx, err := evmSigner.SignWithdrawTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mock.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignCommandTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := TransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, mock.NewZetaCoreBridge(), &txData)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCommandTx CmdWhitelistERC20", func(t *testing.T) {
		txData.cmd = corecommon.CmdWhitelistERC20
		txData.params = ConnectorAddress
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mock.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})

	t.Run("SignCommandTx CmdMigrateTssFunds", func(t *testing.T) {
		txData.cmd = corecommon.CmdMigrateTssFunds
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mock.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_SignERC20WithdrawTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := TransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, mock.NewZetaCoreBridge(), &txData)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignERC20WithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := mock.NewTSSMainnet()
		_, r, s := tx.RawSignatureValues()
		signature := append(r.Bytes(), s.Bytes()...)
		hash := evmSigner.EvmSigner().Hash(tx)

		verified := crypto.VerifySignature(tss.Pubkey(), hash.Bytes(), signature)
		require.True(t, verified)
	})
}

func TestSigner_BroadcastOutTx(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := TransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, mock.NewZetaCoreBridge(), &txData)
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("BroadcastOutTx - should successfully broadcast", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(&txData)
		require.NoError(t, err)

		evmSigner.BroadcastOutTx(tx, cctx, zerolog.Logger{}, sdktypes.AccAddress{}, mock.NewZetaCoreBridge(), &txData)

		//Check if cctx was signed and broadcasted
		list := evmSigner.GetReportedTxList()
		found := false
		for range *list {
			found = true
		}
		require.True(t, found)
	})
}

func TestSigner_SetChainAndSender(t *testing.T) {
	// setup inputs
	cctx, err := getCCTX()
	require.NoError(t, err)

	txData := &TransactionData{}
	logger := zerolog.Logger{}

	t.Run("SetChainAndSender PendingRevert", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		skipTx := SetChainAndSender(cctx, logger, txData)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.InboundTxParams.Sender), txData.to)
		require.Equal(t, big.NewInt(cctx.InboundTxParams.SenderChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender PendingOutBound", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		skipTx := SetChainAndSender(cctx, logger, txData)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver), txData.to)
		require.Equal(t, big.NewInt(cctx.GetCurrentOutTxParam().ReceiverChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender Should skip cctx", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingInbound
		skipTx := SetChainAndSender(cctx, logger, txData)
		require.True(t, skipTx)
	})
}

func TestSigner_SetupGas(t *testing.T) {
	cctx, err := getCCTX()
	require.NoError(t, err)

	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	txData := &TransactionData{}
	logger := zerolog.Logger{}

	t.Run("SetupGas_success", func(t *testing.T) {
		chain := corecommon.BscMainnetChain()
		err := SetupGas(cctx, logger, evmSigner.EvmClient(), &chain, txData)
		require.NoError(t, err)
	})

	t.Run("SetupGas_error", func(t *testing.T) {
		cctx.GetCurrentOutTxParam().OutboundTxGasPrice = "invalidGasPrice"
		chain := corecommon.BtcMainnetChain()
		err := SetupGas(cctx, logger, evmSigner.EvmClient(), &chain, txData)
		require.Error(t, err)
	})
}

func TestSigner_SetTransactionData(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := TransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{}, mock.NewZetaCoreBridge(), &txData)
	require.False(t, skip)
	require.NoError(t, err)
}
