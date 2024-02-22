//go:build integration

package keyring_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	oskey "github.com/zalando/go-keyring"

	"go.nhat.io/otp/keyring"
)

func TestStorage_Get_SecretNotFound(t *testing.T) {
	t.Parallel()

	s := keyring.NewStorage[chan struct{}]()

	r, err := s.Get(t.Name(), "key")

	assert.Nil(t, r)
	require.EqualError(t, err, "failed to read data from keyring: secret not found in keyring")
}

func TestStorage_Get_UnsupportedType(t *testing.T) {
	t.Parallel()

	const key = "key"

	err := oskey.Set(t.Name(), key, "value")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[chan struct{}]()

	r, err := s.Get(t.Name(), key)

	assert.Nil(t, r)
	require.EqualError(t, err, "failed to unmarshal data read from keyring: unsupported type: *chan struct {}")
}

func TestStorage_Get_Success_String(t *testing.T) {
	t.Parallel()

	const key = "key"

	err := oskey.Set(t.Name(), key, "value")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[string]()

	actual, err := s.Get(t.Name(), key)
	require.NoError(t, err)

	expected := "value"

	assert.Equal(t, expected, *actual)
}

func TestStorage_Get_Success_ByteSlice(t *testing.T) {
	t.Parallel()

	const key = "key"

	err := oskey.Set(t.Name(), key, "value")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[[]byte]()

	actual, err := s.Get(t.Name(), key)
	require.NoError(t, err)

	expected := []byte("value")

	assert.Equal(t, expected, *actual)
}

func TestStorage_Get_Success_TextUnmarshaler(t *testing.T) {
	t.Parallel()

	const key = "key"

	err := oskey.Set(t.Name(), key, "42")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[custom]()

	actual, err := s.Get(t.Name(), key)
	require.NoError(t, err)

	expected := custom(42)

	assert.Equal(t, expected, *actual)
}

func TestStorage_Get_Failure_TextUnmarshaler(t *testing.T) {
	t.Parallel()

	const key = "key"

	err := oskey.Set(t.Name(), key, "value")
	require.NoError(t, err)

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[custom]()

	actual, err := s.Get(t.Name(), key)

	require.Nil(t, actual)
	require.EqualError(t, err, `failed to unmarshal data read from keyring: strconv.Atoi: parsing "value": invalid syntax`)
}

func TestStorage_Set_UnsupportedType(t *testing.T) {
	t.Parallel()

	const key = "key"

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.ErrorIs(t, err, oskey.ErrNotFound)
	})

	s := keyring.NewStorage[chan struct{}]()

	err := s.Set(t.Name(), key, make(chan struct{}))

	require.EqualError(t, err, "failed to marshal data for writing to keyring: unsupported type: chan struct {}")
}

func TestStorage_Set_Success_String(t *testing.T) {
	t.Parallel()

	const key = "key"

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[string]()

	err := s.Set(t.Name(), key, "value")
	require.NoError(t, err)
}

func TestStorage_Set_Success_ByteSlice(t *testing.T) {
	t.Parallel()

	const key = "key"

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[[]byte]()

	err := s.Set(t.Name(), key, []byte("value"))
	require.NoError(t, err)
}

func TestStorage_Set_Success_TextMarshaler(t *testing.T) {
	t.Parallel()

	const key = "key"

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.NoError(t, err)
	})

	s := keyring.NewStorage[custom]()

	err := s.Set(t.Name(), key, custom(42))
	require.NoError(t, err)
}

func TestStorage_Set_Failure_TextMarshaler(t *testing.T) {
	t.Parallel()

	const key = "key"

	t.Cleanup(func() {
		err := oskey.Delete(t.Name(), key)
		require.ErrorIs(t, err, oskey.ErrNotFound)
	})

	s := keyring.NewStorage[custom]()

	err := s.Set(t.Name(), key, custom(-1))
	require.EqualError(t, err, `failed to marshal data for writing to keyring: negative value`)
}

func TestStorage_Delete_SecretNotFound(t *testing.T) {
	t.Parallel()

	s := keyring.NewStorage[string]()

	err := s.Delete(t.Name(), "key")

	require.EqualError(t, err, "failed to delete data from keyring: secret not found in keyring")
}

func TestStorage_Delete_Success(t *testing.T) {
	t.Parallel()

	const key = "key"

	err := oskey.Set(t.Name(), key, "value")
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err = oskey.Get(t.Name(), key)
		require.ErrorIs(t, err, oskey.ErrNotFound, "secret was not deleted from keyring")
	})

	s := keyring.NewStorage[any]()

	err = s.Delete(t.Name(), "key")
	require.NoError(t, err)
}

type custom int

func (c custom) MarshalText() (text []byte, err error) {
	if c < 0 {
		return nil, errors.New("negative value")
	}

	return []byte(strconv.Itoa(int(c))), nil
}

func (c *custom) UnmarshalText(text []byte) error {
	r, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}

	*c = custom(r)

	return nil
}
