package main

import (
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"math/big"
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
		MPI_CONTRACT: "0x68Bc806414e743D88436AEB771Be387A55B4df70",
		chain_id:     big.NewInt(5),
	},
	{
		name:         common.Chain("BSC"),
		MPI_CONTRACT: "0xE626402550fB921E4a47c11568F89dF3496fbEF0",
		chain_id:     big.NewInt(97),
	},
	{
		name:         common.Chain("POLYGON"),
		MPI_CONTRACT: "0x18A276F4ecF6B788F805EF265F89C521401B1815",
		chain_id:     big.NewInt(80001),
	},
}

func startAllChainListeners() {
	for _, chain := range ALL_CHAINS {
		chain.Start()
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ethEndPoint := os.Getenv("GOERLI_ENDPOINT")
	if ethEndPoint != "" {
		config.Chains[common.GoerliChain.String()].Endpoint = ethEndPoint
	}
	log.Info().Msgf("GOERLI_ENDPOINT: %s", ethEndPoint)

	bscEndPoint := os.Getenv("BSCTESTNET_ENDPOINT")
	if bscEndPoint != "" {
		config.Chains[common.BSCTestnetChain.String()].Endpoint = bscEndPoint
	}
	log.Info().Msgf("BSCTESTNET_ENDPOINT: %s", bscEndPoint)

	polygonEndPoint := os.Getenv("MUMBAI_ENDPOINT")
	if polygonEndPoint != "" {
		config.Chains[common.MumbaiChain.String()].Endpoint = polygonEndPoint
	}
	log.Info().Msgf("MUMBAI_ENDPOINT: %s", polygonEndPoint)

	var logZetaSentSignature = []byte("ZetaSent(address,uint256,bytes,uint256,uint256,bytes,bytes)")
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
