package keeper_test

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_SetGasPrice(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	_, _, _, _, system := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

	queryGasPrice := func(chainID *big.Int) *big.Int {
		abi, err := systemcontract.SystemContractMetaData.GetAbi()
		assert.NoError(t, err)
		res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, system, keeper.BigIntZero, nil, false, false, "gasPriceByChainId", chainID)
		assert.NoError(t, err)
		unpacked, err := abi.Unpack("gasPriceByChainId", res.Ret)
		assert.NoError(t, err)
		gasPrice, ok := unpacked[0].(*big.Int)
		assert.True(t, ok)
		return gasPrice
	}

	_, err := k.SetGasPrice(ctx, big.NewInt(1), big.NewInt(42))
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(42), queryGasPrice(big.NewInt(1)))
}

func TestKeeper_SetGasCoin(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	gas := sample.EthAddress()

	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	err := k.SetGasCoin(ctx, big.NewInt(1), gas)
	assert.NoError(t, err)

	found, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(1))
	assert.NoError(t, err)
	assert.Equal(t, gas.Hex(), found.Hex())
}

func TestKeeper_SetGasZetaPool(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
	zrc20 := sample.EthAddress()

	_, _, _, _, system := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

	queryZetaPool := func(chainID *big.Int) ethcommon.Address {
		abi, err := systemcontract.SystemContractMetaData.GetAbi()
		assert.NoError(t, err)
		res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, system, keeper.BigIntZero, nil, false, false, "gasZetaPoolByChainId", chainID)
		assert.NoError(t, err)
		unpacked, err := abi.Unpack("gasZetaPoolByChainId", res.Ret)
		assert.NoError(t, err)
		pool, ok := unpacked[0].(ethcommon.Address)
		assert.True(t, ok)
		return pool
	}

	err := k.SetGasZetaPool(ctx, big.NewInt(1), zrc20)
	assert.NoError(t, err)
	assert.NotEqual(t, ethcommon.Address{}, queryZetaPool(big.NewInt(1)))
}
