//go:generate sh -c "solc  ZEVMSwapApp.sol --combined-json abi,bin --base-path ../ --include-path interfaces/ | jq '.contracts.\"zevmswap/ZEVMSwapApp.sol:ZEVMSwapApp\"'  > ZEVMSwapApp.json"
//go:generate sh -c "cat ZEVMSwapApp.json | jq .abi | abigen --abi - --pkg zevmswap --type ZEVMSwapApp --out ZEVMSwapApp.go"

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
