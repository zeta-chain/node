package backend

import (
	"fmt"
	"math/big"

	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/viper"
	"github.com/zeta-chain/ethermint/crypto/ethsecp256k1"
	ethermint "github.com/zeta-chain/ethermint/types"
	"google.golang.org/grpc/metadata"

	"github.com/zeta-chain/node/rpc/backend/mocks"
)

func (suite *BackendTestSuite) TestRPCMinGasPrice() {
	testCases := []struct {
		name           string
		registerMock   func()
		expMinGasPrice int64
		expPass        bool
	}{
		{
			"pass - default gas price",
			func() {
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParamsWithoutHeaderError(queryClient, 1)
			},
			ethermint.DefaultGasPrice,
			true,
		},
		{
			"pass - min gas price is 0",
			func() {
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParamsWithoutHeader(queryClient, 1)
			},
			ethermint.DefaultGasPrice,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			minPrice := suite.backend.RPCMinGasPrice()
			if tc.expPass {
				suite.Require().Equal(tc.expMinGasPrice, minPrice)
			} else {
				suite.Require().NotEqual(tc.expMinGasPrice, minPrice)
			}
		})
	}
}

func (suite *BackendTestSuite) TestSetGasPrice() {
	defaultGasPrice := (*hexutil.Big)(big.NewInt(1))
	testCases := []struct {
		name         string
		registerMock func()
		gasPrice     hexutil.Big
		expOutput    bool
	}{
		{
			"pass - cannot get server config",
			func() {
				suite.backend.clientCtx.Viper = viper.New()
			},
			*defaultGasPrice,
			false,
		},
		{
			"pass - cannot find coin denom",
			func() {
				suite.backend.clientCtx.Viper = viper.New()
				suite.backend.clientCtx.Viper.Set("telemetry.global-labels", []interface{}{})
			},
			*defaultGasPrice,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()
			output := suite.backend.SetGasPrice(tc.gasPrice)
			suite.Require().Equal(tc.expOutput, output)
		})
	}
}

func (suite *BackendTestSuite) TestAccounts() {
	testCases := []struct {
		name                    string
		accountsRetrievalMethod func() ([]common.Address, error)
		registerMock            func()
		expAddr                 []common.Address
		expPass                 bool
	}{
		{
			"pass - returns acc from keyring from ListAccounts",
			suite.backend.ListAccounts,
			func() {},
			[]common.Address{},
			true,
		},
		{
			"pass - returns acc from keyring from Accounts",
			suite.backend.Accounts,
			func() {},
			[]common.Address{},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			output, err := tc.accountsRetrievalMethod()

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expAddr, output)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestSyncing() {
	testCases := []struct {
		name         string
		registerMock func()
		expResponse  interface{}
		expPass      bool
	}{
		{
			"fail - Can't get status",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterStatusError(client)
			},
			false,
			false,
		},
		{
			"pass - Node not catching up",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterStatus(client)
			},
			false,
			true,
		},
		{
			"pass - Node is catching up",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterStatus(client)
				status, _ := client.Status(suite.backend.ctx)
				status.SyncInfo.CatchingUp = true
			},
			map[string]interface{}{
				"startingBlock": hexutil.Uint64(0),
				"currentBlock":  hexutil.Uint64(0),
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			output, err := suite.backend.Syncing()

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expResponse, output)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestSetEtherbase() {
	testCases := []struct {
		name         string
		registerMock func()
		etherbase    common.Address
		expResult    bool
	}{
		{
			"pass - Failed to get coinbase address",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterStatusError(client)
			},
			common.Address{},
			false,
		},
		{
			"pass - the minimum fee is not set",
			func() {
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterStatus(client)
				RegisterValidatorAccount(queryClient, suite.acc)
			},
			common.Address{},
			false,
		},
		{
			"fail - error querying for account ",
			func() {
				var header metadata.MD
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterStatus(client)
				RegisterValidatorAccount(queryClient, suite.acc)
				RegisterParams(queryClient, &header, 1)
				c := sdk.NewDecCoin("azeta", sdk.NewIntFromBigInt(big.NewInt(1)))
				suite.backend.cfg.SetMinGasPrices(sdk.DecCoins{c})
				delAddr, _ := suite.backend.GetCoinbase()
				// account, _ := suite.backend.clientCtx.AccountRetriever.GetAccount(suite.backend.clientCtx, delAddr)
				delCommonAddr := common.BytesToAddress(delAddr.Bytes())
				request := &authtypes.QueryAccountRequest{Address: sdk.AccAddress(delCommonAddr.Bytes()).String()}
				requestMarshal, _ := request.Marshal()
				RegisterABCIQueryWithOptionsError(
					client,
					"/cosmos.auth.v1beta1.Query/Account",
					requestMarshal,
					tmrpcclient.ABCIQueryOptions{Height: int64(1), Prove: false},
				)
			},
			common.Address{},
			false,
		},
		{
			"pass - set the etherbase for the miner",
			func() {
				priv, err := ethsecp256k1.GenerateKey()
				suite.Require().NoError(err)

				armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
				suite.backend.clientCtx.Keyring.ImportPrivKey("test_key", armor, "")

				suite.acc = sdk.AccAddress(priv.PubKey().Address().Bytes())
				accounts := map[string]client.TestAccount{}
				accounts[suite.acc.String()] = client.TestAccount{
					Address: suite.acc,
					Num:     uint64(1),
					Seq:     uint64(1),
				}

				suite.backend.clientCtx = suite.backend.clientCtx.WithAccountRetriever(
					client.TestAccountRetriever{Accounts: accounts},
				)
				var header metadata.MD
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParams(queryClient, &header, 1)
				suite.backend.clientCtx = suite.backend.clientCtx.WithInterfaceRegistry(
					codectypes.NewInterfaceRegistry(),
				)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				RegisterStatus(client)
				RegisterValidatorAccount(queryClient, suite.acc)
				c := sdk.NewDecCoin("azeta", sdk.NewIntFromBigInt(big.NewInt(1)))
				suite.backend.cfg.SetMinGasPrices(sdk.DecCoins{c})
				delAddr, _ := suite.backend.GetCoinbase()
				account, _ := suite.backend.clientCtx.AccountRetriever.GetAccount(suite.backend.clientCtx, delAddr)
				delCommonAddr := common.BytesToAddress(delAddr.Bytes())
				request := &authtypes.QueryAccountRequest{Address: sdk.AccAddress(delCommonAddr.Bytes()).String()}
				requestMarshal, _ := request.Marshal()
				RegisterUnconfirmedTxsEmpty(client, nil)
				RegisterABCIQueryAccount(
					client,
					requestMarshal,
					tmrpcclient.ABCIQueryOptions{Height: int64(1), Prove: false},
					account,
				)
				RegisterABCIQuerySimulate(client, tmrpcclient.ABCIQueryOptions{Height: int64(1), Prove: false})
				RegisterBroadcastTxAny(client)
			},
			common.Address{},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			output := suite.backend.SetEtherbase(tc.etherbase)

			suite.Require().Equal(tc.expResult, output)
		})
	}
}

func (suite *BackendTestSuite) TestImportRawKey() {
	priv, _ := ethsecp256k1.GenerateKey()
	privHex := common.Bytes2Hex(priv.Bytes())
	pubAddr := common.BytesToAddress(priv.PubKey().Address().Bytes())

	testCases := []struct {
		name         string
		registerMock func()
		privKey      string
		password     string
		expAddr      common.Address
		expPass      bool
	}{
		{
			"fail - not a valid private key",
			func() {},
			"",
			"",
			common.Address{},
			false,
		},
		{
			"pass - returning correct address",
			func() {},
			privHex,
			"",
			pubAddr,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			output, err := suite.backend.ImportRawKey(tc.privKey, tc.password)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expAddr, output)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
