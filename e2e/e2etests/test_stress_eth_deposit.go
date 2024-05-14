package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"golang.org/x/sync/errgroup"
)

// TestStressEtherDeposit tests the stressing deposit of ether
func TestStressEtherDeposit(r *runner.E2ERunner, args []string) {
	if len(args) != 2 {
		panic("TestStressEtherDeposit requires exactly two arguments: the deposit amount and the number of deposits.")
	}

	depositAmount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid deposit amount specified for TestMultipleERC20Deposit.")
	}

	numDeposits, err := strconv.Atoi(args[1])
	if err != nil || numDeposits < 1 {
		panic("Invalid number of deposits specified for TestStressEtherDeposit.")
	}

	r.Logger.Print("starting stress test of %d deposits", numDeposits)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	// cancel MonitorTxPriorityInBlocks when test completes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		MonitorTxPriorityInBlocks(r, ctx)
	}()

	// send the deposits
	for i := 0; i < numDeposits; i++ {
		i := i
		hash := r.DepositEtherWithAmount(false, depositAmount)
		r.Logger.Print("index %d: starting deposit, tx hash: %s", i, hash.Hex())

		eg.Go(func() error {
			return MonitorEtherDeposit(r, hash, i, time.Now())
		})
	}

	// wait for all the deposits to complete
	if err := eg.Wait(); err != nil {
		panic(err)
	}

	r.Logger.Print("all deposits completed")
}

// MonitorEtherDeposit monitors the deposit of ether, returns once the deposit is complete
func MonitorEtherDeposit(r *runner.E2ERunner, hash ethcommon.Hash, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, hash.Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: deposit cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Now().Sub(startTime)
	r.Logger.Print("index %d: deposit cctx success in %s", index, timeToComplete.String())

	return nil
}

func MonitorTxPriorityInBlocks(r *runner.E2ERunner, ctx context.Context) {
	url := "http://zetacore0:26657" // TODO: init rpc client in runner?
	rpc, err := rpchttp.New(url, "/websocket")
	if err != nil {
		return
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			block, err := rpc.Block(ctx, nil)
			if err != nil {
				return
			}
			r.Logger.Print("block %d has %d txs", block.Block.Height, len(block.Block.Txs))

			// iterate txs events and check if some MsgEthereumTx is included above crosschain and observer txs
			nonSystemTxFound := false
			for _, tx := range block.Block.Txs {
				txRes, err := rpc.Tx(context.Background(), tx.Hash(), false)
				if err != nil {
					continue
				}

				for _, ev := range txRes.TxResult.Events {
					for _, attr := range ev.Attributes {
						if attr.Key == "msg_type_url" {
							if strings.Contains(attr.Value, "zetachain.zetacore.crosschain.Msg") ||
								strings.Contains(attr.Value, "zetachain.zetacore.observer.Msg") {
								r.Logger.Print("system tx %s", attr.Value)
								if nonSystemTxFound {
									// TODO: communicate failure to main routine
									r.Logger.Print("wrong priority")
									panic("wrong priority")
								}
							}
						}
						if attr.Key == "action" {
							if strings.Contains(attr.Value, "ethermint.evm.v1.MsgEthereumTx") {
								nonSystemTxFound = true
								r.Logger.Print("non system tx %s", attr.Value)
							}
						}
					}
				}
			}
		}
	}
}
