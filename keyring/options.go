package keyring

import "github.com/bool64/ctxd"

// Option configures the services provided by the keyring package.
type Option interface {
	TOTPSecretGetSetterOption
}

type option struct {
	TOTPSecretGetSetterOption
}

// WithLogger sets the logger for the keyring package.
func WithLogger(l ctxd.Logger) Option {
	return option{
		TOTPSecretGetSetterOption: totpSecretGetSetterOptionFunc(func(s *TOTPSecretGetSetter) {
			s.logger = l
		}),
	}
}
