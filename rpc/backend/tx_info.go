package backend

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/cosmos/evm/mempool/txpool"
	servertypes "github.com/cosmos/evm/server/types"
	"github.com/cosmos/evm/utils"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	rpctypes "github.com/zeta-chain/node/rpc/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetTransactionByHash returns the Ethereum format transaction identified by Ethereum transaction hash
func (b *Backend) GetTransactionByHash(txHash common.Hash) (*rpctypes.RPCTransaction, error) {
	res, additional, err := b.GetTxByEthHash(txHash)
	if err != nil {
		return b.GetTransactionByHashPending(txHash)
	}

	block, err := b.CometBlockByNumber(rpctypes.BlockNumber(res.Height))
	if err != nil {
		return nil, err
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &block.Block.Height)
	if err != nil {
		b.Logger.Debug("block result not found", "height", block.Block.Height, "error", err.Error())
		return nil, fmt.Errorf("block result not found: %w", err)
	}

	var ethMsg *evmtypes.MsgEthereumTx
	// if additional fields are empty we can try to get MsgEthereumTx from sdk.Msg array
	if additional == nil {
		// #nosec G115 always in range
		if int(res.TxIndex) >= len(block.Block.Txs) {
			b.Logger.Error("tx out of bounds")
			return nil, fmt.Errorf("tx out of bounds")
		}
		tx, err := b.ClientCtx.TxConfig.TxDecoder()(block.Block.Txs[res.TxIndex])
		if err != nil {
			b.Logger.Debug("decoding failed", "error", err.Error())
			return nil, fmt.Errorf("failed to decode tx: %w", err)
		}

		// the `res.MsgIndex` is inferred from tx index, should be within the bound.
		msg, ok := tx.GetMsgs()[res.MsgIndex].(*evmtypes.MsgEthereumTx)
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
		msgs, _ := b.EthMsgsFromCometBlock(block, blockRes)
		for i := range msgs {
			if msgs[i].Hash == txHash.Hex() {
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
		b.Logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", blockRes.Height, "error", err)
	}

	height := uint64(res.Height)                       //#nosec G115 -- checked for int overflow already
	blockTime := uint64(block.Block.Time.UTC().Unix()) //#nosec G115 -- checked for int overflow already
	index := uint64(res.EthTxIndex)                    //#nosec G115 -- checked for int overflow already
	return rpctypes.NewTransactionFromMsg(
		ethMsg,
		common.BytesToHash(block.BlockID.Hash.Bytes()),
		height,
		blockTime,
		index,
		baseFee,
		b.ChainConfig(),
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

		if msg.Hash == txHash.Hex() {
			// use zero block values since it's not included in a block yet
			return rpctypes.NewTransactionFromMsg(
				msg,
				common.Hash{},
				uint64(0),
				uint64(0),
				uint64(0),
				nil,
				b.ChainConfig(),
				nil,
			)
		}
	}

	b.Logger.Debug("tx not found", "hash", hexTx)
	return nil, nil
}

// GetGasUsed returns gasUsed from transaction
func (b *Backend) GetGasUsed(res *servertypes.TxResult, price *big.Int, gas uint64) uint64 {
	return res.GasUsed
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (b *Backend) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	hexTx := hash.Hex()
	b.Logger.Debug("eth_getTransactionReceipt", "hash", hexTx)

	// Retry logic for transaction lookup with exponential backoff
	maxRetries := 10
	baseDelay := 50 * time.Millisecond

	var res *servertypes.TxResult
	var err error
	var additional *rpctypes.TxResultAdditionalFields

	for attempt := 0; attempt <= maxRetries; attempt++ {
		res, additional, err = b.GetTxByEthHash(hash)
		if err == nil {
			break // Found the transaction
		}

		if attempt == maxRetries/2 && b.Mempool != nil {
			status := b.Mempool.GetTxPool().Status(hash)
			if status == txpool.TxStatusUnknown {
				break
			}
		}

		if attempt < maxRetries {
			// Exponential backoff: 50ms, 100ms, 200ms
			delay := time.Duration(1<<attempt) * baseDelay
			b.Logger.Debug("tx not found, retrying", "hash", hexTx, "attempt", attempt+1, "delay", delay)
			time.Sleep(delay)
		}
	}

	if err != nil {
		b.Logger.Debug("tx not found after retries", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	resBlock, err := b.CometBlockByNumber(rpctypes.BlockNumber(res.Height))
	if err != nil {
		b.Logger.Debug("block not found", "height", res.Height, "error", err.Error())
		return nil, fmt.Errorf("block not found at height %d: %w", res.Height, err)
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &res.Height)
	if err != nil {
		b.Logger.Debug("failed to retrieve block results", "height", res.Height, "error", err.Error())
		return nil, fmt.Errorf("block result not found at height %d: %w", res.Height, err)
	}

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

		ethMsg = tx.GetMsgs()[res.MsgIndex].(*evmtypes.MsgEthereumTx)
		if ethMsg == nil {
			b.Logger.Error("failed to get eth msg")
			return nil, fmt.Errorf("failed to get eth msg")
		}
	} else {
		// if additional fields are not empty try to parse synthetic tx from them
		ethMsg = b.parseSyntethicTxFromAdditionalFields(additional)
		if ethMsg == nil {
			b.Logger.Error("failed to parse tx")
			return nil, fmt.Errorf("failed to parse tx")
		}
	}

	receipts, err := b.ReceiptsFromCometBlock(resBlock, blockRes, []*evmtypes.MsgEthereumTx{ethMsg})
	if err != nil {
		return nil, fmt.Errorf("failed to get receipts from comet block")
	}

	var signer ethtypes.Signer
	ethTx := ethMsg.AsTransaction()
	if ethTx.Protected() {
		signer = ethtypes.LatestSignerForChainID(ethTx.ChainId())
	} else {
		signer = ethtypes.FrontierSigner{}
	}
	from, err := ethMsg.GetSenderLegacy(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}

	return rpctypes.RPCMarshalReceipt(receipts[0], ethTx, from)
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
	height, err := utils.SafeUint64(resBlockResult.Height)
	if err != nil {
		return nil, err
	}
	// parse tx logs from events
	index := int(res.MsgIndex) // #nosec G701
	logs, err := evmtypes.DecodeMsgLogs(
		resBlockResult.TxsResults[res.TxIndex].Data,
		index,
		height,
	)
	if err != nil {
		b.Logger.Debug("failed to parse tx logs", "error", err.Error())
	}

	return logs, nil
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
// func (b *Backend) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
// 	b.Logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash.Hex(), "index", idx)
// 	sc, ok := b.ClientCtx.Client.(cmtrpcclient.SignClient)
// 	if !ok {
// 		return nil, errors.New("invalid rpc client")
// 	}

// 	block, err := sc.BlockByHash(b.Ctx, hash.Bytes())
// 	if err != nil {
// 		b.Logger.Debug("block not found", "hash", hash.Hex(), "error", err.Error())
// 		return nil, nil
// 	}

// 	if block.Block == nil {
// 		b.Logger.Debug("block not found", "hash", hash.Hex())
// 		return nil, nil
// 	}

// 	return b.GetTransactionByBlockAndIndex(block, idx)
// }

// // GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
// func (b *Backend) GetTransactionByBlockNumberAndIndex(blockNum rpctypes.BlockNumber, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
// 	b.Logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)

// 	block, err := b.CometBlockByNumber(blockNum)
// 	if err != nil {
// 		b.Logger.Debug("block not found", "height", blockNum.Int64(), "error", err.Error())
// 		return nil, nil
// 	}

// 	if block.Block == nil {
// 		b.Logger.Debug("block not found", "height", blockNum.Int64())
// 		return nil, nil
// 	}

// 	return b.GetTransactionByBlockAndIndex(block, idx)
// }

// GetTxByEthHash uses `/tx_query` to find transaction by ethereum tx hash
// TODO: Don't need to convert once hashing is fixed on CometBFT
// https://github.com/cometbft/cometbft/issues/6539
func (b *Backend) GetTxByEthHash(hash common.Hash) (*servertypes.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	if b.Indexer != nil {
		txRes, err := b.Indexer.GetByTxHash(hash)
		if err != nil {
			return nil, nil, err
		}
		return txRes, nil, nil
	}

	// fallback to CometBFT tx indexer
	query := fmt.Sprintf("%s.%s='%s'", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash, hash.Hex())
	txResult, txAdditional, err := b.QueryCometTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		return txs.GetTxByHash(hash)
	})
	if err != nil {
		return nil, nil, errorsmod.Wrapf(err, "GetTxByEthHash %s", hash.Hex())
	}
	return txResult, txAdditional, nil
}

// GetTxByTxIndex uses `/tx_query` to find transaction by tx index of valid ethereum txs
func (b *Backend) GetTxByTxIndex(height int64, index uint) (*servertypes.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	int32Index := int32(index) //#nosec G115 -- checked for int overflow already
	if b.Indexer != nil {
		// #nosec G115 always in range
		txRes, err := b.Indexer.GetByBlockAndIndex(height, int32Index)
		if err == nil {
			return txRes, nil, nil
		}
	}

	// fallback to CometBFT tx indexer
	query := fmt.Sprintf("tx.height=%d AND %s.%s=%d",
		height, evmtypes.TypeMsgEthereumTx,
		evmtypes.AttributeKeyTxIndex, index,
	)
	txResult, txAdditional, err := b.QueryCometTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		return txs.GetTxByTxIndex(int(index)) // #nosec G115 -- checked for int overflow already
	})
	if err != nil {
		return nil, nil, errorsmod.Wrapf(err, "GetTxByTxIndex %d %d", height, index)
	}
	return txResult, txAdditional, nil
}

// QueryCometTxIndexer query tx in CometBFT tx indexer
func (b *Backend) QueryCometTxIndexer(query string, txGetter func(*rpctypes.ParsedTxs) *rpctypes.ParsedTx) (*servertypes.TxResult, *rpctypes.TxResultAdditionalFields, error) {
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

// // GetTransactionByBlockAndIndex is the common code shared by `GetTransactionByBlockNumberAndIndex` and `GetTransactionByBlockHashAndIndex`.
// func (b *Backend) GetTransactionByBlockAndIndex(block *cmtrpctypes.ResultBlock, idx hexutil.Uint) (*rpctypes.RPCTransaction, error) {
// 	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &block.Block.Height)
// 	if err != nil {
// 		return nil, nil
// 	}

// 	var msg *evmtypes.MsgEthereumTx
// 	// find in tx indexer
// 	res, err := b.GetTxByTxIndex(block.Block.Height, uint(idx))
// 	if err == nil {
// 		tx, err := b.ClientCtx.TxConfig.TxDecoder()(block.Block.Txs[res.TxIndex])
// 		if err != nil {
// 			b.Logger.Debug("invalid ethereum tx", "height", block.Block.Header, "index", idx)
// 			return nil, nil
// 		}

// 		var ok bool
// 		// msgIndex is inferred from tx events, should be within bound.
// 		msg, ok = tx.GetMsgs()[res.MsgIndex].(*evmtypes.MsgEthereumTx)
// 		if !ok {
// 			b.Logger.Debug("invalid ethereum tx", "height", block.Block.Header, "index", idx)
// 			return nil, nil
// 		}
// 	} else {
// 		i := int(idx) // #nosec G115
// 		ethMsgs := b.EthMsgsFromCometBlock(block, blockRes)
// 		if i >= len(ethMsgs) {
// 			b.Logger.Debug("block txs index out of bound", "index", i)
// 			return nil, nil
// 		}

// 		msg = ethMsgs[i]
// 	}

// 	baseFee, err := b.BaseFee(blockRes)
// 	if err != nil {
// 		// handle the error for pruned node.
// 		b.Logger.Error("failed to fetch Base Fee from prunned block. Check node prunning configuration", "height", block.Block.Height, "error", err)
// 	}

// 	height := uint64(block.Block.Height)               // #nosec G115 -- checked for int overflow already
// 	blockTime := uint64(block.Block.Time.UTC().Unix()) // #nosec G115 -- checked for int overflow already
// 	index := uint64(idx)                               // #nosec G115 -- checked for int overflow already
// 	return rpctypes.NewTransactionFromMsg(
// 		msg,
// 		common.BytesToHash(block.BlockID.Hash),
// 		height,
// 		blockTime,
// 		index,
// 		baseFee,
// 		b.ChainConfig(),
// 	), nil
// }

// // CreateAccessList returns the list of addresses and storage keys used by the transaction (except for the
// // sender account and precompiles), plus the estimated gas if the access list were added to the transaction.
// func (b *Backend) CreateAccessList(
// 	args evmtypes.TransactionArgs,
// 	blockNrOrHash rpctypes.BlockNumberOrHash,
// 	overrides *json.RawMessage,
// ) (*rpctypes.AccessListResult, error) {
// 	accessList, gasUsed, vmErr, err := b.createAccessList(args, blockNrOrHash, overrides)
// 	if err != nil {
// 		return nil, err
// 	}

// 	hexGasUsed := hexutil.Uint64(gasUsed)
// 	result := rpctypes.AccessListResult{
// 		AccessList: &accessList,
// 		GasUsed:    &hexGasUsed,
// 	}
// 	if vmErr != nil {
// 		result.Error = vmErr.Error()
// 	}
// 	return &result, nil
// }

// // createAccessList creates the access list for the transaction.
// // It iteratively expands the access list until it converges.
// // If the access list has converged, the access list is returned.
// // If the access list has not converged, an error is returned.
// // If the transaction itself fails, an vmErr is returned.
// func (b *Backend) createAccessList(
// 	args evmtypes.TransactionArgs,
// 	blockNrOrHash rpctypes.BlockNumberOrHash,
// 	overrides *json.RawMessage,
// ) (ethtypes.AccessList, uint64, error, error) {
// 	args, err := b.SetTxDefaults(args)
// 	if err != nil {
// 		b.Logger.Error("failed to set tx defaults", "error", err)
// 		return nil, 0, nil, err
// 	}

// 	blockNum, err := b.BlockNumberFromComet(blockNrOrHash)
// 	if err != nil {
// 		b.Logger.Error("failed to get block number", "error", err)
// 		return nil, 0, nil, err
// 	}

// 	addressesToExclude, err := b.getAccessListExcludes(args, blockNum)
// 	if err != nil {
// 		b.Logger.Error("failed to get access list excludes", "error", err)
// 		return nil, 0, nil, err
// 	}

// 	prevTracer, traceArgs, err := b.initAccessListTracer(args, blockNum, addressesToExclude)
// 	if err != nil {
// 		b.Logger.Error("failed to init access list tracer", "error", err)
// 		return nil, 0, nil, err
// 	}

// 	// iteratively expand the access list
// 	for {
// 		accessList := prevTracer.AccessList()
// 		traceArgs.AccessList = &accessList
// 		res, err := b.DoCall(*traceArgs, blockNum, overrides)
// 		if err != nil {
// 			b.Logger.Error("failed to apply transaction", "error", err)
// 			return nil, 0, nil, fmt.Errorf("failed to apply transaction: %v err: %v", traceArgs.ToTransaction(ethtypes.LegacyTxType).Hash(), err)
// 		}

// 		// Check if access list has converged (no new addresses/slots accessed)
// 		newTracer := logger.NewAccessListTracer(accessList, addressesToExclude)
// 		if newTracer.Equal(prevTracer) {
// 			b.Logger.Info("access list converged", "accessList", accessList)
// 			var vmErr error
// 			if res.VmError != "" {
// 				b.Logger.Error("vm error after access list converged", "vmError", res.VmError)
// 				vmErr = errors.New(res.VmError)
// 			}
// 			return accessList, res.GasUsed, vmErr, nil
// 		}
// 		prevTracer = newTracer
// 	}
// }

// // getAccessListExcludes returns the addresses to exclude from the access list.
// // This includes the sender account, the target account (if provided), precompiles,
// // and any addresses in the authorization list.
// func (b *Backend) getAccessListExcludes(args evmtypes.TransactionArgs, blockNum rpctypes.BlockNumber) (map[common.Address]struct{}, error) {
// 	header, err := b.HeaderByNumber(blockNum)
// 	if err != nil {
// 		b.Logger.Error("failed to get header by number", "error", err)
// 		return nil, err
// 	}

// 	// exclude sender and precompiles
// 	addressesToExclude := make(map[common.Address]struct{})
// 	addressesToExclude[args.GetFrom()] = struct{}{}
// 	if args.To != nil {
// 		addressesToExclude[*args.To] = struct{}{}
// 	}

// 	isMerge := b.ChainConfig().MergeNetsplitBlock != nil
// 	precompiles := vm.ActivePrecompiles(b.ChainConfig().Rules(header.Number, isMerge, header.Time))
// 	for _, addr := range precompiles {
// 		addressesToExclude[addr] = struct{}{}
// 	}

// 	// check if enough gas was provided to cover all authorization lists
// 	maxAuthorizations := uint64(*args.Gas) / params.CallNewAccountGas
// 	if uint64(len(args.AuthorizationList)) > maxAuthorizations {
// 		b.Logger.Error("insufficient gas to process all authorizations", "maxAuthorizations", maxAuthorizations)
// 		return nil, errors.New("insufficient gas to process all authorizations")
// 	}

// 	for _, auth := range args.AuthorizationList {
// 		// validate authorization (duplicating stateTransition.validateAuthorization() logic from geth: https://github.com/ethereum/go-ethereum/blob/bf8f63dcd27e178bd373bfe41ea718efee2851dd/core/state_transition.go#L575)
// 		nonceOverflow := auth.Nonce+1 < auth.Nonce
// 		invalidChainID := !auth.ChainID.IsZero() && auth.ChainID.CmpBig(b.ChainConfig().ChainID) != 0
// 		if nonceOverflow || invalidChainID {
// 			b.Logger.Error("invalid authorization", "auth", auth)
// 			continue
// 		}
// 		if authority, err := auth.Authority(); err == nil {
// 			addressesToExclude[authority] = struct{}{}
// 		}
// 	}

// 	b.Logger.Debug("access list excludes created", "addressesToExclude", addressesToExclude)
// 	return addressesToExclude, nil
// }

// // initAccessListTracer initializes the access list tracer for the transaction.
// // It sets the default call arguments and creates a new access list tracer.
// // If an access list is provided in args, it uses that instead of creating a new one.
// func (b *Backend) initAccessListTracer(args evmtypes.TransactionArgs, blockNum rpctypes.BlockNumber, addressesToExclude map[common.Address]struct{}) (*logger.AccessListTracer, *evmtypes.TransactionArgs, error) {
// 	header, err := b.HeaderByNumber(blockNum)
// 	if err != nil {
// 		b.Logger.Error("failed to get header by number", "error", err)
// 		return nil, nil, err
// 	}

// 	if args.Nonce == nil {
// 		pending := blockNum == rpctypes.EthPendingBlockNumber
// 		nonce, err := b.getAccountNonce(args.GetFrom(), pending, blockNum.Int64(), b.Logger)
// 		if err != nil {
// 			b.Logger.Error("failed to get account nonce", "error", err)
// 			return nil, nil, err
// 		}
// 		nonce64 := hexutil.Uint64(nonce)
// 		args.Nonce = &nonce64
// 	}
// 	if err = args.CallDefaults(b.RPCGasCap(), header.BaseFee, b.ChainConfig().ChainID); err != nil {
// 		b.Logger.Error("failed to set default call args", "error", err)
// 		return nil, nil, err
// 	}

// 	tracer := logger.NewAccessListTracer(nil, addressesToExclude)
// 	if args.AccessList != nil {
// 		tracer = logger.NewAccessListTracer(*args.AccessList, addressesToExclude)
// 	}

// 	b.Logger.Debug("access list tracer initialized", "tracer", tracer)
// 	return tracer, &args, nil
// }
