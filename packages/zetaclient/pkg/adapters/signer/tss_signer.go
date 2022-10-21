package signer

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/go-tss-ctx/tss"
)

type TSSSigner interface {
	Server() *tss.TssServer
	Pubkey() []byte
	Sign(data []byte) ([65]byte, error)
	Address() ethcommon.Address // TODO : transform in generic address
	InsertPubKey(string) error
	SetCurrentPubKey(string)
	CurrentPubKey() string
}
