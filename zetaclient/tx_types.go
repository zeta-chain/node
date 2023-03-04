package zetaclient

type KeyType int

const (
	ValidatorKey KeyType = 1
	ObserverKey  KeyType = 2
)

type TxType string

const (
	NonceVoter    TxType = "NonceVoter"
	GasPriceVoter TxType = "GasPriceVoter"
)

type AuthZSigner struct {
	KeyType        KeyType
	MsgUrl         string
	GranterAddress string
	GranteeAddress string
}

var SignerList map[TxType]AuthZSigner

func init() {
	SignerList = map[TxType]AuthZSigner{
		NonceVoter: {
			KeyType: ObserverKey,
			MsgUrl:  "/zetachain.zetacore.crosschain.MsgNonceVoter",
		},
		GasPriceVoter: {
			KeyType: ObserverKey,
			MsgUrl:  "/zetachain.zetacore.crosschain.MsgGasPriceVoter",
		},
	}
}
