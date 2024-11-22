package cosmos

import (
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" // nolint
)

const Bech32PubKeyTypeAccPub = legacybech32.AccPK

var (
	GetPubKeyFromBech32 = legacybech32.UnmarshalPubKey
	Bech32ifyPubKey     = legacybech32.MarshalPubKey
)
