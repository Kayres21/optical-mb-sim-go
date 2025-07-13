package infrastructure

type Node struct {
	ID string `json:"id"`
}

type Link struct {
	ID          string              `json:"id"`
	Source      string              `json:"src"`
	Destination string              `json:"dst"`
	Length      int                 `json:"length"`
	Capacities  map[string]Capacity `json:"-"`
}

type Capacity struct {
	Bands []Bands `json:"bands"`
}

type Bands struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
	Slots    []bool `json:"-"`
}

type Network struct {
	Name  string          `json:"name"`
	Alias string          `json:"alias"`
	Nodes map[string]Node `json:"nodes"`
	Links map[string]Link `json:"links"`
}
