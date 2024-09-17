package bank

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	zrc20Depositor common.Address,
	zrc20Token common.Address,
	cosmosAddr string,
	cosmosCoin string,
	amount *big.Int,
) error {
	event := c.Abi().Events[DepositEventName]

	// ZRC20, cosmos coin and depositor.
	topics, err := logs.MakeTopics(
		event,
		[]interface{}{common.BytesToAddress(zrc20Depositor.Bytes())},
		[]interface{}{common.BytesToAddress(zrc20Token.Bytes())},
		[]interface{}{cosmosCoin},
	)
	if err != nil {
		return err
	}

	// Amount and cosmos address are part of event data.
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return err
	}

	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return err
	}

	arguments := abi.Arguments{
		{
			Type: stringType,
		},
		{
			Type: uint256Type,
		},
	}

	data, err := arguments.Pack(cosmosAddr, amount)
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}

func (c *Contract) AddWithdrawLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	zrc20Withdrawer common.Address,
	zrc20Token common.Address,
	cosmosAddr string,
	cosmosCoin string,
	amount *big.Int,
) error {
	event := c.Abi().Events[WithdrawEventName]

	// ZRC20, cosmos coin  and witgdrawer are indexed.
	topics, err := logs.MakeTopics(
		event,
		[]interface{}{common.BytesToAddress(zrc20Withdrawer.Bytes())},
		[]interface{}{common.BytesToAddress(zrc20Token.Bytes())},
		[]interface{}{cosmosCoin},
	)
	if err != nil {
		return err
	}

	// Amount and cosmos address are part of event data.
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return err
	}

	stringType, err := abi.NewType("string", "", nil)
	if err != nil {
		return err
	}

	arguments := abi.Arguments{
		{
			Type: stringType,
		},
		{
			Type: uint256Type,
		},
	}

	data, err := arguments.Pack(cosmosAddr, amount)
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}
