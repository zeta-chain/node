package types_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/gas"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMigrateERC20CustodyFundsCmdCCTX(t *testing.T) {
	t.Run("returns a new CCTX for migrating ERC20 custody funds with unique index", func(t *testing.T) {
		// ARRANGE
		creator := sample.AccAddress()
		erc20Address := sample.EthAddress().String()
		custodyContractAddress := sample.EthAddress().String()
		newCustodyContractAddress := sample.EthAddress().String()
		chainID := int64(42)
		amount := sdkmath.NewUint(1000)
		gasPrice := "100000"
		priorityFee := "100000"
		tssPubKey := sample.PubKeyString()
		currentNonce := uint64(1)

		// ACT
		cctx := types.MigrateERC20CustodyFundsCmdCCTX(
			creator,
			erc20Address,
			custodyContractAddress,
			newCustodyContractAddress,
			chainID,
			amount,
			gasPrice,
			priorityFee,
			tssPubKey,
			currentNonce,
		)
		cctxDifferentERC20Address := types.MigrateERC20CustodyFundsCmdCCTX(
			creator,
			sample.EthAddress().String(),
			custodyContractAddress,
			newCustodyContractAddress,
			chainID,
			amount,
			gasPrice,
			priorityFee,
			tssPubKey,
			currentNonce,
		)
		cctxDifferentNonce := types.MigrateERC20CustodyFundsCmdCCTX(
			creator,
			erc20Address,
			custodyContractAddress,
			newCustodyContractAddress,
			chainID,
			amount,
			gasPrice,
			priorityFee,
			tssPubKey,
			currentNonce+1,
		)
		cctxDifferentTSSPubkey := types.MigrateERC20CustodyFundsCmdCCTX(
			creator,
			erc20Address,
			custodyContractAddress,
			newCustodyContractAddress,
			chainID,
			amount,
			gasPrice,
			priorityFee,
			sample.PubKeyString(),
			currentNonce,
		)

		// ASSERT
		require.NotEmpty(t, cctx.Index)
		require.EqualValues(t, creator, cctx.Creator)
		require.EqualValues(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.EqualValues(t, fmt.Sprintf("%s:%s,%s,1000",
			constant.CmdMigrateERC20CustodyFunds,
			newCustodyContractAddress,
			erc20Address,
		), cctx.RelayedMessage)
		require.EqualValues(t, creator, cctx.InboundParams.Sender)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.InboundParams.CoinType)
		require.Len(t, cctx.OutboundParams, 1)
		require.EqualValues(t, custodyContractAddress, cctx.OutboundParams[0].Receiver)
		require.EqualValues(t, chainID, cctx.OutboundParams[0].ReceiverChainId)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.OutboundParams[0].CoinType)
		require.EqualValues(t, sdkmath.NewUint(0), cctx.OutboundParams[0].Amount)
		require.EqualValues(t, 100_000, cctx.OutboundParams[0].CallOptions.GasLimit)
		require.EqualValues(t, gasPrice, cctx.OutboundParams[0].GasPrice)
		require.EqualValues(t, priorityFee, cctx.OutboundParams[0].GasPriorityFee)
		require.EqualValues(t, tssPubKey, cctx.OutboundParams[0].TssPubkey)

		// check erc20, TSS pubkey and nonce produce unique index
		require.NotEqual(t, cctx.Index, cctxDifferentERC20Address.Index)
		require.NotEqual(t, cctx.Index, cctxDifferentNonce.Index)
		require.NotEqual(t, cctx.Index, cctxDifferentTSSPubkey.Index)
	})
}

func TestGetERC20CustodyMigrationCCTXIndexString(t *testing.T) {
	t.Run("returns the unique index string for the CCTX for migrating ERC20 custody funds", func(t *testing.T) {
		// ARRANGE
		tssPubKey := sample.PubKeyString()
		nonce := uint64(1)
		chainID := int64(42)
		erc20Address := sample.EthAddress().String()

		// ACT
		index := types.GetERC20CustodyMigrationCCTXIndexString(
			tssPubKey,
			nonce,
			chainID,
			erc20Address,
		)
		indexDifferentTSSPubkey := types.GetERC20CustodyMigrationCCTXIndexString(
			sample.PubKeyString(),
			nonce,
			chainID,
			erc20Address,
		)
		indexDifferentNonce := types.GetERC20CustodyMigrationCCTXIndexString(
			tssPubKey,
			nonce+1,
			chainID,
			erc20Address,
		)
		indexDifferentERC20Address := types.GetERC20CustodyMigrationCCTXIndexString(
			tssPubKey,
			nonce,
			chainID,
			sample.EthAddress().String(),
		)
		indexDifferentChainID := types.GetERC20CustodyMigrationCCTXIndexString(
			tssPubKey,
			nonce,
			chainID+1,
			erc20Address,
		)

		// ASSERT
		require.NotEmpty(t, index)
		require.NotEqual(t, index, indexDifferentTSSPubkey)
		require.NotEqual(t, index, indexDifferentNonce)
		require.NotEqual(t, index, indexDifferentERC20Address)
		require.NotEqual(t, index, indexDifferentChainID)
	})
}

func TestUpdateERC20CustodyPauseStatusCmdCCTX(t *testing.T) {
	t.Run("returns a new CCTX to pause ERC20Custody", func(t *testing.T) {
		// ARRANGE
		creator := sample.AccAddress()
		custodyContractAddress := sample.EthAddress().String()
		chainID := int64(42)
		gasPrice := "100000"
		priorityFee := "100000"
		tssPubKey := sample.PubKeyString()
		currentNonce := uint64(1)

		// ACT
		cctx := types.UpdateERC20CustodyPauseStatusCmdCCTX(
			creator,
			custodyContractAddress,
			chainID,
			true,
			gasPrice,
			priorityFee,
			tssPubKey,
			currentNonce,
		)
		cctxDifferentNonce := types.UpdateERC20CustodyPauseStatusCmdCCTX(
			creator,
			custodyContractAddress,
			chainID,
			true,
			gasPrice,
			priorityFee,
			tssPubKey,
			currentNonce+1,
		)
		cctxDifferentTSSPubkey := types.UpdateERC20CustodyPauseStatusCmdCCTX(
			creator,
			custodyContractAddress,
			chainID,
			true,
			gasPrice,
			priorityFee,
			sample.PubKeyString(),
			currentNonce,
		)
		cctxDifferentChainID := types.UpdateERC20CustodyPauseStatusCmdCCTX(
			creator,
			custodyContractAddress,
			chainID+1,
			true,
			gasPrice,
			priorityFee,
			tssPubKey,
			currentNonce,
		)

		// ASSERT
		require.NotEmpty(t, cctx.Index)
		require.EqualValues(t, creator, cctx.Creator)
		require.EqualValues(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.EqualValues(t, fmt.Sprintf("%s:%s",
			constant.CmdUpdateERC20CustodyPauseStatus,
			constant.OptionPause,
		), cctx.RelayedMessage)
		require.EqualValues(t, creator, cctx.InboundParams.Sender)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.InboundParams.CoinType)
		require.Len(t, cctx.OutboundParams, 1)
		require.EqualValues(t, custodyContractAddress, cctx.OutboundParams[0].Receiver)
		require.EqualValues(t, chainID, cctx.OutboundParams[0].ReceiverChainId)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.OutboundParams[0].CoinType)
		require.EqualValues(t, sdkmath.NewUint(0), cctx.OutboundParams[0].Amount)
		require.EqualValues(t, 100_000, cctx.OutboundParams[0].CallOptions.GasLimit)
		require.EqualValues(t, gasPrice, cctx.OutboundParams[0].GasPrice)
		require.EqualValues(t, priorityFee, cctx.OutboundParams[0].GasPriorityFee)
		require.EqualValues(t, tssPubKey, cctx.OutboundParams[0].TssPubkey)

		// check erc20, TSS pubkey and nonce produce unique index
		require.NotEqual(t, cctx.Index, cctxDifferentNonce.Index)
		require.NotEqual(t, cctx.Index, cctxDifferentTSSPubkey.Index)
		require.NotEqual(t, cctx.Index, cctxDifferentChainID.Index)
	})
}

func TestGetERC20CustodyPausingCmdCCTXIndecString(t *testing.T) {
	t.Run("returns the unique index string for the CCTX for updating ERC20 custody pause status", func(t *testing.T) {
		// ARRANGE
		tssPubKey := sample.PubKeyString()
		nonce := uint64(1)
		chainID := int64(42)

		// ACT
		index := types.GetERC20CustodyPausingCmdCCTXIndexString(
			tssPubKey,
			nonce,
			chainID,
		)
		indexDifferentTSSPubkey := types.GetERC20CustodyPausingCmdCCTXIndexString(
			sample.PubKeyString(),
			nonce,
			chainID,
		)
		indexDifferentNonce := types.GetERC20CustodyPausingCmdCCTXIndexString(
			tssPubKey,
			nonce+1,
			chainID,
		)
		indexDifferentChainID := types.GetERC20CustodyPausingCmdCCTXIndexString(
			tssPubKey,
			nonce,
			chainID+1,
		)

		// ASSERT
		require.NotEmpty(t, index)
		require.NotEqual(t, index, indexDifferentTSSPubkey)
		require.NotEqual(t, index, indexDifferentNonce)
		require.NotEqual(t, index, indexDifferentChainID)
	})
}

func TestWhitelistERC20CmdCCTX(t *testing.T) {
	t.Run("returns a new CCTX for whitelisting ERC20 tokens", func(t *testing.T) {
		// ARRANGE
		creator := sample.AccAddress()
		zrc20Address := sample.EthAddress()
		erc20Address := sample.EthAddress().Hex()
		custodyAddress := sample.EthAddress().Hex()
		chainID := int64(42)
		gasPrice := "100000"
		priorityFee := "100000"
		tssPubKey := sample.PubKeyString()

		// ACT
		cctx := types.WhitelistAssetCmdCCTX(
			creator,
			zrc20Address,
			erc20Address,
			custodyAddress,
			chainID,
			gasPrice,
			priorityFee,
			tssPubKey,
		)
		cctxDifferentZRC20Address := types.WhitelistAssetCmdCCTX(
			creator,
			sample.EthAddress(),
			erc20Address,
			custodyAddress,
			chainID,
			gasPrice,
			priorityFee,
			tssPubKey,
		)

		// ASSERT
		require.NotEmpty(t, cctx.Index)
		require.EqualValues(t, creator, cctx.Creator)
		require.EqualValues(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.EqualValues(t, fmt.Sprintf("%s:%s", constant.CmdWhitelistAsset, erc20Address), cctx.RelayedMessage)
		require.EqualValues(t, creator, cctx.InboundParams.Sender)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.InboundParams.CoinType)
		require.Len(t, cctx.OutboundParams, 1)
		require.EqualValues(t, custodyAddress, cctx.OutboundParams[0].Receiver)
		require.EqualValues(t, chainID, cctx.OutboundParams[0].ReceiverChainId)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.OutboundParams[0].CoinType)
		require.EqualValues(t, sdkmath.NewUint(0), cctx.OutboundParams[0].Amount)
		require.EqualValues(t, 100_000, cctx.OutboundParams[0].CallOptions.GasLimit)
		require.EqualValues(t, gasPrice, cctx.OutboundParams[0].GasPrice)
		require.EqualValues(t, priorityFee, cctx.OutboundParams[0].GasPriorityFee)
		require.EqualValues(t, tssPubKey, cctx.OutboundParams[0].TssPubkey)

		// check zrc20 address produces unique index
		require.NotEqual(t, cctx.Index, cctxDifferentZRC20Address.Index)
	})
}

func TestMigrateFundCmdCCTX(t *testing.T) {
	t.Run("returns a new CCTX for migrating funds on EVM", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := chains.Ethereum.ChainId
		amount := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := sample.Tss()
		newTSSPubkey := sample.Tss()

		// ACT
		cctx, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey.TssPubkey,
			newTSSPubkey.TssPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.NoError(t, err)
		require.NotEmpty(t, cctx.Index)
		require.EqualValues(t, creator, cctx.Creator)
		require.EqualValues(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.EqualValues(
			t,
			fmt.Sprintf("%s:%s", constant.CmdMigrateTssFunds, "Funds Migrator Admin Cmd"),
			cctx.RelayedMessage,
		)
		require.NotEmpty(t, cctx.InboundParams.Sender)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.InboundParams.CoinType)
		require.Len(t, cctx.OutboundParams, 1)
		require.NotEmpty(t, cctx.OutboundParams[0].Receiver)
		require.EqualValues(t, chains.Ethereum.ChainId, cctx.OutboundParams[0].ReceiverChainId)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.OutboundParams[0].CoinType)
		require.False(t, cctx.OutboundParams[0].Amount.IsZero())
		require.EqualValues(t, gas.EVMSend, cctx.OutboundParams[0].CallOptions.GasLimit)
		require.NotEmpty(t, cctx.OutboundParams[0].GasPrice)
		require.NotEmpty(t, cctx.OutboundParams[0].GasPriorityFee)
	})

	t.Run("returns a new CCTX for migrating funds on Bitcoin", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := chains.BitcoinMainnet.ChainId
		amount := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := sample.Tss()
		newTSSPubkey := sample.Tss()

		// ACT
		cctx, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey.TssPubkey,
			newTSSPubkey.TssPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.NoError(t, err)
		require.NotEmpty(t, cctx.Index)
		require.EqualValues(t, creator, cctx.Creator)
		require.EqualValues(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.EqualValues(
			t,
			fmt.Sprintf("%s:%s", constant.CmdMigrateTssFunds, "Funds Migrator Admin Cmd"),
			cctx.RelayedMessage,
		)
		require.NotEmpty(t, cctx.InboundParams.Sender)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.InboundParams.CoinType)
		require.Len(t, cctx.OutboundParams, 1)
		require.NotEmpty(t, cctx.OutboundParams[0].Receiver)
		require.EqualValues(t, chains.BitcoinMainnet.ChainId, cctx.OutboundParams[0].ReceiverChainId)
		require.EqualValues(t, coin.CoinType_Cmd, cctx.OutboundParams[0].CoinType)
		require.False(t, cctx.OutboundParams[0].Amount.IsZero())
		require.EqualValues(t, uint64(1_000_000), cctx.OutboundParams[0].CallOptions.GasLimit)
		require.NotEmpty(t, cctx.OutboundParams[0].GasPrice)
		require.NotEmpty(t, cctx.OutboundParams[0].GasPriorityFee)
	})

	t.Run("prevent migration with invalid ETH address for current TSS", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := chains.Ethereum.ChainId
		amount := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := "invalid"
		newTSSPubkey := sample.Tss()

		// ACT
		_, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey,
			newTSSPubkey.TssPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("prevent migration with invalid ETH address for new TSS", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := chains.Ethereum.ChainId
		amount := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := sample.Tss()
		newTSSPubkey := "invalid"

		// ACT
		_, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey.TssPubkey,
			newTSSPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("prevent migration on EVM if fees higher than amount", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := chains.Ethereum.ChainId
		amount := sdkmath.NewUint(100_000_000)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := sample.Tss()
		newTSSPubkey := sample.Tss()

		// ACT
		_, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey.TssPubkey,
			newTSSPubkey.TssPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("prevent migration with invalid Bitcoin address for current TSS", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := chains.BitcoinMainnet.ChainId
		amount := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := "invalid"
		newTSSPubkey := sample.Tss()

		// ACT
		_, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey,
			newTSSPubkey.TssPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("prevent migration with invalid Bitcoin address for new TSS", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := chains.BitcoinMainnet.ChainId
		amount := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := sample.Tss()
		newTSSPubkey := "invalid"

		// ACT
		_, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey.TssPubkey,
			newTSSPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.Error(t, err)
	})

	t.Run("prevent migration if invalid chain ID", func(t *testing.T) {
		// ARRANGE
		blockHeight := int64(1000)
		creator := sample.AccAddress()
		inboundHash := sample.Hash().Hex()
		chainID := int64(1000)
		amount := sdkmath.NewUint(1e18)
		medianGasPrice := sdkmath.NewUint(100000)
		priorityFee := sdkmath.NewUint(100000)
		currentTSSPubkey := sample.Tss()
		newTSSPubkey := sample.Tss()

		// ACT
		_, err := types.MigrateFundCmdCCTX(
			blockHeight,
			creator,
			inboundHash,
			chainID,
			amount,
			medianGasPrice,
			priorityFee,
			currentTSSPubkey.TssPubkey,
			newTSSPubkey.TssPubkey,
			[]chains.Chain{},
		)

		// ASSERT
		require.Error(t, err)
	})
}

func TestGetTssMigrationCCTXIndexString(t *testing.T) {
	t.Run("returns unique index string for the CCTX for migrating funds", func(t *testing.T) {
		// ARRANGE
		currentTSSPubkey := sample.PubKeyString()
		newTSSPubkey := sample.PubKeyString()
		chainID := int64(42)
		amount := sdkmath.NewUint(1000)
		height := int64(1000)

		// ACT
		index := types.GetTssMigrationCCTXIndexString(
			currentTSSPubkey,
			newTSSPubkey,
			chainID,
			amount,
			height,
		)
		indexDifferentCurrentTSSPubkey := types.GetTssMigrationCCTXIndexString(
			sample.PubKeyString(),
			newTSSPubkey,
			chainID,
			amount,
			height,
		)
		indexDifferentNewTSSPubkey := types.GetTssMigrationCCTXIndexString(
			currentTSSPubkey,
			sample.PubKeyString(),
			chainID,
			amount,
			height,
		)
		indexDifferentChainID := types.GetTssMigrationCCTXIndexString(
			currentTSSPubkey,
			newTSSPubkey,
			chainID+1,
			amount,
			height,
		)
		indexDifferentAmount := types.GetTssMigrationCCTXIndexString(
			currentTSSPubkey,
			newTSSPubkey,
			chainID,
			sdkmath.NewUint(1001),
			height,
		)
		indexDifferentHeight := types.GetTssMigrationCCTXIndexString(
			currentTSSPubkey,
			newTSSPubkey,
			chainID,
			amount,
			height+1,
		)

		// ASSERT
		require.NotEmpty(t, index)
		require.NotEqual(t, index, indexDifferentCurrentTSSPubkey)
		require.NotEqual(t, index, indexDifferentNewTSSPubkey)
		require.NotEqual(t, index, indexDifferentChainID)
		require.NotEqual(t, index, indexDifferentAmount)
		require.NotEqual(t, index, indexDifferentHeight)
	})
}
