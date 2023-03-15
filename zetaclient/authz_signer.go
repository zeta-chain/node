package zetaclient

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

type AuthZSigner struct {
	KeyType        common.KeyType
	GranterAddress string
	GranteeAddress sdk.AccAddress
}

func (a AuthZSigner) String() string {
	return a.KeyType.String() + " " + a.GranterAddress + " " + a.GranteeAddress.String()
}

var signers map[string]AuthZSigner

func init() {
	signersList := make(map[string]AuthZSigner)
	for _, tx := range crosschaintypes.GetAllAuthzTxTypes() {
		signersList[tx] = AuthZSigner{KeyType: common.ObserverGranteeKey}
	}
	signers = signersList
}

func SetupAuthZSignerList(granter string, grantee sdk.AccAddress) {
	for k, v := range signers {
		v.GranterAddress = granter
		v.GranteeAddress = grantee
		signers[k] = v
	}
}

func GetSigner(msgUrl string) AuthZSigner {
	return signers[msgUrl]
}
