// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// EmissionBankKeeper is an autogenerated mock type for the EmissionBankKeeper type
type EmissionBankKeeper struct {
	mock.Mock
}

// GetBalance provides a mock function with given fields: ctx, addr, denom
func (_m *EmissionBankKeeper) GetBalance(ctx types.Context, addr types.AccAddress, denom string) types.Coin {
	ret := _m.Called(ctx, addr, denom)

	if len(ret) == 0 {
		panic("no return value specified for GetBalance")
	}

	var r0 types.Coin
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, string) types.Coin); ok {
		r0 = rf(ctx, addr, denom)
	} else {
		r0 = ret.Get(0).(types.Coin)
	}

	return r0
}

// SendCoinsFromAccountToModule provides a mock function with given fields: ctx, senderAddr, recipientModule, amt
func (_m *EmissionBankKeeper) SendCoinsFromAccountToModule(ctx types.Context, senderAddr types.AccAddress, recipientModule string, amt types.Coins) error {
	ret := _m.Called(ctx, senderAddr, recipientModule, amt)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinsFromAccountToModule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress, string, types.Coins) error); ok {
		r0 = rf(ctx, senderAddr, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinsFromModuleToAccount provides a mock function with given fields: ctx, senderModule, recipientAddr, amt
func (_m *EmissionBankKeeper) SendCoinsFromModuleToAccount(ctx types.Context, senderModule string, recipientAddr types.AccAddress, amt types.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientAddr, amt)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinsFromModuleToAccount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, types.AccAddress, types.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientAddr, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinsFromModuleToModule provides a mock function with given fields: ctx, senderModule, recipientModule, amt
func (_m *EmissionBankKeeper) SendCoinsFromModuleToModule(ctx types.Context, senderModule string, recipientModule string, amt types.Coins) error {
	ret := _m.Called(ctx, senderModule, recipientModule, amt)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinsFromModuleToModule")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, string, string, types.Coins) error); ok {
		r0 = rf(ctx, senderModule, recipientModule, amt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SpendableCoins provides a mock function with given fields: ctx, addr
func (_m *EmissionBankKeeper) SpendableCoins(ctx types.Context, addr types.AccAddress) types.Coins {
	ret := _m.Called(ctx, addr)

	if len(ret) == 0 {
		panic("no return value specified for SpendableCoins")
	}

	var r0 types.Coins
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) types.Coins); ok {
		r0 = rf(ctx, addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(types.Coins)
		}
	}

	return r0
}

// NewEmissionBankKeeper creates a new instance of EmissionBankKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEmissionBankKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *EmissionBankKeeper {
	mock := &EmissionBankKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
