package keeper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/memo"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const InCCTXIndexKey = "inCctxIndex"

// HandleEVMDeposit handles a deposit from an inbound tx
// returns (isContractReverted, err)
// (true, non-nil) means CallEVM() reverted
func (k Keeper) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx) (bool, error) {
	to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
	sender := ethcommon.HexToAddress(cctx.InboundParams.Sender)
	var ethTxHash ethcommon.Hash
	inboundAmount := cctx.GetInboundParams().Amount.BigInt()
	inboundSender := cctx.GetInboundParams().Sender
	inboundSenderChainID := cctx.GetInboundParams().SenderChainId
	inboundCoinType := cctx.InboundParams.CoinType

	if len(ctx.TxBytes()) > 0 {
		// add event for tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		ethTxHash = ethcommon.BytesToHash(hash)
		cctx.GetCurrentOutboundParam().Hash = ethTxHash.String()
		// #nosec G115 always positive
		cctx.GetCurrentOutboundParam().ObservedExternalHeight = uint64(ctx.BlockHeight())
	}
	// Refactor HandleEVMDeposit to have two clear branches of logic for V1 and V2 protocol versions
	// TODO : https://github.com/zeta-chain/node/issues/3988
	if inboundCoinType == coin.CoinType_Zeta && cctx.ProtocolContractVersion == types.ProtocolContractVersion_V1 {
		// In case of an error
		// 	- Return true will revert the cctx and create a revert cctx with status PendingRevert
		// 	- Return false will abort the cctx
		indexBytes, err := cctx.GetCCTXIndexBytes()
		if err != nil {
			return false, errors.Wrap(types.ErrUnableToParseCCTXIndexBytes, err.Error())
		}
		data, err := base64.StdEncoding.DecodeString(cctx.RelayedMessage)
		if err != nil {
			return true, errors.Wrap(types.ErrUnableToDecodeMessageString, err.Error())
		}
		// if coin type is Zeta, this is a deposit ZETA to zEVM cctx.
		evmTxResponse, err := k.fungibleKeeper.LegacyZETADepositAndCallContract(
			ctx,
			sender,
			to,
			inboundSenderChainID,
			inboundAmount,
			data,
			indexBytes,
		)

		if fungibletypes.IsContractReverted(evmTxResponse, err) || errShouldRevertCctx(err) {
			return true, err
		} else if err != nil {
			return false, err
		}
	} else {
		var (
			message []byte
			err     error
		)

		// in protocol version 1, the destination of the deposit is the first 20 bytes of the message when the message is not empty
		// in protocol version 2, the destination of the deposit is always the to address, the message is the data to be sent to the contract
		if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V1 {
			var parsedAddress ethcommon.Address
			parsedAddress, message, err = memo.DecodeLegacyMemoHex(cctx.RelayedMessage)
			if err != nil {
				return false, errors.Wrap(types.ErrUnableToParseAddress, err.Error())
			}
			if parsedAddress != (ethcommon.Address{}) {
				to = parsedAddress
			}
		} else if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V2 {
			if len(cctx.RelayedMessage) > 0 {
				message, err = hex.DecodeString(cctx.RelayedMessage)
				if err != nil {
					return false, errors.Wrap(types.ErrUnableToDecodeMessageString, err.Error())
				}
			}
		}

		from, err := chains.DecodeAddressFromChainID(
			inboundSenderChainID,
			inboundSender,
			k.GetAuthorityKeeper().GetAdditionalChainList(ctx),
		)
		if err != nil {
			return false, fmt.Errorf("HandleEVMDeposit: unable to decode address: %w", err)
		}

		// use a temporary context to not commit any state change in case of error
		// note: ZRC20DepositAndCallContract is solely responsible for calling the contract and depositing tokens if needed
		// and does not include any other side effects or any logic that modifies the state directly
		tmpCtx, commit := ctx.CacheContext()

		// contractCall is the same as the isCrossChainCall for V2 protocol version
		// when removing V1 flows this can be simplified
		evmTxResponse, contractCall, err := k.fungibleKeeper.ZRC20DepositAndCallContract(
			tmpCtx,
			from,
			to,
			inboundAmount,
			inboundSenderChainID,
			message,
			inboundCoinType,
			cctx.InboundParams.Asset,
			cctx.ProtocolContractVersion,
			cctx.InboundParams.IsCrossChainCall,
		)
		if fungibletypes.IsContractReverted(evmTxResponse, err) || errShouldRevertCctx(err) {
			// this is a contract revert, we commit the state to save the emitted logs related to revert
			commit()
			return true, err
		} else if err != nil {
			// this should not happen and we don't commit the state to avoid inconsistent state
			return false, err
		}

		// non-empty msg.Message means this is a contract call; therefore, the logs should be processed.
		// a withdrawal event in the logs could generate cctxs for outbound transactions.
		if !evmTxResponse.Failed() && contractCall {
			logs := evmtypes.LogsToEthereum(evmTxResponse.Logs)
			if len(logs) > 0 {
				tmpCtx = tmpCtx.WithValue(InCCTXIndexKey, cctx.Index)
				txOrigin := cctx.InboundParams.TxOrigin
				if txOrigin == "" {
					txOrigin = inboundSender
				}

				// process logs to process cctx events initiated during the contract call
				err = k.ProcessLogs(tmpCtx, logs, to, txOrigin)
				if err != nil {
					// this happens if the cctx events are not processed correctly with invalid withdrawals
					// in this situation we want the CCTX to be reverted, we don't commit the state so the contract call is not persisted
					// the contract call is considered as reverted
					return true, errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
				}
				tmpCtx.EventManager().EmitEvent(
					sdk.NewEvent(sdk.EventTypeMessage,
						sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
						sdk.NewAttribute("action", "CallDepositAndCall"),
						sdk.NewAttribute("contract", to.String()),
						sdk.NewAttribute("data", hex.EncodeToString(message)),
						sdk.NewAttribute("cctxIndex", cctx.Index),
					),
				)
			}
		}

		// commit state change from the deposit and eventual cctx events
		commit()
	}

	return false, nil
}

// errShouldRevertCctx returns true if the cctx should revert from the error of the deposit
// we revert the cctx if a non-contract is tried to be called, if the liquidity cap is reached, or if the zrc20 is paused
func errShouldRevertCctx(err error) bool {
	return errors.Is(err, fungibletypes.ErrForeignCoinCapReached) ||
		errors.Is(err, fungibletypes.ErrCallNonContract) ||
		errors.Is(err, fungibletypes.ErrPausedZRC20)
}
