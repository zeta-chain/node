package testutils

import (
	"fmt"

	"github.com/zeta-chain/node/pkg/coin"
)

// FileNameEVMBlock returns unified archive file name for block
func FileNameEVMBlock(chainID int64, blockNumber uint64, trimmed bool) string {
	if !trimmed {
		return fmt.Sprintf("chain_%d_block_ethrpc_%d.json", chainID, blockNumber)
	}
	return fmt.Sprintf("chain_%d_block_ethrpc_trimmed_%d.json", chainID, blockNumber)
}

// FileNameCctxByInbound returns unified archive cctx file name by inbound
func FileNameCctxByInbound(chainID int64, inboundHash string, coinType coin.CoinType) string {
	return fmt.Sprintf("cctx_inbound_%d_%s_%s.json", chainID, coinType, inboundHash)
}

// FileNameCctxByNonce returns unified archive cctx file name by nonce
func FileNameCctxByNonce(chainID int64, nonce uint64) string {
	return fmt.Sprintf("chain_%d_cctx_%d.json", chainID, nonce)
}

// FileNameEVMInbound returns unified archive file name for inbound tx
func FileNameEVMInbound(chainID int64, inboundHash string, coinType coin.CoinType, donation bool) string {
	if !donation {
		return fmt.Sprintf("chain_%d_inbound_ethrpc_%s_%s.json", chainID, coinType, inboundHash)
	}
	return fmt.Sprintf("chain_%d_inbound_ethrpc_donation_%s_%s.json", chainID, coinType, inboundHash)
}

// FileNameEVMInboundReceipt returns unified archive file name for inbound tx receipt
func FileNameEVMInboundReceipt(chainID int64, inboundHash string, coinType coin.CoinType, donation bool) string {
	if !donation {
		return fmt.Sprintf("chain_%d_inbound_receipt_%s_%s.json", chainID, coinType, inboundHash)
	}
	return fmt.Sprintf("chain_%d_inbound_receipt_donation_%s_%s.json", chainID, coinType, inboundHash)
}

// FileNameBTCInbound returns unified archive file name for inbound tx
func FileNameBTCInbound(chainID int64, inboundHash string, donation bool) string {
	if !donation {
		return fmt.Sprintf("chain_%d_inbound_raw_result_%s.json", chainID, inboundHash)
	}
	return fmt.Sprintf("chain_%d_inbound_raw_result_donation_%s.json", chainID, inboundHash)
}

// FileNameBTCOutbound returns unified archive file name for outbound tx
func FileNameBTCOutbound(chainID int64, nonce uint64) string {
	return fmt.Sprintf("chain_%d_outbound_raw_result_nonce_%d.json", chainID, nonce)
}

// FileNameBTCTxByType returns unified archive file name for tx by type
// txType: "P2TR", "P2WPKH", "P2WSH", "P2PKH", "P2SH
func FileNameBTCTxByType(chainID int64, txType string, txHash string) string {
	return fmt.Sprintf("chain_%d_tx_raw_result_%s_%s.json", chainID, txType, txHash)
}

// FileNameBTCMsgTx returns unified archive file name for btc MsgTx
func FileNameBTCMsgTx(chainID int64, txHash string) string {
	return fmt.Sprintf("chain_%d_msgtx_%s.json", chainID, txHash)
}

// FileNameEVMOutbound returns unified archive file name for outbound tx
func FileNameEVMOutbound(chainID int64, txHash string, coinType coin.CoinType) string {
	return fmt.Sprintf("chain_%d_outbound_%s_%s.json", chainID, coinType, txHash)
}

// FileNameEVMOutboundReceipt returns unified archive file name for outbound tx receipt
func FileNameEVMOutboundReceipt(chainID int64, txHash string, coinType coin.CoinType, eventName string) string {
	// empty eventName is for regular transfer receipt, no event
	if eventName == "" {
		return fmt.Sprintf("chain_%d_outbound_receipt_%s_%s.json", chainID, coinType, txHash)
	}
	return fmt.Sprintf("chain_%d_outbound_receipt_%s_%s_%s.json", chainID, coinType, eventName, txHash)
}

//=================================================================================================
// Solana chain

// FileNameSolanaInbound returns archive file name for inbound tx result
func FileNameSolanaInbound(chainID int64, inboundHash string, donation bool) string {
	if !donation {
		return fmt.Sprintf("chain_%d_inbound_tx_result_%s.json", chainID, inboundHash)
	}
	return fmt.Sprintf("chain_%d_inbound_tx_result_donation_%s.json", chainID, inboundHash)
}

// FileNameSolanaOutbound returns archive file name for outbound tx result
func FileNameSolanaOutbound(chainID int64, txHash string) string {
	return fmt.Sprintf("chain_%d_outbound_tx_result_%s.json", chainID, txHash)
}
