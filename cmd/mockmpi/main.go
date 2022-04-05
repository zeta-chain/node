package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
)

var ALL_CHAINS = []*ChainETHish{
	{
		name:         common.Chain("ETH"),
		MPI_CONTRACT: "0x4740f4051eA6D896C694303228D86Ba3141065ca",
		chain_id:     5,
	},
	{
		name:         common.Chain("BSC"),
		MPI_CONTRACT: "0x4a2d53e16ebe3feC54B407c9e29590951Ce2b6ad",
		chain_id:     97,
	},
	{
		name:         common.Chain("POLYGON"),
		MPI_CONTRACT: "0xD9D3f57800033a1b403c62927398E97FA2Ce0c24",
		chain_id:     8001,
	},
}

func startAllChainListeners() {
	for _, chain := range ALL_CHAINS {
		chain.Start()
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ethEndPoint := os.Getenv("ETH_ENDPOINT")
	if ethEndPoint != "" {
		config.ETH_ENDPOINT = ethEndPoint
		log.Info().Msgf("ETH_ENDPOINT: %s", ethEndPoint)
	}
	bscEndPoint := os.Getenv("BSC_ENDPOINT")
	if bscEndPoint != "" {
		config.BSC_ENDPOINT = bscEndPoint
		log.Info().Msgf("BSC_ENDPOINT: %s", bscEndPoint)
	}
	polygonEndPoint := os.Getenv("POLYGON_ENDPOINT")
	if polygonEndPoint != "" {
		config.POLY_ENDPOINT = polygonEndPoint
		log.Info().Msgf("POLYGON_ENDPOINT: %s", polygonEndPoint)
	}

	var logZetaSentSignature = []byte("ZetaSent(address,uint16,bytes,uint256,uint256,bytes,bytes)")
	logZetaSentSignatureHash := crypto.Keccak256Hash(logZetaSentSignature)
	MAGIC_HASH = logZetaSentSignatureHash.String()
	log.Info().Msgf("Magic Hash: %s", MAGIC_HASH)

	debug := flag.Bool("debug", false, "sets log level to debug")
	onlyChain := flag.String("only-chain", "all", "Uppercase name of a supported chain")
	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if *onlyChain == "all" {
		log.Info().Msg("Starting all chains")
		startAllChainListeners()
	} else {
		log.Info().Msg(fmt.Sprintf("Running 1 chain only: %s", *onlyChain))
		chain, err := FindChainByName(*onlyChain)
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
