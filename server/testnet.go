package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client"
	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	srvflags "github.com/cosmos/evm/server/flags"
	"github.com/spf13/cobra"

	zeta "github.com/zeta-chain/node/app"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestNetCmd(appCreator types.AppCreator) *cobra.Command {
	return TestnetCmdWithOptions(appCreator, StartCmdOptions{
		DBOpener:            openDB,
		StartCommandHandler: start,
	})
}

// TestnetCmdWithOptions creates a command that modifies the local state to create a testnet fork.
// After running this command, the network can be started with the regular start command.
func TestnetCmdWithOptions(testnetAppCreator types.AppCreator, opts StartCmdOptions) *cobra.Command {
	if opts.DBOpener == nil || opts.StartCommandHandler == nil {
		panic("DBOpener and StartCommandHandler must be provided")
	}

	cmd := &cobra.Command{
		Use:   "testnet [newChainID] [operatorAddress]",
		Short: "Modify state to create testnet from current local data",
		Long: `Modify state to create a testnet from current local state. This will set the chain ID to the provided newChainID.
The provided opeartorAddress is used as the operator for the single validator in this network. The existing node key is reused .
`,
		Example: "zetacored testnet testnet_7001-1 zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax",
		Args:    cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			_, err := server.GetPruningOptionsFromFlags(serverCtx.Viper)
			if err != nil {
				return fmt.Errorf("failed to get pruning options from flags: %w", err)
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return fmt.Errorf("failed to get client query context: %w", err)
			}

			newChainID := args[0]
			operatorAddress := args[1]

			skipConfirmation, _ := cmd.Flags().GetBool(FlagSkipConfirmation)

			if !skipConfirmation {
				// Confirmation prompt to prevent accidental modification of state.
				reader := bufio.NewReader(os.Stdin)
				fmt.Println(
					"This operation will modify state in your data folder and cannot be undone. Do you want to continue? (y/n)",
				)
				text, _ := reader.ReadString('\n')
				response := strings.TrimSpace(strings.ToLower(text))
				if response != "y" && response != "yes" {
					fmt.Println("Operation canceled.")
					return nil
				}
			}

			serverCtx.Viper.Set(KeyIsTestnet, true)
			serverCtx.Viper.Set(KeyNewChainID, newChainID)
			serverCtx.Viper.Set(KeyOperatorAddress, operatorAddress)
			withCmt, _ := cmd.Flags().GetBool(srvflags.WithCometBFT)

			err = opts.StartCommandHandler(serverCtx, clientCtx, testnetAppCreator, withCmt, opts)
			if err != nil {
				return fmt.Errorf("failed to start command handler: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().Bool(FlagSkipConfirmation, false, "Skip the confirmation prompt")
	cmd.Flags().Bool(srvflags.WithCometBFT, true, "Run abci app embedded in-process with CometBFT")
	cmd.Flags().String(srvflags.TraceStore, "", "Enable KVStore tracing to an output file")

	return cmd
}

func initAppForTestnet(svrCtx *server.Context, appInterface types.Application) error {
	app, ok := appInterface.(*zeta.App)
	if !ok {
		panic("expected *zeta.ZetaApp")
	}
	err := updateObserverData(*app)
	if err != nil {
		return fmt.Errorf("failed to update observer state: %w", err)
	}
	err = updateValidatorData(svrCtx, *app)
	if err != nil {
		return fmt.Errorf("failed to update staking for testnet: %w", err)
	}
	return nil
}

func updateObserverData(app zeta.App) error {
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

	observerSet, found := app.ObserverKeeper.GetObserverSet(ctx)
	if !found {
		return fmt.Errorf("could not find observer set")
	}

	newObserverSet := observertypes.ObserverSet{
		ObserverList: []string{observerSet.ObserverList[0]},
	}
	app.ObserverKeeper.SetObserverSet(ctx, newObserverSet)
	app.ObserverKeeper.SetLastObserverCount(ctx, &observertypes.LastObserverCount{
		Count:            1,
		LastChangeHeight: ctx.BlockHeight(),
	})
	return nil
}

func updateValidatorData(svrCtx *server.Context, app zeta.App) error {
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	operatorAddrStr := svrCtx.Viper.GetString(KeyOperatorAddress)

	validators, err := app.StakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all validators: %w", err)
	}

	newVal := stakingtypes.Validator{}

	for _, val := range validators {
		accAddr, err := observertypes.GetAccAddressFromOperatorAddress(val.OperatorAddress)
		if err != nil {
			return fmt.Errorf("failed to get account address from operator address: %w", err)
		}
		if accAddr.String() == operatorAddrStr {
			newVal = val
		}
		//val.Status = stakingtypes.Unbonded
		//err = app.StakingKeeper.SetValidator(ctx, val)
		//if err != nil {
		//	return fmt.Errorf("failed to set validator status to unbonded: %w", err)
		//}
	}

	//params, err := app.StakingKeeper.GetParams(ctx)
	//if err != nil {
	//	return fmt.Errorf("failed to get staking params: %w", err)
	//}
	//
	//params.MaxValidators = 1
	//params.UnbondingTime = 5 * time.Second
	//err = app.StakingKeeper.SetParams(ctx, params)
	//if err != nil {
	//	return fmt.Errorf("failed to set staking params: %w", err)
	//}
	if newVal.OperatorAddress == "" {
		return fmt.Errorf("operator address %s not found in validator set", operatorAddrStr)
	}

	stakingKey := app.GetKey(stakingtypes.ModuleName)
	stakingStore := ctx.KVStore(stakingKey)
	iterator, err := app.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return fmt.Errorf("failed to get validators power store iterator: %w", err)
	}

	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	if err := iterator.Close(); err != nil {
		return fmt.Errorf("failed to close validators power store iterator: %w", err)
	}

	svrCtx.Logger.Info("Cleared staking validators by power index")
	iterator, err = app.StakingKeeper.LastValidatorsIterator(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last validators iterator: %w", err)
	}
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	if err := iterator.Close(); err != nil {
		return fmt.Errorf("failed to close last validators iterator: %w", err)
	}

	svrCtx.Logger.Info("Cleared staking last validator power")

	//
	newValPubkeyBytes := svrCtx.Viper.Get(KeyValidatorPubkey).([]byte)
	pubkey := &ed25519.PubKey{Key: newValPubkeyBytes}
	pubkeyAny, err := types2.NewAnyWithValue(pubkey)
	if err != nil {
		return fmt.Errorf("failed to pack pubkey into Any: %w", err)
	}
	fmt.Println("Setting new validator pubkey:", pubkeyAny.String())
	newVal.ConsensusPubkey = pubkeyAny

	pk, ok := newVal.ConsensusPubkey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return fmt.Errorf("failed to get cached pubkey value")
	}

	fmt.Println("Setting new validator pubkey:", pk.String())
	err = app.StakingKeeper.SetValidator(ctx, newVal)
	if err != nil {
		return fmt.Errorf("failed to set validator: %w", err)
	}
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, newVal)
	if err != nil {
		return fmt.Errorf("failed to set validator by consensus address: %w", err)
	}
	err = app.StakingKeeper.SetValidatorByPowerIndex(ctx, newVal)
	if err != nil {
		return fmt.Errorf("failed to set validator by power index: %w", err)
	}

	valAddress, err := sdk.ValAddressFromBech32(newVal.GetOperator())
	if err != nil {
		return fmt.Errorf("failed to parse validator address from bech32: %w", err)
	}
	err = app.StakingKeeper.SetLastValidatorPower(ctx, valAddress, 0)
	if err != nil {
		return fmt.Errorf("failed to set last validator power: %w", err)
	}
	if err := app.StakingKeeper.Hooks().AfterValidatorCreated(ctx, valAddress); err != nil {
		return fmt.Errorf("failed to execute after validator created hook: %w", err)
	}

	err = app.DistrKeeper.SetValidatorHistoricalRewards(
		ctx,
		valAddress,
		0,
		distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1),
	)
	if err != nil {
		return fmt.Errorf("failed to set validator historical rewards: %w", err)
	}
	err = app.DistrKeeper.SetValidatorCurrentRewards(
		ctx,
		valAddress,
		distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1),
	)
	if err != nil {
		return fmt.Errorf("failed to set validator current rewards: %w", err)
	}
	err = app.DistrKeeper.SetValidatorAccumulatedCommission(
		ctx,
		valAddress,
		distrtypes.InitialValidatorAccumulatedCommission(),
	)
	if err != nil {
		return fmt.Errorf("failed to set validator accumulated commission: %w", err)
	}
	err = app.DistrKeeper.SetValidatorOutstandingRewards(
		ctx,
		valAddress,
		distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}},
	)
	if err != nil {
		return fmt.Errorf("failed to set validator outstanding rewards: %w", err)
	}

	newValAddr := svrCtx.Viper.GetString(KeyValidatorAddr)

	newConsAddr := sdk.ConsAddress(newValAddr)
	newValidatorSigningInfo := slashingtypes.ValidatorSigningInfo{
		Address:     newConsAddr.String(),
		StartHeight: app.LastBlockHeight() - 1,
		Tombstoned:  false,
	}
	err = app.SlashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newValidatorSigningInfo)
	if err != nil {
		return fmt.Errorf("failed to set validator signing info: %w", err)
	}

	params, err := app.SlashingKeeper.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get slashing params: %w", err)
	}
	params.MinSignedPerWindow = math.LegacyZeroDec()
	err = app.SlashingKeeper.SetParams(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to set slashing params: %w", err)
	}

	fmt.Println("slashing info set for validator:", newConsAddr.String())

	return nil
}
