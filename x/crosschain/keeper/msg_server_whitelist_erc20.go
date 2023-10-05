package keeper

import (
	"context"
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// WhitelistERC20 deploys a new zrc20, create a foreign coin object for the ERC20
// and emit a crosschain tx to whitelist the ERC20 on the external chain
func (k Keeper) WhitelistERC20(goCtx context.Context, msg *types.MsgWhitelistERC20) (*types.MsgWhitelistERC20Response, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if msg.Creator != k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_group1) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "Deploy can only be executed by the correct policy account")
	}
	erc20Addr := ethcommon.HexToAddress(msg.Erc20Address)
	if erc20Addr == (ethcommon.Address{}) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ERC20 contract address (%s)", msg.Erc20Address)
	}

	// check if the erc20 is already whitelisted
	foreignCoins := k.fungibleKeeper.GetAllForeignCoins(ctx)
	for _, fCoin := range foreignCoins {
		assetAddr := ethcommon.HexToAddress(fCoin.Asset)
		if assetAddr == erc20Addr && fCoin.ForeignChainId == msg.ChainId {
			return nil, errorsmod.Wrapf(
				types.ErrInvalidAddress,
				"ERC20 contract address (%s) already whitelisted on chain (%d)",
				msg.Erc20Address,
				msg.ChainId,
			)
		}
	}

	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCannotFindTSSKeys, "Cannot create new admin cmd of type whitelistERC20")
	}

	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidChainID, "chain id (%d) not supported", msg.ChainId)
	}

	// use a temporary context for the zrc20 deployment
	tmpCtx, commit := ctx.CacheContext()

	// add to the foreign coins. Deploy ZRC20 contract for it.
	zrc20Addr, err := k.fungibleKeeper.DeployZRC20Contract(
		tmpCtx,
		msg.Name,
		msg.Symbol,
		// #nosec G701 always in range
		uint8(msg.Decimals),
		chain.ChainId,
		common.CoinType_ERC20,
		msg.Erc20Address,
		big.NewInt(msg.GasLimit),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrDeployContract,
			"failed to deploy ZRC20 contract for ERC20 contract address (%s) on chain (%d)",
			msg.Erc20Address,
			msg.ChainId,
		)
	}
	if zrc20Addr == (ethcommon.Address{}) {
		return nil, errorsmod.Wrapf(
			types.ErrDeployContract,
			"deployed ZRC20 return 0 address for ERC20 contract address (%s) on chain (%d)",
			msg.Erc20Address,
			msg.ChainId,
		)
	}

	// get necessary parameters to create the cctx
	param, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, msg.ChainId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidChainID, "core params not found for chain id (%d)", msg.ChainId)
	}
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, msg.ChainId)
	if !isFound {
		return nil, errorsmod.Wrapf(types.ErrUnableToGetGasPrice, "median gas price not found for chain id (%d)", msg.ChainId)
	}
	medianGasPrice = medianGasPrice.MulUint64(2) // overpays gas price by 2x

	// calculate the cctx index
	// we use the deployed zrc20 contract address to generate a unique index
	// since other parts of the system may use the zrc20 for the index, we add a message specific suffix
	hash := crypto.Keccak256Hash(zrc20Addr.Bytes(), []byte("WhitelistERC20"))
	index := hash.Hex()

	// create a cmd cctx to whitelist the erc20 on the external chain
	cctx := types.CrossChainTx{
		Creator:        msg.Creator,
		Index:          index,
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
			Amount:                          math.Uint{},
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
				Amount:                           math.NewUint(0),
				OutboundTxTssNonce:               0,
				OutboundTxGasLimit:               100_000,
				OutboundTxGasPrice:               medianGasPrice.String(),
				OutboundTxHash:                   "",
				OutboundTxBallotIndex:            "",
				OutboundTxObservedExternalHeight: 0,
				TssPubkey:                        tss.TssPubkey,
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
		// #nosec G701 always positive
		GasLimit: uint64(msg.GasLimit),
	}
	k.fungibleKeeper.SetForeignCoins(ctx, foreignCoin)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)

	commit()

	return &types.MsgWhitelistERC20Response{}, nil
}
