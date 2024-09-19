package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/zetaclient/config"
)

func waitForZetaCore(config config.Config, logger zerolog.Logger) {
	// wait until zetacore is up
	logger.Debug().Msg("Waiting for zetacore to open 9090 port...")
	for {
		_, err := grpc.Dial(
			fmt.Sprintf("%s:9090", config.ZetaCoreURL),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Warn().Err(err).Msg("grpc dial fail")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

func validatePeer(seedPeer string) error {
	parsedPeer := strings.Split(seedPeer, "/")

	if len(parsedPeer) < 7 {
		return errors.New("seed peer missing IP or ID or both, seed: " + seedPeer)
	}

	seedIP := parsedPeer[2]
	seedID := parsedPeer[6]

	if net.ParseIP(seedIP) == nil {
		return errors.New("invalid seed IP address format, seed: " + seedPeer)
	}

	if len(seedID) == 0 {
		return errors.New("seed id is empty, seed: " + seedPeer)
	}

	return nil
}
