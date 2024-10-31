package e2etests

import (
	"math/rand"
	"sync"

	"cosmossdk.io/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	cc "github.com/zeta-chain/node/x/crosschain/types"
)

func TestTONWithdrawConcurrent(r *runner.E2ERunner, _ []string) {
	// ARRANGE
	// Given a deployer
	_, deployer := r.Ctx, r.TONDeployer

	const recipientsCount = 10
	type withdrawal struct {
		recipient ton.AccountID
		amount    math.Uint
	}

	var (
		testCases []withdrawal
		wg        sync.WaitGroup
	)

	// Given multiple recipients WITHOUT deployed wallet-contracts
	// and sample withdrawal amounts between 1 and 5 TON
	for i := 0; i < recipientsCount; i++ {
		// #nosec G404: it's a test
		amount := 1 + rand.Intn(5)
		testCases = append(testCases, withdrawal{
			// #nosec G115 test - always in range
			amount:    toncontracts.Coins(uint64(amount)),
			recipient: sample.GenerateTONAccountID(),
		})
	}

	// ACT
	// Fire withdrawals. Note that zevm sender is r.ZEVMAuth
	for i, tc := range testCases {
		r.Logger.Info(
			"Withdrawal #%d: sending %s to %s",
			i+1,
			toncontracts.FormatCoins(tc.amount),
			tc.recipient.ToRaw(),
		)

		approvedAmount := tc.amount.Add(toncontracts.Coins(1))
		tx := r.SendWithdrawTONZRC20(tc.recipient, tc.amount.BigInt(), approvedAmount.BigInt())

		wg.Add(1)

		go func(number int, tx *ethtypes.Transaction) {
			defer wg.Done()

			// wait for the cctx to be mined
			cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)

			// ASSERT
			utils.RequireCCTXStatus(r, cctx, cc.CctxStatus_OutboundMined)
			r.Logger.Info("Withdrawal #%d complete! cctx index: %s", number, cctx.Index)

			// Check recipient's balance ON TON
			balance, err := deployer.GetBalanceOf(r.Ctx, tc.recipient, false)
			require.NoError(r, err, "failed to get balance of %s", tc.recipient.ToRaw())
			require.Equal(r, tc.amount.Uint64(), balance.Uint64())
		}(i+1, tx)
	}

	wg.Wait()
}
