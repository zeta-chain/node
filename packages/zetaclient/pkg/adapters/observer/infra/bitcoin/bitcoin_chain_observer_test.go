package bitcoin

import (
	"context"
	"os"
	"testing"

	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/config"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type BitcoinObserverTestSuite struct {
	suite.Suite
	obs *BTCChainObserver
	ctx context.Context
}

func (suite *BitcoinObserverTestSuite) SetupTest() {
	//os.Setenv("BITCOIN_ENDPOINT", "https://btc.getblock.io/testnet/?api_key=14a8dcca-3d41-4e17-a0bb-4c2d4dc2a524")
	//os.Setenv("BITCOIN_ENDPOINT", "https://hidden-holy-firefly.btc-testnet.discover.quiknode.pro/2610716f2259558b46f50a852032b5d09827aeaa/")
	//	os.Setenv("BITCOIN_ENDPOINT", "https://btc.getblock.io/mainnet/?api_key=14a8dcca-3d41-4e17-a0bb-4c2d4dc2a524")
	os.Setenv("BITCOIN_ENDPOINT", "https://btc.getblock.io/testnet/?api_key=14a8dcca-3d41-4e17-a0bb-4c2d4dc2a524")
	suite.ctx = context.Background()
	logger, _ := zap.NewDevelopment()
	cfg := config.MustGetConfig()
	obs, err := NewBTCChainObserver(suite.ctx, cfg, common.Chain("BITCOIN"), logger.Sugar())
	suite.Require().NoError(err)
	suite.obs = obs
}

func (suite *BitcoinObserverTestSuite) TearDownSuite() {
}

func (suite *BitcoinObserverTestSuite) TestAll() {
	blockNumber, err := suite.obs.GetBlockHeight(suite.ctx)
	if err != nil {
		suite.T().Logf("ERR=>%v\n", err)
	} else {
		suite.T().Logf("BLOCK=>%d\n", blockNumber)
	}
}

// TestBitcoinObserver is the entry point of this test suite
func TestBitcoinObserver(t *testing.T) {
	suite.Run(t, new(BitcoinObserverTestSuite))
}
