// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/tonytcb/inventory-management-system/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// TransferNotifierHandler is an autogenerated mock type for the TransferNotifierHandler type
type TransferNotifierHandler struct {
	mock.Mock
}

type TransferNotifierHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *TransferNotifierHandler) EXPECT() *TransferNotifierHandler_Expecter {
	return &TransferNotifierHandler_Expecter{mock: &_m.Mock}
}

// Settlement provides a mock function with given fields: _a0, _a1, _a2
func (_m *TransferNotifierHandler) Settlement(_a0 context.Context, _a1 *domain.Transfer, _a2 *domain.FXRate) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Settlement")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Transfer, *domain.FXRate) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransferNotifierHandler_Settlement_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Settlement'
type TransferNotifierHandler_Settlement_Call struct {
	*mock.Call
}

// Settlement is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *domain.Transfer
//   - _a2 *domain.FXRate
func (_e *TransferNotifierHandler_Expecter) Settlement(_a0 interface{}, _a1 interface{}, _a2 interface{}) *TransferNotifierHandler_Settlement_Call {
	return &TransferNotifierHandler_Settlement_Call{Call: _e.mock.On("Settlement", _a0, _a1, _a2)}
}

func (_c *TransferNotifierHandler_Settlement_Call) Run(run func(_a0 context.Context, _a1 *domain.Transfer, _a2 *domain.FXRate)) *TransferNotifierHandler_Settlement_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Transfer), args[2].(*domain.FXRate))
	})
	return _c
}

func (_c *TransferNotifierHandler_Settlement_Call) Return(_a0 error) *TransferNotifierHandler_Settlement_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransferNotifierHandler_Settlement_Call) RunAndReturn(run func(context.Context, *domain.Transfer, *domain.FXRate) error) *TransferNotifierHandler_Settlement_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransferNotifierHandler creates a new instance of TransferNotifierHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransferNotifierHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransferNotifierHandler {
	mock := &TransferNotifierHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}