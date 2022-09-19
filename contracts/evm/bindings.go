//go:generate abigen --abi Connector.abi --pkg evm --type Connector --out Connector.go
//go:generate sh -c "solc ZRC4.sol --combined-json abi,bin | jq '.contracts.\"ZRC4.sol:ZRC4\"'  > ZRC4.json"
//go:generate sh -c "cat ZRC4.json | jq .abi | abigen --abi - --pkg evm --type ZRC4 --out ZRC4.go"
//go:generate sh -c "solc zeta_deposit_call.sol --combined-json abi,bin | jq '.contracts.\"zeta_deposit_call.sol:ZetaDepositAndCall\"'  > zeta_deposit_call.json"
//go:generate sh -c "cat zeta_deposit_call.json | jq .abi | abigen --abi - --pkg evm --type ZetaDepositAndCall --out ZetaDepositAndCall.go"

package evm

import (
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

var _ = Connector{}

var _ = ZRC4{}

type CompiledContract struct {
	ABI abi.ABI
	Bin evmtypes.HexString
}

var (
	//go:embed ZRC4.json
	ZRC4JSON []byte // nolint: golint

	//go:embed zeta_deposit_call.json
	ZetaDepositAndCallJSON []byte // nolint: golint

	ZRC4Contract               CompiledContract
	ZetaDepositAndCallContract CompiledContract

	// the module address of zetacore; no private exists.
	ZRC4AdminAddress ethcommon.Address
)

func init() {
	ZRC4AdminAddress = fungibletypes.ModuleAddressEVM

	err := json.Unmarshal(ZRC4JSON, &ZRC4Contract)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(ZetaDepositAndCallJSON, &ZetaDepositAndCallContract)
	if err != nil {
		panic(err)
	}

	if len(ZRC4Contract.Bin) == 0 {
		panic("load contract failed")
	}

	if len(ZRC4Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
