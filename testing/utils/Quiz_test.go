package utils

import (
	"testing"

	"github.com/XzerozZ/Kasian_Phrom_BE/pkg/utils"
)

func TestCalculateRisk(t *testing.T) {
	testCases := []struct {
		name       string
		weights    []int
		expected   int
		shouldFail bool
	}{
		{
			name:     "Risk Level 1 (Low Risk)",
			weights:  []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expected: 1,
		},
		{
			name:     "Risk Level 2 (Lower Medium Risk)",
			weights:  []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
			expected: 2,
		},
		{
			name:     "Risk Level 3 (Medium Risk)",
			weights:  []int{3, 3, 3, 3, 3, 3, 3, 3, 2, 2},
			expected: 3,
		},
		{
			name:     "Risk Level 4 (Higher Medium Risk)",
			weights:  []int{4, 4, 4, 4, 3, 3, 3, 3, 3, 3},
			expected: 4,
		},
		{
			name:     "Risk Level 5 (High Risk)",
			weights:  []int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5},
			expected: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk, err := utils.CalculateRisk(tc.weights)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if risk != tc.expected {
				t.Errorf("Expected risk level %d, got %d", tc.expected, risk)
			}
		})
	}
}

func TestCalculateRiskEdgeCases(t *testing.T) {
	t.Run("Insufficient Weights", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for insufficient weights")
			}
		}()

		_, _ = utils.CalculateRisk([]int{1, 2, 3})
	})
}

func BenchmarkCalculateRisk(b *testing.B) {
	weights := []int{3, 3, 3, 3, 3, 3, 3, 3, 3, 3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utils.CalculateRisk(weights)
	}
}
