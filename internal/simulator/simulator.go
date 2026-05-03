package simulator

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections/controller"
	randomvariable "github.com/Kayres21/optical-mb-sim-go/internal/connections/randomVariable"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
	"github.com/Kayres21/optical-mb-sim-go/pkg/helpers"
	"github.com/Kayres21/optical-mb-sim-go/pkg/plotter"
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

func (s *Simulator) AddNewConnectionEvent(event connections.ConnectionEvent) {
	i := sort.Search(len(s.ConnectionsEvents), func(i int) bool {
		return s.ConnectionsEvents[i].Time >= event.Time
	})
	s.ConnectionsEvents = append(s.ConnectionsEvents, connections.ConnectionEvent{})
	copy(s.ConnectionsEvents[i+1:], s.ConnectionsEvents[i:])
	s.ConnectionsEvents[i] = event
}

func (s *Simulator) GetFirstEvent() connections.ConnectionEvent {
	if len(s.ConnectionsEvents) == 0 {
		fmt.Println("No more events to process.")
		return connections.ConnectionEvent{}
	}

	firstElement := s.ConnectionsEvents[0]
	s.ConnectionsEvents = s.ConnectionsEvents[1:]

	return firstElement
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

func (s *Simulator) printBlockingTable(i int, logOn bool) {

	if logOn {
		if i == 0 {
			fmt.Println("+----------+----------+----------+----------+")
			fmt.Println("| progress |  arrives | blocking |  time(s) |")
			fmt.Println("+----------+----------+----------+----------+")
		}

		step := s.GoalConnections / 10

		if step == 0 {
			step = 1
		}

		if i > 0 && i%int(step) == 0 {
			blockingProbability := helpers.ComputeBlockingProbabilities(s.AssignedConnections, s.TotalConnections)
			progress := (float64(i) / float64(s.GoalConnections)) * 100

			elapsedTime := time.Since(s.startTime)

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

	s.RandomVariable = randomVariable

}

func (s *Simulator) connectionsEventsInit(nodeLen int, rv randomvariable.RandomVariable) {
	connectionsEvents := connections.GenerateEvents(nodeLen, rv)
	s.ConnectionsEvents = connectionsEvents
}

func (s *Simulator) NetworkInit(networkPath string, capacitiesPath string) {
	network, err := infrastructure.NetworkGenerate(networkPath, capacitiesPath)

	if err != nil {
		fmt.Printf("Error reading network file: %v\n", err)
	}
	fmt.Println("Network Name:", network.Name)

	s.Network = network
}

func (s *Simulator) BitRateInit(bitRatePath string) {
	bitRate, err := connections.ReadBitRateFile(bitRatePath)

	if err != nil {
		fmt.Printf("Error reading bitrate file: %v\n", err)
	}

	s.BitRateList = bitRate
}

func (s *Simulator) VariablesNumbersInit(numberOfBands int) {
	numberOfNodes := len(s.Network.Nodes)
	numberOfBitrates := len(s.BitRateList.BitRates)

	s.NumberOfNodes = numberOfNodes
	s.NumberOfBitrates = numberOfBitrates

	if numberOfBands > 4 {
		fmt.Println("Warning: Number of bands exceeds 4, setting to 4.")
		numberOfBands = 4
	} else if numberOfBands < 1 {
		fmt.Println("Warning: Number of bands is less than 1, setting to 1.")
		numberOfBands = 1
	}

	s.NumberOfBands = numberOfBands

	s.NumberOfGigabits = 5
}

func New(networkPath string, routesPath string, capacitiesPath string, bitRatePath string, lambda int, mu int, goalConnections float64, allocator allocator.Allocator, numberOfBands int) *Simulator {
	s := &Simulator{}

	s.NetworkInit(networkPath, capacitiesPath)

	s.BitRateInit(bitRatePath)

	s.VariablesNumbersInit(numberOfBands)

	s.RandomVariableInit(lambda, mu, s.NumberOfBitrates, s.NumberOfNodes, s.NumberOfNodes, s.NumberOfBands, s.NumberOfGigabits)

	s.GoalConnections = goalConnections

	s.Time = 0

	s.connectionsEventsInit(s.NumberOfNodes, s.RandomVariable)

	con, err := controller.New(routesPath, s.Network, allocator)
	if err != nil {
		fmt.Printf("Error initializing controller: %v\n", err)
	}
	s.Controller = con

	s.startTime = time.Now()

	return s
}

func (s *Simulator) Start(logOn bool) {
	var countRealease = 0

	fmt.Println("Starting simulation...")

	for i := 0; i <= int(s.GoalConnections); i++ {

		event := s.GetFirstEvent()
		time := event.Time

		rv := s.RandomVariable

		s.printBlockingTable(i, logOn)

		if event.Event == connections.ConnectionEventTypeArrive {

			s.TotalConnections++

			newArriveEvent := connections.ConnectionEvent{
				Id:                   event.Id,
				Source:               event.Source,
				Destination:          event.Destination,
				Bitrate:              event.Bitrate,
				GigabitsSelected:     event.GigabitsSelected,
				Event:                connections.ConnectionEventTypeArrive,
				Time:                 time + rv.GetNetValueExponential("arrive"),
				ConnectionAssignedId: "",
			}
			s.AddNewConnectionEvent(newArriveEvent)

			slot := s.getSlotgigabites(s.BitRateList.BitRates[event.Bitrate], event.GigabitsSelected)

			asigned, con := s.Controller.ConectionAllocation(event.Source, event.Destination, slot, s.Network, s.Controller.Routes, s.NumberOfBands, strconv.Itoa(i))

			if asigned {
				s.Controller.AddConnection(con)
				s.AssignedConnections++

				s.AllocatedConnections = append(s.AllocatedConnections, true)

				newDepartureEvent := connections.ConnectionEvent{
					Id:                   event.Id,
					Source:               event.Source,
					Destination:          event.Destination,
					Bitrate:              event.Bitrate,
					Event:                connections.ConnectionEventTypeRelease,
					GigabitsSelected:     event.GigabitsSelected,
					Time:                 time + rv.GetNetValueExponential("departure"),
					ConnectionAssignedId: strconv.Itoa(i),
				}
				s.AddNewConnectionEvent(newDepartureEvent)
			}

		}

		if event.Event == connections.ConnectionEventTypeRelease {
			countRealease++

			s.Controller.ReleaseConnection(event.ConnectionAssignedId)
		}

	}

	fmt.Println("Simulation completed.", countRealease)

}

func (s *Simulator) Plot(title string, xLabel string, yLabel string) error {
	err := plotter.GenerateScatterPlot(s.Arrives, s.Results, title, xLabel, yLabel)
	return err
}
