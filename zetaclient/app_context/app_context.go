package appcontext

import (
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
)

// AppContext contains global app structs like config, core context and logger
type AppContext struct {
	coreContext *corecontext.ZetaCoreContext
	config      config.Config
}

// NewAppContext creates and returns new AppContext
func NewAppContext(
	coreContext *corecontext.ZetaCoreContext,
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

func (a AppContext) ZetaCoreContext() *corecontext.ZetaCoreContext {
	return a.coreContext
}

// GetBTCChainAndConfig returns btc chain and config if enabled
func (a AppContext) GetBTCChainAndConfig() (common.Chain, config.BTCConfig, bool) {
	btcConfig, configEnabled := a.Config().GetBTCConfig()
	btcChain, _, paramsEnabled := a.ZetaCoreContext().GetBTCChainParams()

	if !configEnabled || !paramsEnabled {
		return common.Chain{}, config.BTCConfig{}, false
	}

	return btcChain, btcConfig, true
}
