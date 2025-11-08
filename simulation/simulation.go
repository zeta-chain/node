package simulation

import (
	"encoding/json"
	"fmt"
	"os"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	evmtypes "github.com/cosmos/evm/x/vm/types"

	zetaapp "github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	"github.com/zeta-chain/node/cmd/zetacored/config"
	serverconfig "github.com/zeta-chain/node/server/config"
)

func NewSimApp(
	logger log.Logger,
	db dbm.DB,
	appOptions servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*zetaapp.App, error) {
	configurator := evmtypes.NewEVMConfigurator()
	configurator.ResetTestConfig()
	err := configurator.
		WithEVMCoinInfo(evmtypes.EvmCoinInfo{
			Denom:         config.BaseDenom,
			ExtendedDenom: config.BaseDenom,
			DisplayDenom:  config.BaseDenom,
			Decimals:      config.BaseDenomUnit,
		}).
		Configure()
	if err != nil {
		panic(err)
	}

	// Set load latest version to false as we manually set it later.
	zetaApp := zetaapp.New(
		logger,
		db,
		nil,
		false,
		map[int64]bool{},
		"",
		5,
		serverconfig.DefaultEVMChainID,
		appOptions,
		baseAppOptions...,
	)

	// use zeta antehandler
	encCdc := zetaapp.MakeEncodingConfig(777)
	options := ante.HandlerOptions{
		AccountKeeper:   zetaApp.AccountKeeper,
		BankKeeper:      zetaApp.BankKeeper,
		EvmKeeper:       zetaApp.EvmKeeper,
		FeeMarketKeeper: zetaApp.FeeMarketKeeper,
		SignModeHandler: encCdc.TxConfig.SignModeHandler(),
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
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
