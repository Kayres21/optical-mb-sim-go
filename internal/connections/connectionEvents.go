package connections

import (
	"sort"
	"strconv"

	randomvariable "github.com/Kayres21/optical-mb-sim-go/internal/connections/randomVariable"
)

type ConnectionEvent struct {
	Id                   string
	Source               int
	Destination          int
	Bitrate              int
	GigabitsSelected     int
	Event                EventsType // "Arrive", "Release"
	Time                 float64
	ConnectionAssignedId string
}

type EventsType string

const ConnectionEventTypeArrive EventsType = "Arrive"
const ConnectionEventTypeRelease EventsType = "Release"

func GenerateEvents(nodeCount int, randomVariable randomvariable.RandomVariable) []ConnectionEvent {
	events := make([]ConnectionEvent, 0, nodeCount*nodeCount)
	id := 0

	for i := range nodeCount {
		for j := range nodeCount {
			if i != j {
				event := ConnectionEvent{
					Id:                   strconv.Itoa(id),
					Source:               i,
					Destination:          j,
					Bitrate:              randomVariable.GetNetValueUniform(randomvariable.KeyBitrate),
					GigabitsSelected:     randomVariable.GetNetValueUniform(randomvariable.KeyGigabits),
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
