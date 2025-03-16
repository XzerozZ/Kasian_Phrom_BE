package utils_test

import (
	"os"
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"
	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeEmail(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "valid email with dots",
			input:    "John.Doe@Example.com",
			expected: "johndoe@example.com",
			hasError: false,
		},
		{
			name:     "valid email with multiple dots",
			input:    "john.middle.doe@example.com",
			expected: "johnmiddledoe@example.com",
			hasError: false,
		},
		{
			name:     "mixed case email",
			input:    "JohnDoe@ExAmPlE.cOm",
			expected: "johndoe@example.com",
			hasError: false,
		},
		{
			name:     "invalid email - no @",
			input:    "invalid-email",
			expected: "",
			hasError: true,
		},
		{
			name:     "invalid email - multiple @",
			input:    "john@doe@example.com",
			expected: "",
			hasError: true,
		},
		{
			name:     "empty email",
			input:    "",
			expected: "",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalized, err := utils.NormalizeEmail(tc.input)

			if tc.hasError {
				assert.Error(t, err)
				assert.Equal(t, "", normalized)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, normalized)
			}
		})
	}
}

func TestSendMail(t *testing.T) {
	if os.Getenv("SKIP_EMAIL_TEST") != "" {
		t.Skip("Skipping email sending test")
	}

	t.Run("invalid template path", func(t *testing.T) {
		user := entities.User{
			Username: "testuser",
			Email:    "test@example.com",
		}
		config := configs.Mail{
			Sender: "sender@example.com",
			Host:   "smtp.example.com",
			Port:   "587",
			Key:    "testkey",
		}

		err := utils.SendMail("/path/to/non/existent/template.html", &user, "123456", config)
		assert.Error(t, err)
	})

	t.Run("invalid port configuration", func(t *testing.T) {
		user := entities.User{
			Username: "testuser",
			Email:    "test@example.com",
		}
		config := configs.Mail{
			Sender: "sender@example.com",
			Host:   "smtp.example.com",
			Port:   "invalid-port",
			Key:    "testkey",
		}

		err := utils.SendMail("", &user, "123456", config)
		assert.Error(t, err)
	})
}
