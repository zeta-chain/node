package utils

import (
	"math/big"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
)

func ParseFloat(t require.TestingT, s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	require.NoError(t, err, "unable to parse float %q", s)
	return f
}

func ParseInt(t require.TestingT, s string) int {
	v, err := strconv.Atoi(s)
	require.NoError(t, err, "unable to parse int from %q", s)

	return v
}

func ParseBigInt(t require.TestingT, s string) *big.Int {
	v, ok := big.NewInt(0).SetString(s, 10)
	require.True(t, ok, "unable to parse big.Int from %q", s)

	return v
}

func ParseUint8Array(t require.TestingT, s string) []uint8 {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	indexes := make([]uint8, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		u, err := strconv.ParseUint(p, 10, 8)
		require.NoError(t, err, "invalid uint8: %q", p)

		indexes = append(indexes, uint8(u))
	}

	return indexes
}

func ParseUint(t require.TestingT, s string) math.Uint {
	return math.NewUintFromBigInt(ParseBigInt(t, s))
}

// BTCAmountFromFloat64 takes float64 (e.g. 0.001) that represents btc amount
// and converts it to big.Int for downstream usage.
func BTCAmountFromFloat64(t require.TestingT, amount float64) *big.Int {
	satoshi, err := btcutil.NewAmount(amount)
	require.NoError(t, err)

	return big.NewInt(int64(satoshi))
}

func ParseBitcoinWithdrawArgs(
	t require.TestingT,
	args []string,
	defaultReceiver string,
	bitcoinChainID int64,
) (btcutil.Address, *big.Int) {
	require.NotEmpty(t, args, "args list is empty")

	receiverRaw := defaultReceiver
	if args[0] != "" {
		receiverRaw = args[0]
	}

	receiver, err := chains.DecodeBtcAddress(receiverRaw, bitcoinChainID)
	require.NoError(t, err, "unable to decode btc address")

	withdrawalAmount := ParseFloat(t, args[1])
	amount := BTCAmountFromFloat64(t, withdrawalAmount)

	return receiver, amount
}
