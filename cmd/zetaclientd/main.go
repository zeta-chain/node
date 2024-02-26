package main

import (
	"path/filepath"

	"github.com/rs/zerolog/log"

	ecdsakeygen "github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/rs/zerolog"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
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

const (
	ComplianceLogFile = "compliance.log"
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

func InitLogger(cfg *config.Config) (clientcommon.ClientLogger, error) {
	// open compliance log file
	file, err := OpenComplianceLogFile(cfg)
	if err != nil {
		return clientcommon.DefaultLoggers(), err
	}

	var logger zerolog.Logger
	var loggerCompliance zerolog.Logger
	switch cfg.LogFormat {
	case "json":
		logger = zerolog.New(os.Stdout).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
		loggerCompliance = zerolog.New(file).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
	case "text":
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
		loggerCompliance = zerolog.New(file).Level(zerolog.Level(cfg.LogLevel)).With().Timestamp().Logger()
	default:
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
		loggerCompliance = zerolog.New(file).With().Timestamp().Logger()
	}

	if cfg.LogSampler {
		logger = logger.Sample(&zerolog.BasicSampler{N: 5})
	}
	log.Logger = logger // set global logger

	return clientcommon.ClientLogger{
		Std:        log.Logger,
		Compliance: loggerCompliance,
	}, nil
}

func OpenComplianceLogFile(cfg *config.Config) (*os.File, error) {
	// use zetacore home as default
	logPath := cfg.ZetaCoreHome
	if cfg.ComplianceConfig != nil && cfg.ComplianceConfig.LogPath != "" {
		logPath = cfg.ComplianceConfig.LogPath
	}

	// clean file name
	name := filepath.Join(logPath, ComplianceLogFile)
	name, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	name = filepath.Clean(name)

	// open (or create) compliance log file
	return os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
