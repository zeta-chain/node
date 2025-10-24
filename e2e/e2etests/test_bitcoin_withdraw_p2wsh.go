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

func TestBitcoinWithdrawP2WSH(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse arguments and withdraw BTC
	defaultReceiver := "bcrt1qm9mzhyky4w853ft2ms6dtqdyyu3z2tmrq8jg8xglhyuv0dsxzmgs2f0sqy"
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressWitnessScriptHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawP2WSH.")

	r.WithdrawBTCAndWaitCCTX(
		receiver,
		amount,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		crosschaintypes.CctxStatus_OutboundMined,
	)
}
