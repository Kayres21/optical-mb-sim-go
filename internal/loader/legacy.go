package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

// LegacyLoader implements ResourceLoader for legacy file formats.
type LegacyLoader struct{}

// LoadNetwork parses legacy network JSON files.
func (l *LegacyLoader) LoadNetwork(networkPath, capacitiesPath string) (infrastructure.Network, error) {
	data, err := os.ReadFile(networkPath)
	if err != nil {
		return infrastructure.Network{}, fmt.Errorf("reading legacy network file: %w", err)
	}

	var raw struct {
		Name  string `json:"Name"`
		Alias string `json:"alias"`
		Nodes []struct {
			ID int `json:"id"`
		} `json:"nodes"`
		Links []struct {
			ID     int             `json:"id"`
			Src    int             `json:"src"`
			Dst    int             `json:"dst"`
			Length int             `json:"length"`
			Slots  json.RawMessage `json:"slots"`
		} `json:"links"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return infrastructure.Network{}, fmt.Errorf("parsing legacy network JSON: %w", err)
	}

	network := infrastructure.Network{
		Name:  raw.Name,
		Alias: raw.Alias,
	}

	for _, n := range raw.Nodes {
		network.Nodes = append(network.Nodes, infrastructure.Node{ID: n.ID})
	}

	for _, rl := range raw.Links {
		link := infrastructure.Link{
			ID:          rl.ID,
			Source:      rl.Src,
			Destination: rl.Dst,
			Length:      rl.Length,
		}

		// Try to parse slots as int (single band)
		var slotsInt int
		if err := json.Unmarshal(rl.Slots, &slotsInt); err == nil {
			link.Capacities = infrastructure.Capacity{
				Bands: []infrastructure.Band{
					{
						ID:       "0",
						Name:     "C",
						SlotsLen: slotsInt,
						Slots:    make([]bool, slotsInt),
					},
				},
			}
		} else {
			// Try to parse as map (multi band)
			var slotsMap map[string]int
			if err := json.Unmarshal(rl.Slots, &slotsMap); err == nil {
				// We need to order the bands. Common order is C, L, S, E or similar.
				// For now, let's just sort them alphabetically to be consistent.
				var keys []string
				for k := range slotsMap {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				var bands []infrastructure.Band
				for i, k := range keys {
					sLen := slotsMap[k]
					bands = append(bands, infrastructure.Band{
						ID:       strconv.Itoa(i),
						Name:     k,
						SlotsLen: sLen,
						Slots:    make([]bool, sLen),
					})
				}
				link.Capacities = infrastructure.Capacity{Bands: bands}
			} else {
				return infrastructure.Network{}, fmt.Errorf("unknown slots format in link %d", rl.ID)
			}
		}

		network.Links = append(network.Links, link)
	}

	return network, nil
}

// LoadBitRate parses legacy bitrate JSON files.
func (l *LegacyLoader) LoadBitRate(bitRatePath string) (connections.BitRateList, error) {
	data, err := os.ReadFile(bitRatePath)
	if err != nil {
		return connections.BitRateList{}, fmt.Errorf("reading legacy bitrate file: %w", err)
	}

	// Legacy format is map[gigabits][]map[modulation]config
	var raw map[string][]map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return connections.BitRateList{}, fmt.Errorf("parsing legacy bitrate JSON: %w", err)
	}

	modMap := make(map[string]*connections.BitRate)

	for gigabits, entries := range raw {
		for _, entry := range entries {
			for modulation, configRaw := range entry {
				br, ok := modMap[modulation]
				if !ok {
					br = &connections.BitRate{Modulation: modulation}
					modMap[modulation] = br
				}

				// Check if configRaw is a single config or a band map
				var singleConfig struct {
					Slots int `json:"slots"`
					Reach int `json:"reach"`
				}

				if err := json.Unmarshal(configRaw, &singleConfig); err == nil {
					// Single band legacy format
					br.Slots = append(br.Slots, connections.Slots{
						Gigabits: gigabits,
						Slots:    singleConfig.Slots,
					})
					addReach(br, 1, "C", singleConfig.Reach)
				} else {
					// Multi band legacy format: [ { "C": {...} }, { "L": {...} } ]
					var bandConfigs []map[string]struct {
						Slots int `json:"slots"`
						Reach int `json:"reach"`
					}
					if err := json.Unmarshal(configRaw, &bandConfigs); err == nil {
						mainSlots := 0
						slotsPerBand := make(map[string]int)
						numBands := len(bandConfigs)
						
						for _, bc := range bandConfigs {
							for bandName, v := range bc {
								if mainSlots == 0 || bandName == "C" {
									mainSlots = v.Slots
								}
								slotsPerBand[bandName] = v.Slots
								addReach(br, numBands, bandName, v.Reach)
							}
						}

						br.Slots = append(br.Slots, connections.Slots{
							Gigabits:     gigabits,
							Slots:        mainSlots,
							SlotsPerBand: slotsPerBand,
						})
					}
				}
			}
		}
	}

	var res connections.BitRateList
	for _, br := range modMap {
		res.BitRates = append(res.BitRates, *br)
	}

	return res, nil
}

func addReach(br *connections.BitRate, numBands int, bandName string, reachVal int) {
	var foundReach *connections.Reach
	for i := range br.Reachs {
		if br.Reachs[i].NumberOfBands == numBands {
			foundReach = &br.Reachs[i]
			break
		}
	}

	if foundReach == nil {
		br.Reachs = append(br.Reachs, connections.Reach{
			NumberOfBands: numBands,
		})
		foundReach = &br.Reachs[len(br.Reachs)-1]
	}

	foundReach.Reach = append(foundReach.Reach, connections.ReachPerBand{
		Band:  bandName,
		Reach: reachVal,
	})
}

// LoadRoutes parses legacy routes JSON files.
func (l *LegacyLoader) LoadRoutes(routesPath string) (connections.Routes, error) {
	data, err := os.ReadFile(routesPath)
	if err != nil {
		return connections.Routes{}, fmt.Errorf("reading legacy routes file: %w", err)
	}

	var routes connections.Routes
	if err := json.Unmarshal(data, &routes); err != nil {
		return connections.Routes{}, fmt.Errorf("parsing legacy routes JSON: %w", err)
	}

	return routes, nil
}
