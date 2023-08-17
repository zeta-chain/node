package keeper

import (
	"context"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) WhitelistERC20(goCtx context.Context, msg *types.MsgWhitelistERC20) (*types.MsgWhitelistERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_deploy_fungible_coin) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Deploy can only be executed by the correct policy account")
	}
	erc20Addr := ethcommon.HexToAddress(msg.Erc20Address)
	if erc20Addr == (ethcommon.Address{}) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ERC20 contract address (%s)", msg.Erc20Address)
	}

	// check if it's already whitelisted
	foreignCoins := k.fungibleKeeper.GetAllForeignCoins(ctx)
	for _, fcoin := range foreignCoins {
		assetAddr := ethcommon.HexToAddress(fcoin.Asset)
		if assetAddr == erc20Addr && fcoin.ForeignChainId == msg.ChainId {
			return nil, sdkerrors.Wrapf(types.ErrInvalidAddress, "ERC20 contract address (%s) already whitelisted on chain (%d)", msg.Erc20Address, msg.ChainId)
		}
	}

	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidChainID, "chain id (%d) not supported", msg.ChainId)
	}

	tmpCtx, commit := ctx.CacheContext()
	// add to the foreign coins. Deploy ZRC20 contract for it.
	zrc20Addr, err := k.fungibleKeeper.DeployZRC20Contract(tmpCtx, msg.Name, msg.Symbol, uint8(msg.Decimals), chain.ChainId, common.CoinType_ERC20, msg.Erc20Address, big.NewInt(msg.GasLimit))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDeployContract, "failed to deploy ZRC20 contract for ERC20 contract address (%s) on chain (%d)", msg.Erc20Address, msg.ChainId)
	}
	if zrc20Addr == (ethcommon.Address{}) {
		return nil, sdkerrors.Wrapf(types.ErrDeployContract, "deployed ZRC20 return 0 address for ERC20 contract address (%s) on chain (%d)", msg.Erc20Address, msg.ChainId)
	}

	param, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, msg.ChainId)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidChainID, "core params not found for chain id (%d)", msg.ChainId)
	}

	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, msg.ChainId)
	if !isFound {
		return nil, sdkerrors.Wrapf(types.ErrUnableToGetGasPrice, "median gas price not found for chain id (%d)", msg.ChainId)
	}
	medianGasPrice = medianGasPrice.MulUint64(2) // overpays gas price by 2x

	hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())

	index := crypto.Keccak256Hash(hash.Bytes())
	cctx := types.CrossChainTx{
		Creator:        msg.Creator,
		Index:          index.String(),
		ZetaFees:       sdk.NewUint(0),
		RelayedMessage: fmt.Sprintf("%s:%s", common.CmdWhitelistERC20, msg.Erc20Address),
		CctxStatus: &types.Status{
			Status:              types.CctxStatus_PendingOutbound,
			StatusMessage:       "",
			LastUpdateTimestamp: 0,
		},
		InboundTxParams: &types.InboundTxParams{
			Sender:                          "",
			SenderChainId:                   0,
			TxOrigin:                        "",
			CoinType:                        common.CoinType_Cmd,
			Asset:                           "",
			Amount:                          sdk.Uint{},
			InboundTxObservedHash:           hash.String(), // all Upper case Cosmos TX HEX, with no 0x prefix
			InboundTxObservedExternalHeight: 0,
			InboundTxBallotIndex:            "",
			InboundTxFinalizedZetaHeight:    0,
		},
		OutboundTxParams: []*types.OutboundTxParams{
			{
				Receiver:                         param.Erc20CustodyContractAddress,
				ReceiverChainId:                  msg.ChainId,
				CoinType:                         common.CoinType_Cmd,
				Amount:                           sdk.NewUint(0),
				OutboundTxTssNonce:               0,
				OutboundTxGasLimit:               100_000,
				OutboundTxGasPrice:               medianGasPrice.String(),
				OutboundTxHash:                   "",
				OutboundTxBallotIndex:            "",
				OutboundTxObservedExternalHeight: 0,
			},
		},
	}
	err = k.UpdateNonce(ctx, msg.ChainId, &cctx)
	if err != nil {
		return nil, err
	}

	// add to the foreign coins
	foreignCoin := fungibletypes.ForeignCoins{
		Zrc20ContractAddress: zrc20Addr.Hex(),
		Asset:                msg.Erc20Address,
		ForeignChainId:       msg.ChainId,
		Decimals:             msg.Decimals,
		Name:                 msg.Name,
		Symbol:               msg.Symbol,
		CoinType:             common.CoinType_ERC20,
		GasLimit:             uint64(msg.GasLimit),
	}
	k.fungibleKeeper.SetForeignCoins(ctx, foreignCoin)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	commit()
	return &types.MsgWhitelistERC20Response{}, nil
}
