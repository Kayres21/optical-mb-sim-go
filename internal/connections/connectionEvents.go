package connections

import (
	"sort"
	"strconv"

	randomvariable "github.com/Kayres21/optical-mb-sim-go/internal/connections/randomVariable"
)

type ConnectionEvent struct {
	Id               string
	Source           int
	Destination      int
	Bitrate          int
	GigabitsSelected int
	Event            EventsType // "Arrive",  "Release"
	Time             float64
}

const ConnectionEventTypeArrive EventsType = "Arrive"
const ConnectionEventTypeRelease EventsType = "Release"

func GenerateEvents(nodes_len int, randomVariable randomvariable.RandomVariable) []ConnectionEvent {

	events := make([]ConnectionEvent, 0)
	id := 0

	for i := range nodes_len {
		for j := range nodes_len {
			if i != j {
				event := ConnectionEvent{
					Id:               strconv.Itoa(id),
					Source:           i,
					Destination:      j,
					Bitrate:          randomVariable.GetNetValueUniform("bitrate"),
					Event:            "Arrive",
					GigabitsSelected: randomVariable.GetNetValueUniform("gigabits"),
					Time:             randomVariable.GetNetValueExponential("arrive"),
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
