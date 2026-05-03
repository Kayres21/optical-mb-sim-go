package helpers

import (
	"math"
	"testing"
)

func TestComputeBlockingProbabilities(t *testing.T) {
	tests := []struct {
		name                string
		assignedConnections int
		totalConnections    int
		expected            float64
		checkNaN            bool
	}{
		{
			name:                "half assigned",
			assignedConnections: 5,
			totalConnections:    10,
			expected:            0.5,
			checkNaN:            false,
		},
		{
			name:                "all assigned",
			assignedConnections: 10,
			totalConnections:    10,
			expected:            0.0,
			checkNaN:            false,
		},
		{
			name:                "none assigned",
			assignedConnections: 0,
			totalConnections:    10,
			expected:            1.0,
			checkNaN:            false,
		},
		{
			name:                "zero total connections (div by zero)",
			assignedConnections: 5,
			totalConnections:    0,
			expected:            math.NaN(),
			checkNaN:            true, // in Go, non-zero / zero is +Inf or -Inf, but here 1 - Inf = -Inf, wait. If float64(5)/float64(0) = +Inf. 1 - Inf = -Inf.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeBlockingProbabilities(tt.assignedConnections, tt.totalConnections)
			if tt.checkNaN {
				// float64(assigned)/0 is +Inf, 1 - (+Inf) is -Inf.
				if !math.IsInf(got, -1) {
					t.Errorf("ComputeBlockingProbabilities() = %v, expected -Inf", got)
				}
			} else {
				if math.Abs(got-tt.expected) > 1e-9 {
					t.Errorf("ComputeBlockingProbabilities() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}
