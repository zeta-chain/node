package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
)

const (
	TestDataPathEVM  = "testdata/evm"
	TestDataPathBTC  = "testdata/btc"
	TestDataPathCctx = "testdata/cctx"
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
func LoadObjectFromJSONFile(obj interface{}, filename string) error {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// read the struct from the file
	decoder := json.NewDecoder(file)
	return decoder.Decode(&obj)
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

func DummyCoreBridge() *zetabridge.ZetaCoreBridge {
	bridge, _ := zetabridge.NewZetaCoreBridge(
		&keys.Keys{OperatorAddress: types.AccAddress{}},
		"127.0.0.1",
		"",
		"zetachain_7000-1",
		false,
		nil)
	return bridge
}
