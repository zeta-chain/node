package main

import (
	ecdsakeygen "github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/rs/zerolog"

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

func InitLogger(level zerolog.Level) zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(level).
		With().
		Timestamp().
		Logger().
		Sample(&zerolog.BasicSampler{N: 5})
	return logger
}
