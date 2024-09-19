package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/ethermint/x/evm/statedb"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/testutil/contracts"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_ZEVMDepositAndCallContract(t *testing.T) {
	t.Run("successfully call ZETADepositAndCallContract on connector contract ", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, contracts.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := dAppContract
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.ZETADepositAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)

		dappAbi, err := contracts.DappMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*dappAbi,
			types.ModuleAddressEVM,
			dAppContract,
			big.NewInt(0),
			nil,
			false,
			false,
			"zetaTxSenderAddress",
		)
		require.NoError(t, err)
		unpacked, err := dappAbi.Unpack("zetaTxSenderAddress", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		valSenderAddress, ok := unpacked[0].([]byte)
		require.True(t, ok)
		require.Equal(t, zetaTxSender.Bytes(), valSenderAddress)
	})

	t.Run("successfully deposit coin if account is not a contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := sample.EthAddress()
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, zetaTxReceiver, statedb.Account{
			Nonce:    0,
			Balance:  big.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)

		_, err = k.ZETADepositAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(zetaTxReceiver.Bytes()), config.BaseDenom)
		require.Equal(t, inboundAmount.Int64(), b.Amount.Int64())
	})

	t.Run("automatically deposit coin  if account not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := sample.EthAddress()
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err := k.ZETADepositAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(zetaTxReceiver.Bytes()), config.BaseDenom)
		require.Equal(t, inboundAmount.Int64(), b.Amount.Int64())
	})

	t.Run("fail ZETADepositAndCallContract if Deposit Fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{UseBankMock: true})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		bankMock := keepertest.GetFungibleBankMock(t, k)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := sample.EthAddress()
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, zetaTxReceiver, statedb.Account{
			Nonce:    0,
			Balance:  big.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)
		errorMint := errors.New("", 10, "error minting coins")
		bankMock.On("MintCoins", ctx, types.ModuleName, mock.Anything).Return(errorMint).Once()

		_, err = k.ZETADepositAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			inboundSenderChainID,
			inboundAmount,
			data,
			cctxIndexBytes,
		)
		require.ErrorIs(t, err, errorMint)
	})
}

func TestKeeper_ZEVMRevertAndCallContract(t *testing.T) {
	t.Run("successfully call ZETARevertAndCallContract if receiver is a contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, contracts.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		zetaTxSender := dAppContract
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		zetaTxReceiver := sample.EthAddress()
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.ZETARevertAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)

		dappAbi, err := contracts.DappMetaData.GetAbi()
		require.NoError(t, err)
		res, err := k.CallEVM(
			ctx,
			*dappAbi,
			types.ModuleAddressEVM,
			dAppContract,
			big.NewInt(0),
			nil,
			false,
			false,
			"zetaTxSenderAddress",
		)
		require.NoError(t, err)
		unpacked, err := dappAbi.Unpack("zetaTxSenderAddress", res.Ret)
		require.NoError(t, err)
		require.NotZero(t, len(unpacked))
		valSenderAddress, ok := unpacked[0].([]byte)
		require.True(t, ok)
		require.Equal(t, zetaTxSender.Bytes(), valSenderAddress)
	})

	t.Run("successfully deposit coin if account is not a contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, zetaTxSender, statedb.Account{
			Nonce:    0,
			Balance:  big.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)

		_, err = k.ZETARevertAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(zetaTxSender.Bytes()), config.BaseDenom)
		require.Equal(t, amount.Int64(), b.Amount.Int64())
	})

	t.Run("automatically deposit coin if account not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err := k.ZETARevertAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(zetaTxSender.Bytes()), config.BaseDenom)
		require.Equal(t, amount.Int64(), b.Amount.Int64())
	})

	t.Run("fail ZETARevertAndCallContract if Deposit Fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{UseBankMock: true})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		bankMock := keepertest.GetFungibleBankMock(t, k)
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		err := sdkk.EvmKeeper.SetAccount(ctx, zetaTxSender, statedb.Account{
			Nonce:    0,
			Balance:  big.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)
		errorMint := errors.New("", 101, "error minting coins")
		bankMock.On("MintCoins", ctx, types.ModuleName, mock.Anything).Return(errorMint).Once()

		_, err = k.ZETARevertAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.ErrorIs(t, err, errorMint)
	})

	t.Run("fail ZETARevertAndCallContract if ZevmOnRevert fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		dAppContract, err := k.DeployContract(ctx, contracts.DappReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		zetaTxSender := dAppContract
		zetaTxReceiver := sample.EthAddress()
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.ZETARevertAndCallContract(
			ctx,
			zetaTxSender,
			zetaTxReceiver,
			senderChainID.Int64(),
			destinationChainID.Int64(),
			amount,
			data,
			cctxIndexBytes,
		)
		require.ErrorIs(t, err, types.ErrContractNotFound)
		require.ErrorContains(t, err, "GetSystemContract address not found")
	})
}
