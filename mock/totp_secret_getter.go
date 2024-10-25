// Code generated by mockery v2.46.3. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	otp "go.nhat.io/otp"
)

// TOTPSecretGetter is an autogenerated mock type for the TOTPSecretGetter type
type TOTPSecretGetter struct {
	mock.Mock
}

// TOTPSecret provides a mock function with given fields: ctx
func (_m *TOTPSecretGetter) TOTPSecret(ctx context.Context) otp.TOTPSecret {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for TOTPSecret")
	}

	var r0 otp.TOTPSecret
	if rf, ok := ret.Get(0).(func(context.Context) otp.TOTPSecret); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(otp.TOTPSecret)
	}

	return r0
}

// NewTOTPSecretGetter creates a new instance of TOTPSecretGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTOTPSecretGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *TOTPSecretGetter {
	mock := &TOTPSecretGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
