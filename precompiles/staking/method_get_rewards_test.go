package staking

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/emissions"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func Test_GetRewards(t *testing.T) {
	t.Run("become azeta staker, distribute ZRC20, get rewards", func(t *testing.T) {
		/* ARRANGE */
		s := newTestSuite(t)

		// Create validator.
		validator := sample.Validator(t, rand.New(rand.NewSource(42)))
		s.sdkKeepers.StakingKeeper.SetValidator(s.ctx, validator)

		// Create staker.
		staker := sample.Bech32AccAddress()

		// Become a staker.
		stakeThroughCosmosAPI(
			t,
			s.ctx,
			s.sdkKeepers.BankKeeper,
			s.sdkKeepers.StakingKeeper,
			validator,
			staker,
			math.NewInt(100),
		)

		/* Distribute 1000 ZRC20 tokens to the staking contract */
		distributeZRC20(t, s, big.NewInt(1000))

		// Produce blocks.
		for i := 0; i < 10; i++ {
			// produce a block
			emissions.BeginBlocker(s.ctx, *s.sdkKeepers.EmissionsKeeper)
			s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 1)
		}

		/* ACT */
		// Call getRewards.
		getRewardsMethod := s.stkContractABI.Methods[GetRewardsMethodName]
	
		fmt.Println(common.HexToAddress(staker.String()))
		fmt.Println(validator.GetOperator().String())

		// Setup method input.
		s.mockVMContract.Input = packInputArgs(
			t,
			getRewardsMethod,
			[]interface{}{common.HexToAddress(staker.String()), validator.GetOperator().String()}...,
		)

		bytes, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
		require.NoError(t, err)

		res, err := s.stkContractABI.Methods[DistributeMethodName].Outputs.Unpack(bytes)
		require.NoError(t, err)
		fmt.Println(res)

		/* ASSERT */
	})
}

func stakeThroughCosmosAPI(
	t *testing.T,
	ctx sdk.Context,
	bankKeeper bankkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	validator stakingtypes.Validator,
	staker sdk.AccAddress,
	amount math.Int,
)  {
	// Coins to stake with default cosmos denom.
	coins := sdk.NewCoins(sdk.NewCoin("stake", amount))

	err := bankKeeper.MintCoins(ctx, fungibletypes.ModuleName, coins)
	require.NoError(t, err)

	err = bankKeeper.SendCoinsFromModuleToAccount(ctx, fungibletypes.ModuleName, staker, coins)
	require.NoError(t, err)

	b := bankKeeper.GetAllBalances(ctx, staker)
	fmt.Println(b)

	shares, err := stakingKeeper.Delegate(
		ctx,
		staker,
		coins.AmountOf(coins.Denoms()[0]),
		validator.Status,
		validator,
		true,
	)
	require.NoError(t, err)
	require.Equal(t, amount.Uint64(), shares.TruncateInt().Uint64())
	b = bankKeeper.GetAllBalances(ctx, staker)
	
	del, found := stakingKeeper.GetDelegation(ctx, staker, validator.GetOperator())
	fmt.Println(found)
	fmt.Println(del)
}

func distributeZRC20(
	t *testing.T,
	s testSuite,
	amount *big.Int,
) {
	distributeMethod := s.stkContractABI.Methods[DistributeMethodName]

	_, err := s.fungibleKeeper.DepositZRC20(s.ctx, s.zrc20Address, s.defaultCaller, amount)
	require.NoError(t, err)
	allowStaking(t, s, amount)

	// Setup method input.
	s.mockVMContract.Input = packInputArgs(
		t,
		distributeMethod,
		[]interface{}{s.zrc20Address, amount}...,
	)

	// Call distribute method.
	success, err := s.stkContract.Run(s.mockEVM, s.mockVMContract, false)
	require.NoError(t, err)
	res, err := distributeMethod.Outputs.Unpack(success)
	require.NoError(t, err)
	ok := res[0].(bool)
	require.True(t, ok)
}
