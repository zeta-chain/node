package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/flags"
	pruningtypes "github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	zetaCoreModuleTypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionsModuleTypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibleModuleTypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverModuleTypes "github.com/zeta-chain/zetacore/x/observer/types"

	"github.com/spf13/cobra"
)

const FlagAppDBBackend = "app-db-backend"

// PruningCmd prunes the sdk root multi store history versions based on the pruning options
// specified by command flags.
func RepairStoreCmd(appCreator servertypes.AppCreator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair-store",
		Short: "Prune app history states by keeping the recent heights and deleting old heights",
		Long: `Prune app history states by keeping the recent heights and deleting old heights.
		The pruning option is provided via the '--pruning' flag or alternatively with '--pruning-keep-recent'
		
		For '--pruning' the options are as follows:
		
		default: the last 362880 states are kept
		nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
		everything: 2 latest states will be kept
		custom: allow pruning options to be manually specified through 'pruning-keep-recent'.
		besides pruning options, database home directory and database backend type should also be specified via flags
		'--home' and '--app-db-backend'.
		valid app-db-backend type includes 'goleveldb', 'cleveldb', 'rocksdb', 'boltdb', and 'badgerdb'.
		`,
		Example: "prune --home './' --app-db-backend 'goleveldb' --pruning 'custom' --pruning-keep-recent 100",
		RunE: func(cmd *cobra.Command, _ []string) error {
			vp := viper.New()
			if err := vp.BindPFlags(cmd.Flags()); err != nil {
				return err
			}

			home := vp.GetString(flags.FlagHome)
			db, err := openDB(home, server.GetAppDBBackend(vp))
			if err != nil {
				return err
			}

			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			app := appCreator(logger, db, nil, vp)
			cms := app.CommitMultiStore()

			rms, ok := cms.(*rootmulti.Store)
			if !ok {
				return fmt.Errorf("currently only support the pruning of rootmulti.Store type")
			}
			qms := sdk.MultiStore(rms)
			latestVersion := qms.LatestVersion()

			latestHeight := rootmulti.GetLatestVersion(db)
			fmt.Printf("latest height of db/qms: %d %d\n", latestHeight, latestVersion)

			cacheMS, err := qms.CacheMultiStoreWithVersion(latestVersion)
			if err != nil {
				panic(err)
			}
			_ = cacheMS

			for name, key := range rms.StoreKeysByName() {
				store := rms.GetCommitKVStore(key)
				fmt.Printf("store name %s\n", name)
				switch store.GetStoreType() {
				case types.StoreTypeIAVL:
					iavlStore, ok := store.(*iavl.Store)
					if ok {
						commitID := iavlStore.LastCommitID()
						fmt.Printf("  version %d\n", commitID.Version)
					}
				}
			}

			keysByName := rms.StoreKeysByName()
			evidenceKey, found := keysByName["evidence"]
			if found {
				store := rms.GetCommitKVStore(evidenceKey)
				if store.GetStoreType() == types.StoreTypeIAVL {
					iavlStore, ok := store.(*iavl.Store)
					if !ok {
						panic("can't convert to iavl store")
					}
					//st, err := iavlStore.GetImmutable(latestVersion)
					_ = iavlStore
				}

			}

			keys := sdk.NewKVStoreKeys(
				authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
				govtypes.StoreKey, paramstypes.StoreKey,
				group.StoreKey,
				upgradetypes.StoreKey,
				evidencetypes.StoreKey,
				zetaCoreModuleTypes.StoreKey,
				zetaObserverModuleTypes.StoreKey,
				evmtypes.StoreKey, feemarkettypes.StoreKey,
				fungibleModuleTypes.StoreKey,
				emissionsModuleTypes.StoreKey,
				authzkeeper.StoreKey,
			)
			_ = keys

			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, "", "The database home directory")
	cmd.Flags().String(FlagAppDBBackend, "", "The type of database for application and snapshots databases")
	cmd.Flags().String(server.FlagPruning, pruningtypes.PruningOptionDefault, "Pruning strategy (default|nothing|everything|custom)")
	cmd.Flags().Uint64(server.FlagPruningKeepRecent, 0, "Number of recent heights to keep on disk (ignored if pruning is not 'custom')")
	cmd.Flags().Uint64(server.FlagPruningInterval, 10,
		`Height interval at which pruned heights are removed from disk (ignored if pruning is not 'custom'), 
		this is not used by this command but kept for compatibility with the complete pruning options`)

	return cmd
}

func openDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", backendType, dataDir)
}

func NewRollbackCosmosCmd(appCreator servertypes.AppCreator, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback-cosmos [height]",
		Short: "rollback cosmos-sdk app state to a specific height",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := server.GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir
			db, err := openDB(home, server.GetAppDBBackend(ctx.Viper))
			if err != nil {
				return err
			}
			app := appCreator(ctx.Logger, db, nil, ctx.Viper)

			// rollback the multistore only. Tendermint state is not rolled back
			height, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			if err := app.CommitMultiStore().RollbackToVersion(height); err != nil {
				return fmt.Errorf("failed to rollback to version: %w", err)
			}

			fmt.Printf("Rolled back state to height %d", height)
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	return cmd
}
