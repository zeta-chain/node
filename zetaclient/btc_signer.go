package zetaclient

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/big"
)

type BTCSigner struct {
	tssSigner    *TestSigner
	logger       zerolog.Logger
	utxos        []*btcutil.Tx   // sorted utxos
	pendingUxtos map[string]bool // need to persist
}

func NewBTCSigner(tssSigner *TestSigner) (*BTCSigner, error) {
	return &BTCSigner{
		tssSigner: tssSigner,
		logger:    log.With().Str("module", "BTCSigner").Logger(),
	}, nil
}

// TODO:
func (signer *BTCSigner) SignWithdrawTx(to btcutil.AddressPubKeyHash, sat *big.Int, satPerKB *big.Int) (*wire.MsgTx, error) {
	// sort utxo by value in ascending order
	// select N utxo sufficient to cover the amount

	return nil, nil
}

// TODO:
func (signer *BTCSigner) Broadcast(signedTx *wire.MsgTx) error {
	return nil
}
