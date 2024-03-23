package otp

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/pquerna/otp/totp"
	"go.nhat.io/clock"
)

// ErrTOTPSecretReadOnly indicates that the TOTP secret is read-only.
var ErrTOTPSecretReadOnly = errors.New("totp secret is read-only")

// ErrNoTOTPSecret indicates that the user has not configured the TOTP secret.
var ErrNoTOTPSecret = errors.New("no totp secret")

// NoTOTPSecret is a TOTP secret that is empty.
const NoTOTPSecret = TOTPSecret("")

// TOTPSecret is a TOTP secret.
type TOTPSecret string

// MarshalText returns the TOTP secret as text.
func (s TOTPSecret) MarshalText() ([]byte, error) { //nolint: unparam
	return []byte(s), nil
}

// UnmarshalText unmarshals the TOTP secret from text.
func (s *TOTPSecret) UnmarshalText(text []byte) error { //nolint: unparam
	*s = TOTPSecret(text)

	return nil
}

// String returns the TOTP secret as a string.
func (s TOTPSecret) String() string {
	return string(s)
}

// TOTPSecret returns the TOTP secret.
func (s TOTPSecret) TOTPSecret(context.Context) TOTPSecret {
	return s
}

// SetTOTPSecret sets the TOTP secret.
func (s TOTPSecret) SetTOTPSecret(context.Context, TOTPSecret, string) error {
	return ErrTOTPSecretReadOnly
}

// DeleteTOTPSecret deletes the TOTP secret.
func (s TOTPSecret) DeleteTOTPSecret(context.Context) error {
	return ErrTOTPSecretReadOnly
}

// TOTPSecretGetter returns TOTPSecretGetter.
func (s TOTPSecret) TOTPSecretGetter() TOTPSecretGetter {
	return s
}

type totpSecreteSurrogate struct {
	secret TOTPSecret
}

// TOTPSecret returns the TOTP secret.
func (s totpSecreteSurrogate) TOTPSecret(context.Context) TOTPSecret {
	return s.secret
}

// SetTOTPSecret sets the TOTP secret.
func (s *totpSecreteSurrogate) SetTOTPSecret(_ context.Context, secret TOTPSecret, _ string) error {
	s.secret = secret

	return nil
}

// DeleteTOTPSecret deletes the TOTP secret.
func (s *totpSecreteSurrogate) DeleteTOTPSecret(context.Context) error {
	s.secret = NoTOTPSecret

	return nil
}

// TOTPSecretProvider is an interface that manages a TOTP secret.
type TOTPSecretProvider interface {
	TOTPSecretGetter
	TOTPSecretSetter
	TOTPSecretDeleter
}

// TOTPSecretGetter is an interface that provides a TOTP secret.
type TOTPSecretGetter interface {
	TOTPSecret(ctx context.Context) TOTPSecret
}

// TOTPSecretSetter is an interface that sets a TOTP secret.
type TOTPSecretSetter interface {
	SetTOTPSecret(ctx context.Context, secret TOTPSecret, issuer string) error
}

// TOTPSecretDeleter is an interface that deletes a TOTP secret.
type TOTPSecretDeleter interface {
	DeleteTOTPSecret(ctx context.Context) error
}

// TOTPSecretGetters is a list of TOTP secret getters.
type TOTPSecretGetters []TOTPSecretGetter

// TOTPSecret returns the first non-empty TOTP secret that it finds from the list of TOTP secret getters.
func (p TOTPSecretGetters) TOTPSecret(ctx context.Context) TOTPSecret {
	for _, sp := range p {
		if s := sp.TOTPSecret(ctx); s != NoTOTPSecret {
			return s
		}
	}

	return NoTOTPSecret
}

// TOTPSecretGetter returns TOTPSecretGetter.
func (p TOTPSecretGetters) TOTPSecretGetter() TOTPSecretGetter {
	return p
}

// ChainTOTPSecretGetters chains the TOTP secret getters.
func ChainTOTPSecretGetters(getters ...TOTPSecretGetter) TOTPSecretGetters {
	result := make(TOTPSecretGetters, 0, len(getters))

	for _, g := range getters {
		switch g := g.(type) {
		case nil:
		case TOTPSecretGetters:
			result = append(result, g...)
		default:
			result = append(result, g)
		}
	}

	return result
}

// TOTPSecretProviders is a list of TOTP secret getters.
type TOTPSecretProviders []TOTPSecretProvider

// TOTPSecret returns the first non-empty TOTP secret that it finds from the list of TOTP secret getters.
func (ps TOTPSecretProviders) TOTPSecret(ctx context.Context) TOTPSecret {
	for _, p := range ps {
		if s := p.TOTPSecret(ctx); s != NoTOTPSecret {
			return s
		}
	}

	return NoTOTPSecret
}

// SetTOTPSecret sets the TOTP secret.
func (ps TOTPSecretProviders) SetTOTPSecret(ctx context.Context, secret TOTPSecret, issuer string) error {
	for _, p := range ps {
		if err := p.SetTOTPSecret(ctx, secret, issuer); err != nil {
			return err
		}
	}

	return nil
}

// DeleteTOTPSecret deletes the TOTP secret.
func (ps TOTPSecretProviders) DeleteTOTPSecret(ctx context.Context) error {
	for _, p := range ps {
		if err := p.DeleteTOTPSecret(ctx); err != nil {
			return err
		}
	}

	return nil
}

// TOTPSecretGetter returns TOTPSecretGetter.
func (ps TOTPSecretProviders) TOTPSecretGetter() TOTPSecretGetter {
	return ps
}

// TOTPSecretSetter returns TOTPSecretSetter.
func (ps TOTPSecretProviders) TOTPSecretSetter() TOTPSecretSetter {
	return ps
}

// TOTPSecretDeleter returns TOTPSecretDeleter.
func (ps TOTPSecretProviders) TOTPSecretDeleter() TOTPSecretDeleter {
	return ps
}

// ChainTOTPSecretProviders chains the TOTP secret providers.
func ChainTOTPSecretProviders(providers ...TOTPSecretProvider) TOTPSecretProviders {
	result := make(TOTPSecretProviders, 0, len(providers))

	for _, p := range providers {
		switch g := p.(type) {
		case nil:
		case TOTPSecret:
			result = append(result, &totpSecreteSurrogate{g})
		case TOTPSecretProviders:
			result = append(result, g...)
		default:
			result = append(result, g)
		}
	}

	return result
}

var _ TOTPSecretProvider = (*EnvTOTPSecret)(nil)

// EnvTOTPSecret is a TOTP secret provider that gets the TOTP secret from the environment.
type EnvTOTPSecret struct {
	env string
}

// TOTPSecret returns the TOTP secret from the environment.
func (e EnvTOTPSecret) TOTPSecret(_ context.Context) TOTPSecret {
	return TOTPSecret(os.Getenv(e.env))
}

// SetTOTPSecret sets the TOTP secret to the environment.
func (e EnvTOTPSecret) SetTOTPSecret(_ context.Context, secret TOTPSecret, _ string) error {
	return os.Setenv(e.env, string(secret))
}

// DeleteTOTPSecret deletes the TOTP secret from the environment.
func (e EnvTOTPSecret) DeleteTOTPSecret(_ context.Context) error {
	return os.Unsetenv(e.env)
}

// TOTPSecretGetter returns TOTPSecretGetter.
func (e EnvTOTPSecret) TOTPSecretGetter() TOTPSecretGetter {
	return e
}

// TOTPSecretSetter returns TOTPSecretSetter.
func (e EnvTOTPSecret) TOTPSecretSetter() TOTPSecretSetter {
	return e
}

// TOTPSecretDeleter returns TOTPSecretDeleter.
func (e EnvTOTPSecret) TOTPSecretDeleter() TOTPSecretDeleter {
	return e
}

// TOTPSecretFromEnv returns a TOTP secret getter that gets the TOTP secret from the environment.
func TOTPSecretFromEnv(env string) EnvTOTPSecret {
	return EnvTOTPSecret{
		env: env,
	}
}

var _ Generator = (*TOTPGenerator)(nil)

// TOTPGenerator is a .TOTPGenerator.
type TOTPGenerator struct {
	secretGetter TOTPSecretGetter
	clock        clock.Clock
}

// GenerateOTP generates a TOTP.
func (g *TOTPGenerator) GenerateOTP(ctx context.Context) (OTP, error) {
	s := g.secretGetter.TOTPSecret(ctx)
	if s == NoTOTPSecret {
		return "", fmt.Errorf("could not generate otp: %w", ErrNoTOTPSecret)
	}

	code, err := totp.GenerateCode(string(s), g.clock.Now())
	if err != nil {
		return "", fmt.Errorf("could not generate otp: %w", err)
	}

	return OTP(code), nil
}

// NewTOTPGenerator initiates a new .TOTPGenerator.
func NewTOTPGenerator(secretGetter TOTPSecretGetter, opts ...TOTPGeneratorOption) *TOTPGenerator {
	g := &TOTPGenerator{
		secretGetter: secretGetter,
		clock:        clock.New(),
	}

	for _, opt := range opts {
		opt.applyTOTPGeneratorOption(g)
	}

	return g
}

// GenerateTOTP generates a TOTP.
func GenerateTOTP(ctx context.Context, secret TOTPSecretGetter, opts ...TOTPGeneratorOption) (OTP, error) {
	return NewTOTPGenerator(secret, opts...).GenerateOTP(ctx)
}

// TOTPGeneratorOption is an option to configure TOTPGenerator.
type TOTPGeneratorOption interface {
	applyTOTPGeneratorOption(g *TOTPGenerator)
}

type totpGeneratorOptionFunc func(g *TOTPGenerator)

func (f totpGeneratorOptionFunc) applyTOTPGeneratorOption(g *TOTPGenerator) {
	f(g)
}
