package otp

import "context"

// OTP is a one-time password.
type OTP string

// Generator is a one-time password generator.
type Generator interface {
	GenerateOTP(ctx context.Context) (OTP, error)
}
