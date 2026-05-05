package randomvariable

import (
	"math/rand"
)

// ExponentialKey identifies which exponential random variable to sample.
type ExponentialKey string

// UniformKey identifies which uniform random variable to sample.
type UniformKey string

const (
	KeyArrive    ExponentialKey = "arrive"
	KeyDeparture ExponentialKey = "departure"
)

const (
	KeyBitrate     UniformKey = "bitrate"
	KeySource      UniformKey = "source"
	KeyDestination UniformKey = "destination"
	KeyBand        UniformKey = "band"
	KeyGigabits    UniformKey = "gigabits"
)

// DefaultGigabitOptions lists the supported gigabit values for connection requests.
var DefaultGigabitOptions = []int{10, 40, 100, 400, 1000}

type RandomVariable struct {
	Arrive                ExponentialVariable
	Departure             ExponentialVariable
	BitrateSelect         UniformVariable
	SourceNodeSelect      UniformVariable
	DestinationNodeSelect UniformVariable
	BandSelect            UniformVariable
	GigabitsSelected      UniformVariable
}

type ExponentialVariable struct {
	Parameter int
	Rng       *rand.Rand
}

type UniformVariable struct {
	Parameter int
	Rng       *rand.Rand
}

func (rv *RandomVariable) SetSeeds(seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits int64) {
	rv.Arrive.Rng = rand.New(rand.NewSource(seedArrive))
	rv.Departure.Rng = rand.New(rand.NewSource(seedDeparture))
	rv.BitrateSelect.Rng = rand.New(rand.NewSource(seedBitrate))
	rv.SourceNodeSelect.Rng = rand.New(rand.NewSource(seedSource))
	rv.DestinationNodeSelect.Rng = rand.New(rand.NewSource(seedDestination))
	rv.BandSelect.Rng = rand.New(rand.NewSource(seedBand))
	rv.GigabitsSelected.Rng = rand.New(rand.NewSource(seedGigabits))
}

func (rv *RandomVariable) SetParameters(lambda, mu, bitrateSelect, sourceNodeSelect, destinationNodeSelect, bandSelect, gigabits int) {
	rv.Arrive.Parameter = lambda
	rv.Departure.Parameter = mu
	rv.BitrateSelect.Parameter = bitrateSelect
	rv.SourceNodeSelect.Parameter = sourceNodeSelect
	rv.DestinationNodeSelect.Parameter = destinationNodeSelect
	rv.BandSelect.Parameter = bandSelect
	rv.GigabitsSelected.Parameter = gigabits
}
