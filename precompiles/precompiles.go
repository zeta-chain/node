package precompiles

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"
	evmkeeper "github.com/zeta-chain/ethermint/x/evm/keeper"

	"github.com/zeta-chain/zetacore/precompiles/prototype"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
)

// EnabledStatefulContracts contains the list of all enabled stateful precompiles.
// This is useful for listing and reading from other packages, such as BlockedAddrs() function.
// Setting to false a contract here will disable it, not being included in the blockchain.
var EnabledStatefulContracts = map[common.Address]bool{
	prototype.ContractAddress: true,
}

// StatefulContracts returns all the registered precompiled contracts.
func StatefulContracts(
	fungibleKeeper *keeper.Keeper,
	cdc codec.Codec,
	gasConfig storetypes.GasConfig,
) (precompiledContracts []evmkeeper.CustomContractFn) {
	// Initialize at 0 the custom compiled contracts and the addresses.
	precompiledContracts = make([]evmkeeper.CustomContractFn, 0)

	// Define the regular contract function.
	if EnabledStatefulContracts[prototype.ContractAddress] {
		prototypeContract := func(_ sdktypes.Context, _ ethparams.Rules) vm.PrecompiledContract {
			return prototype.NewIPrototypeContract(fungibleKeeper, cdc, gasConfig)
		}

		// Append the regular contract to the precompiledContracts slice.
		precompiledContracts = append(precompiledContracts, prototypeContract)
	}

	return precompiledContracts
}