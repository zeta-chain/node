package common

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// EVMSend is the gas limit required to transfer tokens on an EVM based chain
	EVMSend = 21000
	// TODO: Move gas limits from zeta-client to this file
	// https://github.com/zeta-chain/node/issues/1606
)

// MultiplyGasPrice multiplies the median gas price by the given multiplier and returns the truncated value
func MultiplyGasPrice(medianGasPrice sdkmath.Uint, multiplierString string) (sdkmath.Uint, error) {
	multiplier, err := sdk.NewDecFromStr(multiplierString)
	if err != nil {
		return sdkmath.ZeroUint(), err
	}
	gasPrice, err := sdk.NewDecFromStr(medianGasPrice.String())
	if err != nil {
		return sdkmath.ZeroUint(), err
	}
	return sdkmath.NewUintFromString(gasPrice.Mul(multiplier).TruncateInt().String()), nil
}
