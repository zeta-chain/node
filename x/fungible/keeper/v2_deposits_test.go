package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/testdappv2"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
	"math/big"
	"testing"
)

// deployTestDAppV2 deploys the test dapp v2 contract and returns its address
func deployTestDAppV2(t *testing.T, ctx sdk.Context, k *fungiblekeeper.Keeper, evmk types.EVMKeeper) common.Address {
	testDAppV2, err := k.DeployContract(ctx, testdappv2.TestDAppV2MetaData)
	require.NoError(t, err)
	require.NotEmpty(t, testDAppV2)
	assertContractDeployment(t, evmk, ctx, testDAppV2)

	return testDAppV2
}

// assertTestDAppV2MessageAndAmount asserts the message and amount of the test dapp v2 contract
func assertTestDAppV2MessageAndAmount(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	contract common.Address,
	expectedMessage string,
	expectedAmount int64,
) {
	testDAppABI, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(t, err)

	// message
	res, err := k.CallEVM(
		ctx,
		*testDAppABI,
		types.ModuleAddressEVM,
		contract,
		fungiblekeeper.BigIntZero,
		nil,
		false,
		false,
		"getCalledWithMessage",
		expectedMessage,
	)
	require.NoError(t, err)

	unpacked, err := testDAppABI.Unpack("getCalledWithMessage", res.Ret)
	require.NoError(t, err)
	require.Len(t, unpacked, 1)
	found, ok := unpacked[0].(bool)
	require.True(t, ok)
	require.True(t, found)

	// amount
	res, err = k.CallEVM(
		ctx,
		*testDAppABI,
		types.ModuleAddressEVM,
		contract,
		fungiblekeeper.BigIntZero,
		nil,
		false,
		false,
		"getAmountWithMessage",
		expectedMessage,
	)
	require.NoError(t, err)

	unpacked, err = testDAppABI.Unpack("getAmountWithMessage", res.Ret)
	require.NoError(t, err)
	require.Len(t, unpacked, 1)
	amount, ok := unpacked[0].(*big.Int)
	require.True(t, ok)
	require.Equal(t, expectedAmount, amount.Int64())
}

func TestKeeper_ProcessV2Deposit(t *testing.T) {
	t.Run("should process no-call deposit", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId
		receiver := sample.EthAddress()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, contractCall, err := k.ProcessV2Deposit(
			ctx,
			sample.EthAddress().Bytes(),
			chainID,
			zrc20,
			receiver,
			big.NewInt(42),
			[]byte{},
			coin.CoinType_Gas,
		)

		// ASSERT
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, receiver)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("should process deposit and call", func(t *testing.T) {
		// ARRANGE
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainID := chains.DefaultChainsList()[0].ChainId

		// deploy test dapp
		testDapp := deployTestDAppV2(t, ctx, k, sdkk.EvmKeeper)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "foobar", "foobar")

		// ACT
		_, contractCall, err := k.ProcessV2Deposit(
			ctx,
			sample.EthAddress().Bytes(),
			chainID,
			zrc20,
			testDapp,
			big.NewInt(82),
			[]byte("foo"),
			coin.CoinType_Gas,
		)

		// ASSERT
		require.NoError(t, err)
		require.True(t, contractCall)
		balance, err := k.BalanceOfZRC4(ctx, zrc20, testDapp)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(82), balance)
		assertTestDAppV2MessageAndAmount(t, ctx, k, testDapp, "foo", 82)
	})
}
