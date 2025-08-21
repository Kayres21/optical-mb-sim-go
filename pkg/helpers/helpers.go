package helpers

func ComputeBlockingProbabilities(assignedConnections, totalConnections int) float64 {

	result := float64(1.0 - (float64(assignedConnections) / float64(totalConnections)))

	return result
}
