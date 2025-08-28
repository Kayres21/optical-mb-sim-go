package randomvariable

import (
	"math"
	"math/rand"
)

func (ev *ExponentialVariable) SetRng(rng *rand.Rand) {
	ev.Rng = rng
}

func (ev *ExponentialVariable) SetParameter(parameter int) {
	ev.Parameter = parameter
}

func (rv *RandomVariable) GetNetValueExponential(options string) float64 {

	if options == "arrive" {
		return -1 * float64(math.Log(1-rv.Arrive.Rng.Float64())/float64(rv.Arrive.Parameter))
	}
	if options == "departure" {
		return -1 * float64(math.Log(1-rv.Departure.Rng.Float64())/float64(rv.Departure.Parameter))
	}

	return float64(0)

}
