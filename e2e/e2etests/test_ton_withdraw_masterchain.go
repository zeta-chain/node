package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/chains/ton/signer"
)

func TestTONWithdrawMasterchain(r *runner.E2ERunner, args []string) {
	// ARRANGE
	require.Len(r, args, 1)

	// Given gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given zEVM sender
	zevmSender := r.ZEVMAuth.From

	// Given its ZRC-20 balance
	senderZRC20BalanceBefore, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, zevmSender)
	require.NoError(r, err)
	r.Logger.Info("zEVM sender's ZRC20 TON balance before withdraw: %d", senderZRC20BalanceBefore)

	// Given amount to withdraw
	amount := utils.ParseUint(r, args[0])

	r.Logger.Info("Amount to withdraw: %s", toncontracts.FormatCoins(amount))

	// Given a restricted receiver (masterchain -1)
	receiver := ton.MustParseAccountID("-1:fffffbd865df68188ea84d6615086c26a7b5912a60bc55fded2cdb029b67cdef")

	// ACT
	tx := r.SendWithdrawTONZRC20(receiver, amount.BigInt(), gatewayzevm.RevertOptions{
		RevertAddress:    r.EVMAddress(),
		OnRevertGasLimit: big.NewInt(0),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw_restricted")

	// ASSERT
	// Ensure cctx is reverted
	utils.RequireCCTXStatus(r, cctx, cctypes.CctxStatus_Reverted)

	// 1: ton cctx, 2: zevm revert tx
	require.Len(r, cctx.OutboundParams, 2, "expected 2 outbound params")

	// Let's query ton tx
	lt, hash, err := encoder.DecodeTx(cctx.OutboundParams[0].Hash)
	require.NoError(r, err)

	// And ensure that this is an "increase seqno" transaction
	tonTx, err := r.Clients.TON.GetTransaction(r.Ctx, gw.AccountID(), lt, hash)
	require.NoError(r, err)

	gwTx, err := gw.ParseTransaction(tonTx)
	require.NoError(r, err)

	increaseSeqno, err := gwTx.IncreaseSeqno()
	require.NoError(r, err)
	require.Equal(r, uint32(signer.InvalidWorkchain), increaseSeqno.ReasonCode)
}
