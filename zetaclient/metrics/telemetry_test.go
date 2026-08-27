package metrics

import (
	"testing"

	"github.com/libp2p/go-libp2p/core/peer"
	maddr "github.com/multiformats/go-multiaddr"
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
	addrA := maddr.StringCast("/ip4/127.0.0.1/tcp/1111")
	addrB := maddr.StringCast("/ip4/127.0.0.1/tcp/2222")
	addrC := maddr.StringCast("/ip4/127.0.0.1/tcp/3333")

	peers := []peer.AddrInfo{{ID: peerA, Addrs: []maddr.Multiaddr{addrA}}}
	telemetry.SetConnectedPeers(peers)

	peers[0].ID = peerB
	peers[0].Addrs[0] = addrB
	require.Equal(t, peerA, telemetry.GetConnectedPeers()[0].ID)
	require.Equal(t, addrA, telemetry.GetConnectedPeers()[0].Addrs[0])

	snapshot := telemetry.GetConnectedPeers()
	snapshot[0].ID = peerC
	snapshot[0].Addrs[0] = addrC
	require.Equal(t, peerA, telemetry.GetConnectedPeers()[0].ID)
	require.Equal(t, addrA, telemetry.GetConnectedPeers()[0].Addrs[0])
}
