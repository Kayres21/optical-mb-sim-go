package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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

		link.UpdateAllFragmentationRatios()

		network.Links = append(network.Links, link)
	}

	return network, nil
}

// LoadBitRate parses legacy bitrate JSON files, mirroring
// flexnetsim.bitrate.BitRate.read_bit_rate_file / read_bit_rate_file_mb: one
// BitRate is produced per magnitude (gigabits key), holding every modulation
// found for it. Single-band and multi-band (per-band) modulation configs may
// be mixed freely within the same file, detected per modulation entry.
func (l *LegacyLoader) LoadBitRate(bitRatePath string, _ int) (connections.BitRateList, error) {
	data, err := os.ReadFile(bitRatePath)
	if err != nil {
		return connections.BitRateList{}, fmt.Errorf("reading legacy bitrate file: %w", err)
	}

	// Legacy format is map[gigabits][]map[modulation]config
	var raw map[string][]map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return connections.BitRateList{}, fmt.Errorf("parsing legacy bitrate JSON: %w", err)
	}

	keys := make([]string, 0, len(raw))
	for gigabits := range raw {
		keys = append(keys, gigabits)
	}
	sort.Slice(keys, func(i, j int) bool {
		left, _ := strconv.ParseFloat(keys[i], 64)
		right, _ := strconv.ParseFloat(keys[j], 64)
		return left < right
	})

	var res connections.BitRateList
	for _, gigabits := range keys {
		value, err := strconv.ParseFloat(gigabits, 64)
		if err != nil {
			return connections.BitRateList{}, fmt.Errorf("parsing bitrate magnitude %q: %w", gigabits, err)
		}
		br := connections.BitRate{Value: value}

		configs := make(map[string]json.RawMessage)
		modulations := make([]string, 0)
		for _, entry := range raw[gigabits] {
			for modulation, configRaw := range entry {
				modulations = append(modulations, modulation)
				configs[modulation] = configRaw
			}
		}
		sort.Slice(modulations, func(i, j int) bool {
			return modulationSortKey(modulations[i]) < modulationSortKey(modulations[j])
		})

		for _, modulation := range modulations {
			configRaw := configs[modulation]

			// Check if configRaw is a single config or a band map
			var singleConfig struct {
				Slots int     `json:"slots"`
				Reach float64 `json:"reach"`
			}

			if err := json.Unmarshal(configRaw, &singleConfig); err == nil {
				if err := validateSlotsReach(singleConfig.Slots, singleConfig.Reach); err != nil {
					return connections.BitRateList{}, fmt.Errorf("bitrate %q modulation %q: %w", gigabits, modulation, err)
				}
				br.AddModulation(modulation, singleConfig.Slots, singleConfig.Reach, nil, nil, nil)
				continue
			}

			// Multi band legacy format: [ { "C": {...} }, { "L": {...} } ]
			var bandConfigs []map[string]struct {
				Slots int     `json:"slots"`
				Reach float64 `json:"reach"`
			}
			if err := json.Unmarshal(configRaw, &bandConfigs); err != nil {
				return connections.BitRateList{}, fmt.Errorf("parsing modulation %q for bitrate %q: %w", modulation, gigabits, err)
			}

			var bandNames []string
			var slotsPerBand []int
			var reachPerBand []float64
			totalSlots := 0
			totalReach := 0.0
			for _, bc := range bandConfigs {
				for bandName, v := range bc {
					if err := validateSlotsReach(v.Slots, v.Reach); err != nil {
						return connections.BitRateList{}, fmt.Errorf("bitrate %q modulation %q band %q: %w", gigabits, modulation, bandName, err)
					}
					bandNames = append(bandNames, bandName)
					slotsPerBand = append(slotsPerBand, v.Slots)
					reachPerBand = append(reachPerBand, v.Reach)
					totalSlots += v.Slots
					totalReach += v.Reach
				}
			}
			br.AddModulation(modulation, totalSlots, totalReach, bandNames, slotsPerBand, reachPerBand)
		}

		res.BitRates = append(res.BitRates, br)
	}

	return res, nil
}

func modulationSortKey(modulation string) int {
	normalized := strings.ToUpper(strings.ReplaceAll(modulation, "-", ""))
	switch normalized {
	case "BPSK":
		return 0
	case "QPSK":
		return 1
	case "8QAM":
		return 2
	case "16QAM":
		return 3
	default:
		return 1000
	}
}

func validateSlotsReach(slots int, reach float64) error {
	switch {
	case slots < 0 && reach < 0:
		return fmt.Errorf("value entered for slots and reach is less than zero")
	case reach < 0:
		return fmt.Errorf("value entered for reach is less than zero")
	case slots < 0:
		return fmt.Errorf("value entered for slots is less than zero")
	}
	return nil
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
