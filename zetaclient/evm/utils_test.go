package evm

import (
	"fmt"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestCheckEvmTxLog(t *testing.T) {
	// test data
	connectorAddr := ethcommon.HexToAddress("0x00005e3125aba53c5652f9f0ce1a4cf91d8b15ea")
	txHash := "0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626"
	topics := []ethcommon.Hash{
		// https://goerli.etherscan.io/tx/0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626#eventlog
		ethcommon.HexToHash("0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4"),
		ethcommon.HexToHash("0x00000000000000000000000023856df5d563bd893fc7df864102d8bbfe7fc487"),
		ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000061"),
	}

	tests := []struct {
		name string
		vLog *ethtypes.Log
		fail bool
	}{
		{
			name: "chain reorganization",
			vLog: &ethtypes.Log{
				Removed: true,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash(txHash),
				Topics:  topics,
			},
			fail: true,
		},
		{
			name: "emitter address mismatch",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: ethcommon.HexToAddress("0x184ba627DB853244c9f17f3Cb4378cB8B39bf147"),
				TxHash:  ethcommon.HexToHash(txHash),
				Topics:  topics,
			},
			fail: true,
		},
		{
			name: "tx hash mismatch",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash("0x781c018d604af4dad0fe5e3cea4ad9fb949a996d8cd0cd04a92cadd7f08c05f2"),
				Topics:  topics,
			},
			fail: true,
		},
		{
			name: "topics mismatch",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash(txHash),
				Topics: []ethcommon.Hash{
					// https://goerli.etherscan.io/tx/0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626#eventlog
					ethcommon.HexToHash("0x7ec1c94701e09b1652f3e1d307e60c4b9ebf99aff8c2079fd1d8c585e031c4e4"),
					ethcommon.HexToHash("0x00000000000000000000000023856df5d563bd893fc7df864102d8bbfe7fc487"),
				},
			},
			fail: true,
		},
		{
			name: "should pass",
			vLog: &ethtypes.Log{
				Removed: false,
				Address: connectorAddr,
				TxHash:  ethcommon.HexToHash(txHash),
				Topics:  topics,
			},
			fail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("check test: %s\n", tt.name)
			err := CheckEvmTxLog(
				tt.vLog,
				connectorAddr,
				"0xb252c9e77feafdeeae25cc1f037a16c4b50fa03c494754b99a7339d816c79626",
				TopicsZetaSent,
			)
			if tt.fail {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckEvmTransaction(t *testing.T) {
	// use archived intx
	intxHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"

	t.Run("should pass for valid transaction", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		err := CheckEvmTransaction(tx)
		require.NoError(t, err)
	})
	t.Run("should fail for nil transaction", func(t *testing.T) {
		err := CheckEvmTransaction(nil)
		require.ErrorContains(t, err, "transaction is nil")
	})
	t.Run("should fail for empty hash", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.Hash = ""
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "hash is empty")
	})
	t.Run("should fail for negative nonce", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.Nonce = -1
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "nonce -1 is negative")
	})
	t.Run("should fail for invalid from address", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.From = "0x"
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "from 0x is not a valid hex address")
	})
	t.Run("should fail for invalid to address", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.To = "0xinvalid"
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "to 0xinvalid is not a valid hex address")
	})
	t.Run("should fail for negative value", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.Value = *big.NewInt(-1)
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "value -1 is negative")
	})
	t.Run("should fail for negative gas", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.Gas = -1
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "gas -1 is negative")
	})
	t.Run("should fail for negative gas price", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.GasPrice = *big.NewInt(-1)
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "gas price -1 is negative")
	})
	t.Run("should remove '0x' prefix from input data", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		err := CheckEvmTransaction(tx)
		require.NoError(t, err)
		require.Equal(t, "", tx.Input)
	})
	t.Run("nil block number should pass", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.BlockNumber = nil
		err := CheckEvmTransaction(tx)
		require.NoError(t, err)
	})
	t.Run("should fail for negative block number", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		negBlockNumber := -1
		tx.BlockNumber = &negBlockNumber
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "block number -1 is not positive")
	})
	t.Run("should fail for empty block hash", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.BlockHash = ""
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "block hash is empty")
	})
	t.Run("nil transaction index should fail", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.TransactionIndex = nil
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "index is nil")
	})
	t.Run("should fail for negative transaction index", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		negTransactionIndex := -1
		tx.TransactionIndex = &negTransactionIndex
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "index -1 is negative")
	})
	t.Run("should fail for invalid input data", func(t *testing.T) {
		tx := testutils.LoadEVMIntx(t, 1, intxHash, common.CoinType_Gas)
		tx.Input = "03befinvalid"
		err := CheckEvmTransaction(tx)
		require.ErrorContains(t, err, "input data is not hex encoded")
	})
}
