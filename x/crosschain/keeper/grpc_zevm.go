package keeper

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrpctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func (k Keeper) ZEVMGetBlock(c context.Context, req *types.QueryZEVMGetBlockByNumberRequest) (*types.QueryZEVMGetBlockByNumberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	rpcclient := types.ClientCtx.Client
	if rpcclient == nil {
		return nil, status.Error(codes.Internal, "rpc client is not initialized")
	}

	blockResults, err := GetTendermintBlockResultsByNumber(ctx, int64(req.Height))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	block, err := GetTendermintBlockByNumber(ctx, int64(req.Height))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	txDecoder := types.ClientCtx.TxConfig.TxDecoder()
	transactionHashes := make([]string, 0)
	for idx, txResult := range blockResults.TxsResults {
		logs, err := GetEthLogsFromEvents(txResult.Events)
		if err != nil || len(logs) == 0 {
			continue
		}
		txBytes := block.Block.Txs[idx]
		tx, err := txDecoder(txBytes)
		if err != nil {
			continue
		}
		_, ok := tx.GetMsgs()[0].(*evmtypes.MsgEthereumTx)
		if ok { // skip MsgEthereumTx; these txs are handled by ethermint JSON-RPC server
			continue
		}

		transactionHashes = append(transactionHashes, fmt.Sprintf("0x%x", block.Block.Txs[idx].Hash()))
	}
	return &types.QueryZEVMGetBlockByNumberResponse{
		Number:       fmt.Sprintf("0x%x", req.Height),
		Transactions: transactionHashes,
	}, nil
}

func (k Keeper) ZEVMGetTransactionReceipt(c context.Context, req *types.QueryZEVMGetTransactionReceiptRequest) (*types.QueryZEVMGetTransactionReceiptResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	rpcclient := types.ClientCtx.Client
	if rpcclient == nil {
		return nil, status.Error(codes.Internal, "rpc client is not initialized")
	}

	var hash string
	if len(req.Hash) == 66 && req.Hash[:2] == "0x" { // eth format
		hash = strings.ToUpper(req.Hash[2:])
	} else {
		hash = req.Hash
	}
	query := fmt.Sprintf("ethereum_tx.txHash='%s'", hash)
	res, err := rpcclient.TxSearch(c, query, false, nil, nil, "asc")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(res.Txs) == 0 {
		return nil, status.Error(codes.NotFound, "transaction not found")
	}

	txRaw := res.Txs[0]
	block, err := GetTendermintBlockByNumber(ctx, txRaw.Height)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	blockHash := ethcommon.BytesToHash(block.BlockID.Hash.Bytes())
	blockNumber := fmt.Sprintf("0x%x", txRaw.Height)

	tx, err := types.ClientCtx.TxConfig.TxDecoder()(txRaw.Tx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	msg0 := tx.GetMsgs()[0]
	fromAddress := ethcommon.BytesToAddress(msg0.GetSigners()[0].Bytes())

	status0 := "0x0"
	if txRaw.TxResult.Code == 0 { // code 0 means success for cosmos tx; ref https://docs.cosmos.network/main/core/baseapp#delivertx
		status0 = "0x1" // 1 = success in ethereum;
	}
	hash = ethcommon.BytesToHash(txRaw.Hash.Bytes()).Hex()

	logs, err := GetEthLogsFromEvents(txRaw.TxResult.Events)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//logs := make([]*types.Log, )
	return &types.QueryZEVMGetTransactionReceiptResponse{
		BlockHash:         blockHash.Hex(),
		BlockNumber:       blockNumber,
		ContractAddress:   "", // this is the contract created by the transaction, if any
		CumulativeGasUsed: "0x0",
		From:              fromAddress.Hex(),
		GasUsed:           fmt.Sprintf("0x%x", txRaw.TxResult.GasUsed),
		LogsBloom:         "", //FIXME: add proper bloom filter
		Status:            status0,
		To:                "",
		TransactionHash:   hash,
		TransactionIndex:  fmt.Sprintf("0x%x", txRaw.Index), // FIXME: does this make sense?
		Logs:              logs,
	}, nil
}

func (k Keeper) ZEVMGetTransaction(c context.Context, req *types.QueryZEVMGetTransactionRequest) (*types.QueryZEVMGetTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	rpcclient := types.ClientCtx.Client

	if rpcclient == nil {
		return nil, status.Error(codes.Internal, "rpc client is not initialized")
	}
	var hash string
	if len(req.Hash) == 66 && req.Hash[:2] == "0x" { // eth format
		hash = strings.ToUpper(req.Hash[2:])
	} else {
		hash = req.Hash
	}
	query := fmt.Sprintf("ethereum_tx.txHash='%s'", hash)
	res, err := rpcclient.TxSearch(c, query, false, nil, nil, "asc")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(res.Txs) == 0 {
		return nil, status.Error(codes.NotFound, "transaction not found")
	}

	txRaw := res.Txs[0]

	block, err := GetTendermintBlockByNumber(ctx, txRaw.Height)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	tx, err := types.ClientCtx.TxConfig.TxDecoder()(txRaw.Tx)
	msg0 := tx.GetMsgs()[0]
	fromAddress := ethcommon.BytesToAddress(msg0.GetSigners()[0].Bytes())
	chainID, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var blockNumber string
	H := ethcommon.BytesToHash(txRaw.Hash.Bytes())
	hash = H.Hex()
	blockNumber = fmt.Sprintf("0x%x", txRaw.Height)
	blockHash := ethcommon.BytesToHash(block.BlockID.Hash.Bytes())

	return &types.QueryZEVMGetTransactionResponse{
		BlockHash:        blockHash.Hex(),
		BlockNumber:      blockNumber,
		From:             fromAddress.Hex(), // FIXME: this should be the EOA on external chain?
		Gas:              fmt.Sprintf("0x%x", txRaw.TxResult.GasWanted),
		GasPrice:         "",
		Hash:             hash, // Note: This is the cosmos tx hash, in ethereum format (0x prefixed)
		Input:            "",
		Nonce:            "0",
		To:               "",
		TransactionIndex: "0",
		Value:            "0",
		Type:             "0x88",
		AccessList:       nil,
		ChainId:          chainID.String(),
		V:                "",
		R:                "",
		S:                "",
	}, nil
}

func GetTendermintBlockByNumber(ctx sdk.Context, blockNum int64) (*tmrpctypes.ResultBlock, error) {
	rpcclient := types.ClientCtx.Client
	if rpcclient == nil {
		return nil, fmt.Errorf("rpc client is not initialized")
	}
	height := blockNum
	if height <= 0 {
		height = ctx.BlockHeight()
	}
	resBlock, err := rpcclient.Block(sdk.WrapSDKContext(ctx), &height)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by height %d: %w", height, err)
	}

	if resBlock.Block == nil {
		return nil, nil
	}

	return resBlock, nil
}

func GetTendermintBlockResultsByNumber(ctx sdk.Context, blockNum int64) (*tmrpctypes.ResultBlockResults, error) {
	rpcclient := types.ClientCtx.Client
	if rpcclient == nil {
		return nil, fmt.Errorf("rpc client is not initialized")
	}
	height := blockNum
	if height <= 0 {
		height = ctx.BlockHeight()
	}
	resBlock, err := rpcclient.BlockResults(sdk.WrapSDKContext(ctx), &height)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by height %d: %w", height, err)
	}

	if resBlock == nil {
		return nil, nil
	}

	return resBlock, nil
}

func GetEthLogsFromEvents(events []abci.Event) ([]*types.Log, error) {
	logs := make([]*types.Log, 0)
	for _, event := range events {
		if event.Type == evmtypes.EventTypeTxLog {
			for _, attr := range event.Attributes {
				if !bytes.Equal(attr.Key, []byte(evmtypes.AttributeKeyTxLog)) {
					continue
				}

				var log types.Log
				err := json.Unmarshal(attr.Value, &log)
				if err != nil {
					return nil, err
				}
				data, err := base64.StdEncoding.DecodeString(log.Data)
				if err == nil {
					log.Data = "0x" + hex.EncodeToString(data)
				}
				logs = append(logs, &log)

			}
		}
	}
	return logs, nil
}
