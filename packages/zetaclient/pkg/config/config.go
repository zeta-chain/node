package config

import (
	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	EnabledChains []string    `env:"ENABLED_CHAINS" envDefault:"GOERLI,BSCTESTNET"`
	ValidatorName string      `env:"VALIDATOR_NAME" envDefault:"alice"`
	PeerAddress   string      `env:"PEER_ADDRESS" envDefault:""`
	LogConsole    bool        `env:"LOG_CONSOLE" envDefault:"false"`
	PreParamsPath bool        `env:"PRE_PARAMS_PATH" envDefault:""`
	NoKeygen      bool        `env:"NO_KEYGEN" envDefault:"false"`
	ZetaCoreHome  string      `env:"ZETA_CORE_HOME" envDefault:".zetacored"`
	KeygenBlock   int64       `env:"KEYGEN_BLOCK" envDefault:"0"`
	ChainIP       string      `env:"CHAIN_IP" envDefault:"127.0.0.1"`
	Goerli        ChainConfig `envPrefix:"GOERLI_"`
	Ropsten       ChainConfig `envPrefix:"ROPSTEN_"`
	Mumbai        ChainConfig `envPrefix:"MUMBAI_"`
	Baobab        ChainConfig `envPrefix:"BAOBAB_"`
	Bitcoin       ChainConfig `envPrefix:"BITCOIN_"`
}

type ChainConfig struct {
	Endpoint    string `env:"ENDPOINT"`
	MPIAddress  string `env:"MPI_ADDRESS"`
	PoolAddress string `env:"POOL_ADDRESS"`
	ZetaAddress string `env:"ZETA_ADDRESS"`
}

// MustGetConfig returns a configuration
func MustGetConfig() *Configuration {
	// These are the default values, can be overriden by env vars above
	cfg := &Configuration{}
	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
