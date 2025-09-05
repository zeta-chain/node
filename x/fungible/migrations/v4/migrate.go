package v4

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

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
	logAndSkip := func(msg string, err error, fields ...interface{}) error {
		ctx.Logger().Error(msg, append(fields, "error", err)...)
		return nil
	}

	zetachain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return logAndSkip("failed to parse chain ID", err, "chain_id", ctx.ChainID())
	}

	chain, err := GetSuiChain(zetachain.ChainId)
	if err != nil {
		return logAndSkip("failed to get Sui chain", err, "chain_id", zetachain.ChainId)
	}

	suiGasZRC20, err := fungibleKeeper.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chain.ChainId))
	if err != nil {
		return logAndSkip("failed to query SUI gas coin ZRC20 ", err, "chain_id", zetachain.ChainId)
	}

	stabilityPoolAddress := types.GasStabilityPoolAddressEVM()

	suiBalance, err := fungibleKeeper.ZRC20BalanceOf(ctx, suiGasZRC20, stabilityPoolAddress)
	if err != nil {
		return logAndSkip(
			"failed to get SUI balance for stability pool",
			err,
			"stability_pool_address",
			stabilityPoolAddress,
		)
	}

	err = fungibleKeeper.CallZRC20Burn(ctx, stabilityPoolAddress, suiGasZRC20, suiBalance, true)
	if err != nil {
		return logAndSkip(
			"failed to burn SUI gas ZRC20 from stability pool",
			err,
			"stability_pool_address",
			stabilityPoolAddress,
		)
	}

	ctx.Logger().
		Info("SUI gas ZRC20 burned from stability pool", "stability_pool_address", stabilityPoolAddress, "amount", suiBalance)
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
