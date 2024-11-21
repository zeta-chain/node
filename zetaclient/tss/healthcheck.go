package tss

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"gitlab.com/thorchain/tss/go-tss/tss"

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
func HealthcheckWorker(ctx context.Context, server *tss.TssServer, p HealthcheckProps, logger zerolog.Logger) error {
	if p.NumConnectedPeersMetric == nil {
		return errors.New("missing NumConnectedPeersMetric")
	}

	if p.Interval == 0 {
		p.Interval = 30 * time.Second
	}

	logger = logger.With().Str(logs.FieldModule, "tss_healthcheck").Logger()

	// Ping & collect round trip time
	var (
		host    = server.GetP2PHost()
		pingRTT = make(map[peer.ID]int64)
		mu      = sync.Mutex{}
	)

	const pingTimeout = 5 * time.Second

	pinger := func(ctx context.Context, _ *ticker.Ticker) error {
		var wg sync.WaitGroup
		for i := range p.WhitelistPeers {
			peerID := p.WhitelistPeers[i]
			if peerID == host.ID() {
				continue
			}

			wg.Add(1)

			go func() {
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
					logger.Error().Str("peer_id", peerID.String()).Err(result.Error).Msg("ping error")
				}

				mu.Lock()
				pingRTT[peerID] = result.RTT.Nanoseconds()
				mu.Unlock()
			}()

			wg.Wait()
			p.Telemetry.SetPingRTT(pingRTT)
		}

		return nil
	}

	peersCounter := func(_ context.Context, _ *ticker.Ticker) error {
		peers := server.GetKnownPeers()
		p.NumConnectedPeersMetric.Set(float64(len(peers)))
		p.Telemetry.SetConnectedPeers(peers)

		return nil
	}

	runBackgroundTicker(ctx, pinger, p.Interval, "TSSHealthcheckPeersPing", logger)
	runBackgroundTicker(ctx, peersCounter, p.Interval, "TSSHealthcheckPeersCounter", logger)

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
