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

// Config returns the config of the app
func (a AppContext) Config() config.Config {
	return a.config
}

// ZetacoreContext returns the context for ZetaChain
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
