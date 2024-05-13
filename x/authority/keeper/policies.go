package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/authorizations"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

// SetPolicies sets the policies to the store
func (k Keeper) SetPolicies(ctx sdk.Context, policies types.Policies) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoliciesKey))
	b := k.cdc.MustMarshal(&policies)
	store.Set([]byte{0}, b)
}

// GetPolicies returns the policies from the store
func (k Keeper) GetPolicies(ctx sdk.Context) (val types.Policies, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoliciesKey))
	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// IsAuthorized checks if the message has been signed by an authorized policy address
func (k Keeper) CheckAuthorization(ctx sdk.Context, msg sdk.Msg) error {
	// Policy transactions must have only one signer
	if len(msg.GetSigners()) != 1 {
		return errors.Wrap(types.ErrSigners, fmt.Sprintf("msg: %v", sdk.MsgTypeURL(msg)))
	}
	signer := msg.GetSigners()[0].String()
	policyRequired := authorizations.GetRequiredPolicy(sdk.MsgTypeURL(msg))
	if policyRequired == types.PolicyType_emptyPolicyType {
		return errors.Wrap(types.ErrMsgNotAuthorized, fmt.Sprintf("msg: %v", sdk.MsgTypeURL(msg)))
	}
	policies, found := k.GetPolicies(ctx)
	if !found {
		return errors.Wrap(types.ErrPoliciesNotFound, fmt.Sprintf("msg: %v", sdk.MsgTypeURL(msg)))
	}
	for _, policy := range policies.Items {
		if policy.Address == signer && policy.PolicyType == policyRequired {
			return nil
		}
	}
	return errors.Wrap(types.ErrSignerDoesntMatch, fmt.Sprintf("signer: %s, policy required for message: %s , msg %s",
		signer, policyRequired.String(), sdk.MsgTypeURL(msg)))
}
