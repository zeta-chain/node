package local

import (
	"context"
	"errors"
	"strings"

	"path/filepath"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// GetConfig returns config from file from the command line flag
func GetConfig(cmd *cobra.Command) (config.Config, error) {
	configFile, err := cmd.Flags().GetString(FlagConfigFile)
	if err != nil {
		return config.Config{}, err
	}

	// use default config if no config file is specified
	if configFile == "" {
		return config.DefaultConfig(), nil
	}

	configFile, err = filepath.Abs(configFile)
	if err != nil {
		return config.Config{}, err
	}

	return config.ReadConfig(configFile)
}

// setCosmosConfig set account prefix to zeta
func setCosmosConfig() {
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()
}

// initTestRunner initializes a runner form tests
// it creates a runner with an account and copy contracts from deployer runner
func initTestRunner(
	name string,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	userAddress ethcommon.Address,
	userPrivKey string,
	logger *runner.Logger,
) (*runner.E2ERunner, error) {
	// initialize runner for test
	testRunner, err := zetae2econfig.RunnerFromConfig(
		deployerRunner.Ctx,
		name,
		deployerRunner.CtxCancel,
		conf,
		userAddress,
		userPrivKey,
		utils.FungibleAdminName,
		FungibleAdminMnemonic,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// copy contracts from deployer runner
	if err := testRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}

	return testRunner, nil
}

// waitKeygenHeight waits for keygen height
func waitKeygenHeight(
	ctx context.Context,
	cctxClient crosschaintypes.QueryClient,
	logger *runner.Logger,
) {
	// wait for keygen to be completed
	keygenHeight := int64(60)
	logger.Print("⏳ wait height %v for keygen to be completed", keygenHeight)
	for {
		time.Sleep(2 * time.Second)
		response, err := cctxClient.LastZetaHeight(ctx, &crosschaintypes.QueryLastZetaHeightRequest{})
		if err != nil {
			logger.Error("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= keygenHeight {
			break
		}
		logger.Info("Last ZetaHeight: %d", response.Height)
	}
}

func MonitorTxPriorityInBlocks(ctx context.Context, conf config.Config, logger *runner.Logger, errCh chan error) {
	rpc, err := rpchttp.New(conf.RPCs.ZetaCoreRPC, "/websocket")
	if err != nil {
		errCh <- err
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			errCh <- nil
		case <-ticker.C:
			block, err := rpc.Block(ctx, nil)
			if err != nil {
				errCh <- err
			}

			// iterate txs events and check if some MsgEthereumTx is included above crosschain and observer txs
			nonSystemTxFound := false
			for _, tx := range block.Block.Txs {
				txRes, err := rpc.Tx(context.Background(), tx.Hash(), false)
				if err != nil {
					continue
				}

				for _, ev := range txRes.TxResult.Events {
					for _, attr := range ev.Attributes {
						if attr.Key == "msg_type_url" {
							if strings.Contains(attr.Value, "zetachain.zetacore.crosschain.MsgVote") ||
								strings.Contains(attr.Value, "zetachain.zetacore.observer.MsgVote") ||
								strings.Contains(attr.Value, "zetachain.zetacore.observer.MsgAddBlameVote") {
								if nonSystemTxFound {
									errCh <- errors.New("wrong tx priority, system tx not on top")
								}
							}
						}
						if attr.Key == "action" {
							if strings.Contains(attr.Value, "ethermint.evm.v1.MsgEthereumTx") {
								nonSystemTxFound = true
							}
						}
					}
				}
			}
		}
	}
}
