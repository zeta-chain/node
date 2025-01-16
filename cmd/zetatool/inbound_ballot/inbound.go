package inbound_ballot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/node/pkg/chains"
	zetacorerpc "github.com/zeta-chain/node/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/cmd/zetatool/config"
)

func NewFetchInboundBallotCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "inbound",
		Short: "Fetch Inbound ballot from the inbound hash",
		RunE:  InboundGetBallot,
	}
}

func InboundGetBallot(cmd *cobra.Command, args []string) error {
	cobra.ExactArgs(2)

	configFile, err := cmd.Flags().GetString(config.FlagConfig)
	if err != nil {
		return err
	}

	cfg, err := config.GetConfig(configFile)
	if err != nil {
		panic(err)
	}
	inboundHash := args[0]
	fmt.Println("Inbound Hash: ", inboundHash)

	inboundChainID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		panic(err)
	}

	var unsecureGRPC = grpc.WithTransportCredentials(insecure.NewCredentials())
	zetacoreClient, err := zetacorerpc.NewGRPCClients(cfg.ZetaGRPC, unsecureGRPC)
	if err != nil {
		panic(err)
	}

	//zetacoreClient, err := zetacorerpc.NewGRPCClients(
	//	cfg.ZetaGRPC,
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//	grpc.WithBlock(),
	//)

	observationChain, found := chains.GetChainFromChainID(inboundChainID, []chains.Chain{})
	if !found {
		fmt.Println("Chain not found")
	}
	ctx := context.Background()
	ballotIdentifier := ""

	if observationChain.IsEVMChain() {
		ballotIdentifier, err = EvmInboundBallotIdentified(ctx, *cfg, zetacoreClient, inboundHash, observationChain, cfg.ZetaChainID)
		if err != nil {
			return fmt.Errorf("failed to get inbound ballot for evm chain %d, %s", observationChain.ChainId, err.Error())
		}
	}

	if observationChain.IsBitcoinChain() {
		ballotIdentifier, err = BtcInboundBallotIdentified(ctx, *cfg, zetacoreClient, inboundHash, observationChain, cfg.ZetaChainID)
		if err != nil {
			return fmt.Errorf("failed to get inbound ballot for bitcoin chain %d, %s", observationChain.ChainId, err.Error())
		}
	}

	if observationChain.IsSolanaChain() {
		ballotIdentifier, err = SolanaInboundBallotIdentified(ctx, *cfg, zetacoreClient, inboundHash, observationChain, cfg.ZetaChainID)
		if err != nil {
			return fmt.Errorf("failed to get inbound ballot for solana chain %d, %s", observationChain.ChainId, err.Error())
		}
	}

	fmt.Println("Ballot Identifier: ", ballotIdentifier)
	return nil
}
