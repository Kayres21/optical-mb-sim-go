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
