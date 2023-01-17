//go:generate sh -c "cat ZetaConnectorEth.json | jq .abi | abigen --abi - --pkg evm --type Connector --out Connector.go"

//go:generate sh -c "cat ZetaEth.json | jq .abi > ZetaEth.abi"
//go:generate sh -c "cat ZetaEth.json | jq .bytecode | tr -d '\"' > ZetaEth.bin"
//go:generate sh -c "abigen --abi ZetaEth.abi --bin ZetaEth.bin --pkg evm --type ZetaEth --out ZetaEth.go"

//go:generate sh -c "cat ZetaNonEth.json | jq .abi > ZetaNonEth.abi"
//go:generate sh -c "cat ZetaNonEth.json | jq .bytecode | tr -d '\"' > ZetaNonEth.bin"
//go:generate sh -c "abigen --abi ZetaNonEth.abi --bin ZetaNonEth.bin --pkg evm --type ZetaNonEth --out ZetaNonEth.go"

//go:generate sh -c "cat ZetaConnectorEth.json | jq .abi > ZetaConnectorEth.abi"
//go:generate sh -c "cat ZetaConnectorEth.json | jq .bytecode | tr -d '\"' > ZetaConnectorEth.bin"
//go:generate sh -c "abigen --abi ZetaConnectorEth.abi --bin ZetaConnectorEth.bin --pkg evm --type ZetaConnectorEth --out ZetaConnectorEth.go"

//go:generate sh -c "cat ZetaConnectorNonEth.json | jq .abi > ZetaConnectorNonEth.abi"
//go:generate sh -c "cat ZetaConnectorNonEth.json | jq .bytecode | tr -d '\"' > ZetaConnectorNonEth.bin"
//go:generate sh -c "abigen --abi ZetaConnectorNonEth.abi --bin ZetaConnectorNonEth.bin --pkg evm --type ZetaConnectorNonEth --out ZetaConnectorNonEth.go"

package evm

import (
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

var _ = Connector{}

type CompiledContract struct {
	ABI      abi.ABI
	Bytecode evmtypes.HexString
}

var (
	//go:embed ZetaEth.json
	ZETAETHJSON []byte // nolint: golint

	ZetaEthContract CompiledContract
)

func init() {
	err := json.Unmarshal(ZETAETHJSON, &ZetaEthContract)
	if err != nil {
		panic(err)
	}
}
