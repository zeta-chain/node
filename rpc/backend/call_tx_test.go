package backend

import (
	"encoding/json"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	utiltx "github.com/cosmos/evm/testutil/tx"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"

	"github.com/zeta-chain/node/rpc/backend/mocks"
	rpctypes "github.com/zeta-chain/node/rpc/types"
)

func (s *TestSuite) TestResend() {
	txNonce := (hexutil.Uint64)(1)
	baseFee := math.NewInt(1)
	gasPrice := new(hexutil.Big)
	toAddr := utiltx.GenerateAddress()
	evmChainID := (*hexutil.Big)(s.backend.EvmChainID)
	callArgs := evmtypes.TransactionArgs{
		From:                 nil,
		To:                   &toAddr,
		Gas:                  nil,
		GasPrice:             nil,
		MaxFeePerGas:         gasPrice,
		MaxPriorityFeePerGas: gasPrice,
		Value:                gasPrice,
		Nonce:                &txNonce,
		Input:                nil,
		Data:                 nil,
		AccessList:           nil,
		ChainID:              evmChainID,
	}

	testCases := []struct {
		name         string
		registerMock func()
		args         evmtypes.TransactionArgs
		gasPrice     *hexutil.Big
		gasLimit     *hexutil.Uint64
		expHash      common.Hash
		expPass      bool
	}{
		{
			"fail - Missing transaction nonce",
			func() {},
			evmtypes.TransactionArgs{
				Nonce: nil,
			},
			nil,
			nil,
			common.Hash{},
			false,
		},
		{
			"pass - Can't set Tx defaults BaseFee disabled",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(QueryClient, &header, 1)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFeeDisabled(QueryClient)
			},
			evmtypes.TransactionArgs{
				Nonce:   &txNonce,
				ChainID: callArgs.ChainID,
			},
			nil,
			nil,
			common.Hash{},
			true,
		},
		{
			"pass - Can't set Tx defaults",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				feeMarketClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
				RegisterParams(QueryClient, &header, 1)
				RegisterFeeMarketParams(feeMarketClient, 1)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFee(QueryClient, baseFee)
			},
			evmtypes.TransactionArgs{
				Nonce: &txNonce,
			},
			nil,
			nil,
			common.Hash{},
			true,
		},
		{
			"pass - MaxFeePerGas is nil",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(QueryClient, &header, 1)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFeeDisabled(QueryClient)
			},
			evmtypes.TransactionArgs{
				Nonce:                &txNonce,
				MaxPriorityFeePerGas: nil,
				GasPrice:             nil,
				MaxFeePerGas:         nil,
			},
			nil,
			nil,
			common.Hash{},
			true,
		},
		{
			"fail - GasPrice and (MaxFeePerGas or MaxPriorityPerGas specified)",
			func() {},
			evmtypes.TransactionArgs{
				Nonce:                &txNonce,
				MaxPriorityFeePerGas: nil,
				GasPrice:             gasPrice,
				MaxFeePerGas:         gasPrice,
			},
			nil,
			nil,
			common.Hash{},
			false,
		},
		{
			"fail - Block error",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(QueryClient, &header, 1)
				RegisterBlockError(client, 1)
			},
			evmtypes.TransactionArgs{
				Nonce: &txNonce,
			},
			nil,
			nil,
			common.Hash{},
			false,
		},
		{
			"pass - MaxFeePerGas is nil",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(QueryClient, &header, 1)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFee(QueryClient, baseFee)
			},
			evmtypes.TransactionArgs{
				Nonce:                &txNonce,
				GasPrice:             nil,
				MaxPriorityFeePerGas: gasPrice,
				MaxFeePerGas:         gasPrice,
				ChainID:              callArgs.ChainID,
			},
			nil,
			nil,
			common.Hash{},
			true,
		},
		{
			"pass - Chain Id is nil",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(QueryClient, &header, 1)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFee(QueryClient, baseFee)
			},
			evmtypes.TransactionArgs{
				Nonce:                &txNonce,
				MaxPriorityFeePerGas: gasPrice,
				ChainID:              nil,
			},
			nil,
			nil,
			common.Hash{},
			true,
		},
		{
			"fail - Pending transactions error",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterEstimateGas(QueryClient, callArgs)
				RegisterParams(QueryClient, &header, 1)
				RegisterUnconfirmedTxsError(client, nil)
			},
			evmtypes.TransactionArgs{
				Nonce:                &txNonce,
				To:                   &toAddr,
				MaxFeePerGas:         gasPrice,
				MaxPriorityFeePerGas: gasPrice,
				Value:                gasPrice,
				Gas:                  nil,
				ChainID:              callArgs.ChainID,
			},
			gasPrice,
			nil,
			common.Hash{},
			false,
		},
		{
			"fail - Not Ethereum txs",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFee(QueryClient, baseFee)
				RegisterEstimateGas(QueryClient, callArgs)
				RegisterParams(QueryClient, &header, 1)

				RegisterUnconfirmedTxsEmpty(client, nil)
			},
			evmtypes.TransactionArgs{
				Nonce:                &txNonce,
				To:                   &toAddr,
				MaxFeePerGas:         gasPrice,
				MaxPriorityFeePerGas: gasPrice,
				Value:                gasPrice,
				Gas:                  nil,
				ChainID:              callArgs.ChainID,
			},
			gasPrice,
			nil,
			common.Hash{},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock()

			hash, err := s.backend.Resend(tc.args, tc.gasPrice, tc.gasLimit)

			if tc.expPass {
				s.Require().Equal(tc.expHash, hash)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestSendRawTransaction() {
	s.SetupTest() // reset

	ethTx, bz := s.buildEthereumTx()
	ethTx.From = s.from.Bytes()

	emptyEvmChainIDTx := s.buildEthereumTxWithChainID(nil)
	invalidEvmChainIDTx := s.buildEthereumTxWithChainID(big.NewInt(1))

	// Sign the ethTx
	ethSigner := ethtypes.LatestSigner(s.backend.ChainConfig())
	err := ethTx.Sign(ethSigner, s.signer)
	s.Require().NoError(err)

	rlpEncodedBz, _ := rlp.EncodeToBytes(ethTx.AsTransaction())
	evmDenom := evmtypes.GetEVMCoinDenom()

	testCases := []struct {
		name         string
		registerMock func()
		rawTx        func() []byte
		expHash      common.Hash
		expError     string
		expPass      bool
	}{
		{
			"fail - empty bytes",
			func() {},
			func() []byte { return []byte{} },
			common.Hash{},
			"",
			false,
		},
		{
			"fail - no RLP encoded bytes",
			func() {},
			func() []byte { return bz },
			common.Hash{},
			"",
			false,
		},
		{
			"fail - invalid chain-id",
			func() {
				s.backend.AllowUnprotectedTxs = false
			},
			func() []byte {
				from, _ := utiltx.NewAddrKey()
				invalidEvmChainIDTx.From = from.Bytes()
				bytes, _ := rlp.EncodeToBytes(invalidEvmChainIDTx.AsTransaction())
				return bytes
			},
			common.Hash{},
			fmt.Errorf("incorrect chain-id; expected %d, got %d", 262144, big.NewInt(1)).Error(),
			false,
		},
		{
			"fail - unprotected tx",
			func() {
				s.backend.AllowUnprotectedTxs = false
			},
			func() []byte {
				bytes, _ := rlp.EncodeToBytes(emptyEvmChainIDTx.AsTransaction())
				return bytes
			},
			common.Hash{},
			errors.New("only replay-protected (EIP-155) transactions allowed over RPC").Error(),
			false,
		},
		{
			"fail - failed to broadcast transaction",
			func() {
				cosmosTx, _ := ethTx.BuildTx(s.backend.ClientCtx.TxConfig.NewTxBuilder(), evmDenom)
				txBytes, _ := s.backend.ClientCtx.TxConfig.TxEncoder()(cosmosTx)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				s.backend.AllowUnprotectedTxs = true
				RegisterBroadcastTxError(client, txBytes)
			},
			func() []byte {
				bytes, _ := rlp.EncodeToBytes(ethTx.AsTransaction())
				return bytes
			},
			common.HexToHash(ethTx.Hash),
			errortypes.ErrInvalidRequest.Error(),
			false,
		},
		{
			"pass - Gets the correct transaction hash of the eth transaction",
			func() {
				cosmosTx, _ := ethTx.BuildTx(s.backend.ClientCtx.TxConfig.NewTxBuilder(), evmDenom)
				txBytes, _ := s.backend.ClientCtx.TxConfig.TxEncoder()(cosmosTx)

				client := s.backend.ClientCtx.Client.(*mocks.Client)
				s.backend.AllowUnprotectedTxs = true
				RegisterBroadcastTx(client, txBytes)
			},
			func() []byte { return rlpEncodedBz },
			common.HexToHash(ethTx.Hash),
			"",
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock()

			hash, err := s.backend.SendRawTransaction(tc.rawTx())

			if tc.expPass {
				s.Require().Equal(tc.expHash, hash)
			} else {
				s.Require().Errorf(err, tc.expError)
			}
		})
	}
}

func (s *TestSuite) TestDoCall() {
	_, bz := s.buildEthereumTx()
	gasPrice := (*hexutil.Big)(big.NewInt(1))
	toAddr := utiltx.GenerateAddress()
	evmChainID := (*hexutil.Big)(s.backend.EvmChainID)
	callArgs := evmtypes.TransactionArgs{
		From:                 nil,
		To:                   &toAddr,
		Gas:                  nil,
		GasPrice:             nil,
		MaxFeePerGas:         gasPrice,
		MaxPriorityFeePerGas: gasPrice,
		Value:                gasPrice,
		Input:                nil,
		Data:                 nil,
		AccessList:           nil,
		ChainID:              evmChainID,
	}
	argsBz, err := json.Marshal(callArgs)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		registerMock func()
		blockNum     rpctypes.BlockNumber
		callArgs     evmtypes.TransactionArgs
		expEthTx     *evmtypes.MsgEthereumTxResponse
		expPass      bool
	}{
		{
			"fail - Invalid request",
			func() {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				_, err := RegisterBlock(client, 1, bz)
				s.Require().NoError(err)
				RegisterEthCallError(
					QueryClient,
					&evmtypes.EthCallRequest{Args: argsBz, ChainId: s.backend.EvmChainID.Int64()},
				)
			},
			rpctypes.BlockNumber(1),
			callArgs,
			&evmtypes.MsgEthereumTxResponse{},
			false,
		},
		{
			"pass - Returned transaction response",
			func() {
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				_, err := RegisterBlock(client, 1, bz)
				s.Require().NoError(err)
				RegisterEthCall(
					QueryClient,
					&evmtypes.EthCallRequest{Args: argsBz, ChainId: s.backend.EvmChainID.Int64()},
				)
			},
			rpctypes.BlockNumber(1),
			callArgs,
			&evmtypes.MsgEthereumTxResponse{},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock()

			msgEthTx, err := s.backend.DoCall(tc.callArgs, tc.blockNum)

			if tc.expPass {
				s.Require().Equal(tc.expEthTx, msgEthTx)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *TestSuite) TestGasPrice() {
	defaultGasPrice := (*hexutil.Big)(big.NewInt(1))

	testCases := []struct {
		name         string
		registerMock func()
		expGas       *hexutil.Big
		expPass      bool
	}{
		{
			"pass - get the default gas price",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				feeMarketClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
				RegisterFeeMarketParams(feeMarketClient, 1)
				RegisterParams(QueryClient, &header, 1)
				RegisterGlobalMinGasPrice(QueryClient, 1)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFee(QueryClient, math.NewInt(1))
			},
			defaultGasPrice,
			true,
		},
		{
			"fail - can't get gasFee, FeeMarketParams error",
			func() {
				var header metadata.MD
				client := s.backend.ClientCtx.Client.(*mocks.Client)
				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
				feeMarketClient := s.backend.QueryClient.FeeMarket.(*mocks.FeeMarketQueryClient)
				RegisterFeeMarketParamsError(feeMarketClient, 1)
				RegisterParams(QueryClient, &header, 1)
				_, err := RegisterBlock(client, 1, nil)
				s.Require().NoError(err)
				_, err = RegisterBlockResults(client, 1)
				s.Require().NoError(err)
				RegisterBaseFee(QueryClient, math.NewInt(1))
			},
			defaultGasPrice,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("case %s", tc.name), func() {
			s.SetupTest() // reset test and queries
			tc.registerMock()

			gasPrice, err := s.backend.GasPrice()
			if tc.expPass {
				s.Require().Equal(tc.expGas, gasPrice)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
