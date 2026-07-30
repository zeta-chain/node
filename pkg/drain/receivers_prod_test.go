//go:build !drain_localnet

package drain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/drain"
)

func TestResolveAnchorsLocalnetFailsClosed(t *testing.T) {
	// Without the drain_localnet build tag there is no localnet anchor path, so a production
	// drain build cannot be redirected via ZETACLIENT_DRAIN_NETWORK=localnet.
	_, _, err := drain.ResolveAnchors(drain.NetworkLocalnet)
	require.Error(t, err)
}
