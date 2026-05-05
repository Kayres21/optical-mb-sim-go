package helpers

func ComputeBlockingProbabilities(assignedConnections, totalConnections int) float64 {
	return 1.0 - float64(assignedConnections)/float64(totalConnections)
}
