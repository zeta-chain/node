package main

import (
	coreInfra "github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/core/infra"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/signer"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/app"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/config"

	"go.uber.org/zap"
)

func main() {
	var application *app.ZetaclientApp
	cfg := config.MustConfig()

	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	//signer :=
	//bridge :=
	sugarLogger := logger.Sugar()

	signers := make(map[model.Chain]signer.Signer, len(cfg.EnabledChains))
	chains := make(map[model.Chain]chain.Adapter, len(cfg.EnabledChains))

	coreObserver := coreInfra.NewCoreObserver(cfg, bridge, signers, chains, metrics, signer, suggarLogger)
	application = app.NewZetaclientApp(cfg, coreObserver, signers, chains, bridge, logger.Sugar())

	application.Start()
}
