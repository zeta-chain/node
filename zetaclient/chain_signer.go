package zetaclient

import (
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

type ChainSigner interface {
	TryProcessOutTx(send *types.CrossChainTx, outTxMan *OutTxProcessorManager, evmClient ChainClient, zetaBridge *ZetaCoreBridge)
}
