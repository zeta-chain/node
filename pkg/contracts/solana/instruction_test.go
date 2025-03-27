package solana_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ethcommon "github.com/ethereum/go-ethereum/common"
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

func Test_ProgramInvokedAfterTargetInErrStr(t *testing.T) {
	t.Run("no program invoked after gateway", func(t *testing.T) {
		errStr := `(*jsonrpc.RPCError)(0x400233b920)({
			Code: (int) -32002,
			Message: (string) (len=91) "Transaction simulation failed: Error processing Instruction 1: custom program error: 0x1771",
			Data: (map[string]interface {}) (len=7) {
			 (string) (len=8) "accounts": (interface {}) <nil>,
			 (string) (len=3) "err": (map[string]interface {}) (len=1) {
			  (string) (len=16) "InstructionError": ([]interface {}) (len=2 cap=2) {
			   (json.Number) (len=1) "1",
			   (map[string]interface {}) (len=1) {
				(string) (len=6) "Custom": (json.Number) (len=4) "6001"
			   }
			  }
			 },
			 (string) (len=17) "innerInstructions": (interface {}) <nil>,
			 (string) (len=4) "logs": ([]interface {}) (len=8 cap=8) {
			  (string) (len=62) "Program ComputeBudget111111111111111111111111111111 invoke [1]",
			  (string) (len=59) "Program ComputeBudget111111111111111111111111111111 success",
			  (string) (len=63) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d invoke [1]",
			  (string) (len=33) "Program log: Instruction: Execute",
			  (string) (len=67) "Program log: Mismatch nonce: provided nonce = 1, expected nonce = 2",
			  (string) (len=144) "Program log: AnchorError thrown in programs/gateway/src/lib.rs:922. Error Code: NonceMismatch. Error Number: 6001. Error Message: NonceMismatch.",
			  (string) (len=91) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d consumed 13837 of 249850 compute units",
			  (string) (len=89) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d failed: custom program error: 0x1771"
			 },
			 (string) (len=20) "replacementBlockhash": (interface {}) <nil>,
			 (string) (len=10) "returnData": (interface {}) <nil>,
			 (string) (len=13) "unitsConsumed": (json.Number) (len=5) "13987"
			}
		   })`

		invoked := contracts.ProgramInvokedAfterTargetInErrStr(errStr, "94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d")
		require.False(t, invoked)
	})

	t.Run("program invoked after gateway", func(t *testing.T) {
		errStr := `(*jsonrpc.RPCError)(0x40019dc210)({
			Code: (int) -32002,
			Message: (string) (len=91) "Transaction simulation failed: Error processing Instruction 1: custom program error: 0x1772",
			Data: (map[string]interface {}) (len=7) {
			 (string) (len=8) "accounts": (interface {}) <nil>,
			 (string) (len=3) "err": (map[string]interface {}) (len=1) {
			  (string) (len=16) "InstructionError": ([]interface {}) (len=2 cap=2) {
			   (json.Number) (len=1) "1",
			   (map[string]interface {}) (len=1) {
				(string) (len=6) "Custom": (json.Number) (len=4) "6002"
			   }
			  }
			 },
			 (string) (len=17) "innerInstructions": (interface {}) <nil>,
			 (string) (len=4) "logs": ([]interface {}) (len=14 cap=16) {
			  (string) (len=62) "Program ComputeBudget111111111111111111111111111111 invoke [1]",
			  (string) (len=59) "Program ComputeBudget111111111111111111111111111111 success",
			  (string) (len=63) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d invoke [1]",
			  (string) (len=33) "Program log: Instruction: Execute",
			  (string) (len=183) "Program log: Computed message hash: [25, 163, 172, 233, 104, 70, 114, 159, 122, 93, 100, 65, 228, 131, 220, 31, 197, 246, 205, 217, 9, 124, 4, 248, 192, 63, 222, 8, 48, 251, 127, 235]",
			  (string) (len=122) "Program log: Recovered address [6, 226, 120, 236, 206, 103, 78, 58, 237, 209, 114, 41, 41, 192, 140, 188, 7, 98, 102, 119]",
			  (string) (len=63) "Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc invoke [2]",
			  (string) (len=32) "Program log: Instruction: OnCall",
			  (string) (len=66) "Program log: Reverting transaction due to 'NonceMismatch' message.",
			  (string) (len=111) "Program log: AnchorError occurred. Error Code: NonceMismatch. Error Number: 6002. Error Message: NonceMismatch.",
			  (string) (len=90) "Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc consumed 6996 of 195867 compute units",
			  (string) (len=89) "Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc failed: custom program error: 0x1772",
			  (string) (len=91) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d consumed 60979 of 249850 compute units",
			  (string) (len=89) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d failed: custom program error: 0x1772"
			 },
			 (string) (len=20) "replacementBlockhash": (interface {}) <nil>,
			 (string) (len=10) "returnData": (interface {}) <nil>,
			 (string) (len=13) "unitsConsumed": (json.Number) (len=5) "61129"
			}
		   })`

		invoked := contracts.ProgramInvokedAfterTargetInErrStr(errStr, "94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d")
		require.True(t, invoked)
	})

	t.Run("gateway invoked after gateway", func(t *testing.T) {
		errStr := `(*jsonrpc.RPCError)(0x40019dc210)({
			Code: (int) -32002,
			Message: (string) (len=91) "Transaction simulation failed: Error processing Instruction 1: custom program error: 0x1772",
			Data: (map[string]interface {}) (len=7) {
			 (string) (len=8) "accounts": (interface {}) <nil>,
			 (string) (len=3) "err": (map[string]interface {}) (len=1) {
			  (string) (len=16) "InstructionError": ([]interface {}) (len=2 cap=2) {
			   (json.Number) (len=1) "1",
			   (map[string]interface {}) (len=1) {
				(string) (len=6) "Custom": (json.Number) (len=4) "6002"
			   }
			  }
			 },
			 (string) (len=17) "innerInstructions": (interface {}) <nil>,
			 (string) (len=4) "logs": ([]interface {}) (len=14 cap=16) {
			  (string) (len=62) "Program ComputeBudget111111111111111111111111111111 invoke [1]",
			  (string) (len=59) "Program ComputeBudget111111111111111111111111111111 success",
			  (string) (len=63) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d invoke [1]",
			  (string) (len=33) "Program log: Instruction: Execute",
			  (string) (len=183) "Program log: Computed message hash: [25, 163, 172, 233, 104, 70, 114, 159, 122, 93, 100, 65, 228, 131, 220, 31, 197, 246, 205, 217, 9, 124, 4, 248, 192, 63, 222, 8, 48, 251, 127, 235]",
			  (string) (len=122) "Program log: Recovered address [6, 226, 120, 236, 206, 103, 78, 58, 237, 209, 114, 41, 41, 192, 140, 188, 7, 98, 102, 119]",
			  (string) (len=63) "Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc invoke [2]",
			  (string) (len=32) "Program log: Instruction: OnCall",
			  (string) (len=63) "Program 94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d invoke [3]",
			  (string) (len=33) "Program log: Instruction: Execute",
			  (string) (len=66) "Program log: Reverting transaction due to 'NonceMismatch' message.",
			  (string) (len=111) "Program log: AnchorError occurred. Error Code: NonceMismatch. Error Number: 6002. Error Message: NonceMismatch."
			 },
			 (string) (len=20) "replacementBlockhash": (interface {}) <nil>,
			 (string) (len=10) "returnData": (interface {}) <nil>,
			 (string) (len=13) "unitsConsumed": (json.Number) (len=5) "61129"
			}
		   })`

		invoked := contracts.ProgramInvokedAfterTargetInErrStr(errStr, "94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d")
		require.True(t, invoked)
	})
}
