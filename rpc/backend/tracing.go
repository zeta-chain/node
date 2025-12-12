package backend

import (
	"encoding/json"
	"fmt"

	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (b *Backend) TraceTransaction(hash common.Hash, config *rpctypes.TraceConfig) (interface{}, error) {
	// Get transaction by hash
	transaction, _, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.Logger.Debug("tx not found", "hash", hash)
		return nil, err
	}

	// check if block number is 0
	if transaction.Height == 0 {
		return nil, errors.New("genesis is not traceable")
	}

	blk, err := b.TendermintBlockByNumber(rpctypes.BlockNumber(transaction.Height))
	if err != nil {
		b.Logger.Debug("block not found", "height", transaction.Height)
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

	nc, ok := b.ClientCtx.Client.(tmrpcclient.NetworkClient)
	if !ok {
		return nil, errors.New("invalid rpc client")
	}

	cp, err := nc.ConsensusParams(b.Ctx, &blk.Block.Height)
	if err != nil {
		return nil, err
	}

	traceTxRequest := evmtypes.QueryTraceTxRequest{
		Msg:             ethMsg,
		Predecessors:    predecessors,
		BlockNumber:     blk.Block.Height,
		BlockTime:       blk.Block.Time,
		BlockHash:       common.Bytes2Hex(blk.BlockID.Hash),
		ProposerAddress: sdk.ConsAddress(blk.Block.ProposerAddress),
		ChainId:         b.EvmChainID.Int64(),
		BlockMaxGas:     cp.ConsensusParams.Block.MaxGas,
	}

	if config != nil {
		traceTxRequest.TraceConfig, err = toEVMTraceConfig(config)
		if err != nil {
			return nil, err
		}
	}

	// minus one to get the context of block beginning
	contextHeight := transaction.Height - 1
	if contextHeight < 1 {
		// 0 is a special value in `ContextWithHeight`
		contextHeight = 1
	}
	traceResult, err := b.QueryClient.TraceTx(rpctypes.ContextWithHeight(contextHeight), &traceTxRequest)
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
	config *rpctypes.TraceConfig,
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
		b.Logger.Debug("block result not found", "height", block.Block.Height, "error", err.Error())
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

	traceConfig, err := toEVMTraceConfig(config)
	if err != nil {
		return nil, err
	}

	nc, ok := b.ClientCtx.Client.(tmrpcclient.NetworkClient)
	if !ok {
		return nil, errors.New("invalid rpc client")
	}

	cp, err := nc.ConsensusParams(b.Ctx, &block.Block.Height)
	if err != nil {
		return nil, err
	}

	traceBlockRequest := &evmtypes.QueryTraceBlockRequest{
		Txs:             msgs,
		TraceConfig:     traceConfig,
		BlockNumber:     block.Block.Height,
		BlockTime:       block.Block.Time,
		BlockHash:       common.Bytes2Hex(block.BlockID.Hash),
		ProposerAddress: sdk.ConsAddress(block.Block.ProposerAddress),
		ChainId:         b.EvmChainID.Int64(),
		BlockMaxGas:     cp.ConsensusParams.Block.MaxGas,
	}

	res, err := b.QueryClient.TraceBlock(ctxWithHeight, traceBlockRequest)
	if err != nil {
		return nil, err
	}

	decodedResults := make([]*evmtypes.TxTraceResult, txsLength)
	if err := json.Unmarshal(res.Data, &decodedResults); err != nil {
		return nil, err
	}

	return decodedResults, nil
}

// toEVMTraceConfig converts rpctypes.TraceConfig to evmtypes.TraceConfig
func toEVMTraceConfig(config *rpctypes.TraceConfig) (*evmtypes.TraceConfig, error) {
	if config == nil {
		return nil, nil
	}

	cfg := config.TraceConfig

	// if TracerConfig is an object, we need to encode it into JSON string
	if config.TracerConfig != nil && config.TracerJsonConfig == "" {
		switch v := config.TracerConfig.(type) {
		case string:
			// It's already a string, use it directly
			cfg.TracerJsonConfig = v
		case map[string]interface{}:
			// this is the compliant style
			// we need to encode it to a string before passing it to the ethermint side
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("unable to encode traceConfig to JSON: %w", err)
			}
			cfg.TracerJsonConfig = string(jsonBytes)
		default:
			return nil, fmt.Errorf("unexpected traceConfig type: %T", v)
		}
	}

	return &cfg, nil
}
