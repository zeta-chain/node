package txserver

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/e2e/utils"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

// SetZRC20LiquidityCap sets the liquidity cap for given ZRC20 token
func (zts ZetaTxServer) SetZRC20LiquidityCap(
	zrc20Addr ethcommon.Address,
	liquidityCap math.Uint,
) (*sdktypes.TxResponse, error) {
	return zts.updateZRC20LiquidityCap(zrc20Addr, liquidityCap)
}

// RemoveZRC20LiquidityCap removes the liquidity cap for given ZRC20 token
func (zts ZetaTxServer) RemoveZRC20LiquidityCap(zrc20Addr ethcommon.Address) (*sdktypes.TxResponse, error) {
	return zts.updateZRC20LiquidityCap(zrc20Addr, math.ZeroUint())
}

// updateZRC20LiquidityCap updates the liquidity cap for given ZRC20 token
func (zts ZetaTxServer) updateZRC20LiquidityCap(
	zrc20Addr ethcommon.Address,
	liquidityCap math.Uint,
) (*sdktypes.TxResponse, error) {
	// create msg
	msg := fungibletypes.NewMsgUpdateZRC20LiquidityCap(
		zts.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		zrc20Addr.Hex(),
		liquidityCap,
	)

	// broadcast tx
	res, err := zts.BroadcastTx(utils.OperationalPolicyName, msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to set ZRC20 liquidity cap for %s", zrc20Addr)
	}

	return res, nil
}
