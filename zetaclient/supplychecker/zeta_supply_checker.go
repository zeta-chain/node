package supplychecker

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// ZetaSupplyChecker is a utility to check the total supply of Zeta tokens
type ZetaSupplyChecker struct {
	coreContext      *context.ZetaCoreContext
	evmClient        map[int64]*ethclient.Client
	zetaClient       *zetacore.Client
	ticker           *clienttypes.DynamicTicker
	stop             chan struct{}
	logger           zerolog.Logger
	externalEvmChain []chains.Chain
	ethereumChain    chains.Chain
	genesisSupply    sdkmath.Int
}

// NewZetaSupplyChecker creates a new ZetaSupplyChecker
func NewZetaSupplyChecker(
	appContext *context.AppContext,
	zetaClient *zetacore.Client,
	logger zerolog.Logger,
) (ZetaSupplyChecker, error) {
	dynamicTicker, err := clienttypes.NewDynamicTicker("ZETASupplyTicker", 15)
	if err != nil {
		return ZetaSupplyChecker{}, err
	}

	zetaSupplyChecker := ZetaSupplyChecker{
		stop:      make(chan struct{}),
		ticker:    dynamicTicker,
		evmClient: make(map[int64]*ethclient.Client),
		logger: logger.With().
			Str("module", "ZetaSupplyChecker").
			Logger(),
		coreContext: appContext.ZetaCoreContext(),
		zetaClient:  zetaClient,
	}

	for _, evmConfig := range appContext.Config().GetAllEVMConfigs() {
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
		chain := chains.GetChainFromChainID(chainID)
		if chain.IsExternalChain() && chains.IsEVMChain(chain.ChainId) && !chains.IsEthereumChain(chain.ChainId) {
			zetaSupplyChecker.externalEvmChain = append(zetaSupplyChecker.externalEvmChain, *chain)
		}
		if chains.IsEthereumChain(chain.ChainId) {
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
		externalEvmChainParams, ok := zs.coreContext.GetEVMChainParams(chain.ChainId)
		if !ok {
			return fmt.Errorf("externalEvmChainParams not found for chain id %d", chain.ChainId)
		}

		zetaTokenAddressString := externalEvmChainParams.ZetaTokenContractAddress
		zetaTokenAddress := ethcommon.HexToAddress(zetaTokenAddressString)
		zetatokenNonEth, err := observer.FetchZetaZetaNonEthTokenContract(zetaTokenAddress, zs.evmClient[chain.ChainId])
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

	evmChainParams, ok := zs.coreContext.GetEVMChainParams(zs.ethereumChain.ChainId)
	if !ok {
		return fmt.Errorf("eth config not found for chain id %d", zs.ethereumChain.ChainId)
	}

	ethConnectorAddressString := evmChainParams.ConnectorContractAddress
	ethConnectorAddress := ethcommon.HexToAddress(ethConnectorAddressString)
	ethConnectorContract, err := observer.FetchConnectorContractEth(ethConnectorAddress, zs.evmClient[zs.ethereumChain.ChainId])
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

func (zs *ZetaSupplyChecker) AbortedTxAmount() (sdkmath.Int, error) {
	amount, err := zs.zetaClient.GetAbortedZetaAmount()
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "error getting aborted zeta amount")
	}
	amountInt, ok := sdkmath.NewIntFromString(amount)
	if !ok {
		return sdkmath.ZeroInt(), errors.New("error parsing aborted zeta amount")
	}
	return amountInt, nil
}

func (zs *ZetaSupplyChecker) GetAmountOfZetaInTransit() sdkmath.Int {
	chainsToCheck := make([]chains.Chain, len(zs.externalEvmChain)+1)
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

func (zs *ZetaSupplyChecker) GetPendingCCTXInTransit(receivingChains []chains.Chain) []*types.CrossChainTx {
	cctxInTransit := make([]*types.CrossChainTx, 0)
	for _, chain := range receivingChains {
		cctx, _, err := zs.zetaClient.ListPendingCctx(chain.ChainId)
		if err != nil {
			continue
		}
		nonceToCctxMap := make(map[uint64]*types.CrossChainTx)
		for _, c := range cctx {
			if c.InboundTxParams.CoinType == coin.CoinType_Zeta {
				nonceToCctxMap[c.GetCurrentOutTxParam().OutboundTxTssNonce] = c
			}
		}

		trackers, err := zs.zetaClient.GetAllOutTxTrackerByChain(chain.ChainId, interfaces.Ascending)
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
