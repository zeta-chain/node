package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"cosmossdk.io/math"
	evidencetypes "cosmossdk.io/x/evidence/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/app"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const MaxItemsForList = 10

// Copy represents a set of modules for which, the entire state is copied without any modifications
var Copy = map[string]bool{
	slashingtypes.ModuleName:  true,
	crisistypes.ModuleName:    true,
	feemarkettypes.ModuleName: true,
	paramstypes.ModuleName:    true,
	upgradetypes.ModuleName:   true,
	evidencetypes.ModuleName:  true,
	vestingtypes.ModuleName:   true,
	fungibletypes.ModuleName:  true,
	emissionstypes.ModuleName: true,
}

// Skip represents a set of modules for which, the entire state is skipped and nothing gets imported
var Skip = map[string]bool{
	// Skipping evm this is done to reduce the size of the genesis file evm module uses the majority of the space due to smart contract data
	evmtypes.ModuleName: true,
	// Skipping staking as new validators would be created for the new chain
	stakingtypes.ModuleName: true,
	// Skipping genutil as new gentxs would be created
	genutiltypes.ModuleName: true,
	// Skipping auth as new accounts would be created for the new chain. This also needs to be done as we are skipping evm module
	authtypes.ModuleName: true,
	// Skipping bank module as it is not used when starting a new chain this is done to make sure the total supply invariant is maintained.
	// This would need modification but might be possible to add in non evm based modules in the future
	banktypes.ModuleName: true,
	// Skipping distribution module as it is not used when starting a new chain , rewards are based on validators and delegators , and so rewards from a different chain do not hold any value
	distributiontypes.ModuleName: true,
	// Skipping group module as it is not used when starting a new chain, new groups should be created based on the validator operator keys
	group.ModuleName: true,
	// Skipping authz as it is not used when starting a new chain, new grants should be created based on the validator hotkeys abd operator keys
	authz.ModuleName: true,
	// Skipping fungible module as new fungible tokens would be created and system contract would be deployed
	fungibletypes.ModuleName: true,
	// Skipping gov types as new parameters are set for the new chain
	govtypes.ModuleName: true,
}

// Modify represents a set of modules for which, the state is modified before importing. Each Module should have a corresponding Modify function
var Modify = map[string]bool{
	observertypes.ModuleName:   true,
	crosschaintypes.ModuleName: true,
}

func CmdParseGenesisFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse-genesis-file [import-genesis-file] [optional-genesis-file]",
		Short: "Parse the provided genesis file and import the required data into the optionally provided genesis file",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec
			modifyEnabled, err := cmd.Flags().GetBool("modify")
			if err != nil {
				return err
			}
			genesisFilePath := filepath.Join(app.DefaultNodeHome, "config", "genesis.json")
			if len(args) == 2 {
				genesisFilePath = args[1]
			}
			_, genesis, err := genutiltypes.GenesisStateFromGenFile(genesisFilePath)
			if err != nil {
				return err
			}
			_, importData, err := genutiltypes.GenesisStateFromGenFile(args[0])
			if err != nil {
				return err
			}
			err = ImportDataIntoFile(genesis, importData, cdc, modifyEnabled)
			if err != nil {
				return err
			}

			err = genutil.ExportGenesisFile(genesis, genesisFilePath)
			if err != nil {
				return err
			}

			return nil
		},
	}
	cmd.PersistentFlags().Bool("modify", false, "modify the genesis file before importing")
	return cmd
}

func ImportDataIntoFile(
	gen *genutiltypes.AppGenesis,
	importFile *genutiltypes.AppGenesis,
	cdc codec.Codec,
	modifyEnabled bool,
) error {
	appState, err := genutiltypes.GenesisStateFromAppGenesis(gen)
	if err != nil {
		return err
	}

	importAppState, err := genutiltypes.GenesisStateFromAppGenesis(importFile)
	if err != nil {
		return err
	}

	moduleList := app.OrderInitGenesis()
	for _, m := range moduleList {
		if Skip[m] {
			continue
		}
		if Copy[m] {
			appState[m] = importAppState[m]
		}
		if Modify[m] && modifyEnabled {
			switch m {
			case crosschaintypes.ModuleName:
				err := ModifyCrosschainState(appState, importAppState, cdc)
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
	gen.AppState = appStateJSON

	return nil
}

// ModifyCrosschainState modifies the crosschain state before importing
// It truncates the crosschain transactions, inbound transactions and finalized inbounds to MaxItemsForList
func ModifyCrosschainState(
	appState map[string]json.RawMessage,
	importAppState map[string]json.RawMessage,
	cdc codec.Codec,
) error {
	importedCrossChainGenState := crosschaintypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedCrossChainGenState.CrossChainTxs = importedCrossChainGenState.CrossChainTxs[:math.Min(MaxItemsForList, len(importedCrossChainGenState.CrossChainTxs))]
	importedCrossChainGenState.InboundHashToCctxList = importedCrossChainGenState.InboundHashToCctxList[:math.Min(MaxItemsForList, len(importedCrossChainGenState.InboundHashToCctxList))]
	importedCrossChainGenState.FinalizedInbounds = importedCrossChainGenState.FinalizedInbounds[:math.Min(MaxItemsForList, len(importedCrossChainGenState.FinalizedInbounds))]
	importedCrossChainStateBz, err := json.Marshal(importedCrossChainGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal zetacrosschain genesis state: %w", err)
	}
	appState[crosschaintypes.ModuleName] = importedCrossChainStateBz
	return nil
}

// ModifyObserverState modifies the observer state before importing
// It truncates the ballots and nonce to cctx list to MaxItemsForList
func ModifyObserverState(
	appState map[string]json.RawMessage,
	importAppState map[string]json.RawMessage,
	cdc codec.Codec,
) error {
	importedObserverGenState := observertypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedObserverGenState.Ballots = importedObserverGenState.Ballots[:math.Min(MaxItemsForList, len(importedObserverGenState.Ballots))]
	importedObserverGenState.NonceToCctx = importedObserverGenState.NonceToCctx[:math.Min(MaxItemsForList, len(importedObserverGenState.NonceToCctx))]

	currentGenState := observertypes.GetGenesisStateFromAppState(cdc, appState)
	currentGenState.Ballots = importedObserverGenState.Ballots
	currentGenState.NonceToCctx = importedObserverGenState.NonceToCctx
	// zero operational flags as they are network specific
	currentGenState.OperationalFlags = observertypes.OperationalFlags{}

	currentGenStateBz, err := cdc.MarshalJSON(&currentGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal observer genesis state: %w", err)
	}

	appState[observertypes.ModuleName] = currentGenStateBz
	return nil
}
