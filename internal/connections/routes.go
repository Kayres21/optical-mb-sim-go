package connections

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/Kayres21/optical-mb-sim-go/pkg/validator"
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
	schemaPath := filepath.Join(filepath.Dir(routesPath), "schema.json")
	dataBytesNetwork, err := validator.ValidateFile(routesPath, schemaPath)
	if err != nil {
		return Routes{}, fmt.Errorf("validating routes file: %w", err)
	}

	var routes Routes

	if err = json.Unmarshal(dataBytesNetwork, &routes); err != nil {
		return Routes{}, fmt.Errorf("parsing routes file %q: %w", routesPath, err)
	}

	return routes, nil
}

func (r *Routes) GetKshortestPath(kShortestPath, src, dst int) []int {
	if len(r.Paths) == 0 {
		return []int{}
	}

	for _, path := range r.Paths {
		if path.Source == src && path.Destination == dst {
			if kShortestPath < 0 || kShortestPath >= len(path.PathLinks) {
				return []int{}
			}
			return path.PathLinks[kShortestPath]
		}
	}
	return []int{}
}

func (r *Routes) GetPaths(src, dst int) [][]int {
	paths := make([][]int, 0)
	for _, path := range r.Paths {
		if path.Source == src && path.Destination == dst {
			paths = append(paths, path.PathLinks...)
		}
	}
	return paths
}
