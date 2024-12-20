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

// GetByID provides a mock function with given fields: ctx, id
func (_m *TransferRepository) GetByID(ctx context.Context, id int) (*domain.Transfer, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *domain.Transfer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*domain.Transfer, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *domain.Transfer); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Transfer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type TransferRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *TransferRepository_Expecter) GetByID(ctx interface{}, id interface{}) *TransferRepository_GetByID_Call {
	return &TransferRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *TransferRepository_GetByID_Call) Run(run func(ctx context.Context, id int)) *TransferRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *TransferRepository_GetByID_Call) Return(_a0 *domain.Transfer, _a1 error) *TransferRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TransferRepository_GetByID_Call) RunAndReturn(run func(context.Context, int) (*domain.Transfer, error)) *TransferRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
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

// UpdateStatus provides a mock function with given fields: ctx, id, status
func (_m *TransferRepository) UpdateStatus(ctx context.Context, id int, status domain.TransferStatus) error {
	ret := _m.Called(ctx, id, status)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, domain.TransferStatus) error); ok {
		r0 = rf(ctx, id, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TransferRepository_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type TransferRepository_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
//   - status domain.TransferStatus
func (_e *TransferRepository_Expecter) UpdateStatus(ctx interface{}, id interface{}, status interface{}) *TransferRepository_UpdateStatus_Call {
	return &TransferRepository_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, id, status)}
}

func (_c *TransferRepository_UpdateStatus_Call) Run(run func(ctx context.Context, id int, status domain.TransferStatus)) *TransferRepository_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(domain.TransferStatus))
	})
	return _c
}

func (_c *TransferRepository_UpdateStatus_Call) Return(_a0 error) *TransferRepository_UpdateStatus_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransferRepository_UpdateStatus_Call) RunAndReturn(run func(context.Context, int, domain.TransferStatus) error) *TransferRepository_UpdateStatus_Call {
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
