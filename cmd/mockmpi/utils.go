package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/zetaclient"
)

func GetZetaTestSignature() zetaclient.TestSigner {
	pkstring := os.Getenv("PRIVKEY")
	if pkstring == "" {
		log.Fatal().Msg("missing env variable PRIVKEY")
		os.Exit(1)
	}
	privateKey, err := crypto.HexToECDSA(pkstring)
	if err != nil {
		log.Err(err).Msg("TEST private key error")
		os.Exit(1)
	}
	tss := zetaclient.TestSigner{
		PrivKey: privateKey,
	}
	log.Debug().Msg(fmt.Sprintf("tss key address: %s", tss.EVMAddress()))

	return tss
}
