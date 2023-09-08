package keeper

import (
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

/**
 * Crosschain Mocks
 */

//go:generate mockery --name CrosschainAccountKeeper --filename account.go --case underscore --output ./mocks/crosschain
type CrosschainAccountKeeper interface {
	crosschaintypes.AccountKeeper
}

//go:generate mockery --name CrosschainBankKeeper --filename bank.go --case underscore --output ./mocks/crosschain
type CrosschainBankKeeper interface {
	crosschaintypes.BankKeeper
}

//go:generate mockery --name CrosschainStakingKeeper --filename staking.go --case underscore --output ./mocks/crosschain
type CrosschainStakingKeeper interface {
	crosschaintypes.StakingKeeper
}

//go:generate mockery --name CrosschainObserverKeeper --filename observer.go --case underscore --output ./mocks/crosschain
type CrosschainObserverKeeper interface {
	crosschaintypes.ZetaObserverKeeper
}

//go:generate mockery --name CrosschainFungibleKeeper --filename fungible.go --case underscore --output ./mocks/crosschain
type CrosschainFungibleKeeper interface {
	crosschaintypes.FungibleKeeper
}

/**
 * Fungible Mocks
 */

//go:generate mockery --name FungibleAccountKeeper --filename account.go --case underscore --output ./mocks/fungible
type FungibleAccountKeeper interface {
	fungibletypes.AccountKeeper
}

//go:generate mockery --name FungibleBankKeeper --filename bank.go --case underscore --output ./mocks/fungible
type FungibleBankKeeper interface {
	fungibletypes.BankKeeper
}

//go:generate mockery --name FungibleObserverKeeper --filename observer.go --case underscore --output ./mocks/fungible
type FungibleObserverKeeper interface {
	fungibletypes.ObserverKeeper
}

//go:generate mockery --name FungibleEVMKeeper --filename evm.go --case underscore --output ./mocks/fungible
type FungibleEVMKeeper interface {
	fungibletypes.EVMKeeper
}
