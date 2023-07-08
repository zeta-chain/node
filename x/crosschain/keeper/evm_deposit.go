package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k msgServer) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx, msg types.MsgVoteOnObservedInboundTx, senderChain *common.Chain) error {
	to := ethcommon.HexToAddress(msg.Receiver)
	//amount, ok := big.NewInt(0).SetString(msg.ZetaBurnt, 10)

	if msg.CoinType == common.CoinType_Zeta {
		// if coin type is Zeta, this is a deposit ZETA to zEVM cctx.
		err := k.fungibleKeeper.DepositCoinZeta(ctx, to, msg.Amount.BigInt())
		if err != nil {
			return err
		}
		if len(ctx.TxBytes()) > 0 {
			// add event for tendermint transaction hash format
			hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
			ethTxHash := ethcommon.BytesToHash(hash) // NOTE(pwu): use cosmos tx hash as eth tx hash if available
			cctx.GetCurrentOutTxParam().OutboundTxHash = ethTxHash.String()
		}
	} else {
		// cointype is Gas or ERC20; then it could be a ZRC20 deposit/depositAndCall cctx.
		contract, data, err := parseContractAndData(msg.Message, msg.Asset)
		if err != nil {
			return errors.Wrap(types.ErrUnableToParseContract, err.Error())
		}
		evmTxResponse, _, err := k.fungibleKeeper.ZRC20DepositAndCallContract(ctx, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.CoinType, msg.Asset)
		if err != nil {
			return err
		}

		// non-empty msg.Message means this is a contract call; therefore the logs should be processed.
		// a withdrawal event in the logs could generate cctxs for outbound transactions.
		if !evmTxResponse.Failed() && len(msg.Message) > 0 {
			logs := evmtypes.LogsToEthereum(evmTxResponse.Logs)
			ctx = ctx.WithValue("inCctxIndex", cctx.Index)
			txOrigin := msg.TxOrigin
			if txOrigin == "" {
				txOrigin = msg.Sender
			}

			err = k.ProcessLogs(ctx, logs, contract, txOrigin)
			if err != nil {
				return err
			}

			if err != nil {
				return errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
			}
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, "crosschain"),
					sdk.NewAttribute("action", "DepositZRC20AndCallContract"),
					sdk.NewAttribute("contract", contract.String()),
					sdk.NewAttribute("data", hex.EncodeToString(data)),
					sdk.NewAttribute("cctxIndex", cctx.Index),
				),
			)

			if evmTxResponse != nil {
				cctx.GetCurrentOutTxParam().OutboundTxHash = evmTxResponse.Hash
			}
		}
	}
	return nil
}

// message is hex encoded byte array
// [ contractAddress calldata ]
// [ 20B, variable]
func parseContractAndData(message string, asset string) (contractAddress ethcommon.Address, data []byte, err error) {
	if len(message) == 0 {
		return contractAddress, nil, nil
	}
	data, err = hex.DecodeString(message)
	if err != nil {
		return contractAddress, nil, err
	}
	if len(data) < 20 {
		if len(asset) != 42 || asset[:2] != "0x" {
			err = fmt.Errorf("invalid message length")
			return contractAddress, nil, err
		}
		contractAddress = ethcommon.HexToAddress(asset)
	} else {
		contractAddress = ethcommon.BytesToAddress(data[:20])
		data = data[20:]
	}
	return contractAddress, data, nil
}
