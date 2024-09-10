package runner

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/zeta-chain/node/e2e/runner/ton"
)

// SetupTON setups TON deployer and deploys Gateway contract
func (r *E2ERunner) SetupTON() error {
	if r.Clients.TON == nil || r.Clients.TONSidecar == nil {
		return fmt.Errorf("TON clients are not initialized")
	}

	// 1. Setup Deployer (acts as a faucet as well)
	faucetConfig, err := r.Clients.TONSidecar.GetFaucet(r.Ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get faucet config")
	}

	deployer, err := ton.NewDeployer(r.Clients.TON, faucetConfig)
	if err != nil {
		return errors.Wrap(err, "unable to create TON deployer")
	}

	r.TONDeployer = deployer

	gwCode, gwState, err := ton.GetGatewayCodeAndState(r.TSSAddress)
	if err != nil {
		return errors.Wrap(err, "unable to get TON Gateway code and state")
	}

	gw, err := r.TONDeployer.Deploy(r.Ctx, gwCode, gwState)
	if err != nil {
		return errors.Wrap(err, "unable to deploy TON Gateway")
	}

	fmt.Println("TON Gateway deployed", gw)

	return nil
}
