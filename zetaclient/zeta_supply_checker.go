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
	cfg        *config.Config
	evmClient  map[int64]*ethclient.Client
	zetaClient *ZetaCoreBridge
	ticker     *DynamicTicker
	stop       chan struct{}
	logger     zerolog.Logger
}

func NewZetaSupplyChecker(cfg *config.Config, zetaClient *ZetaCoreBridge, logger zerolog.Logger) (ZetaSupplyChecker, error) {
	zetaSupplyChecker := ZetaSupplyChecker{}
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
	zetaSupplyChecker.zetaClient = zetaClient
	zetaSupplyChecker.cfg = cfg
	zetaSupplyChecker.logger = logger.With().
		Str("module", "ZetaSupplyChecker").
		Logger()
	zetaSupplyChecker.stop = make(chan struct{})
	zetaSupplyChecker.ticker = NewDynamicTicker(fmt.Sprintf("ZETASupplyTicker"), 15)
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
	externalEvmChain := make([]common.Chain, 0)
	ethereumChain := common.Chain{}
	for chainID, _ := range zs.evmClient {
		chain := common.GetChainFromChainID(chainID)
		if chain.IsExternalChain() && common.IsEVMChain(chain.ChainId) && !common.IsEthereumChain(chain.ChainId) {
			externalEvmChain = append(externalEvmChain, *chain)
		}
		if common.IsEthereumChain(chain.ChainId) {
			ethereumChain = *chain
		}
	}
	if len(externalEvmChain) == 0 {
		return fmt.Errorf("no external chain found")
	}
	externalChainTotalSupply := big.NewInt(0)
	for _, chain := range externalEvmChain {

		zetaTokenAddressString := zs.cfg.EVMChainConfigs[chain.ChainId].CoreParams.ZetaTokenContractAddress
		zetaTokenAddress := ethcommon.HexToAddress(zetaTokenAddressString)
		zetatokenNonEth, err := FetchZetaZetaNonEthTokenContract(zetaTokenAddress, zs.evmClient[chain.ChainId])
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
	ethConnectorContract, err := FetchConnectorContractEth(ethConnectorAddress, zs.evmClient[ethereumChain.ChainId])
	if err != nil {
		return err
	}

	ethLockedAmount, err := ethConnectorContract.GetLockedAmount(nil)
	if err != nil {
		return err
	}

	zetaInTransit := zs.GetAmountOfZetaInTransit(externalEvmChain)
	zetaTokenSupplyOnNode, err := zs.zetaClient.GetZetaTokenSupplyOnNode()
	if err != nil {
		return err
	}
	genesisAmounts := big.NewInt(0)
	AbortedTxAmounts := big.NewInt(0)
	negativeAmounts := genesisAmounts.Add(genesisAmounts, AbortedTxAmounts).Add(genesisAmounts, zetaInTransit)
	positiveAmounts := externalChainTotalSupply.Add(externalChainTotalSupply, zetaTokenSupplyOnNode.BigInt())
	lhs := ethLockedAmount
	rhs := positiveAmounts.Sub(positiveAmounts, negativeAmounts)

	if lhs.Cmp(rhs) != 0 {
		zs.logger.Error().Msgf("zeta supply mismatch, lhs : %s , rhs : %s", lhs.String(), rhs.String())
		zs.logger.Error().Msgf("aborted tx amounts : %s", AbortedTxAmounts.String())
		zs.logger.Error().Msgf("genesis amounts : %s", genesisAmounts.String())
		zs.logger.Error().Msgf("zeta in transit : %s", zetaInTransit.String())
		zs.logger.Error().Msgf("external chain total supply : %s", externalChainTotalSupply.String())
		zs.logger.Error().Msgf("zeta token supply on node : %s", zetaTokenSupplyOnNode.String())
		zs.logger.Error().Msgf("eth locked amount : %s", ethLockedAmount.String())
		return fmt.Errorf("zeta supply mismatch, lhs : %s , rhs : %s", lhs.String(), rhs.String())
	}
	return nil
}

func (zs *ZetaSupplyChecker) GetGenesistokenAmounts() *big.Int {
	return nil
}

func (zs *ZetaSupplyChecker) AbortedTxAmount() *big.Int {
	return nil
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
