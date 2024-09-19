package e2etests

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/e2e/contracts/teststaking"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestPrecompilesStakingThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	testStakingAddr, testStakingTx, testStaking, err := teststaking.DeployTestStaking(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.WZetaAddr,
	)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, testStakingTx, r.Logger, r.ReceiptTimeout)

	validators, err := testStaking.GetAllValidators(&bind.CallOpts{})
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// shares are 0 for both validators at the start
	sharesBeforeVal1, err := testStaking.GetShares(&bind.CallOpts{}, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal1.Int64())

	sharesBeforeVal2, err := testStaking.GetShares(&bind.CallOpts{}, testStakingAddr, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal2.Int64())

	// stake 3 to validator1 using user_precompile account should fail because sender is not provided staker
	_, err = testStaking.Stake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(3))
	require.Error(r, err)

	// stake 3 to validator1 using testStaking smart contract should fail because it doesn't have any balance yet
	_, err = testStaking.Stake(r.ZEVMAuth, testStakingAddr, validators[0].OperatorAddress, big.NewInt(3))
	require.Error(r, err)

	r.ZEVMAuth.Value = big.NewInt(1000000000000)
	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	// fund testStaking contract with azeta
	tx, err := testStaking.DepositWZETA(r.ZEVMAuth)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.ZEVMAuth.Value = big.NewInt(0)

	stakeAmount := 100000000000
	tx, err = testStaking.WithdrawWZETA(r.ZEVMAuth, big.NewInt(int64(stakeAmount)))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// bank balance at the start
	balanceBefore, err := r.BankClient.Balance(r.Ctx, &banktypes.QueryBalanceRequest{
		Address: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		Denom:   config.BaseDenom,
	})
	require.NoError(r, err)
	require.Equal(r, int64(stakeAmount), balanceBefore.Balance.Amount.Int64())

	// stake 3 to validator1 and revert in same function
	tx, err = testStaking.StakeAndRevert(r.ZEVMAuth, testStakingAddr, validators[0].OperatorAddress, big.NewInt(3))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check that bank balance was not changed because of revert in testStaking contract
	balanceAfterRevert, err := r.BankClient.Balance(r.Ctx, &banktypes.QueryBalanceRequest{
		Address: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		Denom:   config.BaseDenom,
	})
	require.NoError(r, err)
	require.Equal(r, balanceBefore.Balance.Amount.Int64(), balanceAfterRevert.Balance.Amount.Int64())

	// check that counter was not updated
	counter, err := testStaking.Counter(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, int64(0), counter.Int64())

	// check that shares are still 0
	sharesAfterRevert, err := testStaking.GetShares(&bind.CallOpts{}, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesAfterRevert.Int64())

	// stake 1 to validator1 using testStaking smart contract without smart contract state update
	tx, err = testStaking.Stake(r.ZEVMAuth, testStakingAddr, validators[0].OperatorAddress, big.NewInt(1))
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check that stake event was emitted
	stakeEvent, err := testStaking.ParseStake(*receipt.Logs[0])
	require.NoError(r, err)
	expectedValAddr, err := sdk.ValAddressFromBech32(validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1).Uint64(), stakeEvent.Amount.Uint64())
	require.Equal(r, common.BytesToAddress(expectedValAddr.Bytes()), stakeEvent.Validator)
	require.Equal(r, testStakingAddr, stakeEvent.Staker)

	// check that bank balance is reduced by 1
	balanceAfterStake, err := r.BankClient.Balance(r.Ctx, &banktypes.QueryBalanceRequest{
		Address: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		Denom:   config.BaseDenom,
	})
	require.NoError(r, err)
	require.Equal(r, balanceBefore.Balance.Amount.Int64()-1, balanceAfterStake.Balance.Amount.Int64())

	// stake 2 more to validator1 using testStaking smart contract with smart contract state update
	tx, err = testStaking.StakeWithStateUpdate(
		r.ZEVMAuth,
		testStakingAddr,
		validators[0].OperatorAddress,
		big.NewInt(2),
	)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check that stake event was emitted
	stakeEvent, err = testStaking.ParseStake(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, big.NewInt(2).Uint64(), stakeEvent.Amount.Uint64())
	require.Equal(r, common.BytesToAddress(expectedValAddr.Bytes()), stakeEvent.Validator)
	require.Equal(r, testStakingAddr, stakeEvent.Staker)

	// check that bank balance is reduced by 2 more, 3 in total
	balanceAfterStake, err = r.BankClient.Balance(r.Ctx, &banktypes.QueryBalanceRequest{
		Address: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		Denom:   config.BaseDenom,
	})
	require.NoError(r, err)
	require.Equal(r, balanceBefore.Balance.Amount.Int64()-3, balanceAfterStake.Balance.Amount.Int64())

	// check that counter is updated
	counter, err = testStaking.Counter(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, int64(2), counter.Int64())

	// check shares are set to 3
	sharesAfterVal1, err := testStaking.GetShares(&bind.CallOpts{}, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(3e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err := r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(3), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	// unstake 1 from validator1
	tx, err = testStaking.Unstake(r.ZEVMAuth, testStakingAddr, validators[0].OperatorAddress, big.NewInt(1))
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check that unstake event was emitted
	unstakeEvent, err := testStaking.ParseUnstake(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1).Uint64(), unstakeEvent.Amount.Uint64())
	require.Equal(r, common.BytesToAddress(expectedValAddr.Bytes()), unstakeEvent.Validator)
	require.Equal(r, testStakingAddr, unstakeEvent.Staker)

	// check shares are set to 2
	sharesAfterVal1, err = testStaking.GetShares(&bind.CallOpts{}, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(2e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err = r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(2), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	// move 1 stake from validator1 to validator2
	tx, err = testStaking.MoveStake(
		r.ZEVMAuth,
		testStakingAddr,
		validators[0].OperatorAddress,
		validators[1].OperatorAddress,
		big.NewInt(1),
	)
	require.NoError(r, err)
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check that moveStake event was emitted
	moveStake, err := testStaking.ParseMoveStake(*receipt.Logs[0])
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1).Uint64(), moveStake.Amount.Uint64())
	expectedValDstAddr, err := sdk.ValAddressFromBech32(validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, common.BytesToAddress(expectedValAddr.Bytes()), moveStake.ValidatorSrc)
	require.Equal(r, common.BytesToAddress(expectedValDstAddr.Bytes()), moveStake.ValidatorDst)
	require.Equal(r, testStakingAddr, moveStake.Staker)

	// check shares for both validator1 and validator2 are 1
	sharesAfterVal1, err = testStaking.GetShares(&bind.CallOpts{}, testStakingAddr, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal1.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal1, err = r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[0].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(1), delegationAfterVal1.DelegationResponse.Balance.Amount.Int64())

	sharesAfterVal2, err := testStaking.GetShares(&bind.CallOpts{}, testStakingAddr, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal2.String())

	// check delegation amount using staking keeper query client
	delegationAfterVal2, err := r.StakingClient.Delegation(r.Ctx, &types.QueryDelegationRequest{
		DelegatorAddr: sdk.AccAddress(testStakingAddr.Bytes()).String(),
		ValidatorAddr: validators[1].OperatorAddress,
	})
	require.NoError(r, err)
	require.Equal(r, int64(1), delegationAfterVal2.DelegationResponse.Balance.Amount.Int64())
}
