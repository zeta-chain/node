package testutils

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/onrik/ethrpc"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	testcctx "github.com/zeta-chain/node/zetaclient/testdata/cctx"
	testtypes "github.com/zeta-chain/node/zetaclient/testutils/types"
)

const (
	TestDataPathEVM    = "testdata/evm"
	TestDataPathBTC    = "testdata/btc"
	TestDataPathSolana = "testdata/solana"
	TestDataPathCctx   = "testdata/cctx"
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

// LoadJSONRawMessageFromFile loads a raw JSON message from a file in JSON format
func LoadJSONRawMessageFromFile(t *testing.T, filename string) json.RawMessage {
	file, err := os.Open(filepath.Clean(filename))
	require.NoError(t, err)
	defer file.Close()

	// read the JSON raw message from the file
	decoder := json.NewDecoder(file)
	var raw json.RawMessage
	err = decoder.Decode(&raw)
	require.NoError(t, err)
	return raw
}

// LoadCctxByInbound loads archived cctx by inbound
func LoadCctxByInbound(
	t *testing.T,
	chainID int64,
	coinType coin.CoinType,
	inboundHash string,
) *crosschaintypes.CrossChainTx {
	// get cctx
	cctx, found := testcctx.CctxByInboundMap[chainID][coinType][inboundHash]
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

	// load archived block
	jsonMessage := LoadJSONRawMessageFromFile(t, name)
	blockProxy := new(testtypes.ProxyBlockWithTransactions)
	err := json.Unmarshal(jsonMessage, blockProxy)
	require.NoError(t, err)

	return blockProxy.ToBlock()
}

// LoadBTCTxRawResult loads archived Bitcoin tx raw result from file
func LoadBTCTxRawResult(t *testing.T, dir string, chainID int64, txType string, txHash string) *btcjson.TxRawResult {
	name := path.Join(dir, TestDataPathBTC, FileNameBTCTxByType(chainID, txType, txHash))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, name)
	return rawResult
}

// LoadBTCInboundRawResult loads archived Bitcoin inbound raw result from file
func LoadBTCInboundRawResult(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	donation bool,
) *btcjson.TxRawResult {
	name := path.Join(dir, TestDataPathBTC, FileNameBTCInbound(chainID, txHash, donation))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, name)
	return rawResult
}

// LoadBTCTxRawResultNCctx loads archived Bitcoin outbound raw result and corresponding cctx
func LoadBTCTxRawResultNCctx(
	t *testing.T,
	dir string,
	chainID int64,
	nonce uint64,
) (*btcjson.TxRawResult, *crosschaintypes.CrossChainTx) {
	nameTx := path.Join(dir, TestDataPathBTC, FileNameBTCOutbound(chainID, nonce))
	rawResult := &btcjson.TxRawResult{}
	LoadObjectFromJSONFile(t, rawResult, nameTx)

	cctx := LoadCctxByNonce(t, chainID, nonce)
	return rawResult, cctx
}

// LoadEVMInbound loads archived inbound from file
func LoadEVMInbound(
	t *testing.T,
	dir string,
	chainID int64,
	inboundHash string,
	coinType coin.CoinType) *ethrpc.Transaction {
	nameTx := path.Join(dir, TestDataPathEVM, FileNameEVMInbound(chainID, inboundHash, coinType, false))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMInboundReceipt loads archived inbound receipt from file
func LoadEVMInboundReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	inboundHash string,
	coinType coin.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join(dir, TestDataPathEVM, FileNameEVMInboundReceipt(chainID, inboundHash, coinType, false))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMInboundNReceipt loads archived inbound and receipt from file
func LoadEVMInboundNReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	inboundHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived inbound and receipt
	tx := LoadEVMInbound(t, dir, chainID, inboundHash, coinType)
	receipt := LoadEVMInboundReceipt(t, dir, chainID, inboundHash, coinType)

	return tx, receipt
}

// LoadEVMInboundDonation loads archived donation inbound from file
func LoadEVMInboundDonation(
	t *testing.T,
	dir string,
	chainID int64,
	inboundHash string,
	coinType coin.CoinType) *ethrpc.Transaction {
	nameTx := path.Join(dir, TestDataPathEVM, FileNameEVMInbound(chainID, inboundHash, coinType, true))

	tx := &ethrpc.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMInboundReceiptDonation loads archived donation inbound receipt from file
func LoadEVMInboundReceiptDonation(
	t *testing.T,
	dir string,
	chainID int64,
	inboundHash string,
	coinType coin.CoinType) *ethtypes.Receipt {
	nameReceipt := path.Join(dir, TestDataPathEVM, FileNameEVMInboundReceipt(chainID, inboundHash, coinType, true))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMInboundNReceiptDonation loads archived donation inbound and receipt from file
func LoadEVMInboundNReceiptDonation(
	t *testing.T,
	dir string,
	chainID int64,
	inboundHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt) {
	// load archived donation inbound and receipt
	tx := LoadEVMInboundDonation(t, dir, chainID, inboundHash, coinType)
	receipt := LoadEVMInboundReceiptDonation(t, dir, chainID, inboundHash, coinType)

	return tx, receipt
}

// LoadEVMInboundNReceiptNCctx loads archived inbound, receipt and corresponding cctx from file
func LoadEVMInboundNReceiptNCctx(
	t *testing.T,
	dir string,
	chainID int64,
	inboundHash string,
	coinType coin.CoinType) (*ethrpc.Transaction, *ethtypes.Receipt, *crosschaintypes.CrossChainTx) {
	// load archived inbound, receipt and cctx
	tx := LoadEVMInbound(t, dir, chainID, inboundHash, coinType)
	receipt := LoadEVMInboundReceipt(t, dir, chainID, inboundHash, coinType)
	cctx := LoadCctxByInbound(t, chainID, coinType, inboundHash)

	return tx, receipt, cctx
}

// LoadEVMOutbound loads archived evm outbound from file
func LoadEVMOutbound(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	coinType coin.CoinType) *ethtypes.Transaction {
	nameTx := path.Join(dir, TestDataPathEVM, FileNameEVMOutbound(chainID, txHash, coinType))

	tx := &ethtypes.Transaction{}
	LoadObjectFromJSONFile(t, &tx, nameTx)
	return tx
}

// LoadEVMOutboundReceipt loads archived evm outbound receipt from file
func LoadEVMOutboundReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	coinType coin.CoinType,
	eventName string) *ethtypes.Receipt {
	nameReceipt := path.Join(dir, TestDataPathEVM, FileNameEVMOutboundReceipt(chainID, txHash, coinType, eventName))

	receipt := &ethtypes.Receipt{}
	LoadObjectFromJSONFile(t, &receipt, nameReceipt)
	return receipt
}

// LoadEVMOutboundNReceipt loads archived evm outbound and receipt from file
func LoadEVMOutboundNReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	coinType coin.CoinType) (*ethtypes.Transaction, *ethtypes.Receipt) {
	// load archived evm outbound and receipt
	tx := LoadEVMOutbound(t, dir, chainID, txHash, coinType)
	receipt := LoadEVMOutboundReceipt(t, dir, chainID, txHash, coinType, "")

	return tx, receipt
}

// LoadEVMCctxNOutboundNReceipt loads archived cctx, outbound and receipt from file
func LoadEVMCctxNOutboundNReceipt(
	t *testing.T,
	dir string,
	chainID int64,
	nonce uint64,
	eventName string) (*crosschaintypes.CrossChainTx, *ethtypes.Transaction, *ethtypes.Receipt) {
	cctx := LoadCctxByNonce(t, chainID, nonce)
	coinType := cctx.GetCurrentOutboundParam().CoinType
	txHash := cctx.GetCurrentOutboundParam().Hash
	outbound := LoadEVMOutbound(t, dir, chainID, txHash, coinType)
	receipt := LoadEVMOutboundReceipt(t, dir, chainID, txHash, coinType, eventName)

	return cctx, outbound, receipt
}

//==============================================================================
// Solana chain

// LoadSolanaInboundTxResult loads archived Solana inbound tx result from file
func LoadSolanaInboundTxResult(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
	donation bool,
) *rpc.GetTransactionResult {
	name := path.Join(dir, TestDataPathSolana, FileNameSolanaInbound(chainID, txHash, donation))
	txResult := &rpc.GetTransactionResult{}
	LoadObjectFromJSONFile(t, txResult, name)
	return txResult
}

// LoadSolanaOutboundTxResult loads archived Solana outbound tx result from file
func LoadSolanaOutboundTxResult(
	t *testing.T,
	dir string,
	chainID int64,
	txHash string,
) *rpc.GetTransactionResult {
	name := path.Join(dir, TestDataPathSolana, FileNameSolanaOutbound(chainID, txHash))
	txResult := &rpc.GetTransactionResult{}
	LoadObjectFromJSONFile(t, txResult, name)
	return txResult
}

//==============================================================================
// other helpers methods

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
