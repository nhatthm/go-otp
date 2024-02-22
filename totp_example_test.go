package otp_test

import (
	"context"
	"fmt"
	"time"

	"go.nhat.io/clock"

	"go.nhat.io/otp"
)

func ExampleGenerateTOTP() {
	c := clock.Fix(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC))

	result, err := otp.GenerateTOTP(context.Background(), otp.TOTPSecret("NBSWY3DP"), otp.WithClock(c))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	// Output:
	// 191882
}
