package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// NOTE: to be used on live networks to check if multiple deposits using legacy method are reverting
func TestETHMultipleDepositsLegacy(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.Logger.Info("starting eth multiple legacy deposits test")

	oldBalance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	// set value of the payable transactions
	previousValue := r.EVMAuth.Value

	// send multiple legacy deposits through contract with 0 fee should revert
	_, err = r.TestDAppV2EVM.GatewayMultipleDepositsLegacy(r.EVMAuth, r.TestDAppV2ZEVMAddr, big.NewInt(0))
	require.Error(r, err)

	// send multiple legacy deposits through contract with correct fee should also revert
	fee, err := r.GatewayEVM.AdditionalActionFeeWei(nil)
	require.NoError(r, err)
	// add fee to provided amount to pay for 2 inbounds (1st one is free)
	r.EVMAuth.Value = new(big.Int).Add(amount, new(big.Int).Mul(fee, big.NewInt(2)))

	_, err = r.TestDAppV2EVM.GatewayMultipleDepositsLegacy(r.EVMAuth, r.TestDAppV2ZEVMAddr, fee)
	require.Error(r, err)

	r.EVMAuth.Value = previousValue

	// verify balance was not updated
	newBalance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)
	require.True(r, newBalance.Cmp(oldBalance) == 0)
}
