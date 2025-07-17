package simulator

import (
	"fmt"
	"simulator/internal/connections"
	randomvariable "simulator/internal/connections/random_variable"
	"simulator/internal/infrastructure"
)

type Simulator struct {
	RandomVariable  randomvariable.RandomVariable
	Network         infrastructure.Network
	BitRateList     connections.BitRateList
	GoalConnections int
}

func (s *Simulator) SetRandomVariable(randomVariable randomvariable.RandomVariable) {
	s.RandomVariable = randomVariable
}

func (s *Simulator) SetNetwork(network infrastructure.Network) {
	s.Network = network
}

func (s *Simulator) SetBitRateList(bitRateList connections.BitRateList) {
	s.BitRateList = bitRateList
}

func (s *Simulator) SetGoalConnection(goalConnections int) {
	s.GoalConnections = goalConnections
}

func (s *Simulator) SimulatorInit(networkPath string, capacitiesPath string, bitRatePath string, lambda int, mu int, goalConnections int) Simulator {

	network, err := infrastructure.NetworkGenerate(networkPath, capacitiesPath)

	if err != nil {
		fmt.Printf("Error reading network file: %v\n", err)
		return Simulator{}
	}
	fmt.Println("Network Name:", network.Name)

	s.SetNetwork(network)

	bitRate, err := connections.ReadBitRateFile(bitRatePath)

	if err != nil {
		fmt.Printf("Error reading bitrate file: %v\n", err)
		return Simulator{}
	}

	s.SetBitRateList(bitRate)

	var randomVariable randomvariable.RandomVariable

	bitrate := len(bitRate.BitRates)
	source := len(network.Nodes)
	destination := len(network.Nodes)
	band := len(network.Links[0].Capacities.Bands)

	randomVariable.SetParameters(lambda, mu, bitrate, source, destination, band)

	seedArrive := int64(1)
	seedDeparture := int64(2)
	seedBitrate := int64(3)
	seedSource := int64(4)
	seedDestination := int64(5)
	seedBand := int64(6)

	randomVariable.SetSeeds(seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand)

	s.SetRandomVariable(randomVariable)

	s.SetGoalConnection(goalConnections)

	return *s
}
