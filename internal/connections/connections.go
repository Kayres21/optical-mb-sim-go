package connections

import (
	"github.com/Kayres21/optical-mb-sim-go/internal/infrastructure"
)

type Connection struct {
	Id           string
	Links        []*infrastructure.Link
	Source       int
	Destination  int
	Slots        int
	InitialSlot  int
	FinalSlot    int
	BandSelected int
	Allocated    bool
}


