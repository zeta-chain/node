package tss

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	libp2p_network "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/zeta-chain/go-tss/tss"

	"github.com/zeta-chain/node/pkg/bg"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// HealthcheckProps represents options for HealthcheckWorker.
type HealthcheckProps struct {
	Telemetry               Telemetry
	Interval                time.Duration
	WhitelistPeers          []peer.ID
	NumConnectedPeersMetric prometheus.Gauge
}

// HealthcheckWorker checks the health of the TSS server and its peers.
func HealthcheckWorker(ctx context.Context, server *tss.Server, p HealthcheckProps, logger zerolog.Logger) error {
	if p.NumConnectedPeersMetric == nil {
		return errors.New("missing NumConnectedPeersMetric")
	}

	if p.Interval == 0 {
		p.Interval = 30 * time.Second
	}

	logger = logger.With().Str(logs.FieldModule, logs.ModNameTssHealthCheck).Logger()

	// Ping & collect round trip time
	var (
		host    = server.GetP2PHost()
		pingRTT = make(map[peer.ID]int64)
		mu      = sync.Mutex{}
	)

	const pingTimeout = 5 * time.Second

	pinger := func(ctx context.Context, _ *ticker.Ticker) error {
		var wg sync.WaitGroup
		for _, peerID := range p.WhitelistPeers {
			if peerID == host.ID() {
				continue
			}

			wg.Add(1)

			go func(peerID peer.ID) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						logger.Error().
							Str("peer_id", peerID.String()).
							Interface("panic", r).
							Msg("panic during ping")
					}
				}()

				ctx, cancel := context.WithTimeout(ctx, pingTimeout)
				defer cancel()

				result := <-ping.Ping(ctx, host, peerID)
				if result.Error != nil {
					result.RTT = -1 // indicates ping error
					logger.Error().
						Err(result.Error).
						Str("peer_id", peerID.String()).
						Msg("ping error")
				}

				mu.Lock()
				pingRTT[peerID] = result.RTT.Nanoseconds()
				mu.Unlock()
			}(peerID)
		}

		wg.Wait()
		p.Telemetry.SetPingRTT(pingRTT)

		return nil
	}

	connectedPeersCounter := func(_ context.Context, _ *ticker.Ticker) error {
		p2pHost := server.GetP2PHost()
		connectedPeers := lo.Map(p2pHost.Network().Conns(), func(conn libp2p_network.Conn, _ int) peer.AddrInfo {
			return peer.AddrInfo{
				ID:    conn.RemotePeer(),
				Addrs: []maddr.Multiaddr{conn.RemoteMultiaddr()},
			}
		})
		p.Telemetry.SetConnectedPeers(connectedPeers)
		p.NumConnectedPeersMetric.Set(float64(len(connectedPeers)))
		return nil
	}

	runBackgroundTicker(ctx, pinger, p.Interval, "TSSHealthcheckPeersPing", logger)
	runBackgroundTicker(ctx, connectedPeersCounter, p.Interval, "TSSHealthcheckConnectedPeersCounter", logger)

	return nil
}

func runBackgroundTicker(
	ctx context.Context,
	task ticker.Task,
	interval time.Duration,
	name string,
	logger zerolog.Logger,
) {
	bgName := fmt.Sprintf("%sWorker", name)
	tickerName := fmt.Sprintf("%sTicker", name)

	bgTask := func(ctx context.Context) error {
		return ticker.Run(ctx, interval, task, ticker.WithLogger(logger, tickerName))
	}

	bg.Work(ctx, bgTask, bg.WithName(bgName), bg.WithLogger(logger))
}
