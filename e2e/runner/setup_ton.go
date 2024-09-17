package runner

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/e2e/runner/ton"
)

// SetupTON setups TON deployer and deploys Gateway contract
func (r *E2ERunner) SetupTON() error {
	if r.Clients.TON == nil {
		return fmt.Errorf("TON clients are not initialized")
	}

	ctx := r.Ctx

	// 1. Setup Deployer (acts as a faucet as well)
	faucetConfig, err := r.Clients.TON.GetFaucet(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get faucet config")
	}

	deployer, err := ton.NewDeployer(r.Clients.TON, faucetConfig)
	if err != nil {
		return errors.Wrap(err, "unable to create TON deployer")
	}

	depAddr := deployer.GetAddress()
	r.Logger.Print("💎TON Deployer %s (%s)", depAddr.ToRaw(), depAddr.ToHuman(false, true))

	gwAccount, err := ton.ConstructGatewayAccount(r.TSSAddress)
	if err != nil {
		return errors.Wrap(err, "unable to initialize TON gateway")
	}

	// 2. Deploy Gateway
	initStateAmount := ton.TONCoins(10)

	if err := deployer.Deploy(ctx, gwAccount, initStateAmount); err != nil {
		return errors.Wrapf(err, "unable to deploy TON gateway")
	}

	r.Logger.Print("💎TON Gateway deployed %s (%s)", gwAccount.ID.ToRaw(), gwAccount.ID.ToHuman(false, true))

	// 3. Check that the gateway indeed was deployed and has desired TON balance.
	gwBalance, err := deployer.GetBalanceOf(ctx, gwAccount.ID)
	if err != nil {
		return errors.Wrap(err, "unable to get balance of TON gateway")
	}

	if gwBalance.IsZero() {
		return fmt.Errorf("TON gateway balance is zero")
	}

	r.TONDeployer = deployer
	r.TONGateway = gwAccount.ID

	return nil
}
