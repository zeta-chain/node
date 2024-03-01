package evm

import (
	"math/big"
	"path"
	"testing"

	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"

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
	return NewEVMSigner(
		corecommon.BscMainnetChain(),
		stub.EVMRPCEnabled,
		stub.NewTSSMainnet(),
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
	tss := stub.NewTSSMainnet()

	evmcfg := config.EVMConfig{Chain: corecommon.BscMainnetChain(), Endpoint: "http://localhost:8545"}
	cfg.EVMChainConfigs[corecommon.BscMainnetChain().ChainId] = &evmcfg
	coreCTX := corecontext.NewZetaCoreContext(cfg)
	appCTX := appcontext.NewAppContext(coreCTX, cfg)

	return NewEVMChainClient(appCTX, stub.NewZetaCoreBridge(), tss, "", logger, evmcfg, ts)
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

	evmSigner.TryProcessOutTx(cctx, processorManager, "123", mockChainClient, stub.NewZetaCoreBridge(), 123)

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
	txData := BaseTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignOutboundTx - should successfully sign", func(t *testing.T) {
		// Call SignOutboundTx
		tx, err := evmSigner.SignOutboundTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
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
	txData := BaseTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignRevertTx - should successfully sign", func(t *testing.T) {
		// Call SignRevertTx
		tx, err := evmSigner.SignRevertTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
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
	txData := BaseTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignWithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignWithdrawTx
		tx, err := evmSigner.SignWithdrawTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
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
	txData := BaseTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignCommandTx CmdWhitelistERC20", func(t *testing.T) {
		txData.cmd = corecommon.CmdWhitelistERC20
		txData.params = ConnectorAddress.Hex()
		// Call SignCommandTx
		tx, err := evmSigner.SignCommandTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
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
		tss := stub.NewTSSMainnet()
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
	txData := BaseTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("SignERC20WithdrawTx - should successfully sign", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(&txData)
		require.NoError(t, err)

		// Verify Signature
		tss := stub.NewTSSMainnet()
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
	txData := BaseTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)

	t.Run("BroadcastOutTx - should successfully broadcast", func(t *testing.T) {
		// Call SignERC20WithdrawTx
		tx, err := evmSigner.SignERC20WithdrawTx(&txData)
		require.NoError(t, err)

		evmSigner.BroadcastOutTx(tx, cctx, zerolog.Logger{}, sdktypes.AccAddress{}, stub.NewZetaCoreBridge(), &txData)

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

	txData := &BaseTransactionData{}
	logger := zerolog.Logger{}

	t.Run("SetChainAndSender PendingRevert", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		skipTx := txData.SetChainAndSender(cctx, logger)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.InboundTxParams.Sender), txData.to)
		require.Equal(t, big.NewInt(cctx.InboundTxParams.SenderChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender PendingOutBound", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		skipTx := txData.SetChainAndSender(cctx, logger)

		require.False(t, skipTx)
		require.Equal(t, ethcommon.HexToAddress(cctx.GetCurrentOutTxParam().Receiver), txData.to)
		require.Equal(t, big.NewInt(cctx.GetCurrentOutTxParam().ReceiverChainId), txData.toChainID)
	})

	t.Run("SetChainAndSender Should skip cctx", func(t *testing.T) {
		cctx.CctxStatus.Status = types.CctxStatus_PendingInbound
		skipTx := txData.SetChainAndSender(cctx, logger)
		require.True(t, skipTx)
	})
}

func TestSigner_SetupGas(t *testing.T) {
	cctx, err := getCCTX()
	require.NoError(t, err)

	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	txData := &BaseTransactionData{}
	logger := zerolog.Logger{}

	t.Run("SetupGas_success", func(t *testing.T) {
		chain := corecommon.BscMainnetChain()
		err := txData.SetupGas(cctx, logger, evmSigner.EvmClient(), &chain)
		require.NoError(t, err)
	})

	t.Run("SetupGas_error", func(t *testing.T) {
		cctx.GetCurrentOutTxParam().OutboundTxGasPrice = "invalidGasPrice"
		chain := corecommon.BtcMainnetChain()
		err := txData.SetupGas(cctx, logger, evmSigner.EvmClient(), &chain)
		require.Error(t, err)
	})
}

func TestSigner_SetTransactionData(t *testing.T) {
	// Setup evm signer
	evmSigner, err := getNewEvmSigner()
	require.NoError(t, err)

	// Setup txData struct
	txData := BaseTransactionData{}
	cctx, err := getCCTX()
	require.NoError(t, err)
	mockChainClient, err := getNewEvmChainClient()
	require.NoError(t, err)
	skip, err := txData.SetTransactionData(cctx, mockChainClient, evmSigner.EvmClient(), zerolog.Logger{})
	require.False(t, skip)
	require.NoError(t, err)
}
