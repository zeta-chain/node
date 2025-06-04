package e2etests

import (
	"math/big"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestBitcoinWithdrawP2SH(r *runner.E2ERunner, args []string) {
	// check length of arguments
	require.Len(r, args, 2)

	// parse arguments and withdraw BTC
	defaultReceiver := "2N6AoUj3KPS7wNGZXuCckh8YEWcSYNsGbqd"
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressScriptHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawP2SH.")

	r.WithdrawBTCAndWaitCCTX(
		receiver,
		amount,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		crosschaintypes.CctxStatus_OutboundMined,
	)
}
