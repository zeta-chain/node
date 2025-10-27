package backend

// import (
// 	"fmt"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"github.com/spf13/viper"
// 	"google.golang.org/grpc/metadata"

// 	tmrpcclient "github.com/cometbft/cometbft/rpc/client"

// 	"github.com/cosmos/evm/crypto/ethsecp256k1"
// 	"github.com/cosmos/evm/server/config"
// 	"github.com/cosmos/evm/testutil/constants"
// 	evmtypes "github.com/cosmos/evm/x/vm/types"
// 	"github.com/zeta-chain/node/rpc/backend/mocks"

// 	"cosmossdk.io/math"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// )

// func (s *TestSuite) TestRPCMinGasPrice() {
// 	testCases := []struct {
// 		name           string
// 		registerMock   func()
// 		expMinGasPrice *big.Int
// 		expPass        bool
// 	}{
// 		{
// 			"pass - default gas price",
// 			func() {},
// 			big.NewInt(constants.DefaultGasPrice),
// 			true,
// 		},
// 		{
// 			"pass - min gas price is 0",
// 			func() {},
// 			big.NewInt(constants.DefaultGasPrice),
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			minPrice := s.backend.RPCMinGasPrice()
// 			if tc.expPass {
// 				s.Require().Equal(tc.expMinGasPrice, minPrice)
// 			} else {
// 				s.Require().NotEqual(tc.expMinGasPrice, minPrice)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestGenerateMinGasCoin() {
// 	defaultGasPrice := (*hexutil.Big)(big.NewInt(1))
// 	testCases := []struct {
// 		name           string
// 		gasPrice       hexutil.Big
// 		minGas         sdk.DecCoins
// 		expectedOutput sdk.DecCoin
// 	}{
// 		{
// 			"pass - empty min gas Coins (default denom)",
// 			*defaultGasPrice,
// 			sdk.DecCoins{},
// 			sdk.DecCoin{
// 				Denom:  evmtypes.GetEVMCoinDenom(),
// 				Amount: math.LegacyNewDecFromBigInt(defaultGasPrice.ToInt()),
// 			},
// 		},
// 		{
// 			"pass - different min gas Coin",
// 			*defaultGasPrice,
// 			sdk.DecCoins{sdk.NewDecCoin("test", math.NewInt(1))},
// 			sdk.DecCoin{
// 				Denom:  "test",
// 				Amount: math.LegacyNewDecFromBigInt(defaultGasPrice.ToInt()),
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			s.backend.ClientCtx.Viper = viper.New()

// 			appConf := config.DefaultConfig()
// 			appConf.SetMinGasPrices(tc.minGas)

// 			output := s.backend.GenerateMinGasCoin(tc.gasPrice, *appConf)
// 			s.Require().Equal(tc.expectedOutput, output)
// 		})
// 	}
// }

// // TODO: Combine these 2 into one test since the code is identical
// func (s *TestSuite) TestListAccounts() {
// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		expAddr      []common.Address
// 		expPass      bool
// 	}{
// 		{
// 			"pass - returns empty address",
// 			func() {},
// 			[]common.Address{},
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			output, err := s.backend.ListAccounts()

// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expAddr, output)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestAccounts() {
// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		expAddr      []common.Address
// 		expPass      bool
// 	}{
// 		{
// 			"pass - returns empty address",
// 			func() {},
// 			[]common.Address{},
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			output, err := s.backend.Accounts()

// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expAddr, output)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestSyncing() {
// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		expResponse  interface{}
// 		expPass      bool
// 	}{
// 		{
// 			"fail - Can't get status",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				RegisterStatusError(client)
// 			},
// 			false,
// 			false,
// 		},
// 		{
// 			"pass - Node not catching up",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				RegisterStatus(client)
// 			},
// 			false,
// 			true,
// 		},
// 		{
// 			"pass - Node is catching up",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				RegisterStatus(client)
// 				status, _ := client.Status(s.backend.Ctx)
// 				status.SyncInfo.CatchingUp = true
// 			},
// 			map[string]interface{}{
// 				"startingBlock": hexutil.Uint64(0),
// 				"currentBlock":  hexutil.Uint64(0),
// 			},
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			output, err := s.backend.Syncing()

// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expResponse, output)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }

// func (s *TestSuite) TestSetEtherbase() {
// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		etherbase    common.Address
// 		expResult    bool
// 	}{
// 		{
// 			"pass - Failed to get coinbase address",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				RegisterStatusError(client)
// 			},
// 			common.Address{},
// 			false,
// 		},
// 		{
// 			"pass - the minimum fee is not set",
// 			func() {
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterStatus(client)
// 				RegisterValidatorAccount(QueryClient, s.acc)
// 			},
// 			common.Address{},
// 			false,
// 		},
// 		{
// 			"fail - error querying for account",
// 			func() {
// 				var header metadata.MD
// 				client := s.backend.ClientCtx.Client.(*mocks.Client)
// 				QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 				RegisterStatus(client)
// 				RegisterValidatorAccount(QueryClient, s.acc)
// 				RegisterParams(QueryClient, &header, 1)
// 				c := sdk.NewDecCoin(constants.ExampleAttoDenom, math.NewIntFromBigInt(big.NewInt(1)))
// 				s.backend.Cfg.SetMinGasPrices(sdk.DecCoins{c})
// 				delAddr, _ := s.backend.GetCoinbase()
// 				// account, _ := s.backend.ClientCtx.AccountRetriever.GetAccount(s.backend.ClientCtx, delAddr)
// 				delCommonAddr := common.BytesToAddress(delAddr.Bytes())
// 				request := &authtypes.QueryAccountRequest{Address: sdk.AccAddress(delCommonAddr.Bytes()).String()}
// 				requestMarshal, _ := request.Marshal()
// 				RegisterABCIQueryWithOptionsError(
// 					client,
// 					"/cosmos.auth.v1beta1.Query/Account",
// 					requestMarshal,
// 					tmrpcclient.ABCIQueryOptions{Height: int64(1), Prove: false},
// 				)
// 			},
// 			common.Address{},
// 			false,
// 		},
// 		// TODO: Finish this test case once ABCIQuery GetAccount is fixed
// 		// {
// 		//	"pass - set the etherbase for the miner",
// 		//	func() {
// 		//		client := s.backend.ClientCtx.Client.(*mocks.Client)
// 		//		QueryClient := s.backend.QueryClient.QueryClient.(*mocks.EVMQueryClient)
// 		//		RegisterStatus(client)
// 		//		RegisterValidatorAccount(QueryClient, s.acc)
// 		//		c := sdk.NewDecCoin(testconstants.ExampleAttoDenom, math.NewIntFromBigInt(big.NewInt(1)))
// 		//		s.backend.Cfg.SetMinGasPrices(sdk.DecCoins{c})
// 		//		delAddr, _ := s.backend.GetCoinbase()
// 		//		account, _ := s.backend.ClientCtx.AccountRetriever.GetAccount(s.backend.ClientCtx, delAddr)
// 		//		delCommonAddr := common.BytesToAddress(delAddr.Bytes())
// 		//		request := &authtypes.QueryAccountRequest{Address: sdk.AccAddress(delCommonAddr.Bytes()).String()}
// 		//		requestMarshal, _ := request.Marshal()
// 		//		RegisterABCIQueryAccount(
// 		//			client,
// 		//			requestMarshal,
// 		//			tmrpcclient.ABCIQueryOptions{Height: int64(1), Prove: false},
// 		//			account,
// 		//		)
// 		//	},
// 		//	common.Address{},
// 		//	false,
// 		// },
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			output := s.backend.SetEtherbase(tc.etherbase)

// 			s.Require().Equal(tc.expResult, output)
// 		})
// 	}
// }

// func (s *TestSuite) TestImportRawKey() {
// 	priv, _ := ethsecp256k1.GenerateKey()
// 	privHex := common.Bytes2Hex(priv.Bytes())
// 	pubAddr := common.BytesToAddress(priv.PubKey().Address().Bytes())

// 	testCases := []struct {
// 		name         string
// 		registerMock func()
// 		privKey      string
// 		password     string
// 		expAddr      common.Address
// 		expPass      bool
// 	}{
// 		{
// 			"fail - not a valid private key",
// 			func() {},
// 			"",
// 			"",
// 			common.Address{},
// 			false,
// 		},
// 		{
// 			"pass - returning correct address",
// 			func() {},
// 			privHex,
// 			"",
// 			pubAddr,
// 			true,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		s.Run(fmt.Sprintf("case %s", tc.name), func() {
// 			s.SetupTest() // reset test and queries
// 			tc.registerMock()

// 			output, err := s.backend.ImportRawKey(tc.privKey, tc.password)
// 			if tc.expPass {
// 				s.Require().NoError(err)
// 				s.Require().Equal(tc.expAddr, output)
// 			} else {
// 				s.Require().Error(err)
// 			}
// 		})
// 	}
// }
