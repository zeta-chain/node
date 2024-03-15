package zetaclient

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
)

func MockCoreObserver(t *testing.T, evmChain, btcChain common.Chain, evmChainParams, btcChainParams *observertypes.ChainParams) *CoreObserver {
	// create mock signers and clients
	evmSigner := stub.NewEVMSigner(
		evmChain,
		ethcommon.HexToAddress(evmChainParams.ConnectorContractAddress),
		ethcommon.HexToAddress(evmChainParams.Erc20CustodyContractAddress),
	)
	btcSigner := stub.NewBTCSigner()
	evmClient := stub.NewEVMClient(evmChainParams)
	btcClient := stub.NewBTCClient(btcChainParams)

	// create core observer
	observer := &CoreObserver{
		signerMap: map[int64]interfaces.ChainSigner{
			evmChain.ChainId: evmSigner,
			btcChain.ChainId: btcSigner,
		},
		clientMap: map[int64]interfaces.ChainClient{
			evmChain.ChainId: evmClient,
			btcChain.ChainId: btcClient,
		},
	}
	return observer
}

func CreateCoreContext(evmChain, btcChain common.Chain, evmChainParams, btcChainParams *observertypes.ChainParams) *corecontext.ZetaCoreContext {
	// new config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain: evmChain,
	}
	cfg.BitcoinConfig = config.BTCConfig{
		RPCHost: "localhost",
	}
	// new core context
	coreContext := corecontext.NewZetaCoreContext(cfg)
	evmChainParamsMap := make(map[int64]*observertypes.ChainParams)
	evmChainParamsMap[evmChain.ChainId] = evmChainParams

	// feed chain params
	coreContext.Update(
		&observertypes.Keygen{},
		[]common.Chain{evmChain, btcChain},
		evmChainParamsMap,
		btcChainParams,
		"",
		true,
		zerolog.Logger{},
	)
	return coreContext
}

func Test_GetUpdatedSigner(t *testing.T) {
	// initial parameters for core observer creation
	evmChain := common.EthChain()
	btcChain := common.BtcMainnetChain()
	evmChainParams := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConnectorContractAddress:    testutils.ConnectorAddresses[evmChain.ChainId].Hex(),
		Erc20CustodyContractAddress: testutils.CustodyAddresses[evmChain.ChainId].Hex(),
	}
	btcChainParams := &observertypes.ChainParams{}

	// new chain params in core context
	evmChainParamsNew := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConnectorContractAddress:    testutils.OtherAddress1,
		Erc20CustodyContractAddress: testutils.OtherAddress2,
	}

	t.Run("signer should not be found", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, evmChainParams, btcChainParams)
		coreContext := CreateCoreContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		// BSC signer should not be found
		_, err := observer.GetUpdatedSigner(coreContext, common.BscMainnetChain().ChainId)
		require.ErrorContains(t, err, "signer not found")
	})
	t.Run("should be able to update connector and erc20 custody address", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, evmChainParams, btcChainParams)
		coreContext := CreateCoreContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		// update signer with new connector and erc20 custody address
		signer, err := observer.GetUpdatedSigner(coreContext, evmChain.ChainId)
		require.NoError(t, err)
		require.Equal(t, testutils.OtherAddress1, signer.GetZetaConnectorAddress().Hex())
		require.Equal(t, testutils.OtherAddress2, signer.GetERC20CustodyAddress().Hex())
	})
}

func Test_GetUpdatedChainClient(t *testing.T) {
	// initial parameters for core observer creation
	evmChain := common.EthChain()
	btcChain := common.BtcMainnetChain()
	evmChainParams := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConnectorContractAddress:    testutils.ConnectorAddresses[evmChain.ChainId].Hex(),
		Erc20CustodyContractAddress: testutils.CustodyAddresses[evmChain.ChainId].Hex(),
	}
	btcChainParams := &observertypes.ChainParams{
		ChainId: btcChain.ChainId,
	}

	// new chain params in core context
	evmChainParamsNew := &observertypes.ChainParams{
		ChainId:                     evmChain.ChainId,
		ConfirmationCount:           10,
		GasPriceTicker:              11,
		InTxTicker:                  12,
		OutTxTicker:                 13,
		WatchUtxoTicker:             14,
		ZetaTokenContractAddress:    testutils.OtherAddress1,
		ConnectorContractAddress:    testutils.OtherAddress2,
		Erc20CustodyContractAddress: testutils.OtherAddress3,
		OutboundTxScheduleInterval:  15,
		OutboundTxScheduleLookahead: 16,
		BallotThreshold:             sdk.OneDec(),
		MinObserverDelegation:       sdk.OneDec(),
		IsSupported:                 true,
	}
	btcChainParamsNew := &observertypes.ChainParams{
		ChainId:                     btcChain.ChainId,
		ConfirmationCount:           3,
		GasPriceTicker:              300,
		InTxTicker:                  60,
		OutTxTicker:                 60,
		WatchUtxoTicker:             30,
		ZetaTokenContractAddress:    testutils.OtherAddress1,
		ConnectorContractAddress:    testutils.OtherAddress2,
		Erc20CustodyContractAddress: testutils.OtherAddress3,
		OutboundTxScheduleInterval:  60,
		OutboundTxScheduleLookahead: 200,
		BallotThreshold:             sdk.OneDec(),
		MinObserverDelegation:       sdk.OneDec(),
		IsSupported:                 true,
	}

	t.Run("evm chain client should not be found", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, evmChainParams, btcChainParams)
		coreContext := CreateCoreContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		// BSC chain client should not be found
		_, err := observer.GetUpdatedChainClient(coreContext, common.BscMainnetChain().ChainId)
		require.ErrorContains(t, err, "chain client not found")
	})
	t.Run("chain params in evm chain client should be updated successfully", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, evmChainParams, btcChainParams)
		coreContext := CreateCoreContext(evmChain, btcChain, evmChainParamsNew, btcChainParams)
		// update evm chain client with new chain params
		chainOb, err := observer.GetUpdatedChainClient(coreContext, evmChain.ChainId)
		require.NoError(t, err)
		require.NotNil(t, chainOb)
		require.True(t, observertypes.ChainParamsEqual(*evmChainParamsNew, chainOb.GetChainParams()))
	})
	t.Run("btc chain client should not be found", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, evmChainParams, btcChainParams)
		coreContext := CreateCoreContext(btcChain, btcChain, evmChainParams, btcChainParamsNew)
		// BTC testnet chain client should not be found
		_, err := observer.GetUpdatedChainClient(coreContext, common.BtcTestNetChain().ChainId)
		require.ErrorContains(t, err, "chain client not found")
	})
	t.Run("chain params in btc chain client should be updated successfully", func(t *testing.T) {
		observer := MockCoreObserver(t, evmChain, btcChain, evmChainParams, btcChainParams)
		coreContext := CreateCoreContext(btcChain, btcChain, evmChainParams, btcChainParamsNew)
		// update btc chain client with new chain params
		chainOb, err := observer.GetUpdatedChainClient(coreContext, btcChain.ChainId)
		require.NoError(t, err)
		require.NotNil(t, chainOb)
		require.True(t, observertypes.ChainParamsEqual(*btcChainParamsNew, chainOb.GetChainParams()))
	})
}
