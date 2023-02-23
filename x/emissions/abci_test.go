package emissions_test

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	zetaapp "github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/testutil/simapp"
	emissionsModule "github.com/zeta-chain/zetacore/x/emissions"
	emissionsModuleTypes "github.com/zeta-chain/zetacore/x/emissions/types"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"testing"
)

func getaZetaFromString(amount string) sdk.Coins {
	emissionPoolInt, _ := sdk.NewIntFromString(amount)
	return sdk.NewCoins(sdk.NewCoin(config.BaseDenom, emissionPoolInt))
}

func SetupApp(t *testing.T, params emissionsModuleTypes.Params, emissionPoolCoins sdk.Coins) (*zetaapp.App, sdk.Context, *tmtypes.ValidatorSet, *authtypes.BaseAccount) {
	pk1 := ed25519.GenPrivKey().PubKey()
	acc1 := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pk1.Address()))
	// genDelActs and genDelBalances need to have the same addresses
	// bondAmount is specified separately , the Balances here are additional tokens for delegators to have in their accounts
	genDelActs := make(authtypes.GenesisAccounts, 1)
	genDelBalances := make([]banktypes.Balance, 1)
	genDelActs[0] = acc1
	genDelBalances[0] = banktypes.Balance{
		Address: acc1.GetAddress().String(),
		Coins:   emissionPoolCoins,
	}
	delBondAmount := getaZetaFromString("1000000000000000000000000")

	//genBalances := make([]banktypes.Balance, 1)
	//genBalances[0] = banktypes.Balance{
	//	Address: emissionsModuleTypes.EmissionsModuleAddress.String(),
	//	Coins:   emissionPoolCoins,
	//}

	vset := tmtypes.NewValidatorSet([]*tmtypes.Validator{})
	for i := 0; i < 1; i++ {
		privKey := ed25519.GenPrivKey()
		pubKey := privKey.PubKey()
		val := tmtypes.NewValidator(pubKey, 1)
		err := vset.UpdateWithChangeSet([]*tmtypes.Validator{val})
		if err != nil {
			panic("Failed to add validator")
		}
	}

	app := simapp.SetupWithGenesisValSet(t, vset, genDelActs, delBondAmount.AmountOf(config.BaseDenom), params, genDelBalances, nil)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	ctx = ctx.WithBlockHeight(app.LastBlockHeight())
	return app, ctx, vset, acc1
}

type EmissionTestData struct {
	BlockHeight    int64   `json:"blockHeight,omitempty"`
	BondFactor     sdk.Dec `json:"bondFactor"`
	ReservesFactor sdk.Dec `json:"reservesFactor"`
	DurationFactor string  `json:"durationFactor"`
}

func TestAppModule_GetBlockRewardComponents(t *testing.T) {

	tests := []struct {
		name                 string
		startingEmissionPool string
		params               emissionsModuleTypes.Params
		testMaxHeight        int64
		inputFilename        string
		checkValues          []EmissionTestData
	}{
		{
			name:                 "test 1",
			params:               emissionsModuleTypes.DefaultParams(),
			startingEmissionPool: "1000000000000000000000000",
			testMaxHeight:        300,
			inputFilename:        "simulations.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, ctx, _, minter := SetupApp(t, tt.params, getaZetaFromString(tt.startingEmissionPool))
			err := app.BankKeeper.SendCoinsFromAccountToModule(ctx, minter.GetAddress(), emissionsModuleTypes.ModuleName, getaZetaFromString(tt.startingEmissionPool))
			assert.NoError(t, err)
			inputTestData, err := GetInputData(tt.inputFilename)
			assert.NoError(t, err)
			sort.SliceStable(inputTestData, func(i, j int) bool { return inputTestData[i].BlockHeight < inputTestData[j].BlockHeight })
			startHeight := ctx.BlockHeight()
			assert.Equal(t, startHeight, inputTestData[0].BlockHeight, "starting block height should be equal to the first block height in the input data")
			for i := startHeight; i < tt.testMaxHeight; i++ {
				//First distribution will occur only when begin-block is triggered
				reservesFactor, bondFactor, durationFactor := emissionsModule.GetBlockRewardComponents(ctx, app.BankKeeper, app.StakingKeeper, app.EmissionsKeeper)
				//generatedTestData = append(generatedTestData, EmissionTestData{
				//	BlockHeight:    i,
				//	BondFactor:     bondFactor,
				//	ReservesFactor: reservesFactor,
				//	DurationFactor: durationFactor,
				//})
				assert.Equal(t, inputTestData[i-1].ReservesFactor, reservesFactor, "reserves factor should be equal to the input data"+fmt.Sprintf(" , block height: %d", i))
				assert.Equal(t, inputTestData[i-1].BondFactor, bondFactor, "bond factor should be equal to the input data"+fmt.Sprintf(" , block height: %d", i))
				assert.Equal(t, inputTestData[i-1].DurationFactor, durationFactor.String(), "duration factor should be equal to the input data"+fmt.Sprintf(" , block height: %d", i))
				emissionsModule.BeginBlocker(ctx, app.EmissionsKeeper, app.StakingKeeper, app.BankKeeper)
				ctx = ctx.WithBlockHeight(i + 1)
			}
			//GenerateSampleFile("simulations.json", generatedTestData)
		})
	}
}

//fmt.Printf("Params:\n %+v \n", tt.params)
//fmt.Printf("Bond Ratio: %+v \n", app.StakingKeeper.BondedRatio(ctx))
//fmt.Printf("Total Bonded %s: \n", app.StakingKeeper.TotalBondedTokens(ctx))
//fmt.Printf("Emission Pool starting Balance : %s \n", app.BankKeeper.GetBalance(ctx, emissionsModuleTypes.EmissionsModuleAddress, config.BaseDenom))
//fmt.Printf("Total Supply : %s \n", app.BankKeeper.GetSupply(ctx, config.BaseDenom))

func GetInputData(fp string) ([]EmissionTestData, error) {
	data := []EmissionTestData{}
	file, err := filepath.Abs(fp)
	if err != nil {

		return nil, err
	}
	file = filepath.Clean(file)
	input, err := ioutil.ReadFile(file) // #nosec G304
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(input, &data)
	if err != nil {
		return nil, err
	}
	formatedData := make([]EmissionTestData, len(data))
	for i, dd := range data {
		fl, err := strconv.ParseFloat(dd.DurationFactor, 64)
		if err != nil {
			return nil, err
		}
		dd.DurationFactor = fmt.Sprintf("%0.18f", fl)
		formatedData[i] = dd
	}
	return formatedData, nil
}

func GenerateSampleFile(fp string, data []EmissionTestData) {
	file, _ := json.MarshalIndent(data, "", " ")
	//for _, dd := range data {
	//	fmt.Println(dd)
	//}
	_ = ioutil.WriteFile(fp, file, 0600)
}
