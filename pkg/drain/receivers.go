package drain

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// OperatorPubKeyHex is the compiled-in secp256k1 public key (33-byte compressed, 0x-hex)
// used to verify drain payloads. It is a PLACEHOLDER and MUST be replaced with the real
// operator public key (reviewed in the PR) before a drain build is cut.
const OperatorPubKeyHex = "0x000000000000000000000000000000000000000000000000000000000000000000"

// OperatorPubKey returns the compiled-in operator public key bytes.
func OperatorPubKey() ([]byte, error) {
	return hexutil.Decode(OperatorPubKeyHex)
}

// Network selects which hardcoded safe-wallet receiver set to use.
const (
	NetworkLocalnet = "localnet"
	NetworkTestnet  = "testnet"
	NetworkMainnet  = "mainnet"
)

// Receivers holds the hardcoded safe-wallet addresses for one network. These are the
// security anchor: the zetaclient compiles them in and asserts every payload tx sends
// only here, so a compromised generator can change when funds move but never where.
type Receivers struct {
	EVM string
	BTC string
}

// receiversByNetwork is the single source of truth for drain destinations.
//
// The testnet and mainnet values are PLACEHOLDERS and MUST be replaced with the real
// safe-wallet addresses (reviewed in the PR) before a drain build is cut. Localnet
// values are only used as documentation — the e2e test injects fresh receivers directly.
var receiversByNetwork = map[string]Receivers{
	NetworkLocalnet: {
		EVM: "0x0000000000000000000000000000000000000000",
		BTC: "bcrt1qqypqxpq9qcrsszg2pvxq6rs0zqg3yyc5phstwt",
	},
	NetworkTestnet: {
		EVM: "0x0000000000000000000000000000000000000000",
		BTC: "tb1qqypqxpq9qcrsszg2pvxq6rs0zqg3yyc5r7fxez",
	},
	NetworkMainnet: {
		EVM: "0x0000000000000000000000000000000000000000",
		BTC: "bc1qqypqxpq9qcrsszg2pvxq6rs0zqg3yyc5fcj4z3",
	},
}

// ReceiverForNetwork returns the hardcoded receivers for the given network.
func ReceiverForNetwork(network string) (Receivers, error) {
	r, ok := receiversByNetwork[network]
	if !ok {
		return Receivers{}, errors.Errorf("no drain receivers configured for network %q", network)
	}
	return r, nil
}
