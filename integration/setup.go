package integration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/testutil/network"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	txHeight    int64
	queryClient tx.ServiceClient
	txRes       sdk.TxResponse
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)
	s.Require().NotNil(s.network)

	val := s.network.Validators[0]

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.queryClient = tx.NewServiceClient(val.ClientCtx)

	s.Require().NoError(s.network.WaitForNextBlock())
	height, err := s.network.LatestHeight()

	s.txHeight = height
}
