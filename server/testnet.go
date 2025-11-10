package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	srvflags "github.com/cosmos/evm/server/flags"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	zeta "github.com/zeta-chain/node/app"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	DefaultTestnetValidatorTokes = "30000000000000000000000"
	DefaultDelegatorShares       = "30000000000000000000000.000000000000000"
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
The provided operatorAddress is used as the operator for the single validator in this network. The existing node key is reused.
`,
		Example: `  zetacored testnet testnet_7001-1 zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax
  					zetacored testnet testnet_7001-1 zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax --upgrade-version v37.0.0`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			_, err := server.GetPruningOptionsFromFlags(serverCtx.Viper)
			if err != nil {
				return errors.Wrap(err, "failed to get pruning options from flags")
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return errors.Wrap(err, "failed to get client query context")
			}

			newChainID := args[0]
			operatorAddress := args[1]

			_, err = sdk.AccAddressFromBech32(operatorAddress)
			if err != nil {
				return errors.Wrap(err, "invalid operator address")
			}

			skipConfirmation, err := cmd.Flags().GetBool(FlagSkipConfirmation)
			if err != nil {
				return errors.Wrap(err, "failed to get skip-confirmation flag")
			}

			if !skipConfirmation {
				reader := bufio.NewReader(os.Stdin)
				fmt.Println(
					"This operation will modify state in your data folder and cannot be undone. This operation also updates the configuration , so it would not work with read only file systems. Do you want to continue? (y/n)",
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

			upgradeVersion, err := cmd.Flags().GetString(FlagUpgradeVersion)
			if err != nil {
				return fmt.Errorf("failed to get upgrade-version flag: %w", err)
			}
			if upgradeVersion != "" {
				serverCtx.Viper.Set(KeyUpgradeVersion, upgradeVersion)
			}

			withCmt, err := cmd.Flags().GetBool(srvflags.WithCometBFT)
			if err != nil {
				return errors.Wrap(err, "failed to get with-cometbft flag")
			}

			err = opts.StartCommandHandler(serverCtx, clientCtx, testnetAppCreator, withCmt, opts)
			if err != nil {
				return errors.Wrap(err, "failed to start command handler")
			}

			return nil
		},
	}

	cmd.Flags().Bool(FlagSkipConfirmation, false, "Skip the confirmation prompt")
	cmd.Flags().String(FlagUpgradeVersion, "", "Schedule upgrade to this version (e.g., v37.0.0). If empty, no upgrade is scheduled")
	cmd.Flags().Bool(srvflags.WithCometBFT, true, "Run abci app embedded in-process with CometBFT")
	cmd.Flags().String(srvflags.TraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().Duration(server.FlagShutdownGrace, 3*time.Second, "On Shutdown, duration to wait for resource clean up")

	return cmd
}

func initAppForTestnet(svrCtx *server.Context, appInterface types.Application) error {
	app, ok := appInterface.(*zeta.App)
	if !ok {
		return fmt.Errorf("invalid app type: %T", appInterface)
	}
	err := updateObserverData(svrCtx, *app)
	if err != nil {
		return errors.Wrap(err, "failed to update observer data")
	}
	err = updateValidatorData(svrCtx, *app)
	if err != nil {
		return errors.Wrap(err, "failed to update validator data")
	}
	err = updateUpgradeData(svrCtx, *app)
	if err != nil {
		return fmt.Errorf("failed to update upgrade data: %w", err)
	}
	return nil
}

// updateUpgradeData schedules an upgrade if the --upgrade-version flag is provided.
// It detects the current OS/architecture and creates download URLs for both zetacored and zetaclientd binaries.
func updateUpgradeData(svrCtx *server.Context, app zeta.App) error {
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})

	// Check if upgrade version flag was provided
	upgradeVersion := svrCtx.Viper.GetString(KeyUpgradeVersion)
	if upgradeVersion == "" {
		svrCtx.Logger.Info("No upgrade version specified, skipping upgrade scheduling")
		return nil
	}

	// Clear any existing upgrade plan from the source testnet
	existingPlan, err := app.UpgradeKeeper.GetUpgradePlan(ctx)
	if err == nil && existingPlan.Name != "" {
		svrCtx.Logger.Info("Clearing existing upgrade plan", "name", existingPlan.Name, "height", existingPlan.Height)
		if err := app.UpgradeKeeper.ClearUpgradePlan(ctx); err != nil {
			return fmt.Errorf("failed to clear existing upgrade plan: %w", err)
		}
	}
	appBlockHeight := svrCtx.Viper.GetInt64(KeyAppBlockedHeight)
	upgradeHeight := appBlockHeight + 100

	goos := runtime.GOOS
	goarch := runtime.GOARCH
	platform := fmt.Sprintf("%s/%s", goos, goarch)

	downloadInfo := map[string]interface{}{
		"binaries": map[string]string{
			platform: fmt.Sprintf(
				"https://github.com/zeta-chain/node/releases/download/%s/zetacored-ubuntu-22-%s",
				upgradeVersion,
				goarch,
			),
			fmt.Sprintf("zetaclientd-%s", platform): fmt.Sprintf(
				"https://github.com/zeta-chain/node/releases/download/%s/zetaclientd-ubuntu-22-%s",
				upgradeVersion,
				goarch,
			),
		},
	}

	infoBytes, err := json.Marshal(downloadInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal upgrade info: %w", err)
	}

	svrCtx.Logger.Info(
		"Scheduling upgrade",
		"version", upgradeVersion,
		"height", upgradeHeight,
		"platform", platform,
		"info", string(infoBytes),
	)

	err = app.UpgradeKeeper.ScheduleUpgrade(ctx, upgradetypes.Plan{
		Name:   upgradeVersion,
		Info:   string(infoBytes),
		Height: upgradeHeight,
	})
	if err != nil {
		return fmt.Errorf("failed to schedule upgrade: %w", err)
	}

	return nil
}

// updateObserverData updates the observer state to have a single observer: the operator address.
func updateObserverData(svrCtx *server.Context, app zeta.App) error {
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	operatorAddrStr := svrCtx.Viper.GetString(KeyOperatorAddress)
	newObserverSet := observertypes.ObserverSet{
		ObserverList: []string{operatorAddrStr},
	}
	app.ObserverKeeper.SetObserverSet(ctx, newObserverSet)
	app.ObserverKeeper.SetLastObserverCount(ctx, &observertypes.LastObserverCount{
		Count:            1,
		LastChangeHeight: ctx.BlockHeight(),
	})
	return nil
}

// updateValidatorData updates application state to have a single validator with the provided operator address and consensus pubkey.
// this affects staking, slashing, and distribution modules.
func updateValidatorData(svrCtx *server.Context, app zeta.App) error {
	ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
	operatorAddrStr := svrCtx.Viper.GetString(KeyOperatorAddress)

	newValPubkeyBytes, ok := svrCtx.Viper.Get(KeyValidatorConsensusPubkey).([]byte)
	if !ok {
		return errors.New("failed to get validator consensus pubkey as bytes")
	}
	pubkey := &ed25519.PubKey{Key: newValPubkeyBytes}
	pubkeyAny, err := codectypes.NewAnyWithValue(pubkey)
	if err != nil {
		return errors.Wrap(err, "failed to pack pubkey into Any")
	}

	newValAddrBytes, ok := svrCtx.Viper.Get(KeyValidatorConsensusAddr).([]byte)
	if !ok {
		return errors.New("failed to get validator consensus address as bytes")
	}
	newConsAddr := sdk.ConsAddress(newValAddrBytes)

	valAddress, err := observertypes.GetOperatorAddressFromAccAddress(operatorAddrStr)
	if err != nil {
		return errors.Wrap(err, "failed to get operator address from account address")
	}

	tokens, ok := math.NewIntFromString(DefaultTestnetValidatorTokes)
	if !ok {
		return errors.New("failed to parse tokens string to Int")
	}

	newVal := stakingtypes.Validator{
		OperatorAddress: valAddress.String(),
		ConsensusPubkey: pubkeyAny,
		Jailed:          false,
		Status:          stakingtypes.Bonded,
		Tokens:          tokens,
		DelegatorShares: math.LegacyMustNewDecFromStr(DefaultDelegatorShares),
		Description: stakingtypes.Description{
			Moniker: "Testnet Validator",
		},
		Commission: stakingtypes.Commission{
			CommissionRates: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.010000000000000000"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.200000000000000000"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.100000000000000000"),
			},
		},
		MinSelfDelegation: math.OneInt(),
	}

	params, err := app.StakingKeeper.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get staking params")
	}

	params.MaxValidators = 1
	params.UnbondingTime = 5 * time.Second
	err = app.StakingKeeper.SetParams(ctx, params)
	if err != nil {
		return errors.Wrap(err, "failed to set staking params")
	}

	stakingKey := app.GetKey(stakingtypes.ModuleName)
	stakingStore := ctx.KVStore(stakingKey)
	iterator, err := app.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get validators power store iterator")
	}

	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	if err := iterator.Close(); err != nil {
		return errors.Wrap(err, "failed to close validators power store iterator")
	}

	svrCtx.Logger.Info("Cleared staking validators by power index")
	iterator, err = app.StakingKeeper.LastValidatorsIterator(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get last validators iterator")
	}
	for ; iterator.Valid(); iterator.Next() {
		stakingStore.Delete(iterator.Key())
	}
	if err := iterator.Close(); err != nil {
		return errors.Wrap(err, "failed to close last validators iterator")
	}

	svrCtx.Logger.Info("Cleared staking last validator power")

	err = app.StakingKeeper.SetValidator(ctx, newVal)
	if err != nil {
		return errors.Wrap(err, "failed to set validator")
	}
	err = app.StakingKeeper.SetValidatorByConsAddr(ctx, newVal)
	if err != nil {
		return errors.Wrap(err, "failed to set validator by consensus address")
	}
	err = app.StakingKeeper.SetValidatorByPowerIndex(ctx, newVal)
	if err != nil {
		return errors.Wrap(err, "failed to set validator by power index")
	}

	err = app.StakingKeeper.SetLastValidatorPower(ctx, valAddress, 0)
	if err != nil {
		return errors.Wrap(err, "failed to set last validator power")
	}
	if err := app.StakingKeeper.Hooks().AfterValidatorCreated(ctx, valAddress); err != nil {
		return errors.Wrap(err, "failed to execute after validator created hook")
	}

	err = app.DistrKeeper.SetValidatorHistoricalRewards(
		ctx,
		valAddress,
		0,
		distrtypes.NewValidatorHistoricalRewards(sdk.DecCoins{}, 1),
	)
	if err != nil {
		return errors.Wrap(err, "failed to set validator historical rewards")
	}
	err = app.DistrKeeper.SetValidatorCurrentRewards(
		ctx,
		valAddress,
		distrtypes.NewValidatorCurrentRewards(sdk.DecCoins{}, 1),
	)
	if err != nil {
		return errors.Wrap(err, "failed to set validator current rewards")
	}
	err = app.DistrKeeper.SetValidatorAccumulatedCommission(
		ctx,
		valAddress,
		distrtypes.InitialValidatorAccumulatedCommission(),
	)
	if err != nil {
		return errors.Wrap(err, "failed to set validator accumulated commission")
	}
	err = app.DistrKeeper.SetValidatorOutstandingRewards(
		ctx,
		valAddress,
		distrtypes.ValidatorOutstandingRewards{Rewards: sdk.DecCoins{}},
	)
	if err != nil {
		return errors.Wrap(err, "failed to set validator outstanding rewards")
	}

	newValidatorSigningInfo := slashingtypes.ValidatorSigningInfo{
		Address:     newConsAddr.String(),
		StartHeight: app.LastBlockHeight() - 1,
		Tombstoned:  false,
	}
	err = app.SlashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newValidatorSigningInfo)
	if err != nil {
		return errors.Wrap(err, "failed to set validator signing info")
	}

	sp, err := app.SlashingKeeper.GetParams(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get slashing params")
	}
	sp.MinSignedPerWindow = math.LegacyZeroDec()
	err = app.SlashingKeeper.SetParams(ctx, sp)
	if err != nil {
		return errors.Wrap(err, "failed to set slashing params")
	}

	return nil
}
