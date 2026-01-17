package backend

import (
	"context"
	"math/big"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	cmtversion "github.com/cometbft/cometbft/proto/tendermint/version"
	cmtrpcclient "github.com/cometbft/cometbft/rpc/client"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	"github.com/cometbft/cometbft/version"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/rpc/backend/mocks"
	rpc "github.com/zeta-chain/node/rpc/types"
)

// Client defines a mocked object that implements the Tendermint JSON-RPC Client
// interface. It allows for performing Client queries without having to run a
// Tendermint RPC Client server.
//
// To use a mock method it has to be registered in a given test.
var _ cmtrpcclient.Client = &mocks.Client{}

// Tx Search
func RegisterTxSearch(client *mocks.Client, query string, txBz []byte) {
	resulTxs := []*cmtrpctypes.ResultTx{{Tx: txBz}}
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&cmtrpctypes.ResultTxSearch{Txs: resulTxs, TotalCount: 1}, nil)
}

func RegisterTxSearchWithTxResult(client *mocks.Client, query string, txBz []byte, res abci.ExecTxResult) {
	resulTxs := []*cmtrpctypes.ResultTx{{Tx: txBz, Height: 1, TxResult: res}}
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&cmtrpctypes.ResultTxSearch{Txs: resulTxs, TotalCount: 1}, nil)
}

func RegisterTxSearchEmpty(client *mocks.Client, query string) {
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(&cmtrpctypes.ResultTxSearch{}, nil)
}

func RegisterTxSearchError(client *mocks.Client, query string) {
	client.On("TxSearch", rpc.ContextWithHeight(1), query, false, (*int)(nil), (*int)(nil), "").
		Return(nil, errortypes.ErrInvalidRequest)
}

// Broadcast Tx
func RegisterBroadcastTx(client *mocks.Client, tx types.Tx) {
	client.On("BroadcastTxSync", context.Background(), tx).
		Return(&cmtrpctypes.ResultBroadcastTx{}, nil)
}

func RegisterBroadcastTxError(client *mocks.Client, tx types.Tx) {
	client.On("BroadcastTxSync", context.Background(), tx).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Unconfirmed Transactions
func RegisterUnconfirmedTxs(client *mocks.Client, limit *int, txs []types.Tx) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(&cmtrpctypes.ResultUnconfirmedTxs{Txs: txs}, nil)
}

func RegisterUnconfirmedTxsEmpty(client *mocks.Client, limit *int) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(&cmtrpctypes.ResultUnconfirmedTxs{
			Txs: make([]types.Tx, 2),
		}, nil)
}

func RegisterUnconfirmedTxsError(client *mocks.Client, limit *int) {
	client.On("UnconfirmedTxs", rpc.ContextWithHeight(1), limit).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Status
func RegisterStatus(client *mocks.Client) {
	client.On("Status", rpc.ContextWithHeight(1)).
		Return(&cmtrpctypes.ResultStatus{}, nil)
}

func RegisterStatusError(client *mocks.Client) {
	client.On("Status", rpc.ContextWithHeight(1)).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Block
func RegisterBlockMultipleTxs(
	client *mocks.Client,
	height int64,
	txs []types.Tx,
) (*cmtrpctypes.ResultBlock, error) {
	block := types.MakeBlock(height, txs, nil, nil)
	block.ChainID = ChainID.ChainID
	resBlock := &cmtrpctypes.ResultBlock{Block: block}
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil)
	return resBlock, nil
}

func RegisterBlock(
	client *mocks.Client,
	height int64,
	tx []byte,
) (*cmtrpctypes.ResultBlock, error) {
	// without tx
	if tx == nil {
		emptyBlock := types.MakeBlock(height, []types.Tx{}, nil, nil)
		emptyBlock.ChainID = ChainID.ChainID
		blockHash := common.BigToHash(big.NewInt(height)).Bytes()
		resBlock := &cmtrpctypes.ResultBlock{Block: emptyBlock, BlockID: types.BlockID{Hash: bytes.HexBytes(blockHash)}}
		client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil)
		return resBlock, nil
	}

	// with tx
	block := types.MakeBlock(height, []types.Tx{tx}, nil, nil)
	block.ChainID = ChainID.ChainID
	blockHash := common.BigToHash(big.NewInt(height)).Bytes()
	resBlock := &cmtrpctypes.ResultBlock{Block: block, BlockID: types.BlockID{Hash: bytes.HexBytes(blockHash)}}
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).Return(resBlock, nil)
	return resBlock, nil
}

// Block returns error
func RegisterBlockError(client *mocks.Client, height int64) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest)
}

// Block not found
func RegisterBlockNotFound(
	client *mocks.Client,
	height int64,
) (*cmtrpctypes.ResultBlock, error) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(&cmtrpctypes.ResultBlock{Block: nil}, nil)

	return &cmtrpctypes.ResultBlock{Block: nil}, nil
}

func RegisterBlockResultsWithTxResults(
	client *mocks.Client,
	height int64,
	txResults []*abci.ExecTxResult,
) (*cmtrpctypes.ResultBlockResults, error) {
	res := &cmtrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: txResults,
	}

	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil)
	return res, nil
}

// Block panic
func RegisterBlockPanic(client *mocks.Client, height int64) {
	client.On("Block", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(func(context.Context, *int64) *cmtrpctypes.ResultBlock {
			panic("Block call panic")
		}, nil)
}

func TestRegisterBlock(t *testing.T) {
	client := mocks.NewClient(t)
	height := rpc.BlockNumber(1).Int64()
	_, err := RegisterBlock(client, height, nil)
	require.NoError(t, err)

	res, err := client.Block(rpc.ContextWithHeight(height), &height)

	emptyBlock := types.MakeBlock(height, []types.Tx{}, nil, nil)
	emptyBlock.ChainID = ChainID.ChainID
	blockHash := common.BigToHash(big.NewInt(height)).Bytes()
	resBlock := &cmtrpctypes.ResultBlock{Block: emptyBlock, BlockID: types.BlockID{Hash: bytes.HexBytes(blockHash)}}
	require.Equal(t, resBlock, res)
	require.NoError(t, err)
}

// ConsensusParams
func RegisterConsensusParams(client *mocks.Client, height int64) {
	consensusParams := types.DefaultConsensusParams()
	client.On("ConsensusParams", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(&cmtrpctypes.ResultConsensusParams{ConsensusParams: *consensusParams}, nil)
}

func RegisterConsensusParamsError(client *mocks.Client, height int64) {
	client.On("ConsensusParams", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest)
}

func TestRegisterConsensusParams(t *testing.T) {
	client := mocks.NewClient(t)
	height := int64(1)
	RegisterConsensusParams(client, height)

	res, err := client.ConsensusParams(rpc.ContextWithHeight(height), &height)
	consensusParams := types.DefaultConsensusParams()
	require.Equal(t, &cmtrpctypes.ResultConsensusParams{ConsensusParams: *consensusParams}, res)
	require.NoError(t, err)
}

// BlockResults

func RegisterBlockResultsWithEventLog(client *mocks.Client, height int64) (*cmtrpctypes.ResultBlockResults, error) {
	res := &cmtrpctypes.ResultBlockResults{
		Height: height,
		TxsResults: []*abci.ExecTxResult{
			{Code: 0, GasUsed: 0, Events: []abci.Event{{
				Type: evmtypes.EventTypeTxLog,
				Attributes: []abci.EventAttribute{{
					Key:   evmtypes.AttributeKeyTxLog,
					Value: "{\"test\": \"hello\"}", // TODO refactor the value to unmarshall to a evmtypes.Log struct successfully
					Index: true,
				}},
			}}},
		},
	}
	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil)
	return res, nil
}

func RegisterBlockResults(
	client *mocks.Client,
	height int64,
) (*cmtrpctypes.ResultBlockResults, error) {
	res := &cmtrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: []*abci.ExecTxResult{{Code: 0, GasUsed: 0}},
	}

	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(res, nil)
	return res, nil
}

func RegisterBlockResultsError(client *mocks.Client, height int64) {
	client.On("BlockResults", rpc.ContextWithHeight(height), mock.AnythingOfType("*int64")).
		Return(nil, errortypes.ErrInvalidRequest)
}

func TestRegisterBlockResults(t *testing.T) {
	client := mocks.NewClient(t)
	height := int64(1)
	_, err := RegisterBlockResults(client, height)
	require.NoError(t, err)

	res, err := client.BlockResults(rpc.ContextWithHeight(height), &height)
	expRes := &cmtrpctypes.ResultBlockResults{
		Height:     height,
		TxsResults: []*abci.ExecTxResult{{Code: 0, GasUsed: 0}},
	}
	require.Equal(t, expRes, res)
	require.NoError(t, err)
}

// BlockByHash
func RegisterBlockByHash(
	client *mocks.Client,
	_ common.Hash,
	tx []byte,
) (*cmtrpctypes.ResultBlock, error) {
	block := types.MakeBlock(1, []types.Tx{tx}, nil, nil)
	resBlock := &cmtrpctypes.ResultBlock{Block: block}

	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(resBlock, nil)
	return resBlock, nil
}

func RegisterBlockByHashError(client *mocks.Client, _ common.Hash, _ []byte) {
	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, errortypes.ErrInvalidRequest)
}

func RegisterBlockByHashNotFound(client *mocks.Client, _ common.Hash, _ []byte) {
	client.On("BlockByHash", rpc.ContextWithHeight(1), []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, nil)
}

// HeaderByHash
func RegisterHeaderByHash(
	client *mocks.Client,
	_ common.Hash,
	_ []byte,
) (*cmtrpctypes.ResultHeader, error) {
	header := &types.Header{
		Version: cmtversion.Consensus{Block: version.BlockProtocol, App: 0},
		Height:  1,
	}
	resHeader := &cmtrpctypes.ResultHeader{
		Header: header,
	}

	client.On("HeaderByHash", rpc.ContextWithHeight(1), bytes.HexBytes{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(resHeader, nil)
	return resHeader, nil
}

func RegisterHeaderByHashError(client *mocks.Client, _ common.Hash, _ []byte) {
	client.On("HeaderByHash", rpc.ContextWithHeight(1), bytes.HexBytes{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, errortypes.ErrInvalidRequest)
}

func RegisterHeaderByHashNotFound(client *mocks.Client, _ common.Hash, _ []byte) {
	client.On("HeaderByHash", rpc.ContextWithHeight(1), bytes.HexBytes{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}).
		Return(nil, nil)
}

func RegisterABCIQueryWithOptions(
	client *mocks.Client,
	height int64,
	path string,
	data bytes.HexBytes,
	opts cmtrpcclient.ABCIQueryOptions,
) {
	client.On("ABCIQueryWithOptions", context.Background(), path, data, opts).
		Return(&cmtrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  []byte{2}, // TODO replace with data.Bytes(),
				Height: height,
			},
		}, nil)
}

func RegisterABCIQueryWithOptionsError(
	clients *mocks.Client,
	path string,
	data bytes.HexBytes,
	opts cmtrpcclient.ABCIQueryOptions,
) {
	clients.On("ABCIQueryWithOptions", context.Background(), path, data, opts).
		Return(nil, errortypes.ErrInvalidRequest)
}

func RegisterABCIQueryAccount(
	clients *mocks.Client,
	data bytes.HexBytes,
	opts cmtrpcclient.ABCIQueryOptions,
	acc client.Account,
) {
	baseAccount := authtypes.NewBaseAccount(
		acc.GetAddress(),
		acc.GetPubKey(),
		acc.GetAccountNumber(),
		acc.GetSequence(),
	)
	accAny, _ := codectypes.NewAnyWithValue(baseAccount)
	accResponse := authtypes.QueryAccountResponse{Account: accAny}
	respBz, _ := accResponse.Marshal()
	clients.On("ABCIQueryWithOptions", context.Background(), "/cosmos.auth.v1beta1.Query/Account", data, opts).
		Return(&cmtrpctypes.ResultABCIQuery{
			Response: abci.ResponseQuery{
				Value:  respBz,
				Height: 1,
			},
		}, nil)
}
