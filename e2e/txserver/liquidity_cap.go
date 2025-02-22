package txserver

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/e2e/utils"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// SetZRC20LiquidityCap sets the liquidity cap for given ZRC20 token
func (zts ZetaTxServer) SetZRC20LiquidityCap(zrc20Addr string, liquidityCap math.Uint) (*sdktypes.TxResponse, error) {
	// create msg
	msg := fungibletypes.NewMsgUpdateZRC20LiquidityCap(
		zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		zrc20Addr,
		liquidityCap,
	)

	// broadcast tx
	res, err := zts.BroadcastTx(utils.OperationalPolicyName, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to set ZRC20 liquidity cap for %s", zrc20Addr)
	}

	return res, nil
}

// RemoveZRC20LiquidityCap removes the liquidity cap for given ZRC20 token
func (zts ZetaTxServer) RemoveZRC20LiquidityCap(zrc20Addr string) (*sdktypes.TxResponse, error) {
	// create msg with zero liquidity cap
	msg := fungibletypes.NewMsgUpdateZRC20LiquidityCap(
		zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		zrc20Addr,
		math.ZeroUint(),
	)

	// broadcast tx
	res, err := zts.BroadcastTx(utils.OperationalPolicyName, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to remove ZRC20 liquidity cap for %s", zrc20Addr)
	}

	return res, nil
}
