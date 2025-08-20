package connections

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
	"github.com/google/uuid"
)

type Connection struct {
	Id           string
	Links        []infrastructure.Link
	Source       int
	Destination  int
	BitRate      string
	Slots        int
	InitialSlot  int
	FinalSlot    int
	BandSelected int
	Allocated    bool
}

type ConnectionEvent struct {
	Id          string
	Source      int
	Destination int
	Bitrate     int
	Event       EventsType // "Arrive",  "Release"
	Time        float64
}

type BitRateList struct {
	BitRates []BitRate `json:"bitrates"`
}

type BitRate struct {
	Modulation string  `json:"modulation"`
	Slots      []Slots `json:"slots"`
	Reachs     []Reach `json:"reach"`
}

type Slots struct {
	Gigabits string `json:"gigabits"`
	Slots    int    `json:"slots"`
}

type Reach struct {
	NumberOfBands int            `json:"number_of_bands"`
	Reach         []ReachPerBand `json:"reach"`
}

type ReachPerBand struct {
	Band  string `json:"band"`
	Reach int    `json:"reach"`
}

type Routes struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
	Paths []Path `json:"routes"`
}

type Path struct {
	Source      int     `json:"src"`
	Destination int     `json:"dst"`
	PathLinks   [][]int `json:"paths"`
}

type EventsType string

const ConnectionEventTypeArrive EventsType = "Arrive"
const ConnectionEventTypeRelease EventsType = "Release"

func GenerateConnectionID() string {
	// Implement a simple ID generation logic (you can replace this with a more robust method)
	return uuid.New().String()
}
