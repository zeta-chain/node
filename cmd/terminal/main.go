package main

import (
	"bufio"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type SendInput struct {
	GasLimit    *big.Int
	DestChainID *big.Int
	To          ethcommon.Address
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	enabledChains := flag.String("enable-chains", "GOERLI,BSCTESTNET,MUMBAI,ROPSTEN", "enable chains, comma separated list")
	flag.Parse()
	chains := strings.Split(*enabledChains, ",")
	for _, chain := range chains {
		if c, err := common.ParseChain(chain); err == nil {
			config.ChainsEnabled = append(config.ChainsEnabled, c)
		} else {
			log.Error().Err(err).Msgf("invalid chain %s", chain)
			return
		}
	}

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
	for _, name := range config.ChainsEnabled {
		if name == "" || name.String() == common.RopstenChain.String() {
			continue
		}
		chain := config.Chains[name.String()]
		if endpoint := os.Getenv(name.String() + "_ENDPOINT"); endpoint != "" {
			chain.Endpoint = endpoint
		}
		client, err := ethclient.Dial(chain.Endpoint)
		if err != nil {
			log.Fatal().Err(err).Msgf("connect to chain error %s", name)
			continue
		}
		auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chain.ChainID)
		if err != nil {
			log.Fatal().Err(err).Msgf("bind keyed transactor error %s", name)
			continue
		}

		probe := NewProbe(client, &connABI, address, chain.ChainID, chain.ConnectorContractAddress, chain.ZETATokenContractAddress, auth)
		probeMap[name.String()] = probe

	}

	// setup signal handler Ctrl-C: cancel REPEAT command rather than terminate program
	log.Info().Msgf("awaiting the os.Interrupt, syscall.SIGTERM signals...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	log.Info().Msg("start REPL loop...")
	scanner := bufio.NewScanner(os.Stdin)
	probe := probeMap[common.GoerliChain.String()]
	probe.CurrentChain = common.GoerliChain
REPL_LOOP:
	for {
		fmt.Printf("[%s] > ", probe.CurrentChain)
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
			if len(cmdList) == 1 {
				for name, probe := range probeMap {
					log.Info().Msgf("%s: ChainID: %d", name, probe.ChainID)
					bal, err := probe.GetBalance()
					if err != nil {
						log.Error().Err(err).Msg("get balance error")
					} else {
						log.Info().Msgf("chain %s user account balance %s ETH/MATIC/BNB", name, bal)
					}

					bal, err = probe.GetZetaBalance()
					if err != nil {
						log.Error().Err(err).Msg("get zeta balance error")
					} else {
						log.Info().Msgf("chain %s zeta  balance %s ZETA", name, bal)
					}
				}
			}
		case "SET-CHAIN":
			if len(cmdList) != 2 {
				log.Error().Msg("SET-CHAIN <chain>")
				log.Info().Msgf("available chains: %v", config.ChainsEnabled)
				continue
			}
			chain, err := common.ParseChain(cmdList[1])
			if err != nil {
				log.Error().Msg("SET-CHAIN <chain>")
				log.Info().Msgf("available chains: %v", config.ChainsEnabled)
				continue
			}
			probe = probeMap[chain.String()]
		case "SEND":
			probe.evalSend(cmdList)
		case "REPEAT-SEND":
			if len(cmdList) < 4 {
				log.Error().Msg("REPEAT-SEND <repeat times> <duration> SEND-CMD")
				continue
			}
			repeat, err := strconv.ParseInt(cmdList[1], 10, 64)
			if err != nil {
				log.Error().Err(err).Msgf("wrong repeat number %s", cmdList[1])
				continue
			}
			duration, err := time.ParseDuration(cmdList[2])
			if err != nil {
				log.Error().Err(err).Msgf("wrong duration %s", cmdList[2])
				continue
			}
			interval := duration.Seconds() / float64(repeat)
			if interval < 1 {
				interval = 1
				log.Warn().Msgf("interval %f <1s; set to 1s", interval)
			}
			intervalSeconds := int64(interval)
			log.Info().Msgf("interval %d seconds; begin", intervalSeconds)
			ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
		REPEAT:
			for i := 0; i < int(repeat); i++ {
				log.Info().Msgf("iteration %d/%d", i, repeat)
				select {
				case <-ch:
					break REPEAT
				case <-ticker.C:
					probe.evalSend(cmdList[3:])
				}
			}
		default:
			log.Error().Msg("wrong command")
		}

	}
}

func (probe *Probe) evalSend(cmdList []string) {
	sendInput := &SendInput{
		GasLimit:    big.NewInt(90_000),
		DestChainID: config.Chains[probe.CurrentChain.String()].ChainID,
		To:          probe.Address,
	}
	for _, kv := range cmdList[1:] {
		err := SetSendInput(sendInput, kv)
		if err != nil {
			log.Error().Err(err).Msgf("incorrect SEND parameter %s; should be key=value", kv)
		}
	}

	log.Info().Msgf("send %s to %s", probe.ChainID, sendInput.DestChainID)
	log.Info().Msgf("sendInput %+v", sendInput)
	err := probe.SendTransaction(sendInput)
	if err != nil {
		log.Error().Err(err).Msgf("SendTransaction error")
	}
}

func SetSendInput(sendInput *SendInput, kv string) error {
	kvList := strings.Split(kv, "=")
	if len(kvList) != 2 {
		return fmt.Errorf("wrong format")
	}
	switch kvList[0] {
	case "GasLimit":
		if gasLimit, ok := big.NewInt(0).SetString(kvList[1], 10); !ok {
			return fmt.Errorf("wrong GasLimit %s", kvList)
		} else {
			sendInput.GasLimit = gasLimit
		}
	case "DestChainID":
		if chainID, ok := big.NewInt(0).SetString(kvList[1], 10); !ok {
			return fmt.Errorf("wrong GasLimit %s", kvList)
		} else {
			sendInput.DestChainID = chainID
		}
	}
	return nil
}
