package backend

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/proto/tendermint/crypto"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"

	"github.com/cosmos/evm/utils"
	feemarkettypes "github.com/cosmos/evm/x/feemarket/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/node/rpc/types"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
func (b *Backend) getAccountNonce(accAddr common.Address, pending bool, height int64, logger log.Logger) (uint64, error) {
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

	// eip2681 - tx with nonce >= 2^64 is invalid; saturate at 2^64-1
	// if already at max nonce, don't add to pending
	if nonce == math.MaxUint64 {
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
				// saturate - never overflow beyond 2^64-1 when counting pending txs
				if nonce < math.MaxUint64 {
					nonce++
				}
			}
		}
	}

	return nonce, nil
}

// ProcessBlock processes a CometBFT block and calculates fee history data for eth_feeHistory RPC.
// It extracts gas usage, base fees, and transaction reward percentiles from the block data.
//
// The function calculates:
//   - Current block's base fee and next block's base fee (for EIP-1559)
//   - Gas used ratio (gasUsed / gasLimit)
//   - Transaction reward percentiles based on effective gas tip values
//
// Parameters:
//   - cometBlock: The raw CometBFT block containing transaction data
//   - ethBlock: Ethereum-formatted block with gas limit and usage information
//   - rewardPercentiles: Percentile values (0-100) for reward calculation
//   - cometBlockResult: Block execution results containing gas usage per transaction
//   - targetOneFeeHistory: Output parameter to populate with calculated fee history data
//
// Returns an error if block processing fails due to invalid data types or calculation errors.
func (b *Backend) ProcessBlock(
	cometBlock *cmtrpctypes.ResultBlock,
	ethBlock *map[string]interface{},
	rewardPercentiles []float64,
	cometBlockResult *cmtrpctypes.ResultBlockResults,
	targetOneFeeHistory *types.OneFeeHistory,
) error {
	blockHeight := cometBlock.Block.Height
	blockBaseFee, err := b.BaseFee(cometBlockResult)
	if err != nil || blockBaseFee == nil {
		targetOneFeeHistory.BaseFee = big.NewInt(0)
	} else {
		targetOneFeeHistory.BaseFee = blockBaseFee
	}
	cfg := b.ChainConfig()
	gasLimitUint64, ok := (*ethBlock)["gasLimit"].(hexutil.Uint64)
	if !ok {
		return fmt.Errorf("invalid gas limit type: %T", (*ethBlock)["gasLimit"])
	}

	gasUsedBig, ok := (*ethBlock)["gasUsed"].(*hexutil.Big)
	if !ok {
		return fmt.Errorf("invalid gas used type: %T", (*ethBlock)["gasUsed"])
	}
	gasUsedInt := gasUsedBig.ToInt()

	timestampHex, ok := (*ethBlock)["timestamp"].(hexutil.Uint64)
	if !ok {
		return fmt.Errorf("invalid timestamp type: %T", (*ethBlock)["timestamp"])
	}

	header := ethtypes.Header{
		Number:   new(big.Int).SetInt64(blockHeight),
		GasLimit: uint64(gasLimitUint64),
		GasUsed:  gasUsedInt.Uint64(),
		Time:     uint64(timestampHex),
	}
	if baseFee, ok := (*ethBlock)["baseFeePerGas"].(*hexutil.Big); ok && baseFee != nil {
		header.BaseFee = baseFee.ToInt()
	} else {
		header.BaseFee = big.NewInt(0)
	}
	targetOneFeeHistory.BlobBaseFee = big.NewInt(0)
	targetOneFeeHistory.NextBlobBaseFee = big.NewInt(0)
	targetOneFeeHistory.BlobGasUsedRatio = 0

	if cfg.IsLondon(big.NewInt(blockHeight + 1)) {
		ctx := types.ContextWithHeight(blockHeight)
		params, err := b.QueryClient.FeeMarket.Params(ctx, &feemarkettypes.QueryParamsRequest{})
		if err != nil {
			return err
		}
		nextBaseFee, err := types.CalcBaseFee(cfg, &header, params.Params)
		if err != nil {
			return err
		}
		targetOneFeeHistory.NextBaseFee = nextBaseFee
	} else {
		targetOneFeeHistory.NextBaseFee = new(big.Int)
	}
	if cfg.IsCancun(header.Number, header.Time) {
		blobGasUsed := uint64(0)
		if val, ok := (*ethBlock)["blobGasUsed"].(hexutil.Uint64); ok {
			blobGasUsed = uint64(val)
		}
		excessBlobGas := uint64(0)
		if val, ok := (*ethBlock)["excessBlobGas"].(hexutil.Uint64); ok {
			excessBlobGas = uint64(val)
		}
		header.BlobGasUsed = new(uint64)
		*header.BlobGasUsed = blobGasUsed
		header.ExcessBlobGas = new(uint64)
		*header.ExcessBlobGas = excessBlobGas

		targetOneFeeHistory.BlobBaseFee = eip4844.CalcBlobFee(cfg, &header)
		nextExcess := eip4844.CalcExcessBlobGas(cfg, &header, header.Time)
		nextHeader := &ethtypes.Header{
			Number:        header.Number,
			Time:          header.Time,
			ExcessBlobGas: &nextExcess,
		}
		targetOneFeeHistory.NextBlobBaseFee = eip4844.CalcBlobFee(cfg, nextHeader)

		maxBlobGas := eip4844.MaxBlobGasPerBlock(cfg, header.Time)
		targetOneFeeHistory.BlobGasUsedRatio = safeRatio(blobGasUsed, maxBlobGas)
	}
	if gasLimitUint64 <= 0 {
		return fmt.Errorf("gasLimit of block height %d should be bigger than 0 , current gaslimit %d", blockHeight, gasLimitUint64)
	}

	gasUsedUint64 := gasUsedInt.Uint64()
	targetOneFeeHistory.GasUsedRatio = safeRatio(gasUsedUint64, uint64(gasLimitUint64))
	blockGasUsed := float64(gasUsedUint64)

	rewardCount := len(rewardPercentiles)
	targetOneFeeHistory.Reward = make([]*big.Int, rewardCount)
	for i := 0; i < rewardCount; i++ {
		targetOneFeeHistory.Reward[i] = big.NewInt(0)
	}

	// check cometTxs
	cometTxs := cometBlock.Block.Txs
	cometTxResults := cometBlockResult.TxsResults
	CometTxCount := len(cometTxs)

	var sorter sortGasAndReward

	for i := 0; i < CometTxCount; i++ {
		cometTx := cometTxs[i]
		cometTxResult := cometTxResults[i]

		tx, err := b.ClientCtx.TxConfig.TxDecoder()(cometTx)
		if err != nil {
			b.Logger.Debug("failed to decode transaction in block", "height", blockHeight, "error", err.Error())
			continue
		}
		txGasUsed := uint64(cometTxResult.GasUsed) // #nosec G115
		for _, msg := range tx.GetMsgs() {
			ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
			if !ok {
				continue
			}
			tx := ethMsg.AsTransaction()
			reward, err := tx.EffectiveGasTip(blockBaseFee)
			if err != nil {
				b.Logger.Error("failed to calculate effective gas tip", "height", blockHeight, "error", err.Error())
			}
			if reward == nil || reward.Sign() < 0 {
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

func safeRatio(num, denom uint64) float64 {
	if denom == 0 || num == 0 {
		return 0
	}
	rat := new(big.Rat).SetFrac(
		new(big.Int).SetUint64(num),
		new(big.Int).SetUint64(denom),
	)
	value, _ := rat.Float64()
	return value
}

// ShouldIgnoreGasUsed returns true if the gasUsed in result should be ignored
// workaround for issue: https://github.com/cosmos/cosmos-sdk/issues/10832
func ShouldIgnoreGasUsed(res *abci.ExecTxResult) bool {
	return res.GetCode() == 11 && strings.Contains(res.GetLog(), "no block gas left to run tx: out of gas")
}

// GetLogsFromBlockResults returns the list of event logs from the CometBFT block result response
func GetLogsFromBlockResults(blockRes *cmtrpctypes.ResultBlockResults) ([][]*ethtypes.Log, error) {
	height, err := utils.SafeUint64(blockRes.Height)
	if err != nil {
		return nil, err
	}
	blockLogs := [][]*ethtypes.Log{}
	for _, txResult := range blockRes.TxsResults {
		logs, err := evmtypes.DecodeTxLogs(txResult.Data, height)
		if err != nil {
			return nil, err
		}
		blockLogs = append(blockLogs, logs)
	}
	return blockLogs, nil
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
