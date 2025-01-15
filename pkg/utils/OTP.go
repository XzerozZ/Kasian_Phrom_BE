package utils

import (
	"time"
	"errors"
	"math/rand"
)

func GenerateRandomOTP(length int, useDigits bool) (string, error) {
    if length <= 0 {
		return "", errors.New("OTP length must be greater than 0")
	}

	var chars string
    if useDigits {
        chars = "0123456789"
    } else {
        chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" // Use letters and digits
    }

    rand.Seed(time.Now().UnixNano())
    b := make([]byte, length)
    for i := range b {
        b[i] = chars[rand.Intn(len(chars))]
    }

    return string(b), nil
}
