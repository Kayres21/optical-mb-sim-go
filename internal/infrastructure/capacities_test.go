package infrastructure

import (
	"strconv"
	"testing"
)

func TestCapacityBands(t *testing.T) {
	capacity := Capacity{
		Bands: []Band{
			{ID: "band1"},
			{ID: "band2"},
		},
	}
	bands := capacity.Bands
	if len(bands) != 2 {
		t.Errorf("Expected 2 bands, got %d", len(bands))
	}
	if bands[0].ID != "band1" || bands[1].ID != "band2" {
		t.Errorf("Band IDs do not match expected values")
	}
}

func TestBandID(t *testing.T) {
	band := Band{ID: "test-band"}
	if band.ID != "test-band" {
		t.Errorf("Expected ID to be 'test-band', got '%s'", band.ID)
	}
}

func TestBandSlotsLen(t *testing.T) {
	band := Band{SlotsLen: 10}
	if band.SlotsLen != 10 {
		t.Errorf("Expected SlotsLen to be 10, got %d", band.SlotsLen)
	}
}

func TestBandName(t *testing.T) {
	band := Band{Name: "C-Band"}
	if band.Name != "C-Band" {
		t.Errorf("Expected Name to be 'C-Band', got '%s'", band.Name)
	}
}

func TestBandSlots(t *testing.T) {
	band := Band{}
	slots := []bool{false, true, false}
	band.Slots = slots
	retrievedSlots := band.Slots
	if len(retrievedSlots) != len(slots) {
		t.Errorf("Expected slots length to be %d, got %d", len(slots), len(retrievedSlots))
	}
	for i, slot := range retrievedSlots {
		if slot != slots[i] {
			t.Errorf("Expected slot %d to be %v, got %v", i, slots[i], slot)
		}
	}
}

func TestReadCapacityFile(t *testing.T) {

	capacityPath := "../../files/capacities/capacities_test.json"

	capacity, err := ReadCapacityFile(capacityPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(capacity.Bands) == 0 {
		t.Errorf("Expected at least one band, got 0")
	}
	for _, band := range capacity.Bands {
		if len(band.Slots) != band.SlotsLen {
			t.Errorf("Expected slots length to be %d, got %d", band.SlotsLen, len(band.Slots))
		}
	}

	for i, band := range capacity.Bands {
		expectedID := strconv.Itoa(i)
		if band.ID != expectedID {
			t.Errorf("Expected band ID to be '%s', got '%s'", expectedID, band.ID)
		}

		if band.SlotsLen != 344 && band.SlotsLen != 480 {
			t.Errorf("Expected band SlotsLen to be 344 or 480, got %d", band.SlotsLen)
		}

		if band.Name != "C" && band.Name != "L" {
			t.Errorf("Expected band Name to be 'C' or 'L', got '%s'", band.Name)
		}
	}
}
