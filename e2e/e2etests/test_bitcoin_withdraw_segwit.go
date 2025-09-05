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

func TestBitcoinWithdrawSegWit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse arguments
	defaultReceiver := r.GetBtcAddress().EncodeAddress()
	receiver, amount := utils.ParseBitcoinWithdrawArgs(r, args, defaultReceiver, r.GetBitcoinChainID())
	_, ok := receiver.(*btcutil.AddressWitnessPubKeyHash)
	require.True(r, ok, "Invalid receiver address specified for TestBitcoinWithdrawSegWit.")

	r.WithdrawBTCAndWaitCCTX(
		receiver,
		amount,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		crosschaintypes.CctxStatus_OutboundMined,
	)
}
