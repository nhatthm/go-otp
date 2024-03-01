package keyring

import "github.com/bool64/ctxd"

// Option configures the services provided by the keyring package.
type Option interface {
	TOTPSecretProviderOption
}

type option struct {
	TOTPSecretProviderOption
}

// WithLogger sets the logger for the keyring package.
func WithLogger(l ctxd.Logger) Option {
	return option{
		TOTPSecretProviderOption: totpSecretProviderOptionFunc(func(s *TOTPSecretProvider) {
			s.logger = l
		}),
	}
}
