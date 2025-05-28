package e2etests

import (
	"math/big"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinWithdrawRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 3)

	// ARRANGE
	// Given amount, restricted BTC P2WPKH address, revert address
	addressRestricted, err := chains.DecodeBtcAddress(args[0], chains.BitcoinRegtest.ChainId)
	require.NoError(r, err)
	amountStr := utils.ParseFloat(r, args[1])
	amount := utils.BTCAmountFromFloat64(r, amountStr)
	revertAddress := ethcommon.HexToAddress(args[2])

	// balance before
	revertBalanceBefore, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// ACT
	// the cctx should be reverted
	rawTx := r.WithdrawBTCAndWaitCCTX(
		addressRestricted,
		amount,
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
		crosschaintypes.CctxStatus_Reverted,
	)

	// ASSERT
	// normal withdraw tx has 3 outputs: [nonce, payment, change], but a cancel tx should have only 2 outputs: [nonce, change]
	require.Len(r, rawTx.Vout, 2, "BTC cancelled tx should have 2 outputs")

	// receiver balance should not change
	unspent, err := r.BtcRPCClient.ListUnspentMinMaxAddresses(r.Ctx, 1, 9999999, []btcutil.Address{addressRestricted})
	require.NoError(r, err)
	require.Empty(r, unspent)

	// revert address should receive the amount
	revertBalanceAfter, err := r.BTCZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	require.EqualValues(r, new(big.Int).Add(revertBalanceBefore, amount), revertBalanceAfter)
}
