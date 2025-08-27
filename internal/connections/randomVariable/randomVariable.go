package randomvariable

import (
	"fmt"
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

func (rv *RandomVariable) SetGigabits(seed int64) {
	fuente := rand.NewSource(seed)

	rv.GigabitsSelected.SetRng(rand.New(fuente))
}

func (rv *RandomVariable) SetLambda(lambda int) {
	fmt.Println("lambda", lambda)
	rv.Arrive.SetParameter(lambda)
}

func (rv *RandomVariable) SetMu(mu int) {
	fmt.Println("Mu", mu)
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

func (rv *RandomVariable) SetGigabitsSelectedParameter(parameter int) {
	rv.GigabitsSelected.SetParameter(parameter)
}

func (rv *RandomVariable) SetSeeds(seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits int64) {
	rv.SetSeedArrive(seedArrive)
	rv.SetSeedDeparture(seedDeparture)
	rv.SetSeedBitrate(seedBitrate)
	rv.SetSeedSource(seedSource)
	rv.SetSeedDestination(seedDestination)
	rv.SetSeedBand(seedBand)
	rv.SetGigabits(seedGigabits)

}

func (rv *RandomVariable) SetParameters(lambda, mu, bitrateSelect, sourceNodeSelect, destinationNodeSelect, bandSelect, gigabits int) {
	rv.SetLambda(lambda)
	rv.SetMu(mu)
	rv.SetBitrateSelectParameter(bitrateSelect)
	rv.SetSourceNodeSelectParameter(sourceNodeSelect)
	rv.SetDestinationNodeSelectParameter(destinationNodeSelect)
	rv.SetBandSelectParameter(bandSelect)
	rv.SetGigabitsSelectedParameter(gigabits)
}
