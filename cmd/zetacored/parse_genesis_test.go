package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/stretchr/testify/require"

	zetacored "github.com/zeta-chain/node/cmd/zetacored"
	_ "github.com/zeta-chain/node/pkg/sdkconfig/default"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func Test_ModifyCrossChainState(t *testing.T) {
	t.Run("successfully modify cross chain state to reduce data", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appState := sample.AppState(t)
		importData := GetImportData(t, cdc, 100)
		err := zetacored.ModifyCrosschainState(appState, importData, cdc)
		require.NoError(t, err)

		modifiedCrosschainAppState := crosschaintypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, modifiedCrosschainAppState.CrossChainTxs, 10)
		require.Len(t, modifiedCrosschainAppState.InboundHashToCctxList, 10)
		require.Len(t, modifiedCrosschainAppState.FinalizedInbounds, 10)
	})

	t.Run("successfully modify cross chain state without changing data when not needed", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appState := sample.AppState(t)
		importData := GetImportData(t, cdc, 8)
		err := zetacored.ModifyCrosschainState(appState, importData, cdc)
		require.NoError(t, err)

		modifiedCrosschainAppState := crosschaintypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, modifiedCrosschainAppState.CrossChainTxs, 8)
		require.Len(t, modifiedCrosschainAppState.InboundHashToCctxList, 8)
		require.Len(t, modifiedCrosschainAppState.FinalizedInbounds, 8)
	})
}

func Test_ModifyObserverState(t *testing.T) {
	t.Run("successfully modify observer state to reduce data", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appState := sample.AppState(t)
		importData := GetImportData(t, cdc, 100)
		err := zetacored.ModifyObserverState(appState, importData, cdc)
		require.NoError(t, err)

		modifiedObserverAppState := observertypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, modifiedObserverAppState.Ballots, zetacored.MaxItemsForList)
		require.Len(t, modifiedObserverAppState.NonceToCctx, zetacored.MaxItemsForList)
	})

	t.Run("successfully modify observer state without changing data when not needed", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appState := sample.AppState(t)
		importData := GetImportData(t, cdc, 8)
		err := zetacored.ModifyObserverState(appState, importData, cdc)
		require.NoError(t, err)

		modifiedObserverAppState := observertypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, modifiedObserverAppState.Ballots, 8)
		require.Len(t, modifiedObserverAppState.NonceToCctx, 8)
	})

}

func Test_ImportDataIntoFile(t *testing.T) {
	t.Run("successfully import data into file and modify data", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appGenesis := sample.AppGenesis(t)
		importAppGenesis := ImportAppGenesis(t, cdc, 100)

		err := zetacored.ImportDataIntoFile(appGenesis, importAppGenesis, cdc, true)
		require.NoError(t, err)

		types.GenesisDocFromJSON(appGenesis.AppState)

		var appState map[string]json.RawMessage
		err = json.Unmarshal(appGenesis.AppState, &appState)
		require.NoError(t, err)

		// Crosschain module is in Modify list
		crosschainStateAfterImport := crosschaintypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, crosschainStateAfterImport.CrossChainTxs, zetacored.MaxItemsForList)
		require.Len(t, crosschainStateAfterImport.InboundHashToCctxList, zetacored.MaxItemsForList)
		require.Len(t, crosschainStateAfterImport.FinalizedInbounds, zetacored.MaxItemsForList)

		// Bank module is in Skip list
		var bankStateAfterImport banktypes.GenesisState
		if appState[banktypes.ModuleName] != nil {
			err := cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankStateAfterImport)
			if err != nil {
				panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
			}
		}
		// 4 balances were present in the original genesis state
		require.Len(t, bankStateAfterImport.Balances, 4)

		// Emissions module is in Copy list
		var emissionStateAfterImport emissionstypes.GenesisState
		if appState[emissionstypes.ModuleName] != nil {
			err := cdc.UnmarshalJSON(appState[emissionstypes.ModuleName], &emissionStateAfterImport)
			if err != nil {
				panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
			}
		}
		require.Len(t, emissionStateAfterImport.WithdrawableEmissions, 100)
	})

	t.Run("successfully import data into file without modifying data", func(t *testing.T) {
		cdc := keepertest.NewCodec()
		appGenesis := sample.AppGenesis(t)
		importAppGenesis := ImportAppGenesis(t, cdc, 8)

		err := zetacored.ImportDataIntoFile(appGenesis, importAppGenesis, cdc, false)
		require.NoError(t, err)

		var appState map[string]json.RawMessage
		err = json.Unmarshal(appGenesis.AppState, &appState)
		require.NoError(t, err)

		// Crosschain module is in Modify list
		crosschainStateAfterImport := crosschaintypes.GetGenesisStateFromAppState(cdc, appState)
		require.Len(t, crosschainStateAfterImport.CrossChainTxs, 0)
		require.Len(t, crosschainStateAfterImport.InboundHashToCctxList, 0)
		require.Len(t, crosschainStateAfterImport.FinalizedInbounds, 0)

		// Bank module is in Skip list
		var bankStateAfterImport banktypes.GenesisState
		if appState[banktypes.ModuleName] != nil {
			err := cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankStateAfterImport)
			if err != nil {
				panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
			}
		}
		// 4 balances were present in the original genesis state
		require.Len(t, bankStateAfterImport.Balances, 4)

		// Emissions module is in Copy list
		var emissionStateAfterImport emissionstypes.GenesisState
		if appState[emissionstypes.ModuleName] != nil {
			err := cdc.UnmarshalJSON(appState[emissionstypes.ModuleName], &emissionStateAfterImport)
			if err != nil {
				panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
			}
		}
		require.Len(t, emissionStateAfterImport.WithdrawableEmissions, 8)

	})
}

func ImportAppGenesis(t *testing.T, cdc *codec.ProtoCodec, n int) *genutiltypes.AppGenesis {
	importAppGenesis := sample.AppGenesis(t)
	importStateJson, err := json.Marshal(GetImportData(t, cdc, n))
	require.NoError(t, err)
	importAppGenesis.AppState = importStateJson
	return importAppGenesis
}

func GetImportData(t *testing.T, cdc *codec.ProtoCodec, n int) map[string]json.RawMessage {
	importData := sample.AppState(t)

	// Add crosschain data to genesis state
	importedCrossChainGenState := crosschaintypes.GetGenesisStateFromAppState(cdc, importData)
	cctxList := make([]*crosschaintypes.CrossChainTx, n)
	inboundHashToCctxList := make([]crosschaintypes.InboundHashToCctx, n)
	finalLizedInbounds := make([]string, n)
	for i := 0; i < n; i++ {
		cctxList[i] = sample.CrossChainTx(t, fmt.Sprintf("crosschain-%d", i))
		inboundHashToCctxList[i] = sample.InboundHashToCctx(t, fmt.Sprintf("inboundHashToCctxList-%d", i))
		finalLizedInbounds[i] = fmt.Sprintf("finalLizedInbounds-%d", i)
	}
	importedCrossChainGenState.CrossChainTxs = cctxList
	importedCrossChainGenState.InboundHashToCctxList = inboundHashToCctxList
	importedCrossChainGenState.FinalizedInbounds = finalLizedInbounds
	importedCrossChainStateBz, err := json.Marshal(importedCrossChainGenState)
	require.NoError(t, err)
	importData[crosschaintypes.ModuleName] = importedCrossChainStateBz

	// Add observer data to genesis state
	importedObserverGenState := observertypes.GetGenesisStateFromAppState(cdc, importData)
	ballots := make([]*observertypes.Ballot, n)
	nonceToCctx := make([]observertypes.NonceToCctx, n)
	for i := 0; i < n; i++ {
		ballots[i] = sample.Ballot(t, fmt.Sprintf("ballots-%d", i))
		nonceToCctx[i] = sample.NonceToCCTX(t, fmt.Sprintf("nonceToCctx-%d", i))
	}
	importedObserverGenState.Ballots = ballots
	importedObserverGenState.NonceToCctx = nonceToCctx
	importedObserverStateBz, err := cdc.MarshalJSON(&importedObserverGenState)
	require.NoError(t, err)
	importData[observertypes.ModuleName] = importedObserverStateBz

	// Add emission data to genesis state
	var importedEmissionGenesis emissionstypes.GenesisState
	if importData[emissionstypes.ModuleName] != nil {
		err := cdc.UnmarshalJSON(importData[emissionstypes.ModuleName], &importedEmissionGenesis)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	withdrawableEmissions := make([]emissionstypes.WithdrawableEmissions, n)
	for i := 0; i < n; i++ {
		withdrawableEmissions[i] = sample.WithdrawableEmissions(t)
	}
	importedEmissionGenesis.WithdrawableEmissions = withdrawableEmissions
	importedEmissionGenesisBz, err := cdc.MarshalJSON(&importedEmissionGenesis)
	require.NoError(t, err)
	importData[emissionstypes.ModuleName] = importedEmissionGenesisBz

	// Add bank data to genesis state
	var importedBankGenesis banktypes.GenesisState
	if importData[banktypes.ModuleName] != nil {
		err := cdc.UnmarshalJSON(importData[banktypes.ModuleName], &importedBankGenesis)
		if err != nil {
			panic(fmt.Sprintf("Failed to get genesis state from app state: %s", err.Error()))
		}
	}
	balances := make([]banktypes.Balance, n)
	supply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.ZeroInt()))
	for i := 0; i < n; i++ {
		balances[i] = banktypes.Balance{
			Address: sample.AccAddress(),
			Coins:   sample.Coins(),
		}
		supply = supply.Add(balances[i].Coins...)
	}
	importedBankGenesis.Balances = balances
	importedBankGenesis.Supply = supply
	importedBankGenesisBz, err := cdc.MarshalJSON(&importedBankGenesis)
	require.NoError(t, err)
	importData[banktypes.ModuleName] = importedBankGenesisBz

	return importData
}
