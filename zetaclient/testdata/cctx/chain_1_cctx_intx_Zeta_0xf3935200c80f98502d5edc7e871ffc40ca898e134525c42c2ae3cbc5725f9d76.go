package cctx

import (
	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/0x477544c4b8c8be544b23328b21286125c89cd6bb5d1d6d388d91eea8ea1a6f1f
var chain_1_cctx_intx_Zeta_0xf393520 = &crosschaintypes.CrossChainTx{
	Creator:        "zeta1p0uwsq4naus5r4l7l744upy0k8ezzj84mn40nf",
	Index:          "0x477544c4b8c8be544b23328b21286125c89cd6bb5d1d6d388d91eea8ea1a6f1f",
	ZetaFees:       sdkmath.NewUintFromString("0"),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "Remote omnichain contract call completed",
		LastUpdateTimestamp: 1708490549,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0x2f993766e8e1Ef9288B1F33F6aa244911A0A77a7",
		SenderChainId:          1,
		TxOrigin:               "0x2f993766e8e1Ef9288B1F33F6aa244911A0A77a7",
		CoinType:               coin.CoinType_Zeta,
		Asset:                  "",
		Amount:                 sdkmath.NewUintFromString("20000000000000000000"),
		ObservedHash:           "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76",
		ObservedExternalHeight: 19273702,
		BallotIndex:            "0x9039e8b0493d58ee4db876ffd0785970019d7eb515e2f71b00fe670bb6a8e6ce",
		FinalizedZetaHeight:    1851403,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "0x2f993766e8e1ef9288b1f33f6aa244911a0a77a7",
			ReceiverChainId:        7000,
			CoinType:               coin.CoinType_Zeta,
			Amount:                 sdkmath.ZeroUint(),
			TssNonce:               0,
			GasLimit:               100000,
			GasPrice:               "",
			Hash:                   "0x947434364da7c74d7e896a389aa8cb3122faf24bbcba64b141cb5acd7838209c",
			BallotIndex:            "",
			ObservedExternalHeight: 1851403,
			GasUsed:                0,
			EffectiveGasPrice:      sdkmath.ZeroInt(),
			EffectiveGasLimit:      0,
			TssPubkey:              "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
		},
	},
}
