package zetabridge

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestZetaCoreBridge_GetBallot(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(observertypes.RegisterQueryServer),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(l),
		func(s *grpcmock.Server) {
			s.ExpectUnary("/zetachain.zetacore.observer.Query/BallotByIdentifier").
				UnlimitedTimes().
				Return(observertypes.QueryBallotByIdentifierResponse{
					BallotIdentifier: "456",
					Voters:           nil,
					ObservationType:  0,
					BallotStatus:     0,
				})
		},
	)(t)
	server.Serve()
	defer func(server *grpcmock.Server) {
		err := server.Close()
		require.NoError(t, err)
	}(server)

	grpcConn, err := grpc.Dial(
		fmt.Sprintf("%s:9090", "127.0.0.1"),
		grpc.WithInsecure(),
	)
	client := observertypes.NewQueryClient(grpcConn)
	resp, err := client.BallotByIdentifier(context.Background(), &observertypes.QueryBallotByIdentifierRequest{
		BallotIdentifier: "123",
	})

	fmt.Println(resp)
}
