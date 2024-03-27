package cosmos

import (
	"github.com/cosmos/cosmos-sdk/types/bech32/legacybech32" // nolint
)

var (
	GetPubKeyFromBech32    = legacybech32.UnmarshalPubKey
	Bech32ifyPubKey        = legacybech32.MarshalPubKey
	Bech32PubKeyTypeAccPub = legacybech32.AccPK
)
