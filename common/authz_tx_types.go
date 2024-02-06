package common

type TxType string

const (
	InboundVoter  TxType = "InboundVoter"
	OutboundVoter TxType = "OutboundVoter"
	NonceVoter    TxType = "NonceVoter"
	GasPriceVoter TxType = "GasPriceVoter"
	RefundAborted TxType = "RefundAborted"
)

func (t TxType) String() string {
	return string(t)
}

type KeyType string

const (
	TssSignerKey         KeyType = "tss_signer"
	ValidatorGranteeKey  KeyType = "validator_grantee"
	ZetaClientGranteeKey KeyType = "zetaclient_grantee"
)

func GetAllKeyTypes() []KeyType {
	return []KeyType{ValidatorGranteeKey, ZetaClientGranteeKey, TssSignerKey}
}

func (k KeyType) String() string {
	return string(k)
}
