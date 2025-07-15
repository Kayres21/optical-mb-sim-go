package randomvariable

import (
	"math"
	"math/rand"
)

func (rv *RandomVariable) SetSeedArrive(seed int64) {

	fuente := rand.NewSource(seed)

	rv.Arrive.SetRng(rand.New(fuente))

}

func (rv *RandomVariable) SetSeedDeparture(seed int64) {

	fuente := rand.NewSource(seed)

	rv.Departure.SetRng(rand.New(fuente))

}

func (rv *RandomVariable) SetSeedBitrate(seed int64) {

	fuente := rand.NewSource(seed)

	rv.BitrateSelect.SetRng(rand.New(fuente))

}

func (rv *RandomVariable) SetSeedSource(seed int64) {

	fuente := rand.NewSource(seed)

	rv.SourceNodeSelect.SetRng(rand.New(fuente))

}

func (rv *RandomVariable) SetSeedDestination(seed int64) {

	fuente := rand.NewSource(seed)

	rv.DestinationNodeSelect.SetRng(rand.New(fuente))

}

func (rv *RandomVariable) SetSeedBand(seed int64) {

	fuente := rand.NewSource(seed)

	rv.BandSelect.SetRng(rand.New(fuente))

}

func (rv *RandomVariable) SetLambda(lambda int) {
	rv.Arrive.SetParameter(lambda)
}

func (rv *RandomVariable) SetMu(mu int) {
	rv.Departure.SetParameter(mu)
}

func (rv *RandomVariable) SetBitrateSelectParameter(parameter int) {
	rv.BitrateSelect.SetParameter(parameter)
}

func (rv *RandomVariable) SetSourceNodeSelectParameter(parameter int) {
	rv.SourceNodeSelect.SetParameter(parameter)
}
func (rv *RandomVariable) SetDestinationNodeSelectParameter(parameter int) {
	rv.DestinationNodeSelect.SetParameter(parameter)
}

func (rv *RandomVariable) SetBandSelectParameter(parameter int) {
	rv.BandSelect.SetParameter(parameter)
}

func (ev *ExponentialVariable) SetRng(rng *rand.Rand) {
	ev.Rng = rng
}

func (ev *ExponentialVariable) SetParameter(parameter int) {
	ev.Parameter = parameter
}

func (ev *UniformVariable) SetRng(rng *rand.Rand) {
	ev.Rng = rng
}

func (ev *UniformVariable) SetParameter(parameter int) {
	ev.Parameter = parameter
}

func (rv *RandomVariable) GetNetValueExponential(options string) float64 {

	if options == "arrive" {
		return -1 * float64(math.Log(rv.Arrive.Rng.Float64())/float64(rv.Arrive.Parameter))
	}
	return -1 * float64(math.Log(rv.Departure.Rng.Float64())/float64(rv.Arrive.Parameter))

}

func (rv *RandomVariable) GetNetValueUniform(options string) int {
	switch options {
	case "bitrate":
		return rv.BitrateSelect.Rng.Intn(rv.BitrateSelect.Parameter)
	case "source":
		return rv.SourceNodeSelect.Rng.Intn(rv.SourceNodeSelect.Parameter)
	case "destination":
		return rv.DestinationNodeSelect.Rng.Intn(rv.DestinationNodeSelect.Parameter)
	case "band":
		return rv.BandSelect.Rng.Intn(rv.BandSelect.Parameter)
	}
	return -1 // Default case, should not happen
}
