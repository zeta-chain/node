package testutils

import (
	"fmt"

	"github.com/zeta-chain/zetacore/pkg/coin"
)

// FileNameEVMBlock returns unified archive file name for block
func FileNameEVMBlock(chainID int64, blockNumber uint64, trimmed bool) string {
	if !trimmed {
		return fmt.Sprintf("chain_%d_block_ethrpc_%d.json", chainID, blockNumber)
	}
	return fmt.Sprintf("chain_%d_block_ethrpc_trimmed_%d.json", chainID, blockNumber)
}

// FileNameCctxByIntx returns unified archive cctx file name by intx
func FileNameCctxByIntx(chainID int64, intxHash string, coinType coin.CoinType) string {
	return fmt.Sprintf("cctx_intx_%d_%s_%s.json", chainID, coinType, intxHash)
}

// FileNameCctxByNonce returns unified archive cctx file name by nonce
func FileNameCctxByNonce(chainID int64, nonce uint64) string {
	return fmt.Sprintf("chain_%d_cctx_%d.json", chainID, nonce)
}

// FileNameEVMIntx returns unified archive file name for inbound tx
func FileNameEVMIntx(chainID int64, intxHash string, coinType coin.CoinType, donation bool) string {
	if !donation {
		return fmt.Sprintf("chain_%d_intx_ethrpc_%s_%s.json", chainID, coinType, intxHash)
	}
	return fmt.Sprintf("chain_%d_intx_ethrpc_donation_%s_%s.json", chainID, coinType, intxHash)
}

// FileNameEVMIntxReceipt returns unified archive file name for inbound tx receipt
func FileNameEVMIntxReceipt(chainID int64, intxHash string, coinType coin.CoinType, donation bool) string {
	if !donation {
		return fmt.Sprintf("chain_%d_intx_receipt_%s_%s.json", chainID, coinType, intxHash)
	}
	return fmt.Sprintf("chain_%d_intx_receipt_donation_%s_%s.json", chainID, coinType, intxHash)
}

// FileNameBTCIntx returns unified archive file name for inbound tx
func FileNameBTCIntx(chainID int64, intxHash string, donation bool) string {
	if !donation {
		return fmt.Sprintf("chain_%d_intx_raw_result_%s.json", chainID, intxHash)
	}
	return fmt.Sprintf("chain_%d_intx_raw_result_donation_%s.json", chainID, intxHash)
}

// FileNameBTCOuttx returns unified archive file name for outbound tx
func FileNameBTCOuttx(chainID int64, nonce uint64) string {
	return fmt.Sprintf("chain_%d_outtx_raw_result_nonce_%d.json", chainID, nonce)
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

// FileNameEVMOuttx returns unified archive file name for outbound tx
func FileNameEVMOuttx(chainID int64, txHash string, coinType coin.CoinType) string {
	return fmt.Sprintf("chain_%d_outtx_%s_%s.json", chainID, coinType, txHash)
}

// FileNameEVMOuttxReceipt returns unified archive file name for outbound tx receipt
func FileNameEVMOuttxReceipt(chainID int64, txHash string, coinType coin.CoinType, eventName string) string {
	// empty eventName is for regular transfer receipt, no event
	if eventName == "" {
		return fmt.Sprintf("chain_%d_outtx_receipt_%s_%s.json", chainID, coinType, txHash)
	}
	return fmt.Sprintf("chain_%d_outtx_receipt_%s_%s_%s.json", chainID, coinType, eventName, txHash)
}
