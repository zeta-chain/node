package keeper

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k msgServer) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx, msg types.MsgVoteOnObservedInboundTx, senderChain *common.Chain) error {
	to := ethcommon.HexToAddress(msg.Receiver)
	amount, ok := big.NewInt(0).SetString(msg.ZetaBurnt, 10)
	if !ok {
		return errors.Wrap(types.ErrFloatParseError, fmt.Sprintf("cannot parse zetaBurnt: %s", msg.ZetaBurnt))
	}
	if msg.CoinType == common.CoinType_Zeta {
		err := k.fungibleKeeper.DepositCoinZeta(ctx, to, amount)
		if err != nil {
			return err
		}
		cctx.OutboundTxParams.OutboundTxHash = "Mined directly to ZetaEVM without TX"
	} else {
		contract, data, err := parseContractAndData(msg.Message, msg.Asset)
		if err != nil {
			return errors.Wrap(types.ErrUnableToParseContract, err.Error())
		}
		tx, withdrawMessage, err := k.fungibleKeeper.DepositCoin(ctx, to, amount, senderChain.ChainName.String(), msg.Message, contract, data, msg.CoinType, msg.Asset)
		if err != nil {
			return err
		}
		// TODO : Return error if TX failed ?
		if !tx.Failed() && withdrawMessage {
			logs := evmtypes.LogsToEthereum(tx.Logs)
			// TODO: is passing by ctx KV a good choice?
			ctx = ctx.WithValue("inCctxIndex", cctx.Index)
			txOrigin := msg.TxOrigin
			if txOrigin == "" {
				txOrigin = msg.Sender
			}
			err = k.ProcessWithdrawalLogs(ctx, logs, contract, txOrigin)
			if err != nil {
				return err
			}

			if err != nil {
				return errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
			}
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
					sdk.NewAttribute("action", "depositZRC4AndCallContract"),
					sdk.NewAttribute("contract", contract.String()),
					sdk.NewAttribute("data", hex.EncodeToString(data)),
					sdk.NewAttribute("cctxIndex", cctx.Index),
				),
			)

			if tx != nil {
				cctx.OutboundTxParams.OutboundTxHash = tx.Hash
			}
		}
	}
	return nil
}

// message is hex encoded byte array
// [ contractAddress calldata ]
// [ 20B, variable]
func ParseContractAndData(message string) (contractAddress ethcommon.Address, data []byte, err error) {
	var AddressNull ethcommon.Address
	if len(message) == 0 {
		return AddressNull, nil, nil
	}
	data, err = hex.DecodeString(message)
	if err != nil {
		return AddressNull, nil, err
	}
	if len(data) < 20 {
		err = fmt.Errorf("invalid message length")
		return AddressNull, nil, err
	}
	contractAddress = ethcommon.BytesToAddress(data[:20])
	data = data[20:]
	return contractAddress, data, nil
}
