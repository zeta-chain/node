package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestV2ERC20Deposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ERC20Deposit")

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	oldBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	// perform the deposit
	tx := r.V2ERC20Deposit(r.EVMAddress(), amount, gatewayevm.RevertOptions{})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the balance was updated
	newBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r, new(big.Int).Add(oldBalance, amount), newBalance)

	// add liquidity, 50ZETA/50ETH and 50ZETA/50ERC20
	fifty := big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(50))
	r.AddLiquidityETH(fifty, fifty)
	r.AddLiquidityERC20(fifty, fifty)
}
