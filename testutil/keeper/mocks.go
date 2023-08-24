package keeper

import (
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

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
