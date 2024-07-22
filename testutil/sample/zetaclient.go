package sample

import (
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/zetaclient/types"
)

// InboundEvent returns a sample InboundEvent.
func InboundEvent(chainID int64, sender string, receiver string, amount uint64, memo []byte) *types.InboundEvent {
	r := newRandFromSeed(chainID)

	return &types.InboundEvent{
		SenderChainID: chainID,
		Sender:        sender,
		Receiver:      receiver,
		TxOrigin:      sender,
		Amount:        amount,
		Memo:          memo,
		BlockNumber:   r.Uint64(),
		TxHash:        StringRandom(r, 32),
		Index:         0,
		CoinType:      coin.CoinType(r.Intn(100)),
		Asset:         StringRandom(r, 32),
	}
}
