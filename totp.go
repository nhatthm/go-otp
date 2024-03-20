package otp

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/pquerna/otp/totp"
	"go.nhat.io/clock"
)

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
func (s TOTPSecret) TOTPSecret(_ context.Context) TOTPSecret {
	return s
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
	SetTOTPSecret(ctx context.Context, secret TOTPSecret) error
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

// ChainTOTPSecretGetters chains the TOTP secret getters.
func ChainTOTPSecretGetters(getters ...TOTPSecretGetter) TOTPSecretGetter {
	result := make(TOTPSecretGetters, 0, len(getters))

	for _, g := range getters {
		switch g := g.(type) {
		case TOTPSecretGetters:
			result = append(result, g...)
		case nil:
			// ignore nil provider
		default:
			result = append(result, g)
		}
	}

	return result
}

var _ TOTPSecretProvider = (*envTOTPSecret)(nil)

type envTOTPSecret string

// TOTPSecret returns the TOTP secret from the environment.
func (e envTOTPSecret) TOTPSecret(_ context.Context) TOTPSecret {
	return TOTPSecret(os.Getenv(string(e)))
}

// SetTOTPSecret sets the TOTP secret to the environment.
func (e envTOTPSecret) SetTOTPSecret(_ context.Context, secret TOTPSecret) error {
	return os.Setenv(string(e), string(secret))
}

// DeleteTOTPSecret deletes the TOTP secret from the environment.
func (e envTOTPSecret) DeleteTOTPSecret(_ context.Context) error {
	return os.Unsetenv(string(e))
}

// TOTPSecretFromEnv returns a TOTP secret getter that gets the TOTP secret from the environment.
func TOTPSecretFromEnv(env string) TOTPSecretProvider {
	return envTOTPSecret(env)
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
