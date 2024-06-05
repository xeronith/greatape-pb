package utility

import (
	"fmt"
	"math/rand"
)

func GenerateOTP() string {
	// Generate a random number between 100000 and 999999
	otp := rand.Intn(900000) + 100000

	// Convert the OTP to a string
	return fmt.Sprintf("%06d", otp)
}
