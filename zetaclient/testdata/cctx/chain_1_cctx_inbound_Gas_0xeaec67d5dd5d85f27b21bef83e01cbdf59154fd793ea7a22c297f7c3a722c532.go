package cctx

import (
	sdkmath "cosmossdk.io/math"

	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/0x0210925c7dceeff563e6e240d6d650a5f0e8fca6f5b76044a6cf106d21f27098
var chain_1_cctx_inbound_Gas_0xeaec67d = &crosschaintypes.CrossChainTx{
	Creator:        "zeta1hjct6q7npsspsg3dgvzk3sdf89spmlpf7rqmnw",
	Index:          "0x0210925c7dceeff563e6e240d6d650a5f0e8fca6f5b76044a6cf106d21f27098",
	ZetaFees:       sdkmath.NewUint(0),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "Remote omnichain contract call completed",
		LastUpdateTimestamp: 1709177431,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0xF829fa7069680b8C37A8086b37d4a24697E5003b",
		SenderChainId:          1,
		TxOrigin:               "0xF829fa7069680b8C37A8086b37d4a24697E5003b",
		CoinType:               coin.CoinType_Gas,
		Asset:                  "",
		Amount:                 sdkmath.NewUintFromString("4000000000000000"),
		ObservedHash:           "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532",
		ObservedExternalHeight: 19330473,
		BallotIndex:            "0x79f6a3da92d085b2f3c682a2eb1606ef89e53a7db4fdbbb397b3e0f54884cfb0",
		FinalizedZetaHeight:    1965579,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "0xF829fa7069680b8C37A8086b37d4a24697E5003b",
			ReceiverChainId:        7000,
			CoinType:               coin.CoinType_Gas,
			Amount:                 sdkmath.NewUint(0),
			TssNonce:               0,
			GasLimit:               90000,
			GasPrice:               "",
			Hash:                   "0x3b8c1dab5aa21ff90ddb569f2f962ff2d4aa8d914c9177900102e745955e6f35",
			BallotIndex:            "",
			ObservedExternalHeight: 1965579,
			GasUsed:                0,
			EffectiveGasPrice:      sdkmath.NewInt(0),
			EffectiveGasLimit:      0,
			TssPubkey:              "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
		},
	},
}
