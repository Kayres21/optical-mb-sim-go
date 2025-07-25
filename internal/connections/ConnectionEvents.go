package connections

import (
	"log"
	randomvariable "simulator/internal/connections/random_variable"
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
				}
				log.Printf("Generated event: %+v\n", event)
				events = append(events, event)
			}

		}
	}
	return events
}
