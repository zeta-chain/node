package bank

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/zeta-chain/node/precompiles/logs"
)

const (
	DepositEventName  = "Deposit"
	WithdrawEventName = "Withdraw"
)

func (c *Contract) AddDepositLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	depositor common.Address,
	token common.Address,
	amount *big.Int,
) error {
	event := c.Abi().Events[DepositEventName]

	// depositor and ZRC20 address are indexed.
	topics, err := logs.MakeTopics(
		event,
		[]interface{}{common.BytesToAddress(depositor.Bytes())},
		[]interface{}{common.BytesToAddress(token.Bytes())},
	)
	if err != nil {
		return err
	}

	// amount is part of event data.
	data, err := logs.PackBigInt(amount)
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}
