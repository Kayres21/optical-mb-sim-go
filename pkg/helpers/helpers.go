package helpers

import "math"

const confidence95z = 1.96

func ComputeBlockingProbabilities(assignedConnections, totalConnections int) float64 {
	return 1.0 - float64(assignedConnections)/float64(totalConnections)
}

func WaldConfidenceInterval(assignedConnections, totalConnections int) (float64, float64) {
	if totalConnections == 0 {
		return math.NaN(), math.NaN()
	}

	blocked := float64(totalConnections-assignedConnections) / float64(totalConnections)
	se := math.Sqrt(blocked * (1.0 - blocked) / float64(totalConnections))
	lower := blocked - confidence95z*se
	upper := blocked + confidence95z*se
	return clamp(lower, 0.0, 1.0), clamp(upper, 0.0, 1.0)
}

func AgrestiCoullConfidenceInterval(assignedConnections, totalConnections int) (float64, float64) {
	if totalConnections == 0 {
		return math.NaN(), math.NaN()
	}

	blocked := float64(totalConnections - assignedConnections)
	n := float64(totalConnections)
	z2Even := 4.0
	pHat := (blocked + z2Even/2.0) / (n + z2Even)
	se := math.Sqrt(pHat * (1.0 - pHat) / (n + z2Even))
	lower := pHat - confidence95z*se
	upper := pHat + confidence95z*se
	return clamp(lower, 0.0, 1.0), clamp(upper, 0.0, 1.0)
}

func WilsonConfidenceInterval(assignedConnections, totalConnections int) (float64, float64) {
	if totalConnections == 0 {
		return math.NaN(), math.NaN()
	}

	blocked := float64(totalConnections-assignedConnections) / float64(totalConnections)
	n := float64(totalConnections)
	z2 := confidence95z * confidence95z
	center := (blocked + z2/(2.0*n)) / (1.0 + z2/n)
	margin := confidence95z * math.Sqrt((blocked*(1.0-blocked)+z2/(4.0*n))/n) / (1.0 + z2/n)

	lower := center - margin
	upper := center + margin
	return clamp(lower, 0.0, 1.0), clamp(upper, 0.0, 1.0)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
