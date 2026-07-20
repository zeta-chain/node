package drain

import (
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// OperatorPubKeyHex is the compiled-in secp256k1 public key (33-byte compressed, 0x-hex)
// used to verify drain payloads. It is a PLACEHOLDER and MUST be replaced with the real
// operator public key (reviewed in the PR) before a drain build is cut.
const OperatorPubKeyHex = "0x000000000000000000000000000000000000000000000000000000000000000000"

// Localnet-only env overrides. These apply ONLY when network == localnet so the e2e test
// can inject its own test keypair and receivers. Testnet/mainnet always use the compiled
// anchors — the production security model is unchanged.
const (
	EnvLocalnetPubKey      = "ZETACLIENT_DRAIN_PUBKEY"
	EnvLocalnetEVMReceiver = "ZETACLIENT_DRAIN_EVM_RECEIVER"
	EnvLocalnetBTCReceiver = "ZETACLIENT_DRAIN_BTC_RECEIVER"
)

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
// safe-wallet addresses (reviewed in the PR) before a drain build is cut. Localnet values
// are throwaway defaults the e2e test balance-checks; they can be overridden via env.
var receiversByNetwork = map[string]Receivers{
	NetworkLocalnet: {
		EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071",
		BTC: "bcrt1qzyfpx9q4zct3sxg6rvwp68slyqsjygeyuwdjcu",
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

// ResolveAnchors returns the operator public key and receivers for the given network.
// For localnet only, the pubkey and receivers may be overridden by env vars so the e2e
// test can inject its own material; testnet/mainnet always use the compiled anchors.
func ResolveAnchors(network string) (pubKey []byte, receivers Receivers, err error) {
	receivers, err = ReceiverForNetwork(network)
	if err != nil {
		return nil, Receivers{}, err
	}

	pubKeyHex := OperatorPubKeyHex
	if network == NetworkLocalnet {
		if v := os.Getenv(EnvLocalnetPubKey); v != "" {
			pubKeyHex = v
		}
		if v := os.Getenv(EnvLocalnetEVMReceiver); v != "" {
			receivers.EVM = v
		}
		if v := os.Getenv(EnvLocalnetBTCReceiver); v != "" {
			receivers.BTC = v
		}
	}

	pubKey, err = hexutil.Decode(pubKeyHex)
	if err != nil {
		return nil, Receivers{}, errors.Wrap(err, "invalid operator pubkey")
	}
	return pubKey, receivers, nil
}
