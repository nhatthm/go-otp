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

	p := otp.TOTPSecret("secret")

	assert.Equal(t, "secret", string(p.TOTPSecret(context.Background())))
}

func TestNoTOTPSecret(t *testing.T) {
	t.Parallel()

	assert.Empty(t, otp.NoTOTPSecret.TOTPSecret(context.Background()))
}

func TestTOTPSecretFromEnv(t *testing.T) {
	t.Setenv(t.Name(), "secret")

	p := otp.TOTPSecretFromEnv(t.Name())

	assert.Equal(t, "secret", string(p.TOTPSecret(context.Background())))
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
}

func TestChainTOTPSecretGetters_NoSecret(t *testing.T) {
	t.Parallel()

	p := otp.ChainTOTPSecretGetters()
	actual := p.TOTPSecret(context.Background())

	assert.Equal(t, otp.NoTOTPSecret, actual)
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
