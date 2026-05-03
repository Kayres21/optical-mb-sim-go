package connections

import (
	"encoding/json"
	"log"
	"os"
)

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

func ReadRoutesFile(routesPath string) (Routes, error) {
	dataBytesNetwork, err := os.ReadFile(routesPath)

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return Routes{}, err
	}

	var routes Routes

	err = json.Unmarshal(dataBytesNetwork, &routes)

	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return Routes{}, err
	}

	return routes, nil
}

func (r *Routes) GetKshortestPath(kShortestPath, src, dst int) []int {
	if len(r.Paths) == 0 {
		return []int{}
	}

	for _, path := range r.Paths {
		if path.Source == src && path.Destination == dst {
			return path.PathLinks[kShortestPath]
		}
	}
	return []int{}
}
