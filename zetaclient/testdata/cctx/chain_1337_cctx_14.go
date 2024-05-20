package cctx

import (
	sdkmath "cosmossdk.io/math"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// This cctx was generated in local e2e tests, see original json text attached at end of file
var chain_1337_cctx_14 = &crosschaintypes.CrossChainTx{
	Creator:        "zeta1plfrp7ejn0s9tmwufuxvsyn8nlf6a7u9ndgk9m",
	Index:          "0x85d06ac908823d125a919164f0596e3496224b206ebe8125ffe7b4ab766f85df",
	ZetaFees:       sdkmath.NewUintFromString("4000000000009027082"),
	RelayedMessage: "bgGCGUux3roBhJr9PgNaC3DOfLBp5ILuZjUZx2z1abQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ==",
	CctxStatus: &crosschaintypes.Status{
		Status:              crosschaintypes.CctxStatus_Reverted,
		StatusMessage:       "Outbound failed, start revert : Outbound succeeded, revert executed",
		LastUpdateTimestamp: 1712705995,
	},
	InboundParams: &crosschaintypes.InboundParams{
		Sender:                 "0xBFF76e77D56B3C1202107f059425D56f0AEF87Ed",
		SenderChainId:          1337,
		TxOrigin:               "0x5cC2fBb200A929B372e3016F1925DcF988E081fd",
		Amount:                 sdkmath.NewUintFromString("10000000000000000000"),
		ObservedHash:           "0xa5589bf24eca8f108ca35048adc9d5582a303d416c01319391159269ae7e4e6f",
		ObservedExternalHeight: 177,
		BallotIndex:            "0x85d06ac908823d125a919164f0596e3496224b206ebe8125ffe7b4ab766f85df",
		FinalizedZetaHeight:    150,
		TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
	},
	OutboundParams: []*crosschaintypes.OutboundParams{
		{
			Receiver:               "0xbff76e77d56b3c1202107f059425d56f0aef87ed",
			ReceiverChainId:        1337,
			Amount:                 sdkmath.NewUintFromString("7999999999995486459"),
			TssNonce:               13,
			GasLimit:               250000,
			GasPrice:               "18",
			Hash:                   "0x19f99459da6cb08f917f9b0ee2dac94a7be328371dff788ad46e64a24e8c06c9",
			ObservedExternalHeight: 187,
			GasUsed:                67852,
			EffectiveGasPrice:      sdkmath.NewInt(18),
			EffectiveGasLimit:      250000,
			TssPubkey:              "zetapub1addwnpepqggky6z958k7hhxs6k5quuvs27uv5vtmlv330ppt2362p8ejct88w4g64jv",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
		},
		{
			Receiver:               "0xBFF76e77D56B3C1202107f059425D56f0AEF87Ed",
			ReceiverChainId:        1337,
			Amount:                 sdkmath.NewUintFromString("5999999999990972918"),
			TssNonce:               14,
			GasLimit:               250000,
			GasPrice:               "18",
			Hash:                   "0x1487e6a31dd430306667250b72bf15b8390b73108b69f3de5c1b2efe456036a7",
			BallotIndex:            "0xc36c689fdaf09a9b80a614420cd4fea4fec15044790df60080cdefca0090a9dc",
			ObservedExternalHeight: 201,
			GasUsed:                76128,
			EffectiveGasPrice:      sdkmath.NewInt(18),
			EffectiveGasLimit:      250000,
			TssPubkey:              "zetapub1addwnpepqggky6z958k7hhxs6k5quuvs27uv5vtmlv330ppt2362p8ejct88w4g64jv",
			TxFinalizationStatus:   crosschaintypes.TxFinalizationStatus_Executed,
		},
	},
}

// Here is the original cctx json data used to create above chain_1337_cctx_14
/*
{
  "creator": "zeta1plfrp7ejn0s9tmwufuxvsyn8nlf6a7u9ndgk9m",
  "index": "0x85d06ac908823d125a919164f0596e3496224b206ebe8125ffe7b4ab766f85df",
  "zeta_fees": "4000000000009027082",
  "relayed_message": "bgGCGUux3roBhJr9PgNaC3DOfLBp5ILuZjUZx2z1abQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQ==",
  "cctx_status": {
    "status": 5,
    "status_message": "Outbound failed, start revert : Outbound succeeded, revert executed",
    "lastUpdate_timestamp": 1712705995
  },
  "inbound_tx_params": {
    "sender": "0xBFF76e77D56B3C1202107f059425D56f0AEF87Ed",
    "sender_chain_id": 1337,
    "tx_origin": "0x5cC2fBb200A929B372e3016F1925DcF988E081fd",
    "amount": "10000000000000000000",
    "inbound_tx_observed_hash": "0xa5589bf24eca8f108ca35048adc9d5582a303d416c01319391159269ae7e4e6f",
    "inbound_tx_observed_external_height": 177,
    "inbound_tx_ballot_index": "0x85d06ac908823d125a919164f0596e3496224b206ebe8125ffe7b4ab766f85df",
    "inbound_tx_finalized_zeta_height": 150,
    "tx_finalization_status": 2
  },
  "outbound_tx_params": [
    {
      "receiver": "0xbff76e77d56b3c1202107f059425d56f0aef87ed",
      "receiver_chainId": 1337,
      "amount": "7999999999995486459",
      "outbound_tx_tss_nonce": 13,
      "outbound_tx_gas_limit": 250000,
      "outbound_tx_gas_price": "18",
      "outbound_tx_hash": "0x19f99459da6cb08f917f9b0ee2dac94a7be328371dff788ad46e64a24e8c06c9",
      "outbound_tx_observed_external_height": 187,
      "outbound_tx_gas_used": 67852,
      "outbound_tx_effective_gas_price": "18",
      "outbound_tx_effective_gas_limit": 250000,
      "tss_pubkey": "zetapub1addwnpepqggky6z958k7hhxs6k5quuvs27uv5vtmlv330ppt2362p8ejct88w4g64jv",
      "tx_finalization_status": 2
    },
    {
      "receiver": "0xBFF76e77D56B3C1202107f059425D56f0AEF87Ed",
      "receiver_chainId": 1337,
      "amount": "5999999999990972918",
      "outbound_tx_tss_nonce": 14,
      "outbound_tx_gas_limit": 250000,
      "outbound_tx_gas_price": "18",
      "outbound_tx_hash": "0x1487e6a31dd430306667250b72bf15b8390b73108b69f3de5c1b2efe456036a7",
      "outbound_tx_ballot_index": "0xc36c689fdaf09a9b80a614420cd4fea4fec15044790df60080cdefca0090a9dc",
      "outbound_tx_observed_external_height": 201,
      "outbound_tx_gas_used": 76128,
      "outbound_tx_effective_gas_price": "18",
      "outbound_tx_effective_gas_limit": 250000,
      "tss_pubkey": "zetapub1addwnpepqggky6z958k7hhxs6k5quuvs27uv5vtmlv330ppt2362p8ejct88w4g64jv",
      "tx_finalization_status": 2
    }
  ]
}
*/
