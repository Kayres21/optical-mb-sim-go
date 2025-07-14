package infrastructure

type Node struct {
	ID int `json:"id"`
}

type Link struct {
	ID          int      `json:"id"`
	Source      int      `json:"src"`
	Destination int      `json:"dst"`
	Length      int      `json:"length"`
	Capacities  Capacity `json:"-"`
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
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}
