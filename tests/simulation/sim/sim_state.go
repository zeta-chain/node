package sim

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"cosmossdk.io/math"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/zeta-chain/ethermint/x/evm/types"
	"github.com/zeta-chain/node/testutil/sample"
	observertypes "github.com/zeta-chain/node/x/observer/types"

	zetaapp "github.com/zeta-chain/node/app"
)

// Simulation parameter constants
const (
	StakePerAccount           = "stake_per_account"
	InitiallyBondedValidators = "initially_bonded_validators"
)

// AppStateFn returns the initial application state using a genesis or the simulation parameters.
// It panics if the user provides files for both of them.
// If a file is not given for the genesis or the sim params, it creates a randomized one.
func AppStateFn(
	cdc codec.Codec,
	simManager *module.SimulationManager,
	genesisState map[string]json.RawMessage,
) simtypes.AppStateFn {
	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config,
	) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {
		if FlagGenesisTimeValue == 0 {
			genesisTimestamp = simtypes.RandTimestamp(r)
		} else {
			genesisTimestamp = time.Unix(FlagGenesisTimeValue, 0)
		}

		chainID = config.ChainID
		switch {
		case config.ParamsFile != "" && config.GenesisFile != "":
			panic("cannot provide both a genesis file and a params file")

		case config.GenesisFile != "":
			// override the default chain-id from simapp to set it later to the config
			genesisDoc, accounts, err := AppStateFromGenesisFileFn(r, cdc, config.GenesisFile)
			if err != nil {
				panic(err)
			}

			if FlagGenesisTimeValue == 0 {
				// use genesis timestamp if no custom timestamp is provided (i.e no random timestamp)
				genesisTimestamp = genesisDoc.GenesisTime
			}

			appState = genesisDoc.AppState
			chainID = genesisDoc.ChainID
			simAccs = accounts

		case config.ParamsFile != "":
			appParams := make(simtypes.AppParams)
			bz, err := os.ReadFile(config.ParamsFile)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal(bz, &appParams)
			if err != nil {
				panic(err)
			}
			appState, simAccs = AppStateRandomizedFn(
				simManager,
				r,
				cdc,
				accs,
				genesisTimestamp,
				appParams,
				genesisState,
			)

		default:
			appParams := make(simtypes.AppParams)
			appState, simAccs = AppStateRandomizedFn(
				simManager,
				r,
				cdc,
				accs,
				genesisTimestamp,
				appParams,
				genesisState,
			)
		}

		rawState := make(map[string]json.RawMessage)
		err := json.Unmarshal(appState, &rawState)
		if err != nil {
			panic(err)
		}

		stakingStateBz, ok := rawState[stakingtypes.ModuleName]
		if !ok {
			panic("staking genesis state is missing")
		}

		stakingState := new(stakingtypes.GenesisState)
		err = cdc.UnmarshalJSON(stakingStateBz, stakingState)
		if err != nil {
			panic(err)
		}

		// compute not bonded balance
		notBondedTokens := math.ZeroInt()
		for _, val := range stakingState.Validators {
			if val.Status != stakingtypes.Unbonded {
				continue
			}
			notBondedTokens = notBondedTokens.Add(val.GetTokens())
		}
		notBondedCoins := sdk.NewCoin(stakingState.Params.BondDenom, notBondedTokens)

		// edit bank state to make it have the not bonded pool tokens
		bankStateBz, ok := rawState[banktypes.ModuleName]
		if !ok {
			panic("bank genesis state is missing")
		}
		bankState := new(banktypes.GenesisState)
		err = cdc.UnmarshalJSON(bankStateBz, bankState)
		if err != nil {
			panic(err)
		}

		stakingAddr := authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String()
		var found bool
		for _, balance := range bankState.Balances {
			if balance.Address == stakingAddr {
				found = true
				break
			}
		}
		if !found {
			bankState.Balances = append(bankState.Balances, banktypes.Balance{
				Address: stakingAddr,
				Coins:   sdk.NewCoins(notBondedCoins),
			})
		}

		// Set the bond denom in the EVM genesis state
		evmStateBz, ok := rawState[evmtypes.ModuleName]
		if !ok {
			panic("evm genesis state is missing")
		}

		evmState := new(evmtypes.GenesisState)
		cdc.MustUnmarshalJSON(evmStateBz, evmState)

		// we should replace the EvmDenom with BondDenom
		evmState.Params.EvmDenom = stakingState.Params.BondDenom

		observers := []string{}
		for _, validator := range stakingState.Validators {
			accAddress, err := observertypes.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
			if err != nil {
				panic(err)
			}
			observers = append(observers, accAddress.String())
		}
		observers = observers[:5]

		// update the observer genesis state
		observerStateBz, ok := rawState[observertypes.ModuleName]
		if !ok {
			panic("observer genesis state is missing")
		}
		observerState := new(observertypes.GenesisState)
		cdc.MustUnmarshalJSON(observerStateBz, observerState)
		observerState.Observers.ObserverList = observers
		observerState.CrosschainFlags.IsInboundEnabled = true
		observerState.CrosschainFlags.IsOutboundEnabled = true
		tss := sample.Tss()
		observerState.Tss = &tss

		// change appState back
		rawState[evmtypes.ModuleName] = cdc.MustMarshalJSON(evmState)
		rawState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingState)
		rawState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)
		rawState[observertypes.ModuleName] = cdc.MustMarshalJSON(observerState)

		// replace appstate
		appState, err = json.Marshal(rawState)
		if err != nil {
			panic(err)
		}
		return appState, simAccs, chainID, genesisTimestamp
	}
}

// AppStateRandomizedFn creates calls each module's GenesisState generator function
// and creates the simulation params
func AppStateRandomizedFn(
	simManager *module.SimulationManager, r *rand.Rand, cdc codec.Codec,
	accs []simtypes.Account, genesisTimestamp time.Time, appParams simtypes.AppParams,
	genesisState map[string]json.RawMessage,
) (json.RawMessage, []simtypes.Account) {
	numAccs := int64(len(accs))
	// generate a random amount of initial stake coins and a random initial
	// number of bonded accounts
	var (
		numInitiallyBonded int64
		initialStake       math.Int
	)

	appParams.GetOrGenerate(cdc,
		StakePerAccount, &initialStake, r,
		func(r *rand.Rand) { initialStake = math.NewInt(r.Int63n(1e12)) },
	)
	appParams.GetOrGenerate(cdc,
		InitiallyBondedValidators, &numInitiallyBonded, r,
		func(r *rand.Rand) { numInitiallyBonded = int64(r.Intn(300)) },
	)

	if numInitiallyBonded > numAccs {
		numInitiallyBonded = numAccs
	}

	// Set the default power reduction to be one less than the initial stake so that all randomised validators are part of the validator set
	sdk.DefaultPowerReduction = initialStake.Sub(sdk.OneInt())

	fmt.Printf(
		`Selected randomly generated parameters for simulated genesis:
{
  stake_per_account: "%d",
  initially_bonded_validators: "%d"
}
`, initialStake, numInitiallyBonded,
	)

	simState := &module.SimulationState{
		AppParams:    appParams,
		Cdc:          cdc,
		Rand:         r,
		GenState:     genesisState,
		Accounts:     accs,
		InitialStake: initialStake,
		NumBonded:    numInitiallyBonded,
		GenTimestamp: genesisTimestamp,
	}

	simManager.GenerateGenesisStates(simState)

	appState, err := json.Marshal(genesisState)
	if err != nil {
		panic(err)
	}

	return appState, accs
}

// AppStateFromGenesisFileFn util function to generate the genesis AppState
// from a genesis.json file.
func AppStateFromGenesisFileFn(
	r io.Reader,
	cdc codec.JSONCodec,
	genesisFile string,
) (tmtypes.GenesisDoc, []simtypes.Account, error) {
	bytes, err := os.ReadFile(genesisFile) // #nosec G304 -- genesisFile value is controlled
	if err != nil {
		panic(err)
	}

	var genesis tmtypes.GenesisDoc
	// NOTE: Comet uses a custom JSON decoder for GenesisDoc
	err = cmtjson.Unmarshal(bytes, &genesis)
	if err != nil {
		panic(err)
	}

	var appState zetaapp.GenesisState
	err = json.Unmarshal(genesis.AppState, &appState)
	if err != nil {
		panic(err)
	}

	var authGenesis authtypes.GenesisState
	if appState[authtypes.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[authtypes.ModuleName], &authGenesis)
	}

	newAccs := make([]simtypes.Account, len(authGenesis.Accounts))
	for i, acc := range authGenesis.Accounts {
		// Pick a random private key, since we don't know the actual key
		// This should be fine as it's only used for mock Tendermint validators
		// and these keys are never actually used to sign by mock Tendermint.
		privkeySeed := make([]byte, 15)
		if _, err := r.Read(privkeySeed); err != nil {
			panic(err)
		}

		privKey := secp256k1.GenPrivKeyFromSecret(privkeySeed)

		a, ok := acc.GetCachedValue().(authtypes.AccountI)
		if !ok {
			return genesis, nil, fmt.Errorf("expected account")
		}

		// create simulator accounts
		simAcc := simtypes.Account{PrivKey: privKey, PubKey: privKey.PubKey(), Address: a.GetAddress()}
		newAccs[i] = simAcc
	}

	return genesis, newAccs, nil
}
