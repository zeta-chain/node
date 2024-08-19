package e2etests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/precompiles/prototype"
)

func TestPrecompilesRegular(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	iPrototype, err := prototype.NewIPrototype(prototype.ContractAddress, r.ZEVMClient)
	require.NoError(r, err, "Failed to create prototype contract caller")

	res, err := iPrototype.Bech32ify(nil, "zeta", common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"))
	require.NoError(r, err, "Error calling Bech32ify")
	require.Equal(r, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u", res, "Failed to validate Bech32ify result")

	addr, err := iPrototype.Bech32ToHexAddr(nil, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u")
	require.NoError(r, err, "Error calling Bech32ToHexAddr")
	require.Equal(
		r,
		"0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE",
		addr.String(),
		"Failed to validate Bech32ToHexAddr result",
	)

	balance, err := iPrototype.GetGasStabilityPoolBalance(nil, int64(1337))
	require.NoError(r, err, "Error calling GetGasStabilityPoolBalance")
	require.NotNil(r, balance, "GetGasStabilityPoolBalance returned balance is nil")
}
