package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
)

func TestSourceOverridePublicContract(t *testing.T) {
	cfg := testConfig()
	cfg.Universities[0].SourceOverride = &config.SourceOverrideConfig{Provider: "reparser", BaseURL: "https://private-upstream.example.test"}
	srv := newTestServer(t, cfg)
	defer srv.Close()
	for _, path := range []string{"/v3/univs", "/v3/univs/1", "/v3/univs/1/rooms", "/v3/univs/1/rooms/10", "/openapi.json"} {
		resp := get(t, srv, path)
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s: status %d", path, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		for _, hidden := range []string{"sourceOverride", "reparser", "experimental"} {
			if strings.Contains(strings.ToLower(string(body)), strings.ToLower(hidden)) {
				t.Errorf("%s exposes %s", path, hidden)
			}
		}
	}
	resp := get(t, srv, "/v3/univs/1/rooms/999/events")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unknown room: %d", resp.StatusCode)
	}
	for path, want := range map[string]string{
		"/v3/univs/1/rooms/10":    "https://private-upstream.example.test/v3/ics/1/10",
		"/v3/univs/1/groups/1/20": "https://private-upstream.example.test/v3/ics/1/20",
	} {
		var metadata struct {
			AdeURL string `json:"adeUrl"`
		}
		if err := json.NewDecoder(get(t, srv, path).Body).Decode(&metadata); err != nil {
			t.Fatal(err)
		}
		if metadata.AdeURL != want {
			t.Fatalf("%s: subscription URL %q, want %q", path, metadata.AdeURL, want)
		}
	}
	for _, busy := range []bool{false, true} {
		events := []ade.Event{}
		if busy {
			events = append(events, ade.Event{Title: "Occupied", Start: time.Now().Add(-time.Hour), End: time.Now().Add(time.Hour)})
		}
		seedEvents(1, 10, events)
		resp = get(t, srv, "/v3/univs/1/free-rooms?campusId=1")
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("free rooms: %d", resp.StatusCode)
		}
		var rooms []roomResponse
		if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
			t.Fatal(err)
		}
		want := 1
		if busy {
			want = 0
		}
		if len(rooms) != want {
			t.Fatalf("busy=%v: got %d free rooms, want %d", busy, len(rooms), want)
		}
	}
}
