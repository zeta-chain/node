package mocks

import (
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/maintenance"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// Every module in the project that uses the zetacore client specifies its own ZetacoreClient
// interface. The zetacoreClient interface here is used to generate a mock that works for those
// modules.
//
//go:generate mockery --name zetacoreClient --structname ZetacoreClient --filename zetacore.go --output ../
//nolint:unused // used for code gen
type zetacoreClient interface {
	zrepo.ZetacoreClient
	orchestrator.ZetacoreClient
	maintenance.ZetacoreClient
}

var _ zetacoreClient = &zetacore.Client{}
