package randomvariable

func (rv *RandomVariable) GetNetValueUniform(key UniformKey) int {
	switch key {
	case KeyBitrate:
		return rv.BitrateSelect.Rng.Intn(rv.BitrateSelect.Parameter + 1)
	case KeySource:
		return rv.SourceNodeSelect.Rng.Intn(rv.SourceNodeSelect.Parameter + 1)
	case KeyDestination:
		return rv.DestinationNodeSelect.Rng.Intn(rv.DestinationNodeSelect.Parameter + 1)
	case KeyBand:
		return rv.BandSelect.Rng.Intn(rv.BandSelect.Parameter + 1)
	case KeyGigabits:
		selected := rv.GigabitsSelected.Rng.Intn(rv.GigabitsSelected.Parameter + 1)
		return DefaultGigabitOptions[selected]
	}
	return -1
}
