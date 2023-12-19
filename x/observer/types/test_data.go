package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/zeta-chain/zetacore/common"
)

const PrefixOutput = "Output1"

func CreateObserverMapperList(items int, chain common.Chain) (list []*ObserverMapper) {
	SetConfig(false)
	for i := 0; i < items; i++ {
		mapper := &ObserverMapper{
			Index:         "Index" + strconv.Itoa(i),
			ObserverChain: &chain,
			ObserverList: []string{
				sdk.AccAddress(crypto.AddressHash([]byte(PrefixOutput + strconv.Itoa(i)))).String(),
				sdk.AccAddress(crypto.AddressHash([]byte(PrefixOutput + strconv.Itoa(i+1)))).String(),
				sdk.AccAddress(crypto.AddressHash([]byte(PrefixOutput + strconv.Itoa(i+2)))).String(),
				sdk.AccAddress(crypto.AddressHash([]byte(PrefixOutput + strconv.Itoa(i+3)))).String(),
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
