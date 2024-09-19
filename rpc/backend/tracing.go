// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/zeta-chain/ethermint/blob/main/LICENSE
package backend

import (
	"encoding/json"
	"fmt"

	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"

	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (b *Backend) TraceTransaction(hash common.Hash, config *evmtypes.TraceConfig) (interface{}, error) {
	// Get transaction by hash
	transaction, _, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.logger.Debug("tx not found", "hash", hash)
		return nil, err
	}

	// check if block number is 0
	if transaction.Height == 0 {
		return nil, errors.New("genesis is not traceable")
	}

	blk, err := b.TendermintBlockByNumber(rpctypes.BlockNumber(transaction.Height))
	if err != nil {
		b.logger.Debug("block not found", "height", transaction.Height)
		return nil, err
	}

	blockResult, err := b.TendermintBlockResultByNumber(&blk.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("block result not found for height %d", blk.Block.Height)
	}

	predecessors := []*evmtypes.MsgEthereumTx{}
	msgs, _ := b.EthMsgsFromTendermintBlock(blk, blockResult)
	var ethMsg *evmtypes.MsgEthereumTx
	for _, m := range msgs {
		if m.Hash == hash.Hex() {
			ethMsg = m
			break
		}
		predecessors = append(predecessors, m)
	}

	if ethMsg == nil {
		return nil, fmt.Errorf("tx not found in block %d", blk.Block.Height)
	}

	traceTxRequest := evmtypes.QueryTraceTxRequest{
		Msg:             ethMsg,
		Predecessors:    predecessors,
		BlockNumber:     blk.Block.Height,
		BlockTime:       blk.Block.Time,
		BlockHash:       common.Bytes2Hex(blk.BlockID.Hash),
		ProposerAddress: sdk.ConsAddress(blk.Block.ProposerAddress),
		ChainId:         b.chainID.Int64(),
	}

	if config != nil {
		traceTxRequest.TraceConfig = config
	}

	// minus one to get the context of block beginning
	contextHeight := transaction.Height - 1
	if contextHeight < 1 {
		// 0 is a special value in `ContextWithHeight`
		contextHeight = 1
	}
	traceResult, err := b.queryClient.TraceTx(rpctypes.ContextWithHeight(contextHeight), &traceTxRequest)
	if err != nil {
		return nil, err
	}

	// Response format is unknown due to custom tracer config param
	// More information can be found here https://geth.ethereum.org/docs/dapp/tracing-filtered
	var decodedResult interface{}
	err = json.Unmarshal(traceResult.Data, &decodedResult)
	if err != nil {
		return nil, err
	}

	return decodedResult, nil
}

// TraceBlock configures a new tracer according to the provided configuration, and
// executes all the transactions contained within. The return value will be one item
// per transaction, dependent on the requested tracer.
func (b *Backend) TraceBlock(height rpctypes.BlockNumber,
	config *evmtypes.TraceConfig,
	block *tmrpctypes.ResultBlock,
) ([]*evmtypes.TxTraceResult, error) {
	txs := block.Block.Txs
	txsLength := len(txs)

	if txsLength == 0 {
		// If there are no transactions return empty array
		return []*evmtypes.TxTraceResult{}, nil
	}

	blockRes, err := b.TendermintBlockResultByNumber(&block.Block.Height)
	if err != nil {
		b.logger.Debug("block result not found", "height", block.Block.Height, "error", err.Error())
		return nil, nil
	}
	msgs, _ := b.EthMsgsFromTendermintBlock(block, blockRes)

	// minus one to get the context at the beginning of the block
	contextHeight := height - 1
	if contextHeight < 1 {
		// 0 is a special value for `ContextWithHeight`.
		contextHeight = 1
	}
	ctxWithHeight := rpctypes.ContextWithHeight(int64(contextHeight))

	traceBlockRequest := &evmtypes.QueryTraceBlockRequest{
		Txs:             msgs,
		TraceConfig:     config,
		BlockNumber:     block.Block.Height,
		BlockTime:       block.Block.Time,
		BlockHash:       common.Bytes2Hex(block.BlockID.Hash),
		ProposerAddress: sdk.ConsAddress(block.Block.ProposerAddress),
		ChainId:         b.chainID.Int64(),
	}

	res, err := b.queryClient.TraceBlock(ctxWithHeight, traceBlockRequest)
	if err != nil {
		return nil, err
	}

	decodedResults := make([]*evmtypes.TxTraceResult, txsLength)
	if err := json.Unmarshal(res.Data, &decodedResults); err != nil {
		return nil, err
	}

	return decodedResults, nil
}
