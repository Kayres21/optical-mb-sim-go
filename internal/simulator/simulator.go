package simulator

import (
	"container/heap"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"time"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections/controller"
	randomvariable "github.com/Kayres21/optical-mb-sim-go/internal/connections/randomVariable"
	"github.com/Kayres21/optical-mb-sim-go/internal/defragmentator"
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

func (h eventHeap) Len() int           { return len(h) }
func (h eventHeap) Less(i, j int) bool { return h[i].Time < h[j].Time }
func (h eventHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *eventHeap) Push(x any)        { *h = append(*h, x.(connections.ConnectionEvent)) }
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
	DefragMode       string
	DefragDecision   defragmentator.DecisionFunc
	DefragAction     defragmentator.ActionFunc

	events              eventHeap
	assignedConnections int
	totalConnections    int
	startTime           time.Time
	results             []float64
	arrives             []float64

	// Seeds (optional). If zero, defaults will be used.
	SeedArrive      int64
	SeedDeparture   int64
	SeedBitrate     int64
	SeedSource      int64
	SeedDestination int64
	SeedBand        int64
	SeedGigabits    int64

	// For CI calculations (port from C++ implementation)
	zScore     float64
	zScoreEven float64
}

func (s *Simulator) pushEvent(event connections.ConnectionEvent) {
	heap.Push(&s.events, event)
}

func (s *Simulator) popEvent() connections.ConnectionEvent {
	return heap.Pop(&s.events).(connections.ConnectionEvent)
}

func (s *Simulator) getSlotsByGigabits(bitrate connections.BitRate, gigabits int) connections.Slots {
	key := fmt.Sprint(gigabits)
	for _, slot := range bitrate.Slots {
		if slot.Gigabits == key {
			return slot
		}
	}
	return connections.Slots{}
}

func (s *Simulator) shouldRunDefragment(event connections.ConnectionEvent) bool {
	return s.DefragMode != defragmentator.DefragNone && s.DefragDecision(s.Network, s.Controller.Connections, event, s.NumberOfBands)
}

func (s *Simulator) runDefragment() error {
	if s.DefragMode == defragmentator.DefragNone {
		return nil
	}
	return s.DefragAction(s.Network, s.Controller.Connections, s.Controller.Routes, s.Controller.Allocator, s.NumberOfBands)
}

func (s *Simulator) addResult(result float64) {
	s.results = append(s.results, result)
}

func (s *Simulator) addArrive(arrive float64) {
	s.arrives = append(s.arrives, arrive)
}

func (s *Simulator) printBlockingTable(logOn bool) {
	if !logOn {
		return
	}

	if s.totalConnections == 0 {
		fmt.Println("+----------+----------+----------+----------+-------------------+-------------------+-------------------+")
		fmt.Printf("|%10s|%10s|%10s|%10s|%19s|%19s|%19s|\n",
			"progress",
			"arrives",
			"blocking",
			"time(s)",
			"Wald CI",
			"A-C. CI",
			"Wilson CI",
		)
		fmt.Println("+----------+----------+----------+----------+-------------------+-------------------+-------------------+")
		return
	}

	step := s.GoalConnections / 10
	if step == 0 {
		step = 1
	}

	if s.totalConnections%int(step) == 0 {
		blockingProbability := helpers.ComputeBlockingProbabilities(s.assignedConnections, s.totalConnections)
		// Compute metrics
		// Use half-widths (z * sd) like the C++ implementation for parity
		waldHalf := s.waldCIHalfWidth()
		acHalf := s.agrestiCIHalfWidth()
		wilsonHalf := s.wilsonCIHalfWidth()
		progress := (float64(s.totalConnections) / float64(s.GoalConnections)) * 100
		elapsed := time.Since(s.startTime)
		timeFormatted := fmt.Sprintf("%02d:%02d:%02d",
			int(elapsed.Hours()),
			int(elapsed.Minutes())%60,
			int(elapsed.Seconds())%60,
		)

		// Format as single half-width values (scientific, one decimal) to match C++ output
		waldCI := fmt.Sprintf("%9.1e", waldHalf)
		acCI := fmt.Sprintf("%9.1e", acHalf)
		wilsonCI := fmt.Sprintf("%9.1e", wilsonHalf)

		fmt.Printf("|%8.1f %%|%10d|%10.6f|%10s|%19s|%19s|%19s|\n",
			progress,
			s.totalConnections,
			blockingProbability,
			timeFormatted,
			waldCI,
			acCI,
			wilsonCI,
		)
		fmt.Println("+----------+----------+----------+----------+-------------------+-------------------+-------------------+")
		s.addResult(blockingProbability)
		s.addArrive(float64(s.totalConnections))
	}
}

func (s *Simulator) initRandomVariable(lambda, mu float64, seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits int64) {
	var rv randomvariable.RandomVariable
	rv.SetParameters(
		lambda, mu,
		s.NumberOfBitrates,
		s.NumberOfNodes, s.NumberOfNodes,
		s.NumberOfBands,
		s.NumberOfGigabits,
	)

	// If any provided seed is zero, fall back to defaults.
	if seedArrive == 0 {
		seedArrive = defaultSeedArrive
	}
	if seedDeparture == 0 {
		seedDeparture = defaultSeedDeparture
	}
	if seedBitrate == 0 {
		seedBitrate = defaultSeedBitrate
	}
	if seedSource == 0 {
		seedSource = defaultSeedSource
	}
	if seedDestination == 0 {
		seedDestination = defaultSeedDestination
	}
	if seedBand == 0 {
		seedBand = defaultSeedBand
	}
	if seedGigabits == 0 {
		seedGigabits = defaultSeedGigabits
	}

	rv.SetSeeds(
		seedArrive,
		seedDeparture,
		seedBitrate,
		seedSource,
		seedDestination,
		seedBand,
		seedGigabits,
	)
	s.RandomVariable = rv
}

func (s *Simulator) initConnectionEvents() {
	raw := []connections.ConnectionEvent{
		s.createRandomArrival(0, "0"),
	}
	s.events = eventHeap(raw)
	heap.Init(&s.events)
}

func (s *Simulator) initNetwork(network infrastructure.Network) {
	fmt.Println("Network Name:", network.Name)
	s.Network = network
}

func (s *Simulator) initBitRate(bitRate connections.BitRateList) {
	s.BitRateList = bitRate
}

func (s *Simulator) initVariableNumbers(numberOfBands int) {
	s.NumberOfNodes = len(s.Network.Nodes)
	if len(s.BitRateList.BitRates) > 0 {
		s.NumberOfBitrates = len(s.BitRateList.BitRates) * len(s.BitRateList.BitRates[0].Slots)
	} else {
		s.NumberOfBitrates = 0
	}
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
	network infrastructure.Network,
	bitRate connections.BitRateList,
	routes connections.Routes,
	lambda, mu float64,
	goalConnections float64,
	alloc allocator.Allocator,
	numberOfBands int,
	defragMode string,
	defragDecision defragmentator.DecisionFunc,
	defragAction defragmentator.ActionFunc,
) (*Simulator, error) {
	// Delegate to NewWithSeeds using default seeds (preserves API compatibility)
	return NewWithSeeds(
		network,
		bitRate,
		routes,
		lambda,
		mu,
		goalConnections,
		alloc,
		numberOfBands,
		defragMode,
		defragDecision,
		defragAction,
		defaultSeedArrive,
		defaultSeedDeparture,
		defaultSeedBitrate,
		defaultSeedSource,
		defaultSeedDestination,
		defaultSeedBand,
		defaultSeedGigabits,
	)
}

// NewWithSeeds constructs and initialises a Simulator and allows providing explicit RNG seeds.
func NewWithSeeds(
	network infrastructure.Network,
	bitRate connections.BitRateList,
	routes connections.Routes,
	lambda, mu float64,
	goalConnections float64,
	alloc allocator.Allocator,
	numberOfBands int,
	defragMode string,
	defragDecision defragmentator.DecisionFunc,
	defragAction defragmentator.ActionFunc,
	seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits int64,
) (*Simulator, error) {
	s := &Simulator{}

	s.initNetwork(network)
	s.initBitRate(bitRate)

	s.initVariableNumbers(numberOfBands)
	// store seeds in struct for later introspection
	s.SeedArrive = seedArrive
	s.SeedDeparture = seedDeparture
	s.SeedBitrate = seedBitrate
	s.SeedSource = seedSource
	s.SeedDestination = seedDestination
	s.SeedBand = seedBand
	s.SeedGigabits = seedGigabits

	s.initRandomVariable(lambda, mu, seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits)

	s.GoalConnections = goalConnections
	s.Time = 0
	s.DefragMode = defragMode
	if defragDecision == nil {
		defragDecision = defragmentator.DefaultDecision
	}
	if defragAction == nil {
		defragAction = defragmentator.DefaultAction
	}
	s.DefragDecision = defragDecision
	s.DefragAction = defragAction

	s.initConnectionEvents()

	con := controller.Controller{
		Routes:      routes,
		Connections: make(map[string]connections.Connection),
		Network:     s.Network,
		Allocator:   alloc,
	}
	s.Controller = con

	s.startTime = time.Now()
	// Initialize z-score values for CI calculations (mirroring C++ behavior)
	s.initZScore()
	s.initZScoreEven()

	return s, nil
}

func (s *Simulator) getAllocatedProbability() float64 {
	if s.totalConnections == 0 {
		return 0.0
	}
	return float64(s.assignedConnections) / float64(s.totalConnections)
}

func (s *Simulator) waldCIHalfWidth() float64 {
	np := s.getAllocatedProbability()
	p := 1.0 - np
	n := float64(s.totalConnections)
	if n <= 0.0 {
		return 0.0
	}
	sd := math.Sqrt((np * p) / n)
	return s.zScore * sd
}

func (s *Simulator) agrestiCIHalfWidth() float64 {
	np := s.getAllocatedProbability()
	n := float64(s.totalConnections)
	if n <= 0.0 {
		return 0.0
	}

	adjusted := np * ((n * (float64(s.assignedConnections) + (s.zScoreEven / 2.0))) / (float64(s.assignedConnections) * (n + s.zScoreEven)))
	p := 1.0 - adjusted
	sd := math.Sqrt((adjusted * p) / (n + s.zScoreEven))
	return s.zScore * sd
}

func (s *Simulator) wilsonCIHalfWidth() float64 {
	np := s.getAllocatedProbability()
	p := 1.0 - np
	n := float64(s.totalConnections)
	if n <= 0.0 {
		return 0.0
	}

	denom := (1 + (math.Pow(s.zScore, 2) / n))
	sd := math.Sqrt(((np * p) / n) + ((math.Pow(s.zScore, 2)) / (4 * math.Pow(n, 2))))
	return (s.zScore * sd) / denom
}

func (s *Simulator) initZScore() {
	actual := 0.0
	step := 1.0
	covered := 0.0
	objective := 0.95
	if s.GoalConnections > 0 {
	}
	epsilon := 1e-6

	for math.Abs(objective-covered) > epsilon {
		if objective > covered {
			actual += step
			covered = ((1 + math.Erf(actual/math.Sqrt(2))) - (1 + math.Erf(-actual/math.Sqrt(2)))) / 2
			if covered > objective {
				step /= 2
			}
		} else {
			actual -= step
			covered = ((1 + math.Erf(actual/math.Sqrt(2))) - (1 + math.Erf(-actual/math.Sqrt(2)))) / 2
			if covered < objective {
				step /= 2
			}
		}
	}
	s.zScore = actual
}

func (s *Simulator) initZScoreEven() {
	zEven := math.Pow(s.zScore, 2)
	zEven = math.Floor(zEven*1000.0) / 1000.0
	zEven = math.Ceil(zEven/2.0) * 2.0
	s.zScoreEven = zEven
}

func (s *Simulator) createRandomArrival(currentTime float64, id string) connections.ConnectionEvent {
	rv := s.RandomVariable
	source := rv.GetNetValueUniform(randomvariable.KeySource)
	destination := rv.GetNetValueUniform(randomvariable.KeyDestination)
	for source == destination {
		destination = rv.GetNetValueUniform(randomvariable.KeyDestination)
	}

	unifiedIndex := rv.GetNetValueUniform(randomvariable.KeyBitrate)
	var modulationIndex, gigabits int

	if len(s.BitRateList.BitRates) > 0 && len(s.BitRateList.BitRates[0].Slots) > 0 {
		slotsCount := len(s.BitRateList.BitRates[0].Slots)
		modulationIndex = unifiedIndex / slotsCount
		slotIndex := unifiedIndex % slotsCount

		gigaStr := s.BitRateList.BitRates[modulationIndex].Slots[slotIndex].Gigabits
		gigabits, _ = strconv.Atoi(gigaStr)
	} else {
		modulationIndex = 0
		gigabits = 10
	}

	return connections.ConnectionEvent{
		Id:                   id,
		Source:               source,
		Destination:          destination,
		Bitrate:              modulationIndex,
		GigabitsSelected:     gigabits,
		Event:                connections.ConnectionEventTypeArrive,
		Time:                 currentTime + rv.GetNetValueExponential(randomvariable.KeyArrive),
		ConnectionAssignedId: "",
	}
}

func (s *Simulator) Start(logOn bool) {
	countRelease := 0

	fmt.Println("Starting simulation...")
	s.printBlockingTable(logOn) // Print header

	for s.totalConnections < int(s.GoalConnections) {
		event := s.popEvent()
		s.Time = event.Time
		rv := s.RandomVariable

		if event.Event == connections.ConnectionEventTypeArrive {
			s.totalConnections++
			s.printBlockingTable(logOn)

			if s.DefragMode == defragmentator.DefragBeforeArrival && s.shouldRunDefragment(event) {
				if err := s.runDefragment(); err != nil {
					slog.Warn("defragmentation failed before arrival", "err", err)
				}
			}

			// Schedule the next arrival.
			nextArrive := s.createRandomArrival(event.Time, event.Id)
			s.pushEvent(nextArrive)

			selectedBitrate := s.BitRateList.BitRates[event.Bitrate]
			slotsConfig := s.getSlotsByGigabits(selectedBitrate, event.GigabitsSelected)

			getSlot := func(bandIndex int) int {
				band := s.Network.Links[0].Capacities.Bands[bandIndex]
				if sVal, ok := slotsConfig.SlotsPerBand[band.Name]; ok {
					return sVal
				}
				return slotsConfig.Slots
			}

			assigned, con := s.Controller.ConnectionAllocation(event.Source, event.Destination, getSlot, s.NumberOfBands, strconv.Itoa(s.totalConnections))

			if !assigned && s.DefragMode == defragmentator.DefragAfterBlock && s.shouldRunDefragment(event) {
				if err := s.runDefragment(); err != nil {
					slog.Warn("defragmentation failed after block", "err", err)
				} else {
					assigned, con = s.Controller.ConnectionAllocation(event.Source, event.Destination, getSlot, s.NumberOfBands, strconv.Itoa(s.totalConnections))
				}
			}

			if assigned {
				s.Controller.AddConnection(con)
				s.assignedConnections++

				if s.DefragMode == defragmentator.DefragAfterAssign && s.shouldRunDefragment(event) {
					if err := s.runDefragment(); err != nil {
						slog.Warn("defragmentation failed after assign", "err", err)
					}
				}

				departure := connections.ConnectionEvent{
					Id:                     event.Id,
					Source:                 event.Source,
					Destination:            event.Destination,
					Bitrate:                event.Bitrate,
					GigabitsSelected:       event.GigabitsSelected,
					Event:                  connections.ConnectionEventTypeRelease,
					Time:                   event.Time + rv.GetNetValueExponential(randomvariable.KeyDeparture),
					ConnectionAssignedId:   strconv.Itoa(s.totalConnections),
					ConnectionInitialSlot:  con.InitialSlot,
					ConnectionSlots:        con.Slots,
					ConnectionBandSelected: con.BandSelected,
				}
				s.pushEvent(departure)
			}
		}

		if event.Event == connections.ConnectionEventTypeRelease {
			countRelease++
			connection, ok := s.Controller.GetConnectionByAllocation(event.Source, event.Destination, event.ConnectionInitialSlot, event.ConnectionSlots, event.ConnectionBandSelected)
			if !ok {
				slog.Warn("failed to release connection: not found by allocation details", "source", event.Source, "destination", event.Destination, "initialSlot", event.ConnectionInitialSlot, "slots", event.ConnectionSlots, "band", event.ConnectionBandSelected)
				continue
			}

			if err := s.Controller.ReleaseConnection(connection, event.Time); err != nil {
				slog.Warn("failed to release connection", "id", event.ConnectionAssignedId, "err", err)
			}
		}
	}

	fmt.Printf("Simulation completed. Releases processed: %d, Total simulated time: %.2f\n", countRelease, s.Time)
}

func (s *Simulator) Plot(title, xLabel, yLabel string) error {
	return plotter.GenerateScatterPlot(s.arrives, s.results, title, xLabel, yLabel)
}

func (s *Simulator) SetSeeds(seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits int64) {
	s.SeedArrive = seedArrive
	s.SeedDeparture = seedDeparture
	s.SeedBitrate = seedBitrate
	s.SeedSource = seedSource
	s.SeedDestination = seedDestination
	s.SeedBand = seedBand
	s.SeedGigabits = seedGigabits

	if (s.RandomVariable != randomvariable.RandomVariable{}) {
		s.RandomVariable.SetSeeds(seedArrive, seedDeparture, seedBitrate, seedSource, seedDestination, seedBand, seedGigabits)
	}
}

func (s *Simulator) SetSeedArrive(seed int64) {
	s.SetSeeds(seed, s.SeedDeparture, s.SeedBitrate, s.SeedSource, s.SeedDestination, s.SeedBand, s.SeedGigabits)
}

func (s *Simulator) SetSeedDeparture(seed int64) {
	s.SetSeeds(s.SeedArrive, seed, s.SeedBitrate, s.SeedSource, s.SeedDestination, s.SeedBand, s.SeedGigabits)
}

func (s *Simulator) SetSeedBitrate(seed int64) {
	s.SetSeeds(s.SeedArrive, s.SeedDeparture, seed, s.SeedSource, s.SeedDestination, s.SeedBand, s.SeedGigabits)
}

func (s *Simulator) SetSeedSource(seed int64) {
	s.SetSeeds(s.SeedArrive, s.SeedDeparture, s.SeedBitrate, seed, s.SeedDestination, s.SeedBand, s.SeedGigabits)
}

func (s *Simulator) SetSeedDestination(seed int64) {
	s.SetSeeds(s.SeedArrive, s.SeedDeparture, s.SeedBitrate, s.SeedSource, seed, s.SeedBand, s.SeedGigabits)
}

func (s *Simulator) SetSeedBand(seed int64) {
	s.SetSeeds(s.SeedArrive, s.SeedDeparture, s.SeedBitrate, s.SeedSource, s.SeedDestination, seed, s.SeedGigabits)
}

func (s *Simulator) SetSeedGigabits(seed int64) {
	s.SetSeeds(s.SeedArrive, s.SeedDeparture, s.SeedBitrate, s.SeedSource, s.SeedDestination, s.SeedBand, seed)
}

// SetAllocator allows replacing the allocator after constructing the Simulator.
func (s *Simulator) SetAllocator(alloc allocator.Allocator) {
	s.Controller.Allocator = alloc
}

func (s *Simulator) SetDefragDecision(dec defragmentator.DecisionFunc) {
	if dec == nil {
		dec = defragmentator.DefaultDecision
	}
	s.DefragDecision = dec
}

func (s *Simulator) SetDefragAction(act defragmentator.ActionFunc) {
	if act == nil {
		act = defragmentator.DefaultAction
	}
	s.DefragAction = act
}

func (s *Simulator) SetDefragMode(mode string) {
	s.DefragMode = mode
}

// Unassign callback setters to mirror C++ API
func (s *Simulator) SetUnassignCallback(cb func(connections.Connection, float64, infrastructure.Network)) {
	s.Controller.SetUnassignCallback(cb)
}

func (s *Simulator) SetUnassignMB() {
	s.Controller.SetUnassignMB()
}
