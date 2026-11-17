package drain

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// OperatorPubKeyHex is the compiled-in secp256k1 public key (33-byte compressed, 0x-hex)
// used to verify drain payloads. It is a PLACEHOLDER and MUST be replaced with the real
// operator public key (reviewed in the PR) before a drain build is cut. The arm-time guard
// rejects this all-zero placeholder so an unconfigured build fails closed.
const OperatorPubKeyHex = "0x03579d09c8a72ebf96e943c121926f3bfaf7600b9685eda7692786bf3cfca2c9fc"

// unset is the sentinel for a receiver that has not been configured. It is deliberately not
// a valid address, so an unconfigured testnet/mainnet build fails closed rather than
// silently matching a real burn address like 0x0.
const unset = "UNSET"

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
	// EVMByChain overrides EVM for specific chain ids.
	EVMByChain map[int64]string
}

// EVMFor returns the EVM receiver for chainID: its override if set, else EVM. The payload
// generator and the zetaclient assertion both resolve through here so they cannot diverge.
func (r Receivers) EVMFor(chainID int64) string {
	if a, ok := r.EVMByChain[chainID]; ok {
		return a
	}
	return r.EVM
}

// mainnetDrainedEVM is the original mainnet EVM receiver, kept for chains already drained to it.
const mainnetDrainedEVM = "0x0a538985123729f48D70DBCaE82a7f47f1CbA8f8"

// receiversByNetwork is the single source of truth for the production drain destinations.
//
// The testnet and mainnet values are UNSET sentinels and MUST be replaced with the real
// safe-wallet addresses (reviewed in the PR) before a drain build is cut; the arm-time
// guard rejects UNSET so an unconfigured build fails closed. Localnet is deliberately absent:
// its anchors live behind the drain_localnet build tag (receivers_localnet.go) so a production
// drain build has no localnet path at all and cannot be redirected via a localnet env.
var receiversByNetwork = map[string]Receivers{
	NetworkTestnet: {
		EVM: unset,
		BTC: unset,
	},
	NetworkMainnet: {
		// Gnosis Safe for every chain not yet drained (Base 8453, Arbitrum 42161, Avalanche 43114).
		// Deployed on Base only: it must exist at this address on a chain before a drain fires there.
		EVM: "0x919F207971b9Dea34150D095169388d25e167F1A",
		BTC: "bc1qkl02lqffhmpf5hn3khnetc0ay9yyk4eajmlefy",
		// Already drained to the original receiver; kept so a re-drain lands in the same place.
		EVMByChain: map[int64]string{
			1:   mainnetDrainedEVM, // Ethereum
			56:  mainnetDrainedEVM, // BSC
			137: mainnetDrainedEVM, // Polygon
		},
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
	for chainID, a := range r.EVMByChain {
		// overrides must be checksummed: a typo in one hex digit would otherwise still parse.
		if !ethcommon.IsHexAddress(a) || ethcommon.HexToAddress(a).Hex() != a ||
			ethcommon.HexToAddress(a) == (ethcommon.Address{}) {
			return errors.Errorf(
				"EVM drain receiver override for chain %d (%q) is not a valid checksummed address",
				chainID,
				a,
			)
		}
	}
	return nil
}

// ReceiverForNetwork returns the hardcoded receivers for the given network. Localnet resolves
// only when the drain_localnet build tag is set (localnetReceiver); otherwise it fails closed.
func ReceiverForNetwork(network string) (Receivers, error) {
	if r, ok := receiversByNetwork[network]; ok {
		return r, nil
	}
	if r, ok := localnetReceiver(network); ok {
		return r, nil
	}
	return Receivers{}, errors.Errorf("no drain receivers configured for network %q", network)
}

// ResolveAnchors returns the operator public key and receivers for the given network.
// Testnet/mainnet always use the compiled anchors. Localnet anchors (and their env overrides
// for the e2e test) exist only under the drain_localnet build tag via applyLocalnetAnchors; in
// a production build that hook is a no-op and localnet has already failed closed above.
func ResolveAnchors(network string) (pubKey []byte, receivers Receivers, err error) {
	receivers, err = ReceiverForNetwork(network)
	if err != nil {
		return nil, Receivers{}, err
	}

	pubKeyHex := OperatorPubKeyHex
	applyLocalnetAnchors(network, &pubKeyHex, &receivers)

	pubKey, err = hexutil.Decode(pubKeyHex)
	if err != nil {
		return nil, Receivers{}, errors.Wrap(err, "invalid operator pubkey")
	}
	return pubKey, receivers, nil
}
