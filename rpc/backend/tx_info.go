package backend

import (
	"fmt"
	"math"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"

	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// GetTransactionByHash returns the Ethereum format transaction identified by Ethereum transaction hash
func (b *Backend) GetTransactionByHash(txHash common.Hash) (*rpctypes.RPCTransaction, error) {
	res, additional, err := b.GetTxByEthHash(txHash)
	hexTx := txHash.Hex()

	if err != nil {
		return b.GetTransactionByHashPending(txHash)
	}

	resBlock, err := b.TendermintBlockByNumber(rpctypes.BlockNumber(res.Height))
	if err != nil {
		b.Logger.Debug("block not found", "height", res.Height, "error", err.Error())
		return nil, err
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		b.Logger.Debug("block result not found", "height", resBlock.Block.Height, "error", err.Error())
		return nil, fmt.Errorf("block result not found: %w", err)
	}

	var ethMsg *evmtypes.MsgEthereumTx
	// if additional fields are empty we can try to get MsgEthereumTx from sdk.Msg array
	if additional == nil {
		// #nosec G115 always in range
		if int(res.TxIndex) >= len(resBlock.Block.Txs) {
			b.Logger.Error("tx out of bounds")
			return nil, fmt.Errorf("tx out of bounds")
		}
		tx, err := b.ClientCtx.TxConfig.TxDecoder()(resBlock.Block.Txs[res.TxIndex])
		if err != nil {
			b.Logger.Debug("decoding failed", "error", err.Error())
			return nil, fmt.Errorf("failed to decode tx: %w", err)
		}

		// the `res.MsgIndex` is inferred from tx index, should be within the bound.
		msg, ok := b.DecodeMsgEthereumTxFromCosmosTx(tx)
		if !ok {
			return nil, errors.New("invalid ethereum tx")
		}
		ethMsg = msg
	} else {
		// if additional fields are not empty try to parse synthetic tx from them
		ethMsg = b.parseSyntethicTxFromAdditionalFields(additional)
		if ethMsg == nil {
			b.Logger.Error("failed to get synthetic eth msg from additional fields")
			return nil, fmt.Errorf("failed to get synthetic eth msg from additional fields")
		}
	}

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs, _ := b.EthMsgsFromTendermintBlock(resBlock, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hexTx {
				if i > math.MaxInt32 {
					return nil, errors.New("tx index overflow")
				}
				res.EthTxIndex = int32(i) //#nosec G115 -- checked for int overflow already
				break
			}
		}
	}
	// if we still unable to find the eth tx index, return error, shouldn't happen.
	if res.EthTxIndex == -1 && additional == nil {
		return nil, errors.New("can't find index of ethereum tx")
	}
	if res.EthTxIndex == -1 {
		res.EthTxIndex = 0
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.Logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			blockRes.Height,
			"error",
			err,
		)
	}

	height := uint64(res.Height)    //#nosec G115 -- checked for int overflow already
	index := uint64(res.EthTxIndex) //#nosec G115 -- checked for int overflow already
	return rpctypes.NewTransactionFromMsg(
		ethMsg,
		common.BytesToHash(resBlock.BlockID.Hash.Bytes()),
		height,
		index,
		baseFee,
		b.EvmChainID,
		additional,
	)
}

// GetTransactionByHashPending find pending tx from mempool
func (b *Backend) GetTransactionByHashPending(txHash common.Hash) (*rpctypes.RPCTransaction, error) {
	hexTx := txHash.Hex()
	// try to find tx in mempool
	txs, err := b.PendingTransactions()
	if err != nil {
		b.Logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	for _, tx := range txs {
		msg, err := evmtypes.UnwrapEthereumMsg(tx, txHash)
		if err != nil {
			// not ethereum tx
			continue
		}

		if msg.Hash == hexTx {
			// use zero block values since it's not included in a block yet
			rpctx, err := rpctypes.NewTransactionFromMsg(
				msg,
				common.Hash{},
				uint64(0),
				uint64(0),
				nil,
				b.EvmChainID,
				nil,
			)
			if err != nil {
				return nil, err
			}
			return rpctx, nil
		}
	}

	b.Logger.Debug("tx not found", "hash", hexTx)
	return nil, nil
}

// GetGasUsed returns gasUsed from transaction
func (b *Backend) GetGasUsed(res *types.TxResult, price *big.Int, gas uint64) uint64 {
	// patch gasUsed if tx is reverted and happened before height on which fixed was introduced
	// to return real gas charged
	// more info at https://github.com/evmos/ethermint/pull/1557
	if res.Failed && res.Height < b.Cfg.JSONRPC.FixRevertGasRefundHeight {
		return new(big.Int).Mul(price, new(big.Int).SetUint64(gas)).Uint64()
	}
	return res.GasUsed
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (b *Backend) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	hexTx := hash.Hex()
	b.Logger.Debug("eth_getTransactionReceipt", "hash", hexTx)

	res, additional, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.Logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	resBlock, err := b.TendermintBlockByNumber(rpctypes.BlockNumber(res.Height))
	if err != nil {
		b.Logger.Debug("block not found", "height", res.Height, "error", err.Error())
		return nil, fmt.Errorf("block not found at height %d: %w", res.Height, err)
	}

	var txData evmtypes.TxData
	var ethMsg *evmtypes.MsgEthereumTx
	if additional == nil {
		// #nosec G115 always in range
		if int(res.TxIndex) >= len(resBlock.Block.Txs) {
			b.Logger.Error("tx out of bounds")
			return nil, fmt.Errorf("tx out of bounds")
		}

		tx, err := b.ClientCtx.TxConfig.TxDecoder()(resBlock.Block.Txs[res.TxIndex])
		if err != nil {
			b.Logger.Debug("decoding failed", "error", err.Error())
			return nil, fmt.Errorf("failed to decode tx: %w", err)
		}

		ethMsgDecoded, err := DecodeMsgEthereumTxFromCosmosMsg(tx.GetMsgs()[res.MsgIndex], b.ChainConfig().ChainID)
		if err != nil {
			b.Logger.Error("failed to get eth msg", "error", err.Error())
			return nil, err
		}

		ethMsg = ethMsgDecoded
		txData, err = evmtypes.UnpackTxData(ethMsg.Data)
		if err != nil {
			b.Logger.Error("failed to unpack tx data", "error", err.Error())
			return nil, err
		}
	} else {
		// if additional fields are not empty try to parse synthetic tx from them
		ethMsg = b.parseSyntethicTxFromAdditionalFields(additional)
		if ethMsg == nil {
			b.Logger.Error("failed to parse tx")
			return nil, fmt.Errorf("failed to parse tx")
		}
	}

	cumulativeGasUsed := uint64(0)
	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &res.Height)
	if err != nil {
		b.Logger.Debug("failed to retrieve block results", "height", res.Height, "error", err.Error())
		return nil, fmt.Errorf("block result not found at height %d: %w", res.Height, err)
	}

	for _, txResult := range blockRes.TxsResults[0:res.TxIndex] {
		cumulativeGasUsed += uint64(txResult.GasUsed) // #nosec G115 -- checked for int overflow already
	}

	cumulativeGasUsed += res.CumulativeGasUsed

	var status hexutil.Uint
	if res.Failed {
		status = hexutil.Uint(ethtypes.ReceiptStatusFailed)
	} else {
		status = hexutil.Uint(ethtypes.ReceiptStatusSuccessful)
	}

	var from common.Address
	if additional != nil || len(ethMsg.From) != 0 {
		from = common.BytesToAddress(ethMsg.From)
	} else if ethMsg.Data != nil {
		from, err = ethMsg.GetSenderLegacy(ethtypes.LatestSignerForChainID(b.EvmChainID))
		if err != nil {
			b.Logger.Debug("failed to parse from field", "hash", hexTx, "error", err.Error())
		}
	} else {
		return nil, errors.New("failed to parse receipt")
	}

	// parse tx logs from events
	msgIndex := int(res.MsgIndex) // #nosec G115 -- checked for int overflow already
	logs, err := TxLogsFromEvents(blockRes.TxsResults[res.TxIndex].Events, msgIndex)
	if err != nil {
		b.Logger.Debug("failed to parse logs", "hash", hexTx, "error", err.Error())
	}

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs, _ := b.EthMsgsFromTendermintBlock(resBlock, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hexTx {
				res.EthTxIndex = int32(i) // #nosec G115
				break
			}
		}
	}
	// return error if still unable to find the eth tx index
	if res.EthTxIndex == -1 {
		if additional != nil {
			res.EthTxIndex = 0
		} else {
			return nil, errors.New("can't find index of ethereum tx")
		}
	}

	to := &common.Address{}
	var txType uint8
	var effectiveGasPrice *hexutil.Big

	// Get baseFee for effective gas price calculation
	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		b.Logger.Debug("failed to get base fee", "height", res.Height, "error", err.Error())
		baseFee = nil
	}

	// Set transaction type and recipient (required for all blocks, including pruned)
	if txData == nil {
		// #nosec G115 always in range
		txType = uint8(additional.Type)
		*to = additional.Recipient
	} else {
		txType = ethMsg.AsTransaction().Type()
		to = txData.GetTo()
		effectiveGasPrice = (*hexutil.Big)(rpctypes.EffectiveGasPrice(ethMsg.AsTransaction(), baseFee))
	}

	// create the logs bloom
	var bin ethtypes.Bloom
	for _, log := range logs {
		bin.Add(log.Address.Bytes())
		for _, b := range log.Topics {
			bin.Add(b[:])
		}
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"logsBloom":         ethtypes.BytesToBloom(bin.Bytes()),
		"logs":              logs,

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": hash,
		"contractAddress": nil,
		// "gasUsed":         hexutil.Uint64(b.GetGasUsed(res, txData.GetGasPrice(), txData.GetGas())),
		"gasUsed": hexutil.Uint64(res.GasUsed),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        common.BytesToHash(resBlock.Block.Header.Hash()).Hex(),
		"blockNumber":      hexutil.Uint64(res.Height),     //#nosec G115 won't exceed uint64
		"transactionIndex": hexutil.Uint64(res.EthTxIndex), //#nosec G115 no int overflow expected here

		// sender and receiver (contract or EOA) addreses
		"from": from,
		"to":   to,
		"type": hexutil.Uint(txType),

		// https://github.com/foundry-rs/foundry/issues/7640
		"effectiveGasPrice": effectiveGasPrice,
	}

	if logs == nil {
		receipt["logs"] = [][]*ethtypes.Log{}
	}

	if txData != nil {
		// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
		if txData.GetTo() == nil {
			receipt["contractAddress"] = crypto.CreateAddress(from, txData.GetNonce())
		}
	}

	return receipt, nil
}

// GetTransactionLogs returns the transaction logs identified by hash.
func (b *Backend) GetTransactionLogs(hash common.Hash) ([]*ethtypes.Log, error) {
	hexTx := hash.Hex()

	// TODO https://github.com/zeta-chain/node/issues/4079
	// check if additional fields should be used here
	res, _, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.Logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	if res.Failed {
		// failed, return empty logs
		return nil, nil
	}

	resBlockResult, err := b.RPCClient.BlockResults(b.Ctx, &res.Height)
	if err != nil {
		b.Logger.Debug("block result not found", "number", res.Height, "error", err.Error())
		return nil, nil
	}

	// parse tx logs from events
	index := int(res.MsgIndex) // #nosec G701
	return TxLogsFromEvents(resBlockResult.TxsResults[res.TxIndex].Events, index)
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (b *Backend) GetTransactionByBlockHashAndIndex(
	hash common.Hash,
	idx hexutil.Uint,
) (*rpctypes.RPCTransaction, error) {
	b.Logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash.Hex(), "index", idx)
	sc, ok := b.ClientCtx.Client.(tmrpcclient.SignClient)
	if !ok {
		return nil, errors.New("invalid rpc client")
	}

	block, err := sc.BlockByHash(b.Ctx, hash.Bytes())
	if err != nil {
		b.Logger.Debug("block not found", "hash", hash.Hex(), "error", err.Error())
		return nil, nil
	}

	if block.Block == nil {
		b.Logger.Debug("block not found", "hash", hash.Hex())
		return nil, nil
	}

	return b.GetTransactionByBlockAndIndex(block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (b *Backend) GetTransactionByBlockNumberAndIndex(
	blockNum rpctypes.BlockNumber,
	idx hexutil.Uint,
) (*rpctypes.RPCTransaction, error) {
	b.Logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)

	block, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		b.Logger.Debug("block not found", "height", blockNum.Int64(), "error", err.Error())
		return nil, nil
	}

	if block.Block == nil {
		b.Logger.Debug("block not found", "height", blockNum.Int64())
		return nil, nil
	}

	return b.GetTransactionByBlockAndIndex(block, idx)
}

// GetTxByEthHash uses `/tx_query` to find transaction by ethereum tx hash
// TODO: Don't need to convert once hashing is fixed on Tendermint
// https://github.com/cometbft/cometbft/issues/6539
func (b *Backend) GetTxByEthHash(hash common.Hash) (*types.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	if b.Indexer != nil {
		txRes, err := b.Indexer.GetByTxHash(hash)
		if err != nil {
			return nil, nil, err
		}
		return txRes, nil, nil
	}

	// fallback to tendermint tx indexer
	query := fmt.Sprintf("%s.%s='%s'", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash, hash.Hex())
	txResult, txAdditional, err := b.QueryTendermintTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		return txs.GetTxByHash(hash)
	})
	if err != nil {
		return nil, nil, errorsmod.Wrapf(err, "GetTxByEthHash %s", hash.Hex())
	}
	return txResult, txAdditional, nil
}

// GetTxByTxIndex uses `/tx_query` to find transaction by tx index of valid ethereum txs
func (b *Backend) GetTxByTxIndex(
	height int64,
	index uint,
) (*types.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	int32Index := int32(index) //#nosec G115 -- checked for int overflow already
	if b.Indexer != nil {
		// #nosec G115 always in range
		txRes, err := b.Indexer.GetByBlockAndIndex(height, int32Index)
		if err == nil {
			return txRes, nil, nil
		}
	}

	// fallback to tendermint tx indexer
	query := fmt.Sprintf("tx.height=%d AND %s.%s=%d",
		height, evmtypes.TypeMsgEthereumTx,
		evmtypes.AttributeKeyTxIndex, index,
	)
	txResult, txAdditional, err := b.QueryTendermintTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		return txs.GetTxByTxIndex(int(index)) // #nosec G115 -- checked for int overflow already
	})
	if err != nil {
		return nil, nil, errorsmod.Wrapf(err, "GetTxByTxIndex %d %d", height, index)
	}
	return txResult, txAdditional, nil
}

// QueryTendermintTxIndexer query tx in tendermint tx indexer
func (b *Backend) QueryTendermintTxIndexer(
	query string,
	txGetter func(*rpctypes.ParsedTxs) *rpctypes.ParsedTx,
) (*types.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	resTxs, err := b.ClientCtx.Client.TxSearch(b.Ctx, query, false, nil, nil, "")
	if err != nil {
		return nil, nil, err
	}
	if len(resTxs.Txs) == 0 {
		return nil, nil, errors.New("ethereum tx not found")
	}
	txResult := resTxs.Txs[0]
	if !rpctypes.TxSucessOrExpectedFailure(&txResult.TxResult) {
		return nil, nil, errors.New("invalid ethereum tx")
	}

	var tx sdk.Tx
	if txResult.TxResult.Code != 0 {
		// it's only needed when the tx exceeds block gas limit
		tx, err = b.ClientCtx.TxConfig.TxDecoder()(txResult.Tx)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid ethereum tx")
		}
	}

	return rpctypes.ParseTxIndexerResult(txResult, tx, txGetter)
}

// GetTransactionByBlockAndIndex is the common code shared by `GetTransactionByBlockNumberAndIndex` and `GetTransactionByBlockHashAndIndex`.
func (b *Backend) GetTransactionByBlockAndIndex(
	block *tmrpctypes.ResultBlock,
	idx hexutil.Uint,
) (*rpctypes.RPCTransaction, error) {
	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &block.Block.Height)
	if err != nil {
		return nil, nil
	}

	// #nosec G115 always in range
	i := int(idx)
	ethMsgs, additionals := b.EthMsgsFromTendermintBlock(block, blockRes)
	if i >= len(ethMsgs) {
		b.Logger.Debug("block txs index out of bound", "index", i)
		return nil, nil
	}

	msg := ethMsgs[i]
	additional := additionals[i]
	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.Logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			block.Block.Height,
			"error",
			err,
		)
	}

	return rpctypes.NewTransactionFromMsg(
		msg,
		common.BytesToHash(block.Block.Hash()),
		// #nosec G115 always positive
		uint64(block.Block.Height),
		// #nosec G115 always positive
		uint64(idx),
		baseFee,
		b.EvmChainID,
		additional,
	)
}
