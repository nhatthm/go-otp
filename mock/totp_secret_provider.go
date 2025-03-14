// Code generated by mockery v2.53.2. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	otp "go.nhat.io/otp"
)

// TOTPSecretProvider is an autogenerated mock type for the TOTPSecretProvider type
type TOTPSecretProvider struct {
	mock.Mock
}

// DeleteTOTPSecret provides a mock function with given fields: ctx
func (_m *TOTPSecretProvider) DeleteTOTPSecret(ctx context.Context) error {
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

// SetTOTPSecret provides a mock function with given fields: ctx, secret, issuer
func (_m *TOTPSecretProvider) SetTOTPSecret(ctx context.Context, secret otp.TOTPSecret, issuer string) error {
	ret := _m.Called(ctx, secret, issuer)

	if len(ret) == 0 {
		panic("no return value specified for SetTOTPSecret")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, otp.TOTPSecret, string) error); ok {
		r0 = rf(ctx, secret, issuer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TOTPSecret provides a mock function with given fields: ctx
func (_m *TOTPSecretProvider) TOTPSecret(ctx context.Context) otp.TOTPSecret {
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

// NewTOTPSecretProvider creates a new instance of TOTPSecretProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTOTPSecretProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *TOTPSecretProvider {
	mock := &TOTPSecretProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
