package update

import (
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"os"
)

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
