package cctx

import (
	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/1/7260
var chain_1_cctx_7260 = &crosschaintypes.CrossChainTx{
	Creator:        "",
	Index:          "0xbebecbf1d8c12016e38c09d095290df503fe29731722d939433fa47e3ed1f986",
	ZetaFees:       sdkmath.ZeroUint(),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "",
		LastUpdateTimestamp: 1709574082,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0xd91b507F2A3e2D4A32d0C86Ac19FEAD2D461008D",
		SenderChainId:          7000,
		TxOrigin:               "0x8E62e3e6FbFF3E21F725395416A20EA4E2DeF015",
		CoinType:               coin.CoinType_Gas,
		Asset:                  "",
		Amount:                 sdkmath.NewUint(42635427434588308),
		ObservedHash:           "0x2720e3a98f18c288f4197d412bfce57e58f00dc4f8b31e335ffc0bf7208dd3c5",
		ObservedExternalHeight: 2031411,
		BallotIndex:            "0xbebecbf1d8c12016e38c09d095290df503fe29731722d939433fa47e3ed1f986",
		FinalizedZetaHeight:    0,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "0x8E62e3e6FbFF3E21F725395416A20EA4E2DeF015",
			ReceiverChainId:        1,
			CoinType:               coin.CoinType_Gas,
			Amount:                 sdkmath.NewUint(42635427434588308),
			TssNonce:               7260,
			GasLimit:               21000,
			GasPrice:               "236882693686",
			Hash:                   "0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3",
			BallotIndex:            "0xca4ca249ce29a6305ca88eec8957a6b74e74df3a3bdfe7cd14d7e951b7c820c8",
			ObservedExternalHeight: 19363323,
			GasUsed:                21000,
			EffectiveGasPrice:      sdkmath.NewInt(236882693686),
			EffectiveGasLimit:      21000,
			TssPubkey:              "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}
