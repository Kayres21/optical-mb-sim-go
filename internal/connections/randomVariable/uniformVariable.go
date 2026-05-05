package randomvariable

func (rv *RandomVariable) GetNetValueUniform(key UniformKey) int {
	switch key {
	case KeyBitrate:
		return rv.BitrateSelect.Rng.Intn(rv.BitrateSelect.Parameter)
	case KeySource:
		return rv.SourceNodeSelect.Rng.Intn(rv.SourceNodeSelect.Parameter)
	case KeyDestination:
		return rv.DestinationNodeSelect.Rng.Intn(rv.DestinationNodeSelect.Parameter)
	case KeyBand:
		return rv.BandSelect.Rng.Intn(rv.BandSelect.Parameter)
	case KeyGigabits:
		selected := rv.GigabitsSelected.Rng.Intn(rv.GigabitsSelected.Parameter)
		return DefaultGigabitOptions[selected]
	}
	return -1
}
