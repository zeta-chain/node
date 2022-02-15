package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
)

var ALL_CHAINS = []*ChainETHish{
	{
		chain:                        common.Chain("ETH"),
		contract:                     "0xDe8802902Ff3136bdACe5FFC9a2423B1d37F6833",
		DEFAULT_DESTINATION_CONTRACT: "0xFf6B270ac3790589A1Fe90d0303e9D4d9A54FD1A",
	},
	{
		chain:                        common.Chain("BSC"),
		contract:                     "0xCC3e1C9460B7803d4d79F32342b2b27543362536",
		DEFAULT_DESTINATION_CONTRACT: "0xF47bd84B86d1667e7621c38c72C6905Ca8710b0d",
	},
}

func startAllChainListeners() {
	for _, chain := range ALL_CHAINS {
		chain.Init()
		chain.Listen()
	}
}

func main() {
	startAllChainListeners()

	ch3 := make(chan os.Signal, 1)
	signal.Notify(ch3, syscall.SIGINT, syscall.SIGTERM)
	<-ch3
	log.Info().Msg("stop signal received")
}
