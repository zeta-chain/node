//go:build !drain_localnet

package drain

// localnetReceiver is a no-op in production builds: without the drain_localnet tag there is no
// localnet anchor path, so network == localnet fails closed in ReceiverForNetwork.
func localnetReceiver(string) (Receivers, bool) { return Receivers{}, false }

// applyLocalnetAnchors is a no-op in production builds: env overrides of the anchors are
// impossible, so ZETACLIENT_DRAIN_NETWORK=localnet cannot redirect the drain.
func applyLocalnetAnchors(string, *string, *Receivers) {}
