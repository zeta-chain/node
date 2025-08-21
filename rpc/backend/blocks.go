package backend

import (
	"fmt"
	"math/big"
	"strconv"

	abci "github.com/cometbft/cometbft/abci/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	cmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// BlockNumber returns the current block number in abci app state. Because abci
// app state could lag behind from tendermint latest block, it's more stable for
// the client to use the latest block number in abci app state than tendermint
// rpc.
func (b *Backend) BlockNumber() (hexutil.Uint64, error) {
	// do any grpc query, ignore the response and use the returned block height
	var header metadata.MD
	_, err := b.QueryClient.Params(b.Ctx, &evmtypes.QueryParamsRequest{}, grpc.Header(&header))
	if err != nil {
		return hexutil.Uint64(0), err
	}

	blockHeightHeader := header.Get(grpctypes.GRPCBlockHeightHeader)
	if headerLen := len(blockHeightHeader); headerLen != 1 {
		return 0, fmt.Errorf(
			"unexpected '%s' gRPC header length; got %d, expected: %d",
			grpctypes.GRPCBlockHeightHeader,
			headerLen,
			1,
		)
	}

	height, err := strconv.ParseUint(blockHeightHeader[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	return hexutil.Uint64(height), nil
}

// GetBlockByNumber returns the JSON-RPC compatible Ethereum block identified by
// block number. Depending on fullTx it either returns the full transaction
// objects or if false only the hashes of the transactions.
func (b *Backend) GetBlockByNumber(blockNum rpctypes.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, nil
	}

	// return if requested block height is greater than the current one
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		b.Logger.Debug("failed to fetch block result from Tendermint", "height", blockNum, "error", err.Error())
		return nil, nil
	}

	res, err := b.RPCBlockFromTendermintBlock(resBlock, blockRes, fullTx)
	if err != nil {
		b.Logger.Debug("GetEthBlockFromTendermint failed", "height", blockNum, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// GetBlockByHash returns the JSON-RPC compatible Ethereum block identified by
// hash.
func (b *Backend) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := b.TendermintBlockByHash(hash)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		// block not found
		return nil, nil
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		b.Logger.Debug(
			"failed to fetch block result from Tendermint",
			"block-hash",
			hash.String(),
			"error",
			err.Error(),
		)
		return nil, nil
	}

	res, err := b.RPCBlockFromTendermintBlock(resBlock, blockRes, fullTx)
	if err != nil {
		b.Logger.Debug("GetEthBlockFromTendermint failed", "hash", hash, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// GetBlockTransactionCountByHash returns the number of Ethereum transactions in
// the block identified by hash.
func (b *Backend) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	block, err := b.RPCClient.BlockByHash(b.Ctx, hash.Bytes())
	if err != nil {
		b.Logger.Debug("block not found", "hash", hash.Hex(), "error", err.Error())
		return nil
	}

	if block.Block == nil {
		b.Logger.Debug("block not found", "hash", hash.Hex())
		return nil
	}

	return b.GetBlockTransactionCount(block)
}

// GetBlockTransactionCountByNumber returns the number of Ethereum transactions
// in the block identified by number.
func (b *Backend) GetBlockTransactionCountByNumber(blockNum rpctypes.BlockNumber) *hexutil.Uint {
	block, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		b.Logger.Debug("block not found", "height", blockNum.Int64(), "error", err.Error())
		return nil
	}

	if block.Block == nil {
		b.Logger.Debug("block not found", "height", blockNum.Int64())
		return nil
	}

	return b.GetBlockTransactionCount(block)
}

// GetBlockTransactionCount returns the number of Ethereum transactions in a
// given block.
func (b *Backend) GetBlockTransactionCount(block *tmrpctypes.ResultBlock) *hexutil.Uint {
	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &block.Block.Height)
	if err != nil {
		return nil
	}

	ethMsgs, _ := b.EthMsgsFromTendermintBlock(block, blockRes)
	n := hexutil.Uint(len(ethMsgs))
	return &n
}

// TendermintBlockByNumber returns a Tendermint-formatted block for a given
// block number
func (b *Backend) TendermintBlockByNumber(blockNum rpctypes.BlockNumber) (*tmrpctypes.ResultBlock, error) {
	height := blockNum.Int64()
	if height <= 0 {
		// fetch the latest block number from the app state, more accurate than the tendermint block store state.
		n, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		height = int64(n) //#nosec G115 -- checked for int overflow already
	}
	resBlock, err := b.RPCClient.Block(b.Ctx, &height)
	if err != nil {
		b.Logger.Debug("tendermint client failed to get block", "height", height, "error", err.Error())
		return nil, err
	}

	if resBlock.Block == nil {
		b.Logger.Debug("TendermintBlockByNumber block not found", "height", height)
		return nil, nil
	}

	return resBlock, nil
}

// TendermintBlockResultByNumber returns a Tendermint-formatted block result
// by block number
func (b *Backend) TendermintBlockResultByNumber(height *int64) (*tmrpctypes.ResultBlockResults, error) {
	res, err := b.RPCClient.BlockResults(b.Ctx, height)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block result from Tendermint %d: %w", *height, err)
	}

	return res, nil
}

// TendermintBlockByHash returns a Tendermint-formatted block by block number
func (b *Backend) TendermintBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error) {
	resBlock, err := b.RPCClient.BlockByHash(b.Ctx, blockHash.Bytes())
	if err != nil {
		b.Logger.Debug("tendermint client failed to get block", "blockHash", blockHash.Hex(), "error", err.Error())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		b.Logger.Debug("TendermintBlockByHash block not found", "blockHash", blockHash.Hex())
		return nil, fmt.Errorf("block not found for hash %s", blockHash.Hex())
	}

	return resBlock, nil
}

// BlockNumberFromTendermint returns the BlockNumber from BlockNumberOrHash
func (b *Backend) BlockNumberFromTendermint(blockNrOrHash rpctypes.BlockNumberOrHash) (rpctypes.BlockNumber, error) {
	switch {
	case blockNrOrHash.BlockHash == nil && blockNrOrHash.BlockNumber == nil:
		return rpctypes.EthEarliestBlockNumber, fmt.Errorf("types BlockHash and BlockNumber cannot be both nil")
	case blockNrOrHash.BlockHash != nil:
		blockNumber, err := b.BlockNumberFromTendermintByHash(*blockNrOrHash.BlockHash)
		if err != nil {
			return rpctypes.EthEarliestBlockNumber, err
		}
		return rpctypes.NewBlockNumber(blockNumber), nil
	case blockNrOrHash.BlockNumber != nil:
		return *blockNrOrHash.BlockNumber, nil
	default:
		return rpctypes.EthEarliestBlockNumber, nil
	}
}

// BlockNumberFromTendermintByHash returns the block height of given block hash
func (b *Backend) BlockNumberFromTendermintByHash(blockHash common.Hash) (*big.Int, error) {
	resBlock, err := b.RPCClient.HeaderByHash(b.Ctx, blockHash.Bytes())
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		return nil, errors.Errorf("block not found for hash %s", blockHash.Hex())
	}

	return big.NewInt(resBlock.Header.Height), nil
}

// EthMsgsFromTendermintBlock returns all real and synthetic MsgEthereumTxs from a
// Tendermint block. It also ensures consistency over the correct txs indexes
// across RPC endpoints
func (b *Backend) EthMsgsFromTendermintBlock(
	resBlock *tmrpctypes.ResultBlock,
	blockRes *tmrpctypes.ResultBlockResults,
) ([]*evmtypes.MsgEthereumTx, []*rpctypes.TxResultAdditionalFields) {
	var ethMsgs []*evmtypes.MsgEthereumTx
	var txsAdditional []*rpctypes.TxResultAdditionalFields
	block := resBlock.Block

	txResults := blockRes.TxsResults
	for i, tx := range block.Txs {
		// Check if tx exists on EVM by cross checking with blockResults:
		//  - Include unsuccessful tx that exceeds block gas limit
		//  - Include unsuccessful tx that failed when committing changes to stateDB
		//  - Exclude unsuccessful tx with any other error but ExceedBlockGasLimit
		if !rpctypes.TxSucessOrExpectedFailure(txResults[i]) {
			b.Logger.Debug("invalid tx result code", "cosmos-hash", hexutil.Encode(tx.Hash()))
			continue
		}

		tx, err := b.ClientCtx.TxConfig.TxDecoder()(tx)
		// assumption is that if regular ethermint msg is found in tx
		// there should not be synthetic one as well
		shouldCheckForSyntheticTx := true
		// if tx can be decoded, try to find MsgEthereumTx inside
		if err == nil {
			for _, msg := range tx.GetMsgs() {
				ethMsg, ok := msg.(*evmtypes.MsgEthereumTx)
				if ok {
					shouldCheckForSyntheticTx = false
					ethMsg.Hash = ethMsg.AsTransaction().Hash().Hex()
					ethMsgs = append(ethMsgs, ethMsg)
					txsAdditional = append(txsAdditional, nil)
				}
			}
		} else {
			b.Logger.Debug("failed to decode transaction in block", "height", block.Height, "error", err.Error())
		}

		// if tx can not be decoded or MsgEthereumTx was not found, try to parse it from block results
		if shouldCheckForSyntheticTx {
			ethMsg, additional := b.parseSyntheticTxFromBlockResults(txResults, i, tx, block)
			if ethMsg != nil {
				ethMsgs = append(ethMsgs, ethMsg)
				txsAdditional = append(txsAdditional, additional)
			}
		}
	}
	return ethMsgs, txsAdditional
}

func (b *Backend) parseSyntheticTxFromBlockResults(
	txResults []*abci.ExecTxResult,
	i int,
	tx sdk.Tx,
	block *cmtypes.Block,
) (*evmtypes.MsgEthereumTx, *rpctypes.TxResultAdditionalFields) {
	res, additional, err := rpctypes.ParseTxBlockResult(txResults[i], tx, i, block.Height)
	// just skip tx if it can not be parsed, so remaining txs from the block are parsed
	if err != nil {
		b.Logger.Error(err.Error())
		return nil, nil
	}
	if additional == nil || res == nil {
		b.Logger.Debug("synthetic ethereum tx not found in msgs: block %d, index %d", block.Height, i)
		return nil, nil
	}
	return b.parseSyntethicTxFromAdditionalFields(additional), additional
}

func (b *Backend) parseSyntethicTxFromAdditionalFields(
	additional *rpctypes.TxResultAdditionalFields,
) *evmtypes.MsgEthereumTx {
	recipient := additional.Recipient
	// for transactions before v31 this value was mistakenly used for Gas field
	gas := additional.GasUsed
	if additional.GasLimit != nil {
		gas = *additional.GasLimit
	}
	t := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    additional.Nonce,
		Data:     additional.Data,
		Gas:      gas,
		To:       &recipient,
		GasPrice: nil,
		Value:    additional.Value,
		V:        big.NewInt(0),
		R:        big.NewInt(0),
		S:        big.NewInt(0),
	})
	ethMsg := &evmtypes.MsgEthereumTx{}
	err := ethMsg.FromEthereumTx(t)
	if err != nil {
		b.Logger.Error("can not create eth msg", err.Error())
		return nil
	}
	ethMsg.Hash = additional.Hash.Hex()
	ethMsg.From = additional.Sender.Bytes()
	return ethMsg
}

// HeaderByNumber returns the block header identified by height.
func (b *Backend) HeaderByNumber(blockNum rpctypes.BlockNumber) (*ethtypes.Header, error) {
	resBlock, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		return nil, errors.Errorf("block not found for height %d", blockNum)
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		return nil, errors.Errorf("block result not found for height %d", resBlock.Block.Height)
	}

	bloom, err := b.BlockBloom(blockRes)
	if err != nil {
		b.Logger.Debug("HeaderByNumber BlockBloom failed", "height", resBlock.Block.Height)
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.Logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			resBlock.Block.Height,
			"error",
			err,
		)
	}

	ethHeader := rpctypes.EthHeaderFromTendermint(resBlock.Block.Header, bloom, baseFee)
	return ethHeader, nil
}

// HeaderByHash returns the block header identified by hash.
func (b *Backend) HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error) {
	resHeader, err := b.RPCClient.HeaderByHash(b.Ctx, blockHash.Bytes())
	if err != nil {
		return nil, err
	}

	if resHeader == nil {
		return nil, errors.Errorf("header not found for hash %s", blockHash.Hex())
	}

	height := resHeader.Header.Height

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resHeader.Header.Height)
	if err != nil {
		return nil, errors.Errorf("block result not found for height %d", height)
	}

	bloom, err := b.BlockBloom(blockRes)
	if err != nil {
		b.Logger.Debug("HeaderByHash BlockBloom failed", "height", height)
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.Logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			height,
			"error",
			err,
		)
	}

	ethHeader := rpctypes.EthHeaderFromTendermint(*resHeader.Header, bloom, baseFee)
	return ethHeader, nil
}

// BlockBloom query block bloom filter from block results
func (b *Backend) BlockBloom(blockRes *tmrpctypes.ResultBlockResults) (ethtypes.Bloom, error) {
	for _, event := range blockRes.FinalizeBlockEvents {
		if event.Type != evmtypes.EventTypeBlockBloom {
			continue
		}

		for _, attr := range event.Attributes {
			if attr.Key == evmtypes.AttributeKeyEthereumBloom {
				return ethtypes.BytesToBloom([]byte(attr.Value)), nil
			}
		}
	}
	return ethtypes.Bloom{}, errors.New("block bloom event is not found")
}

// RPCBlockFromTendermintBlock returns a JSON-RPC compatible Ethereum block from a
// given Tendermint block and its block result.
func (b *Backend) RPCBlockFromTendermintBlock(
	resBlock *tmrpctypes.ResultBlock,
	blockRes *tmrpctypes.ResultBlockResults,
	fullTx bool,
) (map[string]interface{}, error) {
	ethRPCTxs := []interface{}{}
	block := resBlock.Block

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.Logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			block.Height,
			"error",
			err,
		)
	}

	msgs, txsAdditional := b.EthMsgsFromTendermintBlock(resBlock, blockRes)
	for txIndex, ethMsg := range msgs {
		if !fullTx {
			hash := common.HexToHash(ethMsg.Hash)
			ethRPCTxs = append(ethRPCTxs, hash)
			continue
		}

		var rpcTx *rpctypes.RPCTransaction
		if txsAdditional[txIndex] == nil {
			height := uint64(block.Height) //#nosec G115 -- checked for int overflow already
			index := uint64(txIndex)       //#nosec G115 -- checked for int overflow already
			rpcTx, err = rpctypes.NewRPCTransaction(
				ethMsg,
				common.BytesToHash(block.Hash()),
				height,
				index,
				baseFee,
				b.EvmChainID,
			)
		} else {
			// #nosec G115 non negative value
			rpcTx, err = rpctypes.NewRPCTransactionFromIncompleteMsg(ethMsg, common.BytesToHash(block.Hash()), uint64(block.Height), uint64(txIndex), baseFee, b.EvmChainID, txsAdditional[txIndex])
		}
		if err != nil {
			b.Logger.Debug("NewTransactionFromData for receipt failed", "hash", ethMsg.Hash, "error", err.Error())
			continue
		}
		ethRPCTxs = append(ethRPCTxs, rpcTx)
	}

	bloom, err := b.BlockBloom(blockRes)
	if err != nil {
		b.Logger.Debug("failed to query BlockBloom", "height", block.Height, "error", err.Error())
	}

	req := &evmtypes.QueryValidatorAccountRequest{
		ConsAddress: sdk.ConsAddress(block.Header.ProposerAddress).String(),
	}

	var validatorAccAddr sdk.AccAddress

	ctx := rpctypes.ContextWithHeight(block.Height)
	res, err := b.QueryClient.ValidatorAccount(ctx, req)
	if err != nil {
		b.Logger.Debug(
			"failed to query validator operator address",
			"height", block.Height,
			"cons-address", req.ConsAddress,
			"error", err.Error(),
		)
		// use zero address as the validator operator address
		validatorAccAddr = sdk.AccAddress(common.Address{}.Bytes())
	} else {
		validatorAccAddr, err = sdk.AccAddressFromBech32(res.AccountAddress)
		if err != nil {
			return nil, err
		}
	}

	validatorAddr := common.BytesToAddress(validatorAccAddr)

	gasLimit, err := rpctypes.BlockMaxGasFromConsensusParams(ctx, b.ClientCtx, block.Height)
	if err != nil {
		b.Logger.Error("failed to query consensus params", "error", err.Error())
	}

	gasUsed := uint64(0)

	for _, txsResult := range blockRes.TxsResults {
		// workaround for cosmos-sdk bug. https://github.com/cosmos/cosmos-sdk/issues/10832
		if ShouldIgnoreGasUsed(txsResult) {
			// block gas limit has exceeded, other txs must have failed with same reason.
			break
		}
		gasUsed += uint64(txsResult.GetGasUsed()) // #nosec G115 -- checked for int overflow already
	}

	formattedBlock := rpctypes.FormatBlock(
		block.Header, block.Size(),
		gasLimit, new(big.Int).SetUint64(gasUsed),
		ethRPCTxs, bloom, validatorAddr, baseFee,
	)
	return formattedBlock, nil
}

// EthBlockByNumber returns the Ethereum Block identified by number.
func (b *Backend) EthBlockByNumber(blockNum rpctypes.BlockNumber) (*ethtypes.Block, error) {
	resBlock, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		// block not found
		return nil, fmt.Errorf("block not found for height %d", blockNum)
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("block result not found for height %d", resBlock.Block.Height)
	}

	return b.EthBlockFromTendermintBlock(resBlock, blockRes)
}

// EthBlockFromTendermintBlock returns an Ethereum Block type from Tendermint block
// EthBlockFromTendermintBlock
func (b *Backend) EthBlockFromTendermintBlock(
	resBlock *tmrpctypes.ResultBlock,
	blockRes *tmrpctypes.ResultBlockResults,
) (*ethtypes.Block, error) {
	block := resBlock.Block
	height := block.Height
	bloom, err := b.BlockBloom(blockRes)
	if err != nil {
		b.Logger.Debug("HeaderByNumber BlockBloom failed", "height", height)
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle error for pruned node and log
		b.Logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			height,
			"error",
			err,
		)
	}

	ethHeader := rpctypes.EthHeaderFromTendermint(block.Header, bloom, baseFee)
	msgs, additionals := b.EthMsgsFromTendermintBlock(resBlock, blockRes)

	txs := []*ethtypes.Transaction{}
	for i, ethMsg := range msgs {
		if additionals[i] == nil {
			txs = append(txs, ethMsg.AsTransaction())
		}
	}

	// TODO: add tx receipts
	ethBlock := ethtypes.NewBlock(
		ethHeader,
		&ethtypes.Body{Transactions: txs, Uncles: nil, Withdrawals: nil},
		nil,
		trie.NewStackTrie(nil))
	return ethBlock, nil
}

// TODO https://github.com/zeta-chain/node/issues/4079
// new method, needs refactoring with synthetic txs
// GetBlockReceipts returns the receipts for a given block number or hash.
// func (b *Backend) GetBlockReceipts(
// 	blockNrOrHash rpctypes.BlockNumberOrHash,
// ) ([]map[string]interface{}, error) {
// 	blockNum, err := b.BlockNumberFromTendermint(blockNrOrHash)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get block number from hash: %w", err)
// 	}

// 	resBlock, err := b.TendermintBlockByNumber(blockNum)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get block by number: %w", err)
// 	}

// 	if resBlock == nil {
// 		return nil, fmt.Errorf("block not found for height %d", *blockNum.TmHeight())
// 	}

// 	blockRes, err := b.RPCClient.BlockResults(b.Ctx, blockNum.TmHeight())
// 	if err != nil {
// 		return nil, fmt.Errorf("block result not found for height %d", resBlock.Block.Height)
// 	}

// 	msgs := b.EthMsgsFromTendermintBlock(resBlock, blockRes)
// 	result := make([]map[string]interface{}, len(msgs))
// 	for i, msg := range msgs {
// 		result[i], err = b.formatTxReceipt(
// 			msg,
// 			msgs,
// 			blockRes,
// 			common.BytesToHash(resBlock.Block.Header.Hash()).Hex(),
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get transaction receipt for tx %s: %w", msg.Hash, err)
// 		}
// 	}

// 	return result, nil
// }

// func (b *Backend) formatTxReceipt(ethMsg *evmtypes.MsgEthereumTx, blockMsgs []*evmtypes.MsgEthereumTx, blockRes *tmrpctypes.ResultBlockResults, blockHeaderHash string) (map[string]interface{}, error) {
// 	txResult, err := b.GetTxByEthHash(common.HexToHash(ethMsg.Hash))
// 	if err != nil {
// 		return nil, fmt.Errorf("tx not found: hash=%s, error=%s", ethMsg.Hash, err.Error())
// 	}

// 	txData, err := evmtypes.UnpackTxData(ethMsg.Data)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unpack tx data: %w", err)
// 	}

// 	cumulativeGasUsed := uint64(0)

// 	for _, txResult := range blockRes.TxsResults[0:txResult.TxIndex] {
// 		cumulativeGasUsed += uint64(txResult.GasUsed) // #nosec G115 -- checked for int overflow already
// 	}

// 	cumulativeGasUsed += txResult.CumulativeGasUsed

// 	var status hexutil.Uint
// 	if txResult.Failed {
// 		status = hexutil.Uint(ethtypes.ReceiptStatusFailed)
// 	} else {
// 		status = hexutil.Uint(ethtypes.ReceiptStatusSuccessful)
// 	}

// 	chainID, err := b.ChainID()
// 	if err != nil {
// 		return nil, err
// 	}

// 	from, err := ethMsg.GetSenderLegacy(ethtypes.LatestSignerForChainID(chainID.ToInt()))
// 	if err != nil {
// 		return nil, err
// 	}

// 	// parse tx logs from events
// 	msgIndex := int(txResult.MsgIndex) // #nosec G115 -- checked for int overflow already
// 	logs, err := TxLogsFromEvents(blockRes.TxsResults[txResult.TxIndex].Events, msgIndex)
// 	if err != nil {
// 		b.Logger.Debug("failed to parse logs", "hash", ethMsg.Hash, "error", err.Error())
// 	}

// 	if txResult.EthTxIndex == -1 {
// 		// Fallback to find tx index by iterating all valid eth transactions
// 		for i := range blockMsgs {
// 			if blockMsgs[i].Hash == ethMsg.Hash {
// 				txResult.EthTxIndex = int32(i) // #nosec G115
// 				break
// 			}
// 		}
// 	}
// 	// return error if still unable to find the eth tx index
// 	if txResult.EthTxIndex == -1 {
// 		return nil, fmt.Errorf("can't find index of ethereum tx")
// 	}

// 	receipt := map[string]interface{}{
// 		// Consensus fields: These fields are defined by the Yellow Paper
// 		"status":            status,
// 		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
// 		"logsBloom":         ethtypes.CreateBloom(&ethtypes.Receipt{Logs: logs}),
// 		"logs":              logs,

// 		// Implementation fields: These fields are added by geth when processing a transaction.
// 		// They are stored in the chain database.
// 		"transactionHash": common.HexToHash(ethMsg.Hash),
// 		"contractAddress": nil,
// 		"gasUsed":         hexutil.Uint64(b.GetGasUsed(txResult, txData.GetGasPrice(), txData.GetGas())),

// 		// Inclusion information: These fields provide information about the inclusion of the
// 		// transaction corresponding to this receipt.
// 		"blockHash":        blockHeaderHash,
// 		"blockNumber":      hexutil.Uint64(txResult.Height),     //nolint:gosec // G115 // won't exceed uint64
// 		"transactionIndex": hexutil.Uint64(txResult.EthTxIndex), //nolint:gosec // G115 // no int overflow expected here

// 		// https://github.com/foundry-rs/foundry/issues/7640
// 		"effectiveGasPrice": (*hexutil.Big)(txData.GetGasPrice()),

// 		// sender and receiver (contract or EOA) addreses
// 		"from": from,
// 		"to":   txData.GetTo(),
// 		"type": hexutil.Uint(ethMsg.AsTransaction().Type()),
// 	}

// 	if logs == nil {
// 		receipt["logs"] = [][]*ethtypes.Log{}
// 	}

// 	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
// 	if txData.GetTo() == nil {
// 		receipt["contractAddress"] = crypto.CreateAddress(from, txData.GetNonce())
// 	}

// 	if dynamicTx, ok := txData.(*evmtypes.DynamicFeeTx); ok {
// 		baseFee, err := b.BaseFee(blockRes)
// 		if err != nil {
// 			// tolerate the error for pruned node.
// 			b.Logger.Error("fetch basefee failed, node is pruned?", "height", txResult.Height, "error", err)
// 		} else {
// 			receipt["effectiveGasPrice"] = hexutil.Big(*dynamicTx.EffectiveGasPrice(baseFee))
// 		}
// 	}

// 	return receipt, nil
// }
