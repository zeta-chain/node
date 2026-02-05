package network

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/log"
	cmtcfg "github.com/cometbft/cometbft/config"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/node"
	"github.com/cometbft/cometbft/p2p"
	pvm "github.com/cometbft/cometbft/privval"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/rpc/client/local"
	tmtypes "github.com/cometbft/cometbft/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	servergrpc "github.com/cosmos/cosmos-sdk/server/grpc"
	servercmtlog "github.com/cosmos/cosmos-sdk/server/log"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func startInProcess(cfg Config, val *Validator) error {
	logger := val.Ctx.Logger
	tmCfg := val.Ctx.Config
	tmCfg.Instrumentation.Prometheus = false

	if err := val.AppConfig.ValidateBasic(); err != nil {
		return err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(tmCfg.NodeKeyFile())
	if err != nil {
		return err
	}

	app := cfg.AppConstructor(*val)
	val.app = app

	appGenesisProvider := func() (*tmtypes.GenesisDoc, error) {
		appGenesis, err := genutiltypes.AppGenesisFromFile(tmCfg.GenesisFile())
		if err != nil {
			return nil, err
		}

		return appGenesis.ToGenesisDoc()
	}

	cmtApp := server.NewCometABCIWrapper(app)
	tmNode, err := node.NewNode( //resleak:notresource
		tmCfg,
		pvm.LoadOrGenFilePV(tmCfg.PrivValidatorKeyFile(), tmCfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(cmtApp),
		appGenesisProvider,
		cmtcfg.DefaultDBProvider,
		node.DefaultMetricsProvider(tmCfg.Instrumentation),
		servercmtlog.CometLoggerWrapper{Logger: logger.With("module", val.Moniker)},
	)
	if err != nil {
		return err
	}

	if err := tmNode.Start(); err != nil {
		return err
	}
	val.tmNode = tmNode

	if val.RPCAddress != "" {
		val.RPCClient = local.New(tmNode)
	}

	ctx := context.Background()

	// We'll need a RPC client if the validator exposes a gRPC or REST endpoint.
	if val.APIAddress != "" || val.AppConfig.GRPC.Enable {
		val.ClientCtx = val.ClientCtx.
			WithClient(val.RPCClient)

		app.RegisterTxService(val.ClientCtx)
		app.RegisterTendermintService(val.ClientCtx)
		app.RegisterNodeService(val.ClientCtx, *val.AppConfig)
	}

	if val.APIAddress != "" {
		apiSrv := api.New(val.ClientCtx, logger.With(log.ModuleKey, "api-server"), val.grpc)
		app.RegisterAPIRoutes(apiSrv, val.AppConfig.API)

		errCh := make(chan error)

		go func() {
			if err := apiSrv.Start(ctx, *val.AppConfig); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(5 * time.Second): // assume server started successfully
		}

		val.api = apiSrv
	}

	grpcCfg := val.AppConfig.GRPC
	if val.AppConfig.GRPC.Enable {
		grpcSrv, err := servergrpc.NewGRPCServer(val.ClientCtx, app, grpcCfg)
		if err != nil {
			return err
		}

		errCh := make(chan error)

		// Start the gRPC server in a goroutine. Note, the provided ctx will ensure
		// that the server is gracefully shut down.
		go func() {
			if err := servergrpc.StartGRPCServer(
				ctx,
				logger.With(log.ModuleKey, "grpc-server"),
				grpcCfg,
				grpcSrv,
			); err != nil {
				errCh <- err
			}
		}()

		select {
		case err := <-errCh:
			return err
		case <-time.After(5 * time.Second): // assume server started successfully
		}

		val.grpc = grpcSrv
	}
	return nil
}

func collectGenFiles(cfg Config, vals []*Validator, outputDir string) error {
	genTime := tmtime.Now()

	for i := 0; i < cfg.NumValidators; i++ {
		tmCfg := vals[i].Ctx.Config

		nodeDir := filepath.Join(outputDir, vals[i].Moniker, "simd")
		gentxsDir := filepath.Join(outputDir, "gentxs")

		tmCfg.Moniker = vals[i].Moniker
		tmCfg.SetRoot(nodeDir)

		initCfg := genutiltypes.NewInitConfig(cfg.ChainID, gentxsDir, vals[i].NodeID, vals[i].PubKey)

		genFile := tmCfg.GenesisFile()
		genDoc, err := tmtypes.GenesisDocFromFile(genFile)
		if err != nil {
			return err
		}

		appGenesis, err := genutiltypes.AppGenesisFromFile(genFile)
		if err != nil {
			return err
		}

		appState, err := genutil.GenAppStateFromConfig(
			cfg.Codec,
			cfg.TxConfig,
			tmCfg,
			initCfg,
			appGenesis,
			banktypes.GenesisBalancesIterator{},
			genutiltypes.DefaultMessageValidator,
			cfg.TxConfig.SigningContext().ValidatorAddressCodec(),
		)
		if err != nil {
			return err
		}

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := ExportGenesisFileWithTimeAndConsensusParams(
			genFile,
			cfg.ChainID,
			nil,
			appState,
			genTime,
			*genDoc.ConsensusParams,
		); err != nil {
			return err
		}
	}

	return nil
}

func ExportGenesisFileWithTimeAndConsensusParams(
	genFile, chainID string, validators []tmtypes.GenesisValidator,
	appState json.RawMessage, genTime time.Time, consensusParams tmtypes.ConsensusParams,
) error {
	genDoc := tmtypes.GenesisDoc{
		GenesisTime:     genTime,
		ChainID:         chainID,
		Validators:      validators,
		AppState:        appState,
		ConsensusParams: &consensusParams,
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genFile)
}

func initGenFiles(
	cfg Config,
	genAccounts []authtypes.GenesisAccount,
	genBalances []banktypes.Balance,
	genFiles []string,
) error {
	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = append(authGenState.Accounts, accounts...)
	cfg.GenesisState[authtypes.ModuleName] = cfg.Codec.MustMarshalJSON(&authGenState)

	// set the balances in the genesis state
	var bankGenState banktypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[banktypes.ModuleName], &bankGenState)

	bankGenState.Balances = append(bankGenState.Balances, genBalances...)
	cfg.GenesisState[banktypes.ModuleName] = cfg.Codec.MustMarshalJSON(&bankGenState)

	appGenStateJSON, err := json.MarshalIndent(cfg.GenesisState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := tmtypes.GenesisDoc{
		ConsensusParams: &tmtypes.ConsensusParams{
			Block: tmtypes.BlockParams{
				MaxBytes: 200000,
				MaxGas:   2000000,
			},
			Evidence: tmtypes.EvidenceParams{
				MaxAgeNumBlocks: 302400,
				MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
				MaxBytes:        10000,
			},
			Validator: tmtypes.ValidatorParams{
				PubKeyTypes: []string{
					tmtypes.ABCIPubKeyTypeEd25519,
				},
			},
		},
		ChainID:    cfg.ChainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < cfg.NumValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir) //nolint:gocritic
	file := filepath.Join(writePath, name)

	err := tmos.EnsureDir(writePath, 0o755) // #nosec G301
	if err != nil {
		return err
	}

	err = os.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}
