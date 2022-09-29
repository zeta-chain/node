//go:generate sh -c "solc ZRC4.sol --combined-json abi,bin | jq '.contracts.\"ZRC4.sol:ZRC4\"'  > ZRC4.json"
//go:generate sh -c "cat ZRC4.json | jq .abi | abigen --abi - --pkg zevm --type ZRC4 --out ZRC4.go"
//go:generate sh -c "cat UniswapV2Factory.json | jq .abi | abigen --abi - --pkg zevm --type UniswapV2Factory --out UniswapV2Factory.go"
//go:generate sh -c "cat WZETA.json | jq .abi | abigen --abi - --pkg zevm --type WZETA --out WZETA.go"
//go:generate sh -c "solc SystemContract.sol --combined-json abi,bin | jq '.contracts.\"SystemContract.sol:SystemContract\"'  > SystemContract.json"
//go:generate sh -c "cat SystemContract.json | jq .abi | abigen --abi - --pkg zevm --type SystemContract --out SystemContract.go"

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
var _ = UniswapV2Factory{}
var _ = SystemContract{}

type CompiledContract struct {
	ABI abi.ABI
	Bin evmtypes.HexString
}

var (
	//go:embed ZRC4.json
	ZRC4JSON []byte // nolint: golint
	//go:embed UniswapV2Factory.json
	UniswapV2FactoryJSON []byte // nolint: golint
	//go:embed WZETA.json
	WZETAJSON []byte // nolint: golint
	//go:embed SystemContract.json
	SystemContractJSON []byte // nolint: golint

	ZRC4Contract             CompiledContract
	UniswapV2FactoryContract CompiledContract
	WZETAContract            CompiledContract
	SystemContractContract   CompiledContract

	// the module address of zetacore; no private exists.
	ZRC4AdminAddress ethcommon.Address
)

func init() {
	ZRC4AdminAddress = fungibletypes.ModuleAddressEVM

	err := json.Unmarshal(ZRC4JSON, &ZRC4Contract)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(UniswapV2FactoryJSON, &UniswapV2FactoryContract)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(WZETAJSON, &WZETAContract)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(SystemContractJSON, &SystemContractContract)
	if err != nil {
		panic(err)
	}

	if len(ZRC4Contract.Bin) == 0 {
		panic("load contract failed")
	}

	if len(UniswapV2FactoryContract.Bin) == 0 {
		panic("load contract failed")
	}

	if len(WZETAContract.Bin) == 0 {
		panic("load contract failed")
	}

	if len(SystemContractContract.Bin) == 0 {
		panic("load contract failed")
	}
}
