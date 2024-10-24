package bank

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/zeta-chain/node/precompiles/logs"
)

type eventData struct {
	zrc20Addr  common.Address
	zrc20Token common.Address
	cosmosAddr string
	cosmosCoin string
	amount     *big.Int
}

func (c *Contract) addEventLog(
	ctx sdk.Context,
	stateDB vm.StateDB,
	eventName string,
	eventData eventData,
) error {
	event := c.Abi().Events[eventName]

	topics, err := logs.MakeTopics(
		event,
		[]interface{}{eventData.zrc20Addr},
		[]interface{}{eventData.zrc20Token},
		[]interface{}{eventData.cosmosCoin},
	)
	if err != nil {
		return err
	}

	data, err := logs.PackArguments([]logs.Argument{
		{Type: "string", Value: eventData.cosmosAddr},
		{Type: "uint256", Value: eventData.amount},
	})
	if err != nil {
		return err
	}

	logs.AddLog(ctx, c.Address(), stateDB, topics, data)

	return nil
}
