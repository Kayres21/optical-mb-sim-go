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
}

type ConnectionEvent struct {
	Id          string
	Source      string
	Destination string
	Bitrate     int
}

type BitRate struct {
	Name string
}

type Routes struct {
	Alias string
	Name  string
}
