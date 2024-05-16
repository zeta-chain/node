package testutils

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	testcctx "github.com/zeta-chain/zetacore/zetaclient/testdata/cctx"
)

const (
	TestDataPathEVM          = "testdata/evm"
	TestDataPathBTC          = "testdata/btc"
	TestDataPathCctx         = "testdata/cctx"
	RestrictedEVMAddressTest = "0x8a81Ba8eCF2c418CAe624be726F505332DF119C6"
	RestrictedBtcAddressTest = "bcrt1qzp4gt6fc7zkds09kfzaf9ln9c5rvrzxmy6qmpp"
)

// cloneCctx returns a deep copy of the cctx
func cloneCctx(t *testing.T, cctx *crosschaintypes.CrossChainTx) *crosschaintypes.CrossChainTx {
	data, err := cctx.Marshal()
	require.NoError(t, err)
	cloned := &crosschaintypes.CrossChainTx{}
	err = cloned.Unmarshal(data)
	require.NoError(t, err)
	return cloned
}

// LoadObjectFromJSONFile loads an object from a file in JSON format
func LoadObjectFromJSONFile(t *testing.T, obj interface{}, filename string) {
	file, err := os.Open(filepath.Clean(filename))
	require.NoError(t, err)
	defer file.Close()

	// read the struct from the file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&obj)
	require.NoError(t, err)
}

func ComplianceConfigTest() config.ComplianceConfig {
	return config.ComplianceConfig{
		RestrictedAddresses: []string{RestrictedEVMAddressTest, RestrictedBtcAddressTest},
	}
}

// LoadCctxByIntx loads archived cctx by intx
func LoadCctxByIntx(
	t *testing.T,
	chainID int64,
	coinType coin.CoinType,
	intxHash string,
) *crosschaintypes.CrossChainTx {
	// get cctx
	cctx, found := testcctx.CctxByIntxMap[chainID][coinType][intxHash]
	require.True(t, found)

	// clone cctx for each individual test
	cloned := cloneCctx(t, cctx)
	return cloned
}

// LoadCctxByNonce loads archived cctx by nonce
func LoadCctxByNonce(
	t *testing.T,
	chainID int64,
	nonce uint64,
) *crosschaintypes.CrossChainTx {
	// get cctx
	cctx, found := testcctx.CCtxByNonceMap[chainID][nonce]
	require.True(t, found)

	// clone cctx for each individual test
	cloned := cloneCctx(t, cctx)
	return cloned
}

// LoadEVMBlock loads archived evm block from file
func LoadEVMBlock(t *testing.T, dir string, chainID int64, blockNumber uint64, trimmed bool) *ethrpc.Block {
	name := path.Join(dir, TestDataPathEVM, FileNameEVMBlock(chainID, blockNumber, trimmed))
	block := &ethrpc.Block{}
	LoadObjectFromJSONFile(t, block, name)
	return block
}

// LoadBTCTxRawResult loads archived Bitcoin tx raw result from file
func LoadBTCTxRawResult(t *testing.T, dir string, chainID int64, txType string, txHash string) *btcjson.TxRawResult {
	name := path.Join(dir, TestDataPathBTC, FileNameBTCTxByType(chainID, txType, txHash))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, name)
	return rawResult
}

// LoadBTCIntxRawResult loads archived Bitcoin intx raw result from file
func LoadBTCIntxRawResult(t *testing.T, dir string, chainID int64, txHash string, donation bool) *btcjson.TxRawResult {
	name := path.Join(dir, TestDataPathBTC, FileNameBTCIntx(chainID, txHash, donation))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, name)
	return rawResult
}

// LoadBTCTxRawResultNCctx loads archived Bitcoin outtx raw result and corresponding cctx
func LoadBTCTxRawResultNCctx(t *testing.T, dir string, chainID int64, nonce uint64) (*btcjson.TxRawResult, *crosschaintypes.CrossChainTx) {
	nameTx := path.Join(dir, TestDataPathBTC, FileNameBTCOuttx(chainID, nonce))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, nameTx)

	cctx := LoadCctxByNonce(t, chainID, nonce)
	return rawResult, cctx
}

// LoadEVMIntx loads archived intx from file
func LoadEVMIntx(
	t *testing.T,
	dir string,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethrpc.Transaction {
	nameTx := path.Join(dir, TestDataPathEVM, FileNameEVMIntx(chainID, intxHash, coinType, false))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMIntxReceipt loads archived intx receipt from file
func LoadEVMIntxReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join(dir, TestDataPathEVM, FileNameEVMIntxReceipt(chainID, intxHash, coinType, false))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMIntxNReceipt loads archived intx and receipt from file
func LoadEVMIntxNReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived intx and receipt
	tx := LoadEVMIntx(t, dir, chainID, intxHash, coinType)
	receipt := LoadEVMIntxReceipt(t, dir, chainID, intxHash, coinType)

	return tx, receipt
}

// LoadEVMIntxDonation loads archived donation intx from file
func LoadEVMIntxDonation(
	t *testing.T,
	dir string,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethrpc.Transaction {
	nameTx := path.Join(dir, TestDataPathEVM, FileNameEVMIntx(chainID, intxHash, coinType, true))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMIntxReceiptDonation loads archived donation intx receipt from file
func LoadEVMIntxReceiptDonation(
	t *testing.T,
	dir string,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join(dir, TestDataPathEVM, FileNameEVMIntxReceipt(chainID, intxHash, coinType, true))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMIntxNReceiptDonation loads archived donation intx and receipt from file
func LoadEVMIntxNReceiptDonation(
	t *testing.T,
	dir string,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived donation intx and receipt
	tx := LoadEVMIntxDonation(t, dir, chainID, intxHash, coinType)
	receipt := LoadEVMIntxReceiptDonation(t, dir, chainID, intxHash, coinType)

	return tx, receipt
}

// LoadEVMIntxNReceiptNCctx loads archived intx, receipt and corresponding cctx from file
func LoadEVMIntxNReceiptNCctx(
	t *testing.T,
	dir string,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt, *crosschaintypes.CrossChainTx) {
	// load archived intx, receipt and cctx
	tx := LoadEVMIntx(t, dir, chainID, intxHash, coinType)
	receipt := LoadEVMIntxReceipt(t, dir, chainID, intxHash, coinType)
	cctx := LoadCctxByIntx(t, chainID, coinType, intxHash)

	return tx, receipt, cctx
}

// LoadEVMOuttx loads archived evm outtx from file
func LoadEVMOuttx(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	coinType coin.CoinType) *ethtypes.Transaction {
	nameTx := path.Join(dir, TestDataPathEVM, FileNameEVMOuttx(chainID, txHash, coinType))

	tx := &ethtypes.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMOuttxReceipt loads archived evm outtx receipt from file
func LoadEVMOuttxReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	coinType coin.CoinType,
	eventName string) *ethtypes.Receipt {
	nameReceipt := path.Join(dir, TestDataPathEVM, FileNameEVMOuttxReceipt(chainID, txHash, coinType, eventName))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMOuttxNReceipt loads archived evm outtx and receipt from file
func LoadEVMOuttxNReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	coinType coin.CoinType) (*ethtypes.Transaction, *ethtypes.Receipt) {
	// load archived evm outtx and receipt
	tx := LoadEVMOuttx(t, dir, chainID, txHash, coinType)
	receipt := LoadEVMOuttxReceipt(t, dir, chainID, txHash, coinType, "")

	return tx, receipt
}

// LoadEVMCctxNOuttxNReceipt loads archived cctx, outtx and receipt from file
func LoadEVMCctxNOuttxNReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	nonce uint64,
	eventName string) (*crosschaintypes.CrossChainTx, *ethtypes.Transaction, *ethtypes.Receipt) {
	cctx := LoadCctxByNonce(t, chainID, nonce)
	coinType := cctx.GetCurrentOutTxParam().CoinType
	txHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	outtx := LoadEVMOuttx(t, dir, chainID, txHash, coinType)
	receipt := LoadEVMOuttxReceipt(t, dir, chainID, txHash, coinType, eventName)
	return cctx, outtx, receipt
}

// SaveObjectToJSONFile saves an object to a file in JSON format
// NOTE: this function is not used in the tests but used when creating test data
func SaveObjectToJSONFile(obj interface{}, filename string) error {
	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// write the struct to the file
	encoder := json.NewEncoder(file)
	return encoder.Encode(obj)
}

// SaveEVMBlockTrimTxInput trims tx input data from a block and saves it to a file
// NOTE: this function is not used in the tests but used when creating test data
func SaveEVMBlockTrimTxInput(block *ethrpc.Block, filename string) error {
	for i := range block.Transactions {
		block.Transactions[i].Input = "0x"
	}
	return SaveObjectToJSONFile(block, filename)
}

// SaveBTCBlockTrimTx trims tx data from a block and saves it to a file
// NOTE: this function is not used in the tests but used when creating test data
func SaveBTCBlockTrimTx(blockVb *btcjson.GetBlockVerboseTxResult, filename string) error {
	for i := range blockVb.Tx {
		// reserve one coinbase tx and one non-coinbase tx
		if i >= 2 {
			blockVb.Tx[i].Hex = ""
			blockVb.Tx[i].Vin = nil
			blockVb.Tx[i].Vout = nil
		}
	}
	return SaveObjectToJSONFile(blockVb, filename)
}
