package helpers

import (
	"math"
	"testing"
)

func TestWaldConfidenceInterval(t *testing.T) {
	expectedLower := 0.19009678930349883
	expectedUpper := 0.8099032106965012
	lower, upper := WaldConfidenceInterval(5, 10)

	if math.Abs(lower-expectedLower) > 1e-9 {
		t.Errorf("WaldConfidenceInterval lower = %v, want %v", lower, expectedLower)
	}
	if math.Abs(upper-expectedUpper) > 1e-9 {
		t.Errorf("WaldConfidenceInterval upper = %v, want %v", upper, expectedUpper)
	}
}

func TestAgrestiCoullConfidenceInterval(t *testing.T) {
	expectedLower := 0.23658959361548726
	expectedUpper := 0.7634104063845127
	lower, upper := AgrestiCoullConfidenceInterval(5, 10)

	if math.Abs(lower-expectedLower) > 1e-9 {
		t.Errorf("AgrestiCoullConfidenceInterval lower = %v, want %v", lower, expectedLower)
	}
	if math.Abs(upper-expectedUpper) > 1e-9 {
		t.Errorf("AgrestiCoullConfidenceInterval upper = %v, want %v", upper, expectedUpper)
	}
}

func TestWilsonConfidenceInterval(t *testing.T) {
	expectedLower := 0.2365895936154873
	expectedUpper := 0.7634104063845127
	lower, upper := WilsonConfidenceInterval(5, 10)

	if math.Abs(lower-expectedLower) > 1e-9 {
		t.Errorf("WilsonConfidenceInterval lower = %v, want %v", lower, expectedLower)
	}
	if math.Abs(upper-expectedUpper) > 1e-9 {
		t.Errorf("WilsonConfidenceInterval upper = %v, want %v", upper, expectedUpper)
	}
}
