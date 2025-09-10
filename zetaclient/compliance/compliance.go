// Package compliance provides functions to check for compliance of cross-chain transactions
package compliance

import (
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
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

// PrintComplianceLog prints compliance log with fields
// [chain, sender, receiver, coin_type, cctx/tx] (coinType is optional)
func PrintComplianceLog(
	logger, complianceLogger zerolog.Logger,
	outbound bool,
	chainID int64,
	identifier, sender, receiver string,
	coinType *coin.CoinType,
) {
	var message string
	fields := map[string]any{
		logs.FieldChain: chainID,
		"sender":        sender,
		"receiver":      receiver,
	}

	if coinType != nil {
		fields[logs.FieldCoinType] = *coinType
	}

	if outbound {
		message = "restricted address detected in CCTX"
		fields["indentifier"] = identifier
	} else {
		message = "restricted address detected in inbound"
		fields[logs.FieldTx] = identifier
	}

	logger.Warn().Fields(fields).Msg(message)
	complianceLogger.Warn().Fields(fields).Msg(message)
}
