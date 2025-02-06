package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
)

func GenerateRandomOTP(length int, useDigits bool) (string, error) {
	if length <= 0 {
		return "", errors.New("OTP length must be greater than 0")
	}

	var chars string
	if useDigits {
		chars = "0123456789"
	} else {
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	}

	b := make([]byte, length)
	charsLength := big.NewInt(int64(len(chars)))

	for i := range b {
		n, err := rand.Int(rand.Reader, charsLength)
		if err != nil {
			return "", err
		}
		b[i] = chars[n.Int64()]
	}

	return string(b), nil
}
