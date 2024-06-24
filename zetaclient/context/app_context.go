package context

import (
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// AppContext contains global app structs like config, zetacore context and logger
type AppContext struct {
	coreContext *ZetacoreContext
	config      config.Config
}

// NewAppContext creates and returns new AppContext
func NewAppContext(
	coreContext *ZetacoreContext,
	config config.Config,
) *AppContext {
	return &AppContext{
		coreContext: coreContext,
		config:      config,
	}
}

func (a AppContext) Config() config.Config {
	return a.config
}

func (a AppContext) ZetacoreContext() *ZetacoreContext {
	return a.coreContext
}

// GetBTCChainAndConfig returns btc chain and config if enabled
func (a AppContext) GetBTCChainAndConfig() (chains.Chain, config.BTCConfig, bool) {
	btcConfig, configEnabled := a.Config().GetBTCConfig()
	btcChain, _, paramsEnabled := a.coreContext.GetBTCChainParams()

	if !configEnabled || !paramsEnabled {
		return chains.Chain{}, config.BTCConfig{}, false
	}

	return btcChain, btcConfig, true
}

// GetBTCChainID returns btc chain id if enabled or regnet chain id by default
// Bitcoin chain ID is currently needed by TSS to calculate the correct Bitcoin address
// TODO: we might have multiple BTC chains in the future: https://github.com/zeta-chain/node/issues/1397
func (a AppContext) GetBTCChainID() int64 {
	bitcoinChainID := chains.BitcoinRegtest.ChainId
	btcChain, _, enabled := a.GetBTCChainAndConfig()
	if enabled {
		bitcoinChainID = btcChain.ChainId
	}
	return bitcoinChainID
}
