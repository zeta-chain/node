package common

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

type ClientLogger struct {
	Std        zerolog.Logger
	Compliance zerolog.Logger
}

func DefaultLoggers() ClientLogger {
	return ClientLogger{
		Std:        log.Logger,
		Compliance: log.Logger,
	}
}

// IsCctxRestricted returns true if the cctx involves restricted addresses
func IsCctxRestricted(cctx *crosschaintypes.CrossChainTx) bool {
	sender := cctx.InboundTxParams.Sender
	receiver := cctx.GetCurrentOutTxParam().Receiver
	return config.ContainRestrictedAddress(sender, receiver)
}

// PrintComplianceLog prints compliance log with fields [chain, cctx/intx, chain, sender, receiver, token]
func PrintComplianceLog(
	logger1 zerolog.Logger,
	logger2 zerolog.Logger,
	outbound bool,
	chainID int64,
	identifier, sender, receiver, token string) {
	var logMsg string
	var logWithFields1 zerolog.Logger
	var logWithFields2 zerolog.Logger
	if outbound {
		// we print cctx for outbound tx
		logMsg = "Restricted address detected in cctx"
		logWithFields1 = logger1.With().Int64("chain", chainID).Str("cctx", identifier).Str("sender", sender).Str("receiver", receiver).Str("token", token).Logger()
		logWithFields2 = logger2.With().Int64("chain", chainID).Str("cctx", identifier).Str("sender", sender).Str("receiver", receiver).Str("token", token).Logger()
	} else {
		// we print intx for inbound tx
		logMsg = "Restricted address detected in intx"
		logWithFields1 = logger1.With().Int64("chain", chainID).Str("intx", identifier).Str("sender", sender).Str("receiver", receiver).Str("token", token).Logger()
		logWithFields2 = logger2.With().Int64("chain", chainID).Str("intx", identifier).Str("sender", sender).Str("receiver", receiver).Str("token", token).Logger()
	}
	logWithFields1.Warn().Msg(logMsg)
	logWithFields2.Warn().Msg(logMsg)
}
