package context

import (
	"context"
)

// TODO: https://github.com/zeta-chain/node/issues/4292
// EnableMultipleCallsFeatureFlag returns true if EnableMultipleCalls feature flag is enabled
func EnableMultipleCallsFeatureFlag(ctx context.Context) bool {
	app, err := FromContext(ctx)
	if err != nil {
		app.logger.Warn().Err(err).
			Msg("unable to get feature flag, using default behavior")
		return false
	}

	return app.Config().FeatureFlags.EnableMultipleCalls
}

// EnableSolanaAddressLookupTable returns true if EnableSolanaAddressLookupTable feature flag is enabled
func EnableSolanaAddressLookupTableFeatureFlag(ctx context.Context) bool {
	app, err := FromContext(ctx)
	if err != nil {
		app.logger.Warn().Err(err).
			Msg("unable to get feature flag, using default behavior")
		return false
	}

	return app.Config().FeatureFlags.EnableSolanaAddressLookupTable
}
