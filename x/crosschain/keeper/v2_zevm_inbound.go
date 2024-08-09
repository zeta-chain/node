package keeper

import (
	"encoding/hex"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ProcessZEVMInboundV2 processes the logs emitted by the zEVM contract for V2 protocol contracts
// it parses logs from GatewayZEVM contract and updates the crosschain module state
func (k Keeper) ProcessZEVMInboundV2(
	ctx sdk.Context,
	log *ethtypes.Log,
	gatewayAddr,
	from ethcommon.Address,
	txOrigin string,
) error {
	withdrawalEvent, err := k.parseGatewayWithdrawalEvent(*log, gatewayAddr)
	if err != nil {
		return err
	}

	// find foreign coin object associated to zrc20
	coin, foundCoin := k.fungibleKeeper.GetForeignCoins(ctx, withdrawalEvent.Zrc20.Hex())
	if !foundCoin {
		ctx.Logger().
			Info(fmt.Sprintf("cannot find foreign coin with contract address %s", withdrawalEvent.Raw.Address.Hex()))
		return nil
	}

	// validate data of the withdrawal event
	if err := k.validateZRC20Withdrawal(ctx, coin.ForeignChainId, withdrawalEvent.Value, withdrawalEvent.Receiver); err != nil {
		return err
	}

	// create a new inbound object for the withdrawal
	inbound, err := k.newWithdrawalInbound(ctx, from, txOrigin, coin, withdrawalEvent)
	if err != nil {
		return err
	}

	// validate inbound for processing
	cctx, err := k.ValidateInbound(ctx, inbound, false)
	if err != nil {
		return err
	}
	if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
		return errors.New("cctx aborted")
	}
	EmitZRCWithdrawCreated(ctx, *cctx)

	return nil
}

// parseGatewayWithdrawalEvent parses the GatewayZEVMWithdrawal event from the log
func (k Keeper) parseGatewayWithdrawalEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
) (*gatewayzevm.GatewayZEVMWithdrawal, error) {
	filterer, err := gatewayzevm.NewGatewayZEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, err
	}
	if len(log.Topics) == 0 {
		return nil, errors.New("ParseGatewayWithdrawalEvent: invalid log - no topics")
	}
	event, err := filterer.ParseWithdrawal(log)
	if err != nil {
		return nil, err
	}

	if event.Raw.Address != gatewayAddr {
		return nil, errors.New("ParseGatewayWithdrawalEvent: invalid log - wrong contract address")
	}

	return event, nil
}

// newWithdrawalInbound creates a new inbound object for a withdrawal
// currently inbound data is represented with a MsgVoteInbound message
// TODO: replace with a more appropriate object
// https://github.com/zeta-chain/node/issues/2658
func (k Keeper) newWithdrawalInbound(
	ctx sdk.Context,
	from ethcommon.Address,
	txOrigin string,
	foreignCoin fungibletypes.ForeignCoins,
	event *gatewayzevm.GatewayZEVMWithdrawal,
) (*types.MsgVoteInbound, error) {
	receiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, foreignCoin.ForeignChainId)
	if !found {
		return nil, errorsmod.Wrapf(
			observertypes.ErrSupportedChains,
			"chain with chainID %d not supported",
			foreignCoin.ForeignChainId,
		)
	}

	senderChain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
	if err != nil {
		return nil, errors.Wrapf(err, "ProcessZEVMInboundV2: failed to convert chainID %s", ctx.ChainID())
	}

	toAddr, err := receiverChain.EncodeAddress(event.Receiver)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot encode address %v", event.Receiver)
	}

	// TODO: emit gas limit in the withdrawal event
	// https://github.com/zeta-chain/node/issues/2658
	gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress))
	if err != nil {
		return nil, errors.Wrap(err, "cannot query gas limit")
	}

	return types.NewMsgVoteInbound(
		"",
		from.Hex(),
		senderChain.ChainId,
		txOrigin,
		toAddr,
		foreignCoin.ForeignChainId,
		math.NewUintFromBigInt(event.Value),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit.Uint64(),
		foreignCoin.CoinType,
		foreignCoin.Asset,
		event.Raw.Index,
		types.ProtocolContractVersion_V2,
	), nil
}
