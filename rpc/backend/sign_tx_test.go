package backend

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	goethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/zeta-chain/ethermint/crypto/ethsecp256k1"
	"github.com/zeta-chain/ethermint/tests"
	"google.golang.org/grpc/metadata"

	"github.com/zeta-chain/node/rpc/backend/mocks"
)

func (suite *BackendTestSuite) TestSendTransaction() {
	gasPrice := new(hexutil.Big)
	gas := hexutil.Uint64(1)
	zeroGas := hexutil.Uint64(0)
	toAddr := tests.GenerateAddress()
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	from := common.BytesToAddress(priv.PubKey().Address().Bytes())
	nonce := hexutil.Uint64(1)
	baseFee := sdkmath.NewInt(1)
	validator := sdk.AccAddress(tests.GenerateAddress().Bytes())
	callArgsDefault := evmtypes.TransactionArgs{
		From:     &from,
		To:       &toAddr,
		GasPrice: gasPrice,
		Gas:      &gas,
		Nonce:    &nonce,
	}

	hash := common.Hash{}

	testCases := []struct {
		name         string
		registerMock func()
		args         evmtypes.TransactionArgs
		expHash      common.Hash
		expPass      bool
	}{
		{
			"fail - Can't find account in Keyring",
			func() {},
			evmtypes.TransactionArgs{},
			hash,
			false,
		},
		{
			"fail - Block error can't set Tx defaults",
			func() {
				var header metadata.MD
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
				suite.backend.clientCtx.Keyring.ImportPrivKey("test_key", armor, "")
				RegisterParams(queryClient, &header, 1)
				RegisterBlockError(client, 1)
			},
			callArgsDefault,
			hash,
			false,
		},
		{
			"fail - Cannot validate transaction gas set to 0",
			func() {
				var header metadata.MD
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
				suite.backend.clientCtx.Keyring.ImportPrivKey("test_key", armor, "")
				RegisterParams(queryClient, &header, 1)
				RegisterBlock(client, 1, nil)
				RegisterBlockResults(client, 1)
				RegisterBaseFee(queryClient, baseFee)
				RegisterValidatorAccount(queryClient, validator)
			},
			evmtypes.TransactionArgs{
				From:     &from,
				To:       &toAddr,
				GasPrice: gasPrice,
				Gas:      &zeroGas,
				Nonce:    &nonce,
			},
			hash,
			false,
		},
		{
			"fail - Cannot broadcast transaction",
			func() {
				var header metadata.MD
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
				suite.backend.clientCtx.Keyring.ImportPrivKey("test_key", armor, "")
				RegisterParams(queryClient, &header, 1)
				RegisterBlock(client, 1, nil)
				RegisterBlockResults(client, 1)
				RegisterBaseFee(queryClient, baseFee)
				RegisterParamsWithoutHeader(queryClient, 1)
				ethSigner := ethtypes.LatestSigner(suite.backend.ChainConfig())
				msg := callArgsDefault.ToTransaction()
				msg.Sign(ethSigner, suite.backend.clientCtx.Keyring)
				tx, err := msg.BuildTx(suite.backend.clientCtx.TxConfig.NewTxBuilder(), "azeta")
				suite.Require().NoError(err)
				txEncoder := suite.backend.clientCtx.TxConfig.TxEncoder()
				txBytes, err := txEncoder(tx)
				suite.Require().NoError(err)
				RegisterBroadcastTxError(client, txBytes)
				RegisterValidatorAccount(queryClient, validator)
			},
			callArgsDefault,
			common.Hash{},
			false,
		},
		{
			"pass - Return the transaction hash",
			func() {
				var header metadata.MD
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				client := suite.backend.clientCtx.Client.(*mocks.Client)
				armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
				suite.backend.clientCtx.Keyring.ImportPrivKey("test_key", armor, "")
				RegisterParams(queryClient, &header, 1)
				RegisterBlock(client, 1, nil)
				RegisterBlockResults(client, 1)
				RegisterBaseFee(queryClient, baseFee)
				RegisterParamsWithoutHeader(queryClient, 1)
				ethSigner := ethtypes.LatestSigner(suite.backend.ChainConfig())
				msg := callArgsDefault.ToTransaction()
				msg.Sign(ethSigner, suite.backend.clientCtx.Keyring)
				tx, err := msg.BuildTx(suite.backend.clientCtx.TxConfig.NewTxBuilder(), "azeta")
				suite.Require().NoError(err)
				txEncoder := suite.backend.clientCtx.TxConfig.TxEncoder()
				txBytes, err := txEncoder(tx)
				suite.Require().NoError(err)
				RegisterBroadcastTx(client, txBytes)
				RegisterValidatorAccount(queryClient, validator)
			},
			callArgsDefault,
			hash,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			if tc.expPass {
				// Sign the transaction and get the hash
				queryClient := suite.backend.queryClient.QueryClient.(*mocks.EVMQueryClient)
				RegisterParamsWithoutHeader(queryClient, 1)
				ethSigner := ethtypes.LatestSigner(suite.backend.ChainConfig())
				msg := callArgsDefault.ToTransaction()
				msg.Sign(ethSigner, suite.backend.clientCtx.Keyring)
				tc.expHash = msg.AsTransaction().Hash()
			}
			responseHash, err := suite.backend.SendTransaction(tc.args)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expHash, responseHash)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestSign() {
	from, priv := tests.NewAddrKey()
	testCases := []struct {
		name         string
		registerMock func()
		fromAddr     common.Address
		inputBz      hexutil.Bytes
		expPass      bool
	}{
		{
			"fail - can't find key in Keyring",
			func() {},
			from,
			nil,
			false,
		},
		{
			"pass - sign nil data",
			func() {
				armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
				suite.backend.clientCtx.Keyring.ImportPrivKey("test_key", armor, "")
			},
			from,
			nil,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			responseBz, err := suite.backend.Sign(tc.fromAddr, tc.inputBz)
			if tc.expPass {
				signature, _, err := suite.backend.clientCtx.Keyring.SignByAddress(
					(sdk.AccAddress)(from.Bytes()),
					tc.inputBz,
					signing.SignMode_SIGN_MODE_TEXTUAL,
				)
				signature[goethcrypto.RecoveryIDOffset] += 27
				suite.Require().NoError(err)
				suite.Require().Equal((hexutil.Bytes)(signature), responseBz)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *BackendTestSuite) TestSignTypedData() {
	from, priv := tests.NewAddrKey()
	testCases := []struct {
		name           string
		registerMock   func()
		fromAddr       common.Address
		inputTypedData apitypes.TypedData
		expPass        bool
	}{
		{
			"fail - can't find key in Keyring",
			func() {},
			from,
			apitypes.TypedData{},
			false,
		},
		{
			"fail - empty TypeData",
			func() {
				armor := crypto.EncryptArmorPrivKey(priv, "", "eth_secp256k1")
				suite.backend.clientCtx.Keyring.ImportPrivKey("test_key", armor, "")
			},
			from,
			apitypes.TypedData{},
			false,
		},
		// TODO: Generate a TypedData msg
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("case %s", tc.name), func() {
			suite.SetupTest() // reset test and queries
			tc.registerMock()

			responseBz, err := suite.backend.SignTypedData(tc.fromAddr, tc.inputTypedData)

			if tc.expPass {
				sigHash, _, err := apitypes.TypedDataAndHash(tc.inputTypedData)
				suite.Require().NoError(err)
				signature, _, err := suite.backend.clientCtx.Keyring.SignByAddress(
					(sdk.AccAddress)(from.Bytes()),
					sigHash,
					signing.SignMode_SIGN_MODE_TEXTUAL,
				)
				signature[goethcrypto.RecoveryIDOffset] += 27
				suite.Require().NoError(err)
				suite.Require().Equal((hexutil.Bytes)(signature), responseBz)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
