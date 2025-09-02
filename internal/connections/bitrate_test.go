package connections

import "testing"

func TestReadBirateFile(t *testing.T) {
	routesFile := "../../files/bitrate/bitrate_test.json"
	bitrate, err := ReadBitRateFile(routesFile)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(bitrate.GetBitRates()) != 2 {
		t.Errorf("Expected 2 bitrates, got %d", len(bitrate.BitRates))
	}

}
