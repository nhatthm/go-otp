// Code generated by mockery v2.45.0. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	otp "go.nhat.io/otp"
)

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// GenerateOTP provides a mock function with given fields: ctx
func (_m *Generator) GenerateOTP(ctx context.Context) (otp.OTP, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GenerateOTP")
	}

	var r0 otp.OTP
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (otp.OTP, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) otp.OTP); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(otp.OTP)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewGenerator creates a new instance of Generator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGenerator(t interface {
	mock.TestingT
	Cleanup(func())
}) *Generator {
	mock := &Generator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
