package simapp

import (
	"encoding/json"
	"testing"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	types2 "github.com/zeta-chain/zetacore/x/emissions/types"

	//"github.com/ignite-hq/cli/ignite/pkg/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"
	"github.com/zeta-chain/zetacore/app"
)

func Setup(isCheckTx bool) *app.App {
	app, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				ChainId:         "simnet_101-1",
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: defaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func setup(withGenesis bool, invCheckPeriod uint) (*app.App, app.GenesisState) {
	db := tmdb.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	a := app.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, invCheckPeriod, encCdc, simapp.EmptyAppOptions{})
	if withGenesis {
		return a, app.NewDefaultGenesisState(encCdc.Codec)
	}
	return a, app.GenesisState{}
}

func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genDelAccs []authtypes.GenesisAccount, bondAmt sdk.Int, emissionParams types2.Params, genDelBalances []banktypes.Balance, genBalances []banktypes.Balance) *app.App {
	app, genesisState := setup(true, 5)
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genDelAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))
	// Make all members of valSet as validators
	// Make all members of delSet as delegators to each of the validators
	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		assert.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		assert.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genDelAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))
	}

	emissionsGenesis := types2.DefaultGenesis()
	emissionsGenesis.Params = emissionParams
	genesisState[types2.ModuleName] = app.AppCodec().MustMarshalJSON(emissionsGenesis)
	// set validators and delegations
	params := stakingtypes.DefaultParams()
	params.BondDenom = config.BaseDenom
	stakingGenesis := stakingtypes.NewGenesisState(params, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	// genDelBalances contains additional balances for delegators
	// Add Bond amount and additional coins for these addresses

	for _, b := range genDelBalances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(config.BaseDenom, bondAmt))...)
	}
	// add balances for non delegator accounts
	// Add only external balances
	for _, b := range genBalances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	totalBalances := []banktypes.Balance{}
	// Add extra balance to account for delegator bonded pool
	totalBalances = append(append(append(totalBalances, genBalances...), genDelBalances...), banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(config.BaseDenom, bondAmt)},
	})

	// update total supply

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, totalBalances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	assert.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         "simnet_101-1",
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: defaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()

	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
		ChainID:            "simnet_101-1",
	}})

	return app
}

func SetupWithGenesisAccounts(genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *app.App {
	app, genesisState := setup(true, 0)
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		totalSupply = totalSupply.Add(b.Coins...)
	}

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: defaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	app.Commit()

	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1}})

	return app
}
