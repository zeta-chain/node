package metrics

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

func TestTelemetryServerCopiesPingRTT(t *testing.T) {
	telemetry := NewTelemetryServer()
	peerID := peer.ID("peer-a")

	rtt := map[peer.ID]int64{peerID: 10}
	telemetry.SetPingRTT(rtt)

	rtt[peerID] = 20
	require.Equal(t, int64(10), telemetry.GetPingRTT()[peerID])

	snapshot := telemetry.GetPingRTT()
	snapshot[peerID] = 30
	require.Equal(t, int64(10), telemetry.GetPingRTT()[peerID])
}

func TestTelemetryServerCopiesConnectedPeers(t *testing.T) {
	telemetry := NewTelemetryServer()
	peerA := peer.ID("peer-a")
	peerB := peer.ID("peer-b")
	peerC := peer.ID("peer-c")

	peers := []peer.AddrInfo{{ID: peerA}}
	telemetry.SetConnectedPeers(peers)

	peers[0].ID = peerB
	require.Equal(t, peerA, telemetry.GetConnectedPeers()[0].ID)

	snapshot := telemetry.GetConnectedPeers()
	snapshot[0].ID = peerC
	require.Equal(t, peerA, telemetry.GetConnectedPeers()[0].ID)
}
