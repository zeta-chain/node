package rpcimportable

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/rpc"
	"github.com/zeta-chain/node/pkg/sdkconfig"
)

func TestRPCImportable(t *testing.T) {
	_ = rpc.Clients{}
}

func TestCosmosSdkConfigUntouched(t *testing.T) {
	zetaCfg := sdk.NewConfig()
	sdkconfig.Set(zetaCfg, true)
	if zetaCfg.GetBech32AccountAddrPrefix() != sdkconfig.AccountAddressPrefix {
		t.Logf("zetaCfg account prefix is not %s", sdkconfig.AccountAddressPrefix)
		t.FailNow()
	}

	// ensure that importing/using zeta sdkconfig does not mutate the global config
	globalConfig := sdk.GetConfig()
	if globalConfig.GetBech32AccountAddrPrefix() != sdk.Bech32MainPrefix {
		t.Logf("globalConfig account prefix is not %s", sdk.Bech32MainPrefix)
		t.FailNow()
	}
}
