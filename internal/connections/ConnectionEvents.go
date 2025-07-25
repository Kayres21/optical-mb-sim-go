package connections

import (
	randomvariable "simulator/internal/connections/random_variable"
	"sort"
	"strconv"
)

func GenerateEvents(nodes_len int, randomVariable randomvariable.RandomVariable) []ConnectionEvent {

	events := make([]ConnectionEvent, 0)
	id := 0

	for i := range nodes_len {
		for j := range nodes_len {
			if i != j {
				event := ConnectionEvent{
					Id:          strconv.Itoa(id),
					Source:      strconv.Itoa(i),
					Destination: strconv.Itoa(j),
					Bitrate:     randomVariable.GetNetValueUniform("bitrate"),
					EventType:   "Arrive",
					Time:        randomVariable.GetNetValueExponential("arrive"),
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
