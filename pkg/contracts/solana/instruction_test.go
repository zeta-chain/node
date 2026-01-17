package solana_test

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	contracts "github.com/zeta-chain/node/pkg/contracts/solana"
)

const (
	// testSigner is the address of the signer for unit tests
	testSigner = "0xaD32427bA235a8350b7805C1b85147c8ea03F437"
)

// getTestSignature returns the signature produced by 'testSigner' for the withdraw instruction:
// ChainID: 902
// Nonce: 0
// Amount: 1336000
// To: 37yGiHAnLvWZUNVwu9esp74YQFqxU1qHCbABkDvRddUQ
func getTestSignature() [65]byte {
	return [65]byte{
		57, 160, 150, 241, 113, 78, 5, 205, 104, 97, 176, 136, 113, 84, 183, 119,
		213, 119, 29, 1, 183, 3, 43, 27, 140, 39, 33, 185, 6, 122, 69, 140,
		42, 102, 187, 143, 110, 9, 106, 162, 158, 26, 135, 253, 130, 157, 216, 191,
		117, 23, 179, 243, 109, 175, 101, 19, 95, 192, 16, 240, 40, 99, 105, 216, 0,
	}
}

// getTestmessageHash returns the message hash used to produce 'testSignature'
func getTestmessageHash() [32]byte {
	return [32]byte{
		162, 12, 221, 179, 248, 136, 244, 6, 76, 237, 137, 42, 71, 113, 1, 244,
		84, 105, 168, 197, 15, 120, 59, 150, 109, 63, 236, 36, 85, 136, 124, 5,
	}
}

func Test_SignerWithdraw(t *testing.T) {
	var sigRS [64]byte
	sigTest := getTestSignature()
	copy(sigRS[:], sigTest[:64])

	// create a withdraw instruction
	inst := contracts.WithdrawInstructionParams{
		Signature:   sigRS,
		RecoveryID:  0,
		MessageHash: getTestmessageHash(),
	}

	// recover signer
	signer, err := inst.Signer()
	require.NoError(t, err)
	require.EqualValues(t, testSigner, signer.String())
}

func Test_RecoverSigner(t *testing.T) {
	sigTest := getTestSignature()
	hashTest := getTestmessageHash()

	// recover the signer from the test message hash and signature
	signer, err := contracts.RecoverSigner(hashTest[:], sigTest[:])
	require.NoError(t, err)
	require.EqualValues(t, testSigner, signer.String())

	// slightly modify the signature and recover the signer
	sigFake := sigTest
	sigFake[0]++
	signer, err = contracts.RecoverSigner(hashTest[:], sigFake[:])
	require.Error(t, err)
	require.Equal(t, ethcommon.Address{}, signer)

	// slightly modify the message hash and recover the signer
	hashFake := hashTest
	hashFake[0]++
	signer, err = contracts.RecoverSigner(hashFake[:], sigTest[:])
	require.NoError(t, err)
	require.NotEqual(t, ethcommon.Address{}, signer)
	require.NotEqual(t, testSigner, signer.String())
}
