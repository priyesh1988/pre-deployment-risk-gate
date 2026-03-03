package engine

import (
	"math"
)

func Score(violations int, env string) (float64, string) {
	mult := 1.0
	if env == "prod" {
		mult = 1.2
	}
	total := float64(violations) * mult
	score := 100 * (1 - math.Exp(-total))

	tier := "LOW"
	if score > 80 {
		tier = "CRITICAL"
	} else if score > 60 {
		tier = "HIGH"
	} else if score > 30 {
		tier = "MEDIUM"
	}
	return score, tier
}
