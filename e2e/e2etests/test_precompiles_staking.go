package e2etests

import (
	"fmt"
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/teststaking"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/staking"
)

func TestPrecompilesStakingThroughContract(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	testStakingAddr, testStakingTx, _, err := teststaking.DeployTestStaking(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, testStakingTx, r.Logger, r.ReceiptTimeout)

	testStaking, err := teststaking.NewTestStaking(testStakingAddr, r.ZEVMClient)
	require.NoError(r, err)

	_, err = teststaking.NewTestStaking(testStakingAddr, r.ZEVMClient)
	require.NoError(r, err)

	stakingContract, err := staking.NewIStaking(staking.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create staking contract caller")
	validators, err := stakingContract.GetAllValidators(nil)
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	fmt.Println("from", r.ZEVMAuth.From)
	fmt.Println("test staking contract", testStakingAddr.String(), validators[0].OperatorAddress)

	_, err = stakingContract.GetShares(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)

	_, err = testStaking.Bech32StaticFn(nil, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u")
	require.NoError(r, err)

	_, err = testStaking.Bech32CallFn(r.ZEVMAuth, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u")
	require.NoError(r, err)

	_, err = testStaking.Bech32Fn(nil, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u")
	require.NoError(r, err)

	// _, err = testStaking.GetShares(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	// require.NoError(r, err)
	// tx, err := testStaking.Stake(r.ZEVMAuth, validators[0].OperatorAddress, big.NewInt(3))
	// require.NoError(r, err)
	// receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

}

func TestPrecompilesStaking(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	stakingContract, err := staking.NewIStaking(staking.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create staking contract caller")

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	validators, err := stakingContract.GetAllValidators(nil)
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// shares are 0 for both validators at the start
	sharesBeforeVal1, err := stakingContract.GetShares(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal1.Int64())

	sharesBeforeVal2, err := stakingContract.GetShares(nil, r.ZEVMAuth.From, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), sharesBeforeVal2.Int64())

	// stake 3 to validator1
	tx, err := stakingContract.Stake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(3))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares are set to 3
	sharesAfterVal1, err := stakingContract.GetShares(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(3e18).String(), sharesAfterVal1.String())

	// unstake 1 from validator1
	tx, err = stakingContract.Unstake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(1))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares are set to 2
	sharesAfterVal1, err = stakingContract.GetShares(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(2e18).String(), sharesAfterVal1.String())

	// move 1 stake from validator1 to validator2
	tx, err = stakingContract.MoveStake(
		r.ZEVMAuth,
		r.ZEVMAuth.From,
		validators[0].OperatorAddress,
		validators[1].OperatorAddress,
		big.NewInt(1),
	)
	require.NoError(r, err)

	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check shares for both validator1 and validator2 are 1
	sharesAfterVal1, err = stakingContract.GetShares(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal1.String())

	sharesAfterVal2, err := stakingContract.GetShares(nil, r.ZEVMAuth.From, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), sharesAfterVal2.String())
}
