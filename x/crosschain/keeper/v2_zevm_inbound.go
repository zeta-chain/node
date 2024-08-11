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
	"github.com/zeta-chain/zetacore/pkg/coin"
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
	// try to parse a withdrawal event from the log
	withdrawalEvent, gatewayEvent, err := k.parseGatewayEvent(*log, gatewayAddr)
	if err == nil {
		var inbound *types.MsgVoteInbound

		if withdrawalEvent != nil {
			// find foreign coin object associated to zrc20
			foreignCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, withdrawalEvent.Zrc20.Hex())
			if !found {
				ctx.Logger().
					Info(fmt.Sprintf("cannot find foreign coin with contract address %s", withdrawalEvent.Raw.Address.Hex()))
				return nil
			}

			// validate data of the withdrawal event
			if err := k.validateZRC20Withdrawal(ctx, foreignCoin.ForeignChainId, withdrawalEvent.Value, withdrawalEvent.Receiver); err != nil {
				return err
			}

			// create a new inbound object for the withdrawal
			inbound, err = k.newWithdrawalInbound(ctx, from, txOrigin, foreignCoin, withdrawalEvent)
			if err != nil {
				return err
			}
		} else if gatewayEvent != nil {
			// create a new inbound object for the call
			inbound, err = k.newCallInbound(ctx, from, txOrigin, gatewayEvent)
			if err != nil {
				return err
			}
		}

		if inbound == nil {
			return errors.New("ParseGatewayEvent: invalid log - no event found")
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
	}

	return nil
}

// parseGatewayEvent parses the event from the gateway contract
func (k Keeper) parseGatewayEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
) (*gatewayzevm.GatewayZEVMWithdrawal, *gatewayzevm.GatewayZEVMCall, error) {
	if len(log.Topics) == 0 {
		return nil, nil, errors.New("ParseGatewayCallEvent: invalid log - no topics")
	}
	filterer, err := gatewayzevm.NewGatewayZEVMFilterer(log.Address, bind.ContractFilterer(nil))
	if err != nil {
		return nil, nil, err
	}
	withdrawalEvent, err := k.parseGatewayWithdrawalEvent(log, gatewayAddr, filterer)
	if err == nil {
		return withdrawalEvent, nil, nil
	}
	callEvent, err := k.parseGatewayCallEvent(log, gatewayAddr, filterer)
	if err == nil {
		return nil, callEvent, nil
	}
	return nil, nil, errors.New("ParseGatewayEvent: invalid log - no event found")
}

// parseGatewayWithdrawalEvent parses the GatewayZEVMWithdrawal event from the log
func (k Keeper) parseGatewayWithdrawalEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
	filterer *gatewayzevm.GatewayZEVMFilterer,
) (*gatewayzevm.GatewayZEVMWithdrawal, error) {
	event, err := filterer.ParseWithdrawal(log)
	if err != nil {
		return nil, err
	}
	if event.Raw.Address != gatewayAddr {
		return nil, errors.New("ParseGatewayWithdrawalEvent: invalid log - wrong contract address")
	}
	return event, nil
}

// parseGatewayCallEvent parses the GatewayZEVMCall event from the log
func (k Keeper) parseGatewayCallEvent(
	log ethtypes.Log,
	gatewayAddr ethcommon.Address,
	filterer *gatewayzevm.GatewayZEVMFilterer,
) (*gatewayzevm.GatewayZEVMCall, error) {
	event, err := filterer.ParseCall(log)
	if err != nil {
		return nil, err
	}
	if event.Raw.Address != gatewayAddr {
		return nil, errors.New("ParseGatewayCallEvent: invalid log - wrong contract address")
	}
	return event, nil
}

// newWithdrawalInbound creates a new inbound object for a withdrawal
// currently inbound data is represented with a MsgVoteInbound message
// TODO: replace with a more appropriate object
// https://github.com/zeta-chain/node/issues/2658
// TODO: include revert options
// https://github.com/zeta-chain/node/issues/2660
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

	// temporary hardcoded gas limit for withdraw and call
	// TODO: use gas limit specified by user
	// https://github.com/zeta-chain/node/issues/2668
	gasLimit := uint64(1000000)

	// for simple withdraw without call, we use the specified gas limit in the zrc20 contract
	if len(event.Message) == 0 {
		gasLimitQueried, err := k.fungibleKeeper.QueryGasLimit(
			ctx,
			ethcommon.HexToAddress(foreignCoin.Zrc20ContractAddress),
		)
		if err != nil {
			return nil, errors.Wrap(err, "cannot query gas limit")
		}
		gasLimit = gasLimitQueried.Uint64()
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
		gasLimit,
		foreignCoin.CoinType,
		foreignCoin.Asset,
		event.Raw.Index,
		types.ProtocolContractVersion_V2,
	), nil
}

// newCallInbound creates a new inbound object for a call
// currently inbound data is represented with a MsgVoteInbound message
// TODO: replace with a more appropriate object
// https://github.com/zeta-chain/node/issues/2658
// TODO: include revert options
// https://github.com/zeta-chain/node/issues/2660
func (k Keeper) newCallInbound(
	ctx sdk.Context,
	from ethcommon.Address,
	txOrigin string,
	event *gatewayzevm.GatewayZEVMCall,
) (*types.MsgVoteInbound, error) {
	receiverChainID := event.ChainId.Int64()
	receiverChain, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, receiverChainID)
	if !found {
		return nil, errorsmod.Wrapf(
			observertypes.ErrSupportedChains,
			"chain with chainID %d not supported",
			receiverChainID,
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

	// temporary hardcoded gas limit for withdraw and call
	// TODO: use gas limit specified by user
	// https://github.com/zeta-chain/node/issues/2668
	gasLimit := uint64(1000000)

	return types.NewMsgVoteInbound(
		"",
		from.Hex(),
		senderChain.ChainId,
		txOrigin,
		toAddr,
		receiverChainID,
		math.ZeroUint(),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.String(),
		event.Raw.BlockNumber,
		gasLimit,
		coin.CoinType_NoAssetCall,
		"",
		event.Raw.Index,
		types.ProtocolContractVersion_V2,
	), nil
}
