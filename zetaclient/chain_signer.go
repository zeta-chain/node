package zetaclient

import (
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

type ChainSigner interface {
	TryProcessOutTx(send *types.CrossChainTx, outTxMan *OutTxProcessorManager, outTxID string, evmClient ChainClient, zetaBridge *ZetaCoreBridge, height uint64)
}
