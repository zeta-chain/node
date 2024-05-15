package local

import (
	"context"
	"errors"
	"time"

	"github.com/cometbft/cometbft/abci/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/zeta-chain/zetacore/e2e/config"
)

// MonitorTxPriorityInBlocks checks for transaction priorities in blocks and reports errors
func MonitorTxPriorityInBlocks(ctx context.Context, conf config.Config, errCh chan error) {
	rpcClient, err := rpchttp.New(conf.RPCs.ZetaCoreRPC, "/websocket")
	if err != nil {
		errCh <- err
		return
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			errCh <- nil
		case <-ticker.C:
			checkBlockTransactions(ctx, rpcClient, errCh)
		}
	}
}

// checkBlockTransactions fetches the latest block and evaluates transaction priorities
func checkBlockTransactions(ctx context.Context, rpc *rpchttp.HTTP, errCh chan error) {
	block, err := rpc.Block(ctx, nil)
	if err != nil {
		errCh <- err
		return
	}

	nonSystemTxFound := false
	for _, tx := range block.Block.Txs {
		txResult, err := rpc.Tx(ctx, tx.Hash(), false)
		if err != nil {
			return
		}

		for _, event := range txResult.TxResult.Events {
			for _, attr := range event.Attributes {
				switch attr.Key {
				case "msg_type_url":
					if isMsgTypeUrlSystemTx(attr) {
						// a non system tx has been found in the block before a system tx
						if nonSystemTxFound {
							errCh <- errors.New("wrong tx priority, system tx not on top")
							return
						}
					} else {
						nonSystemTxFound = true
					}
				case "action":
					nonSystemTxFound = isActionNonSystemTx(attr)
				}
			}
		}
	}
}

func isMsgTypeUrlSystemTx(attr types.EventAttribute) bool {
	systemTxsMsgTypeUrls := []string{
		"/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
		"/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
		"/zetachain.zetacore.crosschain.MsgVoteGasPrice",
		"/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
		"/zetachain.zetacore.observer.MsgVoteBlockHeader",
		"/zetachain.zetacore.observer.MsgVoteTSS",
		"/zetachain.zetacore.observer.MsgAddBlameVote",
	}

	for _, url := range systemTxsMsgTypeUrls {
		if url == attr.Value {
			return true
		}
	}

	return false
}

func isActionNonSystemTx(attr types.EventAttribute) bool {
	return attr.Value == "/ethermint.evm.v1.MsgEthereumTx"
}
