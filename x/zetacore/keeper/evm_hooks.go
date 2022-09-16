package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/evm"
	zetacoretypes "github.com/zeta-chain/zetacore/x/zetacore/types"
)

var _ evmtypes.EvmHooks = Hooks{}

type Hooks struct {
	k Keeper
}

func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on
// the module keeper
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

// PostTxProcessing implements EvmHooks.PostTxProcessing.
func (k Keeper) PostTxProcessing(
	ctx sdk.Context,
	msg core.Message,
	receipt *ethtypes.Receipt,
) error {

	var event *contracts.ZRC4Withdrawal
	found := false
	for _, log := range receipt.Logs {
		e, err := ParseWithdrawalEvent(*log)
		if err != nil {
			fmt.Printf("######### skip log %s #########\n", log.Topics[0].String())
			continue
		} else {
			found = true
			event = e
		}
	}
	if found {
		fmt.Printf("#############################\n")
		fmt.Printf("withdrawal to %s amount %d\n", hex.EncodeToString(event.To), event.Value)
		fmt.Printf("#############################\n")

		msg := zetacoretypes.NewMsgSendVoter("", receipt.ContractAddress.Hex(), common.ZETAChain.String(), string(event.To), "GOERLI", event.Value.String(), "", "", event.Raw.TxHash.String(), event.Raw.BlockNumber, 90000, common.CoinType_Gas)
		sendHash := msg.Digest()

		//k.zetacoreKeeper.Get

		send, found := k.GetSend(ctx, sendHash)
		if found { // this is not supposed to happen
			fmt.Printf("send already exists %s\n", sendHash)
			return nil
		}
		fmt.Printf("send %s not found, create it\n", sendHash)

		send.Sender = receipt.ContractAddress.Hex()
		send.SenderChain = "ZETA"
		send.Receiver = "0x" + hex.EncodeToString(event.To)
		send.ReceiverChain = "GOERLI"
		send.FinalizedMetaHeight = uint64(ctx.BlockHeight())
		send.Index = sendHash
		send.ZetaBurnt = event.Value.String()
		send.ZetaMint = event.Value.String() // does not deduct gas fee
		send.InTxHash = event.Raw.TxHash.String()
		send.InBlockHeight = event.Raw.BlockNumber
		send.GasLimit = 90_000
		gasprice, found := k.GetGasPrice(ctx, "GOERLI")
		if !found {
			fmt.Printf("chain nonce not found for GOERLI\n")
			return nil
		}
		send.GasPrice = fmt.Sprintf("%d", gasprice.Prices[gasprice.MedianIndex])
		send.Status = zetacoretypes.SendStatus_PendingOutbound
		chainNonce, found := k.GetChainNonces(ctx, "GOERLI")
		if !found {
			fmt.Printf("chain nonce not found for GOERLI\n")
			return nil
		}
		send.Nonce = chainNonce.Nonce
		chainNonce.Nonce++
		k.SetChainNonces(ctx, chainNonce)

		k.SetSend(ctx, send)
		fmt.Printf("####setting send... ###########\n")
		//fmt.Printf("%v\n", send)
	}

	return nil
}

func ParseWithdrawalEvent(log ethtypes.Log) (*contracts.ZRC4Withdrawal, error) {
	zrc4Abi, err := contracts.ZRC4MetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	event := new(contracts.ZRC4Withdrawal)
	eventName := "Withdrawal"
	if log.Topics[0] != zrc4Abi.Events[eventName].ID {
		return nil, fmt.Errorf("event signature mismatch")
	}
	if len(log.Data) > 0 {
		if err := zrc4Abi.UnpackIntoInterface(event, eventName, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range zrc4Abi.Events[eventName].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	abi.ParseTopics(event, indexed, log.Topics[1:])
	event.Raw = log

	return event, nil
}
