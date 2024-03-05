package mocks

import (
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

/**
 * Crosschain Mocks
 */

//go:generate mockery --name CrosschainAccountKeeper --filename account.go --case underscore --output ./crosschain
type CrosschainAccountKeeper interface {
	crosschaintypes.AccountKeeper
}

//go:generate mockery --name CrosschainBankKeeper --filename bank.go --case underscore --output ./crosschain
type CrosschainBankKeeper interface {
	crosschaintypes.BankKeeper
}

//go:generate mockery --name CrosschainStakingKeeper --filename staking.go --case underscore --output ./crosschain
type CrosschainStakingKeeper interface {
	crosschaintypes.StakingKeeper
}

//go:generate mockery --name CrosschainObserverKeeper --filename observer.go --case underscore --output ./crosschain
type CrosschainObserverKeeper interface {
	crosschaintypes.ObserverKeeper
}

//go:generate mockery --name CrosschainFungibleKeeper --filename fungible.go --case underscore --output ./crosschain
type CrosschainFungibleKeeper interface {
	crosschaintypes.FungibleKeeper
}

//go:generate mockery --name CrosschainAuthorityKeeper --filename authority.go --case underscore --output ./crosschain
type CrosschainAuthorityKeeper interface {
	crosschaintypes.AuthorityKeeper
}

/**
 * Fungible Mocks
 */

//go:generate mockery --name FungibleAccountKeeper --filename account.go --case underscore --output ./fungible
type FungibleAccountKeeper interface {
	fungibletypes.AccountKeeper
}

//go:generate mockery --name FungibleBankKeeper --filename bank.go --case underscore --output ./fungible
type FungibleBankKeeper interface {
	fungibletypes.BankKeeper
}

//go:generate mockery --name FungibleObserverKeeper --filename observer.go --case underscore --output ./fungible
type FungibleObserverKeeper interface {
	fungibletypes.ObserverKeeper
}

//go:generate mockery --name FungibleEVMKeeper --filename evm.go --case underscore --output ./fungible
type FungibleEVMKeeper interface {
	fungibletypes.EVMKeeper
}

//go:generate mockery --name FungibleAuthorityKeeper --filename authority.go --case underscore --output ./fungible
type FungibleAuthorityKeeper interface {
	fungibletypes.AuthorityKeeper
}

/**
 * Emissions Mocks
 */

//go:generate mockery --name EmissionAccountKeeper --filename account.go --case underscore --output ./emissions
type EmissionAccountKeeper interface {
	emissionstypes.AccountKeeper
}

//go:generate mockery --name EmissionBankKeeper --filename bank.go --case underscore --output ./emissions
type EmissionBankKeeper interface {
	emissionstypes.BankKeeper
}

//go:generate mockery --name EmissionStakingKeeper --filename staking.go --case underscore --output ./emissions
type EmissionStakingKeeper interface {
	emissionstypes.StakingKeeper
}

//go:generate mockery --name EmissionObserverKeeper --filename observer.go --case underscore --output ./emissions
type EmissionObserverKeeper interface {
	emissionstypes.ObserverKeeper
}

/**
 * Observer Mocks
 */

//go:generate mockery --name ObserverStakingKeeper --filename staking.go --case underscore --output ./observer
type ObserverStakingKeeper interface {
	observertypes.StakingKeeper
}

//go:generate mockery --name ObserverSlashingKeeper --filename slashing.go --case underscore --output ./observer
type ObserverSlashingKeeper interface {
	observertypes.SlashingKeeper
}

//go:generate mockery --name ObserverAuthorityKeeper --filename authority.go --case underscore --output ./observer
type ObserverAuthorityKeeper interface {
	observertypes.AuthorityKeeper
}
