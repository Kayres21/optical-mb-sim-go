package randomvariable

import (
	"testing"
)

func TestRandomVariable_ParametersAndSeeds(t *testing.T) {
	var rv RandomVariable

	rv.SetParameters(10, 5, 3, 2, 2, 1, 4)
	if rv.Arrive.Parameter != 10 {
		t.Errorf("expected Lambda=10, got %v", rv.Arrive.Parameter)
	}

	rv.SetSeeds(1, 2, 3, 4, 5, 6, 7)
	// Just check if we can call it without panic, we can't easily inspect the rng state directly 
	// unless we generate values and see they are deterministic.

	val1 := rv.GetNetValueExponential("arrive")
	
	// Reset seed and check if deterministic
	rv.SetSeeds(1, 2, 3, 4, 5, 6, 7)
	val2 := rv.GetNetValueExponential("arrive")
	
	if val1 != val2 {
		t.Errorf("expected deterministic exponential values for same seed, got %v and %v", val1, val2)
	}
}

func TestRandomVariable_GetNetValueUniform(t *testing.T) {
	var rv RandomVariable
	rv.SetParameters(10, 5, 3, 2, 2, 1, 4)
	rv.SetSeeds(1, 2, 3, 4, 5, 6, 7)

	for i := 0; i < 100; i++ {
		val := rv.GetNetValueUniform("bitrate")
		if val < 0 || val >= 3 {
			t.Errorf("bitrate uniform value out of bounds: %v", val)
		}
	}
}
