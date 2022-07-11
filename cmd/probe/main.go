package main

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"os"
	"strings"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	privkey := os.Getenv("PRIVATE_KEY")
	if privkey == "" {
		log.Fatal().Msg("envvar PRIVATE_KEY is not set")
		return
	}
	privateKey, err := crypto.HexToECDSA(privkey)
	if err != nil {
		log.Fatal().Err(err).Msg("parse private key error")
		return
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal().Msg("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	connABI, err := abi.JSON(strings.NewReader(config.CONNECTOR_ABI_STRING))
	if err != nil {
		log.Fatal().Err(err).Msg("parse connector abi error")
	}

	probeMap := make(map[string]*Probe)
	for name, chain := range config.Chains {
		if name == "" || name == common.RopstenChain.String() {
			continue
		}
		if endpoint := os.Getenv(name + "_ENDPOINT"); endpoint != "" {
			chain.Endpoint = endpoint
		}
		client, err := ethclient.Dial(chain.Endpoint)
		if err != nil {
			log.Fatal().Err(err).Msgf("connect to chain error %s", name)
			continue
		}
		probe := NewProbe(client, &connABI, address, chain.ChainID, chain.ConnectorContractAddress, chain.ZETATokenContractAddress)
		probeMap[name] = probe

	}

	log.Info().Msg("start REPL loop...")
	scanner := bufio.NewScanner(os.Stdin)
REPL_LOOP:
	for {
		fmt.Printf("> ")
		scanner.Scan()
		cmd := scanner.Text()
		cmdList := strings.Fields(cmd)
		if len(cmdList) == 0 {
			continue
		}

		switch cmdList[0] {
		case "EXIT":
			break REPL_LOOP
		case "INFO":
			for name, probe := range probeMap {
				bal, err := probe.GetBalance()
				if err != nil {
					log.Error().Err(err).Msg("get balance error")
				} else {
					log.Info().Msgf("chain %s user account balance %s ETH/MATIC/BAL", name, bal)
				}

				bal, err = probe.GetZetaBalance()
				if err != nil {
					log.Error().Err(err).Msg("get zeta balance error")
				} else {
					log.Info().Msgf("chain %s zeta  balance %s ZETA", name, bal)
				}
			}
		}

	}
}
