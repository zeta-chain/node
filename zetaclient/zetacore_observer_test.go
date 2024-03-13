package zetaclient

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/bitcoin"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/evm"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
)

func MockCoreObserver(t *testing.T, evmChain, btcChain common.Chain, connectorAddress, erc20CustodyAddress ethcommon.Address) *CoreObserver {
	evmSigner, err := MockEVMSigner(
		evmChain,
		connectorAddress,
		erc20CustodyAddress,
	)
	require.NoError(t, err)
	btcSigner, err := MockBTCSigner(btcChain)
	require.NoError(t, err)

	observer := &CoreObserver{
		signerMap: map[int64]interfaces.ChainSigner{
			evmChain.ChainId: evmSigner,
			btcChain.ChainId: btcSigner,
		},
	}
	return observer
}

func MockEVMSigner(
	chain common.Chain,
	connectorAddress ethcommon.Address,
	erc20CustodyAddress ethcommon.Address,
) (*evm.Signer, error) {
	return evm.NewEVMSigner(
		chain,
		stub.EVMRPCEnabled,
		nil,
		config.GetConnectorABI(),
		config.GetERC20CustodyABI(),
		connectorAddress,
		erc20CustodyAddress,
		clientcommon.ClientLogger{},
		&metrics.TelemetryServer{},
	)
}

func MockBTCSigner(chain common.Chain) (*bitcoin.BTCSigner, error) {
	return bitcoin.NewBTCSigner(
		config.BTCConfig{},
		nil,
		clientcommon.ClientLogger{},
		&metrics.TelemetryServer{},
	)
}

func MockCoreContext(evmChain common.Chain, evmChainParams *observertypes.ChainParams) *corecontext.ZetaCoreContext {
	// new config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain: evmChain,
	}
	// new core context
	coreContext := corecontext.NewZetaCoreContext(cfg)
	evmChainParamsMap := make(map[int64]*observertypes.ChainParams)
	evmChainParamsMap[evmChain.ChainId] = evmChainParams

	// feed with chain params
	coreContext.Update(
		&observertypes.Keygen{},
		[]common.Chain{evmChain},
		evmChainParamsMap,
		nil,
		"",
		true,
		zerolog.Logger{},
	)
	return coreContext
}

func Test_GetUpdatedSigner(t *testing.T) {
	// create core observer
	evmChain := common.EthChain()
	btcChain := common.BtcMainnetChain()
	connectorAdderss := testutils.ConnectorAddresses[evmChain.ChainId]
	erc20CustodyAddress := testutils.CustodyAddresses[evmChain.ChainId]

	// update chain params in core context
	chainParams := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConnectorContractAddress:    testutils.OtherAddress1,
		Erc20CustodyContractAddress: testutils.OtherAddress2,
	}

	t.Run("signer not found", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, connectorAdderss, erc20CustodyAddress)
		coreContext := MockCoreContext(evmChain, chainParams)
		_, err := observer.GetUpdatedSigner(coreContext, common.BscMainnetChain().ChainId)
		require.ErrorContains(t, err, "signer not found")
	})
	t.Run("signer is not an EVM signer", func(t *testing.T) {
		// swap evmChain and btcChain to mess up the signer map
		observer := MockCoreObserver(t, btcChain, evmChain, connectorAdderss, erc20CustodyAddress)
		coreContext := MockCoreContext(evmChain, chainParams)
		_, err := observer.GetUpdatedSigner(coreContext, evmChain.ChainId)
		require.ErrorContains(t, err, "signer is not an EVM signer")
	})
	t.Run("should be able to update connector and erc20 custody address", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, connectorAdderss, erc20CustodyAddress)
		coreContext := MockCoreContext(evmChain, chainParams)
		signer, err := observer.GetUpdatedSigner(coreContext, evmChain.ChainId)
		require.NoError(t, err)
		evmSigner, ok := signer.(*evm.Signer)
		require.True(t, ok)
		require.Equal(t, testutils.OtherAddress1, evmSigner.GetZetaConnectorAddress().Hex())
		require.Equal(t, testutils.OtherAddress2, evmSigner.GetERC20CustodyAddress().Hex())
	})
}
