package randomvariable

import (
	"math/rand"
)

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
