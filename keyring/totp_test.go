//go:build unit || !integration

package keyring_test

import (
	"context"
	"testing"

	"github.com/bool64/ctxd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"go.nhat.io/otp"
	"go.nhat.io/otp/keyring"
	mockotp "go.nhat.io/otp/mock"
)

func TestTOTPSecretGetSetter_TOTPSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockStorage    mockotp.StorageMocker[string]
		account        string
		expectedResult otp.TOTPSecret
	}{
		{
			scenario:       "no account",
			mockStorage:    mockotp.MockStorage[string](),
			account:        "",
			expectedResult: "",
		},
		{
			scenario: "storage error",
			mockStorage: mockotp.MockStorage(func(s *mockotp.Storage[string]) {
				s.On("Get", mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			}),
			account:        "account",
			expectedResult: "",
		},
		{
			scenario: "no secret",
			mockStorage: mockotp.MockStorage(func(s *mockotp.Storage[string]) {
				s.On("Get", "go.nhat.io/totp", "account").
					Return(ptr(""), nil)
			}),
			account:        "account",
			expectedResult: "",
		},
		{
			scenario: "has secret",
			mockStorage: mockotp.MockStorage(func(s *mockotp.Storage[string]) {
				s.On("Get", "go.nhat.io/totp", "account").
					Return(ptr("secret"), nil)
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

func TestTOTPSecretGetSetter_SetTOTPSecret(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockStorage   mockotp.StorageMocker[string]
		account       string
		expectedError string
	}{
		{
			scenario:    "no account",
			mockStorage: mockotp.MockStorage[string](),
			account:     "",
		},
		{
			scenario: "storage error",
			mockStorage: mockotp.MockStorage(func(s *mockotp.Storage[string]) {
				s.On("Set", mock.Anything, mock.Anything, mock.Anything).
					Return(assert.AnError)
			}),
			account:       "account",
			expectedError: "assert.AnError general error for testing",
		},
		{
			scenario: "success",
			mockStorage: mockotp.MockStorage(func(s *mockotp.Storage[string]) {
				s.On("Set", "go.nhat.io/totp", "account", "secret").
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

func ptr[V any](v V) *V {
	return &v
}
