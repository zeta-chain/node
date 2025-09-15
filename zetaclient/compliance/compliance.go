// Package compliance provides functions to check for compliance of cross-chain transactions
package compliance

import (
	"github.com/rs/zerolog"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// IsCCTXRestricted returns true if the cctx involves restricted addresses
func IsCCTXRestricted(cctx *crosschaintypes.CrossChainTx, additionalAddresses ...string) bool {
	additionalAddresses = append(
		additionalAddresses,
		cctx.InboundParams.Sender,
		cctx.GetCurrentOutboundParam().Receiver,
	)

	return config.ContainRestrictedAddress(additionalAddresses...)
}

// PrintComplianceLog prints compliance log with fields [chain, cctx/inbound, chain, sender, receiver, token]
func PrintComplianceLog(
	logger, complianceLogger zerolog.Logger,
	outbound bool,
	chainID int64,
	identifier, sender, receiver, token string,
) {
	var (
		message string
		fields  map[string]any
	)

	if outbound {
		message = "Restricted address detected in cctx"
		fields = map[string]any{
			logs.FieldChain:    chainID,
			logs.FieldCoinType: token,
			"identifier":       identifier,
			"sender":           sender,
			"receiver":         receiver,
		}
	} else {
		message = "Restricted address detected in inbound"
		fields = map[string]any{
			logs.FieldChain:    chainID,
			logs.FieldCoinType: token,
			logs.FieldTx:       identifier,
			"sender":           sender,
			"receiver":         receiver,
		}
	}

	logger.Warn().Fields(fields).Msg(message)
	complianceLogger.Warn().Fields(fields).Msg(message)
}
