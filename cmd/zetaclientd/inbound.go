package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	btcrpc "github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	evmobserver "github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/keys"
	"github.com/zeta-chain/node/zetaclient/orchestrator"
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

	err = orchestrator.UpdateAppContext(ctx, appContext, client, zerolog.Nop())
	if err != nil {
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

	baseLogger := base.Logger{Std: zerolog.Nop(), Compliance: zerolog.Nop()}

	observers, err := orchestrator.CreateChainObserverMap(ctx, client, nil, db.SqliteInMemory, baseLogger, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create chain observer map")
	}

	// get ballot identifier according to the chain type
	if chain.IsEVM() {
		observer, ok := observers[chainID]
		if !ok {
			return fmt.Errorf("observer not found for evm chain %d", chain.ID())
		}

		evmObserver, ok := observer.(*evmobserver.Observer)
		if !ok {
			return fmt.Errorf("observer is not evm observer for chain %d", chain.ID())
		}

		coinType := coin.CoinType_Cmd
		hash := ethcommon.HexToHash(inboundHash)
		tx, isPending, err := evmObserver.TransactionByHash(inboundHash)
		if err != nil {
			return fmt.Errorf("tx not found on chain %s, %d", err.Error(), chain.ID())
		}

		if isPending {
			return fmt.Errorf("tx is still pending")
		}

		receipt, err := evmObserver.TransactionReceipt(ctx, hash)
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
		bitcoinConfig, found := appContext.Config().GetBTCConfig(chain.ID())
		if !found {
			return fmt.Errorf("unable to find btc config")
		}

		rpcClient, err := btcrpc.NewRPCClient(bitcoinConfig)
		if err != nil {
			return errors.Wrap(err, "unable to create rpc client")
		}

		database, err := db.NewFromSqliteInMemory(true)
		if err != nil {
			return errors.Wrap(err, "unable to open database")
		}

		observer, err := btcobserver.NewObserver(
			*chain.RawChain(),
			rpcClient,
			*chain.Params(),
			client,
			nil,
			database,
			baseLogger,
			nil,
		)
		if err != nil {
			return errors.Wrap(err, "unable to create btc observer")
		}

		ballotIdentifier, err = observer.CheckReceiptForBtcTxHash(ctx, inboundHash, false)
		if err != nil {
			return err
		}
	}

	fmt.Println("BallotIdentifier: ", ballotIdentifier)

	// query ballot
	ballot, err := client.GetBallot(ctx, ballotIdentifier)
	if err != nil {
		return err
	}

	for _, vote := range ballot.Voters {
		fmt.Printf("%s: %s\n", vote.VoterAddress, vote.VoteType)
	}

	fmt.Println("BallotStatus: ", ballot.BallotStatus)

	return nil
}
