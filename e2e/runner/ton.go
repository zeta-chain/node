package runner

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/zeta-chain/node/e2e/runner/ton"
)

// EnsureTONBootstrapped waits unless TON node is bootstrapped.
// - Node should be operational
// - Lite server config exists
// - Faucet is deployed
func (r *E2ERunner) EnsureTONBootstrapped(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(r.Ctx, timeout)
	defer cancel()

	for {
		err := r.Clients.TONSidecar.Status(ctx)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return errors.Wrap(err, "timeout waiting for TON to bootstrap")
		case errors.Is(err, ton.ErrNotHealthy):
			// okay, continue
		case err == nil:
			return nil
		}

		r.Logger.Info("Waiting for TON to bootstrap...")
		time.Sleep(2 * time.Second)
	}
}
