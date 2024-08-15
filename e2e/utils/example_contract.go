package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/zetacore/testutil/contracts"
)

// MustHaveCalledExampleContract checks if the contract has been called correctly
func MustHaveCalledExampleContract(
	t require.TestingT,
	contract *testcontract.Example,
	amount *big.Int,
) {
	bar, err := contract.Bar(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(
		t,
		0,
		bar.Cmp(amount),
		"cross-chain call failed bar value %s should be equal to amount %s",
		bar.String(),
		amount.String(),
	)
}
