package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestVerifyInboundBody(t *testing.T) {
	sampleTo := sample.EthAddress()
	sampleEthTx, sampleEthTxBytes := sample.EthTx(t, chains.EthChain.ChainId, sampleTo, 42)

	// NOTE: errContains == "" means no error
	for _, tc := range []struct {
		desc        string
		msg         types.MsgAddInboundTracker
		txBytes     []byte
		chainParams observertypes.ChainParams
		tss         observertypes.QueryGetTssAddressResponse
		errContains string
	}{
		{
			desc: "can't verify btc tx tx body",
			msg: types.MsgAddInboundTracker{
				ChainId: chains.BtcMainnetChain.ChainId,
			},
			txBytes:     sample.Bytes(),
			errContains: "cannot verify inbound body for chain",
		},
		{
			desc: "txBytes can't be unmarshaled",
			msg: types.MsgAddInboundTracker{
				ChainId: chains.EthChain.ChainId,
			},
			txBytes:     []byte("invalid"),
			errContains: "failed to unmarshal transaction",
		},
		{
			desc: "txHash doesn't correspond",
			msg: types.MsgAddInboundTracker{
				ChainId: chains.EthChain.ChainId,
				TxHash:  sample.Hash().Hex(),
			},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid hash",
		},
		{
			desc: "chain id doesn't correspond",
			msg: types.MsgAddInboundTracker{
				ChainId: chains.SepoliaChain.ChainId,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid chain id",
		},
		{
			desc: "invalid coin type",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType(1000),
			},
			txBytes:     sampleEthTxBytes,
			errContains: "coin type not supported",
		},
		{
			desc: "coin types is zeta, but connector contract address is wrong",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Zeta,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{ConnectorContractAddress: sample.EthAddress().Hex()},
			errContains: "receiver is not connector contract for coin type",
		},
		{
			desc: "coin types is zeta, connector contract address is correct",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Zeta,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{ConnectorContractAddress: sampleTo.Hex()},
		},
		{
			desc: "coin types is erc20, but erc20 custody contract address is wrong",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_ERC20,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{Erc20CustodyContractAddress: sample.EthAddress().Hex()},
			errContains: "receiver is not erc20Custory contract for coin type",
		},
		{
			desc: "coin types is erc20, erc20 custody contract address is correct",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_ERC20,
			},
			txBytes:     sampleEthTxBytes,
			chainParams: observertypes.ChainParams{Erc20CustodyContractAddress: sampleTo.Hex()},
		},
		{
			desc: "coin types is gas, but tss address is not found",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Gas,
			},
			txBytes:     sampleEthTxBytes,
			tss:         observertypes.QueryGetTssAddressResponse{},
			errContains: "tss address not found",
		},
		{
			desc: "coin types is gas, but tss address is wrong",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Gas,
			},
			txBytes:     sampleEthTxBytes,
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sample.EthAddress().Hex()},
			errContains: "receiver is not tssAddress contract for coin type",
		},
		{
			desc: "coin types is gas, tss address is correct",
			msg: types.MsgAddInboundTracker{
				ChainId:  chains.EthChain.ChainId,
				TxHash:   sampleEthTx.Hash().Hex(),
				CoinType: coin.CoinType_Gas,
			},
			txBytes: sampleEthTxBytes,
			tss:     observertypes.QueryGetTssAddressResponse{Eth: sampleTo.Hex()},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := types.VerifyInboundBody(tc.msg, tc.txBytes, tc.chainParams, tc.tss)
			if tc.errContains == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.errContains)
			}
		})
	}
}

func TestVerifyOutboundBody(t *testing.T) {

	sampleTo := sample.EthAddress()
	sampleEthTx, sampleEthTxBytes, sampleFrom := sample.EthTxSigned(t, chains.EthChain.ChainId, sampleTo, 42)
	_, sampleEthTxBytesNonSigned := sample.EthTx(t, chains.EthChain.ChainId, sampleTo, 42)

	// NOTE: errContains == "" means no error
	for _, tc := range []struct {
		desc        string
		msg         types.MsgAddOutboundTracker
		txBytes     []byte
		tss         observertypes.QueryGetTssAddressResponse
		errContains string
	}{
		{
			desc: "invalid chain id",
			msg: types.MsgAddOutboundTracker{
				ChainId: int64(1000),
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     sample.Bytes(),
			errContains: "cannot verify outbound body for chain",
		},
		{
			desc: "txBytes can't be unmarshaled",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.EthChain.ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     []byte("invalid"),
			errContains: "failed to unmarshal transaction",
		},
		{
			desc: "can't recover sender address",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.EthChain.ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			txBytes:     sampleEthTxBytesNonSigned,
			errContains: "failed to recover sender",
		},
		{
			desc: "tss address not found",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.EthChain.ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{},
			txBytes:     sampleEthTxBytes,
			errContains: "tss address not found",
		},
		{
			desc: "tss address is wrong",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.EthChain.ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sample.EthAddress().Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "sender is not tss address",
		},
		{
			desc: "chain id doesn't correspond",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.SepoliaChain.ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid chain id",
		},
		{
			desc: "nonce doesn't correspond",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.EthChain.ChainId,
				Nonce:   100,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid nonce",
		},
		{
			desc: "tx hash doesn't correspond",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.EthChain.ChainId,
				Nonce:   42,
				TxHash:  sample.Hash().Hex(),
			},
			tss:         observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes:     sampleEthTxBytes,
			errContains: "invalid tx hash",
		},
		{
			desc: "valid outbound body",
			msg: types.MsgAddOutboundTracker{
				ChainId: chains.EthChain.ChainId,
				Nonce:   42,
				TxHash:  sampleEthTx.Hash().Hex(),
			},
			tss:     observertypes.QueryGetTssAddressResponse{Eth: sampleFrom.Hex()},
			txBytes: sampleEthTxBytes,
		},
		// TODO: Implement tests for verifyOutboundBodyBTC
		// https://github.com/zeta-chain/node/issues/1994
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := types.VerifyOutboundBody(tc.msg, tc.txBytes, tc.tss)
			if tc.errContains == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.errContains)
			}
		})
	}
}
