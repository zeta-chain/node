package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_BaseContract(t *testing.T) {
	addr := common.BytesToAddress([]byte{0x1})

	contract := NewBaseContract(addr)
	require.NotNil(t, contract)

	// By implementing RegistryKey() the baseContract struct implements Registrable.
	// Registrable is the unique requisite to implement BaseContract as well.
	require.Equal(t, addr, contract.RegistryKey())
}

func Test_BytesToBigInt(t *testing.T) {
	data := []byte{0x1}
	intData := BytesToBigInt(data)
	require.NotNil(t, intData)
	require.Equal(t, int64(1), intData.Int64())
}
