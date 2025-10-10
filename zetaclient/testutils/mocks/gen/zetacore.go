package mocks

import (
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/maintenance"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

//go:generate mockery --name zetacoreClient --structname ZetacoreClient --filename zetacore.go --output ../
//nolint:unused // used for code gen
type zetacoreClient interface {
	zrepo.ZetacoreClient
	orchestrator.Zetacore
	maintenance.ZetacoreClient
}

var _ zetacoreClient = &zetacore.Client{}
