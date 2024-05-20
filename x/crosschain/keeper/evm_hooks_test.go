package keeper_test

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// SetupStateForProcessLogsZetaSent sets up additional state required for processing logs for ZetaSent events
// This sets up the gas coin, zrc20 contract, gas price, zrc20 pool.
// This should be used in conjunction with SetupStateForProcessLogs for processing ZetaSent events
func SetupStateForProcessLogsZetaSent(t *testing.T, ctx sdk.Context, k *crosschainkeeper.Keeper, zk keepertest.ZetaKeepers, sdkk keepertest.SDKKeepers, chain chains.Chain, admin string) {

	assetAddress := sample.EthAddress().String()
	gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chain.ChainId, "ethereum", "ETH")
	zrc20Addr := deployZRC20(
		t,
		ctx,
		zk.FungibleKeeper,
		sdkk.EvmKeeper,
		chain.ChainId,
		"ethereum",
		assetAddress,
		"ETH",
	)

	_, err := zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, gasZRC20, big.NewInt(withdrawFee))
	require.NoError(t, err)
	_, err = zk.FungibleKeeper.UpdateZRC20ProtocolFlatFee(ctx, zrc20Addr, big.NewInt(withdrawFee))
	require.NoError(t, err)

	k.SetGasPrice(ctx, crosschaintypes.GasPrice{
		ChainId:     chain.ChainId,
		MedianIndex: 0,
		Prices:      []uint64{gasPrice},
	})
	setupZRC20Pool(
		t,
		ctx,
		zk.FungibleKeeper,
		sdkk.BankKeeper,
		zrc20Addr,
	)
}

// SetupStateForProcessLogs sets up observer state for required for processing logs
// It deploys system contracts, sets up TSS, gas price, chain nonce's, pending nonce's.These are all required to create a cctx from a log
func SetupStateForProcessLogs(t *testing.T, ctx sdk.Context, k *crosschainkeeper.Keeper, zk keepertest.ZetaKeepers, sdkk keepertest.SDKKeepers, chain chains.Chain) {

	deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	k.SetGasPrice(ctx, crosschaintypes.GasPrice{
		ChainId: chain.ChainId,
		Prices:  []uint64{100},
	})

	zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
		Index:   chain.ChainName.String(),
		ChainId: chain.ChainId,
		Nonce:   0,
	})
	zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
		NonceLow:  0,
		NonceHigh: 0,
		ChainId:   chain.ChainId,
		Tss:       tss.TssPubkey,
	})
}

func TestParseZRC20WithdrawalEvent(t *testing.T) {
	t.Run("unable to parse an event with an invalid address in event log", func(t *testing.T) {
		for i, log := range sample.GetInvalidZRC20WithdrawToExternal(t).Logs {
			event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*log)
			if i < 3 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", string(event.To))
		}
	})

	t.Run("successfully parse event for a valid BTC withdrawal", func(t *testing.T) {
		for i, log := range sample.GetValidZRC20WithdrawToBTC(t).Logs {
			event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*log)
			if i < 3 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, "0x33EaD83db0D0c682B05ead61E8d8f481Bb1B4933", event.From.Hex())
			require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", string(event.To))
		}
	})

	t.Run("successfully parse valid event for ETH withdrawal", func(t *testing.T) {
		for i, log := range sample.GetValidZrc20WithdrawToETH(t).Logs {
			event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*log)
			if i != 11 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, "0x5daBFdd153Aaab4a970fD953DcFEEE8BF6Bb946E", ethcommon.BytesToAddress(event.To).String())
			require.Equal(t, "0x8E0f8E7E9E121403e72151d00F4937eACB2D9Ef3", event.From.Hex())
		}
	})

	t.Run("failed to parse event with a valid address but no topic present", func(t *testing.T) {
		for _, log := range sample.GetValidZRC20WithdrawToBTC(t).Logs {
			log.Topics = nil
			event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*log)
			require.ErrorContains(t, err, "invalid log - no topics")
			require.Nil(t, event)
		}
	})
}
func TestValidateZrc20WithdrawEvent(t *testing.T) {
	t.Run("successfully validate a valid event", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZRC20WithdrawToBTC(t).Logs[3])
		require.NoError(t, err)
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, chains.BtcMainnetChain.ChainId)
		require.NoError(t, err)
	})

	t.Run("unable to validate a event with an invalid amount", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZRC20WithdrawToBTC(t).Logs[3])
		require.NoError(t, err)
		btcMainNetWithdrawalEvent.Value = big.NewInt(0)
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, chains.BtcMainnetChain.ChainId)
		require.ErrorContains(t, err, "ParseZRC20WithdrawalEvent: invalid amount")
	})

	t.Run("unable to validate a event with an invalid chain ID", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZRC20WithdrawToBTC(t).Logs[3])
		require.NoError(t, err)
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, chains.BtcTestNetChain.ChainId)
		require.ErrorContains(t, err, "invalid address")
	})

	t.Run("unable to validate an unsupported address type", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZRC20WithdrawToBTC(t).Logs[3])
		require.NoError(t, err)
		btcMainNetWithdrawalEvent.To = []byte("04b2891ba8cb491828db3ebc8a780d43b169e7b3974114e6e50f9bab6ec" +
			"63c2f20f6d31b2025377d05c2a704d3bd799d0d56f3a8543d79a01ab6084a1cb204f260")
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, chains.BtcMainnetChain.ChainId)
		require.ErrorContains(t, err, "unsupported address")
	})
}

func TestKeeper_ProcessZRC20WithdrawalEvent(t *testing.T) {
	t.Run("successfully process ZRC20Withdrawal to BTC chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BtcMainnetChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZRC20WithdrawToBTC(t).Logs[3])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		event.Raw.Address = zrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("successfully process ZRC20Withdrawal to ETH chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = zrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "0x5daBFdd153Aaab4a970fD953DcFEEE8BF6Bb946E", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("unable to process ZRC20Withdrawal if foreign coin is not found", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "cannot find foreign coin with emittingContract address")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process ZRC20Withdrawal if receiver chain is not supported", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = zrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "chain not supported")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process ZRC20Withdrawal if zeta chainID is not correctly set", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = zrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()
		ctx = ctx.WithChainID("test_21-1")

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "failed to convert chainID: chain 21 not found")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process ZRC20Withdrawal if to address is not in correct format", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = zrc20
		event.To = ethcommon.Address{}.Bytes()
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "cannot encode address")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process ZRC20Withdrawal if gaslimit not set on zrc20 contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = zrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()
		fc, _ := zk.FungibleKeeper.GetForeignCoins(ctx, zrc20.Hex())

		fungibleMock.On("GetForeignCoins", mock.Anything, mock.Anything).Return(fc, true)
		fungibleMock.On("QueryGasLimit", mock.Anything, mock.Anything).Return(big.NewInt(0), fmt.Errorf("error querying gas limit"))
		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "error querying gas limit")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process ZRC20Withdrawal if gasprice is not set in crosschain keeper", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		k.RemoveGasPrice(ctx, strconv.FormatInt(chainID, 10))

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = zrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "gasprice not found")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})

	t.Run("unable to process ZRC20Withdrawal if process cctx fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
			Index:   chain.ChainName.String(),
			ChainId: chain.ChainId,
			Nonce:   1,
		})

		event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*sample.GetValidZrc20WithdrawToETH(t).Logs[11])
		require.NoError(t, err)
		zrc20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "ethereum", "ETH")
		event.Raw.Address = zrc20
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZRC20WithdrawalEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "ProcessWithdrawalEvent: update nonce failed")
		require.Empty(t, k.GetAllCrossChainTx(ctx))
	})
}

func TestKeeper_ParseZetaSentEvent(t *testing.T) {
	t.Run("successfully parse a valid event", func(t *testing.T) {
		logs := sample.GetValidZetaSentDestinationExternal(t).Logs
		for i, log := range logs {
			connector := log.Address
			event, err := crosschainkeeper.ParseZetaSentEvent(*log, connector)
			if i < 4 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.Equal(t, chains.EthChain.ChainId, event.DestinationChainId.Int64())
			require.Equal(t, "70000000000000000000", event.ZetaValueAndGas.String())
			require.Equal(t, "0x60983881bdf302dcfa96603A58274D15D5966209", event.SourceTxOriginAddress.String())
			require.Equal(t, "0xF0a3F93Ed1B126142E61423F9546bf1323Ff82DF", event.ZetaTxSenderAddress.String())
		}
	})

	t.Run("unable to parse if topics field is empty", func(t *testing.T) {
		logs := sample.GetValidZetaSentDestinationExternal(t).Logs
		for _, log := range logs {
			connector := log.Address
			log.Topics = nil
			event, err := crosschainkeeper.ParseZetaSentEvent(*log, connector)
			require.ErrorContains(t, err, "ParseZetaSentEvent: invalid log - no topics")
			require.Nil(t, event)
		}
	})

	t.Run("unable to parse if connector address does not match", func(t *testing.T) {
		logs := sample.GetValidZetaSentDestinationExternal(t).Logs
		for i, log := range logs {
			event, err := crosschainkeeper.ParseZetaSentEvent(*log, sample.EthAddress())
			if i < 4 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.ErrorContains(t, err, "does not match connectorZEVM")
			require.Nil(t, event)
		}
	})
}
func TestKeeper_ProcessZetaSentEvent(t *testing.T) {
	t.Run("successfully process ZetaSentEvent", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)

		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolices(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsZetaSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseZetaSentEvent(*sample.GetValidZetaSentDestinationExternal(t).Logs[4], sample.GetValidZetaSentDestinationExternal(t).Logs[4].Address)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZetaSentEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, strings.Compare("0x60983881bdf302dcfa96603a58274d15d5966209", cctxList[0].GetCurrentOutboundParam().Receiver), 0)
		require.Equal(t, chains.EthChain.ChainId, cctxList[0].GetCurrentOutboundParam().ReceiverChainId)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("unable to process ZetaSentEvent if fungible module does not have enough balance", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolices(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsZetaSent(t, ctx, k, zk, sdkk, chain, admin)

		event, err := crosschainkeeper.ParseZetaSentEvent(*sample.GetValidZetaSentDestinationExternal(t).Logs[4], sample.GetValidZetaSentDestinationExternal(t).Logs[4].Address)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZetaSentEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "ProcessZetaSentEvent: failed to burn coins from fungible")
	})

	t.Run("unable to process ZetaSentEvent if receiver chain is not supported", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolices(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsZetaSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseZetaSentEvent(*sample.GetValidZetaSentDestinationExternal(t).Logs[4], sample.GetValidZetaSentDestinationExternal(t).Logs[4].Address)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()
		err = k.ProcessZetaSentEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "chain not supported")
	})

	t.Run("unable to process ZetaSentEvent if zetachain chain id not correctly set in context", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolices(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsZetaSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseZetaSentEvent(*sample.GetValidZetaSentDestinationExternal(t).Logs[4], sample.GetValidZetaSentDestinationExternal(t).Logs[4].Address)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()
		ctx = ctx.WithChainID("test-21-1")
		err = k.ProcessZetaSentEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "ProcessZetaSentEvent: failed to convert chainID")
	})

	t.Run("unable to process ZetaSentEvent if gas pay fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)
		event, err := crosschainkeeper.ParseZetaSentEvent(*sample.GetValidZetaSentDestinationExternal(t).Logs[4], sample.GetValidZetaSentDestinationExternal(t).Logs[4].Address)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()

		err = k.ProcessZetaSentEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "ProcessWithdrawalEvent: pay gas failed")
	})

	t.Run("unable to process ZetaSentEvent if process cctx fails", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolices(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsZetaSent(t, ctx, k, zk, sdkk, chain, admin)

		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
			Index:   chain.ChainName.String(),
			ChainId: chain.ChainId,
			Nonce:   1,
		})
		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)

		event, err := crosschainkeeper.ParseZetaSentEvent(*sample.GetValidZetaSentDestinationExternal(t).Logs[4], sample.GetValidZetaSentDestinationExternal(t).Logs[4].Address)
		require.NoError(t, err)
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()
		tss := sample.Tss()
		err = k.ProcessZetaSentEvent(ctx, event, emittingContract, txOrigin.Hex(), tss)
		require.ErrorContains(t, err, "ProcessWithdrawalEvent: update nonce failed")
	})
}

func TestKeeper_ProcessLogs(t *testing.T) {
	t.Run("successfully parse and process ZRC20Withdrawal to BTC chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BtcMainnetChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.GetValidZRC20WithdrawToBTC(t)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasZRC20
		}
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err := k.ProcessLogs(ctx, block.Logs, emittingContract, txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", cctxList[0].GetCurrentOutboundParam().Receiver)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("successfully parse and process ZetaSentEvent", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.EthChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)
		admin := keepertest.SetAdminPolices(ctx, zk.AuthorityKeeper)
		SetupStateForProcessLogsZetaSent(t, ctx, k, zk, sdkk, chain, admin)

		amount, ok := sdkmath.NewIntFromString("20000000000000000000000")
		require.True(t, ok)
		err := sdkk.BankKeeper.MintCoins(ctx, fungibletypes.ModuleName, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, amount)))
		require.NoError(t, err)
		block := sample.GetValidZetaSentDestinationExternal(t)
		system, found := zk.FungibleKeeper.GetSystemContract(ctx)
		require.True(t, found)
		for _, log := range block.Logs {
			log.Address = ethcommon.HexToAddress(system.ConnectorZevm)
		}
		emittingContract := sample.EthAddress()
		txOrigin := sample.EthAddress()

		err = k.ProcessLogs(ctx, block.Logs, emittingContract, txOrigin.Hex())
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 1)
		require.Equal(t, strings.Compare("0x60983881bdf302dcfa96603a58274d15d5966209", cctxList[0].GetCurrentOutboundParam().Receiver), 0)
		require.Equal(t, chains.EthChain.ChainId, cctxList[0].GetCurrentOutboundParam().ReceiverChainId)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundParams.TxOrigin)
	})

	t.Run("unable to process logs if system contract not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		err := k.ProcessLogs(ctx, sample.GetValidZRC20WithdrawToBTC(t).Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, "cannot find system contract")
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("no cctx created for logs containing no events", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BtcMainnetChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.GetValidZRC20WithdrawToBTC(t)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasZRC20
		}
		block.Logs = block.Logs[:3]

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("no cctx created  for logs containing proper event but not emitted from a known ZRC20 contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		chain := chains.BtcMainnetChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.GetValidZRC20WithdrawToBTC(t)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = sample.EthAddress()
		}

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("no cctx created  for valid logs if Inbound is disabled", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BtcMainnetChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.GetValidZRC20WithdrawToBTC(t)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasZRC20
		}
		zk.ObserverKeeper.SetCrosschainFlags(ctx, observertypes.CrosschainFlags{
			IsInboundEnabled: false,
		})

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, observertypes.ErrInboundDisabled.Error())
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("error returned for invalid event data", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// use the wrong (testnet) chain ID to make the btc address parsing fail
		chain := chains.BtcTestNetChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.GetInvalidZRC20WithdrawToExternal(t)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasZRC20
		}

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, "ParseZRC20WithdrawalEvent: invalid address")
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})

	t.Run("error returned if unable to process an event", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chain := chains.BtcMainnetChain
		chainID := chain.ChainId
		setSupportedChain(ctx, zk, chainID)
		SetupStateForProcessLogs(t, ctx, k, zk, sdkk, chain)

		block := sample.GetValidZRC20WithdrawToBTC(t)
		gasZRC20 := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = gasZRC20
		}
		ctx = ctx.WithChainID("test-21-1")

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, "ProcessZRC20WithdrawalEvent: failed to convert chainID")
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})
}
