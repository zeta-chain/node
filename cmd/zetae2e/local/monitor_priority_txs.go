package local

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/cometbft/cometbft/abci/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"

	"github.com/zeta-chain/node/e2e/config"
)

var errWrongTxPriority = errors.New("wrong tx priority, system tx not on top")

// monitorTxPriorityInBlocks checks for transaction priorities in blocks and reports errors
func monitorTxPriorityInBlocks(ctx context.Context, conf config.Config, errCh chan error) {
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
			processBlockTxs(ctx, rpcClient, errCh)
		}
	}
}

// processBlockTxs fetches the latest block and evaluates transaction priorities
func processBlockTxs(ctx context.Context, rpc *rpchttp.HTTP, errCh chan error) {
	block, err := rpc.Block(ctx, nil)
	if err != nil {
		errCh <- err
		return
	}

	nonSystemTxFound := false
	for _, tx := range block.Block.Txs {
		txResult, err := rpc.Tx(ctx, tx.Hash(), false)
		if err != nil {
			continue
		}
		processTx(txResult, &nonSystemTxFound, errCh)
	}
}

// processTx handles the processing of each transaction and its events
func processTx(txResult *coretypes.ResultTx, nonSystemTxFound *bool, errCh chan error) {
	for _, event := range txResult.TxResult.Events {
		for _, attr := range event.Attributes {
			// skip attrs with empty value
			if attr.Value == "\"\"" {
				continue
			}

			// skip internal events with msg_type_url key, because they are not representing sdk msgs
			if strings.Contains(attr.Value, ".internal.") {
				continue
			}
			switch attr.Key {
			// if attr key is msg_type_url, check if it's one of system txs, otherwise mark it as non system tx
			case "msg_type_url":
				if isMsgTypeURLSystemTx(attr) {
					// a non system tx has been found in the block before a system tx
					if *nonSystemTxFound {
						errCh <- errWrongTxPriority
						return
					}
				} else {
					*nonSystemTxFound = true
				}
			// if attr key is action, check if tx is evm non system tx and if it is, mark it
			case "action":
				if isActionNonSystemTx(attr) {
					*nonSystemTxFound = true
				}
			}
		}
	}
}

func isMsgTypeURLSystemTx(attr types.EventAttribute) bool {
	// type urls in attr.Value are in double quotes, so it needs to be formatted like this
	systemTxsMsgTypeUrls := []string{
		"\"/zetachain.zetacore.crosschain.MsgVoteOutbound\"",
		"\"/zetachain.zetacore.crosschain.MsgVoteInbound\"",
		"\"/zetachain.zetacore.crosschain.MsgVoteGasPrice\"",
		"\"/zetachain.zetacore.crosschain.MsgAddOutboundTracker\"",
		"\"/zetachain.zetacore.crosschain.MsgAddInboundTracker\"",
		"\"/zetachain.zetacore.observer.MsgVoteBlockHeader\"",
		"\"/zetachain.zetacore.observer.MsgVoteTSS\"",
		"\"/zetachain.zetacore.observer.MsgVoteBlame\"",
	}

	return slices.Contains(systemTxsMsgTypeUrls, attr.Value)
}

func isActionNonSystemTx(attr types.EventAttribute) bool {
	return attr.Value == "/cosmos.evm.vm.v1.MsgEthereumTx"
}
