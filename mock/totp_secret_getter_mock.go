package mock

import "testing"

// TOTPSecretGetterMocker is TOTPSecretGetter mocker.
type TOTPSecretGetterMocker func(tb testing.TB) *TOTPSecretGetter

// NopTOTPSecretGetter is no mock TOTPSecretGetter.
var NopTOTPSecretGetter = MockTOTPSecretGetter()

// MockTOTPSecretGetter creates TOTPSecretGetter mock with cleanup to ensure all the expectations are met.
func MockTOTPSecretGetter(mocks ...func(g *TOTPSecretGetter)) TOTPSecretGetterMocker { //nolint: revive
	return func(tb testing.TB) *TOTPSecretGetter {
		tb.Helper()

		g := NewTOTPSecretGetter(tb)

		for _, m := range mocks {
			m(g)
		}

		return g
	}
}
