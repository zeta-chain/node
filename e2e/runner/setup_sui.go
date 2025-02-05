package runner

import (
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
)

func (r *E2ERunner) SetupSui(faucetURL string) {
	r.Logger.Print("⚙️ initializing gateway program on Sui")

	deployerAddr, err := r.Account.SuiAddress()
	require.NoError(r, err, "get deploy address")

	header := map[string]string{}
	err = sui.RequestSuiFromFaucet(faucetURL, deployerAddr, header)
	require.NoError(r, err, "sui faucet request to %s", faucetURL)
}
