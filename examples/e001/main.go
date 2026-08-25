package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/Kayres21/optical-mb-sim-go/internal/allocator"
	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/defragmentator"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
	"github.com/Kayres21/optical-mb-sim-go/internal/loader"
	"github.com/Kayres21/optical-mb-sim-go/internal/simulator"
)

type eventType string

const (
	eventArrive eventType = "ARRIVE"
	eventDepart eventType = "DEPARTURE"
)

type eventHistory struct {
	ID      int
	Time    float64
	Type    eventType
	Source  int
	Dest    int
	BitRate int
}

func projectRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}

func numberOfRoutes(routes connections.Routes, src, dst int) int {
	return len(routes.GetPaths(src, dst))
}

func numberOfLinks(routes connections.Routes, src, dst, routeIdx int) int {
	paths := routes.GetPaths(src, dst)
	if routeIdx < 0 || routeIdx >= len(paths) {
		return 0
	}
	return len(paths[routeIdx])
}

func linkInRoute(route []int, linkIdx int) int {
	if linkIdx < 0 || linkIdx >= len(route) {
		return -1
	}
	return route[linkIdx]
}

func reqRoute(routes connections.Routes, src, dst, routeIdx int) []int {
	paths := routes.GetPaths(src, dst)
	if routeIdx < 0 || routeIdx >= len(paths) {
		return nil
	}
	return paths[routeIdx]
}

func reqRouteLength(network infrastructure.Network, routes connections.Routes, src, dst, routeIdx int) int {
	route := reqRoute(routes, src, dst, routeIdx)
	length := 0
	for _, linkID := range route {
		if link := network.GetLinkByID(linkID); link != nil {
			length += link.Length
		}
	}
	return length
}

func reqDistanceAdaptiveModulation(bitRate connections.BitRate, routeLength int) int {
	return bitRate.GetDistanceAdaptive(float64(routeLength))
}

func linkInRouteID(network infrastructure.Network, routes connections.Routes, src, dst, routeIdx, linkIdx int) int {
	route := reqRoute(routes, src, dst, routeIdx)
	return linkInRoute(route, linkIdx)
}

func firstFit(src, dst int, bitRate connections.BitRate, network infrastructure.Network, routes connections.Routes) (bool, error) {
	for routeIdx := 0; routeIdx < numberOfRoutes(routes, src, dst); routeIdx++ {
		firstLinkID := linkInRouteID(network, routes, src, dst, routeIdx, 0)
		if firstLinkID == -1 {
			continue
		}
		firstLink := network.GetLinkByID(firstLinkID)
		if firstLink == nil {
			continue
		}

		totalSlots := make([]bool, len(firstLink.GetSlotsByBand(0)))
		for linkIdx := 0; linkIdx < numberOfLinks(routes, src, dst, routeIdx); linkIdx++ {
			linkID := linkInRouteID(network, routes, src, dst, routeIdx, linkIdx)
			link := network.GetLinkByID(linkID)
			if link == nil {
				continue
			}
			for slotIdx, occupied := range link.GetSlotsByBand(0) {
				totalSlots[slotIdx] = totalSlots[slotIdx] || occupied
			}
		}

		routeLength := reqRouteLength(network, routes, src, dst, routeIdx)
		modulation := reqDistanceAdaptiveModulation(bitRate, routeLength)
		if modulation == -1 {
			continue
		}

		requiredSlots := bitRate.Slots[modulation]
		currentSlots := 0
		currentSlotIndex := 0

		for slotIdx := range totalSlots {
			if !totalSlots[slotIdx] {
				currentSlots++
			} else {
				currentSlots = 0
				currentSlotIndex = slotIdx + 1
			}

			if currentSlots == requiredSlots {
				for linkIdx := 0; linkIdx < numberOfLinks(routes, src, dst, routeIdx); linkIdx++ {
					linkID := linkInRouteID(network, routes, src, dst, routeIdx, linkIdx)
					link := network.GetLinkByID(linkID)
					if link == nil {
						return false, fmt.Errorf("link %d not found in route %d", linkID, routeIdx)
					}
					if err := link.AssignConnection(currentSlotIndex, requiredSlots, 0); err != nil {
						return false, err
					}
				}
				return true, nil
			}
		}
	}

	return false, nil
}

func importFromCSV(filename string) ([]connections.ConnectionEvent, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open event csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read event csv: %w", err)
	}
	if len(rows) < 2 {
		return nil, nil
	}

	history := make([]connections.ConnectionEvent, 0, len(rows)-1)
	for _, row := range rows[1:] {
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}
		if len(row) < 6 {
			return nil, fmt.Errorf("invalid CSV row: %#v", row)
		}

		timeValue, err := strconv.ParseFloat(strings.TrimSpace(row[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("parse time: %w", err)
		}
		eventTypeValue := strings.TrimSpace(row[2])
		src, err := strconv.Atoi(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, fmt.Errorf("parse source: %w", err)
		}
		dst, err := strconv.Atoi(strings.TrimSpace(row[4]))
		if err != nil {
			return nil, fmt.Errorf("parse destination: %w", err)
		}
		bitrate, err := strconv.Atoi(strings.TrimSpace(row[5]))
		if err != nil {
			return nil, fmt.Errorf("parse bitrate: %w", err)
		}

		event := connections.ConnectionEvent{
			Id:          strings.TrimSpace(row[0]),
			Source:      src,
			Destination: dst,
			Bitrate:     bitrate,
			Time:        timeValue,
		}
		switch strings.ToUpper(eventTypeValue) {
		case "ARRIVE":
			event.Event = connections.ConnectionEventTypeArrive
		case "DEPARTURE":
			event.Event = connections.ConnectionEventTypeRelease
		default:
			return nil, fmt.Errorf("unknown event type %q", eventTypeValue)
		}

		history = append(history, event)
	}

	return history, nil
}

type csvRunStats struct {
	TotalArrivals int
	Accepted      int
	Blocked       int
}

func runSimulationFromCSV(sim *simulator.Simulator, events []connections.ConnectionEvent) csvRunStats {
	stats := csvRunStats{}

	for _, event := range events {
		switch event.Event {
		case connections.ConnectionEventTypeArrive:
			stats.TotalArrivals++
			selectedBitrate := sim.BitRateList.BitRates[event.Bitrate]
			assigned := sim.Controller.ConnectionAllocation(
				event.Source,
				event.Destination,
				selectedBitrate,
				1,
				event.Id,
			)
			if assigned {
				stats.Accepted++
				fmt.Printf("ARRIVE id=%s src=%d dst=%d bitrate=%d time=%.6f -> allocated\n",
					event.Id, event.Source, event.Destination, event.Bitrate, event.Time)
			} else {
				stats.Blocked++
				fmt.Printf("ARRIVE id=%s src=%d dst=%d bitrate=%d time=%.6f -> blocked\n",
					event.Id, event.Source, event.Destination, event.Bitrate, event.Time)
			}
		case connections.ConnectionEventTypeRelease:
			if connection, ok := sim.Controller.GetConnectionById(event.Id); ok {
				if err := sim.Controller.ReleaseConnection(connection, event.Time); err != nil {
					fmt.Printf("DEPARTURE id=%s time=%.6f -> release error: %v\n", event.Id, event.Time, err)
				} else {
					fmt.Printf("DEPARTURE id=%s time=%.6f -> released\n", event.Id, event.Time)
				}
			} else {
				fmt.Printf("DEPARTURE id=%s time=%.6f -> not found in active connections\n", event.Id, event.Time)
			}
		}
	}

	return stats
}

func exportToCSV(filename string, history []eventHistory) error {
	sortedHistory := append([]eventHistory(nil), history...)
	sort.Slice(sortedHistory, func(i, j int) bool { return sortedHistory[i].Time < sortedHistory[j].Time })

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create output csv: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"ID", "Tiempo", "Evento", "Source", "Destination", "BitRate"}); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	for _, event := range sortedHistory {
		row := []string{
			strconv.Itoa(event.ID),
			strconv.FormatFloat(event.Time, 'f', -1, 64),
			string(event.Type),
			strconv.Itoa(event.Source),
			strconv.Itoa(event.Dest),
			strconv.Itoa(event.BitRate),
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush csv: %w", err)
	}
	return nil
}

func main() {
	root := projectRoot()

	networkPath := filepath.Join(root, "examples", "e001", "BDM_UKNet.json")
	routesPath := filepath.Join(root, "examples", "e001", "BDM_UKNet_rutas.json")
	bitratePath := filepath.Join(root, "examples", "e001", "bitrate_iroBand_C.json")
	csvPath := filepath.Join(root, "examples", "e001", "UKNetOrderv2.csv")

	resourceLoader := &loader.LegacyLoader{}

	network, err := resourceLoader.LoadNetwork(networkPath, networkPath)
	if err != nil {
		log.Fatalf("failed to load network: %v", err)
	}

	bitRate, err := resourceLoader.LoadBitRate(bitratePath, 1)
	if err != nil {
		log.Fatalf("failed to load bitrate: %v", err)
	}

	routes, err := resourceLoader.LoadRoutes(routesPath)
	if err != nil {
		log.Fatalf("failed to load routes: %v", err)
	}

	events, err := importFromCSV(csvPath)
	if err != nil {
		log.Fatalf("failed to load event CSV: %v", err)
	}

	fmt.Printf("Loaded %d events from %s (limited to first 100)\n", len(events), filepath.Base(csvPath))

	for i := 0; i < 5 && i < len(events); i++ {
		fmt.Printf("event[%d] = %s @ %.6f src=%d dst=%d bitrate=%d\n", i, events[i].Event, events[i].Time, events[i].Source, events[i].Destination, events[i].Bitrate)
	}

	sim, err := simulator.New(
		network,
		bitRate,
		routes,
		0,
		1,
		float64(10),
		allocator.FirstFit,
		1,
		defragmentator.DefragNone,
		defragmentator.DefaultDecision,
		defragmentator.DefaultAction,
	)
	if err != nil {
		log.Fatalf("failed to create simulator: %v", err)
	}

	stats := runSimulationFromCSV(sim, events)
	if stats.TotalArrivals > 0 {
		blockingProbability := float64(stats.Blocked) / float64(stats.TotalArrivals)
		fmt.Printf("\nFinal summary: arrivals=%d accepted=%d blocked=%d blocking_probability=%.6f\n",
			stats.TotalArrivals, stats.Accepted, stats.Blocked, blockingProbability)
	}
	fmt.Println("Example e001 finished using UKNetOrderv2.csv without random arrivals")
}
