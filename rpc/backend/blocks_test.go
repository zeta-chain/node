package backend

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"google.golang.org/grpc/metadata"

	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/bytes"
	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	cmttypes "github.com/cometbft/cometbft/types"

	utiltx "github.com/cosmos/evm/testutil/tx"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/node/rpc/backend/mocks"
	ethrpc "github.com/zeta-chain/node/rpc/types"
	"github.com/zeta-chain/node/testutil/sample"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// simpleMsg is a test message that doesn't implement MsgEthereumTx
type simpleMsg struct{}

func (m simpleMsg) Route() string                               { return "test" }
func (m simpleMsg) Type() string                                { return "test" }
func (m simpleMsg) ValidateBasic() error                        { return nil }
func (m simpleMsg) GetSigners() []sdk.AccAddress                { return nil }
func (m simpleMsg) ProtoMessage()                               { panic("not implemented") }
func (m simpleMsg) Reset()                                      { panic("not implemented") }
func (m simpleMsg) String() string                              { panic("not implemented") }
func (m simpleMsg) Bytes() []byte                               { panic("not implemented") }
func (m simpleMsg) VerifySignature(msg []byte, sig []byte) bool { panic("not implemented") }

func (s *TestSuite) TestBlockNumber() {
	testCases := []struct {
		name           string
		registerMock   func()
		expBlockNumber hexutil.Uint64
		expPass        bool
	}{
		{
			"fail - invalid block header height",
			func() {
				var header metadata.MD
				height := int64(1)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParamsInvalidHeight(QueryClient, &header, height)
			},
			0x0,
			false,
		},
		{
			"fail - invalid block header",
			func() {
				var header metadata.MD
				height := int64(1)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParamsInvalidHeader(QueryClient, &header, height)
			},
			0x0,
			false,
		},
		{
			"pass - app state header height 1",
			func() {
				var header metadata.MD
				height := int64(1)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(QueryClient, &header, height)
			},
			0x1,
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock()

			blockNumber, err := s.backend.BlockNumber()

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().Equal(tc.expBlockNumber, blockNumber)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestGetBlockByNumber() {
	var (
		blockRes *cmtrpctypes.ResultBlockResults
		resBlock *cmtrpctypes.ResultBlock
	)
	msgEthereumTx, bz := s.buildEthereumTx()

	testCases := []struct {
		name         string
		blockNumber  ethrpc.BlockNumber
		fullTx       bool
		baseFee      *big.Int
		validator    sdk.AccAddress
		tx           *evmtypes.MsgEthereumTx
		txBz         []byte
		registerMock func(ethrpc.BlockNumber, math.Int, sdk.AccAddress, []byte)
		expNoop      bool
		expPass      bool
	}{
		{
			"pass - tendermint block not found",
			ethrpc.BlockNumber(1),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(blockNum ethrpc.BlockNumber, _ math.Int, _ sdk.AccAddress, _ []byte) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, height)
			},
			true,
			true,
		},
		{
			"pass - block not found (e.g. request block height that is greater than current one)",
			ethrpc.BlockNumber(1),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(blockNum ethrpc.BlockNumber, _ math.Int, _ sdk.AccAddress, _ []byte) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resBlock, _ = RegisterBlockNotFound(client, height)
			},
			true,
			true,
		},
		{
			"pass - block results error",
			ethrpc.BlockNumber(1),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(blockNum ethrpc.BlockNumber, _ math.Int, _ sdk.AccAddress, txBz []byte) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resBlock, _ = RegisterBlock(client, height, txBz)
				RegisterBlockResultsError(client, blockNum.Int64())
			},
			true,
			true,
		},
		{
			"pass - without tx",
			ethrpc.BlockNumber(1),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(blockNum ethrpc.BlockNumber, baseFee math.Int, validator sdk.AccAddress, txBz []byte) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resBlock, _ = RegisterBlock(client, height, txBz)
				blockRes, _ = RegisterBlockResults(client, blockNum.Int64())
				RegisterConsensusParams(client, height)

				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)
			},
			false,
			true,
		},
		{
			"pass - with tx",
			ethrpc.BlockNumber(1),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			msgEthereumTx,
			bz,
			func(blockNum ethrpc.BlockNumber, baseFee math.Int, validator sdk.AccAddress, txBz []byte) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resBlock, _ = RegisterBlock(client, height, txBz)
				blockRes, _ = RegisterBlockResults(client, blockNum.Int64())
				RegisterConsensusParams(client, height)

				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)
			},
			false,
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock(tc.blockNumber, math.NewIntFromBigInt(tc.baseFee), tc.validator, tc.txBz)

			block, err := s.backend.GetBlockByNumber(tc.blockNumber, tc.fullTx)

			if tc.expPass {
				if tc.expNoop {
					s.Require().Nil(block)
				} else {
					expBlock := s.buildFormattedBlock(
						blockRes,
						resBlock,
						tc.fullTx,
						tc.tx,
						tc.validator,
						tc.baseFee,
					)
					s.Require().Equal(expBlock, block)
				}
				s.Require().NoError(err)

			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestGetBlockByHash() {
	var (
		blockRes *cmtrpctypes.ResultBlockResults
		resBlock *cmtrpctypes.ResultBlock
	)
	msgEthereumTx, bz := s.buildEthereumTx()

	block := cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil)

	testCases := []struct {
		name         string
		hash         common.Hash
		fullTx       bool
		baseFee      *big.Int
		validator    sdk.AccAddress
		tx           *evmtypes.MsgEthereumTx
		txBz         []byte
		registerMock func(common.Hash, math.Int, sdk.AccAddress, []byte)
		expNoop      bool
		expPass      bool
	}{
		{
			"fail - tendermint failed to get block",
			common.BytesToHash(block.Hash()),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(hash common.Hash, _ math.Int, _ sdk.AccAddress, txBz []byte) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockByHashError(client, hash, txBz)
			},
			false,
			false,
		},
		{
			"fail - tendermint blockres not found",
			common.BytesToHash(block.Hash()),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(hash common.Hash, _ math.Int, _ sdk.AccAddress, txBz []byte) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockByHashNotFound(client, hash, txBz)
			},
			false,
			false,
		},
		{
			"noop - tendermint failed to fetch block result",
			common.BytesToHash(block.Hash()),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(hash common.Hash, _ math.Int, _ sdk.AccAddress, txBz []byte) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resBlock, _ = RegisterBlockByHash(client, hash, txBz)

				RegisterBlockResultsError(client, height)
			},
			true,
			true,
		},
		{
			"pass - without tx",
			common.BytesToHash(block.Hash()),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			nil,
			nil,
			func(hash common.Hash, baseFee math.Int, validator sdk.AccAddress, txBz []byte) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resBlock, _ = RegisterBlockByHash(client, hash, txBz)

				blockRes, _ = RegisterBlockResults(client, height)
				RegisterConsensusParams(client, height)

				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)
			},
			false,
			true,
		},
		{
			"pass - with tx",
			common.BytesToHash(block.Hash()),
			true,
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			msgEthereumTx,
			bz,
			func(hash common.Hash, baseFee math.Int, validator sdk.AccAddress, txBz []byte) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resBlock, _ = RegisterBlockByHash(client, hash, txBz)

				blockRes, _ = RegisterBlockResults(client, height)
				RegisterConsensusParams(client, height)

				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)
			},
			false,
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock(tc.hash, math.NewIntFromBigInt(tc.baseFee), tc.validator, tc.txBz)

			block, err := s.backend.GetBlockByHash(tc.hash, tc.fullTx)

			if tc.expPass {
				if tc.expNoop {
					s.Require().Nil(block)
				} else {
					expBlock := s.buildFormattedBlock(
						blockRes,
						resBlock,
						tc.fullTx,
						tc.tx,
						tc.validator,
						tc.baseFee,
					)
					s.Require().Equal(expBlock, block)
				}
				s.Require().NoError(err)

			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestGetBlockTransactionCountByHash() {
	_, bz := s.buildEthereumTx()
	block := cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil)
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{}, nil, nil)

	testCases := []struct {
		name         string
		hash         common.Hash
		registerMock func(common.Hash)
		expCount     hexutil.Uint
		expPass      bool
	}{
		{
			"fail - block not found",
			common.BytesToHash(emptyBlock.Hash()),
			func(hash common.Hash) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockByHashError(client, hash, nil)
			},
			hexutil.Uint(0),
			false,
		},
		{
			"fail - tendermint client failed to get block result",
			common.BytesToHash(emptyBlock.Hash()),
			func(hash common.Hash) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlockByHash(client, hash, nil)
				s.Require().NoError(err)
				RegisterBlockResultsError(client, height)
			},
			hexutil.Uint(0),
			false,
		},
		{
			"pass - block without tx",
			common.BytesToHash(emptyBlock.Hash()),
			func(hash common.Hash) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlockByHash(client, hash, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, height)
				s.Require().NoError(err)
			},
			hexutil.Uint(0),
			true,
		},
		{
			"pass - block with tx",
			common.BytesToHash(block.Hash()),
			func(hash common.Hash) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlockByHash(client, hash, bz)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, height)
				s.Require().NoError(err)
			},
			hexutil.Uint(1),
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			tc.registerMock(tc.hash)
			count := s.backend.GetBlockTransactionCountByHash(tc.hash)
			if tc.expPass {
				s.Require().Equal(tc.expCount, *count)
			} else {
				s.Require().Nil(count)
			}
		})
	}
}

func (s *TestSuite) TestGetBlockTransactionCountByNumber() {
	_, bz := s.buildEthereumTx()
	block := cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil)
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{}, nil, nil)

	testCases := []struct {
		name         string
		blockNum     ethrpc.BlockNumber
		registerMock func(ethrpc.BlockNumber)
		expCount     hexutil.Uint
		expPass      bool
	}{
		{
			"fail - block not found",
			ethrpc.BlockNumber(emptyBlock.Height),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, height)
			},
			hexutil.Uint(0),
			false,
		},
		{
			"fail - tendermint client failed to get block result",
			ethrpc.BlockNumber(emptyBlock.Height),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlock(client, height, nil)
				s.Require().NoError(err)
				RegisterBlockResultsError(client, height)
			},
			hexutil.Uint(0),
			false,
		},
		{
			"pass - block without tx",
			ethrpc.BlockNumber(emptyBlock.Height),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlock(client, height, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, height)
				s.Require().NoError(err)
			},
			hexutil.Uint(0),
			true,
		},
		{
			"pass - block with tx",
			ethrpc.BlockNumber(block.Height),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlock(client, height, bz)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, height)
				s.Require().NoError(err)
			},
			hexutil.Uint(1),
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			tc.registerMock(tc.blockNum)
			count := s.backend.GetBlockTransactionCountByNumber(tc.blockNum)
			if tc.expPass {
				s.Require().Equal(tc.expCount, *count)
			} else {
				s.Require().Nil(count)
			}
		})
	}
}

func (s *TestSuite) TestTendermintBlockByNumber() {
	var expResultBlock *cmtrpctypes.ResultBlock

	testCases := []struct {
		name         string
		blockNumber  ethrpc.BlockNumber
		registerMock func(ethrpc.BlockNumber)
		found        bool
		expPass      bool
	}{
		{
			"fail - client error",
			ethrpc.BlockNumber(1),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, height)
			},
			false,
			false,
		},
		{
			"noop - block not found",
			ethrpc.BlockNumber(1),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlockNotFound(client, height)
				s.Require().NoError(err)
			},
			false,
			true,
		},
		{
			"fail - blockNum < 0 with app state height error",
			ethrpc.BlockNumber(-1),
			func(_ ethrpc.BlockNumber) {
				var header metadata.MD
				appHeight := int64(1)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParamsError(QueryClient, &header, appHeight)
			},
			false,
			false,
		},
		{
			"pass - blockNum < 0 with app state height >= 1",
			ethrpc.BlockNumber(-1),
			func(ethrpc.BlockNumber) {
				var header metadata.MD
				appHeight := int64(1)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(QueryClient, &header, appHeight)

				tmHeight := appHeight
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultBlock, _ = RegisterBlock(client, tmHeight, nil)
			},
			true,
			true,
		},
		{
			"pass - blockNum = 0 (defaults to blockNum = 1 due to a difference between tendermint heights and geth heights)",
			ethrpc.BlockNumber(0),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultBlock, _ = RegisterBlock(client, height, nil)
			},
			true,
			true,
		},
		{
			"pass - blockNum = 1",
			ethrpc.BlockNumber(1),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultBlock, _ = RegisterBlock(client, height, nil)
			},
			true,
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			tc.registerMock(tc.blockNumber)
			resultBlock, err := s.backend.TendermintBlockByNumber(tc.blockNumber)

			if tc.expPass {
				s.Require().NoError(err)

				if !tc.found {
					s.Require().Nil(resultBlock)
				} else {
					s.Require().Equal(expResultBlock, resultBlock)
					s.Require().Equal(expResultBlock.Block.Header.Height, resultBlock.Block.Header.Height)
				}
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestTendermintBlockResultByNumber() {
	var expBlockRes *cmtrpctypes.ResultBlockResults

	testCases := []struct {
		name         string
		blockNumber  int64
		registerMock func(int64)
		expPass      bool
	}{
		{
			"fail",
			1,
			func(blockNum int64) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockResultsError(client, blockNum)
			},
			false,
		},
		{
			"pass",
			1,
			func(blockNum int64) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlockResults(client, blockNum)
				s.Require().NoError(err)
				expBlockRes = &cmtrpctypes.ResultBlockResults{
					Height:     blockNum,
					TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
				}
			},
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock(tc.blockNumber)

			client := s.backend.ClientCtx.Client.(*mocks.Client)
			blockRes, err := client.BlockResults(s.backend.Ctx, &tc.blockNumber) //#nosec G601 -- fine for tests

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().Equal(expBlockRes, blockRes)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestBlockNumberFromTendermint() {
	var resHeader *cmtrpctypes.ResultHeader

	_, bz := s.buildEthereumTx()
	block := cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil)
	blockNum := ethrpc.NewBlockNumber(big.NewInt(block.Height))
	blockHash := common.BytesToHash(block.Hash())

	testCases := []struct {
		name         string
		blockNum     *ethrpc.BlockNumber
		hash         *common.Hash
		registerMock func(*common.Hash)
		expPass      bool
	}{
		{
			"error - without blockHash or blockNum",
			nil,
			nil,
			func(*common.Hash) {},
			false,
		},
		{
			"error - with blockHash, tendermint client failed to get block",
			nil,
			&blockHash,
			func(hash *common.Hash) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterHeaderByHashError(client, *hash, bz)
			},
			false,
		},
		{
			"pass - with blockHash",
			nil,
			&blockHash,
			func(hash *common.Hash) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resHeader, _ = RegisterHeaderByHash(client, *hash, bz)
			},
			true,
		},
		{
			"pass - without blockHash & with blockNumber",
			&blockNum,
			nil,
			func(*common.Hash) {},
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			blockNrOrHash := ethrpc.BlockNumberOrHash{
				BlockNumber: tc.blockNum,
				BlockHash:   tc.hash,
			}

			tc.registerMock(tc.hash)
			blockNum, err := s.backend.BlockNumberFromTendermint(blockNrOrHash)

			if tc.expPass {
				s.Require().NoError(err)
				if tc.hash == nil {
					s.Require().Equal(*tc.blockNum, blockNum)
				} else {
					expHeight := ethrpc.NewBlockNumber(big.NewInt(resHeader.Header.Height))
					s.Require().Equal(expHeight, blockNum)
				}
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestBlockNumberFromTendermintByHash() {
	var resHeader *cmtrpctypes.ResultHeader

	_, bz := s.buildEthereumTx()
	block := cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil)
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{}, nil, nil)

	testCases := []struct {
		name         string
		hash         common.Hash
		registerMock func(common.Hash)
		expPass      bool
	}{
		{
			"fail - tendermint client failed to get block",
			common.BytesToHash(block.Hash()),
			func(hash common.Hash) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterHeaderByHashError(client, hash, bz)
			},
			false,
		},
		{
			"pass - block without tx",
			common.BytesToHash(emptyBlock.Hash()),
			func(hash common.Hash) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resHeader, _ = RegisterHeaderByHash(client, hash, bz)
			},
			true,
		},
		{
			"pass - block with tx",
			common.BytesToHash(block.Hash()),
			func(hash common.Hash) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				resHeader, _ = RegisterHeaderByHash(client, hash, bz)
			},
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			tc.registerMock(tc.hash)
			blockNum, err := s.backend.BlockNumberFromTendermintByHash(tc.hash)
			if tc.expPass {
				expHeight := big.NewInt(resHeader.Header.Height)
				s.Require().NoError(err)
				s.Require().Equal(expHeight, blockNum)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestBlockBloom() {
	testCases := []struct {
		name          string
		blockRes      *cmtrpctypes.ResultBlockResults
		expBlockBloom ethtypes.Bloom
		expPass       bool
	}{
		{
			"fail - empty block result",
			&cmtrpctypes.ResultBlockResults{},
			ethtypes.Bloom{},
			false,
		},
		{
			"fail - non block bloom event type",
			&cmtrpctypes.ResultBlockResults{
				FinalizeBlockEvents: []types.Event{{Type: evmtypes.EventTypeEthereumTx}},
			},
			ethtypes.Bloom{},
			false,
		},
		{
			"fail - nonblock bloom attribute key",
			&cmtrpctypes.ResultBlockResults{
				FinalizeBlockEvents: []types.Event{
					{
						Type: evmtypes.EventTypeBlockBloom,
						Attributes: []types.EventAttribute{
							{Key: evmtypes.AttributeKeyEthereumTxHash},
						},
					},
				},
			},
			ethtypes.Bloom{},
			false,
		},
		{
			"pass - block bloom attribute key",
			&cmtrpctypes.ResultBlockResults{
				FinalizeBlockEvents: []types.Event{
					{
						Type: evmtypes.EventTypeBlockBloom,
						Attributes: []types.EventAttribute{
							{Key: evmtypes.AttributeKeyEthereumBloom},
						},
					},
				},
			},
			ethtypes.Bloom{},
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			blockBloom, err := s.backend.BlockBloom(tc.blockRes)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().Equal(tc.expBlockBloom, blockBloom)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestGetEthBlockFromTendermint() {
	msgEthereumTx, bz := s.buildEthereumTx()
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{}, nil, nil)

	testCases := []struct {
		name         string
		baseFee      *big.Int
		validator    sdk.AccAddress
		height       int64
		resBlock     *cmtrpctypes.ResultBlock
		blockRes     *cmtrpctypes.ResultBlockResults
		fullTx       bool
		registerMock func(math.Int, sdk.AccAddress, int64)
		expTxs       bool
		expPass      bool
	}{
		{
			"pass - block without tx",
			math.NewInt(1).BigInt(),
			sdk.AccAddress(common.Address{}.Bytes()),
			int64(1),
			&cmtrpctypes.ResultBlock{Block: emptyBlock},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
			},
			false,
			func(baseFee math.Int, validator sdk.AccAddress, height int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterConsensusParams(client, height)
			},
			false,
			true,
		},
		{
			"pass - block with tx - with BaseFee error",
			nil,
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			int64(1),
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
			},
			true,
			func(_ math.Int, validator sdk.AccAddress, height int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFeeError(QueryClient)
				RegisterValidatorAccount(QueryClient, validator)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterConsensusParams(client, height)
			},
			true,
			true,
		},
		{
			"pass - block with tx - with ValidatorAccount error",
			math.NewInt(1).BigInt(),
			sdk.AccAddress(common.Address{}.Bytes()),
			int64(1),
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
			},
			true,
			func(baseFee math.Int, _ sdk.AccAddress, height int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccountError(QueryClient)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterConsensusParams(client, height)
			},
			true,
			true,
		},
		{
			"pass - block with tx - with ConsensusParams error - BlockMaxGas defaults to max uint32",
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			int64(1),
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
			},
			true,
			func(baseFee math.Int, validator sdk.AccAddress, height int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterConsensusParamsError(client, height)
			},
			true,
			true,
		},
		{
			"pass - block with tx - with ShouldIgnoreGasUsed - empty txs",
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			int64(1),
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				Height: 1,
				TxsResults: []*types.ExecTxResult{
					{
						Code:    11,
						GasUsed: 0,
						Log:     "no block gas left to run tx: out of gas",
					},
				},
			},
			true,
			func(baseFee math.Int, validator sdk.AccAddress, height int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterConsensusParams(client, height)
			},
			false,
			true,
		},
		{
			"pass - block with tx - non fullTx",
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			int64(1),
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
			},
			false,
			func(baseFee math.Int, validator sdk.AccAddress, height int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterConsensusParams(client, height)
			},
			true,
			true,
		},
		{
			"pass - block with tx",
			math.NewInt(1).BigInt(),
			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
			int64(1),
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
			},
			true,
			func(baseFee math.Int, validator sdk.AccAddress, height int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterValidatorAccount(QueryClient, validator)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterConsensusParams(client, height)
			},
			true,
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock(math.NewIntFromBigInt(tc.baseFee), tc.validator, tc.height)

			block, err := s.backend.RPCBlockFromTendermintBlock(tc.resBlock, tc.blockRes, tc.fullTx)

			var expBlock map[string]interface{}
			header := tc.resBlock.Block.Header
			gasLimit := int64(
				^uint32(0),
			) // for `MaxGas = -1` (DefaultConsensusParams)
			gasUsed := new(
				big.Int,
			).SetUint64(uint64(tc.blockRes.TxsResults[0].GasUsed)) //#nosec G115 won't exceed uint64

			root := common.Hash{}.Bytes()
			receipt := ethtypes.NewReceipt(root, false, gasUsed.Uint64())
			bloom := ethtypes.CreateBloom(receipt)

			ethRPCTxs := []interface{}{}

			if tc.expTxs {
				if tc.fullTx {
					rpcTx, err := ethrpc.NewRPCTransaction(
						msgEthereumTx,
						common.BytesToHash(header.Hash()),
						uint64(header.Height), //#nosec G115 won't exceed uint64
						uint64(0),
						tc.baseFee,
						s.backend.EvmChainID,
					)
					s.Require().NoError(err)
					ethRPCTxs = []interface{}{rpcTx}
				} else {
					ethRPCTxs = []interface{}{common.HexToHash(msgEthereumTx.Hash)}
				}
			}

			expBlock = ethrpc.FormatBlock(
				header,
				tc.resBlock.Block.Size(),
				gasLimit,
				gasUsed,
				ethRPCTxs,
				bloom,
				common.BytesToAddress(tc.validator.Bytes()),
				tc.baseFee,
			)

			if tc.expPass {
				s.Require().Equal(expBlock, block)
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestEthMsgsFromTendermintBlock() {
	msgEthereumTx, bz := s.buildEthereumTx()

	testCases := []struct {
		name     string
		resBlock *cmtrpctypes.ResultBlock
		blockRes *cmtrpctypes.ResultBlockResults
		expMsgs  []*evmtypes.MsgEthereumTx
	}{
		{
			"tx in not included in block - unsuccessful tx without ExceedBlockGasLimit error",
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				TxsResults: []*types.ExecTxResult{
					{
						Code: 1,
					},
				},
			},
			[]*evmtypes.MsgEthereumTx(nil),
		},
		{
			"tx included in block - unsuccessful tx with ExceedBlockGasLimit error",
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				TxsResults: []*types.ExecTxResult{
					{
						Code: 1,
						Log:  ethrpc.ExceedBlockGasLimitError,
					},
				},
			},
			[]*evmtypes.MsgEthereumTx{msgEthereumTx},
		},
		{
			"pass",
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				TxsResults: []*types.ExecTxResult{
					{
						Code: 0,
						Log:  ethrpc.ExceedBlockGasLimitError,
					},
				},
			},
			[]*evmtypes.MsgEthereumTx{msgEthereumTx},
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			msgs, _ := s.backend.EthMsgsFromTendermintBlock(tc.resBlock, tc.blockRes)
			s.Require().Equal(tc.expMsgs, msgs)
		})
	}
}

func (s *TestSuite) TestHeaderByNumber() {
	var expResultBlock *cmtrpctypes.ResultBlock

	_, bz := s.buildEthereumTx()

	testCases := []struct {
		name         string
		blockNumber  ethrpc.BlockNumber
		baseFee      *big.Int
		registerMock func(ethrpc.BlockNumber, math.Int)
		expPass      bool
	}{
		{
			"fail - tendermint client failed to get block",
			ethrpc.BlockNumber(1),
			math.NewInt(1).BigInt(),
			func(blockNum ethrpc.BlockNumber, _ math.Int) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, height)
			},
			false,
		},
		{
			"fail - block not found for height",
			ethrpc.BlockNumber(1),
			math.NewInt(1).BigInt(),
			func(blockNum ethrpc.BlockNumber, _ math.Int) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlockNotFound(client, height)
				s.Require().NoError(err)
			},
			false,
		},
		{
			"fail - block not found for height",
			ethrpc.BlockNumber(1),
			math.NewInt(1).BigInt(),
			func(blockNum ethrpc.BlockNumber, _ math.Int) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlock(client, height, nil)
				s.Require().NoError(err)
				RegisterBlockResultsError(client, height)
			},
			false,
		},
		{
			"pass - without Base Fee, failed to fetch from prunned block",
			ethrpc.BlockNumber(1),
			nil,
			func(blockNum ethrpc.BlockNumber, _ math.Int) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultBlock, _ = RegisterBlock(client, height, nil)
				_, err := RegisterBlockResults(client, height)
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFeeError(QueryClient)
			},
			true,
		},
		{
			"pass - blockNum = 1, without tx",
			ethrpc.BlockNumber(1),
			math.NewInt(1).BigInt(),
			func(blockNum ethrpc.BlockNumber, baseFee math.Int) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultBlock, _ = RegisterBlock(client, height, nil)
				_, err := RegisterBlockResults(client, height)
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
			},
			true,
		},
		{
			"pass - blockNum = 1, with tx",
			ethrpc.BlockNumber(1),
			math.NewInt(1).BigInt(),
			func(blockNum ethrpc.BlockNumber, baseFee math.Int) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultBlock, _ = RegisterBlock(client, height, bz)
				_, err := RegisterBlockResults(client, height)
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
			},
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			tc.registerMock(tc.blockNumber, math.NewIntFromBigInt(tc.baseFee))
			header, err := s.backend.HeaderByNumber(tc.blockNumber)

			if tc.expPass {
				expHeader := ethrpc.EthHeaderFromTendermint(expResultBlock.Block.Header, ethtypes.Bloom{}, tc.baseFee)
				s.Require().NoError(err)
				s.Require().Equal(expHeader, header)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestHeaderByHash() {
	var expResultHeader *cmtrpctypes.ResultHeader

	_, bz := s.buildEthereumTx()
	block := cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil)
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{}, nil, nil)

	testCases := []struct {
		name         string
		hash         common.Hash
		baseFee      *big.Int
		registerMock func(common.Hash, math.Int)
		expPass      bool
	}{
		{
			"fail - tendermint client failed to get block",
			common.BytesToHash(block.Hash()),
			math.NewInt(1).BigInt(),
			func(hash common.Hash, _ math.Int) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterHeaderByHashError(client, hash, bz)
			},
			false,
		},
		{
			"fail - block not found for height",
			common.BytesToHash(block.Hash()),
			math.NewInt(1).BigInt(),
			func(hash common.Hash, _ math.Int) {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterHeaderByHashNotFound(client, hash, bz)
			},
			false,
		},
		{
			"fail - block not found for height",
			common.BytesToHash(block.Hash()),
			math.NewInt(1).BigInt(),
			func(hash common.Hash, _ math.Int) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterHeaderByHash(client, hash, bz)
				s.Require().NoError(err)
				RegisterBlockResultsError(client, height)
			},
			false,
		},
		{
			"pass - without Base Fee, failed to fetch from prunned block",
			common.BytesToHash(block.Hash()),
			nil,
			func(hash common.Hash, _ math.Int) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultHeader, _ = RegisterHeaderByHash(client, hash, bz)
				_, err := RegisterBlockResults(client, height)
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFeeError(QueryClient)
			},
			true,
		},
		{
			"pass - blockNum = 1, without tx",
			common.BytesToHash(emptyBlock.Hash()),
			math.NewInt(1).BigInt(),
			func(hash common.Hash, baseFee math.Int) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultHeader, _ = RegisterHeaderByHash(client, hash, nil)
				_, err := RegisterBlockResults(client, height)
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
			},
			true,
		},
		{
			"pass - with tx",
			common.BytesToHash(block.Hash()),
			math.NewInt(1).BigInt(),
			func(hash common.Hash, baseFee math.Int) {
				height := int64(1)
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				expResultHeader, _ = RegisterHeaderByHash(client, hash, bz)
				_, err := RegisterBlockResults(client, height)
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
			},
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries

			tc.registerMock(tc.hash, math.NewIntFromBigInt(tc.baseFee))
			header, err := s.backend.HeaderByHash(tc.hash)

			if tc.expPass {
				expHeader := ethrpc.EthHeaderFromTendermint(*expResultHeader.Header, ethtypes.Bloom{}, tc.baseFee)
				s.Require().NoError(err)
				s.Require().Equal(expHeader, header)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestEthBlockByNumber() {
	msgEthereumTx, bz := s.buildEthereumTx()
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{}, nil, nil)

	testCases := []struct {
		name         string
		blockNumber  ethrpc.BlockNumber
		registerMock func(ethrpc.BlockNumber)
		expEthBlock  *ethtypes.Block
		expPass      bool
	}{
		{
			"fail - tendermint client failed to get block",
			ethrpc.BlockNumber(1),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				RegisterBlockError(client, height)
			},
			nil,
			false,
		},
		{
			"fail - block result not found for height",
			ethrpc.BlockNumber(1),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlock(client, height, nil)
				s.Require().NoError(err)
				RegisterBlockResultsError(client, blockNum.Int64())
			},
			nil,
			false,
		},
		{
			"pass - block without tx",
			ethrpc.BlockNumber(1),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlock(client, height, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, blockNum.Int64())
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				baseFee := math.NewInt(1)
				RegisterBaseFee(QueryClient, baseFee)
			},
			ethtypes.NewBlock(
				ethrpc.EthHeaderFromTendermint(
					emptyBlock.Header,
					ethtypes.Bloom{},
					math.NewInt(1).BigInt(),
				),
				&ethtypes.Body{},
				nil,
				nil,
			),
			true,
		},
		{
			"pass - block with tx",
			ethrpc.BlockNumber(1),
			func(blockNum ethrpc.BlockNumber) {
				height := blockNum.Int64()
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				_, err := RegisterBlock(client, height, bz)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, blockNum.Int64())
				s.Require().NoError(err)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				baseFee := math.NewInt(1)
				RegisterBaseFee(QueryClient, baseFee)
			},
			ethtypes.NewBlock(
				ethrpc.EthHeaderFromTendermint(
					emptyBlock.Header,
					ethtypes.Bloom{},
					math.NewInt(1).BigInt(),
				),
				&ethtypes.Body{
					Transactions: []*ethtypes.Transaction{msgEthereumTx.AsTransaction()},
				},
				nil,
				trie.NewStackTrie(nil),
			),
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock(tc.blockNumber)

			ethBlock, err := s.backend.EthBlockByNumber(tc.blockNumber)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().Equal(tc.expEthBlock.Header(), ethBlock.Header())
				s.Require().Equal(tc.expEthBlock.Uncles(), ethBlock.Uncles())
				s.Require().Equal(tc.expEthBlock.ReceiptHash(), ethBlock.ReceiptHash())
				for i, tx := range tc.expEthBlock.Transactions() {
					s.Require().Equal(tx.Data(), ethBlock.Transactions()[i].Data())
				}

			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestEthBlockFromTendermintBlock() {
	msgEthereumTx, bz := s.buildEthereumTx()
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{}, nil, nil)

	testCases := []struct {
		name         string
		baseFee      *big.Int
		resBlock     *cmtrpctypes.ResultBlock
		blockRes     *cmtrpctypes.ResultBlockResults
		registerMock func(math.Int, int64)
		expEthBlock  *ethtypes.Block
		expPass      bool
	}{
		{
			"pass - block without tx",
			math.NewInt(1).BigInt(),
			&cmtrpctypes.ResultBlock{
				Block: emptyBlock,
			},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
			},
			func(baseFee math.Int, _ int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
			},
			ethtypes.NewBlock(
				ethrpc.EthHeaderFromTendermint(
					emptyBlock.Header,
					ethtypes.Bloom{},
					math.NewInt(1).BigInt(),
				),
				&ethtypes.Body{},
				nil,
				nil,
			),
			true,
		},
		{
			"pass - block with tx",
			math.NewInt(1).BigInt(),
			&cmtrpctypes.ResultBlock{
				Block: cmttypes.MakeBlock(1, []cmttypes.Tx{bz}, nil, nil),
			},
			&cmtrpctypes.ResultBlockResults{
				Height:     1,
				TxsResults: []*types.ExecTxResult{{Code: 0, GasUsed: 0}},
				FinalizeBlockEvents: []types.Event{
					{
						Type: evmtypes.EventTypeBlockBloom,
						Attributes: []types.EventAttribute{
							{Key: evmtypes.AttributeKeyEthereumBloom},
						},
					},
				},
			},
			func(baseFee math.Int, _ int64) {
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterBaseFee(QueryClient, baseFee)
			},
			ethtypes.NewBlock(
				ethrpc.EthHeaderFromTendermint(
					emptyBlock.Header,
					ethtypes.Bloom{},
					math.NewInt(1).BigInt(),
				),
				&ethtypes.Body{Transactions: []*ethtypes.Transaction{msgEthereumTx.AsTransaction()}},
				nil,
				trie.NewStackTrie(nil),
			),
			true,
		},
	}
	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock(math.NewIntFromBigInt(tc.baseFee), tc.blockRes.Height)

			ethBlock, err := s.backend.EthBlockFromTendermintBlock(tc.resBlock, tc.blockRes)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().Equal(tc.expEthBlock.Header(), ethBlock.Header())
				s.Require().Equal(tc.expEthBlock.Uncles(), ethBlock.Uncles())
				s.Require().Equal(tc.expEthBlock.ReceiptHash(), ethBlock.ReceiptHash())
				for i, tx := range tc.expEthBlock.Transactions() {
					s.Require().Equal(tx.Data(), ethBlock.Transactions()[i].Data())
				}

			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (suite *TestSuite) TestEthAndSyntheticMsgsFromTendermintBlock() {
	suite.SetupTest() // reset test and queries
	// synthetic tx
	hash := sample.Hash().Hex()
	tx, txRes := suite.buildSyntheticTxResult(hash)

	// real tx
	msgEthereumTx, realTx := suite.buildEthereumTx()

	suite.backend.Indexer = nil
	// block contains block real and synthetic tx
	emptyBlock := cmttypes.MakeBlock(1, []cmttypes.Tx{realTx, tx}, nil, nil)
	emptyBlock.ChainID = ChainID.ChainID
	blockHash := common.BigToHash(big.NewInt(1)).Bytes()
	resBlock := &tmrpctypes.ResultBlock{Block: emptyBlock, BlockID: cmttypes.BlockID{Hash: bytes.HexBytes(blockHash)}}
	blockRes := &tmrpctypes.ResultBlockResults{
		Height:     1,
		TxsResults: []*types.ExecTxResult{{}, &txRes},
	}

	// both real and synthetic should be returned
	msgs, additionals := suite.backend.EthMsgsFromTendermintBlock(resBlock, blockRes)
	suite.Require().Equal(2, len(msgs))
	suite.Require().Equal(2, len(additionals))
	suite.Require().Nil(additionals[0])
	suite.Require().NotNil(additionals[1])
	suite.Require().Equal(msgEthereumTx.Hash, msgs[0].Hash)
	suite.Require().Equal(hash, msgs[1].Hash)
}

func (suite *TestSuite) TestEthAndSyntheticEthBlockByNumber() {
	suite.SetupTest() // reset test and queries
	// synthetic tx
	hash := sample.Hash().Hex()
	tx, txRes := suite.buildSyntheticTxResult(hash)

	// real tx
	msgEthereumTx, realTx := suite.buildEthereumTx()

	suite.backend.Indexer = nil
	client := suite.backend.ClientCtx.Client.(*mocks.Client)
	queryClient := suite.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)

	// block contains block real and synthetic tx
	RegisterBlockMultipleTxs(client, 1, []cmttypes.Tx{realTx, tx})
	RegisterBlockResultsWithTxResults(client, 1, []*types.ExecTxResult{{}, &txRes})
	RegisterBaseFee(queryClient, sdkmath.NewInt(1))

	// only real should be returned
	block, err := suite.backend.EthBlockByNumber(1)
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(block.Transactions()))
	suite.Require().Equal(msgEthereumTx.Hash, block.Transactions()[0].Hash().String())
}

func (suite *TestSuite) TestEthAndSyntheticGetBlockByNumber() {
	suite.SetupTest() // reset test and queries
	// synthetic tx
	hash := sample.Hash().Hex()
	tx, txRes := suite.buildSyntheticTxResult(hash)

	// real tx
	msgEthereumTx, realTx := suite.buildEthereumTx()

	suite.backend.Indexer = nil
	client := suite.backend.ClientCtx.Client.(*mocks.Client)
	queryClient := suite.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)

	// block contains block real and synthetic tx
	RegisterBlockMultipleTxs(client, 1, []cmttypes.Tx{realTx, tx})
	RegisterBlockResultsWithTxResults(client, 1, []*types.ExecTxResult{{}, &txRes})
	RegisterBaseFee(queryClient, sdkmath.NewInt(1))
	RegisterValidatorAccount(queryClient, sdk.AccAddress(common.Address{}.Bytes()))
	RegisterConsensusParams(client, 1)

	// both real and synthetic should be returned
	block, err := suite.backend.GetBlockByNumber(1, false)
	suite.Require().NoError(err)
	transactions := block["transactions"].([]interface{})
	suite.Require().Equal(2, len(transactions))
	suite.Require().Equal(common.HexToHash(msgEthereumTx.Hash), transactions[0])
	suite.Require().Equal(common.HexToHash(hash), transactions[1])

	// both real and synthetic should be returned
	block, err = suite.backend.GetBlockByNumber(1, true)
	suite.Require().NoError(err)
	transactions = block["transactions"].([]interface{})

	suite.Require().Equal(2, len(transactions))
	resRealTx := transactions[0].(*ethrpc.RPCTransaction)
	suite.Require().Equal(common.HexToHash(msgEthereumTx.Hash), resRealTx.Hash)

	resSyntheticTx := transactions[1].(*ethrpc.RPCTransaction)
	suite.Require().Equal(common.HexToHash(hash), resSyntheticTx.Hash)
	suite.Require().Equal(hash, resSyntheticTx.Hash.Hex())
	suite.Require().Equal("0x735b14BB79463307AAcBED86DAf3322B1e6226aB", resSyntheticTx.From.Hex())
	suite.Require().Equal("0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7", resSyntheticTx.To.Hex())
	suite.Require().Equal("0x58", resSyntheticTx.Type.String())
	suite.Require().Equal("0x1", resSyntheticTx.Nonce.String())
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), resSyntheticTx.V)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), resSyntheticTx.R)
	suite.Require().Equal((*hexutil.Big)(big.NewInt(0)), resSyntheticTx.S)
}

func (s *TestSuite) TestDecodeMsgEthereumTxFromCosmosMsg() {
	// Create an evm transaction
	msgEthereumTx, _ := s.buildEthereumTx()
	expectedTx := msgEthereumTx
	expectedTx.Hash = expectedTx.AsTransaction().Hash().Hex()

	// Create a legacy ethermint transaction for testing conversion
	legacyTx := s.buildLegacyEthereumTx()

	testCases := []struct {
		name     string
		msg      sdk.Msg
		chainID  *big.Int
		expTx    *evmtypes.MsgEthereumTx
		expError bool
		errorMsg string
	}{
		{
			"pass - evmtypes.MsgEthereumTx",
			msgEthereumTx,
			big.NewInt(1),
			expectedTx,
			false,
			"",
		},
		{
			"pass - etherminttypes.MsgEthereumTx (legacy conversion)",
			legacyTx,
			big.NewInt(1),
			nil,
			false,
			"",
		},
		{
			"fail - unsupported message type",
			func() sdk.Msg {
				// Create a simple message that doesn't implement MsgEthereumTx
				return simpleMsg{}
			}(),
			big.NewInt(1),
			nil,
			true,
			"can't cast to MsgEthereumTx",
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			result, err := DecodeMsgEthereumTxFromCosmosMsg(tc.msg, tc.chainID)

			if tc.expError {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errorMsg)
				s.Require().Nil(result)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(result)

				s.Require().NotEmpty(result.Hash)
				s.Require().NotNil(result.From)
			}
		})
	}
}
