// Package authz provides a signer object for transactions using grants
// grants are used to allow a hotkey to sign transactions on behalf of the observers
package authz

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/authz"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// Signer represents a signer for a grantee key
type Signer struct {
	KeyType        authz.KeyType
	GranterAddress string
	GranteeAddress sdk.AccAddress
}

// String returns a string representation of a Signer
func (a Signer) String() string {
	return a.KeyType.String() + " " + a.GranterAddress + " " + a.GranteeAddress.String()
}

// signers is a map of all the signers for the different tx types
var signers map[string]Signer

// init initializes the signers map with all the crosschain tx types using the ZetaClientGranteeKey
func init() {
	signersList := make(map[string]Signer)
	for _, tx := range crosschaintypes.GetAllAuthzZetaclientTxTypes() {
		signersList[tx] = Signer{KeyType: authz.ZetaClientGranteeKey}
	}
	signers = signersList
}

// SetupAuthZSignerList sets the granter and grantee for all the signers
func SetupAuthZSignerList(granter string, grantee sdk.AccAddress) {
	for k, v := range signers {
		v.GranterAddress = granter
		v.GranteeAddress = grantee
		signers[k] = v
	}
}

// GetSigner returns the signer for a given msgURL
func GetSigner(msgURL string) Signer {
	return signers[msgURL]
}
