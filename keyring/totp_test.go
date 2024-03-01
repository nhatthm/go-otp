//go:build unit || !integration

package keyring_test

import (
	"context"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mockss "go.nhat.io/secretstorage/mock"

	"go.nhat.io/otp"
	"go.nhat.io/otp/keyring"
)

func TestTOTPSecretProvider_TOTPSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockStorage    mockss.StorageMocker[otp.TOTPSecret]
		account        string
		expectedResult otp.TOTPSecret
	}{
		{
			scenario:       "no account",
			mockStorage:    mockss.MockStorage[otp.TOTPSecret](),
			account:        "",
			expectedResult: "",
		},
		{
			scenario: "storage error",
			mockStorage: mockss.MockStorage(func(s *mockss.Storage[otp.TOTPSecret]) {
				s.On("Get", mock.Anything, mock.Anything).
					Return(otp.NoTOTPSecret, assert.AnError)
			}),
			account:        "account",
			expectedResult: "",
		},
		{
			scenario: "no secret",
			mockStorage: mockss.MockStorage(func(s *mockss.Storage[otp.TOTPSecret]) {
				s.On("Get", "go.nhat.io/totp", "account").
					Return(otp.NoTOTPSecret, nil)
			}),
			account:        "account",
			expectedResult: "",
		},
		{
			scenario: "has secret",
			mockStorage: mockss.MockStorage(func(s *mockss.Storage[otp.TOTPSecret]) {
				s.On("Get", "go.nhat.io/totp", "account").
					Return(otp.TOTPSecret("secret"), nil)
			}),
			account:        "account",
			expectedResult: "secret",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := keyring.TOTPSecretFromKeyring(tc.account,
				keyring.WithStorage(tc.mockStorage(t)),
				keyring.WithLogger(ctxd.NoOpLogger{}),
			)

			actual := s.TOTPSecret(context.Background())

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestTOTPSecretProvider_SetTOTPSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockStorage   mockss.StorageMocker[otp.TOTPSecret]
		account       string
		expectedError string
	}{
		{
			scenario:    "no account",
			mockStorage: mockss.MockStorage[otp.TOTPSecret](),
			account:     "",
		},
		{
			scenario: "storage error",
			mockStorage: mockss.MockStorage(func(s *mockss.Storage[otp.TOTPSecret]) {
				s.On("Set", mock.Anything, mock.Anything, mock.Anything).
					Return(assert.AnError)
			}),
			account:       "account",
			expectedError: "assert.AnError general error for testing",
		},
		{
			scenario: "success",
			mockStorage: mockss.MockStorage(func(s *mockss.Storage[otp.TOTPSecret]) {
				s.On("Set", "go.nhat.io/totp", "account", otp.TOTPSecret("secret")).
					Return(nil)
			}),
			account: "account",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := keyring.TOTPSecretFromKeyring(tc.account,
				keyring.WithStorage(tc.mockStorage(t)),
				keyring.WithLogger(ctxd.NoOpLogger{}),
			)

			err := s.SetTOTPSecret(context.Background(), "secret")

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTOTPSecretProvider_DeleteTOTPSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockStorage   mockss.StorageMocker[otp.TOTPSecret]
		account       string
		expectedError string
	}{
		{
			scenario:    "no account",
			mockStorage: mockss.MockStorage[otp.TOTPSecret](),
			account:     "",
		},
		{
			scenario: "storage error",
			mockStorage: mockss.MockStorage(func(s *mockss.Storage[otp.TOTPSecret]) {
				s.On("Delete", mock.Anything, mock.Anything).
					Return(assert.AnError)
			}),
			account:       "account",
			expectedError: "assert.AnError general error for testing",
		},
		{
			scenario: "success",
			mockStorage: mockss.MockStorage(func(s *mockss.Storage[otp.TOTPSecret]) {
				s.On("Delete", "go.nhat.io/totp", "account").
					Return(nil)
			}),
			account: "account",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			s := keyring.TOTPSecretFromKeyring(tc.account,
				keyring.WithStorage(tc.mockStorage(t)),
				keyring.WithLogger(ctxd.NoOpLogger{}),
			)

			err := s.DeleteTOTPSecret(context.Background())

			if tc.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
