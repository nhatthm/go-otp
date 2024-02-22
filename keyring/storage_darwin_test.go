//go:build integration && darwin

package keyring_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	oskey "github.com/zalando/go-keyring"

	"go.nhat.io/otp/keyring"
)

func TestStorage_Set_ErrDataTooBig(t *testing.T) {
	t.Parallel()

	const key = "key"

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.ErrorIs(t, err, oskey.ErrNotFound)
	})

	data := strings.Repeat("0", 4097)

	s := keyring.NewStorage[string]()

	err := s.Set(t.Name(), key, data)

	require.EqualError(t, err, "failed to write data to keyring: data passed to Set was too big")
}
