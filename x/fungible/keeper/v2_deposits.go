package keeper

import (
	"github.com/zeta-chain/protocol-contracts/v2/pkg/systemcontract.sol"
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
	from []byte,
	senderChainID int64,
	zrc20Addr ethcommon.Address,
	to ethcommon.Address,
	amount *big.Int,
	message []byte,
) (*evmtypes.MsgEthereumTxResponse, bool, error) {
	// simple deposit
	if len(message) == 0 {
		// simple deposit
		res, err := k.DepositZRC20(ctx, zrc20Addr, to, amount)
		return res, false, err
	}

	// deposit and call
	context := systemcontract.ZContext{
		Origin:  from,
		Sender:  ethcommon.Address{},
		ChainID: big.NewInt(senderChainID),
	}
	res, err := k.CallDepositAndCallZRC20(ctx, context, zrc20Addr, amount, to, message)
	return res, true, err
}
