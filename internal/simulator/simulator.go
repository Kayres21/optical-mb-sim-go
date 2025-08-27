package simulator

import (
	"fmt"
	"sort"
	"time"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections/controller"
	randomvariable "github.com/Kayres21/optical-mb-sim-go/internal/connections/randomVariable"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
	"github.com/Kayres21/optical-mb-sim-go/pkg/helpers"
)

type Simulator struct {
	RandomVariable       randomvariable.RandomVariable
	Network              infrastructure.Network
	BitRateList          connections.BitRateList
	ConnectionsEvents    []connections.ConnectionEvent
	Controller           controller.Controller
	GoalConnections      float64
	Time                 float64
	AllocatedConnections []bool
	NumberOfBands        int
	NumberOfBitrates     int
	NumberOfNodes        int
	NumberOfGigabits     int
	AssignedConnections  int
	TotalConnections     int
	startTime            time.Time
	Results              []float64
	Arrives              []float64
}

func (s *Simulator) GetResults() []float64 {
	return s.Results
}

func (s *Simulator) SetResults(results []float64) {
	s.Results = results
}
func (s *Simulator) GetArrives() []float64 {
	return s.Arrives
}
func (s *Simulator) SetArrives(arrives []float64) {
	s.Arrives = arrives
}

func (s *Simulator) GetNumberOfGigabits() int {
	return s.NumberOfGigabits
}

func (s *Simulator) SetNumberOfGigabits(numberOfGigabits int) {
	s.NumberOfGigabits = numberOfGigabits
}

func (s *Simulator) GetNumberOfBitrates() int {
	return s.NumberOfBitrates
}

func (s *Simulator) SetNumberOfBitrates(numberOfBitrates int) {
	s.NumberOfBitrates = numberOfBitrates
}

func (s *Simulator) GetNumberOfNodes() int {
	return s.NumberOfNodes
}

func (s *Simulator) SetNumberOfNodes(numberOfNodes int) {
	s.NumberOfNodes = numberOfNodes
}

func (s *Simulator) GetAssignedConnections() int {
	return s.AssignedConnections
}

func (s *Simulator) SetAssignedConnections(assignedConnections int) {
	s.AssignedConnections = assignedConnections
}
func (s *Simulator) GetTotalConnections() int {
	return s.TotalConnections
}
func (s *Simulator) SetTotalConnections(totalConnections int) {
	s.TotalConnections = totalConnections
}

// Getter and Setter for AllocatedConnections
func (s *Simulator) GetAllocatedConnections() []bool {
	return s.AllocatedConnections
}

func (s *Simulator) SetAllocatedConnections(allocatedConnections []bool) {
	s.AllocatedConnections = allocatedConnections
}

// Getter and Setter for NumberOfBands
func (s *Simulator) GetNumberOfBands() int {
	return s.NumberOfBands
}

func (s *Simulator) SetNumberOfBands(numberOfBands int) {
	s.NumberOfBands = numberOfBands
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

func (s *Simulator) GetStartTime() time.Time {
	return s.startTime
}

func (s *Simulator) SetStartTime(time time.Time) {
	s.startTime = time
}

func (s *Simulator) AddNewConnectionEvent(event connections.ConnectionEvent) {
	events := s.GetConnectionsEvents()

	i := sort.Search(len(events), func(i int) bool {
		return events[i].Time >= event.Time
	})
	events = append(events, connections.ConnectionEvent{})
	copy(events[i+1:], events[i:])
	events[i] = event

	s.SetConnectionsEvents(events)
}

func (s *Simulator) GetFirstEvent() connections.ConnectionEvent {
	connectionsEvents := s.GetConnectionsEvents()

	if len(connectionsEvents) == 0 {
		fmt.Println("No more events to process.")
		return connections.ConnectionEvent{}
	}

	firstElement := connectionsEvents[0]

	s.SetConnectionsEvents(connectionsEvents[1:])

	return firstElement
}

func (s *Simulator) GetGoalConnections() float64 {
	return s.GoalConnections
}

func (s *Simulator) GetTime() float64 {
	return s.Time
}

func (s *Simulator) SetTime(time float64) {
	s.Time = time
}

func (s *Simulator) SetController(controller controller.Controller) {
	s.Controller = controller
}

func (s *Simulator) GetController() controller.Controller {
	return s.Controller
}

func (s *Simulator) getSlotgigabites(bitrate connections.BitRate, gigabites int) int {

	slots := bitrate.Slots

	for _, slot := range slots {

		if slot.Gigabits == fmt.Sprint(gigabites) {
			return slot.Slots
		}

	}

	return 0

}

func (s *Simulator) addResult(result float64) {
	s.Results = append(s.Results, result)
}

func (s *Simulator) addArrive(arrive float64) {
	s.Arrives = append(s.Arrives, arrive)
}

func (s *Simulator) printBlockingTable(i int, blockingProbability float64, logOn bool) {

	if logOn {
		if i == 0 {
			fmt.Println("+----------+----------+----------+----------+")
			fmt.Println("| progress |  arrives | blocking |  time(s) |")
			fmt.Println("+----------+----------+----------+----------+")
		}

		step := s.GetGoalConnections() / 10

		if step == 0 {
			step = 1
		}

		if i > 0 && i%int(step) == 0 {

			progress := (float64(i) / float64(s.GetGoalConnections())) * 100

			elapsedTime := time.Since(s.GetStartTime())

			hours := int(elapsedTime.Hours())
			minutes := int(elapsedTime.Minutes()) % 60
			seconds := int(elapsedTime.Seconds()) % 60
			timeFormatted := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

			fmt.Printf("|%8.1f %%|%10d|%10.6f|%10s|\n", progress, i, blockingProbability, timeFormatted)
			fmt.Println("+----------+----------+----------+----------+")
			s.addResult(blockingProbability)
			s.addArrive(float64(i))
		}
	}

}

func (s *Simulator) RandomVariableInit(lambda int, mu int, bitrate int, source int, destination int, band int, gigabits int) {

	var randomVariable randomvariable.RandomVariable

	randomVariable.SetParameters(lambda, mu, bitrate, source, destination, band, gigabits)

	seedArrive := int64(1)
	seedDeparture := int64(12)
	seedBitrate := int64(123)
	seedSource := int64(1234)
	seedDestination := int64(12345)
	seedBand := int64(123456)
	seedGigabits := int64(1234567)

	randomVariable.SetSeeds(seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits)

	s.SetRandomVariable(randomVariable)

}

func (s *Simulator) ControllerInit(routesPath string, network infrastructure.Network, allocator allocator.Allocator) {
	var controller controller.Controller
	controller.ControllerInit(routesPath, network, allocator)
	s.SetController(controller)
}

func (s *Simulator) connectionsEventsInit(nodeLen int, rv randomvariable.RandomVariable) {
	connectionsEvents := connections.GenerateEvents(nodeLen, rv)
	s.SetConnectionsEvents(connectionsEvents)
}

func (s *Simulator) NetworkInit(networkPath string, capacitiesPath string) {
	network, err := infrastructure.NetworkGenerate(networkPath, capacitiesPath)

	if err != nil {
		fmt.Printf("Error reading network file: %v\n", err)
	}
	fmt.Println("Network Name:", network.Name)

	s.SetNetwork(network)
}

func (s *Simulator) BitRateInit(bitRatePath string) {
	bitRate, err := connections.ReadBitRateFile(bitRatePath)

	if err != nil {
		fmt.Printf("Error reading bitrate file: %v\n", err)
	}

	s.SetBitRateList(bitRate)
}

func (s *Simulator) VariablesNumbersInit(numberOfBands int) {
	numberOfNodes := len(s.GetNetwork().Nodes)
	numberOfBitrates := len(s.GetBitRateList().BitRates)

	s.SetNumberOfNodes(numberOfNodes)
	s.SetNumberOfBitrates(numberOfBitrates)

	if numberOfBands > 4 {
		fmt.Println("Warning: Number of bands exceeds 4, setting to 4.")
		numberOfBands = 4
	} else if numberOfBands < 1 {
		fmt.Println("Warning: Number of bands is less than 1, setting to 1.")
		numberOfBands = 1
	}

	s.SetNumberOfBands(numberOfBands)

	s.SetNumberOfGigabits(5)
}

func (s *Simulator) SimulatorInit(networkPath string, routesPath string, capacitiesPath string, bitRatePath string, lambda int, mu int, goalConnections float64, allocator allocator.Allocator, numberOfBands int) {

	s.NetworkInit(networkPath, capacitiesPath)

	s.BitRateInit(bitRatePath)

	s.VariablesNumbersInit(numberOfBands)

	s.RandomVariableInit(lambda, mu, s.GetNumberOfBitrates(), s.GetNumberOfNodes(), s.GetNumberOfNodes(), s.GetNumberOfBands(), s.GetNumberOfGigabits())

	s.SetGoalConnection(goalConnections)

	s.SetTime(0)

	s.connectionsEventsInit(s.GetNumberOfNodes(), s.GetRandomVariable())

	s.ControllerInit(routesPath, s.GetNetwork(), allocator)

	s.SetStartTime(time.Now())

}

func (s *Simulator) SimulatorStart(logOn bool) {

	fmt.Println("Starting simulation...")

	for i := 0; i <= int(s.GetGoalConnections()); i++ {
		time := s.GetTime()

		event := s.GetFirstEvent()
		blockingProbabilitie := helpers.ComputeBlockingProbabilities(s.GetAssignedConnections(), s.GetTotalConnections())
		s.printBlockingTable(i, blockingProbabilitie, logOn)
		if event.Event == connections.ConnectionEventTypeArrive {
			controller := s.GetController()

			slot := s.getSlotgigabites(s.GetBitRateList().BitRates[event.Bitrate], event.GigabitsSelected)
			asigned, con := controller.ConectionAllocation(event.Source, event.Destination, slot, s.GetNetwork(), controller.Routes, s.GetNumberOfBands())
			s.SetTotalConnections(s.GetTotalConnections() + 1)

			if asigned {
				s.Controller.AddConnection(con)
				s.SetAssignedConnections(s.GetAssignedConnections() + 1)
				s.SetAllocatedConnections(append(s.GetAllocatedConnections(), true))

				rv := s.GetRandomVariable()

				newEvent := connections.ConnectionEvent{
					Id:          event.Id,
					Source:      event.Source,
					Destination: event.Destination,
					Bitrate:     con.BandSelected,
					Event:       connections.ConnectionEventTypeRelease,
					Time:        s.GetTime() + rv.GetNetValueExponential("departure"),
				}
				s.AddNewConnectionEvent(newEvent)

			}
			if !asigned {

				rv := s.GetRandomVariable()

				newEvent := connections.ConnectionEvent{
					Id:          event.Id,
					Source:      event.Source,
					Destination: event.Destination,
					Bitrate:     event.Bitrate,
					Event:       connections.ConnectionEventTypeArrive,
					Time:        s.GetTime() + rv.GetNetValueExponential("arrive"),
				}
				s.AddNewConnectionEvent(newEvent)
			}

		}

		if event.Event == connections.ConnectionEventTypeRelease {
			rv := s.GetRandomVariable()
			controller := s.GetController()

			controller.ReleaseConnection(event.Id)

			newEvent := connections.ConnectionEvent{
				Id:          event.Id,
				Source:      event.Source,
				Destination: event.Destination,
				Bitrate:     event.Bitrate,
				Event:       connections.ConnectionEventTypeArrive,
				Time:        s.GetTime() + rv.GetNetValueExponential("arrive"),
			}
			s.AddNewConnectionEvent(newEvent)

		}

		s.SetTime(time + event.Time)
	}

	fmt.Println("Simulation completed.")

}
