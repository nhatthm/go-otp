package mock

import "testing"

// TOTPSecretDeleterMocker is TOTPSecretDeleter mocker.
type TOTPSecretDeleterMocker func(tb testing.TB) *TOTPSecretDeleter

// NopTOTPSecretDeleter is no mock TOTPSecretDeleter.
var NopTOTPSecretDeleter = MockTOTPSecretDeleter()

// MockTOTPSecretDeleter creates TOTPSecretDeleter mock with cleanup to ensure all the expectations are met.
func MockTOTPSecretDeleter(mocks ...func(g *TOTPSecretDeleter)) TOTPSecretDeleterMocker { //nolint: revive
	return func(tb testing.TB) *TOTPSecretDeleter {
		tb.Helper()

		g := NewTOTPSecretDeleter(tb)

		for _, m := range mocks {
			m(g)
		}

		return g
	}
}
