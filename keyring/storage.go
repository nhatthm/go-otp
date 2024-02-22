package keyring

import (
	"encoding"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"

	"go.nhat.io/otp"
)

var (
	// ErrNotFound is a not found error.
	ErrNotFound = keyring.ErrNotFound
	// ErrUnsupportedType is an unsupported type error.
	ErrUnsupportedType = errors.New("unsupported type")
)

var _ otp.Storage[any] = (*Storage[any])(nil)

// Storage is a storage implementation that uses the OS keyring.
type Storage[V any] struct{}

// Set sets the value for the given key.
func (k *Storage[V]) Set(service string, key string, value V) error {
	d, err := marshalData(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data for writing to keyring: %w", err)
	}

	if err := keyring.Set(service, key, d); err != nil {
		return fmt.Errorf("failed to write data to keyring: %w", err)
	}

	return nil
}

// Get gets the value for the given key.
func (k *Storage[V]) Get(service string, key string) (*V, error) {
	d, err := keyring.Get(service, key)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from keyring: %w", err)
	}

	var dest V

	if err := unmarshalData(d, &dest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data read from keyring: %w", err)
	}

	return &dest, nil
}

// Delete deletes the value for the given key.
func (k *Storage[V]) Delete(service string, key string) error {
	if err := keyring.Delete(service, key); err != nil {
		return fmt.Errorf("failed to delete data from keyring: %w", err)
	}

	return nil
}

// NewStorage creates a new Storage that uses the OS keyring.
func NewStorage[V any]() *Storage[V] {
	return &Storage[V]{}
}

func marshalData(v any) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil

	case []byte:
		return string(v), nil

	case encoding.TextMarshaler:
		b, err := v.MarshalText()
		if err != nil {
			return "", err //nolint: wrapcheck
		}

		return string(b), nil
	}

	return "", fmt.Errorf("%w: %T", ErrUnsupportedType, v)
}

func unmarshalData(v string, dest any) error {
	switch dest := dest.(type) {
	case *string:
		*dest = v

	case *[]byte:
		*dest = []byte(v)

	case encoding.TextUnmarshaler:
		return dest.UnmarshalText([]byte(v)) //nolint: wrapcheck

	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedType, dest)
	}

	return nil
}
