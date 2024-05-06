package cctx

import (
	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/1/9718
var chain_1_cctx_9718 = &crosschaintypes.CrossChainTx{
	Creator:        "",
	Index:          "0xbf7a214cf9868e1c618123ab4df0081da87bade74eeb5aef37843e35f25e67b7",
	ZetaFees:       sdkmath.NewUintFromString("19525506001302763608"),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "",
		LastUpdateTimestamp: 1712336965,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0xF0a3F93Ed1B126142E61423F9546bf1323Ff82DF",
		SenderChainId:          7000,
		TxOrigin:               "0x87257C910a19a3fe64AfFAbFe8cF9AAF2ab148BF",
		CoinType:               coin.CoinType_Zeta,
		Asset:                  "",
		Amount:                 sdkmath.NewUintFromString("20000000000000000000"),
		ObservedHash:           "0xb136652cd58fb6a537b0a1677965983059a2004d98919cdacd52551f877cc44f",
		ObservedExternalHeight: 2492552,
		BallotIndex:            "0xbf7a214cf9868e1c618123ab4df0081da87bade74eeb5aef37843e35f25e67b7",
		FinalizedZetaHeight:    0,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "0x30735c88fa430f11499b0edcfcc25246fb9182e3",
			ReceiverChainId:        1,
			CoinType:               coin.CoinType_Zeta,
			Amount:                 sdkmath.NewUint(474493998697236392),
			TssNonce:               9718,
			GasLimit:               90000,
			GasPrice:               "112217884384",
			Hash:                   "0x81342051b8a85072d3e3771c1a57c7bdb5318e8caf37f5a687b7a91e50a7257f",
			BallotIndex:            "0xff07eaa34ca02a08bca1558e5f6220cbfc734061f083622b24923e032f0c480f",
			ObservedExternalHeight: 19590894,
			GasUsed:                64651,
			EffectiveGasPrice:      sdkmath.NewInt(112217884384),
			EffectiveGasLimit:      100000,
			TssPubkey:              "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}
