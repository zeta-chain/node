package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/zeta-chain/zetacore/zetaclient/metrics"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/go-tss/p2p"
	"github.com/zeta-chain/zetacore/pkg/cosmos"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

func RunDiagnostics(startLogger zerolog.Logger, peers p2p.AddrList, hotkeyPk cryptotypes.PrivKey, cfg config.Config) error {

	startLogger.Warn().Msg("P2P Diagnostic mode enabled")
	startLogger.Warn().Msgf("seed peer: %s", peers)
	priKey := secp256k1.PrivKey(hotkeyPk.Bytes()[:32])
	pubkeyBech32, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, hotkeyPk.PubKey())
	if err != nil {
		startLogger.Error().Err(err).Msg("Bech32ifyPubKey error")
		return err
	}
	startLogger.Warn().Msgf("my pubkey %s", pubkeyBech32)

	var s *metrics.TelemetryServer
	if len(peers) == 0 {
		startLogger.Warn().Msg("No seed peer specified; assuming I'm the host")

	}
	p2pPriKey, err := crypto.UnmarshalSecp256k1PrivateKey(priKey[:])
	if err != nil {
		startLogger.Error().Err(err).Msg("UnmarshalSecp256k1PrivateKey error")
		return err
	}
	listenAddress, err := maddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 6668))
	if err != nil {
		startLogger.Error().Err(err).Msg("NewMultiaddr error")
		return err
	}
	IP := os.Getenv("MYIP")
	if len(IP) == 0 {
		startLogger.Warn().Msg("empty env MYIP")
	}
	var externalAddr Multiaddr
	if len(IP) != 0 {
		externalAddr, err = maddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", IP, 6668))
		if err != nil {
			startLogger.Error().Err(err).Msg("NewMultiaddr error")
			return err
		}
	}

	host, err := libp2p.New(
		libp2p.ListenAddrs(listenAddress),
		libp2p.Identity(p2pPriKey),
		libp2p.AddrsFactory(func(addrs []Multiaddr) []Multiaddr {
			if externalAddr != nil {
				return []Multiaddr{externalAddr}
			}
			return addrs
		}),
		libp2p.DisableRelay(),
	)
	if err != nil {
		startLogger.Error().Err(err).Msg("fail to create host")
		return err
	}
	startLogger.Info().Msgf("host created: ID %s", host.ID().String())
	if len(peers) == 0 {
		s = metrics.NewTelemetryServer()
		s.SetP2PID(host.ID().String())
		go func() {
			startLogger.Info().Msg("Starting TSS HTTP Server...")
			if err := s.Start(); err != nil {
				fmt.Println(err)
			}
		}()
	}

	// create stream handler
	handleStream := func(s network.Stream) {
		defer s.Close()

		// read the message
		buf := make([]byte, 1024)
		n, err := s.Read(buf)
		if err != nil {
			startLogger.Error().Err(err).Msg("read stream error")
			return
		}
		// send the message back
		if _, err := s.Write(buf[:n]); err != nil {
			startLogger.Error().Err(err).Msg("write stream error")
			return
		}
	}
	ProtocolID := "/echo/0.3.0"
	host.SetStreamHandler(protocol.ID(ProtocolID), handleStream)

	kademliaDHT, err := dht.New(context.Background(), host, dht.Mode(dht.ModeServer))
	if err != nil {
		return fmt.Errorf("fail to create DHT: %w", err)
	}
	startLogger.Info().Msg("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(context.Background()); err != nil {
		return fmt.Errorf("fail to bootstrap DHT: %w", err)
	}

	var wg sync.WaitGroup
	for _, peerAddr := range peers {
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			startLogger.Error().Err(err).Msgf("fail to parse peer address %s", peerAddr)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(context.Background(), *peerinfo); err != nil {
				startLogger.Warn().Msgf("Connection failed with bootstrap node: %s", *peerinfo)
			} else {
				startLogger.Info().Msgf("Connection established with bootstrap node: %s", *peerinfo)
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	startLogger.Info().Msgf("Announcing ourselves...")
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(context.Background(), routingDiscovery, "ZetaZetaOpenTheDoor")
	startLogger.Info().Msgf("Successfully announced!")

	// every 1min, print out the p2p diagnostic
	ticker := time.NewTicker(time.Duration(cfg.P2PDiagnosticTicker) * time.Second)
	round := 0
	for {
		select {
		case <-ticker.C:
			round++
			// Now, look for others who have announced
			// This is like your friend telling you the location to meet you.
			startLogger.Info().Msgf("Searching for other peers...")
			peerChan, err := routingDiscovery.FindPeers(context.Background(), "ZetaZetaOpenTheDoor")
			if err != nil {
				panic(err)
			}

			peerCount := 0
			okPingPongCount := 0
			for peer := range peerChan {
				peerCount++
				if peer.ID == host.ID() {
					startLogger.Info().Msgf("Found myself #(%d): %s", peerCount, peer)
					continue
				}
				startLogger.Info().Msgf("Found peer #(%d): %s; pinging the peer...", peerCount, peer)
				stream, err := host.NewStream(context.Background(), peer.ID, protocol.ID(ProtocolID))
				if err != nil {
					startLogger.Error().Err(err).Msgf("fail to create stream to peer %s", peer)
					continue
				}
				message := fmt.Sprintf("round %d %s => %s", round, host.ID().String()[len(host.ID().String())-5:], peer.ID.String()[len(peer.ID.String())-5:])
				_, err = stream.Write([]byte(message))
				if err != nil {
					startLogger.Error().Err(err).Msgf("fail to write to stream to peer %s", peer)
					err = stream.Close()
					if err != nil {
						startLogger.Warn().Err(err).Msgf("fail to close stream to peer %s", peer)
					}
					continue
				}
				//startLogger.Debug().Msgf("wrote %d bytes", nw)
				buf := make([]byte, 1024)
				nr, err := stream.Read(buf)
				if err != nil {
					startLogger.Error().Err(err).Msgf("fail to read from stream to peer %s", peer)
					err = stream.Close()
					if err != nil {
						startLogger.Warn().Err(err).Msgf("fail to close stream to peer %s", peer)
					}
					continue
				}
				//startLogger.Debug().Msgf("read %d bytes", nr)
				startLogger.Debug().Msgf("echoed message: %s", string(buf[:nr]))
				err = stream.Close()
				if err != nil {
					startLogger.Warn().Err(err).Msgf("fail to close stream to peer %s", peer)
				}

				if string(buf[:nr]) != message {
					startLogger.Error().Msgf("ping-pong failed with peer #(%d): %s; want %s got %s", peerCount, peer, message, string(buf[:nr]))
					continue
				}
				startLogger.Info().Msgf("ping-pong success with peer #(%d): %s;", peerCount, peer)
				okPingPongCount++
			}
			startLogger.Info().Msgf("Expect %d peers in total; successful pings (%d/%d)", peerCount, okPingPongCount, peerCount-1)
		}
	}
}
