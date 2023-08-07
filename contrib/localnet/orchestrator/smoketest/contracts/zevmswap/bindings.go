// solc/abigen use version 0.8.7
//go:generate sh -c " solc --evm-version london --allow-paths ..,  --combined-json abi,bin --base-path .. ZEVMSwapApp.sol     | jq '.contracts.\"ZEVMSwapApp.sol:ZEVMSwapApp\"'  > ZEVMSwapApp.json"
//go:generate sh -c "cat ZEVMSwapApp.json | jq .abi > ZEVMSwapApp.abi"
//go:generate sh -c "cat ZEVMSwapApp.json | jq .bin  | tr -d '\"'  > ZEVMSwapApp.bin"

//go:generate sh -c "abigen --abi ZEVMSwapApp.abi --bin ZEVMSwapApp.bin --pkg zevmswap --type ZEVMSwapApp --out ZEVMSwapApp.go"

package zevmswap

import (
	_ "embed"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type CompiledContract struct {
	ABI abi.ABI
	Bin evmtypes.HexString
}

var (
	//go:embed ZEVMSwapApp.json
	ZEVMSwapAppJSON []byte // nolint: golint

	ZEVMSwapAppContract CompiledContract
)

func init() {
	err := json.Unmarshal(ZEVMSwapAppJSON, &ZEVMSwapAppContract)
	if err != nil {
		panic(err)
	}

	if len(ZEVMSwapAppContract.Bin) == 0 {
		panic("load contract failed")
	}
}
