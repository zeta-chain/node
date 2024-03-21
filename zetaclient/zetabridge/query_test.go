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

func TestZetaCoreBridge_GetCrosschainFlags(t *testing.T) {
}

func TestZetaCoreBridge_GetChainParamsForChainID(t *testing.T) {
}

func TestZetaCoreBridge_GetChainParams(t *testing.T) {
}

func TestZetaCoreBridge_GetUpgradePlan(t *testing.T) {
}

func TestZetaCoreBridge_GetAllCctx(t *testing.T) {
}

func TestZetaCoreBridge_GetCctxByHash(t *testing.T) {
}

func TestZetaCoreBridge_GetCctxByNonce(t *testing.T) {
}

func TestZetaCoreBridge_GetObserverList(t *testing.T) {
}

func TestZetaCoreBridge_ListPendingCctx(t *testing.T) {
}

func TestZetaCoreBridge_GetAbortedZetaAmount(t *testing.T) {
}

func TestZetaCoreBridge_GetGenesisSupply(t *testing.T) {
}

func TestZetaCoreBridge_GetZetaTokenSupplyOnNode(t *testing.T) {
}

func TestZetaCoreBridge_GetLastBlockHeight(t *testing.T) {
}

func TestZetaCoreBridge_GetLatestZetaBlock(t *testing.T) {
}

func TestZetaCoreBridge_GetNodeInfo(t *testing.T) {
}

func TestZetaCoreBridge_GetLastBlockHeightByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetZetaBlockHeight(t *testing.T) {
}

func TestZetaCoreBridge_GetBaseGasPrice(t *testing.T) {
}

func TestZetaCoreBridge_GetNonceByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetAllNodeAccounts(t *testing.T) {
}

func TestZetaCoreBridge_GetKeyGen(t *testing.T) {
}

func TestZetaCoreBridge_GetBallotByID(t *testing.T) {
}

func TestZetaCoreBridge_GetInboundTrackersForChain(t *testing.T) {
}

func TestZetaCoreBridge_GetCurrentTss(t *testing.T) {
}

func TestZetaCoreBridge_GetEthTssAddress(t *testing.T) {
}

func TestZetaCoreBridge_GetBtcTssAddress(t *testing.T) {
}

func TestZetaCoreBridge_GetTssHistory(t *testing.T) {
}

func TestZetaCoreBridge_GetOutTxTracker(t *testing.T) {
}

func TestZetaCoreBridge_GetAllOutTxTrackerByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetPendingNoncesByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetBlockHeaderStateByChain(t *testing.T) {
}

func TestZetaCoreBridge_GetSupportedChains(t *testing.T) {
}

func TestZetaCoreBridge_GetPendingNonces(t *testing.T) {
}

func TestZetaCoreBridge_Prove(t *testing.T) {
}

func TestZetaCoreBridge_HasVoted(t *testing.T) {
}

func TestZetaCoreBridge_GetZetaHotKeyBalance(t *testing.T) {
}
