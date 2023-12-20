package smoketests

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

func TestWhitelistERC20(sm *runner.SmokeTestRunner) {
	iter, err := sm.ERC20Custody.FilterWhitelisted(&bind.FilterOpts{
		Start:   0,
		End:     nil,
		Context: context.Background(),
	}, []ethcommon.Address{})
	if err != nil {
		panic(err)
	}
	for iter.Next() {
		sm.Logger.Info("whitelisted: %s", iter.Event.Asset.Hex())
	}
}
