package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"
)

func TestKeeper_ZRC20SetName(t *testing.T) {
	ts := setupChain(t)

	t.Run("should update name", func(t *testing.T) {
		err := ts.fungibleKeeper.ZRC20SetName(ts.ctx, ts.zrc20Address, "NewName")
		require.NoError(t, err)

		name, err := ts.fungibleKeeper.ZRC20Name(ts.ctx, ts.zrc20Address)
		require.NoError(t, err)

		require.Equal(t, "NewName", name)
	})
}

func TestKeeper_ZRC20SetSymbol(t *testing.T) {
	ts := setupChain(t)

	t.Run("should update symbol", func(t *testing.T) {
		err := ts.fungibleKeeper.ZRC20SetSymbol(ts.ctx, ts.zrc20Address, "SYM")
		require.NoError(t, err)

		symbol, err := ts.fungibleKeeper.ZRC20Symbol(ts.ctx, ts.zrc20Address)
		require.NoError(t, err)

		require.Equal(t, "SYM", symbol)
	})
}

func TestKeeper_ZRC20Allowance(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Allowance(
			ts.ctx,
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
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
		)
		require.NoError(t, err)
		require.Equal(t, uint64(0), allowance.Uint64())
	})
}

func TestKeeper_ZRC20BalanceOf(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20BalanceOf(ts.ctx, ts.zrc20Address, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZeroAddress)
	})

	t.Run("should fail when zrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20BalanceOf(ts.ctx, common.Address{}, sample.EthAddress())
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		balance, err := ts.fungibleKeeper.ZRC20BalanceOf(
			ts.ctx,
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
		)
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance.Uint64())
	})
}

func TestKeeper_ZRC20TotalSupply(t *testing.T) {
	ts := setupChain(t)

	t.Run("should fail when zrc20 address is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TotalSupply(ts.ctx, common.Address{})
		require.Error(t, err)
		require.ErrorAs(t, err, &fungibletypes.ErrZRC20ZeroAddress)
	})

	t.Run("should pass with correct input", func(t *testing.T) {
		totalSupply, err := ts.fungibleKeeper.ZRC20TotalSupply(ts.ctx, ts.zrc20Address)
		require.NoError(t, err)
		require.Equal(t, uint64(10000000), totalSupply.Uint64())
	})
}

func TestKeeper_ZRC20Transfer(t *testing.T) {
	ts := setupChain(t)

	// Make sure sample.EthAddress() exists as an evm account in state.
	accAddress := sdk.AccAddress(sample.EthAddress().Bytes())
	acc := ts.fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ts.ctx, accAddress)
	ts.fungibleKeeper.GetAuthKeeper().SetAccount(ts.ctx, acc)

	t.Run("should fail when owner is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20Transfer(
			ts.ctx,
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
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			big.NewInt(10),
		)
		require.NoError(t, err)
		require.True(t, transferred)
	})
}

func TestKeeper_ZRC20TransferFrom(t *testing.T) {
	// Instantiate the ZRC20 ABI only one time.
	// This avoids instantiating it every time deposit or withdraw are called.
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	require.NoError(t, err)

	ts := setupChain(t)

	// Make sure sample.EthAddress() exists as an evm account in state.
	accAddress := sdk.AccAddress(sample.EthAddress().Bytes())
	acc := ts.fungibleKeeper.GetAuthKeeper().NewAccountWithAddress(ts.ctx, accAddress)
	ts.fungibleKeeper.GetAuthKeeper().SetAccount(ts.ctx, acc)

	t.Run("should fail when from is zero address", func(t *testing.T) {
		_, err := ts.fungibleKeeper.ZRC20TransferFrom(
			ts.ctx,
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
			ts.zrc20Address,
			fungibletypes.ModuleAddressEVM,
			sample.EthAddress(),
			fungibletypes.ModuleAddressEVM,
			big.NewInt(10),
		)
		require.Error(t, err)
	})
}
