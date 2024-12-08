// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/tonytcb/inventory-management-system/internal/domain"
)

// TransferNotifier is an autogenerated mock type for the TransferNotifier type
type TransferNotifier struct {
	mock.Mock
}

type TransferNotifier_Expecter struct {
	mock *mock.Mock
}

func (_m *TransferNotifier) EXPECT() *TransferNotifier_Expecter {
	return &TransferNotifier_Expecter{mock: &_m.Mock}
}

// Created provides a mock function with given fields: _a0, _a1
func (_m *TransferNotifier) Created(_a0 context.Context, _a1 *domain.TransferCreatedEvent) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Created")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.TransferCreatedEvent) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransferNotifier_Created_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Created'
type TransferNotifier_Created_Call struct {
	*mock.Call
}

// Created is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *domain.TransferCreatedEvent
func (_e *TransferNotifier_Expecter) Created(_a0 interface{}, _a1 interface{}) *TransferNotifier_Created_Call {
	return &TransferNotifier_Created_Call{Call: _e.mock.On("Created", _a0, _a1)}
}

func (_c *TransferNotifier_Created_Call) Run(run func(_a0 context.Context, _a1 *domain.TransferCreatedEvent)) *TransferNotifier_Created_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.TransferCreatedEvent))
	})
	return _c
}

func (_c *TransferNotifier_Created_Call) Return(_a0 error) *TransferNotifier_Created_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransferNotifier_Created_Call) RunAndReturn(run func(context.Context, *domain.TransferCreatedEvent) error) *TransferNotifier_Created_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransferNotifier creates a new instance of TransferNotifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransferNotifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransferNotifier {
	mock := &TransferNotifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
