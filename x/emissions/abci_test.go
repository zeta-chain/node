package emissions_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	zetaapp "github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/simapp"
	emissionsModule "github.com/zeta-chain/zetacore/x/emissions"
	emissionsModuleTypes "github.com/zeta-chain/zetacore/x/emissions/types"
	"testing"
)

func getaZetaFromString(amount string) sdk.Coins {
	emissionPoolInt, _ := sdk.NewIntFromString(amount)
	return sdk.NewCoins(sdk.NewCoin(config.BaseDenom, emissionPoolInt))
}

func SetupApp(t *testing.T, params emissionsModuleTypes.Params, emissionPoolCoins sdk.Coins) (*zetaapp.App, sdk.Context, *tmtypes.ValidatorSet) {
	pk1 := ed25519.GenPrivKey().PubKey()
	acc1 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pk1.Address()))
	pk2 := ed25519.GenPrivKey().PubKey()
	acc2 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pk2.Address()))
	// genDelActs and genDelBalances need to have the same addresses
	// bondAmount is specified separately , the Balances here are additional tokens for delegators to have in their accounts
	genDelActs := make(authtypes.GenesisAccounts, 1)
	genDelBalances := make([]banktypes.Balance, 1)
	genDelActs[0] = acc1
	genDelBalances[0] = banktypes.Balance{
		Address: acc1.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(1000000))),
	}
	genDelActs[0] = acc2
	genDelBalances[0] = banktypes.Balance{
		Address: acc2.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(1000000))),
	}
	delBondAmount := getaZetaFromString("1000000000000000000000000")

	genBalances := make([]banktypes.Balance, 1)
	genBalances[0] = banktypes.Balance{
		Address: emissionsModuleTypes.EmissionsModuleAddress.String(),
		Coins:   emissionPoolCoins,
	}

	vset := tmtypes.NewValidatorSet([]*tmtypes.Validator{})
	for i := 0; i < 1; i++ {
		privKey := ed25519.GenPrivKey()
		pubKey := privKey.PubKey()
		val := tmtypes.NewValidator(pubKey, 100000000000)
		err := vset.UpdateWithChangeSet([]*tmtypes.Validator{val})
		if err != nil {
			panic("Failed to add validator")
		}
	}

	app := simapp.SetupWithGenesisValSet(t, vset, genDelActs, delBondAmount.AmountOf(config.BaseDenom), params, genDelBalances, genBalances)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	return app, ctx, vset
}

func TestAppModule_GetBlockRewardComponents(t *testing.T) {

	tests := []struct {
		name                             string
		startingEmissionPool             string
		testMaxHeight                    int64
		checkValuesReserver              map[int64]string
		checkValuesUndistributedObserver map[int64]string
		checkValuesUndistributedTss      map[int64]string
	}{
		{
			name:                 "test 1",
			startingEmissionPool: "1000000000000000000000000",
			testMaxHeight:        5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, ctx, _ := SetupApp(t, emissionsModuleTypes.DefaultParams(), getaZetaFromString(tt.startingEmissionPool))
			fmt.Println(app.BankKeeper.GetBalance(ctx, emissionsModuleTypes.EmissionsModuleAddress, config.BaseDenom))
			for i := ctx.BlockHeight(); i < tt.testMaxHeight; i++ {
				ctx = ctx.WithBlockHeight(i)
				//app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
				//	Height:             app.LastBlockHeight() + 1,
				//	AppHash:            app.LastCommitID().Hash,
				//	ValidatorsHash:     vset.Hash(),
				//	NextValidatorsHash: vset.Hash(),
				//	ChainID:            "simnet_101-1",
				//}})
				emissionsModule.BeginBlocker(ctx, app.EmissionsKeeper, app.StakingKeeper, app.BankKeeper)
				bondFactor, reservesFactor, durationFactor := emissionsModule.GetBlockRewardComponents(ctx, app.BankKeeper, app.StakingKeeper, app.EmissionsKeeper)
				fmt.Printf("bondFactor: %d, reservesFactor: %d, durationFactor: %d , blockheight: %d\n", bondFactor, reservesFactor, durationFactor, ctx.BlockHeight())
				//app.Commit()

			}
		})
	}

}
