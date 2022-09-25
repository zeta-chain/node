//go:generate sh -c "solc ZRC4.sol --combined-json abi,bin | jq '.contracts.\"ZRC4.sol:ZRC4\"'  > ZRC4.json"
//go:generate sh -c "cat ZRC4.json | jq .abi | abigen --abi - --pkg zevm --type ZRC4 --out ZRC4.go"
//go:generate sh -c "solc ZetaDepositAndCall.sol --combined-json abi,bin | jq '.contracts.\"ZetaDepositAndCall.sol:ZetaDepositAndCall\"'  > ZetaDepositAndCall.json"
//go:generate sh -c "cat ZetaDepositAndCall.json | jq .abi | abigen --abi - --pkg zevm --type ZetaDepositAndCall --out ZetaDepositAndCall.go"
//go:generate sh -c "solc GasPriceOracle.sol --combined-json abi,bin | jq '.contracts.\"GasPriceOracle.sol:GasPriceOracle\"'  > GasPriceOracle.json"
//go:generate sh -c "cat GasPriceOracle.json | jq .abi | abigen --abi - --pkg zevm --type GasPriceOracle --out GasPriceOracle.go"

package zevm

import (
	_ "embed"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// test the existence of the generated bindings
var _ = ZRC4{}
var _ = GasPriceOracle{}
var _ = ZetaDepositAndCall{}

type CompiledContract struct {
	ABI abi.ABI
	Bin evmtypes.HexString
}

var (
	//go:embed ZRC4.json
	ZRC4JSON []byte // nolint: golint

	//go:embed ZetaDepositAndCall.json
	ZetaDepositAndCallJSON []byte // nolint: golint
	//go:embed GasPriceOracle.json
	GasPriceOracleJSON []byte // nolint: golint

	ZRC4Contract               CompiledContract
	ZetaDepositAndCallContract CompiledContract
	GasPriceOracleContract     CompiledContract

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
	err = json.Unmarshal(GasPriceOracleJSON, &GasPriceOracleContract)
	if err != nil {
		panic(err)
	}

	if len(ZRC4Contract.Bin) == 0 {
		panic("load contract failed")
	}

	if len(ZetaDepositAndCallContract.Bin) == 0 {
		panic("load contract failed")
	}

	if len(GasPriceOracleContract.Bin) == 0 {
		panic("load contract failed")
	}
}
