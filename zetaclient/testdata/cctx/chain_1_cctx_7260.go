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
	InboundTxParams: &crosschaintypes.InboundTxParams{
		Sender:                          "0xd91b507F2A3e2D4A32d0C86Ac19FEAD2D461008D",
		SenderChainId:                   7000,
		TxOrigin:                        "0x8E62e3e6FbFF3E21F725395416A20EA4E2DeF015",
		CoinType:                        coin.CoinType_Gas,
		Asset:                           "",
		Amount:                          sdkmath.NewUint(42635427434588308),
		InboundTxObservedHash:           "0x2720e3a98f18c288f4197d412bfce57e58f00dc4f8b31e335ffc0bf7208dd3c5",
		InboundTxObservedExternalHeight: 2031411,
		InboundTxBallotIndex:            "0xbebecbf1d8c12016e38c09d095290df503fe29731722d939433fa47e3ed1f986",
		InboundTxFinalizedZetaHeight:    0,
		TxFinalizationStatus:            crosschaintypes.TxFinalizationStatus_NotFinalized,
	},
	OutboundTxParams: []*crosschaintypes.OutboundTxParams{
		{
			Receiver:                         "0x8E62e3e6FbFF3E21F725395416A20EA4E2DeF015",
			ReceiverChainId:                  1,
			CoinType:                         coin.CoinType_Gas,
			Amount:                           sdkmath.NewUint(42635427434588308),
			OutboundTxTssNonce:               7260,
			OutboundTxGasLimit:               21000,
			OutboundTxGasPrice:               "236882693686",
			OutboundTxHash:                   "0xd13b593eb62b5500a00e288cc2fb2c8af1339025c0e6bc6183b8bef2ebbed0d3",
			OutboundTxBallotIndex:            "0x689d894606642a2a7964fa906ebf4998c22a00708544fa88e9c56b86c955066b",
			OutboundTxObservedExternalHeight: 19363323,
			OutboundTxGasUsed:                21000,
			OutboundTxEffectiveGasPrice:      sdkmath.NewInt(236882693686),
			OutboundTxEffectiveGasLimit:      21000,
			TssPubkey:                        "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:             crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}
