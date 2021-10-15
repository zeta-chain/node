package metacore

import (
	"github.com/Meta-Protocol/metacore/x/metacore/keeper"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the sendVoter
	for _, elem := range genState.SendVoterList {
		k.SetSendVoter(ctx, *elem)
	}

	// Set all the txoutConfirmation
	for _, elem := range genState.TxoutConfirmationList {
		k.SetTxoutConfirmation(ctx, *elem)
	}

	// Set all the txout
	for _, elem := range genState.TxoutList {
		k.SetTxout(ctx, *elem)
	}

	// Set txout count
	k.SetTxoutCount(ctx, genState.TxoutCount)

	// Set all the nodeAccount
	for _, elem := range genState.NodeAccountList {
		k.SetNodeAccount(ctx, *elem)
	}

	// Set all the txinVoter
	for _, elem := range genState.TxinVoterList {
		k.SetTxinVoter(ctx, *elem)
	}

	// Set all the txin
	for _, elem := range genState.TxinList {
		k.SetTxin(ctx, *elem)
	}

	// this line is used by starport scaffolding # ibc/genesis/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all sendVoter
	sendVoterList := k.GetAllSendVoter(ctx)
	for _, elem := range sendVoterList {
		elem := elem
		genesis.SendVoterList = append(genesis.SendVoterList, &elem)
	}

	// Get all txoutConfirmation
	txoutConfirmationList := k.GetAllTxoutConfirmation(ctx)
	for _, elem := range txoutConfirmationList {
		elem := elem
		genesis.TxoutConfirmationList = append(genesis.TxoutConfirmationList, &elem)
	}

	// Get all txout
	txoutList := k.GetAllTxout(ctx)
	for _, elem := range txoutList {
		elem := elem
		genesis.TxoutList = append(genesis.TxoutList, &elem)
	}

	// Set the current count
	genesis.TxoutCount = k.GetTxoutCount(ctx)

	// Get all nodeAccount
	nodeAccountList := k.GetAllNodeAccount(ctx)
	for _, elem := range nodeAccountList {
		elem := elem
		genesis.NodeAccountList = append(genesis.NodeAccountList, &elem)
	}

	// Get all txinVoter
	txinVoterList := k.GetAllTxinVoter(ctx)
	for _, elem := range txinVoterList {
		elem := elem
		genesis.TxinVoterList = append(genesis.TxinVoterList, &elem)
	}

	// Get all txin
	txinList := k.GetAllTxin(ctx)
	for _, elem := range txinList {
		elem := elem
		genesis.TxinList = append(genesis.TxinList, &elem)
	}

	// this line is used by starport scaffolding # ibc/genesis/export

	return genesis
}
