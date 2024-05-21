package cctx

import (
	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/56/68270
var chain_56_cctx_68270 = &crosschaintypes.CrossChainTx{
	Creator:        "",
	Index:          "0x541b570182950809f9b9077861a0fc7038af9a14ce8af4e151a83adfa308c7a9",
	ZetaFees:       sdkmath.ZeroUint(),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_PendingOutbound,
		StatusMessage:       "",
		LastUpdateTimestamp: 1709145183,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0xd91b507F2A3e2D4A32d0C86Ac19FEAD2D461008D",
		SenderChainId:          7000,
		TxOrigin:               "0xb0C04e07A301927672A8A7a874DB6930576C90B8",
		CoinType:               coin.CoinType_Gas,
		Asset:                  "",
		Amount:                 sdkmath.NewUint(657177295293237048),
		ObservedHash:           "0x093f4ca4c1884df0fd9dd59b75979342ded29d3c9b6861644287a2e1417b9a39",
		ObservedExternalHeight: 1960153,
		BallotIndex:            "0x541b570182950809f9b9077861a0fc7038af9a14ce8af4e151a83adfa308c7a9",
		FinalizedZetaHeight:    0,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "0xb0C04e07A301927672A8A7a874DB6930576C90B8",
			ReceiverChainId:        56,
			CoinType:               coin.CoinType_Gas,
			Amount:                 sdkmath.NewUint(657177295293237048),
			TssNonce:               68270,
			GasLimit:               21000,
			GasPrice:               "6000000000",
			Hash:                   "0xeb2b183ece6638688b9df9223180b13a67208cd744bbdadeab8de0482d7f4e3c",
			BallotIndex:            "0xa4600c952934f797e162d637d70859a611757407908d96bc53e45a81c80b006b",
			ObservedExternalHeight: 36537856,
			GasUsed:                21000,
			EffectiveGasPrice:      sdkmath.NewInt(6000000000),
			EffectiveGasLimit:      21000,
			TssPubkey:              "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}
