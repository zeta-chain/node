package keeper

import (
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// ProcessV2Deposit handles a deposit from an inbound tx with protocol version 2
// returns [txResponse, isContractCall, error]
// isContractCall is true if the message is non empty
func (k Keeper) ProcessV2Deposit(
	ctx sdk.Context,
	zrc20Addr ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	if len(message) == 0 {
		// simple deposit
		res, err := k.DepositZRC20(ctx, zrc20Addr, to, amount)
		return res, false, err
	}
	return nil, true, errors.New("not implemented")
}
