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

var signerList map[string]AuthZSigner

func init() {
	for _, tx := range crosschaintypes.GetAllAuthzTxTypes() {
		signerList[tx] = AuthZSigner{KeyType: common.ObserverGranteeKey}
	}
}

func SetupAuthZSignerList(granter string, grantee sdk.AccAddress) {
	for k, v := range signerList {
		v.GranterAddress = granter
		v.GranteeAddress = grantee
		signerList[k] = v
	}
}

func GetSigner(msgUrl string) AuthZSigner {
	return signerList[msgUrl]
}
