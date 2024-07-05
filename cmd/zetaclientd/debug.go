package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
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
	DebugCmd().Flags().
		StringVar(&debugArgs.zetaCoreHome, "core-home", "/Users/tanmay/.zetacored", "peer address, e.g. /dns/tss1/tcp/6668/ipfs/16Uiu2HAmACG5DtqmQsHtXg4G2sLS65ttv84e7MrL4kapkjfmhxAp")
	DebugCmd().Flags().StringVar(&debugArgs.zetaNode, "node", "46.4.15.110", "public ip address")
	DebugCmd().Flags().StringVar(&debugArgs.zetaChainID, "chain-id", "athens_7001-1", "pre-params file path")
}

func DebugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-inbound-ballot [inboundHash] [chainID]",
		Short: "provide txHash and chainID to get the ballot status for the txHash",
		RunE:  debugCmd,
	}
}

func debugCmd(_ *cobra.Command, args []string) error {
	cobra.ExactArgs(2)
	cfg, err := config.Load(debugArgs.zetaCoreHome)
	if err != nil {
		return err
	}

	appContext := zctx.New(cfg, zerolog.Nop())
	ctx := zctx.WithAppContext(context.Background(), appContext)

	chainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}

	inboundHash := args[0]
	var ballotIdentifier string

	// create a new zetacore client
	client, err := zetacore.NewClient(
		&keys.Keys{OperatorAddress: sdk.MustAccAddressFromBech32(sample.AccAddress())},
		debugArgs.zetaNode,
		"",
		debugArgs.zetaChainID,
		false,
		zerolog.Nop(),
	)
	if err != nil {
		return err
	}
	chainParams, err := client.GetChainParams(ctx)
	if err != nil {
		return err
	}
	tssEthAddress, err := client.GetEVMTSSAddress(ctx)
	if err != nil {
		return err
	}
	chain, found := chains.GetChainFromChainID(chainID, appContext.GetAdditionalChains())
	if !found {
		return fmt.Errorf("invalid chain id")
	}

	// get ballot identifier according to the chain type
	if chains.IsEVMChain(chain.ChainId, appContext.GetAdditionalChains()) {
		evmObserver := evmobserver.Observer{}
		evmObserver.WithZetacoreClient(client)
		var ethRPC *ethrpc.EthRPC
		var client *ethclient.Client
		coinType := coin.CoinType_Cmd
		for chainIDFromConfig, evmConfig := range cfg.GetAllEVMConfigs() {
			if chainIDFromConfig == chainID {
				ethRPC = ethrpc.NewEthRPC(evmConfig.Endpoint)
				client, err = ethclient.Dial(evmConfig.Endpoint)
				if err != nil {
					return err
				}
				evmObserver.WithEvmClient(client)
				evmObserver.WithEvmJSONRPC(ethRPC)
				evmObserver.WithChain(chain)
			}
		}
		hash := ethcommon.HexToHash(inboundHash)
		tx, isPending, err := evmObserver.TransactionByHash(inboundHash)
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
				evmChainParams, found := appContext.GetEVMChainParams(chainID)
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
			ballotIdentifier, err = evmObserver.CheckAndVoteInboundTokenZeta(ctx, tx, receipt, false)
			if err != nil {
				return err
			}

		case coin.CoinType_ERC20:
			ballotIdentifier, err = evmObserver.CheckAndVoteInboundTokenERC20(ctx, tx, receipt, false)
			if err != nil {
				return err
			}

		case coin.CoinType_Gas:
			ballotIdentifier, err = evmObserver.CheckAndVoteInboundTokenGas(ctx, tx, receipt, false)
			if err != nil {
				return err
			}
		default:
			fmt.Println("CoinType not detected")
		}
		fmt.Println("CoinType : ", coinType)
	} else if chains.IsBitcoinChain(chain.ChainId, appContext.GetAdditionalChains()) {
		btcObserver := btcobserver.Observer{}
		btcObserver.WithZetacoreClient(client)
		btcObserver.WithChain(chain)
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
		ballotIdentifier, err = btcObserver.CheckReceiptForBtcTxHash(ctx, inboundHash, false)
		if err != nil {
			return err
		}
	}
	fmt.Println("BallotIdentifier : ", ballotIdentifier)

	// query ballot
	ballot, err := client.GetBallot(ctx, ballotIdentifier)
	if err != nil {
		return err
	}

	for _, vote := range ballot.Voters {
		fmt.Printf("%s : %s \n", vote.VoterAddress, vote.VoteType)
	}
	fmt.Println("BallotStatus : ", ballot.BallotStatus)

	return nil
}
