package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/client"
	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
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
		Example: "zetacored testnet testnet_7001-1 zeta13c7p3xrhd6q2rx3h235jpt8pjdwvacyw6twpax",
		Args:    cobra.ExactArgs(2),
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

			skipConfirmation, err := cmd.Flags().GetBool(FlagSkipConfirmation)
			if err != nil {
				return fmt.Errorf("failed to get skip-confirmation flag: %w", err)
			}

			if !skipConfirmation {
				// Confirmation prompt to prevent accidental modification of state.
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
			withCmt, err := cmd.Flags().GetBool(srvflags.WithCometBFT)
			if err != nil {
				return fmt.Errorf("failed to get with-cometbft flag: %w", err)
			}

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
	cmd.Flags().Duration(server.FlagShutdownGrace, 3*time.Second, "On Shutdown, duration to wait for resource clean up")

	return cmd
}

func initAppForTestnet(svrCtx *server.Context, appInterface types.Application) error {
	app, ok := appInterface.(*zeta.App)
	if !ok {
		panic("expected *zeta.App")
	}
	err := updateObserverData(svrCtx, *app)
	if err != nil {
		return fmt.Errorf("failed to update observer data: %w", err)
	}
	err = updateValidatorData(svrCtx, *app)
	if err != nil {
		return fmt.Errorf("failed to update validator data: %w", err)
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
		return fmt.Errorf("failed to get validator consensus pubkey as bytes")
	}
	pubkey := &ed25519.PubKey{Key: newValPubkeyBytes}
	pubkeyAny, err := types2.NewAnyWithValue(pubkey)
	if err != nil {
		return fmt.Errorf("failed to pack pubkey into Any: %w", err)
	}

	newValAddrBytes, ok := svrCtx.Viper.Get(KeyValidatorConsensusAddr).([]byte)
	if !ok {
		return fmt.Errorf("failed to get validator consensus address as bytes")
	}
	newConsAddr := sdk.ConsAddress(newValAddrBytes)

	valAddress, err := observertypes.GetOperatorAddressFromAccAddress(operatorAddrStr)
	if err != nil {
		return fmt.Errorf("failed to get operator address from account address: %w", err)
	}

	tokens, ok := math.NewIntFromString(DefaultTestnetValidatorTokes)
	if !ok {
		return fmt.Errorf("failed to parse tokens string to Int")
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
		return fmt.Errorf("failed to get staking params: %w", err)
	}

	params.MaxValidators = 1
	params.UnbondingTime = 5 * time.Second
	err = app.StakingKeeper.SetParams(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to set staking params: %w", err)
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

	newValidatorSigningInfo := slashingtypes.ValidatorSigningInfo{
		Address:     newConsAddr.String(),
		StartHeight: app.LastBlockHeight() - 1,
		Tombstoned:  false,
	}
	err = app.SlashingKeeper.SetValidatorSigningInfo(ctx, newConsAddr, newValidatorSigningInfo)
	if err != nil {
		return fmt.Errorf("failed to set validator signing info: %w", err)
	}

	sp, err := app.SlashingKeeper.GetParams(ctx)
	if err != nil {
		return fmt.Errorf("failed to get slashing params: %w", err)
	}
	sp.MinSignedPerWindow = math.LegacyZeroDec()
	err = app.SlashingKeeper.SetParams(ctx, sp)
	if err != nil {
		return fmt.Errorf("failed to set slashing params: %w", err)
	}

	return nil
}
