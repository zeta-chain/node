package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	contracts "github.com/zeta-chain/zetacore/contracts/evm"
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
		fmt.Printf("#############################\n")
		fmt.Printf("withdrawal to %s amount %d\n", hex.EncodeToString(event.To), event.Value)
		fmt.Printf("#############################\n")
		fmt.Printf("#############################\n")
		//zetacoreTypes.NewMsgSendVoter("", receipt.ContractAddress.Hex(), common.ZEVMChain, event.)
		//k.zetacoreKeeper.GetSend()

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
