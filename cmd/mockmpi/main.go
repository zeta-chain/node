package main

import (
	"flag"
	"fmt"
	"github.com/zeta-chain/zetacore/cmd/mockmpi/common"
	"github.com/zeta-chain/zetacore/cmd/mockmpi/eth"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func startAllChainListeners() {
	for _, chain := range common.ALL_CHAINS {
		chain.Start()
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	debug := flag.Bool("debug", false, "sets log level to debug")
	onlyChain := flag.String("only-chain", "all", "Uppercase name of a supported chain")
	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	eth.RegisterChains()

	if *onlyChain == "all" {
		log.Info().Msg("Starting all chains")
		startAllChainListeners()
	} else {
		log.Info().Msg(fmt.Sprintf("Running 1 chain only: %s", *onlyChain))
		chain, err := common.FindChainByName(*onlyChain)
		if err != nil {
			log.Fatal().Err(err).Msg("Chain not found")
			os.Exit(1)
		}
		chain.Start()
	}

	ch3 := make(chan os.Signal, 1)
	signal.Notify(ch3, syscall.SIGINT, syscall.SIGTERM)
	<-ch3
	log.Info().Msg("stop signal received")
}
