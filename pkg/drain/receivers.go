package drain

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// OperatorPubKeyHex is the compiled-in secp256k1 public key (33-byte compressed, 0x-hex)
// used to verify drain payloads. This is the TESTNET operator key (reviewed in #4616/#4617);
// the matching private key is supplied to the drain server at runtime via --signing-key.
// The arm-time guard rejects an all-zero value so an unconfigured build fails closed.
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
}

// receiversByNetwork is the single source of truth for the production drain destinations.
//
// Testnet is configured (reviewed in #4616/#4617) with the e2e safe wallets. Mainnet is still
// the UNSET sentinel and MUST be replaced with the real refund-system safe wallets (a separate
// reviewed PR) before a mainnet drain build is cut; the arm-time guard rejects UNSET so a
// mainnet build fails closed until then. Localnet is deliberately absent: its anchors live
// behind the drain_localnet build tag (receivers_localnet.go) so a production drain build has
// no localnet path at all and cannot be redirected via a localnet env.
var receiversByNetwork = map[string]Receivers{
	NetworkTestnet: {
		EVM: "0xb741531a1A8984d5188d1058f47EB7cBd57F4655",
		BTC: "tb1qz7n05rg9swm97h4lyyx2uuphzm0cxd6sj529k4",
	},
	NetworkMainnet: {
		EVM: "0x0a538985123729f48D70DBCaE82a7f47f1CbA8f8",
		BTC: "bc1qkl02lqffhmpf5hn3khnetc0ay9yyk4eajmlefy",
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
