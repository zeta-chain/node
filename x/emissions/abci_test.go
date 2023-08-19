package emissions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

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
	emissions "github.com/zeta-chain/zetacore/x/emissions"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
)

func getaZetaFromString(amount string) sdk.Coins {
	emissionPoolInt, _ := sdk.NewIntFromString(amount)
	return sdk.NewCoins(sdk.NewCoin(config.BaseDenom, emissionPoolInt))
}

func SetupApp(t *testing.T, params emissionstypes.Params, emissionPoolCoins sdk.Coins) (*zetaapp.App, sdk.Context, *tmtypes.ValidatorSet, *authtypes.BaseAccount) {
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
	//	Address: emissionstypes.EmissionsModuleAddress.String(),
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
		params               emissionstypes.Params
		testMaxHeight        int64
		inputFilename        string
		checkValues          []EmissionTestData
		generateOnly         bool
	}{
		{
			name:                 "default values",
			params:               emissionstypes.DefaultParams(),
			startingEmissionPool: "1000000000000000000000000",
			testMaxHeight:        300,
			inputFilename:        "simulations.json",
			generateOnly:         false,
		},
		{
			name:                 "higher starting pool",
			params:               emissionstypes.DefaultParams(),
			startingEmissionPool: "100000000000000000000000000000000",
			testMaxHeight:        300,
			inputFilename:        "simulations.json",
			generateOnly:         false,
		},
		{
			name:                 "lower starting pool",
			params:               emissionstypes.DefaultParams(),
			startingEmissionPool: "100000000000000000",
			testMaxHeight:        300,
			inputFilename:        "simulations.json",
			generateOnly:         false,
		},
		{
			name: "different distribution percentages",
			params: emissionstypes.Params{
				MaxBondFactor:               "1.25",
				MinBondFactor:               "0.75",
				AvgBlockTime:                "6.00",
				TargetBondRatio:             "00.67",
				ValidatorEmissionPercentage: "00.10",
				ObserverEmissionPercentage:  "00.85",
				TssSignerEmissionPercentage: "00.05",
				DurationFactorConstant:      "0.001877876953694702",
			},
			startingEmissionPool: "1000000000000000000000000",
			testMaxHeight:        300,
			inputFilename:        "simulations.json",
			generateOnly:         false,
		},
		{
			name: "higher block time",
			params: emissionstypes.Params{
				MaxBondFactor:               "1.25",
				MinBondFactor:               "0.75",
				AvgBlockTime:                "20.00",
				TargetBondRatio:             "00.67",
				ValidatorEmissionPercentage: "00.10",
				ObserverEmissionPercentage:  "00.85",
				TssSignerEmissionPercentage: "00.05",
				DurationFactorConstant:      "0.1",
			},
			startingEmissionPool: "1000000000000000000000000",
			testMaxHeight:        300,
			inputFilename:        "simulations.json",
			generateOnly:         false,
		},
		{
			name: "different duration constant",
			params: emissionstypes.Params{
				MaxBondFactor:               "1.25",
				MinBondFactor:               "0.75",
				AvgBlockTime:                "6.00",
				TargetBondRatio:             "00.67",
				ValidatorEmissionPercentage: "00.10",
				ObserverEmissionPercentage:  "00.85",
				TssSignerEmissionPercentage: "00.05",
				DurationFactorConstant:      "0.1",
			},
			startingEmissionPool: "1000000000000000000000000",
			testMaxHeight:        300,
			inputFilename:        "simulations.json",
			generateOnly:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, ctx, _, minter := SetupApp(t, tt.params, getaZetaFromString(tt.startingEmissionPool))
			err := app.BankKeeper.SendCoinsFromAccountToModule(ctx, minter.GetAddress(), emissionstypes.ModuleName, getaZetaFromString(tt.startingEmissionPool))
			assert.NoError(t, err)
			GenerateTestDataMaths(app, ctx, tt.testMaxHeight, tt.inputFilename)
			defer func(t *testing.T, fp string) {
				err := os.RemoveAll(fp)
				assert.NoError(t, err)
			}(t, tt.inputFilename)

			if tt.generateOnly {
				return
			}
			inputTestData, err := GetInputData(tt.inputFilename)
			assert.NoError(t, err)
			sort.SliceStable(inputTestData, func(i, j int) bool { return inputTestData[i].BlockHeight < inputTestData[j].BlockHeight })
			startHeight := ctx.BlockHeight()
			assert.Equal(t, startHeight, inputTestData[0].BlockHeight, "starting block height should be equal to the first block height in the input data")
			for i := startHeight; i < tt.testMaxHeight; i++ {
				//The First distribution will occur only when begin-block is triggered
				reservesFactor, bondFactor, durationFactor := app.EmissionsKeeper.GetBlockRewardComponents(ctx)
				assert.Equal(t, inputTestData[i-1].ReservesFactor, reservesFactor, "reserves factor should be equal to the input data"+fmt.Sprintf(" , block height: %d", i))
				assert.Equal(t, inputTestData[i-1].BondFactor, bondFactor, "bond factor should be equal to the input data"+fmt.Sprintf(" , block height: %d", i))
				assert.Equal(t, inputTestData[i-1].DurationFactor, durationFactor.String(), "duration factor should be equal to the input data"+fmt.Sprintf(" , block height: %d", i))
				emissions.BeginBlocker(ctx, app.EmissionsKeeper)
				ctx = ctx.WithBlockHeight(i + 1)
			}
		})
	}
}

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

func GenerateTestDataMaths(app *zetaapp.App, ctx sdk.Context, testMaxHeight int64, fileName string) {
	var generatedTestData []EmissionTestData
	reserverCoins := app.BankKeeper.GetBalance(ctx, emissionstypes.EmissionsModuleAddress, config.BaseDenom)
	startHeight := ctx.BlockHeight()
	for i := startHeight; i < testMaxHeight; i++ {
		reservesFactor := sdk.NewDecFromInt(reserverCoins.Amount)
		bondFactor := app.EmissionsKeeper.GetBondFactor(ctx, app.StakingKeeper)
		durationFactor := app.EmissionsKeeper.GetDurationFactor(ctx)
		blockRewards := reservesFactor.Mul(bondFactor).Mul(durationFactor)
		generatedTestData = append(generatedTestData, EmissionTestData{
			BlockHeight:    i,
			BondFactor:     bondFactor,
			DurationFactor: durationFactor.String(),
			ReservesFactor: reservesFactor,
		})
		validatorRewards := sdk.MustNewDecFromStr(app.EmissionsKeeper.GetParams(ctx).ValidatorEmissionPercentage).Mul(blockRewards).TruncateInt()
		observerRewards := sdk.MustNewDecFromStr(app.EmissionsKeeper.GetParams(ctx).ObserverEmissionPercentage).Mul(blockRewards).TruncateInt()
		tssSignerRewards := sdk.MustNewDecFromStr(app.EmissionsKeeper.GetParams(ctx).TssSignerEmissionPercentage).Mul(blockRewards).TruncateInt()
		truncatedRewards := validatorRewards.Add(observerRewards).Add(tssSignerRewards)
		reserverCoins = reserverCoins.Sub(sdk.NewCoin(config.BaseDenom, truncatedRewards))
		ctx = ctx.WithBlockHeight(i + 1)
	}
	GenerateSampleFile(fileName, generatedTestData)
}

func GenerateSampleFile(fp string, data []EmissionTestData) {
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(fp, file, 0600)
}
