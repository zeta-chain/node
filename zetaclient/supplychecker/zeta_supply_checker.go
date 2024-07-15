// Package supplychecker provides functionalities to check the total supply of Zeta tokens
// Currently not used in the codebase
package supplychecker

import (
	"context"
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
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// ZetaSupplyChecker is a utility to check the total supply of Zeta tokens
type ZetaSupplyChecker struct {
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
	ctx context.Context,
	zetaClient *zetacore.Client,
	logger zerolog.Logger,
) (*ZetaSupplyChecker, error) {
	dynamicTicker, err := clienttypes.NewDynamicTicker("ZETASupplyTicker", 15)
	if err != nil {
		return nil, err
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	zetaSupplyChecker := &ZetaSupplyChecker{
		stop:      make(chan struct{}),
		ticker:    dynamicTicker,
		evmClient: make(map[int64]*ethclient.Client),
		logger: logger.With().
			Str("module", "ZetaSupplyChecker").
			Logger(),
		zetaClient: zetaClient,
	}

	for _, evmConfig := range app.Config().GetAllEVMConfigs() {
		if evmConfig.Chain.IsZetaChain() {
			continue
		}
		client, err := ethclient.Dial(evmConfig.Endpoint)
		if err != nil {
			return nil, err
		}

		zetaSupplyChecker.evmClient[evmConfig.Chain.ChainId] = client
	}

	for chainID := range zetaSupplyChecker.evmClient {
		chain, found := chains.GetChainFromChainID(chainID, app.GetAdditionalChains())
		if !found {
			return zetaSupplyChecker, fmt.Errorf("chain not found for chain id %d", chainID)
		}
		if chain.IsExternalChain() && chain.IsEVMChain() &&
			chain.Network != chains.Network_eth {
			zetaSupplyChecker.externalEvmChain = append(zetaSupplyChecker.externalEvmChain, chain)
		} else {
			zetaSupplyChecker.ethereumChain = chain
		}
	}

	balances, err := zetaSupplyChecker.zetaClient.GetGenesisSupply(ctx)
	if err != nil {
		return nil, err
	}

	tokensMintedAtBeginBlock, ok := sdkmath.NewIntFromString("200000000000000000")
	if !ok {
		return nil, fmt.Errorf("error parsing tokens minted at begin block")
	}

	zetaSupplyChecker.genesisSupply = balances.Add(tokensMintedAtBeginBlock)

	logger.Info().
		Msgf("zeta supply checker initialized , external chains : %v ,ethereum chain :%v", zetaSupplyChecker.externalEvmChain, zetaSupplyChecker.ethereumChain)

	return zetaSupplyChecker, nil
}

// Start starts the ZetaSupplyChecker
func (zs *ZetaSupplyChecker) Start(ctx context.Context) {
	defer zs.ticker.Stop()
	for {
		select {
		case <-zs.ticker.C():
			err := zs.CheckZetaTokenSupply(ctx)
			if err != nil {
				zs.logger.Error().Err(err).Msgf("ZetaSupplyChecker error")
			}
		case <-zs.stop:
			return
		}
	}
}

// Stop stops the ZetaSupplyChecker
func (zs *ZetaSupplyChecker) Stop() {
	zs.logger.Info().Msgf("ZetaSupplyChecker is stopping")
	close(zs.stop)
}

// CheckZetaTokenSupply checks the total supply of Zeta tokens
func (zs *ZetaSupplyChecker) CheckZetaTokenSupply(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	externalChainTotalSupply := sdkmath.ZeroInt()
	for _, chain := range zs.externalEvmChain {
		externalEvmChainParams, ok := app.GetEVMChainParams(chain.ChainId)
		if !ok {
			return fmt.Errorf("externalEvmChainParams not found for chain id %d", chain.ChainId)
		}

		zetaTokenAddressString := externalEvmChainParams.ZetaTokenContractAddress
		zetaTokenAddress := ethcommon.HexToAddress(zetaTokenAddressString)
		zetatokenNonEth, err := observer.FetchZetaTokenContract(zetaTokenAddress, zs.evmClient[chain.ChainId])
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

	evmChainParams, ok := app.GetEVMChainParams(zs.ethereumChain.ChainId)
	if !ok {
		return fmt.Errorf("eth config not found for chain id %d", zs.ethereumChain.ChainId)
	}

	ethConnectorAddressString := evmChainParams.ConnectorContractAddress
	ethConnectorAddress := ethcommon.HexToAddress(ethConnectorAddressString)
	ethConnectorContract, err := observer.FetchConnectorContractEth(
		ethConnectorAddress,
		zs.evmClient[zs.ethereumChain.ChainId],
	)
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

	zetaInTransit, err := zs.GetAmountOfZetaInTransit(ctx)
	if err != nil {
		return err
	}
	zetaTokenSupplyOnNode, err := zs.zetaClient.GetZetaTokenSupplyOnNode(ctx)
	if err != nil {
		return err
	}

	abortedAmount, err := zs.AbortedTxAmount(ctx)
	if err != nil {
		return err
	}

	ValidateZetaSupply(
		zs.logger,
		abortedAmount,
		zetaInTransit,
		zs.genesisSupply,
		externalChainTotalSupply,
		zetaTokenSupplyOnNode,
		ethLockedAmountInt,
	)

	return nil
}

// AbortedTxAmount returns the amount of Zeta tokens in aborted transactions
func (zs *ZetaSupplyChecker) AbortedTxAmount(ctx context.Context) (sdkmath.Int, error) {
	amount, err := zs.zetaClient.GetAbortedZetaAmount(ctx)
	if err != nil {
		return sdkmath.ZeroInt(), errors.Wrap(err, "error getting aborted zeta amount")
	}
	amountInt, ok := sdkmath.NewIntFromString(amount)
	if !ok {
		return sdkmath.ZeroInt(), errors.New("error parsing aborted zeta amount")
	}
	return amountInt, nil
}

// GetAmountOfZetaInTransit returns the amount of Zeta tokens in transit
func (zs *ZetaSupplyChecker) GetAmountOfZetaInTransit(ctx context.Context) (sdkmath.Int, error) {
	chainsToCheck := make([]chains.Chain, len(zs.externalEvmChain)+1)
	chainsToCheck = append(append(chainsToCheck, zs.externalEvmChain...), zs.ethereumChain)
	cctxs := zs.GetPendingCCTXInTransit(ctx, chainsToCheck)
	amount := sdkmath.ZeroUint()

	for _, cctx := range cctxs {
		amount = amount.Add(cctx.GetCurrentOutboundParam().Amount)
	}
	amountInt, ok := sdkmath.NewIntFromString(amount.String())
	if !ok {
		return sdkmath.ZeroInt(), fmt.Errorf("error parsing amount %s", amount.String())
	}

	return amountInt, nil
}

// GetPendingCCTXInTransit returns the pending CCTX in transit
func (zs *ZetaSupplyChecker) GetPendingCCTXInTransit(
	ctx context.Context,
	receivingChains []chains.Chain,
) []*types.CrossChainTx {
	cctxInTransit := make([]*types.CrossChainTx, 0)
	for _, chain := range receivingChains {
		cctx, _, err := zs.zetaClient.ListPendingCCTX(ctx, chain.ChainId)
		if err != nil {
			continue
		}
		nonceToCctxMap := make(map[uint64]*types.CrossChainTx)
		for _, c := range cctx {
			if c.InboundParams.CoinType == coin.CoinType_Zeta {
				nonceToCctxMap[c.GetCurrentOutboundParam().TssNonce] = c
			}
		}

		trackers, err := zs.zetaClient.GetAllOutboundTrackerByChain(ctx, chain.ChainId, interfaces.Ascending)
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
