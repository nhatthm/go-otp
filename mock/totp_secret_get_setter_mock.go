package mock

import "testing"

// TOTPSecretGetSetterMocker is TOTPSecretGetSetter mocker.
type TOTPSecretGetSetterMocker func(tb testing.TB) *TOTPSecretGetSetter

// NopTOTPSecretGetSetter is no mock TOTPSecretGetSetter.
var NopTOTPSecretGetSetter = MockTOTPSecretGetSetter()

// MockTOTPSecretGetSetter creates TOTPSecretGetSetter mock with cleanup to ensure all the expectations are met.
func MockTOTPSecretGetSetter(mocks ...func(s *TOTPSecretGetSetter)) TOTPSecretGetSetterMocker { //nolint: revive
	return func(tb testing.TB) *TOTPSecretGetSetter {
		tb.Helper()

		s := NewTOTPSecretGetSetter(tb)

		for _, m := range mocks {
			m(s)
		}

		return s
	}
}
