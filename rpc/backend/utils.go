package backend

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/rpc/types"
)

type txGasAndReward struct {
	gasUsed uint64
	reward  *big.Int
}

type sortGasAndReward []txGasAndReward

func (s sortGasAndReward) Len() int { return len(s) }
func (s sortGasAndReward) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortGasAndReward) Less(i, j int) bool {
	return s[i].reward.Cmp(s[j].reward) < 0
}

// getAccountNonce returns the account nonce for the given account address.
// If the pending value is true, it will iterate over the mempool (pending)
// txs in order to compute and return the pending tx sequence.
// Todo: include the ability to specify a blockNumber
func (b *Backend) getAccountNonce(
	accAddr common.Address,
	pending bool,
	height int64,
	logger log.Logger,
) (uint64, error) {
	queryClient := authtypes.NewQueryClient(b.ClientCtx)
	adr := sdk.AccAddress(accAddr.Bytes()).String()
	ctx := types.ContextWithHeight(height)
	res, err := queryClient.Account(ctx, &authtypes.QueryAccountRequest{Address: adr})
	if err != nil {
		st, ok := status.FromError(err)
		// treat as account doesn't exist yet
		if ok && st.Code() == codes.NotFound {
			return 0, nil
		}
		return 0, err
	}
	var acc sdk.AccountI
	if err := b.ClientCtx.InterfaceRegistry.UnpackAny(res.Account, &acc); err != nil {
		return 0, err
	}

	nonce := acc.GetSequence()

	if !pending {
		return nonce, nil
	}

	// the account retriever doesn't include the uncommitted transactions on the nonce so we need to
	// to manually add them.
	pendingTxs, err := b.PendingTransactions()
	if err != nil {
		logger.Error("failed to fetch pending transactions", "error", err.Error())
		return nonce, nil
	}

	// add the uncommitted txs to the nonce counter
	// only supports `MsgEthereumTx` style tx
	for _, tx := range pendingTxs {
		for _, msg := range (*tx).GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				// not ethereum tx
				break
			}

			sender, err := ethMsg.GetSenderLegacy(ethtypes.LatestSignerForChainID(b.EvmChainID))
			if err != nil {
				continue
			}
			if sender == accAddr {
				nonce++
			}
		}
	}

	return nonce, nil
}

func bigMax(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

// CalcBaseFee calculates the basefee of the header.
func CalcBaseFee(config *params.ChainConfig, parent *ethtypes.Header, p feemarkettypes.Params) (*big.Int, error) {
	// If the current block is the first EIP-1559 block, return the InitialBaseFee.
	if !config.IsLondon(parent.Number) {
		return new(big.Int).SetUint64(params.InitialBaseFee), nil
	}
	if p.ElasticityMultiplier == 0 {
		return nil, errors.New("ElasticityMultiplier cannot be 0 as it's checked in the params validation")
	}
	parentGasTarget := parent.GasLimit / uint64(p.ElasticityMultiplier)
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.GasUsed == parentGasTarget {
		return new(big.Int).Set(parent.BaseFee), nil
	}

	var (
		num   = new(big.Int)
		denom = new(big.Int)
	)

	if parent.GasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		// max(1, parentBaseFee * gasUsedDelta / parentGasTarget / baseFeeChangeDenominator)
		num.SetUint64(parent.GasUsed - parentGasTarget)
		num.Mul(num, parent.BaseFee)
		num.Div(num, denom.SetUint64(parentGasTarget))
		num.Div(num, denom.SetUint64(uint64(p.BaseFeeChangeDenominator)))
		baseFeeDelta := bigMax(num, common.Big1)

		return num.Add(parent.BaseFee, baseFeeDelta), nil
	}

	// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
	// max(0, parentBaseFee * gasUsedDelta / parentGasTarget / baseFeeChangeDenominator)
	num.SetUint64(parentGasTarget - parent.GasUsed)
	num.Mul(num, parent.BaseFee)
	num.Div(num, denom.SetUint64(parentGasTarget))
	num.Div(num, denom.SetUint64(uint64(p.BaseFeeChangeDenominator)))
	baseFee := num.Sub(parent.BaseFee, num)
	minGasPrice := p.MinGasPrice.TruncateInt().BigInt()
	return bigMax(baseFee, minGasPrice), nil
}

// ProcessBlock processes a Tendermint block and calculates fee history data for eth_feeHistory RPC.
// It extracts gas usage, base fees, and transaction reward percentiles from the block data.
//
// The function calculates:
//   - Current block's base fee and next block's base fee (for EIP-1559)
//   - Gas used ratio (gasUsed / gasLimit)
//   - Transaction reward percentiles based on effective gas tip values
//
// Parameters:
//   - tendermintBlock: The raw Tendermint block containing transaction data
//   - ethBlock: Ethereum-formatted block with gas limit and usage information
//   - rewardPercentiles: Percentile values (0-100) for reward calculation
//   - tendermintBlockResult: Block execution results containing gas usage per transaction
//   - targetOneFeeHistory: Output parameter to populate with calculated fee history data
//
// Returns an error if block processing fails due to invalid data types or calculation errors.
func (b *Backend) ProcessBlock(
	tendermintBlock *cmtrpctypes.ResultBlock,
	ethBlock *map[string]interface{},
	rewardPercentiles []float64,
	tendermintBlockResult *cmtrpctypes.ResultBlockResults,
	targetOneFeeHistory *types.OneFeeHistory,
) error {
	blockHeight := tendermintBlock.Block.Height
	blockBaseFee, err := b.BaseFee(tendermintBlockResult)
	if err != nil || blockBaseFee == nil {
		targetOneFeeHistory.BaseFee = big.NewInt(0)
	} else {
		targetOneFeeHistory.BaseFee = blockBaseFee
	}
	cfg := b.ChainConfig()
	// set gas used ratio
	gasLimitUint64, ok := (*ethBlock)["gasLimit"].(hexutil.Uint64)
	if !ok {
		return fmt.Errorf("invalid gas limit type: %T", (*ethBlock)["gasLimit"])
	}

	gasUsedBig, ok := (*ethBlock)["gasUsed"].(*hexutil.Big)
	if !ok {
		return fmt.Errorf("invalid gas used type: %T", (*ethBlock)["gasUsed"])
	}

	if cfg.IsLondon(big.NewInt(blockHeight + 1)) {
		var header ethtypes.Header
		header.Number = new(big.Int).SetInt64(blockHeight)
		baseFee, ok := (*ethBlock)["baseFeePerGas"].(*hexutil.Big)
		if !ok || baseFee == nil {
			header.BaseFee = big.NewInt(0)
		} else {
			header.BaseFee = baseFee.ToInt()
		}
		header.GasLimit = uint64(gasLimitUint64)
		header.GasUsed = gasUsedBig.ToInt().Uint64()
		ctx := types.ContextWithHeight(blockHeight)
		params, err := b.QueryClient.FeeMarket.Params(ctx, &feemarkettypes.QueryParamsRequest{})
		if err != nil {
			return err
		}
		nextBaseFee, err := CalcBaseFee(cfg, &header, params.Params)
		if err != nil {
			return err
		}
		targetOneFeeHistory.NextBaseFee = nextBaseFee
	} else {
		targetOneFeeHistory.NextBaseFee = new(big.Int)
	}
	gasusedfloat, _ := new(big.Float).SetInt(gasUsedBig.ToInt()).Float64()

	if gasLimitUint64 <= 0 {
		return fmt.Errorf(
			"gasLimit of block height %d should be bigger than 0 , current gaslimit %d",
			blockHeight,
			gasLimitUint64,
		)
	}

	gasUsedRatio := gasusedfloat / float64(gasLimitUint64)
	blockGasUsed := gasusedfloat
	targetOneFeeHistory.GasUsedRatio = gasUsedRatio

	rewardCount := len(rewardPercentiles)
	targetOneFeeHistory.Reward = make([]*big.Int, rewardCount)
	for i := 0; i < rewardCount; i++ {
		targetOneFeeHistory.Reward[i] = big.NewInt(0)
	}

	// check tendermintTxs
	tendermintTxs := tendermintBlock.Block.Txs
	tendermintTxResults := tendermintBlockResult.TxsResults
	tendermintTxCount := len(tendermintTxs)

	var sorter sortGasAndReward

	for i := 0; i < tendermintTxCount; i++ {
		eachTendermintTx := tendermintTxs[i]
		eachTendermintTxResult := tendermintTxResults[i]

		tx, err := b.ClientCtx.TxConfig.TxDecoder()(eachTendermintTx)
		if err != nil {
			b.Logger.Debug("failed to decode transaction in block", "height", blockHeight, "error", err.Error())
			continue
		}
		txGasUsed := uint64(eachTendermintTxResult.GasUsed) // #nosec G115
		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				continue
			}
			tx := ethMsg.AsTransaction()
			reward := tx.EffectiveGasTipValue(blockBaseFee)
			if reward == nil || reward.Sign() < 0 {
				b.Logger.Debug(
					"negative or nil reward found in transaction",
					"height",
					blockHeight,
					"txHash",
					tx.Hash().Hex(),
					"reward",
					reward,
				)
				reward = big.NewInt(0)
			}
			sorter = append(sorter, txGasAndReward{gasUsed: txGasUsed, reward: reward})
		}
	}

	// return an all zero row if there are no transactions to gather data from
	ethTxCount := len(sorter)
	if ethTxCount == 0 {
		return nil
	}

	sort.Sort(sorter)

	var txIndex int
	sumGasUsed := sorter[0].gasUsed

	for i, p := range rewardPercentiles {
		thresholdGasUsed := uint64(blockGasUsed * p / 100)
		for sumGasUsed < thresholdGasUsed && txIndex < ethTxCount-1 {
			txIndex++
			sumGasUsed += sorter[txIndex].gasUsed
		}
		targetOneFeeHistory.Reward[i] = sorter[txIndex].reward
	}

	return nil
}

// AllTxLogsFromEvents parses all ethereum logs from cosmos events
func AllTxLogsFromEvents(events []abci.Event) ([][]*ethtypes.Log, error) {
	allLogs := make([][]*ethtypes.Log, 0, 4)
	for _, event := range events {
		if event.Type != evmtypes.EventTypeTxLog {
			continue
		}

		logs, err := ParseTxLogsFromEvent(event)
		if err != nil {
			return nil, err
		}

		allLogs = append(allLogs, logs)
	}
	return allLogs, nil
}

// TxLogsFromEvents parses ethereum logs from cosmos events for specific msg index
func TxLogsFromEvents(events []abci.Event, msgIndex int) ([]*ethtypes.Log, error) {
	for _, event := range events {
		if event.Type != evmtypes.EventTypeTxLog {
			continue
		}

		if msgIndex > 0 {
			// not the eth tx we want
			msgIndex--
			continue
		}

		return ParseTxLogsFromEvent(event)
	}
	return nil, fmt.Errorf("eth tx logs not found for message index %d", msgIndex)
}

// ParseTxLogsFromEvent parse tx logs from one event
func ParseTxLogsFromEvent(event abci.Event) ([]*ethtypes.Log, error) {
	logs := make([]*evmtypes.Log, 0, len(event.Attributes))
	for _, attr := range event.Attributes {
		if attr.Key != evmtypes.AttributeKeyTxLog {
			continue
		}

		var txLog evmtypes.Log
		if err := json.Unmarshal([]byte(attr.Value), &txLog); err != nil {
			return nil, err
		}

		logs = append(logs, &txLog)
	}
	return evmtypes.LogsToEthereum(logs), nil
}

// ShouldIgnoreGasUsed returns true if the gasUsed in result should be ignored
// workaround for issue: https://github.com/cosmos/cosmos-sdk/issues/10832
func ShouldIgnoreGasUsed(res *abci.ExecTxResult) bool {
	return res.GetCode() == 11 && strings.Contains(res.GetLog(), "no block gas left to run tx: out of gas")
}

// GetLogsFromBlockResults returns the list of event logs from the tendermint block result response.
// If needed, it re-indexes logs before returning
func GetLogsFromBlockResults(blockRes *cmtrpctypes.ResultBlockResults) ([][]*ethtypes.Log, error) {
	blockLogs := [][]*ethtypes.Log{}
	for _, txResult := range blockRes.TxsResults {
		logs, err := AllTxLogsFromEvents(txResult.Events)
		if err != nil {
			return nil, err
		}
		blockLogs = append(blockLogs, logs...)
	}

	if needsReindexing(blockLogs) {
		reindexLogs(blockLogs)
	}

	return blockLogs, nil
}

// needsReindexing checks if logs in the block need re-indexing.Returns true if:
//  1. TxIndex values are not strictly increasing across transaction groups
//  2. Any log Index values are duplicated
//  3. Log Index values are not monotonically increasing across all groups(transactions)
func needsReindexing(blockLogs [][]*ethtypes.Log) bool {
	seenLogIndices := make(map[uint]bool)
	lastTxIndex := -1
	lastLogIndex := -1

	for _, txLogs := range blockLogs {
		if len(txLogs) == 0 {
			continue
		}

		groupTxIndex := int(txLogs[0].TxIndex)
		if groupTxIndex <= lastTxIndex {
			return true
		}
		lastTxIndex = groupTxIndex

		for _, entry := range txLogs {
			logIndex := int(entry.Index)

			// Check for duplicates
			if seenLogIndices[entry.Index] {
				return true
			}
			seenLogIndices[entry.Index] = true

			// Check monotonic ordering: each log's Index must be greater than the last across all transactions in a block
			if logIndex <= lastLogIndex {
				return true
			}
			lastLogIndex = logIndex
		}
	}
	return false
}

// reindexLogs assigns unique TxIndex and Index values to all logs in the block.
func reindexLogs(blockLogs [][]*ethtypes.Log) {
	globalTxIndex := uint(0)
	globalLogIndex := uint(0)

	for _, txLogs := range blockLogs {
		for _, entry := range txLogs {
			entry.TxIndex = globalTxIndex
			entry.Index = globalLogIndex
			globalLogIndex++
		}
		globalTxIndex++
	}
}

// GetHexProofs returns list of hex data of proof op
func GetHexProofs(proof *crypto.ProofOps) []string {
	if proof == nil {
		return []string{""}
	}
	proofs := []string{}
	// check for proof
	for _, p := range proof.Ops {
		proof := ""
		if len(p.Data) > 0 {
			proof = hexutil.Encode(p.Data)
		}
		proofs = append(proofs, proof)
	}
	return proofs
}
