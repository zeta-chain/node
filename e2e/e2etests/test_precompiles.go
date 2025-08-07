package e2etests

import (
	"math/big"
	"strings"

	"cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/istaking"
	"github.com/zeta-chain/node/e2e/runner"
)

// source: https://github.com/cosmos/evm/blob/main/precompiles/staking/StakingI.sol
const stakingAddr = "0x0000000000000000000000000000000000000800"

func TestPrecompiles(r *runner.E2ERunner, _ []string) {
	precompile, err := istaking.NewIStaking(ethcommon.HexToAddress(stakingAddr), r.ZEVMClient)
	require.NoError(r, err)

	r.Logger.Print("Trying to perform staking precompile test...")

	// get validators
	res, err := precompile.Validators(&bind.CallOpts{}, "", istaking.PageRequest{})
	require.NoError(r, err)
	require.Len(r, res.Validators, 2)
	r.Logger.Print("Validators found: %d", len(res.Validators))
	for i, v := range res.Validators {
		r.Logger.Print("Validator %d: %s, shares: %s", i, v.OperatorAddress, v.DelegatorShares.String())
	}

	validator := getValidatorOpAddress(r)

	// Get the delegation to the first validator
	delegationBefore := getDelegation(r)

	// Deledate 1 ZETA
	// Note: the call fails here with a `no contract code at given address` error, calling the precompile directly works
	// TODO: fix the issue with the precompile call
	// https://github.com/zeta-chain/node/issues/4081
	tx, err := precompile.Delegate(r.ZEVMAuth, r.ZEVMAuth.From, validator, big.NewInt(1e18))
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(tx)

	// Get the delegation to the first validator after delegation
	// The share denomination of the delegation should increase
	delegationAfter := getDelegation(r)
	require.True(
		r,
		delegationAfter.GT(delegationBefore),
		"Delegation shares should increase after delegation, before: %s, after: %s",
		delegationBefore.String(),
		delegationAfter.String(),
	)

	r.Logger.Print("Staking precompile test passed successfully!")
}

// getDelegation returns the delegation to the first validator from the active E2E account
func getDelegation(r *runner.E2ERunner) math.LegacyDec {
	// Get the first validator's operator address
	opAddress := getValidatorOpAddress(r)

	// Get the delegation address for the first validator
	res, err := r.StakingClient.Delegation(r.Ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: r.Account.RawBech32Address.String(),
		ValidatorAddr: opAddress,
	})
	if err != nil {
		// no delegation found, return zero shares
		if stakingtypes.ErrNoDelegation.Is(err) || strings.Contains(err.Error(), "not found") {
			return math.LegacyZeroDec()
		}
		require.Fail(r, "Failed to get delegation", "error: %v", err)
	}

	return res.DelegationResponse.Delegation.Shares
}

// getValidatorOpAddress returns the operator address of the first validator to perform staking tests
func getValidatorOpAddress(r *runner.E2ERunner) string {
	res, err := r.StakingClient.Validators(r.Ctx, &stakingtypes.QueryValidatorsRequest{})
	require.NoError(r, err)
	require.NotEmpty(r, res.Validators)

	// Ensure the validator is in bonded state for delegation
	validator := res.Validators[0]
	require.Equal(r, stakingtypes.Bonded, validator.Status)

	return res.Validators[0].OperatorAddress
}
