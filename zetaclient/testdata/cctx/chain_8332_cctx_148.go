package cctx

import (
	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/8332/148
var chain_8332_cctx_148 = &crosschaintypes.CrossChainTx{
	Creator:        "",
	Index:          "0xb3f5f3cf2ed2e0c3fa64c8297c9e50fbc07351fb2d26d8eae4cfbbd45e47a524",
	ZetaFees:       sdkmath.ZeroUint(),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "",
		LastUpdateTimestamp: 1708608895,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0x13A0c5930C028511Dc02665E7285134B6d11A5f4",
		SenderChainId:          7000,
		TxOrigin:               "0xe99174F08e1186134830f8511De06bd010978533",
		CoinType:               coin.CoinType_Gas,
		Asset:                  "",
		Amount:                 sdkmath.NewUint(12000),
		ObservedHash:           "0x06455013319acb1b027461134853c77b003d8eab162b1f37673da5ad8a50b74f",
		ObservedExternalHeight: 1870408,
		BallotIndex:            "0xb3f5f3cf2ed2e0c3fa64c8297c9e50fbc07351fb2d26d8eae4cfbbd45e47a524",
		FinalizedZetaHeight:    0,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "bc1qpsdlklfcmlcfgm77c43x65ddtrt7n0z57hsyjp",
			ReceiverChainId:        8332,
			CoinType:               coin.CoinType_Gas,
			Amount:                 sdkmath.NewUint(12000),
			TssNonce:               148,
			GasLimit:               254,
			GasPrice:               "46",
			Hash:                   "030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0",
			ObservedExternalHeight: 150,
			GasUsed:                0,
			EffectiveGasPrice:      sdkmath.NewInt(0),
			EffectiveGasLimit:      0,
			TssPubkey:              "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}
