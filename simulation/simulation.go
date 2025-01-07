package simulation

import (
	"encoding/json"
	"fmt"
	"os"

	"cosmossdk.io/log"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/zeta-chain/ethermint/app"
	evmante "github.com/zeta-chain/ethermint/app/ante"

	zetaapp "github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
)

func NewSimApp(
	logger log.Logger,
	db dbm.DB,
	appOptions servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*zetaapp.App, error) {
	encCdc := zetaapp.MakeEncodingConfig()

	// Set load latest version to false as we manually set it later.
	zetaApp := zetaapp.New(
		logger,
		db,
		nil,
		false,
		map[int64]bool{},
		app.DefaultNodeHome,
		5,
		encCdc,
		appOptions,
		baseAppOptions...,
	)

	// use zeta antehandler
	options := ante.HandlerOptions{
		AccountKeeper:   zetaApp.AccountKeeper,
		BankKeeper:      zetaApp.BankKeeper,
		EvmKeeper:       zetaApp.EvmKeeper,
		FeeMarketKeeper: zetaApp.FeeMarketKeeper,
		SignModeHandler: encCdc.TxConfig.SignModeHandler(),
		SigGasConsumer:  evmante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:  0,
		ObserverKeeper:  zetaApp.ObserverKeeper,
	}

	anteHandler, err := ante.NewAnteHandler(options)
	if err != nil {
		panic(err)
	}

	zetaApp.SetAnteHandler(anteHandler)
	if err := zetaApp.LoadLatestVersion(); err != nil {
		return nil, err
	}
	return zetaApp, nil
}

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
