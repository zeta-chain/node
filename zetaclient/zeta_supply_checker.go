package zetaclient

import (
	"fmt"

	"github.com/pkg/errors"

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
	genesisSupply    sdkmath.Int
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
		zetaSupplyChecker.evmClient[evmConfig.Chain.ChainId] = client
	}

	for chainID := range zetaSupplyChecker.evmClient {
		chain := common.GetChainFromChainID(chainID)
		if chain.IsExternalChain() && common.IsEVMChain(chain.ChainId) && !common.IsEthereumChain(chain.ChainId) {
			zetaSupplyChecker.externalEvmChain = append(zetaSupplyChecker.externalEvmChain, *chain)
		}
		if common.IsEthereumChain(chain.ChainId) {
			zetaSupplyChecker.ethereumChain = *chain
		}
	}
	balances, err := zetaSupplyChecker.zetaClient.GetGenesisSupply()
	if err != nil {
		return zetaSupplyChecker, err
	}
	tokensMintedAtBeginBlock, ok := sdkmath.NewIntFromString("200000000000000000")
	if !ok {
		return zetaSupplyChecker, fmt.Errorf("error parsing tokens minted at begin block")
	}
	zetaSupplyChecker.genesisSupply = balances.Add(tokensMintedAtBeginBlock)

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

func (zs *ZetaSupplyChecker) Stop() {
	zs.logger.Info().Msgf("ZetaSupplyChecker is stopping")
	close(zs.stop)
}

func (zs *ZetaSupplyChecker) CheckZetaTokenSupply() error {

	externalChainTotalSupply := sdkmath.ZeroInt()
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
		totalSupplyInt, ok := sdkmath.NewIntFromString(totalSupply.String())
		if !ok {
			zs.logger.Error().Msgf("error parsing total supply for chain %d", chain.ChainId)
			continue
		}
		externalChainTotalSupply = externalChainTotalSupply.Add(totalSupplyInt)
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
	ethLockedAmountInt, ok := sdkmath.NewIntFromString(ethLockedAmount.String())
	if !ok {
		return fmt.Errorf("error parsing eth locked amount")
	}

	zetaInTransit := zs.GetAmountOfZetaInTransit()
	zetaTokenSupplyOnNode, err := zs.zetaClient.GetZetaTokenSupplyOnNode()
	if err != nil {
		return err
	}
	abortedAmount, err := zs.AbortedTxAmount()
	if err != nil {
		return err
	}
	ValidateZetaSupply(zs.logger, abortedAmount, zetaInTransit, zs.genesisSupply, externalChainTotalSupply, zetaTokenSupplyOnNode, ethLockedAmountInt)
	return nil
}

func ValidateZetaSupply(logger zerolog.Logger, abortedTxAmounts, zetaInTransit, genesisAmounts, externalChainTotalSupply, zetaTokenSupplyOnNode, ethLockedAmount sdkmath.Int) bool {
	lhs := ethLockedAmount.Sub(abortedTxAmounts)
	rhs := zetaTokenSupplyOnNode.Add(zetaInTransit).Add(externalChainTotalSupply).Sub(genesisAmounts)

	copyZetaTokenSupplyOnNode := zetaTokenSupplyOnNode
	copyGenesisAmounts := genesisAmounts
	nodeAmounts := copyZetaTokenSupplyOnNode.Sub(copyGenesisAmounts)
	logger.Info().Msgf("--------------------------------------------------------------------------------")
	logger.Info().Msgf("aborted tx amounts : %s", abortedTxAmounts.String())
	logger.Info().Msgf("zeta in transit : %s", zetaInTransit.String())
	logger.Info().Msgf("external chain total supply : %s", externalChainTotalSupply.String())
	logger.Info().Msgf("effective zeta supply on node : %s", nodeAmounts.String())
	logger.Info().Msgf("zeta token supply on node : %s", zetaTokenSupplyOnNode.String())
	logger.Info().Msgf("genesis amounts : %s", genesisAmounts.String())
	logger.Info().Msgf("eth locked amount : %s", ethLockedAmount.String())
	if !lhs.Equal(rhs) {
		logger.Error().Msgf("zeta supply mismatch, lhs : %s , rhs : %s", lhs.String(), rhs.String())
		return false
	}
	logger.Info().Msgf("zeta supply check passed, lhs : %s , rhs : %s", lhs.String(), rhs.String())
	logger.Info().Msgf("--------------------------------------------------------------------------------")
	return true
}

// TODO : Add genesis state for Aborted amount in zeta
// TODO : Add tests for keeper functions
// TODO : Add cli commands for querying aborted amount
func (zs *ZetaSupplyChecker) AbortedTxAmount() (sdkmath.Int, error) {
	amount, err := zs.zetaClient.GetAbortedZetaAmount()
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "error getting aborted zeta amount")
	}
	amountInt, ok := sdkmath.NewIntFromString(amount.String())
	if !ok {
		return sdkmath.ZeroInt(), errors.New("error parsing aborted zeta amount")
	}
	return amountInt, nil
}

func (zs *ZetaSupplyChecker) GetAmountOfZetaInTransit() sdkmath.Int {
	chainsToCheck := make([]common.Chain, len(zs.externalEvmChain)+1)
	chainsToCheck = append(append(chainsToCheck, zs.externalEvmChain...), zs.ethereumChain)
	cctxs := zs.GetPendingCCTXInTransit(chainsToCheck)
	amount := sdkmath.ZeroUint()
	for _, cctx := range cctxs {
		amount = amount.Add(cctx.GetCurrentOutTxParam().Amount)
	}
	amountInt, ok := sdkmath.NewIntFromString(amount.String())
	if !ok {
		panic("error parsing amount")
	}
	return amountInt
}
func (zs *ZetaSupplyChecker) GetPendingCCTXInTransit(receivingChains []common.Chain) []*types.CrossChainTx {
	cctxInTransit := make([]*types.CrossChainTx, 0)
	for _, chain := range receivingChains {
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
			zs.logger.Info().Msgf("tracker exists for nonce: %d , removing from supply checks", tracker.Nonce)
			delete(nonceToCctxMap, tracker.Nonce)
		}
		for _, c := range nonceToCctxMap {
			if c != nil {
				cctxInTransit = append(cctxInTransit, c)
			}
		}
	}
	return cctxInTransit
}
