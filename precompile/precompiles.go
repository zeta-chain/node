package precompiles

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"
	ethparams "github.com/ethereum/go-ethereum/params"
	evmkeeper "github.com/zeta-chain/ethermint/x/evm/keeper"
	"github.com/zeta-chain/zetacore/precompile/regular"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
)

func PrecompiledContracts(fungibleKeeper keeper.Keeper, cdc codec.Codec, gasConfig storetypes.GasConfig) []evmkeeper.CustomContractFn {
	// Initialize at 0 the custom compiled contracts.
	precompiledContracts := make([]evmkeeper.CustomContractFn, 0)

	// Define the regular contract function.
	regularContract := func(ctx sdktypes.Context, rules ethparams.Rules) vm.PrecompiledContract {
		return regular.NewRegularContract(fungibleKeeper, cdc, gasConfig)
	}

	// Append all the precompiled contracts to the precompiledContracts slice.
	precompiledContracts = append(precompiledContracts, regularContract)

	return precompiledContracts
}