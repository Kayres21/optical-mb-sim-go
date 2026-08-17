package loader

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLegacyLoader_LoadBitRateUsesBitrateClasses(t *testing.T) {
	root := filepath.Join("..", "..")
	bitratePath := filepath.Join(root, "legacy_files", "bitrates", "bitrate_iroBand_C.json")

	loader := &LegacyLoader{}
	bitrateList, err := loader.LoadBitRate(bitratePath, 1)
	if err != nil {
		t.Fatalf("LoadBitRate returned error: %v", err)
	}

	if len(bitrateList.BitRates) != 5 {
		t.Fatalf("expected 5 bitrate classes, got %d", len(bitrateList.BitRates))
	}

	if got := bitrateList.BitRates[0].Value; got != 10 {
		t.Fatalf("expected first bitrate class to be 10, got %v", got)
	}
}

func TestLegacyLoader_LoadBitRateIsDeterministic(t *testing.T) {
	root := filepath.Join("..", "..")
	bitratePath := filepath.Join(root, "legacy_files", "bitrates", "bitrate_iroBand_C.json")

	loader := &LegacyLoader{}
	var first []string

	for i := 0; i < 10; i++ {
		bitrateList, err := loader.LoadBitRate(bitratePath, 1)
		if err != nil {
			t.Fatalf("LoadBitRate iteration %d returned error: %v", i, err)
		}

		var summary []string
		for _, br := range bitrateList.BitRates {
			var modSummary []string
			for m, modulation := range br.Modulation {
				modSummary = append(modSummary, fmt.Sprintf("%s:%d", modulation, br.Slots[m]))
			}
			summary = append(summary, fmt.Sprintf("%v[%s]", br.Value, strings.Join(modSummary, ",")))
		}

		if i == 0 {
			first = summary
			continue
		}

		if !reflect.DeepEqual(summary, first) {
			t.Fatalf("legacy bitrate loader changed between runs:\nfirst: %v\nnext:  %v", first, summary)
		}
	}
}
