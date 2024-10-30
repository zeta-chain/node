package simulation

import (
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
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
