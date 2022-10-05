package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"strconv"
)

func CreateObserverMapperList(items int, chain ObserverChain, observationType ObservationType) (list []*ObserverMapper) {
	SetConfig(false)
	for i := 0; i < items; i++ {
		mapper := &ObserverMapper{
			Index:           "Index" + strconv.Itoa(i),
			ObserverChain:   chain,
			ObservationType: observationType,
			ObserverList: []string{
				sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i)))).String(),
				sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i+1)))).String(),
				sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i+2)))).String(),
				sdk.AccAddress(crypto.AddressHash([]byte("Output1" + strconv.Itoa(i+3)))).String(),
			},
		}
		list = append(list, mapper)
	}
	return
}

const (
	AccountAddressPrefix = "zeta"
)

var (
	AccountPubKeyPrefix    = AccountAddressPrefix + "pub"
	ValidatorAddressPrefix = AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix  = AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix  = AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix   = AccountAddressPrefix + "valconspub"
)

func SetConfig(seal bool) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	if seal {
		config.Seal()
	}
}
