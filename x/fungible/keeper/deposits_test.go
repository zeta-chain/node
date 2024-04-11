package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/contracts"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_ZRC20DepositAndCallContract(t *testing.T) {
	t.Run("can deposit gas coin for transfers", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("can deposit non-gas coin for transfers", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := deployZRC20(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", assetAddress, "foobar")

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_ERC20,
			assetAddress,
		)
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)
	})

	t.Run("should fail if trying to call a contract with data to a EOC", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		deployZRC20(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", assetAddress, "foobar")

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte("DEADBEEF"),
			coin.CoinType_ERC20,
			assetAddress,
		)
		require.ErrorIs(t, err, types.ErrCallNonContract)
	})

	t.Run("can deposit coin for transfers with liquidity cap not reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// there is an initial total supply minted during gas pool setup
		initialTotalSupply, err := k.TotalSupplyZRC4(ctx, zrc20)
		require.NoError(t, err)

		// set a liquidity cap
		foreignCoin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		foreignCoin.LiquidityCap = math.NewUint(initialTotalSupply.Uint64() + 1000)
		k.SetForeignCoins(ctx, foreignCoin)

		// increase total supply
		_, err = k.DepositZRC20(ctx, zrc20, sample.EthAddress(), big.NewInt(500))
		require.NoError(t, err)

		// deposit
		to := sample.EthAddress()
		_, contractCall, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(500),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.NoError(t, err)
		require.False(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, to)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(500), balance)
	})

	t.Run("should fail if coin paused", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// pause the coin
		foreignCoin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		foreignCoin.Paused = true
		k.SetForeignCoins(ctx, foreignCoin)

		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.ErrorIs(t, err, types.ErrPausedZRC20)
	})

	t.Run("should fail if liquidity cap reached", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		// there is an initial total supply minted during gas pool setup
		initialTotalSupply, err := k.TotalSupplyZRC4(ctx, zrc20)
		require.NoError(t, err)

		// set a liquidity cap
		foreignCoin, found := k.GetForeignCoins(ctx, zrc20.String())
		require.True(t, found)
		foreignCoin.LiquidityCap = math.NewUint(initialTotalSupply.Uint64() + 1000)
		k.SetForeignCoins(ctx, foreignCoin)

		// increase total supply
		_, err = k.DepositZRC20(ctx, zrc20, sample.EthAddress(), big.NewInt(500))
		require.NoError(t, err)

		// deposit (500 + 501 > 1000)
		to := sample.EthAddress()
		_, _, err = k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(501),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.ErrorIs(t, err, types.ErrForeignCoinCapReached)
	})

	t.Run("should fail if gas coin not found", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.ErrorIs(t, err, crosschaintypes.ErrGasCoinNotFound)
	})

	t.Run("should fail if zrc20 not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId
		assetAddress := sample.EthAddress().String()

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deposit
		to := sample.EthAddress()
		_, _, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			to,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_ERC20,
			assetAddress,
		)
		require.ErrorIs(t, err, crosschaintypes.ErrForeignCoinNotFound)
	})

	t.Run("should return contract call if receiver is a contract", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		example, err := k.DeployContract(ctx, contracts.ExampleMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, example)

		// deposit
		_, contractCall, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			example,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.NoError(t, err)
		require.True(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, example)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(42), balance)

		// check onCrossChainCall() hook was called
		assertExampleBarValue(t, ctx, k, example, 42)
	})

	t.Run("should fail if call contract fails", func(t *testing.T) {
		// setup gas coin
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		chainList := chains.DefaultChainsList()
		chain := chainList[0].ChainId

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chain, "foobar", "foobar")

		reverter, err := k.DeployContract(ctx, contracts.ReverterMetaData)
		require.NoError(t, err)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, reverter)

		// deposit
		_, contractCall, err := k.ZRC20DepositAndCallContract(
			ctx,
			sample.EthAddress().Bytes(),
			reverter,
			big.NewInt(42),
			chain,
			[]byte{},
			coin.CoinType_Gas,
			sample.EthAddress().String(),
		)
		require.Error(t, err)
		require.True(t, contractCall)

		balance, err := k.BalanceOfZRC4(ctx, zrc20, reverter)
		require.NoError(t, err)
		require.EqualValues(t, int64(0), balance.Int64())
	})
}

func TestKeeper_DepositCoinZeta(t *testing.T) {
	k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
	to := sample.EthAddress()
	amount := big.NewInt(1)
	zetaToAddress := sdk.AccAddress(to.Bytes())

	b := sdkk.BankKeeper.GetBalance(ctx, zetaToAddress, config.BaseDenom)
	require.Equal(t, int64(0), b.Amount.Int64())

	err := k.DepositCoinZeta(ctx, to, amount)
	require.NoError(t, err)
	b = sdkk.BankKeeper.GetBalance(ctx, zetaToAddress, config.BaseDenom)
	require.Equal(t, amount.Int64(), b.Amount.Int64())
}

func Test_AddressConvertion(t *testing.T) {
	addressCosmosString := sample.AccAddress()
	addressCosmsosAccAddress := sdk.MustAccAddressFromBech32(addressCosmosString)
	// Logic used in depositCoins function
	addressEth := ethcommon.HexToAddress(addressCosmosString)
	//https://github.com/zeta-chain/zeta-node/blob/zevm-message-passing/x/fungible/keeper/deposits.go#L17-L17
	depositAddress := sdk.AccAddress(addressEth.Bytes())
	depositAddressString := depositAddress.String()

	require.Equal(t, addressCosmsosAccAddress, depositAddress)
	require.Equal(t, addressCosmosString, depositAddressString)
}
