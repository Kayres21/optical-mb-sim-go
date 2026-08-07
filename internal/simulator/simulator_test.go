package simulator

import (
	"os"
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

	if slot := s.getSlotsByGigabits(br, 100); slot.Slots != 2 {
		t.Errorf("expected 2 slots, got %d", slot.Slots)
	}

	if slot := s.getSlotsByGigabits(br, 200); slot.Slots != 4 {
		t.Errorf("expected 4 slots, got %d", slot.Slots)
	}

	if slot := s.getSlotsByGigabits(br, 300); slot.Slots != 0 {
		t.Errorf("expected 0 slots for unknown gigabits, got %d", slot.Slots)
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

func TestSimulator_ResolveConnectionID(t *testing.T) {
	s := &Simulator{}
	event := connections.ConnectionEvent{Id: "42", ConnectionAssignedId: "99"}

	if got := s.resolveConnectionID(event); got != "42" {
		t.Fatalf("expected event ID to be used as the connection identifier, got %s", got)
	}
}

func TestSimulator_SaveEventsCSV(t *testing.T) {
	s := &Simulator{}
	s.generatedEvents = []connections.ConnectionEvent{
		{Id: "1", Source: 3, Destination: 7, Bitrate: 2, Event: connections.ConnectionEventTypeArrive, Time: 0.25, ConnectionAssignedId: "1"},
		{Id: "1", Event: connections.ConnectionEventTypeRelease, Time: 0.42, ConnectionAssignedId: "1"},
	}

	path := t.TempDir() + "/events.csv"
	if err := s.SaveEventsCSV(path); err != nil {
		t.Fatalf("SaveEventsCSV returned error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading generated CSV: %v", err)
	}

	got := string(data)
	want := "ID,Tiempo,Evento,Source,Destination,BitRate\r\n1,0.25,ARRIVE,3,7,2\r\n1,0.42,DEPARTURE,0,0,0\r\n"
	if got != want {
		t.Fatalf("unexpected CSV content:\n got: %q\nwant: %q", got, want)
	}
}
