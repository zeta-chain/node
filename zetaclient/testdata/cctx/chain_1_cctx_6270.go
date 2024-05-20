package cctx

import (
	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// https://zetachain-mainnet-archive.allthatnode.com:1317/zeta-chain/crosschain/cctx/1/6270
var chain_1_cctx_6270 = &crosschaintypes.CrossChainTx{
	Creator:        "",
	Index:          "0xe930f363591b348a07e0a6d309b4301b84f702e3e81e0d0902340c7f7da4b5af",
	ZetaFees:       sdkmath.ZeroUint(),
	RelayedMessage: "",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_OutboundMined,
		StatusMessage:       "",
		LastUpdateTimestamp: 1708464433,
		IsAbortRefunded:     false,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0xd91b507F2A3e2D4A32d0C86Ac19FEAD2D461008D",
		SenderChainId:          7000,
		TxOrigin:               "0x18D0E2c38b4188D8Ae07008C3BeeB1c80748b41c",
		CoinType:               coin.CoinType_Gas,
		Asset:                  "",
		Amount:                 sdkmath.NewUint(9831832641427386),
		ObservedHash:           "0x8bd0df31e512c472e3162a41281b740b518216cc8eb787c2eb59c81e0cffbe89",
		ObservedExternalHeight: 1846989,
		BallotIndex:            "0xe930f363591b348a07e0a6d309b4301b84f702e3e81e0d0902340c7f7da4b5af",
		FinalizedZetaHeight:    0,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_NotFinalized,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "0x18D0E2c38b4188D8Ae07008C3BeeB1c80748b41c",
			ReceiverChainId:        1,
			CoinType:               coin.CoinType_Gas,
			Amount:                 sdkmath.NewUint(9831832641427386),
			TssNonce:               6270,
			GasLimit:               21000,
			GasPrice:               "69197693654",
			Hash:                   "0x20104d41e042db754cf7908c5441914e581b498eedbca40979c9853f4b7f8460",
			BallotIndex:            "0x346a1d00a4d26a2065fe1dc7d5af59a49ad6a8af25853ae2ec976c07349f48c1",
			ObservedExternalHeight: 19271550,
			GasUsed:                21000,
			EffectiveGasPrice:      sdkmath.NewInt(69197693654),
			EffectiveGasLimit:      21000,
			TssPubkey:              "zetapub1addwnpepqtadxdyt037h86z60nl98t6zk56mw5zpnm79tsmvspln3hgt5phdc79kvfc",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}
