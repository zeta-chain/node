package context

import (
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// AppContext contains global app structs like config, core context and logger
type AppContext struct {
	coreContext *ZetaCoreContext
	config      config.Config
}

// NewAppContext creates and returns new AppContext
func NewAppContext(
	coreContext *ZetaCoreContext,
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

func (a AppContext) ZetaCoreContext() *ZetaCoreContext {
	return a.coreContext
}

// GetBTCChainAndConfig returns btc chain and config if enabled
func (a AppContext) GetBTCChainAndConfig() (chains.Chain, config.BTCConfig, bool) {
	btcConfig, configEnabled := a.Config().GetBTCConfig()
	btcChain, _, paramsEnabled := a.ZetaCoreContext().GetBTCChainParams()

	if !configEnabled || !paramsEnabled {
		return chains.Chain{}, config.BTCConfig{}, false
	}

	return btcChain, btcConfig, true
}
