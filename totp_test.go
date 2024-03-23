//go:build unit || !integration

package otp_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.nhat.io/clock"

	"go.nhat.io/otp"
	"go.nhat.io/otp/mock"
)

func TestTOTPSecret_MarshalText(t *testing.T) {
	t.Parallel()

	s := otp.TOTPSecret("secret")

	data, err := s.MarshalText()
	require.NoError(t, err)

	var s2 otp.TOTPSecret

	err = s2.UnmarshalText(data)
	require.NoError(t, err)

	assert.Equal(t, s, s2)
}

func TestTOTPSecret_String(t *testing.T) {
	t.Parallel()

	s := otp.TOTPSecret("secret")

	assert.Equal(t, "secret", s.String())
}

func TestTOTPSecret_TOTPSecret(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := otp.TOTPSecret("secret")

	assert.Equal(t, "secret", string(p.TOTPSecret(ctx)))

	err := p.SetTOTPSecret(ctx, "changed", "")
	require.ErrorIs(t, err, otp.ErrTOTPSecretReadOnly)

	assert.Equal(t, "secret", string(p.TOTPSecret(ctx)))

	err = p.DeleteTOTPSecret(ctx)
	require.ErrorIs(t, err, otp.ErrTOTPSecretReadOnly)

	assert.Equal(t, "secret", string(p.TOTPSecret(ctx)))
}

func TestTOTPSecret_TOTPSecretGetter(t *testing.T) {
	t.Parallel()

	s := otp.TOTPSecret("secret")

	assert.Equal(t, s, s.TOTPSecretGetter())
}

func TestNoTOTPSecret(t *testing.T) {
	t.Parallel()

	assert.Empty(t, otp.NoTOTPSecret.TOTPSecret(context.Background()))
	assert.Empty(t, otp.NoTOTPSecret, otp.NoTOTPSecret.TOTPSecretGetter())
}

func TestTOTPSecretFromEnv(t *testing.T) {
	t.Setenv(t.Name(), "secret")

	ctx := context.Background()
	p := otp.TOTPSecretFromEnv(t.Name())

	assert.Equal(t, "secret", string(p.TOTPSecret(ctx)))

	err := p.SetTOTPSecret(ctx, "changed", "")
	require.NoError(t, err)

	assert.Equal(t, "changed", string(p.TOTPSecret(ctx)))

	err = p.DeleteTOTPSecret(ctx)
	require.NoError(t, err)

	assert.Empty(t, string(p.TOTPSecret(ctx)))

	assert.Equal(t, p, p.TOTPSecretGetter())
	assert.Equal(t, p, p.TOTPSecretSetter())
	assert.Equal(t, p, p.TOTPSecretDeleter())
}

func TestChainTOTPSecretGetters_HasSecret(t *testing.T) {
	t.Parallel()

	p := otp.ChainTOTPSecretGetters(
		otp.NoTOTPSecret,
		otp.ChainTOTPSecretGetters(
			nil,
			otp.TOTPSecret("secret"),
		),
	)

	assert.Equal(t, "secret", string(p.TOTPSecret(context.Background())))
	assert.Equal(t, p, p.TOTPSecretGetter())
}

func TestChainTOTPSecretGetters_NoSecret(t *testing.T) {
	t.Parallel()

	p := otp.ChainTOTPSecretGetters()
	actual := p.TOTPSecret(context.Background())

	assert.Equal(t, otp.NoTOTPSecret, actual)
}

func TestChainTOTPSecretProviders_HasSecret(t *testing.T) {
	t.Parallel()

	p := otp.ChainTOTPSecretProviders(
		otp.NoTOTPSecret,
		otp.ChainTOTPSecretProviders(
			nil,
			otp.TOTPSecret("secret"),
		),
	)

	assert.Equal(t, "secret", string(p.TOTPSecret(context.Background())))
}

func TestChainTOTPSecretProviders_NoSecret(t *testing.T) {
	t.Parallel()

	p := otp.ChainTOTPSecretProviders()
	actual := p.TOTPSecret(context.Background())

	assert.Equal(t, otp.NoTOTPSecret, actual)
}

func TestChainTOTPSecretProviders_SetSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockProvider  mock.TOTPSecretProviderMocker
		expectedError string
	}{
		{
			scenario: "failed",
			mockProvider: mock.MockTOTPSecretProvider(func(p *mock.TOTPSecretProvider) {
				p.On("SetTOTPSecret", mock.Anything, otp.TOTPSecret("secret"), "issuer").
					Return(assert.AnError)
			}),
			expectedError: `assert.AnError general error for testing`,
		},
		{
			scenario: "readonly",
			mockProvider: mock.MockTOTPSecretProvider(func(p *mock.TOTPSecretProvider) {
				p.On("SetTOTPSecret", mock.Anything, otp.TOTPSecret("secret"), "issuer").
					Return(otp.ErrTOTPSecretReadOnly)
			}),
			expectedError: `totp secret is read-only`,
		},
		{
			scenario: "success",
			mockProvider: mock.MockTOTPSecretProvider(func(p *mock.TOTPSecretProvider) {
				p.On("SetTOTPSecret", mock.Anything, otp.TOTPSecret("secret"), "issuer").
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := otp.ChainTOTPSecretProviders(tc.mockProvider(t))

			err := p.SetTOTPSecret(context.Background(), "secret", "issuer")

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestChainTOTPSecretProviders_DeleteSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockProvider  mock.TOTPSecretProviderMocker
		expectedError string
	}{
		{
			scenario: "failed",
			mockProvider: mock.MockTOTPSecretProvider(func(p *mock.TOTPSecretProvider) {
				p.On("DeleteTOTPSecret", mock.Anything).
					Return(assert.AnError)
			}),
			expectedError: `assert.AnError general error for testing`,
		},
		{
			scenario: "readonly",
			mockProvider: mock.MockTOTPSecretProvider(func(p *mock.TOTPSecretProvider) {
				p.On("DeleteTOTPSecret", mock.Anything).
					Return(otp.ErrTOTPSecretReadOnly)
			}),
			expectedError: `totp secret is read-only`,
		},
		{
			scenario: "success",
			mockProvider: mock.MockTOTPSecretProvider(func(p *mock.TOTPSecretProvider) {
				p.On("DeleteTOTPSecret", mock.Anything).
					Return(nil)
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := otp.ChainTOTPSecretProviders(tc.mockProvider(t))

			err := p.DeleteTOTPSecret(context.Background())

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestChainTOTPSecretProviders(t *testing.T) { //nolint: paralleltest
	t.Setenv(t.Name(), "env")

	providers := otp.TOTPSecretProviders{
		otp.NoTOTPSecret,
		otp.TOTPSecret("secret"),
		otp.TOTPSecretFromEnv(t.Name()),
	}

	ctx := context.Background()
	p := otp.ChainTOTPSecretProviders(providers...)

	assert.Equal(t, p, p.TOTPSecretGetter())
	assert.Equal(t, p, p.TOTPSecretSetter())
	assert.Equal(t, p, p.TOTPSecretDeleter())

	assert.Equal(t, "secret", string(p.TOTPSecret(ctx)))

	err := p.SetTOTPSecret(ctx, "changed", "")
	require.NoError(t, err)

	// Change secret.
	assert.Equal(t, "changed", string(p.TOTPSecret(ctx)))

	actualSecrets := make([]otp.TOTPSecret, len(providers))
	expectedSecrets := []otp.TOTPSecret{"", "secret", "changed"}
	surrogateSecrets := make([]otp.TOTPSecret, len(providers))
	expectedSurrogateSecrets := []otp.TOTPSecret{"changed", "changed", "changed"}

	for i := 0; i < len(providers); i++ {
		actualSecrets[i] = providers[i].TOTPSecret(ctx)
		surrogateSecrets[i] = p[i].TOTPSecret(ctx)
	}

	assert.Equal(t, expectedSecrets, actualSecrets)
	assert.Equal(t, expectedSurrogateSecrets, surrogateSecrets)

	// Delete secret.
	err = p.DeleteTOTPSecret(ctx)
	require.NoError(t, err)

	actualSecrets = make([]otp.TOTPSecret, len(providers))
	expectedSecrets = []otp.TOTPSecret{"", "secret", ""}
	surrogateSecrets = make([]otp.TOTPSecret, len(providers))
	expectedSurrogateSecrets = []otp.TOTPSecret{"", "", ""}

	for i := 0; i < len(providers); i++ {
		actualSecrets[i] = providers[i].TOTPSecret(ctx)
		surrogateSecrets[i] = p[i].TOTPSecret(ctx)
	}

	assert.Equal(t, expectedSecrets, actualSecrets)
	assert.Equal(t, expectedSurrogateSecrets, surrogateSecrets)
}

func TestTOTPGenerator_GenerateOTP(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario         string
		mockSecretGetter mock.TOTPSecretGetterMocker
		expectedResult   otp.OTP
		expectedError    string
	}{
		{
			scenario: "could not get totp secret",
			mockSecretGetter: mock.MockTOTPSecretGetter(func(g *mock.TOTPSecretGetter) {
				g.On("TOTPSecret", context.Background()).
					Return(otp.NoTOTPSecret)
			}),
			expectedError: "could not generate otp: no totp secret",
		},
		{
			scenario: "could not generate otp",
			mockSecretGetter: mock.MockTOTPSecretGetter(func(g *mock.TOTPSecretGetter) {
				g.On("TOTPSecret", context.Background()).
					Return(otp.TOTPSecret("secret"))
			}),
			expectedError: "could not generate otp: Decoding of secret as base32 failed.",
		},
		{
			scenario: "success",
			mockSecretGetter: mock.MockTOTPSecretGetter(func(g *mock.TOTPSecretGetter) {
				g.On("TOTPSecret", context.Background()).
					Return(otp.TOTPSecret("NBSWY3DP"))
			}),
			expectedResult: "191882",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			c := clock.Fix(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))
			g := otp.NewTOTPGenerator(tc.mockSecretGetter(t), otp.WithClock(c))

			result, err := g.GenerateOTP(context.Background())

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGenerateTOTP(t *testing.T) {
	t.Parallel()

	c := clock.Fix(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))

	result, err := otp.GenerateTOTP(context.Background(), otp.TOTPSecret("NBSWY3DP"), otp.WithClock(c))

	require.NoError(t, err)
	assert.Equal(t, otp.OTP("191882"), result)
}
