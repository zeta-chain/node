package filters

// import (
// 	"context"
// 	"errors"
// 	"math/big"
// 	"testing"

// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/eth/filters"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"

// 	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"
// 	comettypes "github.com/cometbft/cometbft/types"

// 	filtermocks "github.com/zeta-chain/node/rpc/namespaces/ethereum/eth/filters/mocks"
// 	rpctypes "github.com/zeta-chain/node/rpc/types"

// 	"cosmossdk.io/log"
// )

// type MockBackend struct {
// 	mock.Mock
// }

// func (m *MockBackend) GetBlockByNumber(blockNum rpctypes.BlockNumber, fullTx bool) (map[string]interface{}, error) {
// 	panic("implement me")
// }

// func (m *MockBackend) HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error) {
// 	panic("implement me")
// }

// func (m *MockBackend) GetLogs(blockHash common.Hash) ([][]*ethtypes.Log, error) {
// 	panic("implement me")
// }

// func (m *MockBackend) GetLogsByHeight(i *int64) ([][]*ethtypes.Log, error) {
// 	panic("implement me")
// }

// func (m *MockBackend) BloomStatus() (uint64, uint64) {
// 	panic("implement me")
// }

// func (m *MockBackend) RPCFilterCap() int32 {
// 	panic("implement me")
// }

// func (m *MockBackend) RPCLogsCap() int32 {
// 	panic("implement me")
// }

// func (m *MockBackend) RPCBlockRangeCap() int32 {
// 	panic("implement me")
// }

// func (m *MockBackend) CometBlockByHash(hash common.Hash) (*cmtrpctypes.ResultBlock, error) {
// 	args := m.Called(hash)
// 	return args.Get(0).(*cmtrpctypes.ResultBlock), args.Error(1)
// }

// func (m *MockBackend) CometBlockResultByNumber(height *int64) (*cmtrpctypes.ResultBlockResults, error) {
// 	args := m.Called(height)
// 	return args.Get(0).(*cmtrpctypes.ResultBlockResults), args.Error(1)
// }

// func (m *MockBackend) BlockBloomFromCometBlock(blockRes *cmtrpctypes.ResultBlockResults) (ethtypes.Bloom, error) {
// 	args := m.Called(blockRes)
// 	return args.Get(0).(ethtypes.Bloom), args.Error(1)
// }

// func (m *MockBackend) HeaderByNumber(blockNum rpctypes.BlockNumber) (*ethtypes.Header, error) {
// 	args := m.Called(blockNum)
// 	return args.Get(0).(*ethtypes.Header), args.Error(1)
// }

// func TestLogs(t *testing.T) {
// 	blockHeight := int64(100)
// 	fakeHeader := &ethtypes.Header{Number: big.NewInt(blockHeight)}
// 	fakeBlockRes := &cmtrpctypes.ResultBlockResults{Height: blockHeight}
// 	fakeBloom := ethtypes.Bloom{}
// 	blockHash := common.HexToHash("0xabc")
// 	fakeBlock := &cmtrpctypes.ResultBlock{Block: &comettypes.Block{Header: comettypes.Header{Height: blockHeight}}}

// 	tests := []struct {
// 		name      string
// 		errorStep string
// 		prepare   func() *MockBackend
// 		criteria  filters.FilterCriteria
// 		expectErr bool
// 		expectMsg string
// 	}{
// 		{
// 			name:      "HeaderByNumber returns error",
// 			errorStep: "HeaderByNumber",
// 			prepare: func() *MockBackend {
// 				backend := &MockBackend{}
// 				backend.On("HeaderByNumber", mock.Anything).Return((*ethtypes.Header)(nil), errors.New("header error"))
// 				return backend
// 			},
// 			criteria: filters.FilterCriteria{
// 				FromBlock: big.NewInt(blockHeight),
// 				ToBlock:   big.NewInt(blockHeight),
// 			},
// 			expectErr: true,
// 			expectMsg: "header error",
// 		},
// 		{
// 			name:      "CometBlockResultByNumber returns error",
// 			errorStep: "CometBlockResultByNumber",
// 			prepare: func() *MockBackend {
// 				backend := &MockBackend{}
// 				backend.On("HeaderByNumber", mock.Anything).Return(fakeHeader, nil)
// 				backend.On("CometBlockResultByNumber", &blockHeight).Return((*cmtrpctypes.ResultBlockResults)(nil), errors.New("block result error"))
// 				return backend
// 			},
// 			criteria: filters.FilterCriteria{
// 				FromBlock: big.NewInt(blockHeight),
// 				ToBlock:   big.NewInt(blockHeight),
// 			},
// 			expectErr: true,
// 			expectMsg: "block result error",
// 		},
// 		{
// 			name:      "BlockBloom returns error",
// 			errorStep: "BlockBloom",
// 			prepare: func() *MockBackend {
// 				backend := &MockBackend{}
// 				backend.On("HeaderByNumber", mock.Anything).Return(fakeHeader, nil)
// 				backend.On("CometBlockResultByNumber", &blockHeight).Return(fakeBlockRes, nil)
// 				backend.On("BlockBloomFromCometBlock", fakeBlockRes).Return(ethtypes.Bloom{}, errors.New("bloom error"))
// 				return backend
// 			},
// 			criteria: filters.FilterCriteria{
// 				FromBlock: big.NewInt(blockHeight),
// 				ToBlock:   big.NewInt(blockHeight),
// 			},
// 			expectErr: true,
// 			expectMsg: "bloom error",
// 		},
// 		{
// 			name:      "Single block by BlockHash",
// 			errorStep: "none",
// 			prepare: func() *MockBackend {
// 				backend := &MockBackend{}
// 				backend.On("CometBlockByHash", blockHash).Return(fakeBlock, nil)
// 				backend.On("CometBlockResultByNumber", &blockHeight).Return(fakeBlockRes, nil)
// 				backend.On("BlockBloomFromCometBlock", fakeBlockRes).Return(fakeBloom, nil)
// 				return backend
// 			},
// 			criteria: filters.FilterCriteria{
// 				BlockHash: &blockHash,
// 			},
// 			expectErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			logger := log.NewNopLogger()
// 			backend := tt.prepare()

// 			var filter *Filter
// 			if tt.criteria.BlockHash != nil && *tt.criteria.BlockHash != (common.Hash{}) {
// 				filter = NewBlockFilter(logger, backend, tt.criteria)
// 			} else {
// 				filter = NewRangeFilter(logger, backend, blockHeight, blockHeight, nil, nil)
// 			}

// 			logs, err := filter.Logs(context.Background(), 1000, 100)
// 			if tt.expectErr {
// 				require.Error(t, err)
// 				require.Contains(t, err.Error(), tt.expectMsg)
// 				require.Nil(t, logs)
// 			} else {
// 				require.NoError(t, err)
// 				require.NotNil(t, logs)
// 			}

// 			backend.AssertExpectations(t)
// 		})
// 	}
// }

// func TestFilter(t *testing.T) {
// 	logger := log.NewNopLogger()
// 	testCases := []struct {
// 		name         string
// 		filter       filters.FilterCriteria
// 		expectations func(b *filtermocks.Backend)
// 		expLogs      []*ethtypes.Log
// 		expErr       string
// 	}{
// 		{
// 			name:   "invalid block range returns error",
// 			filter: filters.FilterCriteria{FromBlock: big.NewInt(100), ToBlock: big.NewInt(110)},
// 			expectations: func(b *filtermocks.Backend) {
// 				b.EXPECT().HeaderByNumber(rpctypes.EthLatestBlockNumber).Return(&ethtypes.Header{Number: big.NewInt(5)}, nil)
// 			},
// 			expErr: "invalid block range params",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			backend := filtermocks.NewBackend(t)
// 			f := newFilter(logger, backend, tc.filter, nil)
// 			tc.expectations(backend)
// 			logs, err := f.Logs(context.Background(), 15, 50)
// 			if tc.expErr != "" {
// 				require.ErrorContains(t, err, tc.expErr)
// 			} else {
// 				require.NoError(t, err)
// 			}

// 			if tc.expLogs != nil {
// 				require.Equal(t, tc.expLogs, logs)
// 			}
// 		})
// 	}
// }
