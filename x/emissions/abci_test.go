package emissions_test

import (
	"fmt"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/simapp"
	"testing"
)

func TestAppModule_BeginBlock(t *testing.T) {
	app := simapp.Setup(false)
	tmtypes.NewValidatorSet()
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	fmt.Println(app.EmissionsKeeper.GetParams(ctx).String())
	fmt.Println(app.BankKeeper.GetSupply(ctx, config.BaseDenom))
}
