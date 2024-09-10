package staking

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/zeta-chain/node/precompiles/logs"
)

func (c *Contract) AddStakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validator string,
	amount *big.Int,
) error {
	event := c.Abi().Events["Stake"]

	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// staker and validator are indexed event params
	topics, err := logs.MakeTopics(event, []interface{}{staker}, []interface{}{common.BytesToAddress(valAddr.Bytes())})
	if err != nil {
		return err
	}

	// amount is part of event data
	data, err := logs.PackBigInt(amount)
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}

func (c *Contract) AddUnstakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validator string,
	amount *big.Int,
) error {
	event := c.Abi().Events["Unstake"]
	valAddr, err := sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// staker and validator are indexed event params
	topics, err := logs.MakeTopics(event, []interface{}{staker}, []interface{}{common.BytesToAddress(valAddr.Bytes())})
	if err != nil {
		return err
	}

	// amount is part of event data
	data, err := logs.PackBigInt(amount)
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}

func (c *Contract) AddMoveStakeLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	staker common.Address,
	validatorSrc string,
	validatorDst string,
	amount *big.Int,
) error {
	event := c.Abi().Events["MoveStake"]
	validatorSrcAddr, err := sdk.ValAddressFromBech32(validatorSrc)
	if err != nil {
		return err
	}

	validatorDstAddr, err := sdk.ValAddressFromBech32(validatorDst)
	if err != nil {
		return err
	}

	// staker and validators are indexed event params
	topics, err := logs.MakeTopics(
		event,
		[]interface{}{staker},
		[]interface{}{common.BytesToAddress(validatorSrcAddr.Bytes())},
		[]interface{}{common.BytesToAddress(validatorDstAddr.Bytes())},
	)
	if err != nil {
		return err
	}

	// amount is part of event data
	data, err := logs.PackBigInt(amount)
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}
