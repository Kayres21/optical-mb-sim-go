package randomvariable

import "math"

func (rv *RandomVariable) GetNetValueExponential(key ExponentialKey) float64 {
	switch key {
	case KeyArrive:
		return -math.Log(1-rv.Arrive.Rng.Float64()) / rv.Arrive.Parameter
	case KeyDeparture:
		return -math.Log(1-rv.Departure.Rng.Float64()) / rv.Departure.Parameter
	}
	return 0
}
