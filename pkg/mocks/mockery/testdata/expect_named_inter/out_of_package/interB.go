// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// InterB is an autogenerated mock type for the interB type
type InterB struct {
	mock.Mock
}

type InterB_Expecter struct {
	mock *mock.Mock
}

func (_m *InterB) EXPECT() *InterB_Expecter {
	return &InterB_Expecter{mock: &_m.Mock}
}

// FuncA provides a mock function with given fields:
func (_m *InterB) FuncA() {
	_m.Called()
}

// InterB_FuncA_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FuncA'
type InterB_FuncA_Call struct {
	*mock.Call
}

// FuncA is a helper method to define mock.On call
func (_e *InterB_Expecter) FuncA() *InterB_FuncA_Call {
	return &InterB_FuncA_Call{Call: _e.mock.On("FuncA")}
}

func (_c *InterB_FuncA_Call) Run(run func()) *InterB_FuncA_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InterB_FuncA_Call) Return() *InterB_FuncA_Call {
	_c.Call.Return()
	return _c
}

func (_c *InterB_FuncA_Call) RunAndReturn(run func()) *InterB_FuncA_Call {
	_c.Call.Return(run)
	return _c
}

// FuncB provides a mock function with given fields:
func (_m *InterB) FuncB() {
	_m.Called()
}

// InterB_FuncB_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FuncB'
type InterB_FuncB_Call struct {
	*mock.Call
}

// FuncB is a helper method to define mock.On call
func (_e *InterB_Expecter) FuncB() *InterB_FuncB_Call {
	return &InterB_FuncB_Call{Call: _e.mock.On("FuncB")}
}

func (_c *InterB_FuncB_Call) Run(run func()) *InterB_FuncB_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InterB_FuncB_Call) Return() *InterB_FuncB_Call {
	_c.Call.Return()
	return _c
}

func (_c *InterB_FuncB_Call) RunAndReturn(run func()) *InterB_FuncB_Call {
	_c.Call.Return(run)
	return _c
}

// NewInterB creates a new instance of InterB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInterB(t interface {
	mock.TestingT
	Cleanup(func())
}) *InterB {
	mock := &InterB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
