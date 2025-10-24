package e2etests

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TODO: This test is similar to TestCrosschainSwap
// purpose is to test similar scenario with v2 contracts where there is swap + withdraw in onCall
// to showcase that it's not reverting with gas limit issues
// this test should be removed when this issue is completed: https://github.com/zeta-chain/node/issues/2711
func TestDepositAndCallSwap(r *runner.E2ERunner, _ []string) {
	// create tokens pair (erc20 and eth)
	tx, err := r.UniswapV2Factory.CreatePair(r.ZEVMAuth, r.ERC20ZRC20Addr, r.ETHZRC20Addr)
	if err != nil {
		r.Logger.Print("ℹ️ create pair error %s", err.Error())
	} else {
		utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	}

	// approve router to spend tokens being swapped
	tx, err = r.ERC20ZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx, err = r.ETHZRC20.Approve(r.ZEVMAuth, r.UniswapV2RouterAddr, big.NewInt(1e18))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// fund ZEVMSwapApp with gas ZRC20s for withdraw
	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, r.ZEVMSwapAppAddr, big.NewInt(1e10))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	tx, err = r.ERC20ZRC20.Transfer(r.ZEVMAuth, r.ZEVMSwapAppAddr, big.NewInt(1e6))
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// temporarily increase gas limit to 400000
	previousGasLimit := r.ZEVMAuth.GasLimit
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	// add liquidity for swap
	r.ZEVMAuth.GasLimit = 400000
	tx, err = r.UniswapV2Router.AddLiquidity(
		r.ZEVMAuth,
		r.ERC20ZRC20Addr,
		r.ETHZRC20Addr,
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e8),
		big.NewInt(1e5),
		r.EVMAddress(),
		big.NewInt(time.Now().Add(10*time.Minute).Unix()),
	)
	require.NoError(r, err)
	utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)

	// memobytes is dApp specific; see the contracts/ZEVMSwapApp.sol for details
	// it is [targetZRC20, receiver]
	memobytes, err := r.ZEVMSwapApp.EncodeMemo(
		&bind.CallOpts{},
		r.ETHZRC20Addr,
		r.EVMAddress().Bytes(),
	)
	require.NoError(r, err)

	// perform the deposit and call
	r.ApproveERC20OnEVM(r.GatewayEVMAddr)
	tx = r.ERC20DepositAndCall(
		r.ZEVMSwapAppAddr,
		big.NewInt(8e7),
		memobytes,
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
}
