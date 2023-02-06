package main

import (
	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	RPCHost        string `env:"BTC_RPC_URL" envDefault:"107.20.255.203:18332"`
	RPCUser        string `env:"BTC_RPC_USER" envDefault:"user"`
	RPCPass        string `env:"BTC_RPC_PASS" envDefault:"pass"`
	WalletAddress  string `env:"BTC_WALLET_ADDRESS" envDefault:"tb1q8ev0a9c0khvumur5w6dw9szuzk9a6f7lh6jlhz"`
	WalletPK       string `env:"BTC_WALLET_PK,required"`
	MinConf        int    `env:"BTC_MIN_CONF" envDefault:"6"`
	MaxConf        int    `env:"BTC_MAX_CONF" envDefault:"10000"`
	TickerInterval int    `env:"TICKER_INTERVAL_SEG" envDefault:"60"`
	PrevCount      int    `env:"BTC_PREV_COUNT" envDefault:"10"` // how many utxos include on each consolidated Tx
}

// MustGetConfig returns a configuration
func MustGetConfig() *Configuration {
	// These are the default values, can be overridden by env vars above
	cfg := &Configuration{}
	err := env.Parse(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
