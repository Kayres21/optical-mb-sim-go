package simulator

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
)

func TestSimulator_PushAndPopEvent(t *testing.T) {
	s := &Simulator{}

	event1 := connections.ConnectionEvent{Id: "1", Time: 10.0}
	event2 := connections.ConnectionEvent{Id: "2", Time: 5.0}
	event3 := connections.ConnectionEvent{Id: "3", Time: 15.0}

	s.pushEvent(event1)
	s.pushEvent(event2)
	s.pushEvent(event3)

	if s.events.Len() != 3 {
		t.Fatalf("expected 3 events in heap, got %d", s.events.Len())
	}

	// Should pop the one with lowest time (event2, time 5.0)
	first := s.popEvent()
	if first.Id != "2" {
		t.Errorf("expected first event ID to be 2, got %s", first.Id)
	}

	second := s.popEvent()
	if second.Id != "1" {
		t.Errorf("expected second event ID to be 1, got %s", second.Id)
	}

	third := s.popEvent()
	if third.Id != "3" {
		t.Errorf("expected third event ID to be 3, got %s", third.Id)
	}

	if s.events.Len() != 0 {
		t.Errorf("expected empty heap, got %d elements", s.events.Len())
	}
}

func TestSimulator_getSlotsByGigabits(t *testing.T) {
	s := &Simulator{}

	br := connections.BitRate{
		Slots: []connections.Slots{
			{Gigabits: "100", Slots: 2},
			{Gigabits: "200", Slots: 4},
		},
	}

	if slots := s.getSlotsByGigabits(br, 100); slots != 2 {
		t.Errorf("expected 2 slots, got %d", slots)
	}

	if slots := s.getSlotsByGigabits(br, 200); slots != 4 {
		t.Errorf("expected 4 slots, got %d", slots)
	}

	if slots := s.getSlotsByGigabits(br, 300); slots != 0 {
		t.Errorf("expected 0 slots for unknown gigabits, got %d", slots)
	}
}

func TestSimulator_AddResultsAndArrives(t *testing.T) {
	s := &Simulator{}

	s.addResult(0.5)
	if len(s.results) != 1 || s.results[0] != 0.5 {
		t.Errorf("expected results to have [0.5], got %v", s.results)
	}

	s.addArrive(10.0)
	if len(s.arrives) != 1 || s.arrives[0] != 10.0 {
		t.Errorf("expected arrives to have [10.0], got %v", s.arrives)
	}
}
