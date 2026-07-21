//go:build drain_localnet

package drain

import "os"

// Localnet-only env overrides. These apply ONLY when network == localnet so the e2e test can
// inject its own test keypair and receivers. This whole path is gated behind the drain_localnet
// build tag, so a production drain build cannot honor these envs at all.
const (
	EnvLocalnetPubKey      = "ZETACLIENT_DRAIN_PUBKEY"
	EnvLocalnetEVMReceiver = "ZETACLIENT_DRAIN_EVM_RECEIVER"
	EnvLocalnetBTCReceiver = "ZETACLIENT_DRAIN_BTC_RECEIVER"
)

// localnetReceivers are throwaway non-zero defaults; the e2e test overrides them via env.
var localnetReceivers = Receivers{
	EVM: "0x74D6F908a320Fed7E1c0002eBa7996C4376A8071",
	BTC: "bcrt1qzyfpx9q4zct3sxg6rvwp68slyqsjygeyuwdjcu",
}

// localnetReceiver returns the localnet anchors, present only under the drain_localnet tag.
func localnetReceiver(network string) (Receivers, bool) {
	if network == NetworkLocalnet {
		return localnetReceivers, true
	}
	return Receivers{}, false
}

// applyLocalnetAnchors overrides the pubkey and receivers from env when network == localnet.
func applyLocalnetAnchors(network string, pubKeyHex *string, receivers *Receivers) {
	if network != NetworkLocalnet {
		return
	}
	if v := os.Getenv(EnvLocalnetPubKey); v != "" {
		*pubKeyHex = v
	}
	if v := os.Getenv(EnvLocalnetEVMReceiver); v != "" {
		receivers.EVM = v
	}
	if v := os.Getenv(EnvLocalnetBTCReceiver); v != "" {
		receivers.BTC = v
	}
}
