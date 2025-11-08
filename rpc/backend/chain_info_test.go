package backend

// import (
// 	"fmt"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"github.com/ethereum/go-ethereum/common/math"
// 	ethrpc "github.com/ethereum/go-ethereum/rpc"
// 	"google.golang.org/grpc/metadata"

// 	"github.com/cometbft/cometbft/abci/types"
// 	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"

// 	utiltx "github.com/cosmos/evm/testutil/tx"
// 	evmtypes "github.com/cosmos/evm/x/vm/types"
// 	"github.com/zeta-chain/node/rpc/backend/mocks"
// 	rpc "github.com/zeta-chain/node/rpc/types"

// 	sdkmath "cosmossdk.io/math"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (s *TestSuite) TestBaseFee() {
// 	baseFee := sdkmath.NewInt(1)

// 	testCases := []struct {
// 		name         string
// 		blockRes     *tmrpctypes.ResultBlockResults
// 		registerMock func()
// 		expBaseFee   *big.Int
// 		expPass      bool
// 	}{
// 		{
// 			"fail - grpc BaseFee error",
// 			&tmrpctypes.ResultBlockResults{Height: 1},
// 			func() {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBaseFeeError(QueryClient)
// 			},
// 			nil,
// 			false,
// 		},
// 		{
// 			"fail - grpc BaseFee error - with non feemarket block event",
// 			&tmrpctypes.ResultBlockResults{
// 				Height: 1,
// 				FinalizeBlockEvents: []types.Event{
// 					{
// 						Type: evmtypes.EventTypeBlockBloom,
// 					},
// 				},
// 			},
// 			func() {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBaseFeeError(QueryClient)
// 			},
// 			nil,
// 			false,
// 		},
// 		{
// 			"fail - grpc BaseFee error - with feemarket block event",
// 			&tmrpctypes.ResultBlockResults{
// 				Height: 1,
// 				FinalizeBlockEvents: []types.Event{
// 					{
// 						Type: evmtypes.EventTypeFeeMarket,
// 					},
// 				},
// 			},
// 			func() {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBaseFeeError(QueryClient)
// 			},
// 			nil,
// 			false,
// 		},
// 		{
// 			"fail - grpc BaseFee error - with feemarket block event with wrong attribute value",
// 			&tmrpctypes.ResultBlockResults{
// 				Height: 1,
// 				FinalizeBlockEvents: []types.Event{
// 					{
// 						Type: evmtypes.EventTypeFeeMarket,
// 						Attributes: []types.EventAttribute{
// 							{Value: "/1"},
// 						},
// 					},
// 				},
// 			},
// 			func() {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBaseFeeError(QueryClient)
// 			},
// 			nil,
// 			false,
// 		},
// 		{
// 			"fail - grpc baseFee error - with feemarket block event with baseFee attribute value",
// 			&tmrpctypes.ResultBlockResults{
// 				Height: 1,
// 				FinalizeBlockEvents: []types.Event{
// 					{
// 						Type: evmtypes.EventTypeFeeMarket,
// 						Attributes: []types.EventAttribute{
// 							{Value: baseFee.String()},
// 						},
// 					},
// 				},
// 			},
// 			func() {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBaseFeeError(QueryClient)
// 			},
// 			baseFee.BigInt(),
// 			true,
// 		},
// 		{
// 			"fail - base fee or london fork not enabled",
// 			&tmrpctypes.ResultBlockResults{Height: 1},
// 			func() {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBaseFeeDisabled(QueryClient)
// 			},
// 			nil,
// 			true,
// 		},
// 		{
// 			"pass",
// 			&tmrpctypes.ResultBlockResults{Height: 1},
// 			func() {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBaseFee(QueryClient, baseFee)
// 			},
// 			baseFee.BigInt(),
// 			true,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			baseFee, err := s.backend.BaseFee(tc.blockRes)

// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expBaseFee, baseFee)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestChainID() {
// 	expChainID := (*hexutil.Big)(big.NewInt(int64(ChainID.EVMChainID)))
// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		expChainID   *hexutil.Big
// 		expPass      bool
// 	}{
// 		{
// 			"pass - block is at or past the EIP-155 replay-protection fork block, return chainID from config ",
// 			func() {
// 				var header metadata.MD
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterParamsInvalidHeight(QueryClient, &header, int64(1))
// 			},
// 			expChainID,
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			chainID, err := s.backend.ChainID()
// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expChainID, chainID)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestGetCoinbase() {
// 	validatorAcc := sdk.AccAddress(utiltx.GenerateAddress().Bytes())
// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		accAddr      sdk.AccAddress
// 		expPass      bool
// 	}{
// 		{
// 			"fail - Can't retrieve status from node",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				RegisterStatusError(client)
// 			},
// 			validatorAcc,
// 			false,
// 		},
// 		{
// 			"fail - Can't query validator account",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterStatus(client)
// 				RegisterValidatorAccountError(QueryClient)
// 			},
// 			validatorAcc,
// 			false,
// 		},
// 		{
// 			"pass - Gets coinbase account",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterStatus(client)
// 				RegisterValidatorAccount(QueryClient, validatorAcc)
// 			},
// 			validatorAcc,
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			accAddr, err := s.backend.GetCoinbase()

// 			if tc.expPass {
// 				s.Require().Equal(tc.accAddr, accAddr)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestSuggestGasTipCap() {
// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		baseFee      *big.Int
// 		expGasTipCap *big.Int
// 		expPass      bool
// 	}{
// 		{
// 			"pass - London hardfork not enabled or feemarket not enabled ",
// 			func() {},
// 			nil,
// 			big.NewInt(0),
// 			true,
// 		},
// 		{
// 			"pass - Gets the suggest gas tip cap ",
// 			func() {},
// 			nil,
// 			big.NewInt(0),
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			maxDelta, err := s.backend.SuggestGasTipCap(tc.baseFee)

// 			if tc.expPass {
// 				s.Require().Equal(tc.expGasTipCap, maxDelta)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestGlobalMinGasPrice() {
// 	testCases := []struct {
// 		name           string
// 		registerMock   func()
// 		expMinGasPrice *big.Int
// 		expPass        bool
// 	}{
// 		{
// 			"pass - get GlobalMinGasPrice",
// 			func() {
// 				qc := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterGlobalMinGasPrice(qc, 1)
// 			},
// 			big.NewInt(1),
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			globalMinGasPrice, err := s.backend.GlobalMinGasPrice()

// 			if tc.expPass {
// 				s.Require().Equal(tc.expMinGasPrice, globalMinGasPrice)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestFeeHistory() {
// 	testCases := []struct {
// 		name              string
// 		registerMock      func(validator sdk.AccAddress)
// 		userBlockCount    math.HexOrDecimal64
// 		latestBlock       ethrpc.BlockNumber
// 		expFeeHistory     *rpc.FeeHistoryResult
// 		validator         sdk.AccAddress
// 		expPass           bool
// 		targetNewBaseFees []*big.Int
// 	}{
// 		{
// 			"fail - can't get params ",
// 			func(_ sdk.AccAddress) {
// 				var header metadata.MD
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 0
// 				RegisterParamsError(QueryClient, &header, ethrpc.BlockNumber(1).Int64())
// 			},
// 			1,
// 			-1,
// 			nil,
// 			nil,
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - user block count higher than max block count ",
// 			func(_ sdk.AccAddress) {
// 				var header metadata.MD
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 0
// 				RegisterParams(QueryClient, &header, ethrpc.BlockNumber(1).Int64())
// 			},
// 			1,
// 			-1,
// 			nil,
// 			nil,
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - Tendermint block fetching error ",
// 			func(_ sdk.AccAddress) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				var header metadata.MD
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterParams(queryClient, &header, 1)
// 				RegisterBlockError(client, ethrpc.BlockNumber(1).Int64())
// 			},
// 			1,
// 			1,
// 			nil,
// 			nil,
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - Tendermint block fetching panic",
// 			func(_ sdk.AccAddress) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				var header metadata.MD
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterParams(queryClient, &header, 1)
// 				RegisterBlockPanic(client, ethrpc.BlockNumber(1).Int64())
// 			},
// 			1,
// 			1,
// 			nil,
// 			nil,
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - Eth block fetching error",
// 			func(_ sdk.AccAddress) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				var header metadata.MD
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterParams(queryClient, &header, 1)
// 				_, err := RegisterBlock(client, ethrpc.BlockNumber(1).Int64(), nil)
// 				s.Require().NoError(err)
// 				RegisterBlockResultsError(client, 1)
// 			},
// 			1,
// 			1,
// 			nil,
// 			nil,
// 			true,
// 			nil,
// 		},
// 		{
// 			"pass - skip invalid base fee",
// 			func(validator sdk.AccAddress) {
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				var header metadata.MD
// 				RegisterParams(queryClient, &header, 1)
// 				_, err := RegisterBlock(client, ethrpc.BlockNumber(1).Int64(), nil)
// 				s.Require().NoError(err)
// 				_, err = RegisterBlockResults(client, 1)
// 				s.Require().NoError(err)
// 				RegisterBaseFeeError(queryClient)
// 				RegisterValidatorAccount(queryClient, validator)
// 				RegisterConsensusParams(client, 1)
// 				fQueryClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
// 				RegisterFeeMarketParams(fQueryClient, 1)
// 			},
// 			1,
// 			1,
// 			&rpc.FeeHistoryResult{
// 				OldestBlock:  (*hexutil.Big)(big.NewInt(1)),
// 				BaseFee:      []*hexutil.Big{(*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(new(big.Int).SetBits([]big.Word{}))},
// 				GasUsedRatio: []float64{0},
// 				Reward:       [][]*hexutil.Big{{(*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0))}},
// 			},
// 			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
// 			true,
// 			nil,
// 		},
// 		{
// 			"pass - Valid FeeHistoryResults object",
// 			func(validator sdk.AccAddress) {
// 				var header metadata.MD
// 				baseFee := sdkmath.NewInt(1)
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				fQueryClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				_, err := RegisterBlock(client, ethrpc.BlockNumber(1).Int64(), nil)
// 				s.Require().NoError(err)
// 				_, err = RegisterBlockResults(client, 1)
// 				s.Require().NoError(err)
// 				RegisterBaseFee(queryClient, baseFee)
// 				RegisterValidatorAccount(queryClient, validator)
// 				RegisterConsensusParams(client, 1)
// 				RegisterParams(queryClient, &header, 1)
// 				RegisterFeeMarketParams(fQueryClient, 1)
// 			},
// 			1,
// 			1,
// 			&rpc.FeeHistoryResult{
// 				OldestBlock:  (*hexutil.Big)(big.NewInt(1)),
// 				BaseFee:      []*hexutil.Big{(*hexutil.Big)(big.NewInt(1)), (*hexutil.Big)(big.NewInt(1))},
// 				GasUsedRatio: []float64{0},
// 				Reward:       [][]*hexutil.Big{{(*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0))}},
// 			},
// 			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
// 			true,
// 			nil,
// 		},
// 		{
// 			"pass - Concurrent FeeHistoryResults object",
// 			func(validator sdk.AccAddress) {
// 				var header metadata.MD
// 				baseFee := sdkmath.NewInt(1)
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				fQueryClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				_, err := RegisterBlock(client, ethrpc.BlockNumber(1).Int64(), nil)
// 				s.Require().NoError(err)
// 				_, err = RegisterBlockResults(client, 1)
// 				s.Require().NoError(err)
// 				RegisterBaseFee(queryClient, baseFee)
// 				RegisterValidatorAccount(queryClient, validator)
// 				RegisterConsensusParams(client, 1)
// 				RegisterParams(queryClient, &header, 1)
// 				RegisterFeeMarketParams(fQueryClient, 1)
// 			},
// 			1,
// 			1,
// 			&rpc.FeeHistoryResult{
// 				OldestBlock:  (*hexutil.Big)(big.NewInt(1)),
// 				BaseFee:      []*hexutil.Big{(*hexutil.Big)(big.NewInt(1)), (*hexutil.Big)(big.NewInt(0))},
// 				GasUsedRatio: []float64{0},
// 				Reward:       [][]*hexutil.Big{{(*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0))}},
// 			},
// 			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
// 			true,
// 			[]*big.Int{
// 				big.NewInt(0), // for overwrite overlap
// 			},
// 		},
// 		{
// 			"pass - EarliestBlockNumber(0x0)",
// 			func(validator sdk.AccAddress) {
// 				var header metadata.MD
// 				baseFee := sdkmath.NewInt(1)
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				fQueryClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				_, err := RegisterBlock(client, ethrpc.BlockNumber(1).Int64(), nil)
// 				s.Require().NoError(err)
// 				_, err = RegisterBlockResults(client, 1)
// 				s.Require().NoError(err)
// 				RegisterBaseFee(queryClient, baseFee)
// 				RegisterValidatorAccount(queryClient, validator)
// 				RegisterConsensusParams(client, 1)
// 				RegisterParams(queryClient, &header, 1)
// 				RegisterFeeMarketParams(fQueryClient, 1)
// 			},
// 			1,
// 			0,
// 			&rpc.FeeHistoryResult{
// 				OldestBlock:  (*hexutil.Big)(big.NewInt(0)),
// 				BaseFee:      []*hexutil.Big{(*hexutil.Big)(big.NewInt(1)), (*hexutil.Big)(big.NewInt(1))},
// 				GasUsedRatio: []float64{0},
// 				Reward:       [][]*hexutil.Big{{(*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0))}},
// 			},
// 			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
// 			true,
// 			nil,
// 		},
// 		{
// 			"pass - EarliestBlockNumber(tag)",
// 			func(validator sdk.AccAddress) {
// 				var header metadata.MD
// 				baseFee := sdkmath.NewInt(1)
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				fQueryClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				s.backend.Cfg.JSONRPC.FeeHistoryCap = 2
// 				_, err := RegisterBlock(client, ethrpc.BlockNumber(1).Int64(), nil)
// 				s.Require().NoError(err)
// 				_, err = RegisterBlockResults(client, 1)
// 				s.Require().NoError(err)
// 				RegisterBaseFee(queryClient, baseFee)
// 				RegisterValidatorAccount(queryClient, validator)
// 				RegisterConsensusParams(client, 1)
// 				RegisterParams(queryClient, &header, 1)
// 				RegisterFeeMarketParams(fQueryClient, 1)
// 			},
// 			1,
// 			ethrpc.EarliestBlockNumber,
// 			&rpc.FeeHistoryResult{
// 				OldestBlock:  (*hexutil.Big)(big.NewInt(0)),
// 				BaseFee:      []*hexutil.Big{(*hexutil.Big)(big.NewInt(1)), (*hexutil.Big)(big.NewInt(1))},
// 				GasUsedRatio: []float64{0},
// 				Reward:       [][]*hexutil.Big{{(*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0)), (*hexutil.Big)(big.NewInt(0))}},
// 			},
// 			sdk.AccAddress(utiltx.GenerateAddress().Bytes()),
// 			true,
// 			nil,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock(tc.validator)

// 			called := 0
// 			if len(tc.targetNewBaseFees) > 0 {
// 				s.backend.ProcessBlocker = func(
// 					tendermintBlock *tmrpctypes.ResultBlock,
// 					ethBlock *map[string]interface{},
// 					rewardPercentiles []float64,
// 					tendermintBlockResult *tmrpctypes.ResultBlockResults,
// 					targetOneFeeHistory *rpc.OneFeeHistory,
// 				) error {
// 					err := s.backend.ProcessBlock(tendermintBlock, ethBlock, rewardPercentiles, tendermintBlockResult, targetOneFeeHistory)
// 					s.Require().NoError(err)
// 					targetOneFeeHistory.NextBaseFee = tc.targetNewBaseFees[called]
// 					called++
// 					return nil
// 				}
// 			}

// 			feeHistory, err := s.backend.FeeHistory(tc.userBlockCount, tc.latestBlock, []float64{25, 50, 75, 100})
// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(feeHistory, tc.expFeeHistory)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }
