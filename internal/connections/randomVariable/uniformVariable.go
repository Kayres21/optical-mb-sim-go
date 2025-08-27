package randomvariable

import "math/rand"

func (ev *UniformVariable) SetRng(rng *rand.Rand) {
	ev.Rng = rng
}

func (ev *UniformVariable) SetParameter(parameter int) {
	ev.Parameter = parameter
}

func (rv *RandomVariable) GetNetValueUniform(options string) int {
	switch options {
	case "bitrate":
		return rv.BitrateSelect.Rng.Intn(rv.BitrateSelect.Parameter)
	case "source":
		return rv.SourceNodeSelect.Rng.Intn(rv.SourceNodeSelect.Parameter)
	case "destination":
		return rv.DestinationNodeSelect.Rng.Intn(rv.DestinationNodeSelect.Parameter)
	case "band":
		return rv.BandSelect.Rng.Intn(rv.BandSelect.Parameter)
	case "gigabits":
		var gigabitsArray [5]int = [5]int{10, 40, 100, 400, 1000}
		selected := rv.GigabitsSelected.Rng.Intn(rv.GigabitsSelected.Parameter)
		return gigabitsArray[selected]
	}

	return -1 // Default case, should not happen
}
