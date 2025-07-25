package simulator

import (
	"fmt"
	"simulator/internal/connections"
	controller "simulator/internal/connections/controller"
	randomvariable "simulator/internal/connections/random_variable"
	"simulator/internal/infrastructure"
)

type Simulator struct {
	RandomVariable    randomvariable.RandomVariable
	Network           infrastructure.Network
	BitRateList       connections.BitRateList
	ConnectionsEvents []connections.ConnectionEvent
	Controller        controller.Controller
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

func (s *Simulator) SetController(controller controller.Controller) {
	s.Controller = controller
}

func (s *Simulator) GetController() controller.Controller {
	return s.Controller
}

func (s *Simulator) SimulatorInit(networkPath string, capacitiesPath string, bitRatePath string, lambda int, mu int, goalConnections float64, allocator controller.Allocator, numberOfBands int) {

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

	node_len := len(network.Nodes)

	bitrate := len(bitRate.BitRates)
	source := node_len
	destination := node_len

	if numberOfBands > 4 {
		fmt.Println("Warning: Number of bands exceeds 4, setting to 4.")
		numberOfBands = 4
	} else if numberOfBands < 1 {
		fmt.Println("Warning: Number of bands is less than 1, setting to 1.")
		numberOfBands = 1
	}

	band := numberOfBands

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

	s.SetTime(0)

	connectionsEvents := connections.GenerateEvents(node_len, randomVariable)

	s.SetConnectionsEvents(connectionsEvents)

	var controller controller.Controller
	controller.ControllerInit("files/routes/UKNet_routes.json", s.Network, allocator)

	s.SetController(controller)

}
