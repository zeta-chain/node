package zetaclient

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

type ZetaSupplyChecker struct {
	chainClientMap map[common.Chain]ChainClient
	cfg            config.Config
	evmClient      *ethclient.Client
	zetaClient     *ZetaCoreBridge
	ticker         DynamicTicker
	stop           chan struct{}
}

func (zs *ZetaSupplyChecker) Start() {
	defer zs.ticker.Stop()
	for {
		select {
		case <-zs.ticker.C():
			fmt.Println("Checking Zeta Supply")
		case <-zs.stop:
			return
		}
	}
}

func (zs *ZetaSupplyChecker) CheckZetaTokenSupply() error {
	externalChainEvmChain := make([]common.Chain, 0)
	ethereumChain := common.Chain{}
	for chain, _ := range zs.chainClientMap {
		if chain.IsExternalChain() && common.IsEVMChain(chain.ChainId) && !common.IsEthereumChain(chain.ChainId) {
			externalChainEvmChain = append(externalChainEvmChain, chain)
		}
		if common.IsEthereumChain(chain.ChainId) {
			ethereumChain = chain
		}
	}
	// Todo add checks for ethereum chain
	externalChainTotalSupply := big.NewInt(0)
	for _, chain := range externalChainEvmChain {
		zetaTokenAddressString := zs.cfg.EVMChainConfigs[chain.ChainId].CoreParams.ZetaTokenContractAddress
		zetaTokenAddress := ethcommon.HexToAddress(zetaTokenAddressString)
		zetatokenNonEth, err := FetchZetaZetaNonEthTokenContract(zetaTokenAddress, zs.evmClient)
		if err != nil {
			return err
		}
		totalSupply, err := zetatokenNonEth.TotalSupply(nil)
		if err != nil {
			return err
		}
		externalChainTotalSupply.Add(externalChainTotalSupply, totalSupply)
	}
	ethConnectorAddressString := zs.cfg.EVMChainConfigs[ethereumChain.ChainId].CoreParams.ConnectorContractAddress
	ethConnectorAddress := ethcommon.HexToAddress(ethConnectorAddressString)
	ethConnectorContract, err := FetchConnectorContractEth(ethConnectorAddress, zs.evmClient)
	if err != nil {
		return err
	}
	ethLockedAmount, err := ethConnectorContract.GetLockedAmount(nil)
	if err != nil {
		return err
	}
	pendingCCTX
	zs.zetaClient.GetAllPendingCctx()
	return nil
}
