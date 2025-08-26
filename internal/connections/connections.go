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
	Slots        int
	InitialSlot  int
	FinalSlot    int
	BandSelected int
	Allocated    bool
}

type ConnectionEvent struct {
	Id               string
	Source           int
	Destination      int
	Bitrate          int
	GigabitsSelected int
	Event            EventsType // "Arrive",  "Release"
	Time             float64
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
	return uuid.New().String()
}

func (c *Connection) GetId() string {
	return c.Id
}

func (c *Connection) GetLinks() []infrastructure.Link {
	return c.Links
}

func (c *Connection) GetSource() int {
	return c.Source
}

func (c *Connection) GetDestination() int {
	return c.Destination
}

func (c *Connection) GetSlots() int {
	return c.Slots
}

func (c *Connection) GetInitialSlot() int {
	return c.InitialSlot
}

func (c *Connection) GetFinalSlot() int {
	return c.FinalSlot
}

func (c *Connection) GetBandSelected() int {
	return c.BandSelected
}

func (c *Connection) IsAllocated() bool {
	return c.Allocated
}

func (c *Connection) SetId(id string) {
	c.Id = id
}

func (c *Connection) SetLinks(links []infrastructure.Link) {
	c.Links = links
}

func (c *Connection) SetSource(source int) {
	c.Source = source
}

func (c *Connection) SetDestination(destination int) {
	c.Destination = destination
}

func (c *Connection) SetSlots(slots int) {
	c.Slots = slots
}

func (c *Connection) SetInitialSlot(initialSlot int) {
	c.InitialSlot = initialSlot
}

func (c *Connection) SetFinalSlot(finalSlot int) {
	c.FinalSlot = finalSlot
}

func (c *Connection) SetBandSelected(band int) {
	c.BandSelected = band
}

func (c *Connection) SetAllocated(allocated bool) {
	c.Allocated = allocated
}
