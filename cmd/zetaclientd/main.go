package main

import (
	ecdsakeygen "github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	//mcconfig "github.com/Meta-Protocol/zetacore/metaclient/config"
	"github.com/cosmos/cosmos-sdk/types"

	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/app"
)

var (
	preParams *ecdsakeygen.LocalPreParams
)

func main() {
	if err := svrcmd.Execute(RootCmd, "", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func SetupConfigForTest() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	//config.SetCoinType(cmd.MetaChainCoinType)
	config.SetFullFundraiserPath(cmd.ZetaChainHDPath)
	types.SetCoinDenomRegex(func() string {
		return cmd.DenomRegex
	})

	rand.Seed(time.Now().UnixNano())

}

func initLogLevel(debug bool) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		log.Info().Msgf("zerolog global log level: DEBUG")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func updateMPIAddress(chain common.Chain, envvar string) {
	mpi := os.Getenv(envvar)
	if mpi != "" {
		config.ChainConfigs[chain.ChainName.String()].ConnectorContractAddress = mpi
		log.Info().Msgf("MPI: %s", mpi)
	}
}

func updateEndpoint(chain common.Chain, envvar string) {
	endpoint := os.Getenv(envvar)
	if endpoint != "" {
		config.ChainConfigs[chain.ChainName.String()].Endpoint = endpoint
		log.Info().Msgf("ENDPOINT: %s", endpoint)
	}
}

func updateTokenAddress(chain common.Chain, envvar string) {
	token := os.Getenv(envvar)
	if token != "" {
		config.ChainConfigs[chain.String()].ZETATokenContractAddress = token
		log.Info().Msgf("TOKEN: %s", token)
	}
}
