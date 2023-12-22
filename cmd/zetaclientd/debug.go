package main

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/btcsuite/btcd/rpcclient"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

var debugArgs = debugArguments{}

type debugArguments struct {
	zetaCoreHome string
	zetaNode     string
	zetaChainID  string
}

func init() {
	RootCmd.AddCommand(DebugCmd())
	DebugCmd().Flags().StringVar(&debugArgs.zetaCoreHome, "core-home", "/Users/tanmay/.zetacored", "peer address, e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp")
	DebugCmd().Flags().StringVar(&debugArgs.zetaNode, "node", "46.4.15.110", "public ip address")
	DebugCmd().Flags().StringVar(&debugArgs.zetaChainID, "chain-id", "athens_7001-1", "pre-params file path")
}

func DebugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-ballot-from-intx [txHash] [chainID]",
		Short: "provide txHash and chainID to get the ballot status for the txHash",
		RunE: func(cmd *cobra.Command, args []string) error {
			cobra.ExactArgs(2)
			cfg, err := config.Load(debugArgs.zetaCoreHome)
			if err != nil {
				return err
			}
			chainID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			txHash := args[0]
			var ballotIdentifier string
			chainLogger := zerolog.New(io.Discard).Level(zerolog.Disabled)

			telemetryServer := zetaclient.NewTelemetryServer()
			go func() {
				err := telemetryServer.Start()
				if err != nil {
					panic("telemetryServer error")
				}
			}()

			bridge, err := zetaclient.NewZetaCoreBridge(
				&zetaclient.Keys{OperatorAddress: sdk.MustAccAddressFromBech32(sample.AccAddress())},
				debugArgs.zetaNode,
				"",
				debugArgs.zetaChainID,
				false,
				telemetryServer)

			if err != nil {
				return err
			}
			coreParams, err := bridge.GetCoreParams()
			if err != nil {
				return err
			}
			tssEthAddress, err := bridge.GetEthTssAddress()
			if err != nil {
				return err
			}

			chain := common.GetChainFromChainID(chainID)
			if chain == nil {
				return fmt.Errorf("invalid chain id")
			}

			if common.IsEVMChain(chain.ChainId) {

				ob := zetaclient.EVMChainClient{
					Mu: &sync.Mutex{},
				}
				ob.WithZetaClient(bridge)
				ob.WithLogger(chainLogger)
				client := &ethclient.Client{}
				coinType := common.CoinType_Cmd
				for chain, evmConfig := range cfg.GetAllEVMConfigs() {
					if chainID == chain {
						client, err = ethclient.Dial(evmConfig.Endpoint)
						if err != nil {
							return err
						}
						ob.WithEvmClient(client)
						ob.WithChain(*common.GetChainFromChainID(chainID))
					}
				}
				hash := ethcommon.HexToHash(txHash)
				tx, isPending, err := client.TransactionByHash(context.Background(), hash)
				if err != nil {
					return fmt.Errorf("tx not found on chain %s , %d", err.Error(), chain.ChainId)
				}
				if isPending {
					return fmt.Errorf("tx is still pending")
				}

				for _, chainCoreParams := range coreParams {
					if chainCoreParams.ChainId == chainID {
						ob.WithParams(observertypes.CoreParams{
							ChainId:                     chainID,
							ConnectorContractAddress:    chainCoreParams.ConnectorContractAddress,
							ZetaTokenContractAddress:    chainCoreParams.ZetaTokenContractAddress,
							Erc20CustodyContractAddress: chainCoreParams.Erc20CustodyContractAddress,
						})
						cfg.EVMChainConfigs[chainID].ZetaTokenContractAddress = chainCoreParams.ZetaTokenContractAddress
						ob.SetConfig(cfg)
						if strings.EqualFold(tx.To().Hex(), chainCoreParams.ConnectorContractAddress) {
							coinType = common.CoinType_Zeta
						} else if strings.EqualFold(tx.To().Hex(), chainCoreParams.Erc20CustodyContractAddress) {
							coinType = common.CoinType_ERC20
						} else if strings.EqualFold(tx.To().Hex(), tssEthAddress) {
							coinType = common.CoinType_Gas
						}

					}
				}

				switch coinType {
				case common.CoinType_Zeta:
					ballotIdentifier, err = ob.CheckReceiptForCoinTypeZeta(txHash, false)
					if err != nil {
						return err
					}

				case common.CoinType_ERC20:
					ballotIdentifier, err = ob.CheckReceiptForCoinTypeERC20(txHash, false)
					if err != nil {
						return err
					}

				case common.CoinType_Gas:
					ballotIdentifier, err = ob.CheckReceiptForCoinTypeGas(txHash, false)
					if err != nil {
						return err
					}
				default:
					fmt.Println("CoinType not detected")
				}
				fmt.Println("CoinType : ", coinType)
			} else if common.IsBitcoinChain(chain.ChainId) {
				obBtc := zetaclient.BitcoinChainClient{
					Mu: &sync.Mutex{},
				}
				obBtc.WithZetaClient(bridge)
				obBtc.WithLogger(chainLogger)
				obBtc.WithChain(*common.GetChainFromChainID(chainID))
				connCfg := &rpcclient.ConnConfig{
					Host:         cfg.BitcoinConfig.RPCHost,
					User:         cfg.BitcoinConfig.RPCUsername,
					Pass:         cfg.BitcoinConfig.RPCPassword,
					HTTPPostMode: true,
					DisableTLS:   true,
					Params:       cfg.BitcoinConfig.RPCParams,
				}

				btcClient, err := rpcclient.New(connCfg, nil)
				if err != nil {
					return err
				}
				obBtc.WithBtcClient(btcClient)
				ballotIdentifier, err = obBtc.CheckReceiptForBtcTxHash(txHash, false)
				if err != nil {
					return err
				}

			}
			fmt.Println("BallotIdentifier : ", ballotIdentifier)

			ballot, err := bridge.GetBallot(ballotIdentifier)
			if err != nil {
				return err
			}

			for _, vote := range ballot.Voters {
				fmt.Printf("%s : %s \n", vote.VoterAddress, vote.VoteType)
			}
			fmt.Println("BallotStatus : ", ballot.BallotStatus)

			return nil
		},
	}

	return cmd
}
