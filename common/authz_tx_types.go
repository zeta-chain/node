package common

type TxType string

const (
	InboundVoter  TxType = "InboundVoter"
	OutboundVoter TxType = "OutboundVoter"
	NonceVoter    TxType = "NonceVoter"
	GasPriceVoter TxType = "GasPriceVoter"
)

func (t TxType) String() string {
	return string(t)
}

type KeyType string

const (
	TssSignerKey        KeyType = "tss_signer"
	ValidatorGranteeKey KeyType = "grantee_validator"
	ObserverGranteeKey  KeyType = "grantee_observer"
)

func GetAllKeyTypes() []KeyType {
	return []KeyType{ValidatorGranteeKey, ObserverGranteeKey, TssSignerKey}
}

func (k KeyType) String() string {
	return string(k)
}
