package e2etests

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/precompiles/regular"
)

func TestPrecompilesRegular(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 0, "No arguments expected")

	dummyBech32Addr := "1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u"

	// Call the Regular contract in the static precompile address.
	contract, err := regular.NewRegular(regular.ContractAddress, r.EVMClient)
	require.NoError(r, err, "Failed to create Regular contract caller")

	addr, err := contract.Bech32ToHexAddr(
		nil,
		common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE").String(),
	)
	require.NoError(r, err, "Failed to call Bech32ToHexAddr in Regular precompiled contract")

	require.Equal(r, dummyBech32Addr, addr.String(), "Expected address %s, got %s", dummyBech32Addr, addr.String())
}
