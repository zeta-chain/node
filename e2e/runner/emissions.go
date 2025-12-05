package runner

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/e2e/txserver"
	e2eutils "github.com/zeta-chain/node/e2e/utils"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// FundEmissionsPool funds the emissions pool on ZetaChain with the same value as used originally on mainnet (20M ZETA)
func (r *E2ERunner) FundEmissionsPool() error {
	r.Logger.Print("⚙️ funding the emissions pool on ZetaChain with 20M ZETA (%s)", txserver.EmissionsPoolAddress)

	return r.ZetaTxServer.FundEmissionsPool(e2eutils.OperationalPolicyName, EmissionsPoolFunding)
}

func (r *E2ERunner) RemoveObserver() error {
	observerSet, err := r.ObserverClient.ObserverSet(r.Ctx, &observertypes.QueryObserverSet{})
	if err != nil {
		return err
	}
	r.Logger.Print("🏃 Removing observer from the set : %s", observerSet.Observers[len(observerSet.Observers)-1])
	err = r.ZetaTxServer.RemoveObserver(observerSet.Observers[len(observerSet.Observers)-1])
	return err
}

// WithdrawEmissions withdraws emissions from the emission pool on ZetaChain for all observers
// This functions uses the UserEmissionsWithdrawName to create the withdraw tx.
// UserEmissionsWithdraw can sign the authz transactions because the necessary permissions are granted in the genesis file
func (r *E2ERunner) WithdrawEmissions() error {
	observerSet, err := r.ObserverClient.ObserverSet(r.Ctx, &observertypes.QueryObserverSet{})
	if err != nil {
		return err
	}

	for _, observer := range observerSet.Observers {
		r.Logger.Print("🏃 Withdrawing emissions for observer %s", observer)
		var (
			baseDenom            = config.BaseDenom
			queryObserverBalance = &banktypes.QueryBalanceRequest{
				Address: observer,
				Denom:   baseDenom,
			}
		)

		balanceBefore, err := r.BankClient.Balance(r.Ctx, queryObserverBalance)
		if err != nil {
			return errors.Wrapf(err, "failed to get balance for observer before withdrawing emissions %s", observer)
		}

		availableCoin, err := r.FetchWithdrawableEmissions(observer)
		if err != nil {
			return err
		}

		if availableCoin.Amount.IsZero() {
			r.Logger.Print("no emissions to withdraw for observer %s", observer)
			continue
		}

		if err := r.ZetaTxServer.WithdrawAllEmissions(availableCoin.Amount, e2eutils.UserEmissionsWithdrawName, observer); err != nil {
			r.Logger.Error("failed to withdraw emissions for observer %s: %s", observer, err)
			r.Logger.Error("Withdraw amount: %s", availableCoin.Amount)
			availableCoinAfter, fetchErr := r.FetchWithdrawableEmissions(observer)
			if fetchErr != nil {
				r.Logger.Error("failed to fetch available emissions for observer %s: %s", observer, fetchErr)
				return err
			}
			r.Logger.Error("Available emissions after failed withdrawal: %s", availableCoinAfter.Amount)
			return err
		}

		balanceAfter, err := r.BankClient.Balance(r.Ctx, queryObserverBalance)
		if err != nil {
			return errors.Wrapf(err, "failed to get balance for observer after withdrawing emissions %s", observer)
		}

		changeInBalance := balanceAfter.Balance.Sub(*balanceBefore.Balance).Amount
		if !changeInBalance.Equal(availableCoin.Amount) {
			return fmt.Errorf(
				"invalid balance change for observer %s, expected %s, got %s",
				observer,
				availableCoin.Amount,
				changeInBalance,
			)
		}
	}

	return nil
}

func (r *E2ERunner) FetchWithdrawableEmissions(observer string) (sdk.Coin, error) {
	availableAmount, err := r.EmissionsClient.ShowAvailableEmissions(
		r.Ctx,
		&emissionstypes.QueryShowAvailableEmissionsRequest{
			Address: observer,
		},
	)
	if err != nil {
		return sdk.Coin{}, errors.Wrapf(err, "failed to get available emissions for observer %s", observer)
	}

	availableCoin, err := sdk.ParseCoinNormalized(availableAmount.Amount)
	if err != nil {
		return sdk.Coin{}, errors.Wrap(err, "failed to parse coin amount")
	}
	return availableCoin, nil
}
