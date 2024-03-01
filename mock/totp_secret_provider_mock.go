package mock

import "testing"

// TOTPSecretProviderMocker is TOTPSecretProvider mocker.
type TOTPSecretProviderMocker func(tb testing.TB) *TOTPSecretProvider

// NopTOTPSecretProvider is no mock TOTPSecretProvider.
var NopTOTPSecretProvider = MockTOTPSecretProvider()

// MockTOTPSecretProvider creates TOTPSecretProvider mock with cleanup to ensure all the expectations are met.
func MockTOTPSecretProvider(mocks ...func(s *TOTPSecretProvider)) TOTPSecretProviderMocker { //nolint: revive
	return func(tb testing.TB) *TOTPSecretProvider {
		tb.Helper()

		s := NewTOTPSecretProvider(tb)

		for _, m := range mocks {
			m(s)
		}

		return s
	}
}
