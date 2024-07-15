package observer

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// GetTxID returns a unique id for Solana outbound
func (ob *Observer) GetTxID(_ uint64) string {
	//TODO implement me
	panic("implement me")
}

// IsOutboundProcessed checks outbound status and returns (isIncluded, isConfirmed, error)
func (ob *Observer) IsOutboundProcessed(
	_ context.Context,
	_ *types.CrossChainTx,
	_ zerolog.Logger,
) (bool, bool, error) {
	//TODO implement me
	panic("implement me")
}
