package otp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.nhat.io/otp"
)

func TestOTP_String(t *testing.T) {
	t.Parallel()

	o := otp.OTP("123456")

	assert.Equal(t, "123456", o.String())
}
