package errors

import (
	"fmt"

	"github.com/zeta-chain/node/x/crosschain/types"
)

// ErrTxMonitor represents an error from the inbound vote monitoring goroutine.
type ErrTxMonitor struct {
	Err        error
	ZetaTxHash string
	Msg        types.MsgVoteInbound
}

func (m ErrTxMonitor) Error() string {
	if m.Err == nil {
		return "monitoring completed without error"
	}
	return fmt.Sprintf("monitoring error: %v, inbound hash: %s, zeta tx hash: %s, ballot index: %s",
		m.Err, m.Msg.InboundHash, m.ZetaTxHash, m.Msg.Digest())
}
