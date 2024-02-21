package keeper_test

import (
	"encoding/json"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

//

func TestValidateZrc20WithdrawEvent(t *testing.T) {
	t.Run("valid event", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*GetValidSampleBlock(t).Logs[3])
		require.NoError(t, err)
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, common.BtcMainnetChain().ChainId)
		require.NoError(t, err)
	})
	t.Run("invalid value", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*GetValidSampleBlock(t).Logs[3])
		require.NoError(t, err)
		btcMainNetWithdrawalEvent.Value = big.NewInt(0)
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, common.BtcMainnetChain().ChainId)
		require.ErrorContains(t, err, "ParseZRC20WithdrawalEvent: invalid amount")
	})
	t.Run("invalid chain ID", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*GetValidSampleBlock(t).Logs[3])
		require.NoError(t, err)
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, common.BtcTestNetChain().ChainId)
		require.ErrorContains(t, err, "address is not for network testnet3")
	})
	t.Run("invalid address type", func(t *testing.T) {
		btcMainNetWithdrawalEvent, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*GetValidSampleBlock(t).Logs[3])
		require.NoError(t, err)
		btcMainNetWithdrawalEvent.To = []byte("1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3")
		err = crosschainkeeper.ValidateZrc20WithdrawEvent(btcMainNetWithdrawalEvent, common.BtcTestNetChain().ChainId)
		require.ErrorContains(t, err, "decode address failed: unknown address type")
	})
}
func TestParseZRC20WithdrawalEvent(t *testing.T) {
	t.Run("invalid address", func(t *testing.T) {
		for i, log := range GetSampleBlockLegacyToAddress(t).Logs {
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
	t.Run("valid address", func(t *testing.T) {
		for i, log := range GetValidSampleBlock(t).Logs {
			event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*log)
			if i < 3 {
				require.ErrorContains(t, err, "event signature mismatch")
				require.Nil(t, event)
				continue
			}
			require.NoError(t, err)
			require.NotNil(t, event)
			require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", string(event.To))
		}
	})
	t.Run("valid address remove topics", func(t *testing.T) {
		for _, log := range GetValidSampleBlock(t).Logs {
			log.Topics = nil
			event, err := crosschainkeeper.ParseZRC20WithdrawalEvent(*log)
			require.ErrorContains(t, err, "invalid log - no topics")
			require.Nil(t, event)
		}
	})
}

func TestKeeper_ProcessLogs(t *testing.T) {
	t.Run("create ZRC20Withdrawal to BTC chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := common.BtcMainnetChain().ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupProcessLogsTest(t, ctx, k, zk, sdkk, chainID)

		block := GetValidSampleBlock(t)
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
		require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", cctxList[0].GetCurrentOutTxParam().Receiver)
		require.Equal(t, emittingContract.Hex(), cctxList[0].InboundTxParams.Sender)
		require.Equal(t, txOrigin.Hex(), cctxList[0].InboundTxParams.TxOrigin)
	})
	t.Run("unable to process logs if system contract not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		err := k.ProcessLogs(ctx, GetValidSampleBlock(t).Logs, sample.EthAddress(), "")
		require.ErrorContains(t, err, "cannot find system contract")
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})
	t.Run("no cctx created for logs containing no events", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := common.BtcMainnetChain().ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupProcessLogsTest(t, ctx, k, zk, sdkk, chainID)

		block := GetValidSampleBlock(t)
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
	t.Run("no cctx created for logs containing proper event but not emitted from a known ZRC20 contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := common.BtcMainnetChain().ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupProcessLogsTest(t, ctx, k, zk, sdkk, chainID)

		block := GetValidSampleBlock(t)
		setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "bitcoin", "BTC")
		for _, log := range block.Logs {
			log.Address = sample.EthAddress()
		}

		err := k.ProcessLogs(ctx, block.Logs, sample.EthAddress(), "")
		require.NoError(t, err)
		cctxList := k.GetAllCrossChainTx(ctx)
		require.Len(t, cctxList, 0)
	})
	t.Run("no cctx created for for valid logs if Inbound is disabled", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := common.BtcMainnetChain().ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupProcessLogsTest(t, ctx, k, zk, sdkk, chainID)

		block := GetValidSampleBlock(t)
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

		chainID := common.BtcMainnetChain().ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupProcessLogsTest(t, ctx, k, zk, sdkk, chainID)

		block := GetSampleBlockLegacyToAddress(t)
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

		chainID := common.BtcMainnetChain().ChainId
		setSupportedChain(ctx, zk, chainID)

		SetupProcessLogsTest(t, ctx, k, zk, sdkk, chainID)

		block := GetValidSampleBlock(t)
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

func SetupProcessLogsTest(t *testing.T, ctx sdk.Context, k *keeper.Keeper, zk keepertest.ZetaKeepers, sdkk keepertest.SDKKeepers, chainID int64) {

	deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
	tss := sample.Tss()
	zk.ObserverKeeper.SetTSS(ctx, tss)
	k.SetGasPrice(ctx, crosschaintypes.GasPrice{
		ChainId: chainID,
		Prices:  []uint64{100},
	})

	zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{
		Index:   common.BtcMainnetChain().ChainName.String(),
		ChainId: chainID,
		Nonce:   0,
	})
	zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
		NonceLow:  0,
		NonceHigh: 0,
		ChainId:   chainID,
		Tss:       tss.TssPubkey,
	})
}

// receiver is 1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3
func GetSampleBlockLegacyToAddress(t *testing.T) (receipt ethtypes.Receipt) {
	block := "{\n  \"type\": \"0x2\",\n  \"root\": \"0x\",\n  \"status\": \"0x1\",\n  \"cumulativeGasUsed\": \"0x4e7a38\",\n  \"logsBloom\": \"0x00000000000000000000010000020000000000000000000000000000000000020000000100000000000000000000000080000000000000000000000400200000200000000002000000000008000000000000000000000000000000000000000000000000020000000000000000800800000040000000000000000010000000000000000000000000000000000000000000000000000004000000000000000000020000000000000000000000000000000000000000000000000000000000010000000002000000000000000000000000000000000000000000000000000020000010000000000000000001000000000000000000040200000000000000000000\",\n  \"logs\": [\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x000000000000000000000000735b14bb79463307aacbed86daf3322b1e6226ab\"\n      ],\n      \"data\": \"0x0000000000000000000000000000000000000000000000000000000000003790\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x46\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x00000000000000000000000013a0c5930c028511dc02665e7285134b6d11a5f4\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000000000006a1217\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x47\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x0000000000000000000000000000000000000000000000000000000000000000\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000000000006a0c70\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x48\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\"\n      ],\n      \"data\": \"0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000006a0c700000000000000000000000000000000000000000000000000000000000003790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000223145595676584c7573437874567545776f59765752794e35455a5458775056766f33000000000000000000000000000000000000000000000000000000000000\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x49\",\n      \"removed\": false\n    }\n  ],\n  \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n  \"contractAddress\": \"0x0000000000000000000000000000000000000000\",\n  \"gasUsed\": \"0x12521\",\n  \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n  \"blockNumber\": \"0x1a2ad3\",\n  \"transactionIndex\": \"0x10\"\n}\n"
	err := json.Unmarshal([]byte(block), &receipt)
	require.NoError(t, err)
	return
}

// receiver is bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu
func GetValidSampleBlock(t *testing.T) (receipt ethtypes.Receipt) {
	block := "{\"type\":\"0x2\",\"root\":\"0x\",\"status\":\"0x1\",\"cumulativeGasUsed\":\"0x1f25ed\",\"logsBloom\":\"0x00000000000000000000000000020000000000000000000000000000000000020000000100000000000000000040000080000000000000000000000400200000200000000002000000000008000000000000000000000000000000000000000000000000020000000000000000800800000000000000000000000010000000000000000000000000000000000000000000000000000004000000000000000000020000000001000000000000000000000000000000000000000000000000010000000002000000000000000010000000000000000000000000000000000020000010000000000000000000000000000000000000040200000000000000000000\",\"logs\":[{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x000000000000000000000000735b14bb79463307aacbed86daf3322b1e6226ab\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000003d84\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x28\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x00000000000000000000000013a0c5930c028511dc02665e7285134b6d11a5f4\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000978c98\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x29\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x0000000000000000000000000000000000000000000000000000000000000000\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000003039\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x2a\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\"],\"data\":\"0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000030390000000000000000000000000000000000000000000000000000000000003d840000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a626331717973643473703971386d793539756c3977736635727673397033383768663876667761747a7500000000000000000000000000000000000000000000\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x2b\",\"removed\":false}],\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"contractAddress\":\"0x0000000000000000000000000000000000000000\",\"gasUsed\":\"0x12575\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"blockNumber\":\"0x1a00f3\",\"transactionIndex\":\"0x8\"}\n"
	err := json.Unmarshal([]byte(block), &receipt)
	require.NoError(t, err)
	return
}
