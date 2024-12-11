package main

import (
	"fmt"
	"net"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
	"gitlab.com/thorchain/tss/go-tss/conversion"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/zeta-chain/node/pkg/rpc"
	"github.com/zeta-chain/node/pkg/sdkconfig"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func NewGetZetaclientBootstrap() *cobra.Command {
	var ConfigureZetaclientBootstrapCmd = &cobra.Command{
		Use:   "get-zetaclient-bootstrap",
		Short: "get bootstrap address book entries for zetaclient",
		RunE:  getZetaclientBootstrap,
	}

	return ConfigureZetaclientBootstrapCmd
}

func getZetaclientBootstrap(cmd *cobra.Command, _ []string) error {
	sdkconfig.SetDefault(true)
	rpcClient, err := rpc.NewGRPCClients(
		"zetacore0:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("get zetacore rpc client: %w", err)
	}
	var res *observertypes.QueryAllNodeAccountResponse
	for {
		res, err = rpcClient.Observer.NodeAccountAll(cmd.Context(), &observertypes.QueryAllNodeAccountRequest{})
		if err != nil {
			return fmt.Errorf("get all node accounts: %w", err)
		}
		if len(res.NodeAccount) > 1 {
			break
		}
		fmt.Fprintln(cmd.OutOrStderr(), "waiting for node accounts")
	}

	// note that we deliberately do not filter ourselfs/localhost
	// to mirror the production configuration
	for _, account := range res.NodeAccount {
		accAddr, err := sdk.AccAddressFromBech32(account.Operator)
		if err != nil {
			return err
		}
		valAddr := sdk.ValAddress(accAddr).String()
		validatorRes, err := rpcClient.Staking.Validator(cmd.Context(), &stakingtypes.QueryValidatorRequest{
			ValidatorAddr: valAddr,
		})
		if err != nil {
			return fmt.Errorf("getting validator info for %s: %w", account.Operator, err)
		}
		// in localnet, moniker is also the hostname
		moniker := validatorRes.Validator.Description.Moniker

		peerID, err := conversion.Bech32PubkeyToPeerID(account.GranteePubkey.Secp256k1.String())
		if err != nil {
			return fmt.Errorf("converting pubkey to peerID: %w", err)
		}
		zetaclientHostname := strings.ReplaceAll(moniker, "zetacore", "zetaclient")

		// resolve the hostname
		// something in libp2p/go-tss requires /ip4/<ip> and doesn't tolerate /dns4/<hostname>
		ipAddresses, err := net.LookupIP(zetaclientHostname)
		if err != nil {
			return fmt.Errorf("failed to resolve hostname %s: %w", zetaclientHostname, err)
		}
		if len(ipAddresses) == 0 {
			return fmt.Errorf("no IP addresses found for hostname %s", zetaclientHostname)
		}
		ipv4Address := ""
		for _, ip := range ipAddresses {
			if ip.To4() != nil {
				ipv4Address = ip.String()
				break
			}
		}
		if ipv4Address == "" {
			return fmt.Errorf("no IPv4 address found for hostname %s", zetaclientHostname)
		}
		fmt.Printf("/ip4/%s/tcp/6668/p2p/%s\n", ipv4Address, peerID.String())
	}

	return nil
}
