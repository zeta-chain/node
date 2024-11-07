package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/errors"
	"github.com/btcsuite/btcd/rpcclient"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	evmobserver "github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

type inboundOptions struct {
	Node    string
	ChainID string
}

var inboundOpts inboundOptions

func setupInboundOptions() {
	f, cfg := InboundCmd.PersistentFlags(), &inboundOpts

	f.StringVar(&cfg.Node, "node", "46.4.15.110", "zeta public ip address")
	f.StringVar(&cfg.ChainID, "chain-id", "athens_7001-1", "zeta chain id")
}

func InboundGetBallot(_ *cobra.Command, args []string) error {
	cobra.ExactArgs(2)

	cfg, err := config.Load(globalOpts.ZetacoreHome)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	inboundHash := args[0]

	chainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to parse chain id")
	}

	// create a new zetacore client
	client, err := zetacore.NewClient(
		&keys.Keys{OperatorAddress: sdk.MustAccAddressFromBech32(sample.AccAddress())},
		inboundOpts.Node,
		"",
		inboundOpts.ChainID,
		zerolog.Nop(),
	)
	if err != nil {
		return err
	}

	appContext := zctx.New(cfg, nil, zerolog.Nop())
	ctx := zctx.WithAppContext(context.Background(), appContext)

	if err := client.UpdateAppContext(ctx, appContext, zerolog.Nop()); err != nil {
		return errors.Wrap(err, "failed to update app context")
	}

	var ballotIdentifier string

	tssEthAddress, err := client.GetEVMTSSAddress(ctx)
	if err != nil {
		return err
	}

	chain, err := appContext.GetChain(chainID)
	if err != nil {
		return err
	}

	chainProto := chain.RawChain()

	// get ballot identifier according to the chain type
	if chain.IsEVM() {
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
				evmObserver.WithChain(*chainProto)
			}
		}
		hash := ethcommon.HexToHash(inboundHash)
		tx, isPending, err := evmObserver.TransactionByHash(inboundHash)
		if err != nil {
			return fmt.Errorf("tx not found on chain %s, %d", err.Error(), chain.ID())
		}

		if isPending {
			return fmt.Errorf("tx is still pending")
		}

		receipt, err := client.TransactionReceipt(context.Background(), hash)
		if err != nil {
			return fmt.Errorf("tx receipt not found on chain %s, %d", err.Error(), chain.ID())
		}

		params := chain.Params()

		evmObserver.SetChainParams(*params)

		if strings.EqualFold(tx.To, params.ConnectorContractAddress) {
			coinType = coin.CoinType_Zeta
		} else if strings.EqualFold(tx.To, params.Erc20CustodyContractAddress) {
			coinType = coin.CoinType_ERC20
		} else if strings.EqualFold(tx.To, tssEthAddress) {
			coinType = coin.CoinType_Gas
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
	} else if chain.IsBitcoin() {
		btcObserver := btcobserver.Observer{}
		btcObserver.WithZetacoreClient(client)
		btcObserver.WithChain(*chainProto)
		btcConfig, found := cfg.GetBTCConfig(chainID)
		if !found {
			return fmt.Errorf("unable to find config for BTC chain %d", chainID)
		}
		connCfg := &rpcclient.ConnConfig{
			Host:         btcConfig.RPCHost,
			User:         btcConfig.RPCUsername,
			Pass:         btcConfig.RPCPassword,
			HTTPPostMode: true,
			DisableTLS:   true,
			Params:       btcConfig.RPCParams,
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
