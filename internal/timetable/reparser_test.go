package timetable

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
)

const sampleEvents = `[{"title":"Algorithms","teacher":"Alice","description":"TP","location":"Room A","start":"2026-10-13T08:00:01.903+02:00","end":"2026-10-13T09:00:01.903+02:00","groups":["A"]}]`

func testUniversity(baseURL string) *config.UniversityConfig {
	return &config.UniversityConfig{
		ID: 1, AdeProjectId: 1,
		SourceOverride: &config.SourceOverrideConfig{Provider: "reparser", BaseURL: baseURL},
		Rooms:          []config.RoomConfig{{AdeResources: 10}},
		Groups:         []config.GroupConfig{{ID: 2, Timetables: []config.TimetableConfig{{AdeResources: 20}}}},
	}
}

func TestReparserRoutesAndMapping(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.RequestURI()
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method %s", r.Method)
		}
		_, _ = fmt.Fprint(w, sampleEvents)
	}))
	defer srv.Close()
	u := testUniversity(srv.URL)
	client := NewReparser(srv.Client())
	for resource, want := range map[int]string{20: "/v3/univs/1/groups/2/20/events", 10: "/v3/univs/1/rooms/10/events"} {
		events, err := client.Fetch(context.Background(), Request{u, resource})
		if err != nil {
			t.Fatal(err)
		}
		if path != want {
			t.Fatalf("path = %s, want %s", path, want)
		}
		if len(events) != 1 {
			t.Fatalf("events = %v", events)
		}
		event := events[0]
		if event.Title != "Algorithms" || event.Teacher != "Alice" || event.Description != "TP" || event.Location != "Room A" || event.Start.Format(time.RFC3339Nano) != "2026-10-13T06:00:01.903Z" {
			t.Fatalf("incorrect mapping: %+v", event)
		}
	}
	path = ""
	if _, err := client.Fetch(context.Background(), Request{u, 999}); err == nil || path != "" {
		t.Fatal("unknown resource reached upstream")
	}
}

func TestDecodeEventsValidation(t *testing.T) {
	for _, body := range []string{"", "null", "{}", "<html>error</html>", "[null]", `[{"title":"A"}]`, `[{"title":"A","start":"2026-01-01T10:00:00Z","end":"2026-01-01T09:00:00Z"}]`} {
		if _, err := decodeEvents([]byte(body)); err == nil {
			t.Errorf("accepted %q", body)
		}
	}
	events, err := decodeEvents([]byte("[]"))
	if err != nil || events == nil || len(events) != 0 {
		t.Fatalf("empty calendar: %v, %v", events, err)
	}
}

func TestReparserDoesNotFollowRedirects(t *testing.T) {
	var reached atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { reached.Store(true) }))
	defer target.Close()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, target.URL, http.StatusFound) }))
	defer srv.Close()
	_, err := NewReparser(nil).Fetch(context.Background(), Request{testUniversity(srv.URL), 20})
	if err != ErrUnavailable || reached.Load() {
		t.Fatalf("redirect followed or leaked error: %v", err)
	}
}

func TestReparserRetriesTransientFailure(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = fmt.Fprint(w, "[]")
	}))
	defer srv.Close()
	_, err := NewReparser(nil).Fetch(context.Background(), Request{testUniversity(srv.URL), 20})
	if err != nil || calls.Load() != 2 {
		t.Fatalf("calls %d, error %v", calls.Load(), err)
	}
}
