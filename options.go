package otp

import "go.nhat.io/clock"

// Option configures the apis of the authenticator package.
type Option interface {
	TOTPGeneratorOption
}

type option struct {
	TOTPGeneratorOption
}

// WithClock sets the clock of the TOTPGenerator.
func WithClock(c clock.Clock) Option {
	return option{
		TOTPGeneratorOption: totpGeneratorOptionFunc(func(g *TOTPGenerator) {
			g.clock = c
		}),
	}
}
