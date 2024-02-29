package keyring

import (
	"context"
	"sync"

	"github.com/bool64/ctxd"
	"go.nhat.io/secretstorage"

	"go.nhat.io/otp"
)

const keyringServiceTOTP = "go.nhat.io/totp"

var (
	_ otp.TOTPSecretGetter = (*TOTPSecretGetSetter)(nil)
	_ otp.TOTPSecretSetter = (*TOTPSecretGetSetter)(nil)
)

// TOTPSecretGetSetter is a TOTP secret getter and setter that uses the keyring to store the TOTP secret.
type TOTPSecretGetSetter struct {
	storage secretstorage.Storage[otp.TOTPSecret]
	logger  ctxd.Logger

	account   string
	secret    otp.TOTPSecret
	fetchOnce sync.Once
}

func (s *TOTPSecretGetSetter) fetch(ctx context.Context) otp.TOTPSecret {
	if s.account == "" {
		return otp.NoTOTPSecret
	}

	secret, err := s.storage.Get(keyringServiceTOTP, s.account)
	if err != nil {
		s.logger.Error(ctx, "could not get totp secret from keyring", "error", err, "service", keyringServiceTOTP, "account", s.account)

		return otp.NoTOTPSecret
	}

	return secret
}

// TOTPSecret returns the TOTP secret from the keyring.
func (s *TOTPSecretGetSetter) TOTPSecret(ctx context.Context) otp.TOTPSecret {
	s.fetchOnce.Do(func() {
		s.secret = s.fetch(ctx)
	})

	return s.secret
}

// SetTOTPSecret persists the TOTP secret to the keyring.
func (s *TOTPSecretGetSetter) SetTOTPSecret(ctx context.Context, secret otp.TOTPSecret) error {
	if s.account == "" {
		return nil
	}

	if err := s.storage.Set(keyringServiceTOTP, s.account, secret); err != nil {
		s.logger.Error(ctx, "could not persist totp secret to keyring", "error", err, "service", keyringServiceTOTP, "account", s.account)

		return err
	}

	return nil
}

// TOTPSecretFromKeyring returns a TOTP secret getter and setter that uses the keyring to store the TOTP secret.
func TOTPSecretFromKeyring(account string, opts ...TOTPSecretGetSetterOption) *TOTPSecretGetSetter {
	s := &TOTPSecretGetSetter{
		storage: secretstorage.NewKeyringStorage[otp.TOTPSecret](),
		logger:  ctxd.NoOpLogger{},

		account: account,
	}

	for _, opt := range opts {
		opt.applyTOTPSecretGetSetterOption(s)
	}

	return s
}

// TOTPSecretGetSetterOption is an option to configure TOTPSecretGetSetter.
type TOTPSecretGetSetterOption interface {
	applyTOTPSecretGetSetterOption(s *TOTPSecretGetSetter)
}

type totpSecretGetSetterOptionFunc func(s *TOTPSecretGetSetter)

func (f totpSecretGetSetterOptionFunc) applyTOTPSecretGetSetterOption(s *TOTPSecretGetSetter) {
	f(s)
}

// WithStorage sets the storage for the TOTP secret getter and setter.
func WithStorage(storage secretstorage.Storage[otp.TOTPSecret]) TOTPSecretGetSetterOption {
	return totpSecretGetSetterOptionFunc(func(s *TOTPSecretGetSetter) {
		s.storage = storage
	})
}
