package main

import (
	"reflect"
	"testing"
)

func TestFilterConfigs(t *testing.T) {
	configPaths := []string{
		"configs/NET_A/config-defrag-before-arrival.json",
		"configs/NET_A/config-defrag-multiband-same-band.json",
		"configs/NET_A/config.json",
		"configs/NET_B/config-defrag-before-arrival.json",
		"configs/NET_B/config-defrag-multiband-same-band.json",
		"configs/NET_B/config.json",
	}

	tests := []struct {
		name     string
		selected map[string]bool
		want     []string
	}{
		{
			name:     "no selection returns all configs",
			selected: map[string]bool{},
			want:     configPaths,
		},
		{
			name:     "select one variant",
			selected: map[string]bool{"config.json": true},
			want: []string{
				"configs/NET_A/config.json",
				"configs/NET_B/config.json",
			},
		},
		{
			name: "select multiple variants",
			selected: map[string]bool{
				"config-defrag-before-arrival.json":      true,
				"config-defrag-multiband-same-band.json": true,
			},
			want: []string{
				"configs/NET_A/config-defrag-before-arrival.json",
				"configs/NET_A/config-defrag-multiband-same-band.json",
				"configs/NET_B/config-defrag-before-arrival.json",
				"configs/NET_B/config-defrag-multiband-same-band.json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterConfigs(configPaths, tt.selected)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}
