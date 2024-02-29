package evm

import (
	"path"
	"testing"

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
	evmcfg := config.EVMConfig{Endpoint: "http://localhost:8545"}
	return NewEVMChainClient(mock.NewZetaCoreBridge(), mock.NewTSSMainnet(), "", logger, cfg, evmcfg, ts)
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

	//Check if cctx was processed
	list := evmSigner.GetReportedTxList()
	found := false
	for range *list {
		found = true
	}
	require.True(t, found)
}

func TestSigner_SignOutboundTx(t *testing.T) {
}

func TestSigner_SignRevertTx(t *testing.T) {
}

func TestSigner_SignWithdrawTx(t *testing.T) {
}

func TestSigner_SignCommandTx(t *testing.T) {
}

func TestSigner_SignERC20WithdrawTx(t *testing.T) {
}

func TestSigner_BroadcastOutTx(t *testing.T) {
}

func TestSigner_SetChainAndSender(t *testing.T) {
}

func TestSigner_SetupGas(t *testing.T) {
}

func TestSigner_SetTransactionData(t *testing.T) {
}
