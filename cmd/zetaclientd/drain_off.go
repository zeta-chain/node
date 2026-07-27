//go:build !drain

package main

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/zetaclient/maintenance"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
)

// startDrainIfArmed is a no-op in production builds. The emergency drain poller is only
// compiled under the `drain` build tag.
func startDrainIfArmed(
	_ context.Context,
	_ maintenance.ZetacoreClient,
	_ *orchestrator.Orchestrator,
	_ zerolog.Logger,
) {
}
