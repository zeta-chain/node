package base

import (
	"github.com/zeta-chain/node/zetaclient/context"
)

// CheckSkipInbound returns true if inbound related observations should be skipped.
func CheckSkipInbound(ob *Observer, app *context.AppContext) bool {
	isSupported := ob.ChainParams().IsSupported
	isInboundEnabled := app.IsInboundObservationEnabled()
	isMempoolCongested := app.IsMempoolCongested()
	isMaxFeeExceeded := app.IsMaxFeeExceeded()

	if !isSupported || !isInboundEnabled || isMempoolCongested || isMaxFeeExceeded {
		ob.Logger().
			Sampled.Info().
			Bool("is_supported", isSupported).
			Bool("is_enabled", isInboundEnabled).
			Bool("is_congested", isMempoolCongested).
			Bool("is_max_fee_exceeded", isMaxFeeExceeded).
			Msg("skip inbound observation")
		return true
	}
	return false
}

// CheckSkipOutbound returns true if outbound related observations should be skipped.
func CheckSkipOutbound(ob *Observer, app *context.AppContext) bool {
	isSupported := ob.ChainParams().IsSupported
	isOutboundEnabled := app.IsOutboundObservationEnabled()
	isMempoolCongested := app.IsMempoolCongested()

	if !isSupported || !isOutboundEnabled || isMempoolCongested {
		ob.Logger().
			Sampled.Info().
			Bool("is_supported", isSupported).
			Bool("is_enabled", isOutboundEnabled).
			Bool("is_congested", isMempoolCongested).
			Msg("skip outbound observation")
		return true
	}
	return false
}

// CheckSkipGasPrice returns true if gas price observation should be skipped.
func CheckSkipGasPrice(ob *Observer, app *context.AppContext) bool {
	isSupported := ob.ChainParams().IsSupported
	isMempoolCongested := app.IsMempoolCongested()
	isMaxFeeExceeded := app.IsMaxFeeExceeded()

	if !isSupported || isMempoolCongested || isMaxFeeExceeded {
		ob.Logger().
			Sampled.Info().
			Bool("is_supported", isSupported).
			Bool("is_congested", isMempoolCongested).
			Bool("is_max_fee_exceeded", isMaxFeeExceeded).
			Msg("skip gas price observation")
		return true
	}
	return false
}
