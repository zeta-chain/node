package supplychecker

import (
	"fmt"

	appcontext "github.com/zeta-chain/zetacore/zetaclient/app_context"
	"github.com/zeta-chain/zetacore/zetaclient/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"

	"github.com/zeta-chain/zetacore/zetaclient/evm"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

type ZetaSupplyChecker struct {
	coreContext      *corecontext.ZetaCoreContext
	evmClient        map[int64]*ethclient.Client
	zetaClient       *zetabridge.ZetaCoreBridge
	ticker           *clienttypes.DynamicTicker
	stop             chan struct{}
	logger           zerolog.Logger
	externalEvmChain []common.Chain
	ethereumChain    common.Chain
	genesisSupply    sdkmath.Int
}

func NewZetaSupplyChecker(appContext *appcontext.AppContext, zetaClient *zetabridge.ZetaCoreBridge, logger zerolog.Logger) (ZetaSupplyChecker, error) {
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
		externalEvmChainParams, ok := zs.coreContext.GetEVMChainParams(chain.ChainId)
		if !ok {
			return fmt.Errorf("externalEvmChainParams not found for chain id %d", chain.ChainId)
		}

		zetaTokenAddressString := externalEvmChainParams.ZetaTokenContractAddress
		zetaTokenAddress := ethcommon.HexToAddress(zetaTokenAddressString)
		zetatokenNonEth, err := evm.FetchZetaZetaNonEthTokenContract(zetaTokenAddress, zs.evmClient[chain.ChainId])
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

	ethConfig, ok := zs.coreContext.GetEVMChainParams(zs.ethereumChain.ChainId)
	if !ok {
		return fmt.Errorf("eth config not found for chain id %d", zs.ethereumChain.ChainId)
	}
	ethConnectorAddressString := ethConfig.ConnectorContractAddress
	ethConnectorAddress := ethcommon.HexToAddress(ethConnectorAddressString)
	ethConnectorContract, err := evm.FetchConnectorContractEth(ethConnectorAddress, zs.evmClient[zs.ethereumChain.ChainId])
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

type ZetaSupplyCheckLogs struct {
	Logger                   zerolog.Logger
	AbortedTxAmounts         sdkmath.Int `json:"aborted_tx_amounts"`
	ZetaInTransit            sdkmath.Int `json:"zeta_in_transit"`
	ExternalChainTotalSupply sdkmath.Int `json:"external_chain_total_supply"`
	ZetaTokenSupplyOnNode    sdkmath.Int `json:"zeta_token_supply_on_node"`
	EthLockedAmount          sdkmath.Int `json:"eth_locked_amount"`
	NodeAmounts              sdkmath.Int `json:"node_amounts"`
	LHS                      sdkmath.Int `json:"LHS"`
	RHS                      sdkmath.Int `json:"RHS"`
	SupplyCheckSuccess       bool        `json:"supply_check_success"`
}

func (z ZetaSupplyCheckLogs) LogOutput() {
	output, err := bitcoin.PrettyPrintStruct(z)
	if err != nil {
		z.Logger.Error().Err(err).Msgf("error pretty printing struct")
	}
	z.Logger.Info().Msgf(output)
}

func ValidateZetaSupply(logger zerolog.Logger, abortedTxAmounts, zetaInTransit, genesisAmounts, externalChainTotalSupply, zetaTokenSupplyOnNode, ethLockedAmount sdkmath.Int) bool {
	lhs := ethLockedAmount.Sub(abortedTxAmounts)
	rhs := zetaTokenSupplyOnNode.Add(zetaInTransit).Add(externalChainTotalSupply).Sub(genesisAmounts)

	copyZetaTokenSupplyOnNode := zetaTokenSupplyOnNode
	copyGenesisAmounts := genesisAmounts
	nodeAmounts := copyZetaTokenSupplyOnNode.Sub(copyGenesisAmounts)
	logs := ZetaSupplyCheckLogs{
		Logger:                   logger,
		AbortedTxAmounts:         abortedTxAmounts,
		ZetaInTransit:            zetaInTransit,
		ExternalChainTotalSupply: externalChainTotalSupply,
		NodeAmounts:              nodeAmounts,
		ZetaTokenSupplyOnNode:    zetaTokenSupplyOnNode,
		EthLockedAmount:          ethLockedAmount,
		LHS:                      lhs,
		RHS:                      rhs,
	}
	defer logs.LogOutput()
	if !lhs.Equal(rhs) {
		logs.SupplyCheckSuccess = false
		return false
	}
	logs.SupplyCheckSuccess = true
	return true
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
		cctx, _, err := zs.zetaClient.ListPendingCctx(chain.ChainId)
		if err != nil {
			continue
		}
		nonceToCctxMap := make(map[uint64]*types.CrossChainTx)
		for _, c := range cctx {
			if c.GetInboundTxParams().CoinType == common.CoinType_Zeta {
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
