package drain

import (
	"os"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// OperatorPubKeyHex is the compiled-in secp256k1 public key (33-byte compressed, 0x-hex)
// used to verify drain payloads. It is a PLACEHOLDER and MUST be replaced with the real
// operator public key (reviewed in the PR) before a drain build is cut. The arm-time guard
// rejects this all-zero placeholder so an unconfigured build fails closed.
const OperatorPubKeyHex = "0x000000000000000000000000000000000000000000000000000000000000000000"

// unset is the sentinel for a receiver that has not been configured. It is deliberately not
// a valid address, so an unconfigured testnet/mainnet build fails closed rather than
// silently matching a real burn address like 0x0.
const unset = "UNSET"

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
// The testnet and mainnet values are UNSET sentinels and MUST be replaced with the real
// safe-wallet addresses (reviewed in the PR) before a drain build is cut; the arm-time
// guard rejects UNSET so an unconfigured build fails closed. The localnet values are
// throwaway non-zero defaults; the e2e test overrides them with a fresh receiver via env.
var receiversByNetwork = map[string]Receivers{
	NetworkLocalnet: {
		EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071",
		BTC: "bcrt1qzyfpx9q4zct3sxg6rvwp68slyqsjygeyuwdjcu",
	},
	NetworkTestnet: {
		EVM: unset,
		BTC: unset,
	},
	NetworkMainnet: {
		EVM: unset,
		BTC: unset,
	},
}

// Validate rejects unconfigured or zero-address receivers so the poller fails closed.
func (r Receivers) Validate() error {
	switch {
	case r.EVM == "" || r.EVM == unset:
		return errors.New("EVM drain receiver is not configured")
	case !ethcommon.IsHexAddress(r.EVM):
		return errors.Errorf("EVM drain receiver %q is not a valid address", r.EVM)
	case ethcommon.HexToAddress(r.EVM) == (ethcommon.Address{}):
		return errors.New("EVM drain receiver is the zero address")
	case r.BTC == "" || r.BTC == unset:
		return errors.New("BTC drain receiver is not configured")
	}
	return nil
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
