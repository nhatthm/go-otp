package mock

import "testing"

// GeneratorMocker is Generator mocker.
type GeneratorMocker func(tb testing.TB) *Generator

// NopGenerator is no mock Generator.
var NopGenerator = MockGenerator()

// MockGenerator creates Generator mock with cleanup to ensure all the expectations are met.
func MockGenerator(mocks ...func(g *Generator)) GeneratorMocker { //nolint: revive
	return func(tb testing.TB) *Generator {
		tb.Helper()

		g := NewGenerator(tb)

		for _, m := range mocks {
			m(g)
		}

		return g
	}
}
