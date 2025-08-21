package simulator

import (
	"fmt"
	"sort"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections/controller"
	randomvariable "github.com/Kayres21/optical-mb-sim-go/internal/connections/random_variable"
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
	AssignedConnections  int
	TotalConnections     int
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

func (s *Simulator) AddNewConnectionEvent(event connections.ConnectionEvent) {
	events := s.GetConnectionsEvents()

	events = append(events, event)

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time < events[j].Time
	})

	s.SetConnectionsEvents(events)
}

func (s *Simulator) GetFirstEvent() connections.ConnectionEvent {
	connectionsEvents := s.GetConnectionsEvents()

	if len(connectionsEvents) == 0 {
		fmt.Println("No more events to process.")
		return connections.ConnectionEvent{}
	}

	// Tomar el primer evento
	firstElement := connectionsEvents[0]

	// Actualizar la lista quitando el primer elemento
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

func (s *Simulator) printBlockingTable(i int, blockingProbabilitie float64) {
	// Calcula el 10% del valor de s.GetGoalConnections()
	// La división de enteros trunca el resultado, que es lo que queremos en este caso.
	step := int(s.GetGoalConnections()) / 10

	// Imprime la tabla solo si i es un múltiplo del 'step'
	if i > 0 && i%step == 0 {
		fmt.Println("--------------------------------")
		fmt.Println("Tabla de Bloqueo")
		fmt.Println("--------------------------------")
		fmt.Printf("Conexiones: %d\n", i)
		fmt.Printf("Probabilidad de bloqueo: %.4f\n", blockingProbabilitie)
		fmt.Println("--------------------------------")
		fmt.Println() // Salto de línea para mejor legibilidad
	}
}

func (s *Simulator) SimulatorInit(networkPath string, routesPath string, capacitiesPath string, bitRatePath string, lambda int, mu int, goalConnections float64, allocator allocator.Allocator, numberOfBands int) {

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

	gigabits := 5

	s.SetNumberOfBands(band)

	randomVariable.SetParameters(lambda, mu, bitrate, source, destination, band, gigabits)

	seedArrive := int64(1)
	seedDeparture := int64(2)
	seedBitrate := int64(3)
	seedSource := int64(4)
	seedDestination := int64(5)
	seedBand := int64(6)
	seedGigabits := int64(7)

	randomVariable.SetSeeds(seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits)

	s.SetRandomVariable(randomVariable)

	s.SetGoalConnection(goalConnections)

	s.SetTime(0)

	connectionsEvents := connections.GenerateEvents(node_len, randomVariable)

	s.SetConnectionsEvents(connectionsEvents)

	var controller controller.Controller
	controller.ControllerInit(routesPath, s.Network, allocator)

	s.SetController(controller)

}

func (s *Simulator) SimulatorStart() {

	fmt.Println("Starting simulation...")

	for i := 1; i <= int(s.GetGoalConnections()); i++ {
		time := s.GetTime()

		event := s.GetFirstEvent()
		blockingProbabilitie := helpers.ComputeBlockingProbabilities(s.GetAssignedConnections(), s.GetTotalConnections())
		s.printBlockingTable(i, blockingProbabilitie)
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
			controller := s.GetController()

			controller.ReleaseConnection(event.Id)

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

		s.SetTime(time + event.Time)
	}

	fmt.Println("Simulation completed.")

}
