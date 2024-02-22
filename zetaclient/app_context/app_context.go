package appcontext

import (
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
)

type AppContext struct {
	coreContext *corecontext.ZetaCoreContext
	config      *config.Config
	logger      zerolog.Logger
}

// NewAppContext creates a new AppContext, containing global app structs
// like config, core context and logger
func NewAppContext(
	coreContext *corecontext.ZetaCoreContext,
	config *config.Config,
	logger zerolog.Logger,
) *AppContext {
	return &AppContext{
		coreContext: coreContext,
		config:      config,
		logger:      logger,
	}
}

func (a *AppContext) Config() *config.Config {
	return a.config
}

func (a *AppContext) ZetaCoreContext() *corecontext.ZetaCoreContext {
	return a.coreContext
}

func (a *AppContext) Logger() zerolog.Logger {
	return a.logger
}

// GetBTCChainAndConfig returns btc chain and config if enabled
func (a *AppContext) GetBTCChainAndConfig() (common.Chain, config.BTCConfig, bool) {
	btcConfig, configEnabled := a.Config().GetBTCConfig()
	btcChain, _, paramsEnabled := a.ZetaCoreContext().GetBTCChainParams()

	if !configEnabled || !paramsEnabled {
		return common.Chain{}, config.BTCConfig{}, false
	}

	return btcChain, btcConfig, true
}
