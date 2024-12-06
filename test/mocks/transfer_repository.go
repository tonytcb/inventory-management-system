// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	domain "github.com/tonytcb/inventory-management-system/internal/domain"
)

// TransferRepository is an autogenerated mock type for the TransferRepository type
type TransferRepository struct {
	mock.Mock
}

type TransferRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *TransferRepository) EXPECT() *TransferRepository_Expecter {
	return &TransferRepository_Expecter{mock: &_m.Mock}
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *TransferRepository) Save(_a0 context.Context, _a1 *domain.Transfer) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Transfer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransferRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type TransferRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *domain.Transfer
func (_e *TransferRepository_Expecter) Save(_a0 interface{}, _a1 interface{}) *TransferRepository_Save_Call {
	return &TransferRepository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *TransferRepository_Save_Call) Run(run func(_a0 context.Context, _a1 *domain.Transfer)) *TransferRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.Transfer))
	})
	return _c
}

func (_c *TransferRepository_Save_Call) Return(_a0 error) *TransferRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransferRepository_Save_Call) RunAndReturn(run func(context.Context, *domain.Transfer) error) *TransferRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransferRepository creates a new instance of TransferRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransferRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransferRepository {
	mock := &TransferRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
