package keeper

import (
	"encoding/hex"
	"fmt"

	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// HandleEVMDeposit handles a deposit from an inbound tx
// returns (isContractReverted, err)
// (true, non-nil) means CallEVM() reverted
func (k Keeper) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx, msg types.MsgVoteOnObservedInboundTx, senderChain *common.Chain) (bool, error) {
	to := ethcommon.HexToAddress(msg.Receiver)
	var ethTxHash ethcommon.Hash
	if len(ctx.TxBytes()) > 0 {
		// add event for tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		ethTxHash = ethcommon.BytesToHash(hash)
		cctx.GetCurrentOutTxParam().OutboundTxHash = ethTxHash.String()
		// #nosec G701 always positive
		cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = uint64(ctx.BlockHeight())
	}

	if msg.CoinType == common.CoinType_Zeta {
		// if coin type is Zeta, this is a deposit ZETA to zEVM cctx.
		err := k.fungibleKeeper.DepositCoinZeta(ctx, to, msg.Amount.BigInt())
		if err != nil {
			return false, err
		}
	} else {
		// cointype is Gas or ERC20; then it could be a ZRC20 deposit/depositAndCall cctx.
		parsedAddress, data, err := parseAddressAndData(msg.Message, msg.Asset)
		if err != nil {
			return false, errors.Wrap(types.ErrUnableToParseContract, err.Error())
		}
		if parsedAddress != (ethcommon.Address{}) {
			to = parsedAddress
		}

		from, err := senderChain.DecodeAddress(msg.Sender)
		if err != nil {
			return false, fmt.Errorf("HandleEVMDeposit: unable to decode address: %s", err.Error())
		}

		evmTxResponse, contractCall, err := k.fungibleKeeper.ZRC20DepositAndCallContract(
			ctx,
			from,
			to,
			msg.Amount.BigInt(),
			senderChain,
			data,
			msg.CoinType,
			msg.Asset,
		)
		if err != nil {
			isContractReverted := false

			// consider the contract as reverted if foreign coin liquidity cap is reached
			if (evmTxResponse != nil && evmTxResponse.Failed()) || errors.Is(err, fungibletypes.ErrForeignCoinCapReached) {
				isContractReverted = true
			}

			return isContractReverted, err
		}

		// non-empty msg.Message means this is a contract call; therefore the logs should be processed.
		// a withdrawal event in the logs could generate cctxs for outbound transactions.
		if !evmTxResponse.Failed() && contractCall {
			logs := evmtypes.LogsToEthereum(evmTxResponse.Logs)
			if len(logs) > 0 {
				ctx = ctx.WithValue("inCctxIndex", cctx.Index)
				txOrigin := msg.TxOrigin
				if txOrigin == "" {
					txOrigin = msg.Sender
				}

				err = k.ProcessLogs(ctx, logs, to, txOrigin)
				if err != nil {
					// ProcessLogs should not error; error indicates exception, should abort
					return false, errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
				}
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(sdk.EventTypeMessage,
						sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
						sdk.NewAttribute("action", "DepositZRC20AndCallContract"),
						sdk.NewAttribute("contract", to.String()),
						sdk.NewAttribute("data", hex.EncodeToString(data)),
						sdk.NewAttribute("cctxIndex", cctx.Index),
					),
				)
			}
		}
	}
	return false, nil
}

// parseAddressAndData parses the message string into an address and data
// message is hex encoded byte array
// [ contractAddress calldata ]
// [ 20B, variable]
func parseAddressAndData(message string, asset string) (address ethcommon.Address, data []byte, err error) {
	if len(message) == 0 {
		return ethcommon.Address{}, nil, nil
	}
	data, err = hex.DecodeString(message)
	if err != nil {
		return ethcommon.Address{}, nil, err
	}
	if len(data) < 20 {
		if len(asset) != 42 || asset[:2] != "0x" {
			return ethcommon.Address{}, nil, fmt.Errorf("invalid message length")
		}
		address = ethcommon.HexToAddress(asset)
	} else {
		address = ethcommon.BytesToAddress(data[:20])
		data = data[20:]
	}
	return address, data, nil
}
