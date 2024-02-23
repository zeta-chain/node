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
