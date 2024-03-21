package zetabridge

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"go.nhat.io/grpcmock"
	"go.nhat.io/grpcmock/planner"
)

func setupMockServer(t *testing.T, serviceFunc any, method string, input any, expectedOutput any) *grpcmock.Server {
	l, err := net.Listen("tcp", "127.0.0.1:9090")
	require.NoError(t, err)

	server := grpcmock.MockUnstartedServer(
		grpcmock.RegisterService(serviceFunc),
		grpcmock.WithPlanner(planner.FirstMatch()),
		grpcmock.WithListener(l),
		func(s *grpcmock.Server) {
			s.ExpectUnary(method).
				UnlimitedTimes().
				WithPayload(input).
				Return(expectedOutput)
		},
	)(t)

	return server
}

func closeMockServer(t *testing.T, server *grpcmock.Server) {
	err := server.Close()
	require.NoError(t, err)
}

func setupCorBridge() (*ZetaCoreBridge, error) {
	return NewZetaCoreBridge(
		&keys.Keys{},
		"127.0.0.1",
		"",
		"zetachain_7000-1",
		false,
		&metrics.TelemetryServer{})
}

func TestZetaCoreBridge_GetBallot(t *testing.T) {
	expectedOutput := observertypes.QueryBallotByIdentifierResponse{
		BallotIdentifier: "456",
		Voters:           nil,
		ObservationType:  0,
		BallotStatus:     0,
	}
	input := observertypes.QueryBallotByIdentifierRequest{BallotIdentifier: "123"}
	method := "/zetachain.zetacore.observer.Query/BallotByIdentifier"
	server := setupMockServer(t, observertypes.RegisterQueryServer, method, input, expectedOutput)
	server.Serve()
	defer closeMockServer(t, server)

	zetabridge, err := setupCorBridge()
	require.NoError(t, err)

	resp, err := zetabridge.GetBallotByID("123")
	require.NoError(t, err)
	require.Equal(t, expectedOutput, *resp)
}
