// Code generated by mockery v2.53.2. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TOTPSecretDeleter is an autogenerated mock type for the TOTPSecretDeleter type
type TOTPSecretDeleter struct {
	mock.Mock
}

// DeleteTOTPSecret provides a mock function with given fields: ctx
func (_m *TOTPSecretDeleter) DeleteTOTPSecret(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteTOTPSecret")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTOTPSecretDeleter creates a new instance of TOTPSecretDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTOTPSecretDeleter(t interface {
	mock.TestingT
	Cleanup(func())
}) *TOTPSecretDeleter {
	mock := &TOTPSecretDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
