package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

const defaultReceiver = "mxpYha3UJKUgSwsAz2qYRqaDSwAkKZ3YEY"

func WithdrawBitcoinMultipleTimes(r *runner.E2ERunner, args []string) {
	// ARRANGE
	// Given amount and repeat count
	require.Len(r, args, 2)
	var (
		amount = utils.BTCAmountFromFloat64(r, utils.ParseFloat(r, args[0]))
		times  = utils.ParseInt(r, args[1])
	)

	// Given a receiver
	receiver, err := chains.DecodeBtcAddress(defaultReceiver, r.GetBitcoinChainID())
	require.NoError(r, err)

	// ACT
	for i := 0; i < times; i++ {
		r.WithdrawBTCAndWaitCCTX(
			receiver,
			amount,
			gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
			crosschaintypes.CctxStatus_OutboundMined,
		)
	}
}
