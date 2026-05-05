package simulator

import (
	"container/heap"
	"fmt"
	"log/slog"
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

// Default seeds for reproducible simulation runs.
const (
	defaultSeedArrive      int64 = 1
	defaultSeedDeparture   int64 = 12
	defaultSeedBitrate     int64 = 123
	defaultSeedSource      int64 = 1234
	defaultSeedDestination int64 = 12345
	defaultSeedBand        int64 = 123456
	defaultSeedGigabits    int64 = 1234567
)

// Band count limits.
const (
	minBands = 1
	maxBands = 4
)

// ─── Event heap (min-heap by Time) ───────────────────────────────────────────

type eventHeap []connections.ConnectionEvent

func (h eventHeap) Len() int            { return len(h) }
func (h eventHeap) Less(i, j int) bool  { return h[i].Time < h[j].Time }
func (h eventHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *eventHeap) Push(x any)         { *h = append(*h, x.(connections.ConnectionEvent)) }
func (h *eventHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// ─── Simulator ────────────────────────────────────────────────────────────────

type Simulator struct {
	RandomVariable   randomvariable.RandomVariable
	Network          infrastructure.Network
	BitRateList      connections.BitRateList
	Controller       controller.Controller
	GoalConnections  float64
	Time             float64
	NumberOfBands    int
	NumberOfBitrates int
	NumberOfNodes    int
	NumberOfGigabits int

	events              eventHeap
	assignedConnections int
	totalConnections    int
	startTime           time.Time
	results             []float64
	arrives             []float64
}

func (s *Simulator) pushEvent(event connections.ConnectionEvent) {
	heap.Push(&s.events, event)
}

func (s *Simulator) popEvent() connections.ConnectionEvent {
	return heap.Pop(&s.events).(connections.ConnectionEvent)
}

func (s *Simulator) getSlotsByGigabits(bitrate connections.BitRate, gigabits int) int {
	key := fmt.Sprint(gigabits)
	for _, slot := range bitrate.Slots {
		if slot.Gigabits == key {
			return slot.Slots
		}
	}
	return 0
}

func (s *Simulator) addResult(result float64) {
	s.results = append(s.results, result)
}

func (s *Simulator) addArrive(arrive float64) {
	s.arrives = append(s.arrives, arrive)
}

func (s *Simulator) printBlockingTable(i int, logOn bool) {
	if !logOn {
		return
	}

	if i == 0 {
		fmt.Println("+----------+----------+----------+----------+")
		fmt.Println("| progress |  arrives | blocking |  time(s) |")
		fmt.Println("+----------+----------+----------+----------+")
		return
	}

	step := s.GoalConnections / 10
	if step == 0 {
		step = 1
	}

	if i%int(step) == 0 {
		blockingProbability := helpers.ComputeBlockingProbabilities(s.assignedConnections, s.totalConnections)
		progress := (float64(i) / float64(s.GoalConnections)) * 100
		elapsed := time.Since(s.startTime)
		timeFormatted := fmt.Sprintf("%02d:%02d:%02d",
			int(elapsed.Hours()),
			int(elapsed.Minutes())%60,
			int(elapsed.Seconds())%60,
		)

		fmt.Printf("|%8.1f %%|%10d|%10.6f|%10s|\n", progress, i, blockingProbability, timeFormatted)
		fmt.Println("+----------+----------+----------+----------+")
		s.addResult(blockingProbability)
		s.addArrive(float64(i))
	}
}

func (s *Simulator) initRandomVariable(lambda, mu int) {
	var rv randomvariable.RandomVariable
	rv.SetParameters(
		lambda, mu,
		s.NumberOfBitrates,
		s.NumberOfNodes, s.NumberOfNodes,
		s.NumberOfBands,
		s.NumberOfGigabits,
	)
	rv.SetSeeds(
		defaultSeedArrive,
		defaultSeedDeparture,
		defaultSeedBitrate,
		defaultSeedSource,
		defaultSeedDestination,
		defaultSeedBand,
		defaultSeedGigabits,
	)
	s.RandomVariable = rv
}

func (s *Simulator) initConnectionEvents() {
	raw := connections.GenerateEvents(s.NumberOfNodes, s.RandomVariable)
	s.events = eventHeap(raw)
	heap.Init(&s.events)
}

func (s *Simulator) initNetwork(networkPath, capacitiesPath string) error {
	network, err := infrastructure.NetworkGenerate(networkPath, capacitiesPath)
	if err != nil {
		return fmt.Errorf("initializing network: %w", err)
	}
	fmt.Println("Network Name:", network.Name)
	s.Network = network
	return nil
}

func (s *Simulator) initBitRate(bitRatePath string) error {
	bitRate, err := connections.ReadBitRateFile(bitRatePath)
	if err != nil {
		return fmt.Errorf("initializing bitrates: %w", err)
	}
	s.BitRateList = bitRate
	return nil
}

func (s *Simulator) initVariableNumbers(numberOfBands int) {
	s.NumberOfNodes = len(s.Network.Nodes)
	s.NumberOfBitrates = len(s.BitRateList.BitRates)
	s.NumberOfGigabits = len(randomvariable.DefaultGigabitOptions)

	switch {
	case numberOfBands > maxBands:
		fmt.Printf("Warning: number of bands exceeds %d, setting to %d.\n", maxBands, maxBands)
		numberOfBands = maxBands
	case numberOfBands < minBands:
		fmt.Printf("Warning: number of bands is less than %d, setting to %d.\n", minBands, minBands)
		numberOfBands = minBands
	}

	s.NumberOfBands = numberOfBands
}

// New constructs and initialises a Simulator. Returns an error if any
// resource file cannot be loaded or the controller cannot be created.
func New(
	networkPath, routesPath, capacitiesPath, bitRatePath string,
	lambda, mu int,
	goalConnections float64,
	alloc allocator.Allocator,
	numberOfBands int,
) (*Simulator, error) {
	s := &Simulator{}

	if err := s.initNetwork(networkPath, capacitiesPath); err != nil {
		return nil, err
	}

	if err := s.initBitRate(bitRatePath); err != nil {
		return nil, err
	}

	s.initVariableNumbers(numberOfBands)
	s.initRandomVariable(lambda, mu)

	s.GoalConnections = goalConnections
	s.Time = 0

	s.initConnectionEvents()

	con, err := controller.New(routesPath, s.Network, alloc)
	if err != nil {
		return nil, fmt.Errorf("initializing controller: %w", err)
	}
	s.Controller = con

	s.startTime = time.Now()

	return s, nil
}

func (s *Simulator) Start(logOn bool) {
	countRelease := 0

	fmt.Println("Starting simulation...")

	for i := 0; i <= int(s.GoalConnections); i++ {
		event := s.popEvent()
		rv := s.RandomVariable

		s.printBlockingTable(i, logOn)

		if event.Event == connections.ConnectionEventTypeArrive {
			s.totalConnections++

			// Schedule the next arrival for this traffic stream.
			nextArrive := connections.ConnectionEvent{
				Id:                   event.Id,
				Source:               event.Source,
				Destination:          event.Destination,
				Bitrate:              event.Bitrate,
				GigabitsSelected:     event.GigabitsSelected,
				Event:                connections.ConnectionEventTypeArrive,
				Time:                 event.Time + rv.GetNetValueExponential(randomvariable.KeyArrive),
				ConnectionAssignedId: "",
			}
			s.pushEvent(nextArrive)

			slot := s.getSlotsByGigabits(s.BitRateList.BitRates[event.Bitrate], event.GigabitsSelected)
			assigned, con := s.Controller.ConnectionAllocation(event.Source, event.Destination, slot, s.NumberOfBands, strconv.Itoa(i))

			if assigned {
				s.Controller.AddConnection(con)
				s.assignedConnections++

				departure := connections.ConnectionEvent{
					Id:                   event.Id,
					Source:               event.Source,
					Destination:          event.Destination,
					Bitrate:              event.Bitrate,
					GigabitsSelected:     event.GigabitsSelected,
					Event:                connections.ConnectionEventTypeRelease,
					Time:                 event.Time + rv.GetNetValueExponential(randomvariable.KeyDeparture),
					ConnectionAssignedId: strconv.Itoa(i),
				}
				s.pushEvent(departure)
			}
		}

		if event.Event == connections.ConnectionEventTypeRelease {
			countRelease++
			if err := s.Controller.ReleaseConnection(event.ConnectionAssignedId); err != nil {
				slog.Warn("failed to release connection", "id", event.ConnectionAssignedId, "err", err)
			}
		}
	}

	fmt.Printf("Simulation completed. Releases processed: %d\n", countRelease)
}

func (s *Simulator) Plot(title, xLabel, yLabel string) error {
	return plotter.GenerateScatterPlot(s.arrives, s.results, title, xLabel, yLabel)
}
