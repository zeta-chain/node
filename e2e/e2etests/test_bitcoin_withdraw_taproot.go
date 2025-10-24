package e2etests

import (
	"math/big"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinWithdrawTaproot(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1pqqqsyqcyq5rqwzqfpg9scrgwpugpzysnzs23v9ccrydpk8qarc0sj9hjuh"
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressTaproot)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawTaproot.")

	r.WithdrawBTCAndWaitCCTX(
		receiver,
		amount,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		crosschaintypes.CctxStatus_OutboundMined,
	)
}
