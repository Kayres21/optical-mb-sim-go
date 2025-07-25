package simulator

import (
	"fmt"
	"simulator/internal/connections"
	randomvariable "simulator/internal/connections/random_variable"
	"simulator/internal/infrastructure"
)

type Simulator struct {
	RandomVariable    randomvariable.RandomVariable
	Network           infrastructure.Network
	BitRateList       connections.BitRateList
	ConnectionsEvents []connections.ConnectionEvent
	GoalConnections   float64
	time              float64
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

func (s *Simulator) SetGoalConnection(goalConnections float64) {
	s.GoalConnections = goalConnections
}

func (s *Simulator) SetConnectionsEvents(connectionsEvents []connections.ConnectionEvent) {
	s.ConnectionsEvents = connectionsEvents
}

func (s *Simulator) GetRandomVariable() randomvariable.RandomVariable {
	return s.RandomVariable
}

func (s *Simulator) GetNetwork() infrastructure.Network {
	return s.Network
}

func (s *Simulator) GetBitRateList() connections.BitRateList {
	return s.BitRateList
}

func (s *Simulator) GetConnectionsEvents() []connections.ConnectionEvent {
	return s.ConnectionsEvents
}

func (s *Simulator) GetGoalConnections() float64 {
	return s.GoalConnections
}

func (s *Simulator) GetTime() float64 {
	return s.time
}

func (s *Simulator) SetTime(time float64) {
	s.time = time
}

func (s *Simulator) SimulatorInit(networkPath string, capacitiesPath string, bitRatePath string, lambda int, mu int, goalConnections float64) {

	network, err := infrastructure.NetworkGenerate(networkPath, capacitiesPath)

	if err != nil {
		fmt.Printf("Error reading network file: %v\n", err)
	}
	fmt.Println("Network Name:", network.Name)

	s.SetNetwork(network)

	bitRate, err := connections.ReadBitRateFile(bitRatePath)

	if err != nil {
		fmt.Printf("Error reading bitrate file: %v\n", err)
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
	s.RandomVariable.GetNetValueUniform("source")

	s.SetTime(0)
}
