package context

import (
	"context"
)

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
