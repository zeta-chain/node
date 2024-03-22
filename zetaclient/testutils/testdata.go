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
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const (
	TestDataPathEVM          = "testdata/evm"
	TestDataPathBTC          = "testdata/btc"
	TestDataPathCctx         = "testdata/cctx"
	RestrictedEVMAddressTest = "0x8a81Ba8eCF2c418CAe624be726F505332DF119C6"
	RestrictedBtcAddressTest = "bcrt1qzp4gt6fc7zkds09kfzaf9ln9c5rvrzxmy6qmpp"
)

// SaveObjectToJSONFile saves an object to a file in JSON format
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

// LoadObjectFromJSONFile loads an object from a file in JSON format
func LoadObjectFromJSONFile(obj interface{}, filename string) {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// read the struct from the file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&obj)
	if err != nil {
		panic(err)
	}
}

func ComplianceConfigTest() config.ComplianceConfig {
	return config.ComplianceConfig{
		RestrictedAddresses: []string{RestrictedEVMAddressTest, RestrictedBtcAddressTest},
	}
}

// SaveTrimedEVMBlockTrimTxInput trims tx input data from a block and saves it to a file
func SaveEVMBlockTrimTxInput(block *ethrpc.Block, filename string) error {
	for i := range block.Transactions {
		block.Transactions[i].Input = "0x"
	}
	return SaveObjectToJSONFile(block, filename)
}

// SaveTrimedBTCBlockTrimTx trims tx data from a block and saves it to a file
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

// LoadEVMBlock loads archived evm block from file
func LoadEVMBlock(chainID int64, blockNumber uint64, trimmed bool) *ethrpc.Block {
	name := path.Join("../", TestDataPathEVM, FileNameEVMBlock(chainID, blockNumber, trimmed))
	block := &ethrpc.Block{}
	LoadObjectFromJSONFile(block, name)
	return block
}

// LoadBTCTxRawResultNCctx loads archived Bitcoin outtx raw result and corresponding cctx
func LoadBTCTxRawResultNCctx(chainID int64, nonce uint64) (*btcjson.TxRawResult, *crosschaintypes.CrossChainTx) {
	//nameTx := FileNameBTCOuttx(chainID, nonce)
	nameTx := path.Join("../", TestDataPathBTC, FileNameBTCOuttx(chainID, nonce))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(rawResult, nameTx)

	nameCctx := path.Join("../", TestDataPathCctx, FileNameCctxByNonce(chainID, nonce))
	cctx := &crosschaintypes.CrossChainTx{}
	LoadObjectFromJSONFile(cctx, nameCctx)
	return rawResult, cctx
}

// LoadEVMIntx loads archived intx from file
func LoadEVMIntx(
	_ *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) *ethrpc.Transaction {
	nameTx := path.Join("../", TestDataPathEVM, FileNameEVMIntx(chainID, intxHash, coinType, false))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(&tx, nameTx)
	return tx
}

// LoadEVMIntxReceipt loads archived intx receipt from file
func LoadEVMIntxReceipt(
	_ *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join("../", TestDataPathEVM, FileNameEVMIntxReceipt(chainID, intxHash, coinType, false))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(&receipt, nameReceipt)
	return receipt
}

// LoadEVMIntxCctx loads archived intx cctx from file
func LoadEVMIntxCctx(
	_ *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) *crosschaintypes.CrossChainTx {
	nameCctx := path.Join("../", TestDataPathCctx, FileNameEVMIntxCctx(chainID, intxHash, coinType))

	cctx := &crosschaintypes.CrossChainTx{}
	LoadObjectFromJSONFile(&cctx, nameCctx)
	return cctx
}

// LoadCctxByNonce loads archived cctx by nonce from file
func LoadCctxByNonce(
	_ *testing.T,
	chainID int64,
	nonce uint64) *crosschaintypes.CrossChainTx {
	nameCctx := path.Join("../", TestDataPathCctx, FileNameCctxByNonce(chainID, nonce))

	cctx := &crosschaintypes.CrossChainTx{}
	LoadObjectFromJSONFile(&cctx, nameCctx)
	return cctx
}

// LoadEVMIntxNReceipt loads archived intx and receipt from file
func LoadEVMIntxNReceipt(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived intx and receipt
	tx := LoadEVMIntx(t, chainID, intxHash, coinType)
	receipt := LoadEVMIntxReceipt(t, chainID, intxHash, coinType)

	return tx, receipt
}

// LoadEVMIntxDonation loads archived donation intx from file
func LoadEVMIntxDonation(
	_ *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) *ethrpc.Transaction {
	nameTx := path.Join("../", TestDataPathEVM, FileNameEVMIntx(chainID, intxHash, coinType, true))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(&tx, nameTx)
	return tx
}

// LoadEVMIntxReceiptDonation loads archived donation intx receipt from file
func LoadEVMIntxReceiptDonation(
	_ *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join("../", TestDataPathEVM, FileNameEVMIntxReceipt(chainID, intxHash, coinType, true))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(&receipt, nameReceipt)
	return receipt
}

// LoadEVMIntxNReceiptDonation loads archived donation intx and receipt from file
func LoadEVMIntxNReceiptDonation(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived donation intx and receipt
	tx := LoadEVMIntxDonation(t, chainID, intxHash, coinType)
	receipt := LoadEVMIntxReceiptDonation(t, chainID, intxHash, coinType)

	return tx, receipt
}

// LoadTxNReceiptNCctx loads archived intx, receipt and corresponding cctx from file
func LoadEVMIntxNReceiptNCctx(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt, *crosschaintypes.CrossChainTx) {
	// load archived intx, receipt and cctx
	tx := LoadEVMIntx(t, chainID, intxHash, coinType)
	receipt := LoadEVMIntxReceipt(t, chainID, intxHash, coinType)
	cctx := LoadEVMIntxCctx(t, chainID, intxHash, coinType)

	return tx, receipt, cctx
}

// LoadEVMOuttx loads archived evm outtx from file
func LoadEVMOuttx(
	_ *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) *ethtypes.Transaction {
	nameTx := path.Join("../", TestDataPathEVM, FileNameEVMOuttx(chainID, intxHash, coinType))

	tx := &ethtypes.Transaction{}
	LoadObjectFromJSONFile(&tx, nameTx)
	return tx
}

// LoadEVMOuttxReceipt loads archived evm outtx receipt from file
func LoadEVMOuttxReceipt(
	_ *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join("../", TestDataPathEVM, FileNameEVMOuttxReceipt(chainID, intxHash, coinType))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(&receipt, nameReceipt)
	return receipt
}

// LoadEVMOuttxNReceipt loads archived evm outtx and receipt from file
func LoadEVMOuttxNReceipt(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType common.CoinType) (*ethtypes.Transaction, *ethtypes.Receipt) {
	// load archived evm outtx and receipt
	tx := LoadEVMOuttx(t, chainID, intxHash, coinType)
	receipt := LoadEVMOuttxReceipt(t, chainID, intxHash, coinType)

	return tx, receipt
}
