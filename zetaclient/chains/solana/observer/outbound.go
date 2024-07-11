package observer

import (
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// GetTxID returns a unique id for Solana outbound
func (ob *Observer) GetTxID(_ uint64) string {
	//TODO implement me
	panic("implement me")
}

func (ob *Observer) IsOutboundProcessed(_ *types.CrossChainTx, _ zerolog.Logger) (bool, bool, error) {
	//TODO implement me
	panic("implement me")
}
