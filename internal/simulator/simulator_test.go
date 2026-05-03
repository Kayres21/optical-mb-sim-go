package simulator

import (
	"testing"

	"github.com/Kayres21/optical-mb-sim-go/internal/connections"
)

func TestSimulator_AddNewConnectionEvent_And_GetFirstEvent(t *testing.T) {
	s := &Simulator{
		ConnectionsEvents: []connections.ConnectionEvent{},
	}

	event1 := connections.ConnectionEvent{Id: "1", Time: 10.0}
	event2 := connections.ConnectionEvent{Id: "2", Time: 5.0}
	event3 := connections.ConnectionEvent{Id: "3", Time: 15.0}

	s.AddNewConnectionEvent(event1)
	s.AddNewConnectionEvent(event2)
	s.AddNewConnectionEvent(event3)

	if len(s.ConnectionsEvents) != 3 {
		t.Fatalf("expected 3 events, got %d", len(s.ConnectionsEvents))
	}

	// Should pop the one with lowest time (event2, time 5.0)
	first := s.GetFirstEvent()
	if first.Id != "2" {
		t.Errorf("expected first event ID to be 2, got %s", first.Id)
	}

	second := s.GetFirstEvent()
	if second.Id != "1" {
		t.Errorf("expected second event ID to be 1, got %s", second.Id)
	}

	third := s.GetFirstEvent()
	if third.Id != "3" {
		t.Errorf("expected third event ID to be 3, got %s", third.Id)
	}

	// Queue should be empty now
	if len(s.ConnectionsEvents) != 0 {
		t.Errorf("expected empty queue, got %d elements", len(s.ConnectionsEvents))
	}

	// Popping from empty should not panic and return zero value
	empty := s.GetFirstEvent()
	if empty.Id != "" {
		t.Errorf("expected empty event, got %+v", empty)
	}
}

func TestSimulator_getSlotgigabites(t *testing.T) {
	s := &Simulator{}

	br := connections.BitRate{
		Slots: []connections.Slots{
			{Gigabits: "100", Slots: 2},
			{Gigabits: "200", Slots: 4},
		},
	}

	slots := s.getSlotgigabites(br, 100)
	if slots != 2 {
		t.Errorf("expected 2 slots, got %d", slots)
	}

	slots = s.getSlotgigabites(br, 200)
	if slots != 4 {
		t.Errorf("expected 4 slots, got %d", slots)
	}

	slots = s.getSlotgigabites(br, 300)
	if slots != 0 {
		t.Errorf("expected 0 slots for unknown gigabits, got %d", slots)
	}
}

func TestSimulator_AddResultsAndArrives(t *testing.T) {
	s := &Simulator{}

	s.addResult(0.5)
	if len(s.Results) != 1 || s.Results[0] != 0.5 {
		t.Errorf("expected Results to have [0.5], got %v", s.Results)
	}

	s.addArrive(10.0)
	if len(s.Arrives) != 1 || s.Arrives[0] != 10.0 {
		t.Errorf("expected Arrives to have [10.0], got %v", s.Arrives)
	}
}
