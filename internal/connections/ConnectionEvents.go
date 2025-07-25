package connections

import (
	"log"
	randomvariable "simulator/internal/connections/random_variable"
	"sort"
	"strconv"
)

func GenerateEvents(nodes_len int, randomVariable randomvariable.RandomVariable) []ConnectionEvent {

	events := make([]ConnectionEvent, 0)

	for i := range nodes_len {
		for j := range nodes_len {
			if i != j {
				event := ConnectionEvent{
					Id:          strconv.Itoa(i) + strconv.Itoa(j),
					Source:      strconv.Itoa(i),
					Destination: strconv.Itoa(j),
					Bitrate:     randomVariable.GetNetValueUniform("bitrate"),
					EventType:   "Arrive",
					Time:        randomVariable.GetNetValueExponential("arrive"),
				}
				log.Printf("Generated event: %+v\n", event)
				events = append(events, event)
			}

		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time < events[j].Time
	})
	return events
}
