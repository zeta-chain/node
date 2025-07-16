package v4

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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

// MigrateStore migrates the x/fungible module state from the consensus version 2 to 3
// It updates all existing address in ForeignCoin to use checksum format if the address is EVM type
func MigrateStore(ctx sdk.Context, fungibleKeeper fungibleKeeper) error {
	chainID, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		// Its fine to return nil here and not try to execute the migration at all if the parsing fails
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
		return errors.Wrapf(err, "failed to query SUI gas coin ZRC20 for chain ID %d", chain.ChainId)
	}
	stabilityPoolAddress := types.GasStabilityPoolAddressEVM()

	suiBalance, err := fungibleKeeper.ZRC20BalanceOf(ctx, suiGasZRC20, stabilityPoolAddress)
	if err != nil {
		return errors.Wrapf(err, "failed to get SUI balance for stability pool %s", stabilityPoolAddress)
	}

	err = fungibleKeeper.CallZRC20Burn(ctx, stabilityPoolAddress, suiGasZRC20, suiBalance, true)
	if err != nil {
		return err
	}
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
