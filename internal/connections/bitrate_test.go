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

func TestTransformIntToModulation(t *testing.T) {

	if TrasnformIntToModulation(0) != "BPSK" {
		t.Errorf("Expected BPSK, got %s", TrasnformIntToModulation(0))
	}
	if TrasnformIntToModulation(1) != "QPSK" {
		t.Errorf("Expected QPSK, got %s", TrasnformIntToModulation(0))
	}
	if TrasnformIntToModulation(2) != "8-QAM" {
		t.Errorf("Expected 8-QAM, got %s", TrasnformIntToModulation(0))
	}
	if TrasnformIntToModulation(3) != "16-QAM" {
		t.Errorf("Expected 16-QAM, got %s", TrasnformIntToModulation(0))
	}

}
