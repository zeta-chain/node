package cctx

import (
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// CCtxByNonceMap maps the [chainID, nonce] to the cross chain transaction
var CCtxByNonceMap = map[int64]map[uint64]*crosschaintypes.CrossChainTx{
	// Ethereum mainnet
	1: {
		chain_1_cctx_6270.GetCurrentOutboundParam().TssNonce: chain_1_cctx_6270,
		chain_1_cctx_7260.GetCurrentOutboundParam().TssNonce: chain_1_cctx_7260,
		chain_1_cctx_8014.GetCurrentOutboundParam().TssNonce: chain_1_cctx_8014,
		chain_1_cctx_9718.GetCurrentOutboundParam().TssNonce: chain_1_cctx_9718,
	},
	// BSC mainnet
	56: {
		chain_56_cctx_68270.GetCurrentOutboundParam().TssNonce: chain_56_cctx_68270,
	},
	// local goerli testnet
	1337: {
		chain_1337_cctx_14.GetCurrentOutboundParam().TssNonce: chain_1337_cctx_14,
	},
	// Bitcoin mainnet
	8332: {
		chain_8332_cctx_148.GetCurrentOutboundParam().TssNonce: chain_8332_cctx_148,
	},
}

// CctxByIntxMap maps the [chainID, coinType, intxHash] to the cross chain transaction
var CctxByIntxMap = map[int64]map[coin.CoinType]map[string]*crosschaintypes.CrossChainTx{
	// Ethereum mainnet
	1: {
		coin.CoinType_Zeta: {
			chain_1_cctx_intx_Zeta_0xf393520.InboundParams.ObservedHash: chain_1_cctx_intx_Zeta_0xf393520,
		},
		coin.CoinType_ERC20: {
			chain_1_cctx_intx_ERC20_0x4ea69a0.InboundParams.ObservedHash: chain_1_cctx_intx_ERC20_0x4ea69a0,
		},
		coin.CoinType_Gas: {
			chain_1_cctx_intx_Gas_0xeaec67d.InboundParams.ObservedHash: chain_1_cctx_intx_Gas_0xeaec67d,
		},
	},
	// BSC mainnet
	56: {},
	// Bitcoin mainnet
	8332: {},
}
