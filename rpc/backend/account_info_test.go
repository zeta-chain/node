package backend

// import (
// 	"fmt"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"google.golang.org/grpc/metadata"

// 	"github.com/cometbft/cometbft/libs/bytes"
// 	cmtrpcclient "github.com/cometbft/cometbft/rpc/client"

// 	utiltx "github.com/cosmos/evm/testutil/tx"
// 	evmtypes "github.com/cosmos/evm/x/vm/types"
// 	"github.com/zeta-chain/node/rpc/backend/mocks"
// 	rpctypes "github.com/zeta-chain/node/rpc/types"

// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// )

// func (s *TestSuite) TestGetCode() {
// 	blockNr := rpctypes.NewBlockNumber(big.NewInt(1))
// 	contractCode := []byte("0xef616c92f3cfc9e92dc270d6acff9cea213cecc7020a76ee4395af09bdceb4837a1ebdb5735e11e7d3adb6104e0c3ac55180b4ddf5e54d022cc5e8837f6a4f971b")

// 	testCases := []struct {
// 		name          string
// 		addr          common.Address
// 		blockNrOrHash rpctypes.BlockNumberOrHash
// 		registerMock  func(common.Address)
// 		expPass       bool
// 		expCode       hexutil.Bytes
// 	}{
// 		{
// 			"fail - BlockHash and BlockNumber are both nil ",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{},
// 			func(_ common.Address) {},
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - query client errors on getting Code",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(addr common.Address) {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterCodeError(QueryClient, addr)
// 			},
// 			false,
// 			nil,
// 		},
// 		{
// 			"pass",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(addr common.Address) {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterCode(QueryClient, addr, contractCode)
// 			},
// 			true,
// 			contractCode,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
// 			s.SetupTest() // reset
// 			tc.registerMock(tc.addr)

// 			code, err := s.backend.GetCode(tc.addr, tc.blockNrOrHash)
// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expCode, code)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestGetProof() {
// 	blockNrInvalid := rpctypes.NewBlockNumber(big.NewInt(1))
// 	blockNr := rpctypes.NewBlockNumber(big.NewInt(4))
// 	blockNrZero := rpctypes.NewBlockNumber(big.NewInt(0))
// 	address1 := utiltx.GenerateAddress()

// 	testCases := []struct {
// 		name          string
// 		addr          common.Address
// 		storageKeys   []string
// 		blockNrOrHash rpctypes.BlockNumberOrHash
// 		registerMock  func(rpctypes.BlockNumber, common.Address)
// 		expPass       bool
// 		expAccRes     *rpctypes.AccountResult
// 	}{
// 		{
// 			"fail - BlockNumeber = 1 (invalidBlockNumber)",
// 			address1,
// 			[]string{},
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNrInvalid},
// 			func(bn rpctypes.BlockNumber, addr common.Address) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				height := bn.Int64()
// 				RegisterHeader(client, &height, nil)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterAccount(QueryClient, addr, blockNrInvalid.Int64())
// 			},
// 			false,
// 			&rpctypes.AccountResult{},
// 		},
// 		{
// 			"fail - Block doesn't exist",
// 			address1,
// 			[]string{},
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNrInvalid},
// 			func(bn rpctypes.BlockNumber, _ common.Address) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				height := bn.Int64()
// 				RegisterHeaderError(client, &height)
// 			},
// 			false,
// 			&rpctypes.AccountResult{},
// 		},
// 		{
// 			"pass",
// 			address1,
// 			[]string{"0x0"},
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(bn rpctypes.BlockNumber, addr common.Address) {
// 				height := bn.Int64()
// 				s.backend.Ctx = rpctypes.ContextWithHeight(height)
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				RegisterHeader(client, &height, nil)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterAccount(QueryClient, addr, height)

// 				// Use the IAVL height if a valid CometBFT height is passed in.
// 				iavlHeight := bn.Int64()
// 				RegisterABCIQueryWithOptions(
// 					client,
// 					bn.Int64(),
// 					"store/evm/key",
// 					evmtypes.StateKey(address1, common.HexToHash("0x0").Bytes()),
// 					cmtrpcclient.ABCIQueryOptions{Height: iavlHeight, Prove: true},
// 				)
// 				RegisterABCIQueryWithOptions(
// 					client,
// 					bn.Int64(),
// 					"store/acc/key",
// 					bytes.HexBytes(append(authtypes.AddressStoreKeyPrefix, address1.Bytes()...)),
// 					cmtrpcclient.ABCIQueryOptions{Height: iavlHeight, Prove: true},
// 				)
// 			},
// 			true,
// 			&rpctypes.AccountResult{
// 				Address:      address1,
// 				AccountProof: []string{""},
// 				Balance:      (*hexutil.Big)(big.NewInt(0)),
// 				CodeHash:     common.HexToHash(""),
// 				Nonce:        0x0,
// 				StorageHash:  common.Hash{},
// 				StorageProof: []rpctypes.StorageResult{
// 					{
// 						Key:   "0x0",
// 						Value: (*hexutil.Big)(big.NewInt(2)),
// 						Proof: []string{""},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			"pass, 0 height",
// 			address1,
// 			[]string{"0x0"},
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNrZero},
// 			func(bn rpctypes.BlockNumber, addr common.Address) {
// 				height := int64(4)
// 				s.backend.Ctx = rpctypes.ContextWithHeight(height)
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				RegisterHeader(client, &height, nil)
// 				queryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterAccount(queryClient, addr, height)
// 				var header metadata.MD
// 				RegisterParams(queryClient, &header, height)

// 				// Use the IAVL height if a valid CometBFT height is passed in.
// 				iavlHeight := height
// 				RegisterABCIQueryWithOptions(
// 					client,
// 					iavlHeight,
// 					"store/evm/key",
// 					evmtypes.StateKey(address1, common.HexToHash("0x0").Bytes()),
// 					cmtrpcclient.ABCIQueryOptions{Height: iavlHeight, Prove: true},
// 				)
// 				RegisterABCIQueryWithOptions(
// 					client,
// 					iavlHeight,
// 					"store/acc/key",
// 					bytes.HexBytes(append(authtypes.AddressStoreKeyPrefix, address1.Bytes()...)),
// 					cmtrpcclient.ABCIQueryOptions{Height: iavlHeight, Prove: true},
// 				)
// 			},
// 			true,
// 			&rpctypes.AccountResult{
// 				Address:      address1,
// 				AccountProof: []string{""},
// 				Balance:      (*hexutil.Big)(big.NewInt(0)),
// 				CodeHash:     common.HexToHash(""),
// 				Nonce:        0x0,
// 				StorageHash:  common.Hash{},
// 				StorageProof: []rpctypes.StorageResult{
// 					{
// 						Key:   "0x0",
// 						Value: (*hexutil.Big)(big.NewInt(2)),
// 						Proof: []string{""},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
// 			s.SetupTest()
// 			tc.registerMock(*tc.blockNrOrHash.BlockNumber, tc.addr)

// 			accRes, err := s.backend.GetProof(tc.addr, tc.storageKeys, tc.blockNrOrHash)

// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expAccRes, accRes)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestGetStorageAt() {
// 	blockNr := rpctypes.NewBlockNumber(big.NewInt(1))

// 	testCases := []struct {
// 		name          string
// 		addr          common.Address
// 		key           string
// 		blockNrOrHash rpctypes.BlockNumberOrHash
// 		registerMock  func(common.Address, string, string)
// 		expPass       bool
// 		expStorage    hexutil.Bytes
// 	}{
// 		{
// 			"fail - BlockHash and BlockNumber are both nil",
// 			utiltx.GenerateAddress(),
// 			"0x0",
// 			rpctypes.BlockNumberOrHash{},
// 			func(common.Address, string, string) {},
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - query client errors on getting Storage",
// 			utiltx.GenerateAddress(),
// 			"0x0",
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(addr common.Address, key string, _ string) {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterStorageAtError(QueryClient, addr, key)
// 			},
// 			false,
// 			nil,
// 		},
// 		{
// 			"pass",
// 			utiltx.GenerateAddress(),
// 			"0x0",
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(addr common.Address, key string, storage string) {
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterStorageAt(QueryClient, addr, key, storage)
// 			},
// 			true,
// 			hexutil.Bytes{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
// 		},
// 	}
// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
// 			s.SetupTest()
// 			tc.registerMock(tc.addr, tc.key, tc.expStorage.String())

// 			storage, err := s.backend.GetStorageAt(tc.addr, tc.key, tc.blockNrOrHash)
// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expStorage, storage)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestGetBalance() {
// 	blockNr := rpctypes.NewBlockNumber(big.NewInt(1))

// 	testCases := []struct {
// 		name          string
// 		addr          common.Address
// 		blockNrOrHash rpctypes.BlockNumberOrHash
// 		registerMock  func(rpctypes.BlockNumber, common.Address)
// 		expPass       bool
// 		expBalance    *hexutil.Big
// 	}{
// 		{
// 			"fail - BlockHash and BlockNumber are both nil",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{},
// 			func(rpctypes.BlockNumber, common.Address) {
// 			},
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - CometBFT client failed to get block",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(bn rpctypes.BlockNumber, _ common.Address) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				height := bn.Int64()
// 				RegisterHeaderError(client, &height)
// 			},
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - query client failed to get balance",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(bn rpctypes.BlockNumber, addr common.Address) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				height := bn.Int64()
// 				RegisterHeader(client, &height, nil)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBalanceError(QueryClient, addr, height)
// 			},
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - invalid balance",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(bn rpctypes.BlockNumber, addr common.Address) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				height := bn.Int64()
// 				RegisterHeader(client, &height, nil)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBalanceInvalid(QueryClient, addr, height)
// 			},
// 			false,
// 			nil,
// 		},
// 		{
// 			"fail - pruned node state",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(bn rpctypes.BlockNumber, addr common.Address) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				height := bn.Int64()
// 				RegisterHeader(client, &height, nil)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBalanceNegative(QueryClient, addr, height)
// 			},
// 			false,
// 			nil,
// 		},
// 		{
// 			"pass",
// 			utiltx.GenerateAddress(),
// 			rpctypes.BlockNumberOrHash{BlockNumber: &blockNr},
// 			func(bn rpctypes.BlockNumber, addr common.Address) {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				height := bn.Int64()
// 				RegisterHeader(client, &height, nil)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterBalance(QueryClient, addr, height)
// 			},
// 			true,
// 			(*hexutil.Big)(big.NewInt(1)),
// 		},
// 	}
// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
// 			s.SetupTest()

// 			// avoid nil pointer reference
// 			if tc.blockNrOrHash.BlockNumber != nil {
// 				tc.registerMock(*tc.blockNrOrHash.BlockNumber, tc.addr)
// 			}

// 			balance, err := s.backend.GetBalance(tc.addr, tc.blockNrOrHash)
// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expBalance, balance)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestGetTransactionCount() {
// 	testCases := []struct {
// 		name         string
// 		accExists    bool
// 		blockNum     rpctypes.BlockNumber
// 		registerMock func(common.Address, rpctypes.BlockNumber)
// 		expPass      bool
// 		expTxCount   hexutil.Uint64
// 	}{
// 		{
// 			"pass - account doesn't exist",
// 			false,
// 			rpctypes.NewBlockNumber(big.NewInt(1)),
// 			func(common.Address, rpctypes.BlockNumber) {
// 				var header metadata.MD
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterParams(QueryClient, &header, 1)
// 			},
// 			true,
// 			hexutil.Uint64(0),
// 		},
// 		{
// 			"fail - block height is in the future",
// 			false,
// 			rpctypes.NewBlockNumber(big.NewInt(10000)),
// 			func(common.Address, rpctypes.BlockNumber) {
// 				var header metadata.MD
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterParams(QueryClient, &header, 1)
// 			},
// 			false,
// 			hexutil.Uint64(0),
// 		},
// 		// TODO: Error mocking the GetAccount call - problem with Any type
// 		// {
// 		//	"pass - returns the number of transactions at the given address up to the given block number",
// 		//	true,
// 		//	rpctypes.NewBlockNumber(big.NewInt(1)),
// 		//	func(addr common.Address, bn rpctypes.BlockNumber) {
// 		//		client := s.backend.ClientCtx.Client.(*mocks.Client)
// 		//		account, err := s.backend.ClientCtx.AccountRetriever.GetAccount(s.backend.ClientCtx, s.acc)
// 		//		s.Require().NoError(err)
// 		//		request := &authtypes.QueryAccountRequest{Address: sdk.AccAddress(s.acc.Bytes()).String()}
// 		//		requestMarshal, _ := request.Marshal()
// 		//		RegisterABCIQueryAccount(
// 		//			client,
// 		//			requestMarshal,
// 		//			cmtrpcclient.ABCIQueryOptions{Height: int64(1), Prove: false},
// 		//			account,
// 		//		)
// 		//	},
// 		//	true,
// 		//	hexutil.Uint64(0),
// 		// },
// 	}
// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
// 			s.SetupTest()

// 			addr := utiltx.GenerateAddress()
// 			if tc.accExists {
// 				addr = common.BytesToAddress(s.acc.Bytes())
// 			}

// 			tc.registerMock(addr, tc.blockNum)

// 			txCount, err := s.backend.GetTransactionCount(addr, tc.blockNum)
// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expTxCount, *txCount)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }
