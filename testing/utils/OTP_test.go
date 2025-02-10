package utils

import (
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
)

func TestGenerateRandomOTP(t *testing.T) {
	t.Run("Generate Digits OTP", func(t *testing.T) {
		otp, err := utils.GenerateRandomOTP(6, true)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(otp) != 6 {
			t.Errorf("Expected OTP length 6, got %d", len(otp))
		}

		for _, char := range otp {
			if char < '0' || char > '9' {
				t.Errorf("Expected only digits, got %c", char)
			}
		}
	})

	t.Run("Generate Alphanumeric OTP", func(t *testing.T) {
		otp, err := utils.GenerateRandomOTP(8, false)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(otp) != 8 {
			t.Errorf("Expected OTP length 8, got %d", len(otp))
		}

		validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		for _, char := range otp {
			found := false
			for _, validChar := range validChars {
				if char == validChar {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("Invalid character in OTP: %c", char)
			}
		}
	})

	t.Run("Zero Length OTP", func(t *testing.T) {
		_, err := utils.GenerateRandomOTP(0, true)
		if err == nil {
			t.Error("Expected error for zero length OTP, got nil")
		}
	})

	t.Run("Negative Length OTP", func(t *testing.T) {
		_, err := utils.GenerateRandomOTP(-5, true)
		if err == nil {
			t.Error("Expected error for negative length OTP, got nil")
		}
	})

	t.Run("Randomness Check", func(t *testing.T) {
		const iterations = 1000
		const otpLength = 6
		otpSet := make(map[string]bool)

		for i := 0; i < iterations; i++ {
			otp, err := utils.GenerateRandomOTP(otpLength, true)
			if err != nil {
				t.Errorf("Unexpected error in iteration %d: %v", i, err)
			}

			if otpSet[otp] {
				t.Logf("Duplicate OTP found: %s (might be unlikely but possible)", otp)
			}

			otpSet[otp] = true
		}

		if len(otpSet) < iterations/2 {
			t.Errorf("Low uniqueness in generated OTPs: %d out of %d", len(otpSet), iterations)
		}
	})
}

func BenchmarkGenerateRandomOTP(b *testing.B) {
	b.Run("Digits OTP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.GenerateRandomOTP(6, true)
		}
	})

	b.Run("Alphanumeric OTP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			utils.GenerateRandomOTP(8, false)
		}
	})
}
