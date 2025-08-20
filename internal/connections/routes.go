package connections

import (
	"encoding/json"
	"log"
	"os"
)

func (r *Routes) ReadRoutesFile(routesPath string) (Routes, error) {
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
