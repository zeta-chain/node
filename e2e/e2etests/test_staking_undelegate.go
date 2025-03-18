package e2etests

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// UndelegateToBelowMinimumObserverDelegation  undelegates the validator self delegation to below the minimum observer delegation
func UndelegateToBelowMinimumObserverDelegation(r *runner.E2ERunner, _ []string) {
	r.Logger.Print("running staking tests")
	// Arrange
	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()

	// List keys to verify that the keyring is working
	validatorKeys, err := validatorsKeyring.List()
	require.NoError(r, err, "failed to list validator keys")

	require.Greater(r, len(validatorKeys), 0, "no validator keys found")

	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	// Pick the first validator
	undelegateFrom := validatorKeys[0]
	undelegateFromAddress, err := undelegateFrom.GetAddress()
	require.NoError(r, err, "failed to get address for depositor")

	valAddress, err := observertypes.GetOperatorAddressFromAccAddress(undelegateFromAddress.String())
	require.NoError(r, err, "failed to get operator address from account address")

	// Fetch the current self delegation
	resGetDelegation, err := r.StakingClient.Delegation(r.Ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: undelegateFromAddress.String(),
		ValidatorAddr: valAddress.String(),
	})
	require.NoError(r, err, "failed to fetch self delegation")

	// We need to undelegate to just below the minimum observer delegation to trigger the hooks which would remove the observer
	// This works as expected , however the hardcoded value of 1zeta (1000000000000000000azeta), for the `min_observer_delegation` is not ideal, for this check
	// The minimum accepted value for `MIN_SELF_DELEGATION` is also 1zeta , therefore this `MsgUndelegate` message ends up removing the validator from the validator set as well.
	// NOTE : although the MIN_SELF_DELEGATION is set to 1, it does not mean 1azeta , when calculating the default power reduction is accounted for, therefore 1ZETA = 1 unit of voting power and not 1 azeta
	// Ideally we should be able to test removing the observer from the observer set, without affecting the validator set.
	// This can be improved by adding the MinObserverDelegation value to params and making it dynamic
	// https://github.com/zeta-chain/node/issues/3550
	delegation := resGetDelegation.DelegationResponse.Balance.Amount
	minDelegation, _ := observertypes.GetMinObserverDelegation()
	undelegateAmount := delegation.Sub(minDelegation).Add(sdkmath.NewInt(1))

	// Act
	msg := stakingtypes.MsgUndelegate{
		DelegatorAddress: undelegateFromAddress.String(),
		ValidatorAddress: valAddress.String(),
		Amount:           sdk.NewCoin(resGetDelegation.DelegationResponse.Balance.Denom, undelegateAmount),
	}

	_, err = validatorsTxServer.BroadcastTx(undelegateFrom.Name, &msg)
	require.NoError(r, err, "failed to broadcast transaction for proposal")

	// Wait for the transaction to be included in a block
	r.WaitForBlocks(2)

	// Assert
	// Check if the observer has been removed from the observer set
	resGetObserverSet, err := r.ObserverClient.ObserverSet(r.Ctx, &observertypes.QueryObserverSet{})
	require.NoError(r, err, "failed to fetch observer set")

	for _, observer := range resGetObserverSet.Observers {
		require.NotEqual(r, observer, undelegateFromAddress.String(), "observer not removed from observer set")
	}

	// Check TSS keygen is updated
	resGetKeygen, err := r.ObserverClient.Keygen(r.Ctx, &observertypes.QueryGetKeygenRequest{})
	require.NoError(r, err, "failed to fetch keygen")

	require.Equal(r, observertypes.KeygenStatus_PendingKeygen, resGetKeygen.Keygen.Status, "keygen status not updated")

	// Check inbound is disabled
	resGetCrosschainFlags, err := r.ObserverClient.CrosschainFlags(
		r.Ctx,
		&observertypes.QueryGetCrosschainFlagsRequest{},
	)
	require.NoError(r, err, "failed to fetch crosschain flags")

	require.False(r, resGetCrosschainFlags.CrosschainFlags.IsInboundEnabled, "inbound not disabled")
}
