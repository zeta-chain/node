package e2etests

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
)

// randomPayload generates a random payload to be used in gateway calls for testing purposes
func randomPayload(r *runner.E2ERunner) string {
	bytes := make([]byte, 50)
	_, err := rand.Read(bytes)
	require.NoError(r, err)

	return hex.EncodeToString(bytes)
}

// withdrawBTCZRC20 is a helper function to withdraw BTC using the gateway contract
func withdrawBTCZRC20(r *runner.E2ERunner, to btcutil.Address, amount *big.Int) *ethtypes.Transaction {
	// approve gateway contract to spend BTC
	r.ApproveBTCZRC20(r.GatewayZEVMAddr)

	// withdraw BTC
	return r.BTCWithdraw(to, amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})
}

// bigAdd is shorthand for new(big.Int).Add(x, y)
func bigAdd(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

// bigSub is shorthand for new(big.Int).Sub(x, y)
func bigSub(x *big.Int, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}
