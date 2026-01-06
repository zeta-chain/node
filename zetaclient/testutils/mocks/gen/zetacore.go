package mocks

import (
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/maintenance"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

//go:generate mockery --name zetacoreClient --structname ZetacoreClient --filename zetacore.go --output ../

// Every module in the project that uses the zetacore client specifies its own ZetacoreClient
// interface. The zetacoreClient interface defined here is used to generate a mock that works for
// all those modules.
//
// The interface is unexported on purpose, since we ONLY use it for mock generation.
//
//nolint:unused // used for code gen
type zetacoreClient interface {
	zrepo.ZetacoreClient
	orchestrator.ZetacoreClient
	maintenance.ZetacoreClient
}

var _ zetacoreClient = &zetacore.Client{}
