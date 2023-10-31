package zetaclient

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

type ZetaSupplyChecker struct {
	cfg              *config.Config
	evmClient        map[int64]*ethclient.Client
	zetaClient       *ZetaCoreBridge
	ticker           *DynamicTicker
	stop             chan struct{}
	logger           zerolog.Logger
	externalEvmChain []common.Chain
	ethereumChain    common.Chain
}

func NewZetaSupplyChecker(cfg *config.Config, zetaClient *ZetaCoreBridge, logger zerolog.Logger) (ZetaSupplyChecker, error) {
	zetaSupplyChecker := ZetaSupplyChecker{
		stop:      make(chan struct{}),
		ticker:    NewDynamicTicker(fmt.Sprintf("ZETASupplyTicker"), 15),
		evmClient: make(map[int64]*ethclient.Client),
		logger: logger.With().
			Str("module", "ZetaSupplyChecker").
			Logger(),
		cfg:        cfg,
		zetaClient: zetaClient,
	}
	for _, evmConfig := range cfg.GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		client, err := ethclient.Dial(evmConfig.Endpoint)
		if err != nil {
			return zetaSupplyChecker, err
		}
		fmt.Println("Adding evmConfig.Chain.ChainId", evmConfig.Chain.ChainId)
		zetaSupplyChecker.evmClient[evmConfig.Chain.ChainId] = client
	}

	for chainID, _ := range zetaSupplyChecker.evmClient {
		chain := common.GetChainFromChainID(chainID)
		if chain.IsExternalChain() && common.IsEVMChain(chain.ChainId) && !common.IsEthereumChain(chain.ChainId) {
			zetaSupplyChecker.externalEvmChain = append(zetaSupplyChecker.externalEvmChain, *chain)
		}
		if common.IsEthereumChain(chain.ChainId) {
			zetaSupplyChecker.ethereumChain = *chain
		}
	}

	logger.Info().Msgf("zeta supply checker initialized , external chains : %v ,ethereum chain :%v", zetaSupplyChecker.externalEvmChain, zetaSupplyChecker.ethereumChain)

	return zetaSupplyChecker, nil
}
func (zs *ZetaSupplyChecker) Start() {
	defer zs.ticker.Stop()
	for {
		select {
		case <-zs.ticker.C():
			err := zs.CheckZetaTokenSupply()
			if err != nil {
				zs.logger.Error().Err(err).Msgf("ZetaSupplyChecker error")
			}
		case <-zs.stop:
			return
		}
	}
}

func (b *ZetaSupplyChecker) Stop() {
	b.logger.Info().Msgf("ZetaSupplyChecker is stopping")
	close(b.stop)
}

func (zs *ZetaSupplyChecker) CheckZetaTokenSupply() error {

	externalChainTotalSupply := big.NewInt(0)
	for _, chain := range zs.externalEvmChain {
		externalEvmChainConfig, ok := zs.cfg.GetEVMConfig(chain.ChainId)
		if !ok {
			return fmt.Errorf("externalEvmChainConfig not found for chain id %d", chain.ChainId)
		}
		zetaTokenAddressString := externalEvmChainConfig.ZetaTokenContractAddress
		zetaTokenAddress := ethcommon.HexToAddress(zetaTokenAddressString)
		zetatokenNonEth, err := FetchZetaZetaNonEthTokenContract(zetaTokenAddress, zs.evmClient[chain.ChainId])
		if err != nil {
			return err
		}
		totalSupply, err := zetatokenNonEth.TotalSupply(nil)
		if err != nil {
			return err
		}
		fmt.Println("Adding to external chain total supply", totalSupply.String())
		externalChainTotalSupply.Add(externalChainTotalSupply, totalSupply)
	}

	ethConfig, ok := zs.cfg.GetEVMConfig(zs.ethereumChain.ChainId)
	if !ok {
		return fmt.Errorf("eth config not found for chain id %d", zs.ethereumChain.ChainId)
	}
	ethConnectorAddressString := ethConfig.ConnectorContractAddress
	ethConnectorAddress := ethcommon.HexToAddress(ethConnectorAddressString)
	ethConnectorContract, err := FetchConnectorContractEth(ethConnectorAddress, zs.evmClient[zs.ethereumChain.ChainId])
	if err != nil {
		return err
	}

	ethLockedAmount, err := ethConnectorContract.GetLockedAmount(nil)
	if err != nil {
		return err
	}

	zetaInTransit := zs.GetAmountOfZetaInTransit(zs.externalEvmChain)
	zetaTokenSupplyOnNode, err := zs.zetaClient.GetZetaTokenSupplyOnNode()
	if err != nil {
		return err
	}
	genesisAmounts := zs.GetGenesistokenAmounts()
	abortedTxAmounts := zs.AbortedTxAmount()
	negativeAmounts := big.NewInt(0)
	negativeAmounts.Add(genesisAmounts, abortedTxAmounts)
	negativeAmounts.Add(negativeAmounts, zetaInTransit)

	positiveAmounts := big.NewInt(0)
	positiveAmounts.Add(externalChainTotalSupply, zetaTokenSupplyOnNode.BigInt())

	rhs := big.NewInt(0)
	lhs := ethLockedAmount
	rhs.Sub(positiveAmounts, negativeAmounts)
	copyZetaTokenSupplyOnNode := big.NewInt(0).Set(zetaTokenSupplyOnNode.BigInt())
	copyGenesisAmounts := big.NewInt(0).Set(genesisAmounts)
	nodeAmounts := big.NewInt(0).Sub(copyZetaTokenSupplyOnNode, copyGenesisAmounts)
	zs.logger.Info().Msgf("--------------------------------------------------------------------------------")
	zs.logger.Info().Msgf("aborted tx amounts : %s", abortedTxAmounts.String())
	zs.logger.Info().Msgf("zeta in transit : %s", zetaInTransit.String())
	zs.logger.Info().Msgf("external chain total supply : %s", externalChainTotalSupply.String())
	zs.logger.Info().Msgf("zeta token on node : %s", nodeAmounts.String())
	zs.logger.Info().Msgf("eth locked amount : %s", ethLockedAmount.String())
	if lhs.Cmp(rhs) != 0 {
		zs.logger.Error().Msgf("zeta supply mismatch, lhs : %s , rhs : %s", lhs.String(), rhs.String())
		return fmt.Errorf("zeta supply mismatch, lhs : %s , rhs : %s", lhs.String(), rhs.String())
	}

	zs.logger.Info().Msgf("zeta supply check passed, lhs : %s , rhs : %s", lhs.String(), rhs.String())
	zs.logger.Info().Msgf("--------------------------------------------------------------------------------")
	return nil
}

func (zs *ZetaSupplyChecker) GetGenesistokenAmounts() *big.Int {
	i, ok := big.NewInt(0).SetString("108402000200000000000000000", 10)
	if !ok {
		panic("error parsing genesis amount")
	}
	return i
}

func (zs *ZetaSupplyChecker) AbortedTxAmount() *big.Int {
	cctxList, err := zs.zetaClient.GetCctxByStatus(types.CctxStatus_Aborted)
	if err != nil {
		panic(err)
	}
	amount := sdkmath.ZeroUint()
	for _, cctx := range cctxList {
		amount = amount.Add(cctx.GetCurrentOutTxParam().Amount)
	}
	return amount.BigInt()
}

func (zs *ZetaSupplyChecker) GetAmountOfZetaInTransit(externalEvmchain []common.Chain) *big.Int {
	cctxs := zs.GetPendingCCTXNotAwaitingConfirmation(externalEvmchain)
	amount := sdkmath.ZeroUint()
	for _, cctx := range cctxs {
		amount = amount.Add(cctx.GetCurrentOutTxParam().Amount)
	}
	return amount.BigInt()
}
func (zs *ZetaSupplyChecker) GetPendingCCTXNotAwaitingConfirmation(externalEvmchain []common.Chain) []*types.CrossChainTx {
	ccTxNotAwaitngconfirmation := make([]*types.CrossChainTx, 0)
	for _, chain := range externalEvmchain {
		cctx, err := zs.zetaClient.GetAllPendingCctx(chain.ChainId)
		if err != nil {
			continue
		}
		nonceToCctxMap := make(map[uint64]*types.CrossChainTx)
		for _, c := range cctx {
			if c.GetInboundTxParams().CoinType == common.CoinType_Zeta {
				nonceToCctxMap[c.GetCurrentOutTxParam().OutboundTxTssNonce] = c
			}
		}
		trackers, err := zs.zetaClient.GetAllOutTxTrackerByChain(chain, Ascending)
		if err != nil {
			continue
		}
		for _, tracker := range trackers {
			if _, ok := nonceToCctxMap[tracker.Nonce]; !ok {
				ccTxNotAwaitngconfirmation = append(ccTxNotAwaitngconfirmation, nonceToCctxMap[tracker.Nonce])
			}
		}

	}
	return ccTxNotAwaitngconfirmation
}
