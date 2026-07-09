package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AntoninHuaut/EtuEDT-Back/internal/ade"
	"github.com/AntoninHuaut/EtuEDT-Back/internal/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// newTestServer creates an httptest.Server with the full v3 routing,
// seeded with the given config.
func newTestServer(t *testing.T, cfg config.Config) *httptest.Server {
	t.Helper()
	config.AppConfig = cfg

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	humaConfig := huma.DefaultConfig("EtuEDT API", "3.0.0")
	humaAPI := humachi.New(r, humaConfig)
	registerV3Handlers(humaAPI)

	return httptest.NewServer(r)
}

// testConfig returns a minimal valid AppConfig for integration tests.
func testConfig() config.Config {
	return config.Config{
		Universities: []config.UniversityConfig{
			{
				ID:           1,
				Name:         "Test University",
				AdeUrl:       "https://ade.example.com",
				AdeProjectId: 42,
				Rooms: []config.RoomConfig{
					{AdeResources: 10, Label: "Room A"},
				},
				Groups: []config.GroupConfig{
					{
						ID:   1,
						Name: "Group 1",
						Timetables: []config.TimetableConfig{
							{AdeResources: 20, Year: 2025, Label: "Timetable A"},
						},
					},
				},
			},
		},
	}
}

// get is a test helper for making GET requests.
func get(t *testing.T, server *httptest.Server, path string) *http.Response {
	t.Helper()
	resp, err := http.Get(server.URL + path)
	if err != nil {
		t.Fatalf("GET %s failed: %v", path, err)
	}
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("GET %s: failed to close body: %v", path, err)
		}
	})
	return resp
}

// --- v3 university endpoints ---

func TestV3_ListUniversities_ReturnsAll(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	// Huma serializes list Body fields as a plain JSON array
	var body []universityResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body) != 1 {
		t.Errorf("expected 1 university, got %d", len(body))
	}
	if body[0].ID != 1 || body[0].Name != "Test University" {
		t.Errorf("unexpected university: %+v", body[0])
	}
}

func TestV3_GetUniversity_KnownID_Returns200(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	// Huma serializes single-object Body fields as a flat JSON object
	var body universityResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.ID != 1 {
		t.Errorf("ID: got %d, want 1", body.ID)
	}
	if body.Name != "Test University" {
		t.Errorf("Name: got %q", body.Name)
	}
}

func TestV3_GetUniversity_UnknownID_Returns404(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/999")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", resp.StatusCode)
	}
}

// --- v3 group endpoints ---

func TestV3_ListGroups_KnownUniv_Returns200(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/groups")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body []groupResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body) != 1 || body[0].ID != 1 {
		t.Errorf("expected 1 group with ID=1, got %v", body)
	}
}

func TestV3_ListGroups_UnknownUniv_Returns404(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/999/groups")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", resp.StatusCode)
	}
}

// --- v3 timetable endpoints ---

func TestV3_ListTimetables_KnownUnivAndGroup_Returns200(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/groups/1")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body []timetableResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body) != 1 || body[0].AdeResources != 20 {
		t.Errorf("expected 1 timetable with AdeResources=20, got %v", body)
	}
}

func TestV3_ListTimetables_UnknownGroup_Returns404(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/groups/999")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", resp.StatusCode)
	}
}

func TestV3_GetTimetableMetadata_KnownResources_Returns200(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/groups/1/20")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body timetableResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.AdeResources != 20 {
		t.Errorf("AdeResources: got %d, want 20", body.AdeResources)
	}
	if body.Label != "Timetable A" {
		t.Errorf("Label: got %q, want %q", body.Label, "Timetable A")
	}
}

func TestV3_GetTimetableMetadata_UnknownResources_Returns404(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/groups/1/999")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", resp.StatusCode)
	}
}

// TestV3_GetTimetableEvents_ServesCachedEvents verifies the events endpoint
// returns cached events (200) when the cache is warm.
func TestV3_GetTimetableEvents_ServesCachedEvents(t *testing.T) {
	// Pre-warm cache with known events for univID=1, adeResources=20
	events := []ade.Event{
		{Title: "Integration Maths", Location: "Amphi B"},
	}
	ade.SetTimetableByAdeResources(1, 20, events)

	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/groups/1/20/events")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body []ade.Event
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body) != 1 || body[0].Title != "Integration Maths" {
		t.Errorf("expected cached event, got %v", body)
	}
}

// --- v3 room endpoints ---

func TestV3_ListRooms_KnownUniv_Returns200(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/rooms")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body []roomResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body) != 1 || body[0].AdeResources != 10 {
		t.Errorf("expected 1 room with AdeResources=10, got %v", body)
	}
}

func TestV3_GetRoomMetadata_KnownResources_Returns200(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/rooms/10")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body roomResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Label != "Room A" {
		t.Errorf("Label: got %q, want %q", body.Label, "Room A")
	}
}

func TestV3_GetRoomMetadata_UnknownResources_Returns404(t *testing.T) {
	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/rooms/999")

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", resp.StatusCode)
	}
}

// TestV3_GetRoomEvents_ServesCachedEvents verifies room events endpoint
// returns cached events (200) when the cache is warm.
func TestV3_GetRoomEvents_ServesCachedEvents(t *testing.T) {
	events := []ade.Event{
		{Title: "Room Integration Test", Location: "Room A"},
	}
	ade.SetTimetableByAdeResources(1, 10, events)

	srv := newTestServer(t, testConfig())
	defer srv.Close()

	resp := get(t, srv, "/v3/univs/1/rooms/10/events")

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", resp.StatusCode)
	}

	var body []ade.Event
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body) < 1 || body[0].Title != "Room Integration Test" {
		t.Errorf("expected cached room event, got %v", body)
	}
}
