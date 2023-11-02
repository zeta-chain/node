package main

import (
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/rs/zerolog"
	ecdsakeygen "github.com/zeta-chain/tss-lib/ecdsa/keygen"
	"github.com/zeta-chain/zetacore/zetaclient/config"

	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common/cosmos"

	//mcconfig "github.com/Meta-Protocol/zetacore/metaclient/config"
	"github.com/cosmos/cosmos-sdk/types"

	"math/rand"
	"os"
	"time"

	"github.com/zeta-chain/zetacore/app"
)

var (
	preParams *ecdsakeygen.LocalPreParams
)

func main() {
	if err := svrcmd.Execute(RootCmd, "", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func SetupConfigForTest() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	//config.SetCoinType(cmd.MetaChainCoinType)
	config.SetFullFundraiserPath(cmd.ZetaChainHDPath)
	types.SetCoinDenomRegex(func() string {
		return cmd.DenomRegex
	})

	rand.Seed(time.Now().UnixNano())

}

func InitLogger(cfg *config.Config) zerolog.Logger {
	var logger zerolog.Logger
	switch cfg.LogFormat {
	case "json":
		logger = zerolog.New(os.Stdout).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
	case "text":
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
	default:
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	if cfg.LogSampler {
		logger = logger.Sample(&zerolog.BasicSampler{N: 5})
	}
	return logger
}
