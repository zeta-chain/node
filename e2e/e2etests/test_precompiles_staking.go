package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/precompiles/staking"
)

func TestPrecompilesStaking(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	stakingContract, err := staking.NewIStaking(staking.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create staking contract caller")

	r.ZEVMAuth.GasLimit = 10000000

	validators, err := stakingContract.GetAllValidators(nil)
	require.NoError(r, err)
	require.GreaterOrEqual(r, len(validators), 2)

	// stakes are 0 for both validators at the start
	stakesBeforeVal1, err := stakingContract.GetStakes(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), stakesBeforeVal1.Int64())

	stakesBeforeVal2, err := stakingContract.GetStakes(nil, r.ZEVMAuth.From, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, int64(0), stakesBeforeVal2.Int64())

	// stake 3 to validator1
	tx, err := stakingContract.Stake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(3))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check stakes are set to 3
	stakesAfterVal1, err := stakingContract.GetStakes(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(3e18).String(), stakesAfterVal1.String())

	// unstake 1 from validator1
	tx, err = stakingContract.Unstake(r.ZEVMAuth, r.ZEVMAuth.From, validators[0].OperatorAddress, big.NewInt(1))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check stakes are set to 2
	stakesAfterVal1, err = stakingContract.GetStakes(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(2e18).String(), stakesAfterVal1.String())

	// transfer 1 stake from validator1 to validator2
	tx, err = stakingContract.TransferStake(
		r.ZEVMAuth,
		r.ZEVMAuth.From,
		validators[0].OperatorAddress,
		validators[1].OperatorAddress,
		big.NewInt(1),
	)
	require.NoError(r, err)

	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// check stakes for both validator1 and validator2 are 1
	stakesAfterVal1, err = stakingContract.GetStakes(nil, r.ZEVMAuth.From, validators[0].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), stakesAfterVal1.String())

	stakesAfterVal2, err := stakingContract.GetStakes(nil, r.ZEVMAuth.From, validators[1].OperatorAddress)
	require.NoError(r, err)
	require.Equal(r, big.NewInt(1e18).String(), stakesAfterVal2.String())
}
