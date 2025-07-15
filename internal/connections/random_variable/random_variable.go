package randomvariable

import (
	"math/rand"
)

type RandomVariable struct {
	Arrive                ExponentialVariable
	Departure             ExponentialVariable
	BitrateSelect         ExponentialVariable
	SourceNodeSelect      ExponentialVariable
	DestinationNodeSelect ExponentialVariable
	BandSelect            ExponentialVariable
}

type ExponentialVariable struct {
	Parameter int
	Rng       *rand.Rand
}
