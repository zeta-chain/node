package zetacore

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/keeper"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the zetaConversionRate
	for _, elem := range genState.ZetaConversionRateList {
		k.SetZetaConversionRate(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	// Set if defined
	if genState.Keygen != nil {
		k.SetKeygen(ctx, *genState.Keygen)
	}

	// Set all the tSSVoter
	for _, elem := range genState.TSSVoterList {
		k.SetTSSVoter(ctx, *elem)
	}

	// Set all the tSS
	for _, elem := range genState.TSSList {
		k.SetTSS(ctx, *elem)
	}

	// Set all the inTx
	for _, elem := range genState.InTxList {
		k.SetInTx(ctx, *elem)
	}

	// Set if defined
	if genState.TxList != nil {
		k.SetTxList(ctx, *genState.TxList)
	}

	// Set all the gasBalance
	for _, elem := range genState.GasBalanceList {
		k.SetGasBalance(ctx, *elem)
	}

	// Set all the gasPrice
	for _, elem := range genState.GasPriceList {
		k.SetGasPrice(ctx, *elem)
	}

	// Set all the chainNonces
	for _, elem := range genState.ChainNoncesList {
		k.SetChainNonces(ctx, *elem)
	}

	// Set all the lastBlockHeight
	for _, elem := range genState.LastBlockHeightList {
		k.SetLastBlockHeight(ctx, *elem)
	}

	// Set all the receive
	for _, elem := range genState.ReceiveList {
		k.SetReceive(ctx, *elem)
	}

	// Set all the send
	for _, elem := range genState.SendList {
		k.SetSend(ctx, *elem)
	}

	// Set all the nodeAccount
	for _, elem := range genState.NodeAccountList {
		k.SetNodeAccount(ctx, *elem)
	}

	// this line is used by starport scaffolding # ibc/genesis/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.ZetaConversionRateList = k.GetAllZetaConversionRate(ctx)
	// this line is used by starport scaffolding # genesis/module/export
	// Get all keygen
	keygen, found := k.GetKeygen(ctx)
	if found {
		genesis.Keygen = &keygen
	}

	// Get all tSSVoter
	tSSVoterList := k.GetAllTSSVoter(ctx)
	for _, elem := range tSSVoterList {
		elem := elem
		genesis.TSSVoterList = append(genesis.TSSVoterList, &elem)
	}

	// Get all tSS
	tSSList := k.GetAllTSS(ctx)
	for _, elem := range tSSList {
		elem := elem
		genesis.TSSList = append(genesis.TSSList, &elem)
	}

	// Get all inTx
	inTxList := k.GetAllInTx(ctx)
	for _, elem := range inTxList {
		elem := elem
		genesis.InTxList = append(genesis.InTxList, &elem)
	}

	// Get all txList
	txList, found := k.GetTxList(ctx)
	if found {
		genesis.TxList = &txList
	}

	// Get all gasBalance
	gasBalanceList := k.GetAllGasBalance(ctx)
	for _, elem := range gasBalanceList {
		elem := elem
		genesis.GasBalanceList = append(genesis.GasBalanceList, &elem)
	}

	// Get all gasPrice
	gasPriceList := k.GetAllGasPrice(ctx)
	for _, elem := range gasPriceList {
		elem := elem
		genesis.GasPriceList = append(genesis.GasPriceList, &elem)
	}

	// Get all chainNonces
	chainNoncesList := k.GetAllChainNonces(ctx)
	for _, elem := range chainNoncesList {
		elem := elem
		genesis.ChainNoncesList = append(genesis.ChainNoncesList, &elem)
	}

	// Get all lastBlockHeight
	lastBlockHeightList := k.GetAllLastBlockHeight(ctx)
	for _, elem := range lastBlockHeightList {
		elem := elem
		genesis.LastBlockHeightList = append(genesis.LastBlockHeightList, &elem)
	}

	// Get all receive
	receiveList := k.GetAllReceive(ctx)
	for _, elem := range receiveList {
		elem := elem
		genesis.ReceiveList = append(genesis.ReceiveList, &elem)
	}

	// Get all send
	sendList := k.GetAllSend(ctx)
	for _, elem := range sendList {
		elem := elem
		genesis.SendList = append(genesis.SendList, &elem)
	}

	// Get all nodeAccount
	nodeAccountList := k.GetAllNodeAccount(ctx)
	for _, elem := range nodeAccountList {
		elem := elem
		genesis.NodeAccountList = append(genesis.NodeAccountList, &elem)
	}

	// this line is used by starport scaffolding # ibc/genesis/export

	return genesis
}
