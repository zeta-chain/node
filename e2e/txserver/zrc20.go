package txserver

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/e2e/utils"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// UpdateZRC20GasLimit updates the gas limit for given ZRC20 token
func (zts ZetaTxServer) UpdateZRC20GasLimit(
	zrc20Addr ethcommon.Address,
	newGasLimit math.Uint,
) (*sdktypes.TxResponse, error) {
	// create msg
	msg := fungibletypes.NewMsgUpdateZRC20WithdrawFee(
		zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		zrc20Addr.Hex(),
		math.ZeroUint(), // 0 flat fee
		newGasLimit,
	)

	// broadcast tx
	res, err := zts.BroadcastTx(utils.OperationalPolicyName, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to update ZRC20 gas limit for %s", zrc20Addr)
	}

	return res, nil
}
