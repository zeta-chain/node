package runner

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// UnStakeToBelowMinimumObserverDelegation  unstakes the validator self delegation to below the minimum observer delegation
func (r *E2ERunner) UnStakeToBelowMinimumObserverDelegation() error {
	// Arrange
	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()

	// List keys to verify that the keyring is working
	validatorKeys, err := validatorsKeyring.List()
	if err != nil {
		return fmt.Errorf("failed to list validator keys: %w", err)
	}
	if len(validatorKeys) == 0 {
		return fmt.Errorf("no validator keys found")
	}
	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	// Pick the first validator
	unstakeFrom := validatorKeys[0]
	unstakeFromAddress, err := unstakeFrom.GetAddress()
	if err != nil {
		return fmt.Errorf("failed to get address for depositor: %w", err)
	}

	valAddress, err := observertypes.GetOperatorAddressFromAccAddress(unstakeFromAddress.String())
	if err != nil {
		return fmt.Errorf("failed to get operator address from account address: %w", err)
	}

	// Fetch the current self delegation
	resGetDelegation, err := r.StakingClient.Delegation(r.Ctx, &stakingtypes.QueryDelegationRequest{
		DelegatorAddr: unstakeFromAddress.String(),
		ValidatorAddr: valAddress.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to fetch self delegation   %w", err)
	}

	// We need to unstake to just below the minimum observer delegation to trigger the hooks which would remove the observer
	// This works as expected , however the hardcoded value of 1zeta (1000000000000000000azeta), for the `min_observer_delegation` is not ideal, for this check
	// The minimum accepted value for `MIN_SELF_DELEGATION` is also 1zeta , therefore this unstake message ends up removing the validator from the validator set as well.
	// NOTE : although the MIN_SELF_DELEGATION is set to 1, it does not mean 1azeta , when calculating the default power reduction is accounted for, therefore 1ZETA = 1 unit of voting power and not 1 azeta
	// Ideally we should be able to test removing the observer from the observer set, without affecting the validator set.
	// This can be improved by adding the MinObserverDelegation value to params and making it dynamic
	// https://github.com/zeta-chain/node/issues/3550
	delegation := resGetDelegation.DelegationResponse.Balance.Amount
	minDelegation, _ := observertypes.GetMinObserverDelegation()
	unstakeAmount := delegation.Sub(minDelegation).Add(sdkmath.NewInt(1))

	// Act
	msg := stakingtypes.MsgUndelegate{
		DelegatorAddress: unstakeFromAddress.String(),
		ValidatorAddress: valAddress.String(),
		Amount:           sdk.NewCoin(resGetDelegation.DelegationResponse.Balance.Denom, unstakeAmount),
	}

	_, err = validatorsTxServer.BroadcastTx(unstakeFrom.Name, &msg)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction for proposal  %w", err)
	}

	// Wait for the transaction to be included in a block
	r.WaitForBlocks(2)

	// Assert
	// Check if the observer has been removed from the observer set
	resGetObserverSet, err := r.ObserverClient.ObserverSet(r.Ctx, &observertypes.QueryObserverSet{})
	if err != nil {
		return fmt.Errorf("failed to fetch observer set %w", err)
	}

	for _, observer := range resGetObserverSet.Observers {
		if observer == unstakeFromAddress.String() {
			return fmt.Errorf("observer not removed from observer set")
		}
	}

	// Check TSS keygen is updated
	resGetKeygen, err := r.ObserverClient.Keygen(r.Ctx, &observertypes.QueryGetKeygenRequest{})
	if err != nil {
		return fmt.Errorf("failed to fetch keygen %w", err)
	}

	if resGetKeygen.Keygen.Status != observertypes.KeygenStatus_PendingKeygen {
		return fmt.Errorf("keygen status not updated")
	}

	// Check inbound is disabled
	resGetCrosschainFlags, err := r.ObserverClient.CrosschainFlags(
		r.Ctx,
		&observertypes.QueryGetCrosschainFlagsRequest{},
	)
	if err != nil {
		return fmt.Errorf("failed to fetch crosschain flags %w", err)
	}

	if resGetCrosschainFlags.CrosschainFlags.IsInboundEnabled {
		return fmt.Errorf("inbound not disabled")
	}
	return nil
}
