package helpers

func ComputeBlockingProbabilities(assignedConnections, totalConnections int) int {

	return 1 - (assignedConnections / totalConnections)
}
