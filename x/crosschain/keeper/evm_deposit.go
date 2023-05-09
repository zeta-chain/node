package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k msgServer) HandleEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx, msg types.MsgVoteOnObservedInboundTx, senderChain *common.Chain) error {
	to := ethcommon.HexToAddress(msg.Receiver)
	//amount, ok := big.NewInt(0).SetString(msg.ZetaBurnt, 10)

	if msg.CoinType == common.CoinType_Zeta {
		err := k.fungibleKeeper.DepositCoinZeta(ctx, to, msg.Amount.BigInt())
		if err != nil {
			return err
		}
		cctx.GetCurrentOutTxParam().OutboundTxHash = "Mined directly to ZetaEVM without TX"
	} else { // cointype is Gas or ERC20
		contract, data, err := parseContractAndData(msg.Message, msg.Asset)
		if err != nil {
			return errors.Wrap(types.ErrUnableToParseContract, err.Error())
		}
		tx, _, err := k.fungibleKeeper.DepositCoin(ctx, to, msg.Amount.BigInt(), senderChain, msg.Message, contract, data, msg.CoinType, msg.Asset)
		if err != nil {
			return err
		}
		// TODO : Return error if TX failed ?

		if !tx.Failed() && len(msg.Message) > 0 {
			logs := evmtypes.LogsToEthereum(tx.Logs)
			// TODO: is passing by ctx KV a good choice?
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
					sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
					sdk.NewAttribute("action", "depositZRC4AndCallContract"),
					sdk.NewAttribute("contract", contract.String()),
					sdk.NewAttribute("data", hex.EncodeToString(data)),
					sdk.NewAttribute("cctxIndex", cctx.Index),
				),
			)

			if tx != nil {
				cctx.GetCurrentOutTxParam().OutboundTxHash = tx.Hash
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
