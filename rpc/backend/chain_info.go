package backend

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	cmtrpcclient "github.com/cometbft/cometbft/rpc/client"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// ChainID is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (b *Backend) ChainID() (*hexutil.Big, error) {
	// if current block is at or past the EIP-155 replay-protection fork block, return EvmChainID from config
	bn, err := b.BlockNumber()
	if err != nil {
		b.Logger.Debug("failed to fetch latest block number", "error", err.Error())
		return (*hexutil.Big)(b.EvmChainID), nil
	}

	if config := b.ChainConfig(); config.IsEIP155(new(big.Int).SetUint64(uint64(bn))) {
		return (*hexutil.Big)(config.ChainID), nil
	}

	return nil, fmt.Errorf("chain not synced beyond EIP-155 replay-protection fork block")
}

// ChainConfig returns the latest ethereum chain configuration
func (b *Backend) ChainConfig() *params.ChainConfig {
	return evmtypes.GetEthChainConfig()
}

// GlobalMinGasPrice returns MinGasPrice param from FeeMarket
func (b *Backend) GlobalMinGasPrice() (*big.Int, error) {
	res, err := b.QueryClient.GlobalMinGasPrice(b.Ctx, &evmtypes.QueryGlobalMinGasPriceRequest{})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, fmt.Errorf("GlobalMinGasPrice query returned a nil response")
	}
	return res.MinGasPrice.BigInt(), nil
}

// BaseFee returns the base fee tracked by the Fee Market module.
// If the base fee is not enabled globally, the query returns nil.
// If the London hard fork is not activated at the current height, the query will
// return nil.
func (b *Backend) BaseFee(blockRes *cmtrpctypes.ResultBlockResults) (*big.Int, error) {
	// return BaseFee if London hard fork is activated and feemarket is enabled
	res, err := b.QueryClient.BaseFee(rpctypes.ContextWithHeight(blockRes.Height), &evmtypes.QueryBaseFeeRequest{})
	if err != nil || res.BaseFee == nil {
		// we can't tell if it's london HF not enabled or the state is pruned,
		// in either case, we'll fallback to parsing from begin blocker event,
		// faster to iterate reversely
		for i := len(blockRes.FinalizeBlockEvents) - 1; i >= 0; i-- {
			evt := blockRes.FinalizeBlockEvents[i]
			if evt.Type == evmtypes.EventTypeFeeMarket && len(evt.Attributes) > 0 {
				baseFee, ok := math.NewIntFromString(evt.Attributes[0].Value)
				if ok {
					return baseFee.BigInt(), nil
				}
				break
			}
		}
		return nil, err
	}

	if res.BaseFee == nil {
		return nil, nil
	}

	return res.BaseFee.BigInt(), nil
}

// CurrentHeader returns the latest block header
// This will return error as per node configuration
// if the ABCI responses are discarded ('discard_abci_responses' config param)
func (b *Backend) CurrentHeader() (*ethtypes.Header, error) {
	return b.HeaderByNumber(rpctypes.EthLatestBlockNumber)
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (b *Backend) PendingTransactions() ([]*sdk.Tx, error) {
	mc, ok := b.ClientCtx.Client.(cmtrpcclient.MempoolClient)
	if !ok {
		return nil, errors.New("invalid rpc client")
	}

	res, err := mc.UnconfirmedTxs(b.Ctx, nil)
	if err != nil {
		return nil, err
	}

	result := make([]*sdk.Tx, 0, len(res.Txs))
	for _, txBz := range res.Txs {
		tx, err := b.ClientCtx.TxConfig.TxDecoder()(txBz)
		if err != nil {
			return nil, err
		}
		result = append(result, &tx)
	}

	return result, nil
}

// GetCoinbase is the address that staking rewards will be send to (alias for Etherbase).
func (b *Backend) GetCoinbase() (sdk.AccAddress, error) {
	node, err := b.ClientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	status, err := node.Status(b.Ctx)
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryValidatorAccountRequest{
		ConsAddress: sdk.ConsAddress(status.ValidatorInfo.Address).String(),
	}

	res, err := b.QueryClient.ValidatorAccount(b.Ctx, req)
	if err != nil {
		return nil, err
	}

	address, _ := sdk.AccAddressFromBech32(res.AccountAddress) // #nosec G703
	return address, nil
}

// FeeHistory returns data relevant for fee estimation based on the specified range of blocks.
func (b *Backend) FeeHistory(
	userBlockCount, // number blocks to fetch, maximum is 100
	lastBlock rpc.BlockNumber, // the block to start search , to oldest
	rewardPercentiles []float64, // percentiles to fetch reward
) (*rpctypes.FeeHistoryResult, error) {
	blockEnd := int64(lastBlock) //#nosec G115 -- checked for int overflow already

	if blockEnd < 0 {
		blockNumber, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		blockEnd = int64(blockNumber) //#nosec G115 -- checked for int overflow already
	}

	blocks := int64(userBlockCount)                     // #nosec G115 -- checked for int overflow already
	maxBlockCount := int64(b.Cfg.JSONRPC.FeeHistoryCap) // #nosec G115 -- checked for int overflow already
	if blocks > maxBlockCount {
		return nil, fmt.Errorf("FeeHistory user block count %d higher than %d", blocks, maxBlockCount)
	}

	if blockEnd+1 < blocks {
		blocks = blockEnd + 1
	}
	// Ensure not trying to retrieve before genesis.
	blockStart := blockEnd + 1 - blocks
	oldestBlock := (*hexutil.Big)(big.NewInt(blockStart))

	// prepare space
	reward := make([][]*hexutil.Big, blocks)
	rewardCount := len(rewardPercentiles)
	for i := 0; i < int(blocks); i++ {
		reward[i] = make([]*hexutil.Big, rewardCount)
	}

	thisBaseFee := make([]*hexutil.Big, blocks+1)
	thisGasUsedRatio := make([]float64, blocks)

	// rewards should only be calculated if reward percentiles were included
	calculateRewards := rewardCount != 0

	// fetch block
	for blockID := blockStart; blockID <= blockEnd; blockID++ {
		index := int32(blockID - blockStart) // #nosec G115
		// tendermint block
		tendermintblock, err := b.TendermintBlockByNumber(rpctypes.BlockNumber(blockID))
		if tendermintblock == nil {
			return nil, err
		}

		// eth block
		ethBlock, err := b.GetBlockByNumber(rpctypes.BlockNumber(blockID), true)
		if ethBlock == nil {
			return nil, err
		}

		// tendermint block result
		tendermintBlockResult, err := b.RPCClient.BlockResults(b.Ctx, &tendermintblock.Block.Height)
		if tendermintBlockResult == nil {
			b.Logger.Debug("block result not found", "height", tendermintblock.Block.Height, "error", err.Error())
			return nil, err
		}

		oneFeeHistory := rpctypes.OneFeeHistory{}
		err = b.processBlock(tendermintblock, &ethBlock, rewardPercentiles, tendermintBlockResult, &oneFeeHistory)
		if err != nil {
			return nil, err
		}

		// copy
		thisBaseFee[index] = (*hexutil.Big)(oneFeeHistory.BaseFee)
		thisBaseFee[index+1] = (*hexutil.Big)(oneFeeHistory.NextBaseFee)
		thisGasUsedRatio[index] = oneFeeHistory.GasUsedRatio
		if calculateRewards {
			for j := 0; j < rewardCount; j++ {
				reward[index][j] = (*hexutil.Big)(oneFeeHistory.Reward[j])
				if reward[index][j] == nil {
					reward[index][j] = (*hexutil.Big)(big.NewInt(0))
				}
			}
		}
	}

	feeHistory := rpctypes.FeeHistoryResult{
		OldestBlock:  oldestBlock,
		BaseFee:      thisBaseFee,
		GasUsedRatio: thisGasUsedRatio,
	}

	if calculateRewards {
		feeHistory.Reward = reward
	}

	return &feeHistory, nil
}

// SuggestGasTipCap returns the suggested tip cap
// Although we don't support tx prioritization yet, but we return a positive value to help client to
// mitigate the base fee changes.
func (b *Backend) SuggestGasTipCap(baseFee *big.Int) (*big.Int, error) {
	if baseFee == nil {
		// london hardfork not enabled or feemarket not enabled
		return big.NewInt(0), nil
	}

	params, err := b.QueryClient.FeeMarket.Params(b.Ctx, &feemarkettypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	// calculate the maximum base fee delta in current block, assuming all block gas limit is consumed
	// ```
	// GasTarget = GasLimit / ElasticityMultiplier
	// Delta = BaseFee * (GasUsed - GasTarget) / GasTarget / Denominator
	// ```
	// The delta is at maximum when `GasUsed` is equal to `GasLimit`, which is:
	// ```
	// MaxDelta = BaseFee * (GasLimit - GasLimit / ElasticityMultiplier) / (GasLimit / ElasticityMultiplier) / Denominator
	//          = BaseFee * (ElasticityMultiplier - 1) / Denominator
	// ```t
	maxDelta := baseFee.Int64() * (int64(params.Params.ElasticityMultiplier) - 1) / int64(
		params.Params.BaseFeeChangeDenominator,
	) // #nosec G115
	if maxDelta < 0 {
		// impossible if the parameter validation passed.
		maxDelta = 0
	}
	return big.NewInt(maxDelta), nil
}
