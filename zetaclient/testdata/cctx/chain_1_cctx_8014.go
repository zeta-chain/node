package cctx

import (
	sdkmath "cosmossdk.io/math"

	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/1/8014
var chain_1_cctx_8014 = &crosschaintypes.CrossChainTx{
	Creator:        "",
	Index:          "0x5a100fdb426da35ad4c95520d7a4f1fd2f38c88067c9e80ba209d3a655c6e06e",
	ZetaFees:       sdkmath.ZeroUint(),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "",
		LastUpdateTimestamp: 1710834402,
		IsAbortRefunded:     false,
	},
	InboundTxParams: &crosschaintypes.InboundTxParams{
		Sender:                          "0x7c8dDa80bbBE1254a7aACf3219EBe1481c6E01d7",
		SenderChainId:                   7000,
		TxOrigin:                        "0x8d8D67A8B71c141492825CAE5112Ccd8581073f2",
		CoinType:                        coin.CoinType_ERC20,
		Asset:                           "0xdac17f958d2ee523a2206206994597c13d831ec7",
		Amount:                          sdkmath.NewUint(23726342442),
		InboundTxObservedHash:           "0x114ed9d327b6afc068c3fa891b82f7c7f2d42ae25a571f7dc004c05e77af592a",
		InboundTxObservedExternalHeight: 2241077,
		InboundTxBallotIndex:            "0x5a100fdb426da35ad4c95520d7a4f1fd2f38c88067c9e80ba209d3a655c6e06e",
		InboundTxFinalizedZetaHeight:    0,
		TxFinalizationStatus:            crosschaintypes.TxFinalizationStatus_NotFinalized,
	},
	OutboundTxParams: []*crosschaintypes.OutboundTxParams{
		{
			Receiver:                         "0x8d8D67A8B71c141492825CAE5112Ccd8581073f2",
			ReceiverChainId:                  1,
			CoinType:                         coin.CoinType_ERC20,
			Amount:                           sdkmath.NewUint(23726342442),
			OutboundTxTssNonce:               8014,
			OutboundTxGasLimit:               100000,
			OutboundTxGasPrice:               "58619665744",
			OutboundTxHash:                   "0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146",
			OutboundTxBallotIndex:            "0x4213f2c335758301b8bbb09d9891949ed6ffeea5dd95e5d9eaa8d410baaa0884",
			OutboundTxObservedExternalHeight: 19467367,
			OutboundTxGasUsed:                60625,
			OutboundTxEffectiveGasPrice:      sdkmath.NewInt(58619665744),
			OutboundTxEffectiveGasLimit:      100000,
			TssPubkey:                        "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:             crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}
