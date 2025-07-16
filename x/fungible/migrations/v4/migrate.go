package v4

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethermint "github.com/zeta-chain/ethermint/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/fungible/types"
)

type fungibleKeeper interface {
	GetAllForeignCoins(ctx sdk.Context) (list []types.ForeignCoins)
	QuerySystemContractGasCoinZRC20(ctx sdk.Context, chainid *big.Int) (ethcommon.Address, error)
	ZRC20BalanceOf(ctx sdk.Context, zrc20Address, owner ethcommon.Address) (*big.Int, error)
	CallZRC20Burn(
		ctx sdk.Context,
		sender ethcommon.Address,
		zrc20address ethcommon.Address,
		amount *big.Int,
		noEthereumTxEvent bool,
	) error
}

// MigrateStore migrates the store from consensus version 3 to 4.
// It burns the SUI gas ZRC20 from the stability pool address.
func MigrateStore(ctx sdk.Context, fungibleKeeper fungibleKeeper) error {
	chainID, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		// It's fine to return nil here and not try to execute the migration at all if the parsing fails
		ctx.Logger().Error("failed to parse chain ID", "chain_id", ctx.ChainID(), "error", err)
		return nil
	}
	chain, err := GetSuiChain(chainID.Int64())
	if err != nil {
		ctx.Logger().Error("failed to get Sui chain", "chain_id", chainID.Int64(), "error", err)
		return nil
	}

	suiGasZRC20, err := fungibleKeeper.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chain.ChainId))
	if err != nil {
		ctx.Logger().Error("failed to query SUI gas coin ZRC20 ", "chain_id", chainID.Int64(), "error", err)
		return nil
	}
	stabilityPoolAddress := types.GasStabilityPoolAddressEVM()

	suiBalance, err := fungibleKeeper.ZRC20BalanceOf(ctx, suiGasZRC20, stabilityPoolAddress)
	if err != nil {
		ctx.Logger().Error("failed to get SUI balance for stability pool", "stability_pool_address", stabilityPoolAddress, "error", err)
		return nil
	}

	err = fungibleKeeper.CallZRC20Burn(ctx, stabilityPoolAddress, suiGasZRC20, suiBalance, true)
	if err != nil {
		ctx.Logger().Error("failed to burn SUI gas ZRC20 from stability pool", "stability_pool_address", stabilityPoolAddress, "error", err)
	}

	ctx.Logger().Info("SUI gas ZRC20 burned from stability pool", "stability_pool_address", stabilityPoolAddress, "amount", suiBalance)
	return nil
}

func GetSuiChain(chainID int64) (chains.Chain, error) {
	switch chainID {
	case chains.ZetaChainMainnet.ChainId:
		return chains.SuiMainnet, nil
	case chains.ZetaChainTestnet.ChainId:
		return chains.SuiTestnet, nil
	case chains.ZetaChainPrivnet.ChainId:
		return chains.SuiLocalnet, nil
	default:
		return chains.Chain{}, fmt.Errorf("unsupported chain ID: %d", chainID)
	}
}
