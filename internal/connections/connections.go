package connections

import (
	"simulator/internal/infrastructure"
)

type Connection struct {
	Id           string
	Links        []infrastructure.Link
	Slots        int
	InitialSlot  int
	FinalSlot    int
	BandSelected string
	Allocated    bool
}

type ConnectionEvent struct {
	Id          string
	Source      string
	Destination string
	Bitrate     int
	EventType   string // "Arrive",  "Release"
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
	PathLinks   [][]int `json:"path"`
}
