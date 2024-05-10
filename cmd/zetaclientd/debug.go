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
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"
	evmobserver "github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	clientcontext "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
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
		RunE: func(_ *cobra.Command, args []string) error {
			cobra.ExactArgs(2)
			cfg, err := config.Load(debugArgs.zetaCoreHome)
			if err != nil {
				return err
			}
			coreContext := clientcontext.NewZetaCoreContext(cfg)
			chainID, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}
			txHash := args[0]
			var ballotIdentifier string
			chainLogger := zerolog.New(io.Discard).Level(zerolog.Disabled)

			telemetryServer := metrics.NewTelemetryServer()
			go func() {
				err := telemetryServer.Start()
				if err != nil {
					panic("telemetryServer error")
				}
			}()

			client, err := zetacore.NewClient(
				&keys.Keys{OperatorAddress: sdk.MustAccAddressFromBech32(sample.AccAddress())},
				debugArgs.zetaNode,
				"",
				debugArgs.zetaChainID,
				false,
				telemetryServer)

			if err != nil {
				return err
			}
			chainParams, err := client.GetChainParams()
			if err != nil {
				return err
			}
			tssEthAddress, err := client.GetEthTssAddress()
			if err != nil {
				return err
			}

			chain := chains.GetChainFromChainID(chainID)
			if chain == nil {
				return fmt.Errorf("invalid chain id")
			}

			if chains.IsEVMChain(chain.ChainId) {

				evmObserver := evmobserver.Observer{
					Mu: &sync.Mutex{},
				}
				evmObserver.WithZetacoreClient(client)
				evmObserver.WithLogger(chainLogger)
				var ethRPC *ethrpc.EthRPC
				var client *ethclient.Client
				coinType := coin.CoinType_Cmd
				for chain, evmConfig := range cfg.GetAllEVMConfigs() {
					if chainID == chain {
						ethRPC = ethrpc.NewEthRPC(evmConfig.Endpoint)
						client, err = ethclient.Dial(evmConfig.Endpoint)
						if err != nil {
							return err
						}
						evmObserver.WithEvmClient(client)
						evmObserver.WithEvmJSONRPC(ethRPC)
						evmObserver.WithChain(*chains.GetChainFromChainID(chainID))
					}
				}
				hash := ethcommon.HexToHash(txHash)
				tx, isPending, err := evmObserver.TransactionByHash(txHash)
				if err != nil {
					return fmt.Errorf("tx not found on chain %s , %d", err.Error(), chain.ChainId)
				}
				if isPending {
					return fmt.Errorf("tx is still pending")
				}
				receipt, err := client.TransactionReceipt(context.Background(), hash)
				if err != nil {
					return fmt.Errorf("tx receipt not found on chain %s, %d", err.Error(), chain.ChainId)
				}

				for _, chainParams := range chainParams {
					if chainParams.ChainId == chainID {
						evmObserver.SetChainParams(observertypes.ChainParams{
							ChainId:                     chainID,
							ConnectorContractAddress:    chainParams.ConnectorContractAddress,
							ZetaTokenContractAddress:    chainParams.ZetaTokenContractAddress,
							Erc20CustodyContractAddress: chainParams.Erc20CustodyContractAddress,
						})
						evmChainParams, found := coreContext.GetEVMChainParams(chainID)
						if !found {
							return fmt.Errorf("missing chain params for chain %d", chainID)
						}
						evmChainParams.ZetaTokenContractAddress = chainParams.ZetaTokenContractAddress
						if strings.EqualFold(tx.To, chainParams.ConnectorContractAddress) {
							coinType = coin.CoinType_Zeta
						} else if strings.EqualFold(tx.To, chainParams.Erc20CustodyContractAddress) {
							coinType = coin.CoinType_ERC20
						} else if strings.EqualFold(tx.To, tssEthAddress) {
							coinType = coin.CoinType_Gas
						}
					}
				}

				switch coinType {
				case coin.CoinType_Zeta:
					ballotIdentifier, err = evmObserver.CheckAndVoteInboundTokenZeta(tx, receipt, false)
					if err != nil {
						return err
					}

				case coin.CoinType_ERC20:
					ballotIdentifier, err = evmObserver.CheckAndVoteInboundTokenERC20(tx, receipt, false)
					if err != nil {
						return err
					}

				case coin.CoinType_Gas:
					ballotIdentifier, err = evmObserver.CheckAndVoteInboundTokenGas(tx, receipt, false)
					if err != nil {
						return err
					}
				default:
					fmt.Println("CoinType not detected")
				}
				fmt.Println("CoinType : ", coinType)
			} else if chains.IsBitcoinChain(chain.ChainId) {
				btcObserver := btcobserver.Observer{
					Mu: &sync.Mutex{},
				}
				btcObserver.WithZetacoreClient(client)
				btcObserver.WithLogger(chainLogger)
				btcObserver.WithChain(*chains.GetChainFromChainID(chainID))
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
				btcObserver.WithBtcClient(btcClient)
				ballotIdentifier, err = btcObserver.CheckReceiptForBtcTxHash(txHash, false)
				if err != nil {
					return err
				}

			}
			fmt.Println("BallotIdentifier : ", ballotIdentifier)

			ballot, err := client.GetBallot(ballotIdentifier)
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
