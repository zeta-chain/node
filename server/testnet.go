package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client"
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
		Use:   "testnet [newChainID]",
		Short: "Modify state to create testnet from current local data",
		Long: `Modify state to create a testnet from current local state.
This command modifies the chain ID and validator set to allow the local validator
to control the network. After running this command, use the regular "start" command
to start the network.


WARNING: This operation modifies state in your data folder and cannot be undone.

Example usage:
  zetacored testnet testnet_7001-1
  zetacored start

The first block may take up to one minute to be committed, depending on how old
the state is. If using old snapshots, pending state (expiring locks, etc.) may
need to be committed first.
`,
		Example: "zetacored testnet testnet_7001-1",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			_, err := server.GetPruningOptionsFromFlags(serverCtx.Viper)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			newChainID := args[0]
			operatorAddress := args[1]

			// TODO add validation
			skipConfirmation, _ := cmd.Flags().GetBool(FlagSkipConfirmation)

			if !skipConfirmation {
				// Confirmation prompt to prevent accidental modification of state.
				reader := bufio.NewReader(os.Stdin)
				fmt.Println("This operation will modify state in your data folder and cannot be undone. Do you want to continue? (y/n)")
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

			return err
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
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

	observerSet, found := app.ObserverKeeper.GetObserverSet(ctx)
	if !found {
		panic("no observer set")
	}

	newObserverSet := observertypes.ObserverSet{
		ObserverList: []string{observerSet.ObserverList[0]},
	}
	app.ObserverKeeper.SetObserverSet(ctx, newObserverSet)

	err := updateStakingForTestnet(svrCtx, *app)
	if err != nil {
		panic(err)
	}
	return nil
}

func updateStakingForTestnet(svrCtx *server.Context, app zeta.App) error {
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	operatorAddrStr := svrCtx.Viper.GetString(KeyOperatorAddress)

	validators, err := app.StakingKeeper.GetAllValidators(ctx)
	if err != nil {
		panic(err)
	}

	newVal := stakingtypes.Validator{}

	for _, val := range validators {
		accAddr, err := observertypes.GetAccAddressFromOperatorAddress(val.OperatorAddress)
		if err != nil {
			panic(err)
		}
		if accAddr.String() == operatorAddrStr {
			newVal = val
		}
		val.Status = stakingtypes.Unbonded
		err = app.StakingKeeper.SetValidator(ctx, val)
		if err != nil {
			panic(err)
		}
	}

	params, err := app.StakingKeeper.GetParams(ctx)
	if err != nil {
		return err
	}

	params.MaxValidators = 1
	params.UnbondingTime = 5 * time.Second
	err = app.StakingKeeper.SetParams(ctx, params)
	if err != nil {
		return err
	}
	if newVal.OperatorAddress == "" {
		return fmt.Errorf("operator address %s not found in validator set", operatorAddrStr)
	}

	stakingKey := app.GetKey(stakingtypes.ModuleName)
	stakingStore := ctx.KVStore(stakingKey)
	iterator, err := app.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return err
	}

	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	svrCtx.Logger.Info("Cleared staking validators by power index")
	iterator, err = app.StakingKeeper.LastValidatorsIterator(ctx)
	if err != nil {
		return err
	}
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	iterator.Close()

	svrCtx.Logger.Info("Cleared staking last validator power")

	//
	err = app.StakingKeeper.SetValidator(ctx, newVal)
	if err != nil {
		return err
	}
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, newVal)
	if err != nil {
		return err
	}
	err = app.StakingKeeper.SetValidatorByPowerIndex(ctx, newVal)
	if err != nil {
		return err
	}

	valAddress, err := sdk.ValAddressFromBech32(newVal.GetOperator())
	if err != nil {
		return err
	}
	err = app.StakingKeeper.SetLastValidatorPower(ctx, valAddress, 0)
	if err != nil {
		return err
	}
	if err := app.StakingKeeper.Hooks().AfterValidatorCreated(ctx, valAddress); err != nil {
		return err
	}

	err = app.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddress, 0, distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1))
	if err != nil {
		return err
	}
	err = app.DistrKeeper.SetValidatorCurrentRewards(ctx, valAddress, distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1))
	if err != nil {
		return err
	}
	err = app.DistrKeeper.SetValidatorAccumulatedCommission(ctx, valAddress, distrtypes.InitialValidatorAccumulatedCommission())
	if err != nil {
		return err
	}
	err = app.DistrKeeper.SetValidatorOutstandingRewards(ctx, valAddress, distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}})
	if err != nil {
		return err
	}

	newValAddr := svrCtx.Viper.GetString(KeyNewValAddr)

	newConsAddr := sdk.ConsAddress(newValAddr)
	newValidatorSigningInfo := slashingtypes.ValidatorSigningInfo{
		Address:     newConsAddr.String(),
		StartHeight: app.LastBlockHeight() - 1,
		Tombstoned:  false,
	}
	err = app.SlashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newValidatorSigningInfo)
	if err != nil {
		return err
	}

	return nil
}
