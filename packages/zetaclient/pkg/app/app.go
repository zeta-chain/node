package app

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/core"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/observer"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/signer"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/config"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/logger"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/model"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ZetaclientApp struct {
	cfg          *config.Configuration
	coreObserver *core.Observer
	signers      map[model.Chain]signer.Signer
	chains       map[model.Chain]observer.Observer
	bridge       bridge.Adapter
	log          logger.Logger
}

func NewZetaclientApp(cfg *config.Configuration, coreObserver *core.Observer, signers map[model.Chain]signer.Signer, chains map[model.Chain]observer.Observer, bridge bridge.Adapter, log logger.Logger) {
	return &ZetaClientApp{
		cfg:          cfg,
		coreObserver: coreObserver,
		signers:      signers,
		chains:       chains,
		bridge:       bridge.Adapter,
		log:          logger,
	}
}

func (app *ZetaclientApp) Start() {
	app.setupConfigForTest()
	app.waitForZetacore()

	// per chain setup
	for chain, observer := range app.chains {
		signer, ok := signers[chain]
		if !ok {
			app.log.Infow("setup - no signer for chain", "chain", chain)
			continue
		}
		app.setupTSS(chain, signer)
		// start chain observer
		observer.Start()
	}

	log.Infow("starting zetacore observer...")
	coreObserver.MonitorCore()

	// wait....
	log.Infow("awaiting the os.Interrupt, syscall.SIGTERM signals...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Infow("stop signal received", "signal", sig)

	// stop zetacore observer
	for _, chain := range app.chains {
		chain.Stop()
	}
}

func (app *ZetaclientApp) setupTSS(chain model.Chain, sgn signer.Signer) {
	// Per chain initialization
	zetaTx, err := app.bridge.SetTSS(chain, sgn.Address(), sgn.CurrentPubKey())
	if err != nil {
		app.log.Errorw("Set TSS failed", "error", err.Error())
		continue
	}
	app.log.Infow("Set TSS ", "to", zetaTx, "chain", chain, "address", signer.Address())
}

func (app *ZetaclientApp) waitForZetacore() {
	// wait until zetacore is up
	app.log.Infow("Waiting for ZetaCore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("%s:9090", app.cfg.ChainIP),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	app.log.Infow("ZetaCore to open 9090 port...")
}

func (app *ZetaclientApp) setupConfigForTest() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(model.Bech32PrefixAccAddr, model.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(model.Bech32PrefixValAddr, model.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(model.Bech32PrefixConsAddr, model.Bech32PrefixConsPub)
	//config.SetCoinType(cmd.MetaChainCoinType)
	config.SetFullFundraiserPath(model.ZetaChainHDPath)
	types.SetCoinDenomRegex(func() string {
		return model.DenomRegex
	})

	rand.Seed(time.Now().UnixNano())
}
