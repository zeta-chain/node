// Code generated by mockery v2.32.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	common "github.com/zeta-chain/zetacore/common"

	observertypes "github.com/zeta-chain/zetacore/x/observer/types"

	types "github.com/cosmos/cosmos-sdk/types"
)

// CrosschainObserverKeeper is an autogenerated mock type for the CrosschainObserverKeeper type
type CrosschainObserverKeeper struct {
	mock.Mock
}

// AddBallotToList provides a mock function with given fields: ctx, ballot
func (_m *CrosschainObserverKeeper) AddBallotToList(ctx types.Context, ballot observertypes.Ballot) {
	_m.Called(ctx, ballot)
}

// AddVoteToBallot provides a mock function with given fields: ctx, ballot, address, observationType
func (_m *CrosschainObserverKeeper) AddVoteToBallot(ctx types.Context, ballot observertypes.Ballot, address string, observationType observertypes.VoteType) (observertypes.Ballot, error) {
	ret := _m.Called(ctx, ballot, address, observationType)

	var r0 observertypes.Ballot
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, observertypes.Ballot, string, observertypes.VoteType) (observertypes.Ballot, error)); ok {
		return rf(ctx, ballot, address, observationType)
	}
	if rf, ok := ret.Get(0).(func(types.Context, observertypes.Ballot, string, observertypes.VoteType) observertypes.Ballot); ok {
		r0 = rf(ctx, ballot, address, observationType)
	} else {
		r0 = ret.Get(0).(observertypes.Ballot)
	}

	if rf, ok := ret.Get(1).(func(types.Context, observertypes.Ballot, string, observertypes.VoteType) error); ok {
		r1 = rf(ctx, ballot, address, observationType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckIfFinalizingVote provides a mock function with given fields: ctx, ballot
func (_m *CrosschainObserverKeeper) CheckIfFinalizingVote(ctx types.Context, ballot observertypes.Ballot) (observertypes.Ballot, bool) {
	ret := _m.Called(ctx, ballot)

	var r0 observertypes.Ballot
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, observertypes.Ballot) (observertypes.Ballot, bool)); ok {
		return rf(ctx, ballot)
	}
	if rf, ok := ret.Get(0).(func(types.Context, observertypes.Ballot) observertypes.Ballot); ok {
		r0 = rf(ctx, ballot)
	} else {
		r0 = ret.Get(0).(observertypes.Ballot)
	}

	if rf, ok := ret.Get(1).(func(types.Context, observertypes.Ballot) bool); ok {
		r1 = rf(ctx, ballot)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// FindBallot provides a mock function with given fields: ctx, index, chain, observationType
func (_m *CrosschainObserverKeeper) FindBallot(ctx types.Context, index string, chain *common.Chain, observationType observertypes.ObservationType) (observertypes.Ballot, bool, error) {
	ret := _m.Called(ctx, index, chain, observationType)

	var r0 observertypes.Ballot
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(types.Context, string, *common.Chain, observertypes.ObservationType) (observertypes.Ballot, bool, error)); ok {
		return rf(ctx, index, chain, observationType)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string, *common.Chain, observertypes.ObservationType) observertypes.Ballot); ok {
		r0 = rf(ctx, index, chain, observationType)
	} else {
		r0 = ret.Get(0).(observertypes.Ballot)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string, *common.Chain, observertypes.ObservationType) bool); ok {
		r1 = rf(ctx, index, chain, observationType)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(types.Context, string, *common.Chain, observertypes.ObservationType) error); ok {
		r2 = rf(ctx, index, chain, observationType)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetAllBallots provides a mock function with given fields: ctx
func (_m *CrosschainObserverKeeper) GetAllBallots(ctx types.Context) []*observertypes.Ballot {
	ret := _m.Called(ctx)

	var r0 []*observertypes.Ballot
	if rf, ok := ret.Get(0).(func(types.Context) []*observertypes.Ballot); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*observertypes.Ballot)
		}
	}

	return r0
}

// GetAllNodeAccount provides a mock function with given fields: ctx
func (_m *CrosschainObserverKeeper) GetAllNodeAccount(ctx types.Context) []observertypes.NodeAccount {
	ret := _m.Called(ctx)

	var r0 []observertypes.NodeAccount
	if rf, ok := ret.Get(0).(func(types.Context) []observertypes.NodeAccount); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]observertypes.NodeAccount)
		}
	}

	return r0
}

// GetAllObserverMappers provides a mock function with given fields: ctx
func (_m *CrosschainObserverKeeper) GetAllObserverMappers(ctx types.Context) []*observertypes.ObserverMapper {
	ret := _m.Called(ctx)

	var r0 []*observertypes.ObserverMapper
	if rf, ok := ret.Get(0).(func(types.Context) []*observertypes.ObserverMapper); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*observertypes.ObserverMapper)
		}
	}

	return r0
}

// GetBallot provides a mock function with given fields: ctx, index
func (_m *CrosschainObserverKeeper) GetBallot(ctx types.Context, index string) (observertypes.Ballot, bool) {
	ret := _m.Called(ctx, index)

	var r0 observertypes.Ballot
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (observertypes.Ballot, bool)); ok {
		return rf(ctx, index)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) observertypes.Ballot); ok {
		r0 = rf(ctx, index)
	} else {
		r0 = ret.Get(0).(observertypes.Ballot)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, index)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetBlockHeader provides a mock function with given fields: ctx, hash
func (_m *CrosschainObserverKeeper) GetBlockHeader(ctx types.Context, hash []byte) (observertypes.BlockHeader, bool) {
	ret := _m.Called(ctx, hash)

	var r0 observertypes.BlockHeader
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, []byte) (observertypes.BlockHeader, bool)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(types.Context, []byte) observertypes.BlockHeader); ok {
		r0 = rf(ctx, hash)
	} else {
		r0 = ret.Get(0).(observertypes.BlockHeader)
	}

	if rf, ok := ret.Get(1).(func(types.Context, []byte) bool); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetCoreParamsByChainID provides a mock function with given fields: ctx, chainID
func (_m *CrosschainObserverKeeper) GetCoreParamsByChainID(ctx types.Context, chainID int64) (*observertypes.CoreParams, bool) {
	ret := _m.Called(ctx, chainID)

	var r0 *observertypes.CoreParams
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, int64) (*observertypes.CoreParams, bool)); ok {
		return rf(ctx, chainID)
	}
	if rf, ok := ret.Get(0).(func(types.Context, int64) *observertypes.CoreParams); ok {
		r0 = rf(ctx, chainID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*observertypes.CoreParams)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Context, int64) bool); ok {
		r1 = rf(ctx, chainID)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetKeygen provides a mock function with given fields: ctx
func (_m *CrosschainObserverKeeper) GetKeygen(ctx types.Context) (observertypes.Keygen, bool) {
	ret := _m.Called(ctx)

	var r0 observertypes.Keygen
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context) (observertypes.Keygen, bool)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(types.Context) observertypes.Keygen); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(observertypes.Keygen)
	}

	if rf, ok := ret.Get(1).(func(types.Context) bool); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetNodeAccount provides a mock function with given fields: ctx, address
func (_m *CrosschainObserverKeeper) GetNodeAccount(ctx types.Context, address string) (observertypes.NodeAccount, bool) {
	ret := _m.Called(ctx, address)

	var r0 observertypes.NodeAccount
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, string) (observertypes.NodeAccount, bool)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string) observertypes.NodeAccount); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Get(0).(observertypes.NodeAccount)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string) bool); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetObserverMapper provides a mock function with given fields: ctx, chain
func (_m *CrosschainObserverKeeper) GetObserverMapper(ctx types.Context, chain *common.Chain) (observertypes.ObserverMapper, bool) {
	ret := _m.Called(ctx, chain)

	var r0 observertypes.ObserverMapper
	var r1 bool
	if rf, ok := ret.Get(0).(func(types.Context, *common.Chain) (observertypes.ObserverMapper, bool)); ok {
		return rf(ctx, chain)
	}
	if rf, ok := ret.Get(0).(func(types.Context, *common.Chain) observertypes.ObserverMapper); ok {
		r0 = rf(ctx, chain)
	} else {
		r0 = ret.Get(0).(observertypes.ObserverMapper)
	}

	if rf, ok := ret.Get(1).(func(types.Context, *common.Chain) bool); ok {
		r1 = rf(ctx, chain)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetParams provides a mock function with given fields: ctx
func (_m *CrosschainObserverKeeper) GetParams(ctx types.Context) observertypes.Params {
	ret := _m.Called(ctx)

	var r0 observertypes.Params
	if rf, ok := ret.Get(0).(func(types.Context) observertypes.Params); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(observertypes.Params)
	}

	return r0
}

// IsAuthorized provides a mock function with given fields: ctx, address, chain
func (_m *CrosschainObserverKeeper) IsAuthorized(ctx types.Context, address string, chain *common.Chain) (bool, error) {
	ret := _m.Called(ctx, address, chain)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, string, *common.Chain) (bool, error)); ok {
		return rf(ctx, address, chain)
	}
	if rf, ok := ret.Get(0).(func(types.Context, string, *common.Chain) bool); ok {
		r0 = rf(ctx, address, chain)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(types.Context, string, *common.Chain) error); ok {
		r1 = rf(ctx, address, chain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsInboundEnabled provides a mock function with given fields: ctx
func (_m *CrosschainObserverKeeper) IsInboundEnabled(ctx types.Context) bool {
	ret := _m.Called(ctx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// IsOutboundAllowed provides a mock function with given fields: ctx
func (_m *CrosschainObserverKeeper) IsOutboundAllowed(ctx types.Context) bool {
	ret := _m.Called(ctx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(types.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetBallot provides a mock function with given fields: ctx, ballot
func (_m *CrosschainObserverKeeper) SetBallot(ctx types.Context, ballot *observertypes.Ballot) {
	_m.Called(ctx, ballot)
}

// SetKeygen provides a mock function with given fields: ctx, keygen
func (_m *CrosschainObserverKeeper) SetKeygen(ctx types.Context, keygen observertypes.Keygen) {
	_m.Called(ctx, keygen)
}

// SetLastObserverCount provides a mock function with given fields: ctx, lbc
func (_m *CrosschainObserverKeeper) SetLastObserverCount(ctx types.Context, lbc *observertypes.LastObserverCount) {
	_m.Called(ctx, lbc)
}

// SetNodeAccount provides a mock function with given fields: ctx, nodeAccount
func (_m *CrosschainObserverKeeper) SetNodeAccount(ctx types.Context, nodeAccount observertypes.NodeAccount) {
	_m.Called(ctx, nodeAccount)
}

// SetObserverMapper provides a mock function with given fields: ctx, om
func (_m *CrosschainObserverKeeper) SetObserverMapper(ctx types.Context, om *observertypes.ObserverMapper) {
	_m.Called(ctx, om)
}

// SetPermissionFlags provides a mock function with given fields: ctx, permissionFlags
func (_m *CrosschainObserverKeeper) SetPermissionFlags(ctx types.Context, permissionFlags observertypes.PermissionFlags) {
	_m.Called(ctx, permissionFlags)
}

// NewCrosschainObserverKeeper creates a new instance of CrosschainObserverKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCrosschainObserverKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *CrosschainObserverKeeper {
	mock := &CrosschainObserverKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
