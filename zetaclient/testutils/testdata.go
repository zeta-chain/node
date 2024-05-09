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
	// nameCctx := path.Join("../", TestDataPathCctx, FileNameCctxByIntx(chainID, intxHash, coinType))

	// cctx := &crosschaintypes.CrossChainTx{}
	// LoadObjectFromJSONFile(t, &cctx, nameCctx)
	// return cctx

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
	// nameCctx := path.Join("../", TestDataPathCctx, FileNameCctxByNonce(chainID, nonce))

	// cctx := &crosschaintypes.CrossChainTx{}
	// LoadObjectFromJSONFile(t, &cctx, nameCctx)

	// get cctx
	cctx, found := testcctx.CCtxByNonceMap[chainID][nonce]
	require.True(t, found)

	// clone cctx for each individual test
	cloned := cloneCctx(t, cctx)
	return cloned
}

// LoadEVMBlock loads archived evm block from file
func LoadEVMBlock(t *testing.T, chainID int64, blockNumber uint64, trimmed bool) *ethrpc.Block {
	name := path.Join("../", TestDataPathEVM, FileNameEVMBlock(chainID, blockNumber, trimmed))
	block := &ethrpc.Block{}
	LoadObjectFromJSONFile(t, block, name)
	return block
}

// LoadBTCInboundRawResult loads archived Bitcoin intx raw result from file
func LoadBTCInboundRawResult(t *testing.T, chainID int64, txHash string, donation bool) *btcjson.TxRawResult {
	name := path.Join("../", TestDataPathBTC, FileNameBTCInbound(chainID, txHash, donation))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, name)
	return rawResult
}

// LoadBTCTxRawResultNCctx loads archived Bitcoin outtx raw result and corresponding cctx
func LoadBTCTxRawResultNCctx(t *testing.T, chainID int64, nonce uint64) (*btcjson.TxRawResult, *crosschaintypes.CrossChainTx) {
	//nameTx := FileNameBTCOutbound(chainID, nonce)
	nameTx := path.Join("../", TestDataPathBTC, FileNameBTCOutbound(chainID, nonce))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, nameTx)

	cctx := LoadCctxByNonce(t, chainID, nonce)
	return rawResult, cctx
}

// LoadEVMInbound loads archived intx from file
func LoadEVMInbound(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethrpc.Transaction {
	nameTx := path.Join("../", TestDataPathEVM, FileNameEVMInbound(chainID, intxHash, coinType, false))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMInboundReceipt loads archived intx receipt from file
func LoadEVMInboundReceipt(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join("../", TestDataPathEVM, FileNameEVMInboundReceipt(chainID, intxHash, coinType, false))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMInboundNReceipt loads archived intx and receipt from file
func LoadEVMInboundNReceipt(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived intx and receipt
	tx := LoadEVMInbound(t, chainID, intxHash, coinType)
	receipt := LoadEVMInboundReceipt(t, chainID, intxHash, coinType)

	return tx, receipt
}

// LoadEVMInboundDonation loads archived donation intx from file
func LoadEVMInboundDonation(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethrpc.Transaction {
	nameTx := path.Join("../", TestDataPathEVM, FileNameEVMInbound(chainID, intxHash, coinType, true))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMInboundReceiptDonation loads archived donation intx receipt from file
func LoadEVMInboundReceiptDonation(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join("../", TestDataPathEVM, FileNameEVMInboundReceipt(chainID, intxHash, coinType, true))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMInboundNReceiptDonation loads archived donation intx and receipt from file
func LoadEVMInboundNReceiptDonation(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived donation intx and receipt
	tx := LoadEVMInboundDonation(t, chainID, intxHash, coinType)
	receipt := LoadEVMInboundReceiptDonation(t, chainID, intxHash, coinType)

	return tx, receipt
}

// LoadEVMInboundNReceiptNCctx loads archived intx, receipt and corresponding cctx from file
func LoadEVMInboundNReceiptNCctx(
	t *testing.T,
	chainID int64,
	intxHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt, *crosschaintypes.CrossChainTx) {
	// load archived intx, receipt and cctx
	tx := LoadEVMInbound(t, chainID, intxHash, coinType)
	receipt := LoadEVMInboundReceipt(t, chainID, intxHash, coinType)
	cctx := LoadCctxByIntx(t, chainID, coinType, intxHash)

	return tx, receipt, cctx
}

// LoadEVMOutbound loads archived evm outtx from file
func LoadEVMOutbound(
	t *testing.T,
	chainID int64,
	txHash string,
	coinType coin.CoinType) *ethtypes.Transaction {
	nameTx := path.Join("../", TestDataPathEVM, FileNameEVMOutbound(chainID, txHash, coinType))

	tx := &ethtypes.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMOutboundReceipt loads archived evm outtx receipt from file
func LoadEVMOutboundReceipt(
	t *testing.T,
	chainID int64,
	txHash string,
	coinType coin.CoinType,
	eventName string) *ethtypes.Receipt {
	nameReceipt := path.Join("../", TestDataPathEVM, FileNameEVMOutboundReceipt(chainID, txHash, coinType, eventName))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMOutboundNReceipt loads archived evm outtx and receipt from file
func LoadEVMOutboundNReceipt(
	t *testing.T,
	chainID int64,
	txHash string,
	coinType coin.CoinType) (*ethtypes.Transaction, *ethtypes.Receipt) {
	// load archived evm outtx and receipt
	tx := LoadEVMOutbound(t, chainID, txHash, coinType)
	receipt := LoadEVMOutboundReceipt(t, chainID, txHash, coinType, "")

	return tx, receipt
}

// LoadEVMCctxNOutboundNReceipt loads archived cctx, outtx and receipt from file
func LoadEVMCctxNOutboundNReceipt(
	t *testing.T,
	chainID int64,
	nonce uint64,
	eventName string) (*crosschaintypes.CrossChainTx, *ethtypes.Transaction, *ethtypes.Receipt) {
	cctx := LoadCctxByNonce(t, chainID, nonce)
	coinType := cctx.GetCurrentOutboundParam().CoinType
	txHash := cctx.GetCurrentOutboundParam().Hash
	outtx := LoadEVMOutbound(t, chainID, txHash, coinType)
	receipt := LoadEVMOutboundReceipt(t, chainID, txHash, coinType, eventName)
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
