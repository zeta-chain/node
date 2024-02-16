package appcontext

import (
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
)

type AppContext struct {
	coreContext *corecontext.ZeraCoreContext
	config      *config.Config
}

func NewAppContext(coreContext *corecontext.ZeraCoreContext, config *config.Config) *AppContext {
	return &AppContext{
		coreContext: coreContext,
		config:      config,
	}
}

func (a *AppContext) Config() *config.Config {
	return a.config
}

func (a *AppContext) ZetaCoreContext() *corecontext.ZeraCoreContext {
	return a.coreContext
}
