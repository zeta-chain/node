//go:generate sh -c "solc ZRC20.sol --combined-json abi,bin | jq '.contracts.\"ZRC20.sol:ZRC20\"'  > ZRC20.json"
//go:generate sh -c "cat ZRC20.json | jq .abi | abigen --abi - --pkg zevm --type ZRC20 --out ZRC20.go"
//go:generate sh -c "cat UniswapV2Factory.json | jq .abi | abigen --abi - --pkg zevm --type UniswapV2Factory --out UniswapV2Factory.go"
//go:generate sh -c "cat UniswapV2Router02.json | jq .abi | abigen --abi - --pkg zevm --type UniswapV2Router02 --out UniswapV2Router02.go"
//go:generate sh -c "cat WZETA.json | jq .abi | abigen --abi - --pkg zevm --type WZETA --out WZETA.go"
//go:generate sh -c "solc SystemContract.sol --combined-json abi,bin | jq '.contracts.\"SystemContract.sol:SystemContract\"'  > SystemContract.json"
//go:generate sh -c "cat SystemContract.json | jq .abi | abigen --abi - --pkg zevm --type SystemContract --out SystemContract.go"
//go:generate sh -c "cat UniswapV2Pair.json | jq .abi | abigen --abi - --pkg zevm --type UniswapV2Pair --out UniswapV2Pair.go"
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
var _ = ZRC20{}
var _ = UniswapV2Factory{}
var _ = SystemContract{}
var _ = UniswapV2Router02{}

type CompiledContract struct {
	ABI abi.ABI
	Bin evmtypes.HexString
}

var (
	//go:embed ZRC20.json
	ZRC20JSON []byte // nolint: golint
	//go:embed UniswapV2Factory.json
	UniswapV2FactoryJSON []byte // nolint: golint
	//go:embed WZETA.json
	WZETAJSON []byte // nolint: golint
	//go:embed SystemContract.json
	SystemContractJSON []byte // nolint: golint
	//go:embed UniswapV2Router02.json
	UniswapV2Router02JSON []byte // nolint: golint

	ZRC20Contract             CompiledContract
	UniswapV2FactoryContract  CompiledContract
	WZETAContract             CompiledContract
	SystemContractContract    CompiledContract
	UniswapV2Router02Contract CompiledContract

	// the module address of zetacore; no private exists.
	ZRC20AdminAddress ethcommon.Address
)

func init() {
	ZRC20AdminAddress = fungibletypes.ModuleAddressEVM

	err := json.Unmarshal(ZRC20JSON, &ZRC20Contract)
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
	err = json.Unmarshal(UniswapV2Router02JSON, &UniswapV2Router02Contract)
	if err != nil {
		panic(err)
	}

	if len(ZRC20Contract.Bin) == 0 {
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

	if len(UniswapV2Router02Contract.Bin) == 0 {
		panic("load contract failed")
	}
}
