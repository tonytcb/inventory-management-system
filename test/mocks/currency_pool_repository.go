// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	decimal "github.com/shopspring/decimal"
	domain "github.com/tonytcb/inventory-management-system/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// CurrencyPoolRepository is an autogenerated mock type for the CurrencyPoolRepository type
type CurrencyPoolRepository struct {
	mock.Mock
}

type CurrencyPoolRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *CurrencyPoolRepository) EXPECT() *CurrencyPoolRepository_Expecter {
	return &CurrencyPoolRepository_Expecter{mock: &_m.Mock}
}

// Debit provides a mock function with given fields: _a0, _a1, _a2
func (_m *CurrencyPoolRepository) Debit(_a0 context.Context, _a1 domain.Currency, _a2 decimal.Decimal) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for Debit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Currency, decimal.Decimal) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CurrencyPoolRepository_Debit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Debit'
type CurrencyPoolRepository_Debit_Call struct {
	*mock.Call
}

// Debit is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 domain.Currency
//   - _a2 decimal.Decimal
func (_e *CurrencyPoolRepository_Expecter) Debit(_a0 interface{}, _a1 interface{}, _a2 interface{}) *CurrencyPoolRepository_Debit_Call {
	return &CurrencyPoolRepository_Debit_Call{Call: _e.mock.On("Debit", _a0, _a1, _a2)}
}

func (_c *CurrencyPoolRepository_Debit_Call) Run(run func(_a0 context.Context, _a1 domain.Currency, _a2 decimal.Decimal)) *CurrencyPoolRepository_Debit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Currency), args[2].(decimal.Decimal))
	})
	return _c
}

func (_c *CurrencyPoolRepository_Debit_Call) Return(_a0 error) *CurrencyPoolRepository_Debit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CurrencyPoolRepository_Debit_Call) RunAndReturn(run func(context.Context, domain.Currency, decimal.Decimal) error) *CurrencyPoolRepository_Debit_Call {
	_c.Call.Return(run)
	return _c
}

// NewCurrencyPoolRepository creates a new instance of CurrencyPoolRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCurrencyPoolRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CurrencyPoolRepository {
	mock := &CurrencyPoolRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
