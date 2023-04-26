package main

import (
	etherminttypes "github.com/evmos/ethermint/types"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common"
	mc "github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"strings"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Configuration",
	RunE:  Initialize,
}

var initArgs = initArguments{}

type initArguments struct {
	enabledChains string
	peer          string
	preParamsPath string
	keygen        int64
	chainID       string
	zetacoreURL   string
	authzGranter  string
	authzHotkey   string
	devMode       bool
	level         int8

	p2pDiagnostic bool
}

func init() {
	RootCmd.AddCommand(InitCmd)

	InitCmd.Flags().StringVar(&initArgs.enabledChains, "enable-chains", "GOERLI,BSCTESTNET,MUMBAI,ROPSTEN,BAOBAB", "enable chains, comma separated list")
	InitCmd.Flags().StringVar(&initArgs.peer, "peer", "", "peer address, e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp")
	InitCmd.Flags().StringVar(&initArgs.preParamsPath, "pre-params", "", "pre-params file path")
	InitCmd.Flags().Int64Var(&initArgs.keygen, "keygen-block", 0, "keygen at block height (default: 0 means no keygen")
	InitCmd.Flags().StringVar(&initArgs.chainID, "chain-id", "athens-1", "chain id")
	InitCmd.Flags().StringVar(&initArgs.zetacoreURL, "zetacore-url", "127.0.0.1", "zetacore node URL")
	InitCmd.Flags().StringVar(&initArgs.authzGranter, "operator", "", "granter for the authorization , this should be operator address")
	InitCmd.Flags().StringVar(&initArgs.authzHotkey, "hotkey", "", "hotkey for zetaclient this key is used for TSS and ZetaClient operations")
	InitCmd.Flags().BoolVar(&initArgs.devMode, "dev", false, "dev mode: geth private network as goerli testnet")
	InitCmd.Flags().Int8Var(&initArgs.level, "log-level", int8(zerolog.InfoLevel), "log level (0:debug, 1:info, 2:warn, 3:error, 4:fatal, 5:panic , 6: NoLevel , 7: Disable)")
	InitCmd.Flags().BoolVar(&initArgs.p2pDiagnostic, "p2p-diagnostic", false, "enable p2p diagnostic")
}

func Initialize(_ *cobra.Command, _ []string) error {
	setHomeDir()

	//Create new config struct
	configData := config.New()

	//Populate new struct with cli arguments
	initEnabledChains(&configData)
	initChainID(&configData)
	configData.Peer = initArgs.peer
	configData.PreParamsPath = initArgs.preParamsPath
	configData.KeygenBlock = initArgs.keygen
	configData.ChainID = initArgs.chainID
	configData.ZetaCoreURL = initArgs.zetacoreURL
	configData.AuthzHotkey = initArgs.authzHotkey
	configData.AuthzGranter = initArgs.authzGranter
	configData.LogLevel = zerolog.Level(initArgs.level)
	configData.P2PDiagnostic = initArgs.p2pDiagnostic

	//Save config file
	return config.Save(&configData, rootArgs.zetaCoreHome)
}

func initEnabledChains(configData *config.Config) {
	chains := strings.Split(initArgs.enabledChains, ",")
	chainList := []common.Chain{}
	supportedChains := mc.GetSupportedChains()
	for _, chain := range chains {
		for _, supportedChain := range supportedChains {
			if supportedChain.ChainName.String() == chain {
				if !initArgs.devMode && supportedChain.ChainId == 1337 {
					panic("GoerliLocalNetChain can only be enabled in Dev Mode ")
				}
				chainList = append(chainList, *supportedChain)
			}
		}
	}
	configData.ChainsEnabled = chainList
}

func initChainID(configData *config.Config) {
	ZEVMChainID, err := etherminttypes.ParseChainID(cmd.CHAINID)
	if err != nil {
		panic(err)
	}
	configData.ChainConfigs[common.ZetaChain().ChainName.String()].Chain.ChainId = ZEVMChainID.Int64()
}
