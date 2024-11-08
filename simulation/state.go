package simulation

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
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"

	zetaapp "github.com/zeta-chain/node/app"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// Simulation parameter constants
const (
	StakePerAccount           = "stake_per_account"
	InitiallyBondedValidators = "initially_bonded_validators"
)

// AppStateFn returns the initial application state using a genesis or the simulation parameters.
// It panics if the user provides files for both of them.
// If a file is not given for the genesis or the sim params, it creates a randomized one.
// All modifications to the genesis state should be done in this function.
func AppStateFn(
	cdc codec.Codec,
	simManager *module.SimulationManager,
	genesisState map[string]json.RawMessage,
	exportedState json.RawMessage,
) simtypes.AppStateFn {
	return func(r *rand.Rand, accs []simtypes.Account, config simtypes.Config,
	) (appState json.RawMessage, simAccs []simtypes.Account, chainID string, genesisTimestamp time.Time) {
		if FlagGenesisTimeValue == 0 {
			genesisTimestamp = simtypes.RandTimestamp(r)
		} else {
			genesisTimestamp = time.Unix(FlagGenesisTimeValue, 0)
		}

		chainID = config.ChainID

		// If exported state is provided then use it
		if exportedState != nil {
			return exportedState, accs, chainID, genesisTimestamp
		}

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

		observers := make([]string, 0)
		// Get all the operator addresses of the validators.
		// The observer set can be a subset of the validator set
		for _, validator := range stakingState.Validators {
			accAddress, err := observertypes.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
			if err != nil {
				continue
			}
			observers = append(observers, accAddress.String())
		}

		// Shuffle the observers list
		r.Shuffle(len(observers), func(i, j int) {
			observers[i], observers[j] = observers[j], observers[i]
		})

		// Pick a random number of observers to add to the observer set
		numObservers := r.Intn(11) + 5
		if numObservers > len(observers) {
			numObservers = len(observers)
		}
		observers = observers[:numObservers]

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
		tss := observertypes.TSS{
			TssPubkey:           "cosmospub1addwnpepq27ldhn924mtwylm2r0vja3fcv3nv6gme0e2jnr96l0fkkqw6guscgqfsk0",
			KeyGenZetaHeight:    100,
			FinalizedZetaHeight: 110,
			TssParticipantList:  []string{},
			OperatorAddressList: observers,
		}
		observerState.Tss = &tss

		// Pick a random account to be the admin of all policies
		randomAccount := accs[r.Intn(len(accs))]
		authorityStateBz, ok := rawState[authoritytypes.ModuleName]
		if !ok {
			panic("authority genesis state is missing")
		}

		// update the authority genesis state
		authorityState := new(authoritytypes.GenesisState)
		cdc.MustUnmarshalJSON(authorityStateBz, authorityState)
		policies := authoritytypes.Policies{
			Items: []*authoritytypes.Policy{
				{
					Address:    randomAccount.Address.String(),
					PolicyType: authoritytypes.PolicyType_groupEmergency,
				},
				{
					Address:    randomAccount.Address.String(),
					PolicyType: authoritytypes.PolicyType_groupAdmin,
				},
				{
					Address:    randomAccount.Address.String(),
					PolicyType: authoritytypes.PolicyType_groupOperational,
				},
			},
		}
		authorityState.Policies = policies

		//Update the fungible genesis state
		fungibleStateBz, ok := rawState[fungibletypes.ModuleName]
		if !ok {
			panic("fungible genesis state is missing")
		}
		fungibleState := new(fungibletypes.GenesisState)
		cdc.MustUnmarshalJSON(fungibleStateBz, fungibleState)
		// TOODO generate ethereum address from r
		fungibleState.SystemContract = &fungibletypes.SystemContract{
			SystemContract: sample.EthAddressRandom(r).String(),
			ConnectorZevm:  sample.EthAddressRandom(r).String(),
			Gateway:        sample.EthAddressRandom(r).String(),
		}

		// change appState back
		rawState[evmtypes.ModuleName] = cdc.MustMarshalJSON(evmState)
		rawState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(stakingState)
		rawState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)
		rawState[observertypes.ModuleName] = cdc.MustMarshalJSON(observerState)
		rawState[authoritytypes.ModuleName] = cdc.MustMarshalJSON(authorityState)
		rawState[fungibletypes.ModuleName] = cdc.MustMarshalJSON(fungibleState)

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
