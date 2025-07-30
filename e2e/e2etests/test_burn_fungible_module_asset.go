package e2etests

import (
	"math/big"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/constant"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

func TestBurnFungibleModuleAsset(r *runner.E2ERunner, _ []string) {
	testBurnFungibleModuleAssetZRC20(r)
	testBurnFungibleModuleAssetZETA(r)
}

func testBurnFungibleModuleAssetZRC20(r *runner.E2ERunner) {
	// get fungible module balance
	balance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, fungibletypes.ModuleAddressEVM)
	require.NoError(r, err)

	// if the balance is zero, we need to deposit some tokens
	if balance.Uint64() == 0 {
		tx, err := r.ETHZRC20.Transfer(r.ZEVMAuth, fungibletypes.ModuleAddressEVM, big.NewInt(1000))
		require.NoError(r, err)

		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)

		balance, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, fungibletypes.ModuleAddressEVM)
		require.NoError(r, err)
		require.NotZero(r, balance.Uint64())
	}

	r.Logger.Info("Sending message to burn fungible module asset")
	msg := fungibletypes.NewMsgBurnFungibleModuleAsset(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ETHZRC20Addr.Hex(),
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)
	r.Logger.Info("Burn fungible module asset tx hash: %s", res.TxHash)

	// check the balance of the fungible module after burn
	balance, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, fungibletypes.ModuleAddressEVM)
	require.NoError(r, err)
	require.Zero(r, balance.Uint64())
}

func testBurnFungibleModuleAssetZETA(r *runner.E2ERunner) {
	// get fungible module balance
	res, err := r.BankClient.SpendableBalanceByDenom(r.Ctx, &banktypes.QuerySpendableBalanceByDenomRequest{
		Address: fungibletypes.ModuleAddress.String(),
		Denom:   config.BaseDenom,
	})
	require.NoError(r, err)
	balance := res.Balance.Amount

	// if the balance is zero, we need to deposit some tokens
	if balance.IsZero() {
		err := r.ZetaTxServer.TransferZETA(fungibletypes.ModuleAddress, 1000)
		require.NoError(r, err)

		// check balance
		res, err := r.BankClient.SpendableBalanceByDenom(r.Ctx, &banktypes.QuerySpendableBalanceByDenomRequest{
			Address: fungibletypes.ModuleAddress.String(),
			Denom:   config.BaseDenom,
		})
		require.NoError(r, err)
		balance := res.Balance.Amount
		require.NotZero(r, balance.Int64())
	}

	r.Logger.Info("Sending message to burn fungible module asset")
	msg := fungibletypes.NewMsgBurnFungibleModuleAsset(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		constant.EVMZeroAddress,
	)
	resBurn, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)
	r.Logger.Info("Burn fungible module asset tx hash: %s", resBurn.TxHash)

	// check balance
	res, err = r.BankClient.SpendableBalanceByDenom(r.Ctx, &banktypes.QuerySpendableBalanceByDenomRequest{
		Address: fungibletypes.ModuleAddress.String(),
		Denom:   config.BaseDenom,
	})
	require.NoError(r, err)
	balance = res.Balance.Amount
	require.Zero(r, balance.Int64())
}
