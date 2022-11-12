package zetaclient

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type BTCSigner struct {
	//client              *ethclient.Client
	//chain               common.Chain
	//chainID             *big.Int
	tssSigner *TestSigner
	//ethSigner           ethtypes.Signer
	logger zerolog.Logger
}

func NewBTCSigner(tssSigner *TestSigner) (*BTCSigner, error) {
	return &BTCSigner{
		tssSigner: tssSigner,
		logger:    log.With().Str("module", "BTCSigner").Logger(),
	}, nil
}
