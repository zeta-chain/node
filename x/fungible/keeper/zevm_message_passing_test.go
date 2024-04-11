package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/x/evm/statedb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/contracts"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_ZEVMDepositAndCallContract(t *testing.T) {
	t.Run("successfully call ZEVMDepositAndCallContract on connector contract ", func(t *testing.T) {
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

		_, err = k.ZEVMDepositAndCallContract(ctx, zetaTxSender, zetaTxReceiver, inboundSenderChainID, inboundAmount, data, cctxIndexBytes)
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

		_, err = k.ZEVMDepositAndCallContract(ctx, zetaTxSender, zetaTxReceiver, inboundSenderChainID, inboundAmount, data, cctxIndexBytes)
		require.NoError(t, err)
		b := sdkk.BankKeeper.GetBalance(ctx, sdk.AccAddress(zetaTxReceiver.Bytes()), config.BaseDenom)
		require.Equal(t, inboundAmount.Int64(), b.Amount.Int64())
	})

	t.Run("fail ZEVMDepositAndCallContract if account not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		zetaTxSender := sample.EthAddress()
		zetaTxReceiver := sample.EthAddress()
		inboundSenderChainID := int64(1)
		inboundAmount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err := k.ZEVMDepositAndCallContract(ctx, zetaTxSender, zetaTxReceiver, inboundSenderChainID, inboundAmount, data, cctxIndexBytes)
		require.ErrorIs(t, err, types.ErrAccountNotFound)
		require.ErrorContains(t, err, "account not found")
	})

	t.Run("fail ZEVMDepositAndCallContract id Deposit Fails", func(t *testing.T) {
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

		_, err = k.ZEVMDepositAndCallContract(ctx, zetaTxSender, zetaTxReceiver, inboundSenderChainID, inboundAmount, data, cctxIndexBytes)
		require.ErrorIs(t, err, errorMint)
	})
}
func TestKeeper_ZevmOnReceive(t *testing.T) {
	t.Run("successfully call ZevmOnReceive on connector contract ", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, contracts.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		zetaTxSender := sample.EthAddress().Bytes()
		senderChainID := big.NewInt(1)
		zetaTxReceiver := dAppContract
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.ZevmOnReceive(ctx, zetaTxSender, zetaTxReceiver, senderChainID, amount, data, cctxIndexBytes)
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
		require.Equal(t, zetaTxSender, valSenderAddress)
	})

	t.Run("fail to call ZevmOnReceive if CallOnReceiveZevmConnector fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		dAppContract, err := k.DeployContract(ctx, contracts.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		zetaTxSender := sample.EthAddress().Bytes()
		senderChainID := big.NewInt(1)
		zetaTxReceiver := dAppContract
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, err = k.ZevmOnReceive(ctx, zetaTxSender, zetaTxReceiver, senderChainID, amount, data, cctxIndexBytes)
		require.ErrorIs(t, err, types.ErrContractNotFound)
		require.ErrorContains(t, err, "GetSystemContract address not found")
	})
}

func TestKeeper_ZevmOnRevert(t *testing.T) {
	t.Run("successfully call ZevmOnRevert on connector contract ", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		dAppContract, err := k.DeployContract(ctx, contracts.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		zetaTxSender := dAppContract
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		zetaTxReceiver := sample.EthAddress().Bytes()
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, isContract, err := k.ZevmOnRevert(ctx, zetaTxSender, zetaTxReceiver, senderChainID, destinationChainID, amount, data, cctxIndexBytes)
		require.NoError(t, err)
		require.True(t, isContract)

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

	t.Run("fail to call ZevmOnRevert if account is not a contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		zetaTxSender := sample.EthAddress()
		err := sdkk.EvmKeeper.SetAccount(ctx, zetaTxSender, statedb.Account{
			Nonce:    0,
			Balance:  big.NewInt(100),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)

		_, isContract, err := k.ZevmOnRevert(ctx, zetaTxSender, sample.EthAddress().Bytes(),
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(45),
			[]byte("message"),
			[32]byte{})
		require.ErrorIs(t, err, types.ErrCallNonContract)
		require.False(t, isContract)
	})

	t.Run("fail to call ZevmOnRevert if CallOnRevertZevmConnector fails", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		dAppContract, err := k.DeployContract(ctx, contracts.DappMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, dAppContract)

		zetaTxSender := dAppContract
		senderChainID := big.NewInt(1)
		destinationChainID := big.NewInt(2)
		zetaTxReceiver := sample.EthAddress().Bytes()
		amount := big.NewInt(45)
		data := []byte("message")
		cctxIndexBytes := [32]byte{}

		_, isContract, err := k.ZevmOnRevert(ctx, zetaTxSender, zetaTxReceiver, senderChainID, destinationChainID, amount, data, cctxIndexBytes)
		require.ErrorIs(t, err, types.ErrContractNotFound)
		require.ErrorContains(t, err, "GetSystemContract address not found")
		require.True(t, isContract)
	})

	t.Run("fail to call ZevmOnRevert if account not found for sender address", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		_, isContract, err := k.ZevmOnRevert(ctx, sample.EthAddress(),
			sample.EthAddress().Bytes(),
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(45),
			[]byte("message"),
			[32]byte{})
		require.ErrorIs(t, err, types.ErrAccountNotFound)
		require.False(t, isContract)
	})

}
