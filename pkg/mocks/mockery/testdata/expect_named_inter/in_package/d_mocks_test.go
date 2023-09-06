// Code generated by mockery. DO NOT EDIT.

package inter_test

import mock "github.com/stretchr/testify/mock"

// InterD is an autogenerated mock type for the interD type
type InterD struct {
	mock.Mock
}

type InterD_Expecter struct {
	mock *mock.Mock
}

func (_m *InterD) EXPECT() *InterD_Expecter {
	return &InterD_Expecter{mock: &_m.Mock}
}

// FuncA provides a mock function with given fields:
func (_m *InterD) FuncA() {
	_m.Called()
}

// InterD_FuncA_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FuncA'
type InterD_FuncA_Call struct {
	*mock.Call
}

// FuncA is a helper method to define mock.On call
func (_e *InterD_Expecter) FuncA() *InterD_FuncA_Call {
	return &InterD_FuncA_Call{Call: _e.mock.On("FuncA")}
}

func (_c *InterD_FuncA_Call) Run(run func()) *InterD_FuncA_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *InterD_FuncA_Call) Return() *InterD_FuncA_Call {
	_c.Call.Return()
	return _c
}

func (_c *InterD_FuncA_Call) RunAndReturn(run func()) *InterD_FuncA_Call {
	_c.Call.Return(run)
	return _c
}

// NewInterD creates a new instance of InterD. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInterD(t interface {
	mock.TestingT
	Cleanup(func())
}) *InterD {
	mock := &InterD{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}