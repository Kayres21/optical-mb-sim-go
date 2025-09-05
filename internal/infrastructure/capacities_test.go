package infrastructure

import (
	"strconv"
	"testing"
)

func TestCapacityGetBands(t *testing.T) {
	capacity := Capacity{
		Bands: []Band{
			{ID: "band1"},
			{ID: "band2"},
		},
	}
	bands := capacity.GetBands()
	if len(bands) != 2 {
		t.Errorf("Expected 2 bands, got %d", len(bands))
	}
	if bands[0].GetID() != "band1" || bands[1].GetID() != "band2" {
		t.Errorf("Band IDs do not match expected values")
	}
}

func TestBandGetID(t *testing.T) {
	band := Band{ID: "test-band"}
	if band.GetID() != "test-band" {
		t.Errorf("Expected ID to be 'test-band', got '%s'", band.GetID())
	}
}

func TestBandGetSlotsLen(t *testing.T) {
	band := Band{SlotsLen: 10}
	if band.GetSlotsLen() != 10 {
		t.Errorf("Expected SlotsLen to be 10, got %d", band.GetSlotsLen())
	}
}

func TestBandGetName(t *testing.T) {
	band := Band{Name: "C-Band"}
	if band.GetName() != "C-Band" {
		t.Errorf("Expected Name to be 'C-Band', got '%s'", band.GetName())
	}
}

func TestBandGetAndSetSlots(t *testing.T) {
	band := Band{}
	slots := []bool{false, true, false}
	band.SetSlots(slots)
	retrievedSlots := band.GetSlots()
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
	if len(capacity.GetBands()) == 0 {
		t.Errorf("Expected at least one band, got 0")
	}
	for _, band := range capacity.Bands {
		if len(band.Slots) != band.SlotsLen {
			t.Errorf("Expected slots length to be %d, got %d", band.SlotsLen, len(band.Slots))
		}
	}

	for i, band := range capacity.Bands {
		expectedID := strconv.Itoa(i)
		if band.GetID() != expectedID {
			t.Errorf("Expected band ID to be '%s', got '%s'", expectedID, band.ID)
		}

		if band.GetSlotsLen() != 344 && band.GetSlotsLen() != 480 {
			t.Errorf("Expected band SlotsLen to be 344 or 480, got %d", band.SlotsLen)
		}

		if band.GetName() != "C" && band.GetName() != "L" {
			t.Errorf("Expected band Name to be 'C' or 'L', got '%s'", band.Name)
		}
	}
}
