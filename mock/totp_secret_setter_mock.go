package mock

import "testing"

// TOTPSecretSetterMocker is TOTPSecretSetter mocker.
type TOTPSecretSetterMocker func(tb testing.TB) *TOTPSecretSetter

// NopTOTPSecretSetter is no mock TOTPSecretSetter.
var NopTOTPSecretSetter = MockTOTPSecretSetter()

// MockTOTPSecretSetter creates TOTPSecretSetter mock with cleanup to ensure all the expectations are met.
func MockTOTPSecretSetter(mocks ...func(s *TOTPSecretSetter)) TOTPSecretSetterMocker { //nolint: revive
	return func(tb testing.TB) *TOTPSecretSetter {
		tb.Helper()

		s := NewTOTPSecretSetter(tb)

		for _, m := range mocks {
			m(s)
		}

		return s
	}
}
