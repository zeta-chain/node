package keeper

import (
	"context"
	"math/big"

	cosmoserrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/zetacore/pkg/coin"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// DeployFungibleCoinZRC20 deploys a fungible coin from a connected chains as a ZRC20 on ZetaChain.
//
// If this is a gas coin, the following happens:
//
// * ZRC20 contract for the coin is deployed
// * contract address of ZRC20 is set as a token address in the system
// contract
// * ZETA tokens are minted and deposited into the module account
// * setGasZetaPool is called on the system contract to add the information
// about the pool to the system contract
// * addLiquidityETH is called to add liquidity to the pool
//
// If this is a non-gas coin, the following happens:
//
// * ZRC20 contract for the coin is deployed
// * The coin is added to the list of foreign coins in the module's state
//
// Authorized: admin policy group 2.
func (k msgServer) DeployFungibleCoinZRC20(
	goCtx context.Context,
	msg *types.MsgDeployFungibleCoinZRC20,
) (*types.MsgDeployFungibleCoinZRC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var address common.Address
	var err error

	if err = msg.ValidateBasic(); err != nil {
		return nil, err
	}

	err = k.GetAuthorityKeeper().CheckAuthorization(ctx, msg)
	if err != nil {
		return nil, cosmoserrors.Wrap(authoritytypes.ErrUnauthorized, err.Error())
	}

	if msg.CoinType == coin.CoinType_Gas {
		// #nosec G115 always in range
		address, err = k.SetupChainGasCoinAndPool(
			ctx,
			msg.ForeignChainId,
			msg.Name,
			msg.Symbol,
			uint8(msg.Decimals),
			big.NewInt(msg.GasLimit),
		)
		if err != nil {
			return nil, cosmoserrors.Wrapf(err, "failed to setupChainGasCoinAndPool")
		}
	} else {
		// #nosec G115 always in range
		address, err = k.DeployZRC20Contract(ctx, msg.Name, msg.Symbol, uint8(msg.Decimals), msg.ForeignChainId, msg.CoinType, msg.ERC20, big.NewInt(msg.GasLimit))
		if err != nil {
			return nil, err
		}
	}

	err = ctx.EventManager().EmitTypedEvent(
		&types.EventZRC20Deployed{
			MsgTypeUrl: sdk.MsgTypeURL(&types.MsgDeployFungibleCoinZRC20{}),
			ChainId:    msg.ForeignChainId,
			Contract:   address.String(),
			Name:       msg.Name,
			Symbol:     msg.Symbol,
			// #nosec G115 always in range
			Decimals: int64(msg.Decimals),
			CoinType: msg.CoinType,
			Erc20:    msg.ERC20,
			GasLimit: msg.GasLimit,
		},
	)
	if err != nil {
		return nil, cosmoserrors.Wrapf(err, "failed to emit event")
	}

	return &types.MsgDeployFungibleCoinZRC20Response{
		Address: address.Hex(),
	}, nil
}
