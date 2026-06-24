package connections

import (
	"sort"
	"strconv"

	randomvariable "github.com/Kayres21/optical-mb-sim-go/internal/connections/randomVariable"
)

type ConnectionEvent struct {
	Id                     string
	Source                 int
	Destination            int
	Bitrate                int
	GigabitsSelected       int
	Event                  EventsType // "Arrive", "Release"
	Time                   float64
	ConnectionAssignedId   string
	ConnectionInitialSlot  int
	ConnectionSlots        int
	ConnectionBandSelected int
}

type EventsType string

const ConnectionEventTypeArrive EventsType = "Arrive"
const ConnectionEventTypeRelease EventsType = "Release"

func GenerateEvents(nodeCount int, randomVariable randomvariable.RandomVariable, bitRateList BitRateList) []ConnectionEvent {
	events := make([]ConnectionEvent, 0, nodeCount*nodeCount)
	id := 0

	for i := range nodeCount {
		for j := range nodeCount {
			if i != j {
				unifiedIndex := randomVariable.GetNetValueUniform(randomvariable.KeyBitrate)
				var modulationIndex, gigabits int

				if len(bitRateList.BitRates) > 0 && len(bitRateList.BitRates[0].Slots) > 0 {
					slotsCount := len(bitRateList.BitRates[0].Slots)
					modulationIndex = unifiedIndex / slotsCount
					slotIndex := unifiedIndex % slotsCount

					gigaStr := bitRateList.BitRates[modulationIndex].Slots[slotIndex].Gigabits
					gigabits, _ = strconv.Atoi(gigaStr)
				} else {
					modulationIndex = 0
					gigabits = 10
				}

				event := ConnectionEvent{
					Id:                   strconv.Itoa(id),
					Source:               i,
					Destination:          j,
					Bitrate:              modulationIndex,
					GigabitsSelected:     gigabits,
					Event:                ConnectionEventTypeArrive,
					Time:                 randomVariable.GetNetValueExponential(randomvariable.KeyArrive),
					ConnectionAssignedId: "",
				}
				events = append(events, event)
				id++
			}
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time < events[j].Time
	})

	return events
}
