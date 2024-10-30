package simulation

import (
	"encoding/json"
	"fmt"
	"os"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/runtime"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

// PrintStats prints the corresponding statistics from the app DB.
func PrintStats(db dbm.DB) {
	fmt.Println("\nDB Stats")
	fmt.Println(db.Stats()["leveldb.stats"])
	fmt.Println("GoLevelDB cached block size", db.Stats()["leveldb.cachedblock"])
}

// CheckExportSimulation exports the app state and simulation parameters to JSON
// if the export paths are defined.
func CheckExportSimulation(app runtime.AppI, config simtypes.Config, params simtypes.Params) error {
	if config.ExportStatePath != "" {
		exported, err := app.ExportAppStateAndValidators(false, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to export app state: %w", err)
		}

		if err := os.WriteFile(config.ExportStatePath, exported.AppState, 0o600); err != nil {
			return err
		}
	}

	if config.ExportParamsPath != "" {
		paramsBz, err := json.MarshalIndent(params, "", " ")
		if err != nil {
			return fmt.Errorf("failed to write app state to %s: %w", config.ExportStatePath, err)
		}

		if err := os.WriteFile(config.ExportParamsPath, paramsBz, 0o600); err != nil {
			return err
		}
	}
	return nil
}
