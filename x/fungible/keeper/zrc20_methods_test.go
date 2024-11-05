package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"
)

func Test_ZRC20Allowance(t *testing.T) {
	// Instantiate the ZRC20 ABI only one time.
	// This avoids instantiating it every time deposit or withdraw are called.
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	t.Run("should fail when ZRC20ABI is nil", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Allowance(ts.ctx, nil, ts.zrc20Address, common.Address{}, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20NilABI)
	})

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Allowance(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			common.Address{},
			sample.EthAddress(),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when spender is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Allowance(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			sample.EthAddress(),
			common.Address{},
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when zrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Allowance(
			ts.ctx,
			zrc20ABI,
			common.Address{},
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		allowance, err := ts.fungibleKeeper.ZRC20Allowance(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
		)
		require.NoError(t, err)
		require.Equal(t, uint64(0), allowance.Uint64())
	})
}

func Test_ZRC20BalanceOf(t *testing.T) {
	// Instantiate the ZRC20 ABI only one time.
	// This avoids instantiating it every time deposit or withdraw are called.
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	t.Run("should fail when ZRC20ABI is nil", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20BalanceOf(ts.ctx, nil, ts.zrc20Address, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20NilABI)
	})

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20BalanceOf(ts.ctx, zrc20ABI, ts.zrc20Address, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when zrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20BalanceOf(ts.ctx, zrc20ABI, common.Address{}, sample.EthAddress())
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		balance, err := ts.fungibleKeeper.ZRC20BalanceOf(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
		)
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Uint64())
	})
}

func Test_ZRC20TotalSupply(t *testing.T) {
	// Instantiate the ZRC20 ABI only one time.
	// This avoids instantiating it every time deposit or withdraw are called.
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	t.Run("should fail when ZRC20ABI is nil", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TotalSupply(ts.ctx, nil, ts.zrc20Address)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20NilABI)
	})

	t.Run("should fail when zrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TotalSupply(ts.ctx, zrc20ABI, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		totalSupply, err := ts.fungibleKeeper.ZRC20TotalSupply(ts.ctx, zrc20ABI, ts.zrc20Address)
		require.NoError(t, err)
		require.Equal(t, uint64(10000000), totalSupply.Uint64())
	})
}

func Test_ZRC20Transfer(t *testing.T) {
	// Instantiate the ZRC20 ABI only one time.
	// This avoids instantiating it every time deposit or withdraw are called.
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	// Make sure sample.EthAddress() exists as an ethermint account in state.
	accAddress := sdk.AccAddress(sample.EthAddress().Bytes())
	ts.fungibleKeeper.GetAccountKeeper().SetAccount(ts.ctx, authtypes.NewBaseAccount(accAddress, nil, 0, 0))

	t.Run("should fail when ZRC20ABI is nil", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Transfer(
			ts.ctx,
			nil,
			ts.zrc20Address,
			common.Address{},
			common.Address{},
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20NilABI)
	})

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Transfer(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			common.Address{},
			sample.EthAddress(),
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when spender is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Transfer(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			sample.EthAddress(),
			common.Address{},
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when zrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Transfer(
			ts.ctx,
			zrc20ABI,
			common.Address{},
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, fungibletypes.ModuleAddressEVM, big.NewInt(10))
		transferred, err := ts.fungibleKeeper.ZRC20Transfer(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			big.NewInt(10),
		)
		require.NoError(t, err)
		require.True(t, transferred)
	})
}

func Test_ZRC20TransferFrom(t *testing.T) {
	// Instantiate the ZRC20 ABI only one time.
	// This avoids instantiating it every time deposit or withdraw are called.
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	// Make sure sample.EthAddress() exists as an ethermint account in state.
	accAddress := sdk.AccAddress(sample.EthAddress().Bytes())
	ts.fungibleKeeper.GetAccountKeeper().SetAccount(ts.ctx, authtypes.NewBaseAccount(accAddress, nil, 0, 0))

	t.Run("should fail when ZRC20ABI is nil", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
			nil,
			ts.zrc20Address,
			common.Address{},
			common.Address{},
			common.Address{},
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20NilABI)
	})

	t.Run("should fail when from is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			sample.EthAddress(),
			common.Address{},
			sample.EthAddress(),
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when to is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			sample.EthAddress(),
			sample.EthAddress(),
			common.Address{},
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when spender is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			common.Address{},
			sample.EthAddress(),
			sample.EthAddress(),
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when zrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
			zrc20ABI,
			common.Address{},
			sample.EthAddress(),
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(0),
		)
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20ZeroAddress)
	})

	t.Run("should fail without an allowance approval", func(t *testing.T) {
		// Deposit ZRC20 into fungible EOA.
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, fungibletypes.ModuleAddressEVM, big.NewInt(1000))

		// Transferring the tokens with transferFrom without approval should fail.
		_, err = ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(10),
		)
		require.Error(t, err)
	})

	t.Run("should success with an allowance approval", func(t *testing.T) {
		// Deposit ZRC20 into fungible EOA.
		ts.fungibleKeeper.DepositZRC20(ts.ctx, ts.zrc20Address, fungibletypes.ModuleAddressEVM, big.NewInt(1000))

		// Approve allowance to sample.EthAddress() to spend 10 ZRC20 tokens.
		approveAllowance(t, ts, zrc20ABI, fungibletypes.ModuleAddressEVM, sample.EthAddress(), big.NewInt(10))

		// Transferring the tokens with transferFrom without approval should fail.
		_, err = ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
			zrc20ABI,
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(10),
		)
		require.Error(t, err)
	})
}
