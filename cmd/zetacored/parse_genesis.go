package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/app"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionsModuleTypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibleModuleTypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

var Copy = map[string]bool{
	slashingtypes.ModuleName:        true,
	govtypes.ModuleName:             true,
	crisistypes.ModuleName:          true,
	feemarkettypes.ModuleName:       true,
	paramstypes.ModuleName:          true,
	upgradetypes.ModuleName:         true,
	evidencetypes.ModuleName:        true,
	vestingtypes.ModuleName:         true,
	fungibleModuleTypes.ModuleName:  true,
	emissionsModuleTypes.ModuleName: true,
	authz.ModuleName:                true,
}
var Skip = map[string]bool{
	evmtypes.ModuleName:          true,
	stakingtypes.ModuleName:      true,
	genutiltypes.ModuleName:      true,
	authtypes.ModuleName:         true,
	banktypes.ModuleName:         true,
	distributiontypes.ModuleName: true,
	group.ModuleName:             true,
}

var Modify = map[string]bool{
	crosschaintypes.ModuleName: true,
	observertypes.ModuleName:   true,
}

func CmdParseGenesisFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse-genesis-file [import-genesis-file] [optional-genesis-file]",
		Short: "Parse the genesis file",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec
			genesisFilePath := filepath.Join(app.DefaultNodeHome, "config", "genesis.json")
			if len(args) == 2 {
				genesisFilePath = args[1]
			}
			genDoc, err := GetGenDoc(genesisFilePath)
			importData, err := GetGenDoc(args[0])

			err = ImportDataIntoFile(genDoc, importData, cdc)
			if err != nil {
				return err
			}

			err = genutil.ExportGenesisFile(genDoc, genesisFilePath)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

func ImportDataIntoFile(genDoc *types.GenesisDoc, importFile *types.GenesisDoc, cdc codec.Codec) error {

	appState, err := genutiltypes.GenesisStateFromGenDoc(*genDoc)
	if err != nil {
		return err
	}
	importAppState, err := genutiltypes.GenesisStateFromGenDoc(*importFile)
	if err != nil {
		return err
	}
	moduleList := app.InitGenesisModuleList()
	for _, m := range moduleList {
		if Skip[m] {
			continue
		}
		if Copy[m] {
			appState[m] = importAppState[m]
		}
		if Modify[m] {
			switch m {
			case crosschaintypes.ModuleName:
				err := ModifyCrossChainState(appState, importAppState, cdc)
				if err != nil {
					return err
				}
			case observertypes.ModuleName:
				err := ModifyObserverState(appState, importAppState, cdc)
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("modify function for %s not found", m)
			}
		}
	}

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}
	genDoc.AppState = appStateJSON

	return nil
}

func ModifyCrossChainState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedCrossChainGenState := crosschaintypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedCrossChainGenState.CrossChainTxs = importedCrossChainGenState.CrossChainTxs[:math.Min(10, len(importedCrossChainGenState.CrossChainTxs))]
	importedCrossChainGenState.InTxHashToCctxList = importedCrossChainGenState.InTxHashToCctxList[:math.Min(10, len(importedCrossChainGenState.InTxHashToCctxList))]
	importedCrossChainGenState.FinalizedInbounds = importedCrossChainGenState.FinalizedInbounds[:math.Min(10, len(importedCrossChainGenState.FinalizedInbounds))]
	importedCrossChainStateBz, err := json.Marshal(importedCrossChainGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal zetacrosschain genesis state: %w", err)
	}
	appState[crosschaintypes.ModuleName] = importedCrossChainStateBz
	return nil
}

func ModifyObserverState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedObserverGenState := observertypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedObserverGenState.Ballots = importedObserverGenState.Ballots[:math.Min(10, len(importedObserverGenState.Ballots))]
	importedObserverGenState.NonceToCctx = importedObserverGenState.NonceToCctx[:math.Min(10, len(importedObserverGenState.NonceToCctx))]

	currentGenState := observertypes.GetGenesisStateFromAppState(cdc, appState)
	currentGenState.Ballots = importedObserverGenState.Ballots
	currentGenState.NonceToCctx = importedObserverGenState.NonceToCctx

	currentGenStateBz, err := cdc.MarshalJSON(&currentGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal observer genesis state: %w", err)
	}

	appState[observertypes.ModuleName] = currentGenStateBz
	return nil
}

func GetGenDoc(fp string) (*types.GenesisDoc, error) {
	path, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	jsonBlob, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	genData, err := types.GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, err
	}
	return genData, nil
}
