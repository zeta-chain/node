package keeper_test

import (
	"encoding/json"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	zetacommon "github.com/zeta-chain/zetacore/common"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestTransactionReceipt(t *testing.T) {
	t.Run("should not process transaction for legacy address with correct chain params", func(t *testing.T) {
		k, ctx, _, zk := testkeeper.CrosschainKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// load receipt
		receipt := ethtypes.Receipt{}
		err := json.Unmarshal([]byte(GetSampleBlockLegacyToAddress()), &receipt)
		require.NoError(t, err)

		// set foreign coins for the contract
		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: receipt.Logs[0].Address.String(),
			ForeignChainId:       zetacommon.BtcMainnetChain().ChainId,
		})
		// logs 1, 2, 3 are not valid
		// log 4 is valid , but the to address is not valid
		for i, log := range receipt.Logs {
			eventWithdrawal, err := k.ParseZRC20WithdrawalEvent(ctx, *log)
			if i < 3 {
				require.Error(t, err)
				require.Nil(t, eventWithdrawal)
				continue
			}
			require.ErrorContains(t, err, "not P2WPKH address")
			require.NotNil(t, eventWithdrawal)
			require.Equal(t, "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", string(eventWithdrawal.To))
		}
	})
	t.Run("should not process transaction for legacy mainnet address with Testnet params ,", func(t *testing.T) {
		k, ctx, _, zk := testkeeper.CrosschainKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// load receipt
		receipt := ethtypes.Receipt{}
		err := json.Unmarshal([]byte(GetSampleBlockLegacyToAddress()), &receipt)
		require.NoError(t, err)

		// set foreign coins for the contract
		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: receipt.Logs[0].Address.String(),
			ForeignChainId:       zetacommon.BtcTestNetChain().ChainId,
		})
		// logs 1, 2, 3 are not valid
		// log 4 is valid , but the `to` address is not valid
		for i, log := range receipt.Logs {
			eventWithdrawal, err := k.ParseZRC20WithdrawalEvent(ctx, *log)
			if i < 3 {
				require.Error(t, err)
				require.Nil(t, eventWithdrawal)
				continue
			}
			require.ErrorContains(t, err, "decode address failed")
			require.NotNil(t, eventWithdrawal)
			require.Equal(t, "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", string(eventWithdrawal.To))
		}
	})

	t.Run("should not process transaction for legacy mainnet address with RegTest params ,", func(t *testing.T) {
		k, ctx, _, zk := testkeeper.CrosschainKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// load receipt
		receipt := ethtypes.Receipt{}
		err := json.Unmarshal([]byte(GetSampleBlockLegacyToAddress()), &receipt)
		require.NoError(t, err)

		// set foreign coins for the contract
		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: receipt.Logs[0].Address.String(),
			ForeignChainId:       zetacommon.BtcRegtestChain().ChainId,
		})
		// logs 1, 2, 3 are not valid
		// log 4 is valid , but the `to` address is not valid
		for i, log := range receipt.Logs {
			eventWithdrawal, err := k.ParseZRC20WithdrawalEvent(ctx, *log)
			if i < 3 {
				require.Error(t, err)
				require.Nil(t, eventWithdrawal)
				continue
			}
			require.ErrorContains(t, err, "decode address failed")
			require.NotNil(t, eventWithdrawal)
			require.Equal(t, "1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3", string(eventWithdrawal.To))
		}
	})

	t.Run("should process valid tx", func(t *testing.T) {
		k, ctx, _, zk := testkeeper.CrosschainKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		// load receipt
		receipt := ethtypes.Receipt{}
		err := json.Unmarshal([]byte(GetValidSampleBlock()), &receipt)
		require.NoError(t, err)

		// set foreign coins for the contract
		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: receipt.Logs[0].Address.String(),
			ForeignChainId:       zetacommon.BtcMainnetChain().ChainId,
		})
		// logs 1, 2, 3 are not valid
		// log 4 is valid
		for i, log := range receipt.Logs {
			eventWithdrawal, err := k.ParseZRC20WithdrawalEvent(ctx, *log)
			if i < 3 {
				require.Error(t, err)
				require.Nil(t, eventWithdrawal)
				continue
			}

			require.NoError(t, err)
			require.NotNil(t, eventWithdrawal)
			require.Equal(t, "bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu", string(eventWithdrawal.To))
		}
	})
}

// receiver is 1EYVvXLusCxtVuEwoYvWRyN5EZTXwPVvo3
func GetSampleBlockLegacyToAddress() string {
	return "{\n  \"type\": \"0x2\",\n  \"root\": \"0x\",\n  \"status\": \"0x1\",\n  \"cumulativeGasUsed\": \"0x4e7a38\",\n  \"logsBloom\": \"0x00000000000000000000010000020000000000000000000000000000000000020000000100000000000000000000000080000000000000000000000400200000200000000002000000000008000000000000000000000000000000000000000000000000020000000000000000800800000040000000000000000010000000000000000000000000000000000000000000000000000004000000000000000000020000000000000000000000000000000000000000000000000000000000010000000002000000000000000000000000000000000000000000000000000020000010000000000000000001000000000000000000040200000000000000000000\",\n  \"logs\": [\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x000000000000000000000000735b14bb79463307aacbed86daf3322b1e6226ab\"\n      ],\n      \"data\": \"0x0000000000000000000000000000000000000000000000000000000000003790\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x46\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x00000000000000000000000013a0c5930c028511dc02665e7285134b6d11a5f4\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000000000006a1217\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x47\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\",\n        \"0x0000000000000000000000000000000000000000000000000000000000000000\"\n      ],\n      \"data\": \"0x00000000000000000000000000000000000000000000000000000000006a0c70\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x48\",\n      \"removed\": false\n    },\n    {\n      \"address\": \"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\n      \"topics\": [\n        \"0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955\",\n        \"0x000000000000000000000000313e74f7755afbae4f90e02ca49f8f09ff934a37\"\n      ],\n      \"data\": \"0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000006a0c700000000000000000000000000000000000000000000000000000000000003790000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000223145595676584c7573437874567545776f59765752794e35455a5458775056766f33000000000000000000000000000000000000000000000000000000000000\",\n      \"blockNumber\": \"0x1a2ad3\",\n      \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n      \"transactionIndex\": \"0x10\",\n      \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n      \"logIndex\": \"0x49\",\n      \"removed\": false\n    }\n  ],\n  \"transactionHash\": \"0x81126c18c7ca7d1fb7ded6644a87802e91bf52154ee4af7a5b379354e24fb6e0\",\n  \"contractAddress\": \"0x0000000000000000000000000000000000000000\",\n  \"gasUsed\": \"0x12521\",\n  \"blockHash\": \"0x5cb338544f64a226f4bfccb7a8d977f861c13ad73f7dd4317b66b00dd95de51c\",\n  \"blockNumber\": \"0x1a2ad3\",\n  \"transactionIndex\": \"0x10\"\n}\n"
}

// receiver is bc1qysd4sp9q8my59ul9wsf5rvs9p387hf8vfwatzu
func GetValidSampleBlock() string {
	return "{\"type\":\"0x2\",\"root\":\"0x\",\"status\":\"0x1\",\"cumulativeGasUsed\":\"0x1f25ed\",\"logsBloom\":\"0x00000000000000000000000000020000000000000000000000000000000000020000000100000000000000000040000080000000000000000000000400200000200000000002000000000008000000000000000000000000000000000000000000000000020000000000000000800800000000000000000000000010000000000000000000000000000000000000000000000000000004000000000000000000020000000001000000000000000000000000000000000000000000000000010000000002000000000000000010000000000000000000000000000000000020000010000000000000000000000000000000000000040200000000000000000000\",\"logs\":[{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x000000000000000000000000735b14bb79463307aacbed86daf3322b1e6226ab\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000003d84\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x28\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x00000000000000000000000013a0c5930c028511dc02665e7285134b6d11a5f4\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000978c98\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x29\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\",\"0x0000000000000000000000000000000000000000000000000000000000000000\"],\"data\":\"0x0000000000000000000000000000000000000000000000000000000000003039\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x2a\",\"removed\":false},{\"address\":\"0x13a0c5930c028511dc02665e7285134b6d11a5f4\",\"topics\":[\"0x9ffbffc04a397460ee1dbe8c9503e098090567d6b7f4b3c02a8617d800b6d955\",\"0x00000000000000000000000033ead83db0d0c682b05ead61e8d8f481bb1b4933\"],\"data\":\"0x000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000030390000000000000000000000000000000000000000000000000000000000003d840000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a626331717973643473703971386d793539756c3977736635727673397033383768663876667761747a7500000000000000000000000000000000000000000000\",\"blockNumber\":\"0x1a00f3\",\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"transactionIndex\":\"0x8\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"logIndex\":\"0x2b\",\"removed\":false}],\"transactionHash\":\"0x9aaefece38fd2bd87077038a63fffb7c84cc8dd1ed01de134a8504a1f9a410c3\",\"contractAddress\":\"0x0000000000000000000000000000000000000000\",\"gasUsed\":\"0x12575\",\"blockHash\":\"0x9517356f0b3877990590421266f02a4ff349b7476010ee34dd5f0dfc85c2684f\",\"blockNumber\":\"0x1a00f3\",\"transactionIndex\":\"0x8\"}\n"
}
