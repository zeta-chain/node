package e2etests

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/contracts/opcode"
	"github.com/zeta-chain/node/e2e/runner"
)

func TestOpcodes(r *runner.E2ERunner, _ []string) {
	// deploy the opcode contract and run function using opcode
	addr, tx, opcodeCaller, err := opcode.DeployOpcode(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(tx)
	r.Logger.Print("Deployed Opcode contract at %s", addr.Hex())

	// check push 0
	r.Logger.Print("Testing PUSH0 opcode...")
	tx, err = opcodeCaller.TestPUSH0(r.ZEVMAuth)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(tx)
	r.Logger.Print("PUSH0 opcode verified")
}
