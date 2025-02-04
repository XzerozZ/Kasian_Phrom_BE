package utils

import (
	"errors"
)

func CalculateRisk(weight []int) (int, error) {
	totalScore := 0
	for i := 0; i < 10; i++ {
		totalScore += weight[i]
	}

	switch {
	case totalScore < 15:
		return 1, nil

	case totalScore >= 15 && totalScore <= 21:
		return 2, nil

	case totalScore >= 22 && totalScore <= 29:
		return 3, nil

	case totalScore >= 30 && totalScore <= 36:
		return 4, nil

	case totalScore >= 37:
		return 5, nil

	default:
		return 0, errors.New("unable to calculate risk")
	}
}
