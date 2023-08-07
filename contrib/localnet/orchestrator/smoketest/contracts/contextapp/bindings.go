// solc use version 0.8.7
//go:generate sh -c " solc --evm-version london --allow-paths ..,  --combined-json abi,bin --base-path .. ContextApp.sol     | jq '.contracts.\"ContextApp.sol:ContextApp\"'  > ContextApp.json"
//go:generate sh -c "cat ContextApp.json | jq .abi > ContextApp.abi"
//go:generate sh -c "cat ContextApp.json | jq .bin  | tr -d '\"'  > ContextApp.bin"

//go:generate sh -c "abigen --abi ContextApp.abi --bin ContextApp.bin --pkg contextapp --type ContextApp --out ContextApp.go"

package contextapp

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
	//go:embed ContextApp.json
	ContextAppJSON []byte // nolint: golint

	ContextAppContract CompiledContract
)

func init() {
	err := json.Unmarshal(ContextAppJSON, &ContextAppContract)
	if err != nil {
		panic(err)
	}

	if len(ContextAppContract.Bin) == 0 {
		panic("load contract failed")
	}
}
