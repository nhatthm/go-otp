package otp

import "context"

// OTP is a one-time password.
type OTP string

// String returns the string representation of the one-time password.
func (o OTP) String() string {
	return string(o)
}

// Generator is a one-time password generator.
type Generator interface {
	GenerateOTP(ctx context.Context) (OTP, error)
}
