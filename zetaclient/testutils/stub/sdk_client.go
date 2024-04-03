package stub

import (
	"context"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/client/mock"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type MockSDKClient struct {
	mock.Client
	err error
}

func (c MockSDKClient) BroadcastTxCommit(_ context.Context, _ tmtypes.Tx) (*coretypes.ResultBroadcastTxCommit, error) {
	return nil, c.err
}

func (c MockSDKClient) BroadcastTxAsync(_ context.Context, _ tmtypes.Tx) (*coretypes.ResultBroadcastTx, error) {
	return nil, c.err
}

func (c MockSDKClient) BroadcastTxSync(_ context.Context, _ tmtypes.Tx) (*coretypes.ResultBroadcastTx, error) {
	return nil, c.err
}

func (c MockSDKClient) Tx(_ context.Context, _ []byte, _ bool) (*coretypes.ResultTx, error) {
	return &coretypes.ResultTx{
		Hash:   bytes.HexBytes{},
		Height: 0,
		Index:  0,
		TxResult: abci.ResponseDeliverTx{
			Log: "",
		},
		Tx:    []byte{},
		Proof: tmtypes.TxProof{},
	}, c.err
}

func (c MockSDKClient) Block(_ context.Context, _ *int64) (*coretypes.ResultBlock, error) {
	return &coretypes.ResultBlock{Block: &types.Block{
		Header:   types.Header{},
		Data:     types.Data{},
		Evidence: types.EvidenceData{},
	}}, c.err
}

func NewMockSDKClientWithErr(err error) *MockSDKClient {
	return &MockSDKClient{
		Client: mock.Client{},
		err:    err,
	}
}
